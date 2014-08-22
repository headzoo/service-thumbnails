package handlers

import (
	"github.com/dulo-tech/thumbnailer/core"
	"net/http"
	"strconv"

	"github.com/dulo-tech/go-pulse/pulse"
)

// PulseHandler is an HTTP handler for the pulse protocol.
type PulseHandler struct {
	Handler
}

// NewPulse creates and returns a new PulseHandler instance.
func NewPulse() *PulseHandler {
	return &PulseHandler{}
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (h *PulseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := pulse.New(r.RemoteAddr, core.AppVersion)
	p.WhiteList = pulseIPWhiteList
	p.RequestHeaders = make(pulse.Headers, len(r.Header))
	for key, headers := range r.Header {
		p.RequestHeaders[key] = headers[0]
	}

	numRequests++
	p.Set("num-requests", strconv.Itoa(numRequests))
	p.Set("num-errors", strconv.Itoa(numErrors))
	for key, value := range p.ResponseHeaders() {
		w.Header().Set(key, value)
	}
	w.WriteHeader(p.StatusCode())
	w.Write([]byte(p.ResponseBody()))
}
