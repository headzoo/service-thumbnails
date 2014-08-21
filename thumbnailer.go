package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
)

// Signals a job is finished.
type ChannelFinished chan bool

// Signals an error.
type ChannelError chan error

var opts = thumbnailer.Options{}
var chanFinished = make(ChannelFinished)
var chanError = make(ChannelError)

func main() {
	parseFlags()
	if opts.PrintHelp || opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		printHelp()
	}
	if opts.ThumbType != "sprite" && opts.ThumbType != "simple" {
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
	for i, file := range inFiles {
		out := opts.OutFile
		if strings.Contains(out, "%") {
			out = fmt.Sprintf(out, i)
		}
		if opts.ThumbType == "sprite" {
			go createSpriteThumbnail(file, out)
		} else {
			go createSimpleThumbnail(file, out)
		}
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

// createSimpleThumbnail creates a simple thumbnail.
func createSimpleThumbnail(inFile, outFile string) {
	defer func() {
		chanFinished <- true
	}()

	f := ffmpeg.New(inFile)
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

	f := ffmpeg.New(inFile)
	f.SkipSeconds = opts.SkipSeconds

	len := int(f.Length())
	interval := 0
	if len < opts.Count {
		interval = len
	} else {
		interval = len / opts.Count
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

// parseFlags parses the command line option flags.
func parseFlags() {
	flag.BoolVar(
		&opts.PrintHelp,
		"help",
		thumbnailer.OPT_PRINT_HELP,
		"Display command help.")
	flag.BoolVar(
		&opts.Verbose,
		"v",
		thumbnailer.OPT_VERBOSE,
		"Verbose output.")
	flag.StringVar(
		&opts.ThumbType,
		"t",
		thumbnailer.OPT_THUMB_TYPE,
		"The type of thumbnail to generate. 'simple' is the default.")
	flag.StringVar(
		&opts.InFile,
		"i",
		thumbnailer.OPT_IN_FILE,
		"The input video file. Separate multiple files with a comma.")
	flag.StringVar(
		&opts.OutFile,
		"o",
		thumbnailer.OPT_OUT_FILE,
		"The output image file.")
	flag.IntVar(
		&opts.Width,
		"w", thumbnailer.OPT_WIDTH,
		"The thumbnail width. Overrides the built in defaults.")
	flag.IntVar(
		&opts.SkipSeconds,
		"s",
		thumbnailer.OPT_SKIP_SECONDS,
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&opts.Count,
		"c",
		thumbnailer.OPT_COUNT,
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.Parse()
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer v%s\n", thumbnailer.VERSION)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("thumbnailer -t <sprite|simple> -i <video> -o <image>")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("thumbnailer -v -t sprite -i source.mp4 -o thumb.jpg")
	fmt.Println("thumbnailer -v -i source1.mp4,source2.mp4 -o out%02d.jpg")

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
