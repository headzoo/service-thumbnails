package commands

import (
	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
)

// SpriteCommand is used to generate sprite thumbnails from the command line.
type SpriteCommand struct {
	Command
}

// NewSprite creates and returns a new SpriteCommand instance.
func NewSprite(opts *thumbnailer.Options) *SpriteCommand {
	return &SpriteCommand{
		Command: *newCommand(opts),
	}
}

// Execute processes a command instruction.
func (c *SpriteCommand) Execute(inFile, outFile string) {
	defer func() {
		(*c.chanFinished) <- true
	}()

	f := ffmpeg.New(inFile)
	f.SkipSeconds = c.opts.SkipSeconds

	len := int(f.Length())
	interval := 0
	if len < c.opts.Count {
		interval = len
	} else {
		interval = len / c.opts.Count
	}

	width := 180
	if c.opts.Width != 0 {
		width = c.opts.Width
	}

	err := f.CreateThumbnailSprite(interval, width, outFile)
	if err != nil {
		(*c.chanError) <- err
		return
	}

	thumbnailer.Verbose("Sprite thumbnail for video %q written to %q.", inFile, outFile)
}