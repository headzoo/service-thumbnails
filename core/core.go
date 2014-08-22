package core

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"
)

const (
	// AppVersion is the thumbnailer application version.
	AppVersion = "0.1"
	// ThumbCountPerSprite is the default number of thumbs to include in each sprite.
	ThumbCountPerSprite = 30
)

// Default values for command line options.
const (
	OptDefaultMode        = "cli"
	OptDefaultHost        = "127.0.0.1"
	OptDefaultPort        = 8080
	OptDefaultThumbType   = "simple"
	OptDefaultInFile      = ""
	OptDefaultOutFile     = ""
	OptDefaultWidth       = 0
	OptDefaultSkipSeconds = 0
	OptDefaultCount       = ThumbCountPerSprite
	OptDefaultQuiet       = false
	OptDefaultPrintHelp   = false
)

// Options stores the command line options.
type Options struct {
	Mode        string
	Host        string
	Port        int
	ThumbType   string
	InFile      string
	OutFile     string
	Width       int
	SkipSeconds int
	Count       int
	Quiet       bool
	PrintHelp   bool
}

// opts stores the command line options.
var Opts = &Options{}

// VPrintf prints the given message when verbose output is turned on.
func VPrintf(msg string, a ...interface{}) {
	if !Opts.Quiet {
		fmt.Printf(msg+"\n", a...)
	}
}

// VPrintfError prints the given message to stderr when verbose output is turned on.
func VPrintfError(msg string, a ...interface{}) {
	if !Opts.Quiet {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
	}
}

// Init initializes global application variables.
func Init() {
	flag.StringVar(
		&Opts.Mode,
		"m",
		OptDefaultMode,
		"Running mode, either 'cli' or 'http'. Defaults to 'cli'.")
	flag.BoolVar(
		&Opts.PrintHelp,
		"help",
		OptDefaultPrintHelp,
		"Display command help.")
	flag.BoolVar(
		&Opts.Quiet,
		"q",
		OptDefaultQuiet,
		"Run in quiet mode.")
	flag.IntVar(
		&Opts.SkipSeconds,
		"s",
		OptDefaultSkipSeconds,
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&Opts.Count,
		"c",
		OptDefaultCount,
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.IntVar(
		&Opts.Width,
		"w",
		OptDefaultWidth,
		"The thumbnail width. Overrides the built in defaults.")
	flag.StringVar(
		&Opts.ThumbType,
		"t",
		OptDefaultThumbType,
		"The type of thumbnail to generate. 'simple' is the default.")
	flag.StringVar(
		&Opts.InFile,
		"i",
		OptDefaultInFile,
		"The input video file. Separate multiple files with a comma.")
	flag.StringVar(
		&Opts.OutFile,
		"o",
		OptDefaultOutFile,
		"The output image file.")
	flag.StringVar(
		&Opts.Host,
		"h",
		OptDefaultHost,
		"The host name to listen on.")
	flag.IntVar(
		&Opts.Port,
		"p",
		OptDefaultPort,
		"The port to listen on.")
	flag.Parse()

	if Opts.PrintHelp {
		ExecuteHelpTemplate()
	} else if Opts.Mode == "cli" {
		if Opts.InFile == "" || Opts.OutFile == "" || Opts.ThumbType == "" {
			VPrintfError("Missing -i, -o, or -t.")
			ExecuteHelpTemplate()
		}
		if Opts.ThumbType != "sprite" && Opts.ThumbType != "simple" {
			VPrintfError("Invalid thumbnail type.")
			ExecuteHelpTemplate()
		}
	} else if Opts.Mode != "http" {
		VPrintfError("Invalid mode.")
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
		AppVersion,
		buff.String(),
	}

	t, _ := template.New("help").Parse(helpTemplate)
	t.Execute(os.Stdout, data)

	os.Exit(1)
}



// helpTemplate is the template used for displaying command line help.
const helpTemplate = `Thumbnailer v{{.Version}} - Used to generate thumbnails from videos.

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

SERVER USAGE:
	thumbnailer -m http -h <host> -p <port>

SERVER EXAMPLES:
	thumbnailer -m http -h 127.0.0.1 -p 3366
`
