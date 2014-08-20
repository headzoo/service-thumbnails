package main

import (
	"flag"
	"fmt"
	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"os"
)

// Misc constants.
const (
	VERSION = "0.1"
	NUM_THUMBNAILS = 30
)

// Default values for command line options.
const (
	OPT_THUMB_TYPE = ""
	OPT_IN_FILE = ""
	OPT_OUT_FILE = ""
	OPT_WIDTH = 0
	OPT_SKIP_SECONDS = 0
	OPT_TEMP_DIR = "/tmp"
	OPT_VERBOSE = false
	OPT_PRINT_HELP = false
)

// Signals for the app to exit.
type ChannelExit chan bool

// Signals an error.
type ChannelError chan error

// Options stores the command line options.
type Options struct {
	ThumbType string
	InFile    string
	OutFile   string
	Width     int
	SkipSeconds int
	TempDirectory string
	Verbose   bool
	PrintHelp bool
}

var opts = Options{}

func main() {
	flag.BoolVar(&opts.PrintHelp, "h", OPT_PRINT_HELP, "Display command help.")
	flag.BoolVar(&opts.Verbose, "v", OPT_VERBOSE, "Verbose output.")
	flag.StringVar(&opts.ThumbType, "t", OPT_THUMB_TYPE, "The type of thumbnail to generate.")
	flag.StringVar(&opts.InFile, "i", OPT_IN_FILE, "The input video file.")
	flag.StringVar(&opts.OutFile, "o", OPT_OUT_FILE, "The output image file.")
	flag.IntVar(&opts.Width, "w", OPT_WIDTH, "The thumbnail width. Overrides the built in defaults.")
	flag.IntVar(&opts.SkipSeconds, "s", OPT_SKIP_SECONDS, "Skip this number of seconds into the video before thubmnailing.")
	flag.StringVar(&opts.TempDirectory, "d", OPT_TEMP_DIR, "Temp directory.")
	flag.Parse()

	if opts.PrintHelp || opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		printHelp()
	}
	if opts.ThumbType != "sprite" && opts.ThumbType != "big" {
		printHelp()
	}
	if _, err := os.Stat(opts.InFile); os.IsNotExist(err) {
		fmt.Println("The input video file does not exist.")
		os.Exit(1)
	}

	chanExit := make(ChannelExit)
	chanError := make(ChannelError)

	if opts.ThumbType == "big" {
		go createBigThumbnail(opts.InFile, chanExit, chanError)
	} else if opts.ThumbType == "sprite" {
		go createSpriteThumbnail(opts.InFile, chanExit, chanError)
	} else {
		printHelp()
	}

	for {
		select {
		case err := <-chanError:
			verbose(err.Error())
			os.Exit(1)
		case <-chanExit:
			verbose(fmt.Sprintf("Thumbnail written to file %q.", opts.OutFile))
			os.Exit(0)
		}
	}
}

// createBigThumbnail creates a big thumbnail.
func createBigThumbnail(inFile string, chanExit ChannelExit, chanError ChannelError) {
	verbose("Creating 'big' thumbnail...")
	defer func() {
		close(chanExit)
	}()

	f := ffmpeg.NewFFmpeg(inFile)
	f.TmpDirectory = opts.TempDirectory
	f.SkipSeconds = opts.SkipSeconds

	err := f.CreateThumbnail(opts.Width, opts.OutFile)
	if err != nil {
		chanError <- err
	}
}

// createSpriteThumbnail creates a sprite thumbnail.
func createSpriteThumbnail(inFile string, chanExit ChannelExit, chanError ChannelError) {
	verbose("Creating 'strip' thumbnail...")
	defer func() {
		close(chanExit)
	}()

	f := ffmpeg.NewFFmpeg(inFile)
	f.TmpDirectory = opts.TempDirectory
	f.SkipSeconds = opts.SkipSeconds

	len := int(f.Length())
	interval := len / NUM_THUMBNAILS
	width := 180
	if opts.Width != 0 {
		width = opts.Width
	}

	err := f.CreateThumbnailSprite(interval, width, opts.OutFile)
	if err != nil {
		chanError <- err
	}
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer v%s\n", VERSION)
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

	os.Exit(1)
}

// verbose prints the given message when verbose output is turned on.
func verbose(msg string) {
	if opts.Verbose {
		fmt.Println(msg)
	}
}
