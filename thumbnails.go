package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/dulo-tech/service-thumbnails/cli"
	"github.com/dulo-tech/service-thumbnails/core"
	"github.com/dulo-tech/service-thumbnails/http"
)

func main() {
	parseFlags()
	switch core.Opts.Mode {
	case "cli":
		cli.Go()
	case "http":
		http.Go()
	}
}

// parseFlags sets up the command line flag parser and parses the options.
func parseFlags() {
	flag.StringVar(
		&core.Opts.Mode,
		"m",
		core.OptDefaultMode,
		"Running mode, either 'cli' or 'http'. Defaults to 'cli'.")
	flag.BoolVar(
		&core.Opts.PrintHelp,
		"help",
		core.OptDefaultPrintHelp,
		"Display command help.")
	flag.BoolVar(
		&core.Opts.Quiet,
		"q",
		core.OptDefaultQuiet,
		"Run in quiet mode.")
	flag.IntVar(
		&core.Opts.SkipSeconds,
		"s",
		core.OptDefaultSkipSeconds,
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&core.Opts.Count,
		"c",
		core.OptDefaultCount,
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.IntVar(
		&core.Opts.Width,
		"w",
		core.OptDefaultWidth,
		"The thumbnail width. Overrides the built in defaults.")
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
	flag.StringVar(
		&core.Opts.Host,
		"h",
		core.OptDefaultHost,
		"The host name to listen on.")
	flag.IntVar(
		&core.Opts.Port,
		"p",
		core.OptDefaultPort,
		"The port to listen on.")
	flag.Parse()

	if core.Opts.PrintHelp {
		ExecuteHelpTemplate()
	} else if core.Opts.Mode == "cli" {
		if core.Opts.InFile == "" || core.Opts.OutFile == "" || core.Opts.ThumbType == "" {
			core.VPrintfError("Missing -i, -o, or -t.")
			ExecuteHelpTemplate()
		}
		if core.Opts.ThumbType != "sprite" && core.Opts.ThumbType != "simple" {
			core.VPrintfError("Invalid thumbnail type.")
			ExecuteHelpTemplate()
		}
	} else if core.Opts.Mode != "http" {
		core.VPrintfError("Invalid mode.")
		ExecuteHelpTemplate()
	}
}

// ExecuteHelpTemplate() prints the command line help using the given template and exits.
func ExecuteHelpTemplate() {
	buff := bytes.Buffer{}
	flag.VisitAll(func(f *flag.Flag) {
		buff.WriteString(fmt.Sprintf("\t-%-8s%s\n", f.Name, f.Usage))
	})

	data := struct {
		Version string
		Flags   string
	}{
		core.AppVersion,
		buff.String(),
	}

	t, _ := template.New("help").Parse(helpTemplate)
	t.Execute(os.Stdout, data)

	os.Exit(1)
}

// helpTemplate is the template used for displaying command line help.
const helpTemplate = `Thumbnailer v{{.Version}} - Used to generate thumbnails from videos.

Thumbnailer can run as a command line app or as an http server. Use the -m
switch to change modes. See usage and examples of each below.

OPTIONS:

{{.Flags}}
CLI USAGE:
	thumbnailer -t <sprite|simple> -i <video> -o <image>

	<sprite|simple> determines the type of thumbnail being generated. Either
	a sprite or a simple thumbnail. Simple is the default when not specified.

	<video> is one or more source videos. Separate multiple videos with commas.

	<image> may contain the place holders {name} and {type} which correspond
	to the name of the source video (without file extension) and the type of
	of thumbnail. One of 'sprite' or 'simple'. The <image> may also contain
	the verb %d which will be replaced with the file number. See the fmt package
	for more information on verbs.

CLI EXAMPLES:

	thumbnailer -t sprite -i source.mp4 -o thumb.jpg
	thumbnailer -i source1.mp4,source2.mp4 -o out%02d.jpg
	thumbnailer -t sprite -i source.mp4 -o thumb{name}{type}.jpg

HTTP USAGE:
	thumbnailer -m http -h <host> -p <port>

HTTP EXAMPLES:
	thumbnailer -m http -h 127.0.0.1 -p 3366
`
