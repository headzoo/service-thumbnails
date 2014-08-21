package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"github.com/gorilla/mux"
	"github.com/dulo-tech/thumbnailer/handlers"
)

func main() {
	opts := parseFlags()
	
	router := mux.NewRouter()
	router.Handle("/thumbnail/simple", handlers.NewSimple(opts)).Methods("POST")
	router.Handle("/thumbnail/sprite", handlers.NewSprite(opts)).Methods("POST")
	router.Handle("/help", handlers.NewHelp(opts)).Methods("GET")
	router.Handle("/pulse", handlers.NewPulse(opts)).Methods("GET")

	conn := opts.Host + ":" + strconv.Itoa(opts.Port)
	log.Println("Listening for requests on", conn)
	err := http.ListenAndServe(conn, router)
	if err != nil {
		panic(err)
	}
}

// parseFlags parses the command line option flags.
func parseFlags() *thumbnailer.Options {
	opts := &thumbnailer.Options{}
	
	flag.BoolVar(
		&opts.PrintHelp, 
		"help", 
		thumbnailer.OPT_PRINT_HELP, 
		"Display command help.")
	flag.StringVar(
		&opts.Host, 
		"h", 
		thumbnailer.OPT_HOST, 
		"The host name to listen on.")
	flag.IntVar(
		&opts.Port, 
		"p", 
		thumbnailer.OPT_PORT, 
		"The port to listen on.")
	flag.IntVar(
		&opts.SkipSeconds, 
		"s", 
		thumbnailer.OPT_SKIP_SECONDS, 
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&opts.Count, 
		"c", 
		thumbnailer.OPT_COUNT, 
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.Parse()

	if opts.PrintHelp {
		printHelp()
	}
	
	return opts
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer Server v%s\n", thumbnailer.VERSION)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("thumbnailer-server -h <host> -p <port>")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("\t-%s\t%s\n", f.Name, f.Usage)
	})
	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("server -h 127.0.0.1 -p 3366")

	os.Exit(1)
}
