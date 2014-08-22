package core

import (
	"fmt"
	"os"
)

const (
	// AppVersion is the thumbnailer application version.
	AppVersion = "0.2"
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
