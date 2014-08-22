package commands

import (
	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/core"
)

// SimpleCommand is used to generate simple thumbnails from the command line.
type SimpleCommand struct {
	Command
}

// NewSimple creates and returns a new SimpleCommand instance.
func NewSimple(opts *core.Options) *SimpleCommand {
	return &SimpleCommand{
		Command: *newCommand(opts),
	}
}

// Execute processes a command instruction.
func (c *SimpleCommand) Execute(inFile, outFile string) {
	defer func() {
		(*c.chanFinished) <- true
	}()

	f := ffmpeg.New(inFile)
	f.SkipSeconds = c.opts.SkipSeconds

	err := f.CreateThumbnail(c.opts.Width, outFile)
	if err != nil {
		(*c.chanError) <- err
		return
	}

	core.VPrintf("Simple thumbnail for video %q written to %q.", inFile, outFile)
}
