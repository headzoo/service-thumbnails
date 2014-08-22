package main

import (
	"flag"
	"os"
	"strings"

	"github.com/dulo-tech/thumbnailer/commands"
	"github.com/dulo-tech/thumbnailer/core"
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
func parseFlags() *core.Options {
	opts := core.FlagOptions()
	flag.StringVar(
		&opts.ThumbType,
		"t",
		core.OptDefaultThumbType,
		"The type of thumbnail to generate. 'simple' is the default.")
	flag.StringVar(
		&opts.InFile,
		"i",
		core.OptDefaultInFile,
		"The input video file. Separate multiple files with a comma.")
	flag.StringVar(
		&opts.OutFile,
		"o",
		core.OptDefaultOutFile,
		"The output image file.")
	flag.Parse()

	if opts.PrintHelp {
		core.ExecuteHelpTemplate(opts, thumbnailerHelpTemplate)
	}
	if opts.InFile == "" || opts.OutFile == "" || opts.ThumbType == "" {
		core.ExecuteHelpTemplate(opts, thumbnailerHelpTemplate)
	}
	if opts.ThumbType != "sprite" && opts.ThumbType != "simple" {
		core.ExecuteHelpTemplate(opts, thumbnailerHelpTemplate)
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
			core.VPrintfError("The input file %q does not exist.", file)
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
