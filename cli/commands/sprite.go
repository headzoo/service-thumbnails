package commands

import (
	"github.com/dulo-tech/service-thumbnails/core"
	"github.com/dulo-tech/service-thumbnails/ffmpeg"
)

// SpriteCommand is used to generate sprite thumbnails from the command line.
type SpriteCommand struct {
	Command
}

// NewSprite creates and returns a new SpriteCommand instance.
func NewSprite() *SpriteCommand {
	return &SpriteCommand{
		Command: *newCommand(),
	}
}

// Execute processes a command instruction.
func (c *SpriteCommand) Execute(inFile, outFile string) {
	defer func() {
		(*c.chanFinished) <- true
	}()

	f := ffmpeg.New(inFile)
	f.SkipSeconds = core.Opts.SkipSeconds

	len := int(f.Length())
	interval := 0
	if len < core.Opts.Count {
		interval = len
	} else {
		interval = len / core.Opts.Count
	}

	width := 180
	if core.Opts.Width != 0 {
		width = core.Opts.Width
	}

	err := f.CreateThumbnailSprite(interval, width, outFile)
	if err != nil {
		(*c.chanError) <- err
		return
	}

	core.VPrintf("Sprite thumbnail for video %q written to %q.", inFile, outFile)
}
