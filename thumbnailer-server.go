package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/dulo-tech/thumbnailer/handlers"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"github.com/gorilla/mux"
)

func main() {
	opts := parseFlags()

	router := mux.NewRouter()
	router.Handle("/thumbnail/simple", handlers.NewSimple(opts)).Methods("POST")
	router.Handle("/thumbnail/sprite", handlers.NewSprite(opts)).Methods("POST")
	router.Handle("/help", handlers.NewHelp(opts)).Methods("GET")
	router.Handle("/pulse", handlers.NewPulse(opts)).Methods("GET")

	thumbnailer.VPrintf("Listening for requests on %s:%d...", opts.Host, opts.Port)
	conn := opts.Host + ":" + strconv.Itoa(opts.Port)
	err := http.ListenAndServe(conn, router)
	if err != nil {
		panic(err)
	}
}

// parseFlags parses the command line option flags.
func parseFlags() *thumbnailer.Options {
	opts := thumbnailer.FlagOptions()
	flag.StringVar(
		&opts.Host,
		"h",
		thumbnailer.OptDefaultHost,
		"The host name to listen on.")
	flag.IntVar(
		&opts.Port,
		"p",
		thumbnailer.OptDefaultPort,
		"The port to listen on.")
	flag.Parse()

	if opts.PrintHelp {
		thumbnailer.ExecuteHelpTemplate(opts, serverHelpTemplate)
	}
	return opts
}

const serverHelpTemplate = `Thumbnailer HTTP Server v{{.Version}} - Video thumbnail generating HTTP server.

USAGE:
	thumbnailer-server -h <host> -p <port>

OPTIONS:

{{.Flags}}
EXAMPLES:
	server -h 127.0.0.1 -p 3366
`
