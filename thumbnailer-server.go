package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
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

	thumbnailer.Verbose("Listening for requests on %s:%d...", opts.Host, opts.Port)
	conn := opts.Host + ":" + strconv.Itoa(opts.Port)
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
	flag.BoolVar(
		&opts.Verbose,
		"v",
		thumbnailer.OPT_VERBOSE,
		"Verbose output.")
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

	thumbnailer.VerboseOutput = opts.Verbose
	if opts.PrintHelp {
		printHelp(opts)
	}

	return opts
}

// printHelp() prints the command line help and exits.
func printHelp(opts *thumbnailer.Options) {
	if thumbnailer.VerboseOutput || opts.PrintHelp {
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
	}

	os.Exit(1)
}
