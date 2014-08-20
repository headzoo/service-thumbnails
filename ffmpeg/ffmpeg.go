package ffmpeg

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"fmt"
	"io/ioutil"
)

// FFmpeg is used to create thumbnails from videos.
type FFmpeg struct {
	SkipSeconds int
	Video string
	TmpDirectory string

	cmdFFprobe string
	cmdFFmpeg string
	cmdConvert string
}

// Creates and returns a new FFmpeg instance.
func NewFFmpeg(video string) *FFmpeg {
	//ffmpeg_path := os.Getenv("HW_FFMPEG_PATH")
	//if ffmpeg_path == "" {
	//	panic("Environment variables HW_FFMPEG_PATH not set.")
	//}

	return &FFmpeg{
		SkipSeconds: 0,
		Video: video,
		TmpDirectory: "/tmp",

		cmdFFprobe: "/home/sean/lib/FFmpeg/ffprobe",
		cmdFFmpeg: "/home/sean/lib/FFmpeg/ffmpeg",
		cmdConvert: "/usr/bin/convert",
	}
}

// Length returns the length of the video in seconds.
func (f *FFmpeg) Length() float64 {
	output, err := exec.Command(
		f.cmdFFprobe,
		"-i",
		f.Video,
		"-v",
		"quiet",
		"-show_entries",
		"format=duration",
		"-of",
		"csv=p=0",
	).Output()
	if err != nil {
		return 0.0
	}

	len, err := strconv.ParseFloat(strings.Trim(string(output), "\n"), 64)
	if err != nil {
		return 0.0
	}

	return len
}

// CreateThumbnail creates a single thumbnail from the video.
// The frame is determined by the value of FFmpeg.SkipSeconds.
// When 0 is given for the 'width' argument, the thumbnail will have the same
// width of the video.
func (f *FFmpeg) CreateThumbnail(width int, outFile string) error {
	os.Remove(outFile)

	args := []string{
		"-ss",
		SecondsToTime(f.SkipSeconds),
		"-i",
		f.Video,
		"-f",
		"image2",
		"-vframes",
		"1",
	}
	if width != 0 {
		args = append(args, "-vf")
		args = append(args, fmt.Sprintf("scale='min(%d\\,iw)':-1", width))
	}
	args = append(args, outFile)

	err := exec.Command(f.cmdFFmpeg, args...).Run()
	if err != nil {
		return err
	}

	return nil
}

// CreateThumbnailSprite creates thumbnails from the video at the given interval,
// and stitches them together into a single sprite.
// A thumbnail is generated every 'interval' seconds with a max width of 'width'.
// The thumbnails are then stitched together into a single image written to 'outFile'.
func (f *FFmpeg) CreateThumbnailSprite(interval, width int, outFile string) error {
	tmp, err := ioutil.TempDir(f.TmpDirectory, "thumb")
	if err != nil {
		return err
	}
	err = os.MkdirAll(tmp, os.FileMode(0755))
	if err != nil {
		return err
	}
	os.Remove(outFile)

	filters := []string{
		fmt.Sprintf("fps=fps=1/%d", interval),
		fmt.Sprintf("scale='min(%d\\,iw)':-1", width),
	}

	err = exec.Command(
		f.cmdFFmpeg,
		"-i",
		f.Video,
		"-ss",
		SecondsToTime(f.SkipSeconds),
		"-f",
		"image2",
		"-vf",
		strings.Join(filters, ","),
		tmp + "/frames%04d.jpg",
	).Run()
	if err != nil {
		return err
	}

	err = exec.Command(
		f.cmdConvert,
		tmp + "/*.jpg",
		"+append",
		outFile,
	).Run()
	if err != nil {
		return err
	}

	return nil
}

// SecondsToTime converts seconds into "00:00:00" format.
func SecondsToTime(secs int) string {
	if secs == 0 {
		return "00:00:00"
	}

	hours  := secs / 3600;
	minutes := (secs - (hours * 3600)) / 60;
	seconds := secs - (hours * 3600) -  (minutes * 60);

	return fmt.Sprintf("%.2d:%.2d:%.2d", hours, minutes, seconds)
}
