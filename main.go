package main

import (
	 "github.com/dulo-tech/thumbnailer/ffmpeg"
)

// Default number of thumbs to create for each video.
const NUM_THUMBNAILS = 30

func main() {
	srcFileName := "./big_buck_bunny.mp4"

	f := ffmpeg.NewFFmpeg(srcFileName)
	f.TmpDirectory = "/tmp"

	// Single frame thumbnail.
	f.SkipSeconds = 0
	err := f.CreateThumbnail(0, "./thumbnail-big.jpg")
	if err != nil {
		panic(err)
	}

	// Thumbnail sprite.
	f.SkipSeconds = 10
	len := int(f.Length())
	interval := len / NUM_THUMBNAILS
	err = f.CreateThumbnailSprite(interval, 180, "./thumbnail-sprite.jpg")
	if err != nil {
		panic(err)
	}
}
