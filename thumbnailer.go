package main

import (
	"flag"
	"os"
	"strings"

	"github.com/dulo-tech/thumbnailer/commands"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
)

func main() {
	opts := parseFlags()

	router := commands.NewRouter(splitFiles(opts.InFile), opts.OutFile)
	router.Command("simple", commands.NewSimple(opts))
	router.Command("sprite", commands.NewSprite(opts))
	err := router.Route(opts.ThumbType)
	if err != nil {
		panic(err)
	}
}

// parseFlags parses the command line option flags.
func parseFlags() *thumbnailer.Options {
	var opts = &thumbnailer.Options{}

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

	thumbnailer.VerboseOutput = opts.Verbose
	if opts.PrintHelp {
		thumbnailer.PrintHelp(opts, thumbnailerHelpTemplate)
	}
	if opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		thumbnailer.PrintHelp(opts, thumbnailerHelpTemplate)
	}
	if opts.ThumbType != "sprite" && opts.ThumbType != "simple" {
		thumbnailer.PrintHelp(opts, thumbnailerHelpTemplate)
	}

	return opts
}

// splitFiles converts a comma separated list of files into an array of file names.
func splitFiles(inFiles string) []string {
	files := strings.Split(inFiles, ",")
	for i, f := range files {
		files[i] = strings.Trim(f, " ")
	}

	for _, file := range files {
		if !fileExists(file) {
			thumbnailer.VerboseError("The input file %q does not exist.", file)
			os.Exit(1)
		}
	}

	return files
}

// fileExists returns whether the given file exists.
func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

const thumbnailerHelpTemplate = `Thumbnailer v{{.Version}} - Used to generate thumbnails from videos.

USAGE:
	thumbnailer -t <sprite|simple> -i <video> -o <image>

	<sprite|simple> determines the type of thumbnail being generated. Either
	a sprite or a simple thumbnail. Simple is the default when not specified.
	
	<video> is one or more source videos. Separate multiple videos with commas.
	
	<image> may contain the place holders {name} and {type} which correspond
	to the name of the source video (without file extension) and the type of
	of thumbnail. One of 'sprite' or 'simple'. The <image> may also contain
	the verb %d which will be replaced with the file number. See the fmt package
	for more information on verbs.

OPTIONS:

{{.Flags}}
EXAMPLES:

	thumbnailer -t sprite -i source.mp4 -o thumb.jpg
	thumbnailer -i source1.mp4,source2.mp4 -o out%02d.jpg
	thumbnailer -t sprite -i source.mp4 -o thumb{name}{type}.jpg
`
