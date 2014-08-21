package handlers

import (
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"net/http"
	"github.com/dulo-tech/thumbnailer/ffmpeg"
)

// SpriteHandler is an HTTP handler for creating sprite thumbnails.
type SpriteHandler struct {
	opts *thumbnailer.Options
}

// NewSprite creates and returns a new SpriteHandler instance.
func NewSprite(opts *thumbnailer.Options) *SpriteHandler {
	return &SpriteHandler{
		opts: opts,
	}
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h *SpriteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	width := DEFAULT_WIDTH_SPRITE
	skip := h.opts.SkipSeconds
	count := h.opts.Count

	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}
	if s, ok := query["skip"]; ok {
		skip = atoi(s[0])
	}
	if s, ok := query["count"]; ok {
		count = atoi(s[0])
	}

	temp := getTempFile()
	ff := ffmpeg.New(file.Temp)
	ff.SkipSeconds = skip

	interval := int(ff.Length())
	if interval > count {
		interval = interval / count
	}

	err := ff.CreateThumbnailSprite(interval, width, temp)
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
