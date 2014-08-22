package thumbnailer

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
	OptDefaultHost        = "127.0.0.1"
	OptDefaultPort        = 8080
	OptDefaultThumbType   = "simple"
	OptDefaultInFile      = ""
	OptDefaultOutFile     = ""
	OptDefaultWidth       = 0
	OptDefaultSkipSeconds = 0
	OptDefaultCount       = ThumbCountPerSprite
	OptDefaultVerbose     = false
	OptDefaultPrintHelp   = false
)

// Options stores the command line options.
type Options struct {
	Host        string
	Port        int
	ThumbType   string
	InFile      string
	OutFile     string
	Width       int
	SkipSeconds int
	Count       int
	Verbose     bool
	PrintHelp   bool
}

// opts stores the command line options.
var opts = &Options{}

// verbose prints the given message when verbose output is turned on.
func VPrintf(msg string, a ...interface{}) {
	if opts.Verbose {
		fmt.Printf(msg+"\n", a...)
	}
}

// verboseError prints the given message to stderr when verbose output is turned on.
func VPrintfError(msg string, a ...interface{}) {
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
	}
}

// FlagOptions initializes the command flags for both the cli app and server.
func FlagOptions() *Options {
	flag.BoolVar(
		&opts.PrintHelp,
		"help",
		OptDefaultPrintHelp,
		"Display command help.")
	flag.BoolVar(
		&opts.Verbose,
		"v",
		OptDefaultVerbose,
		"Verbose output.")
	flag.IntVar(
		&opts.SkipSeconds,
		"s",
		OptDefaultSkipSeconds,
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&opts.Count,
		"c",
		OptDefaultCount,
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.IntVar(
		&opts.Width,
		"w",
		OptDefaultWidth,
		"The thumbnail width. Overrides the built in defaults.")

	return opts
}

// ExecuteHelpTemplate() prints the command line help using the given template and exits.
func ExecuteHelpTemplate(opts *Options, t string) {
	if opts.Verbose || opts.PrintHelp {
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

		t, _ := template.New("help").Parse(t)
		t.Execute(os.Stdout, data)
	}

	os.Exit(1)
}
