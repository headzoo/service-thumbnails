package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/dulo-tech/thumbnailer/core"
	"github.com/dulo-tech/thumbnailer/handlers"
	"github.com/gorilla/mux"
)

func main() {
	Init()
	router := mux.NewRouter()
	router.Handle("/thumbnail/simple", handlers.NewSimple()).Methods("POST")
	router.Handle("/thumbnail/sprite", handlers.NewSprite()).Methods("POST")
	router.Handle("/help", handlers.NewHelp()).Methods("GET")
	router.Handle("/pulse", handlers.NewPulse()).Methods("GET")

	core.VPrintf("Listening for requests on %s:%d...", core.Opts.Host, core.Opts.Port)
	conn := core.Opts.Host + ":" + strconv.Itoa(core.Opts.Port)
	err := http.ListenAndServe(conn, router)
	if err != nil {
		panic(err)
	}
}

// Init parses the command line option flags.
func Init() {
	core.Init()
	flag.StringVar(
		&core.Opts.Host,
		"h",
		core.OptDefaultHost,
		"The host name to listen on.")
	flag.IntVar(
		&core.Opts.Port,
		"p",
		core.OptDefaultPort,
		"The port to listen on.")
	flag.Parse()

	if core.Opts.PrintHelp {
		core.ExecuteHelpTemplate(core.Opts, serverHelpTemplate)
	}
}

const serverHelpTemplate = `Thumbnailer HTTP Server v{{.Version}} - Video thumbnail generating HTTP server.

USAGE:
	thumbnailer-server -h <host> -p <port>

OPTIONS:

{{.Flags}}
EXAMPLES:
	server -h 127.0.0.1 -p 3366
`
