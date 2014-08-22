package main

import (
	"flag"
	"os"
	"strings"

	"github.com/dulo-tech/thumbnailer/commands"
	"github.com/dulo-tech/thumbnailer/core"
)

func main() {
	Init()
	router := commands.NewRouter(splitFiles(core.Opts.InFile), core.Opts.OutFile)
	router.Command("simple", commands.NewSimple())
	router.Command("sprite", commands.NewSprite())
	err := router.Route(core.Opts.ThumbType)
	if err != nil {
		panic(err)
	}
}

// Init parses the command line option flags.
func Init() {
	core.Init()
	flag.StringVar(
		&core.Opts.ThumbType,
		"t",
		core.OptDefaultThumbType,
		"The type of thumbnail to generate. 'simple' is the default.")
	flag.StringVar(
		&core.Opts.InFile,
		"i",
		core.OptDefaultInFile,
		"The input video file. Separate multiple files with a comma.")
	flag.StringVar(
		&core.Opts.OutFile,
		"o",
		core.OptDefaultOutFile,
		"The output image file.")
	flag.Parse()
	
	if core.Opts.PrintHelp {
		core.ExecuteHelpTemplate(core.Opts, thumbnailerHelpTemplate)
	}
	if core.Opts.InFile == "" || core.Opts.OutFile == "" || core.Opts.ThumbType == "" {
		core.ExecuteHelpTemplate(core.Opts, thumbnailerHelpTemplate)
	}
	if core.Opts.ThumbType != "sprite" && core.Opts.ThumbType != "simple" {
		core.ExecuteHelpTemplate(core.Opts, thumbnailerHelpTemplate)
	}
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
