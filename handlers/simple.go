package handlers

import (
	"github.com/dulo-tech/thumbnailer/core"
	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"net/http"
)

// SimpleHandler is an HTTP handler for creating simple thumbnails.
type SimpleHandler struct {
	Handler
}

// NewPulse creates and returns a new SimpleHandler instance.
func NewSimple() *SimpleHandler {
	return &SimpleHandler{}
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h *SimpleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	width := DEFAULT_WIDTH_SIMPLE
	skip := core.Opts.SkipSeconds

	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}
	if s, ok := query["skip"]; ok {
		skip = atoi(s[0])
	}

	temp := getTempFile()
	ff := ffmpeg.New(file.Temp)
	ff.SkipSeconds = skip

	err := ff.CreateThumbnail(width, temp)
	if err != nil {
		numErrors++
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	numRequests++
	w.Header().Set("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	writeFileToResponse(temp, w)
}
