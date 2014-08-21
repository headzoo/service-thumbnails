package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"log"
	"os"
	"strconv"

	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"github.com/rakyll/magicmime"
)

// Misc constants.
const (
	MAX_MEMORY           = 1 * 1024 * 1024
	DEFAULT_MIME_TYPE    = "binary/octet-stream"
	DEFAULT_WIDTH_BIG    = 0
	DEFAULT_WIDTH_SPRITE = 180
)

// Upload stores the values of an uploaded file.
type Upload struct {
	Name     string
	Size     int64
	MimeType string
	Temp     string
}

func main() {
	http.HandleFunc("/thumbnailer/big", createBigThumbnail)
	http.HandleFunc("/thumbnailer/sprite", createSpriteThumbnail)
	http.ListenAndServe(":3000", nil)
}

// createBigThumbnail is the http callback to create a big thumbnail.
func createBigThumbnail(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	temp := getTempFile()
	width := DEFAULT_WIDTH_BIG
	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}

	err := ffmpeg.NewFFmpeg(file.Temp).CreateThumbnail(width, temp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	writeFileToResponse(temp, w)
}

// createSpriteThumbnail is the http callback to create a sprite thumbnail.
func createSpriteThumbnail(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	temp := getTempFile()
	width := DEFAULT_WIDTH_SPRITE
	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}

	ff := ffmpeg.NewFFmpeg(file.Temp)
	interval := int(ff.Length())
	if interval > thumbnailer.NUM_THUMBNAILS {
		interval = interval / thumbnailer.NUM_THUMBNAILS
	}

	err := ff.CreateThumbnailSprite(interval, width, temp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	writeFileToResponse(temp, w)
}

// getFile returns the uploaded file.
func getFile(w http.ResponseWriter, r *http.Request) *Upload {
	files, _ := writeUploadedFiles(r)
	if len(files) > 1 {
		w.WriteHeader(400)
		w.Write([]byte("Only a single file allowed."))
		return nil
	} else if len(files) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No files uploaded."))
		return nil
	}

	log.Printf("Got upload %#v\n", files[0])
	return &files[0]
}

// writeUploadedFiles writes all uploaded files to the temp dir.
func writeUploadedFiles(r *http.Request) ([]Upload, error) {
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		return nil, err
	}

	files := []Upload{}
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			fin, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}

			fout, err := ioutil.TempFile("/tmp", "upload")
			if err != nil {
				return nil, err
			}
			defer fout.Close()

			size, err := io.Copy(fout, fin)
			if err != nil {
				return nil, err
			}
			fout.Close()

			files = append(files, Upload{
				Name:     fileHeader.Filename,
				Size:     size,
				MimeType: getMimeType(fout.Name()),
				Temp:     fout.Name(),
			})
		}
	}

	return files, nil
}

// writeFileToResponse writes a file to the http response.
func writeFileToResponse(file string, w http.ResponseWriter) error {
	fout, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fout.Close()
	io.Copy(w, fout)

	return nil
}

// getMimeType returns the file mime type.
func getMimeType(file string) string {
	mm, err := magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
	if err != nil {
		return DEFAULT_MIME_TYPE
	}

	mimetype, err := mm.TypeByFile(file)
	if err != nil {
		return DEFAULT_MIME_TYPE
	}

	return mimetype
}

// getTempFile returns the name of a new temp file.
func getTempFile() string {
	temp, _ := ioutil.TempFile("/tmp", "thumb")
	temp.Close()
	return temp.Name()
}

// atoi converts a string to an integer.
func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		i = 0
	}

	return i
}
