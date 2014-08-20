package main

import (
	 "github.com/dulo-tech/thumbnailer/ffmpeg"
	"fmt"
)

const NUM_THUMBNAILS = 30

func main() {
	srcFileName := "/home/sean/Downloads/big_buck_bunny.mp4"

	f := ffmpeg.NewFFmpeg(srcFileName)
	f.SkipSeconds = 10

	len := int(f.Length())
	interval := len / NUM_THUMBNAILS
	fmt.Printf("len = %d, thumbs = %d\n", len, interval)
	err := f.CreateThumbnailSprite(interval, 180, "./thumbnail.jpg")
	if err != nil {
		panic(err)
	}

	f.SkipSeconds = 0
	err = f.CreateThumbnail(0, "./thumbnail-big.jpg")
	if err != nil {
		panic(err)
	}
}
