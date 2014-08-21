package main

import (
	"flag"
	"fmt"

	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"os"
	"strings"
)

// Default values for command line options.
const (
	OPT_THUMB_TYPE   = ""
	OPT_IN_FILE      = ""
	OPT_OUT_FILE     = ""
	OPT_WIDTH        = 0
	OPT_SKIP_SECONDS = 0
	OPT_TEMP_DIR     = "/tmp"
	OPT_VERBOSE      = false
	OPT_PRINT_HELP   = false
)

// Signals a job is finished.
type ChannelFinished chan bool

// Signals an error.
type ChannelError chan error

// Options stores the command line options.
type Options struct {
	ThumbType     string
	InFile        string
	OutFile       string
	Width         int
	SkipSeconds   int
	TempDirectory string
	Verbose       bool
	PrintHelp     bool
}

var opts = Options{}
var chanFinished = make(ChannelFinished)
var chanError = make(ChannelError)

func main() {
	flag.BoolVar(&opts.PrintHelp, "help", OPT_PRINT_HELP, "Display command help.")
	flag.BoolVar(&opts.Verbose, "v", OPT_VERBOSE, "Verbose output.")
	flag.StringVar(&opts.ThumbType, "t", OPT_THUMB_TYPE, "The type of thumbnail to generate.")
	flag.StringVar(&opts.InFile, "i", OPT_IN_FILE, "The input video file. Separate multiple files with a comma.")
	flag.StringVar(&opts.OutFile, "o", OPT_OUT_FILE, "The output image file.")
	flag.IntVar(&opts.Width, "w", OPT_WIDTH, "The thumbnail width. Overrides the built in defaults.")
	flag.IntVar(&opts.SkipSeconds, "s", OPT_SKIP_SECONDS, "Skip this number of seconds into the video before thumbnailing.")
	flag.StringVar(&opts.TempDirectory, "d", OPT_TEMP_DIR, "Temp directory.")
	flag.Parse()

	if opts.PrintHelp || opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		printHelp()
	}
	if opts.ThumbType != "sprite" && opts.ThumbType != "big" {
		printHelp()
	}

	inFiles := strings.Split(opts.InFile, ",")
	for i, f := range inFiles {
		inFiles[i] = strings.Trim(f, " ")
	}
	for _, file := range inFiles {
		if !fileExists(file) {
			verboseError(fmt.Sprintf("The input file %q does not exist.", file))
			os.Exit(1)
		}
	}

	verbose(fmt.Sprintf("Thumbnailing %d video(s).", len(inFiles)))
	if opts.ThumbType == "big" {
		for i, file := range inFiles {
			go createBigThumbnail(file, fmt.Sprintf(opts.OutFile, i))
		}
	} else if opts.ThumbType == "sprite" {
		for i, file := range inFiles {
			go createSpriteThumbnail(file, fmt.Sprintf(opts.OutFile, i))
		}
	} else {
		printHelp()
	}

	finished := 0
	numJobs := len(inFiles)
	for {
		select {
		case err := <-chanError:
			{
				verbose(err.Error())
				os.Exit(1)
			}
		case <-chanFinished:
			{
				finished++
				if finished == numJobs {
					verbose("Finished")
					os.Exit(0)
				}
			}
		}
	}
}

// createBigThumbnail creates a big thumbnail.
func createBigThumbnail(inFile, outFile string) {
	defer func() {
		chanFinished <- true
	}()

	f := ffmpeg.NewFFmpeg(inFile)
	f.TmpDirectory = opts.TempDirectory
	f.SkipSeconds = opts.SkipSeconds

	err := f.CreateThumbnail(opts.Width, outFile)
	if err != nil {
		chanError <- err
		return
	}
	verbose(fmt.Sprintf("Thumbnail written to file %q.", outFile))
}

// createSpriteThumbnail creates a sprite thumbnail.
func createSpriteThumbnail(inFile, outFile string) {
	defer func() {
		chanFinished <- true
	}()

	f := ffmpeg.NewFFmpeg(inFile)
	f.TmpDirectory = opts.TempDirectory
	f.SkipSeconds = opts.SkipSeconds

	len := int(f.Length())
	interval := 0
	if len < thumbnailer.NUM_THUMBNAILS {
		interval = len
	} else {
		interval = len / thumbnailer.NUM_THUMBNAILS
	}

	width := 180
	if opts.Width != 0 {
		width = opts.Width
	}

	err := f.CreateThumbnailSprite(interval, width, outFile)
	if err != nil {
		chanError <- err
		return
	}
	verbose(fmt.Sprintf("Thumbnail written to file %q.", outFile))
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer v%s\n", thumbnailer.VERSION)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("thumbnailer -t <sprite|big> -i <video> -o <image>")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("thumbnailer -v -t strip -i source.mp4 -o thumb.jpg")
	fmt.Println("thumbnailer -v -t big -i source1.mp4,source2.mp4 -o out%02d.jpg")

	os.Exit(1)
}

// verbose prints the given message when verbose output is turned on.
func verbose(msg string) {
	if opts.Verbose {
		fmt.Println(msg)
	}
}

// verbose prints the given message when verbose output is turned on.
func verboseError(msg string) {
	if opts.Verbose {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// fileExists returns whether the given file exists.
func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
