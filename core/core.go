package core

import (
	"fmt"
	"os"
)

const (
	// AppName is the full name of the application.
	AppName = "server-thumbnails"
	// ThumbCountPerSprite is the default number of thumbs to include in each sprite.
	ThumbCountPerSprite = 30
)

// Default values for command line options.
const (
	OptDefaultMode         = "cli"
	OptDefaultHost         = "127.0.0.1"
	OptDefaultPort         = 8080
	OptDefaultThumbType    = "simple"
	OptDefaultInFile       = ""
	OptDefaultOutFile      = ""
	OptDefaultWidth        = 0
	OptDefaultSkipSeconds  = 0
	OptDefaultCount        = ThumbCountPerSprite
	OptDefaultQuiet        = false
	OptDefaultPrintHelp    = false
	OptDefaultPrintVersion = false
)

// ThumbTypes stores the possible thumbnail types that may be generated.
var ValidThumbTypes = []string{"sprite", "simple"}

// Options stores the command line options.
type Options struct {
	Mode         string
	Host         string
	Port         int
	ThumbType    string
	InFile       string
	OutFile      string
	Width        int
	SkipSeconds  int
	Count        int
	Quiet        bool
	PrintHelp    bool
	PrintVersion bool
}

// opts stores the command line options.
var Opts = &Options{
	Mode:         OptDefaultMode,
	Host:         OptDefaultHost,
	Port:         OptDefaultPort,
	ThumbType:    OptDefaultThumbType,
	InFile:       OptDefaultInFile,
	OutFile:      OptDefaultOutFile,
	Width:        OptDefaultWidth,
	SkipSeconds:  OptDefaultSkipSeconds,
	Count:        OptDefaultCount,
	Quiet:        OptDefaultQuiet,
	PrintHelp:    OptDefaultPrintHelp,
	PrintVersion: OptDefaultPrintVersion,
}

// BuildInfo returns a string with the build information.
func BuildInfo() string {
	return fmt.Sprintf(
		"%s %s.b%s %s %s",
		AppName,
		AppVersion,
		AppBuildNumber,
		AppBuildDate,
		AppBuildArch)
}

// VPrintf prints the given message when verbose output is turned on.
func VPrintf(msg string, a ...interface{}) {
	if !Opts.Quiet {
		fmt.Printf(msg+"\n", a...)
	}
}

// VPrintfError prints the given message to stderr when verbose output is turned on.
func VErrorf(msg string, a ...interface{}) {
	if !Opts.Quiet {
		fmt.Fprintf(os.Stderr, msg+"\n", a...)
	}
}

// FileExists returns whether the given file exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
