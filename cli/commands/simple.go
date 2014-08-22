package commands

import (
	"github.com/dulo-tech/service-thumbnails/core"
	"github.com/dulo-tech/service-thumbnails/ffmpeg"
)

// SimpleCommand is used to generate simple thumbnails from the command line.
type SimpleCommand struct {
	Command
}

// NewSimple creates and returns a new SimpleCommand instance.
func NewSimple() *SimpleCommand {
	return &SimpleCommand{
		Command: *newCommand(),
	}
}

// Execute processes a command instruction.
func (c *SimpleCommand) Execute(inFile, outFile string) {
	defer func() {
		(*c.chanFinished) <- true
	}()

	f := ffmpeg.New(inFile)
	f.SkipSeconds = core.Opts.SkipSeconds

	err := f.CreateThumbnail(core.Opts.Width, outFile)
	if err != nil {
		(*c.chanError) <- err
		return
	}

	core.VPrintf("Simple thumbnail for video %q written to %q.", inFile, outFile)
}
