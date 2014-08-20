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
	OPT_VERBOSE = false
	OPT_PRINT_HELP = false
)

// Options stores the command line options.
type Options struct {
	ThumbType string
	InFile    string
	OutFile   string
	Width     int
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
	flag.Parse()

	if opts.PrintHelp || opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		printHelp()
	}
	if opts.ThumbType != "strip" && opts.ThumbType != "big" {
		printHelp()
	}
	if _, err := os.Stat(opts.InFile); os.IsNotExist(err) {
		fmt.Println("The input video file does not exist.")
		os.Exit(1)
	}

	f := ffmpeg.NewFFmpeg(opts.InFile)
	f.TmpDirectory = "/tmp"

	if opts.ThumbType == "big" {
		verbose("Creating 'big' thumbnail...")
		f.SkipSeconds = 0
		width := opts.Width
		err := f.CreateThumbnail(width, opts.OutFile)
		if err != nil {
			panic(err)
		}
		verbose("Finished!")
	} else if opts.ThumbType == "strip" {
		verbose("Creating 'strip' thumbnail...")
		f.SkipSeconds = 10

		len := int(f.Length())
		interval := len / NUM_THUMBNAILS
		width := 180
		if opts.Width != 0 {
			width = opts.Width
		}

		err := f.CreateThumbnailSprite(interval, width, opts.OutFile)
		if err != nil {
			panic(err)
		}
		verbose("Finished!")
	} else {
		printHelp()
	}
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer v%s\n", VERSION)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("thumbnailer -t <strip|big> -i <video> -o <image>")
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
