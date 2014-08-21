package handlers

import (
	"html/template"
	"net/http"

	"github.com/dulo-tech/thumbnailer/thumbnailer"
)

// HelpData stores template variables for the help page.
type HelpData struct {
	DefaultCount int
	DefaultSkip  int
}

// HelpHandler is an HTTP handler for displaying a help page using HTML.
type HelpHandler struct {
	Handler
}

// NewHelp creates and returns a new HelpHandler instance.
func NewHelp(opts *thumbnailer.Options) *HelpHandler {
	return &HelpHandler{
		Handler: *New(opts),
	}
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h *HelpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := HelpData{
		DefaultCount: h.opts.Count,
		DefaultSkip:  h.opts.SkipSeconds,
	}

	t, err := template.New("help").Parse(helpTemplate)
	if err != nil {
		numErrors++
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	numRequests++
	t.Execute(w, data)
}

const helpTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <title>Help</title>
    </head>
    <body>
        <h1>Thumbnailer Help</h1>
        <p>End Points:</p>
        <ul>
            <li>
                POST <a href="/thumbnail/simple">/thumbnail/simple</a>
                <p>
                    Generates a simple thumbnail from an uploaded video. A single video must be uploaded.
                    <br/>Possible query arguments:
                    <ul>
                        <li>width - The width of the thumbnail. Defaults to the width of the video.</li>
                        <li>skip - Skip this number of seconds into the video. Defaults to {{.DefaultSkip}}.</li>
                    </ul>
                </p>
            </li>
            <li>
                POST <a href="/thumbnail/sprite">/thumbnail/sprite</a>
                <p>
                    Generates a sprite thumbnail from an uploaded video. A single video must be uploaded.
                    <br/>Possible query arguments:
                    <ul>
                        <li>width - The width of the thumbnail. Defaults to 180px wide maintaining aspect ratio.</li>
                        <li>skip - Skip this number of seconds into the video. Defaults to {{.DefaultSkip}}.</li>
                        <li>count - The number of thumbnails to include in the sprite. Defaults to {{.DefaultCount}}.</li>
                    </ul>
                </p>
            </li>
            <li>
                GET <a href="/help">/help</a>
                <p>
                    Returns this help page.
                </p>
            </li>
            <li>
                GET <a href="/pulse">/pulse</a>
                <p>
                    Returns health information for the server. See <a href="https://github.com/dulo-tech/amsterdam/wiki/Specification:-Pulse-Protocol">Pulse Protocol</a>.
                </p>
            </li>
        </ul>
    </body>
</html>
`
