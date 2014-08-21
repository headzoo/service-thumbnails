package thumbnailer

import (
	"fmt"
	"os"
)

const (
	// The main application version.
	VERSION = "0.1"

	// The number of thumbs to include in a sprite.
	NUM_THUMBNAILS = 30
)

// Default values for command line options.
const (
	OPT_HOST         = "127.0.0.1"
	OPT_PORT         = 8080
	OPT_THUMB_TYPE   = "simple"
	OPT_IN_FILE      = ""
	OPT_OUT_FILE     = ""
	OPT_WIDTH        = 0
	OPT_SKIP_SECONDS = 0
	OPT_COUNT        = NUM_THUMBNAILS
	OPT_VERBOSE      = false
	OPT_PRINT_HELP   = false
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

// Verbose stores whether to use verbose output or not.
var VerboseOutput bool

// verbose prints the given message when verbose output is turned on.
func Verbose(msg string, a ...interface{}) {
	if VerboseOutput {
		fmt.Printf(msg+"\n", a...)
	}
}

// verboseError prints the given message to stderr when verbose output is turned on.
func VerboseError(msg string, a ...interface{}) {
	if VerboseOutput {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
	}
}
