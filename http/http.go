package http

import (
	"net/http"
	"strconv"

	"github.com/dulo-tech/thumbnailer/core"
	"github.com/dulo-tech/thumbnailer/http/handlers"
	"github.com/gorilla/mux"
)

func Go() {
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
