package handlers

import (
	"net/http"
	"html/template"
	"os"
	
	"github.com/dulo-tech/thumbnailer/thumbnailer"
)

// HelpData stores template variables for the help page.
type HelpData struct {
	DefaultCount int
	DefaultSkip  int
}

// HelpHandler is an HTTP handler for displaying a help page using HTML.
type HelpHandler struct {
	opts *thumbnailer.Options
}

// NewHelp creates and returns a new HelpHandler instance.
func NewHelp(opts *thumbnailer.Options) *HelpHandler {
	return &HelpHandler{
		opts: opts,
	}
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h *HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := HelpData{
		DefaultCount: h.opts.Count,
		DefaultSkip:  h.opts.SkipSeconds,
	}
	
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(dir + "/handlers/templates/help.html")
	if err != nil {
		//numErrors++
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	//numRequests++
	t.Execute(w, data)
}
