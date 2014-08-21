package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"flag"
	"fmt"
	"github.com/dulo-tech/go-pulse/pulse"
	"github.com/dulo-tech/thumbnailer/ffmpeg"
	"github.com/dulo-tech/thumbnailer/thumbnailer"
	"github.com/rakyll/magicmime"
)

// Misc constants.
const (
	MAX_MEMORY           = 1 * 1024 * 1024
	DEFAULT_MIME_TYPE    = "binary/octet-stream"
	DEFAULT_WIDTH_BIG    = 0
	DEFAULT_WIDTH_SPRITE = 180
)

// Default values for command line options.
const (
	OPT_HOST         = "127.0.0.1"
	OPT_PORT         = 8080
	OPT_TEMP_DIR     = "/tmp"
	OPT_SKIP_SECONDS = 0
	OPT_VERBOSE      = false
	OPT_PRINT_HELP   = false
)

// Upload stores the values of an uploaded file.
type Upload struct {
	Name     string
	Size     int64
	MimeType string
	Temp     string
}

// Options stores the command line options.
type Options struct {
	Host          string
	Port          int
	TempDirectory string
	SkipSeconds   int
	PrintHelp     bool
}

var opts = Options{}

func main() {
	flag.BoolVar(&opts.PrintHelp, "help", OPT_PRINT_HELP, "Display command help.")
	flag.StringVar(&opts.Host, "h", OPT_HOST, "The host name to listen on.")
	flag.IntVar(&opts.Port, "p", OPT_PORT, "The port to listen on.")
	flag.IntVar(&opts.SkipSeconds, "s", OPT_SKIP_SECONDS, "Skip this number of seconds into the video before thumbnailing.")
	flag.StringVar(&opts.TempDirectory, "d", OPT_TEMP_DIR, "Temp directory.")
	flag.Parse()
	if opts.PrintHelp {
		printHelp()
	}

	http.HandleFunc("/thumbnail/big", handleBigThumbnail)
	http.HandleFunc("/thumbnail/sprite", handleSpriteThumbnail)
	http.HandleFunc("/help", handleHelp)
	http.HandleFunc("/pulse", handlePulse)

	conn := opts.Host + ":" + strconv.Itoa(opts.Port)
	log.Println("Listening for requests on", conn)
	err := http.ListenAndServe(conn, nil)
	if err != nil {
		panic(err)
	}
}

// handleBigThumbnail is the http callback to create a big thumbnail.
func handleBigThumbnail(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	width := DEFAULT_WIDTH_BIG
	skip := opts.SkipSeconds

	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}
	if s, ok := query["skip"]; ok {
		skip = atoi(s[0])
	}

	temp := getTempFile()
	ff := ffmpeg.NewFFmpeg(file.Temp)
	ff.TmpDirectory = opts.TempDirectory
	ff.SkipSeconds = skip

	err := ff.CreateThumbnail(width, temp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	writeFileToResponse(temp, w)
}

// handleSpriteThumbnail is the http callback to create a sprite thumbnail.
func handleSpriteThumbnail(w http.ResponseWriter, r *http.Request) {
	file := getFile(w, r)
	if file == nil {
		return
	}

	width := DEFAULT_WIDTH_SPRITE
	skip := opts.SkipSeconds

	query := r.URL.Query()
	if w, ok := query["width"]; ok {
		width = atoi(w[0])
	}
	if s, ok := query["skip"]; ok {
		skip = atoi(s[0])
	}

	temp := getTempFile()
	ff := ffmpeg.NewFFmpeg(file.Temp)
	ff.TmpDirectory = opts.TempDirectory
	ff.SkipSeconds = skip

	interval := int(ff.Length())
	if interval > thumbnailer.NUM_THUMBNAILS {
		interval = interval / thumbnailer.NUM_THUMBNAILS
	}

	err := ff.CreateThumbnailSprite(interval, width, temp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=thumbnail.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	writeFileToResponse(temp, w)
}

// handleHelp is the http callback to display the help page.
func handleHelp(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/help.html")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	t.Execute(w, "T")
}

// handleHelp is the http callback to handle Pulse Protocol requests.
func handlePulse(w http.ResponseWriter, r *http.Request) {
	p := pulse.New(r.RemoteAddr, thumbnailer.VERSION)
	p.WhiteList = []string{
		"127.*",
		"10.0.*",
		"192.168.*",
	}
	p.RequestHeaders = make(pulse.Headers, len(r.Header))
	for key, headers := range r.Header {
		p.RequestHeaders[key] = headers[0]
	}

	for key, value := range p.ResponseHeaders() {
		w.Header().Set(key, value)
	}
	w.WriteHeader(p.StatusCode())
	w.Write([]byte(p.ResponseBody()))
}

// getFile returns the uploaded file.
func getFile(w http.ResponseWriter, r *http.Request) *Upload {
	files, _ := writeUploadedFiles(r)
	if len(files) > 1 {
		w.WriteHeader(400)
		w.Write([]byte("Only a single file allowed."))
		return nil
	} else if len(files) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("No files uploaded."))
		return nil
	}

	log.Printf("Got upload %#v\n", files[0])
	return &files[0]
}

// writeUploadedFiles writes all uploaded files to the temp dir.
func writeUploadedFiles(r *http.Request) ([]Upload, error) {
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		return nil, err
	}

	files := []Upload{}
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			fin, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}

			fout, err := ioutil.TempFile("/tmp", "upload")
			if err != nil {
				return nil, err
			}
			defer fout.Close()

			size, err := io.Copy(fout, fin)
			if err != nil {
				return nil, err
			}
			fout.Close()

			files = append(files, Upload{
				Name:     fileHeader.Filename,
				Size:     size,
				MimeType: getMimeType(fout.Name()),
				Temp:     fout.Name(),
			})
		}
	}

	return files, nil
}

// writeFileToResponse writes a file to the http response.
func writeFileToResponse(file string, w http.ResponseWriter) error {
	fout, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fout.Close()
	io.Copy(w, fout)

	return nil
}

// getMimeType returns the file mime type.
func getMimeType(file string) string {
	mm, err := magicmime.New(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR)
	if err != nil {
		return DEFAULT_MIME_TYPE
	}

	mimetype, err := mm.TypeByFile(file)
	if err != nil {
		return DEFAULT_MIME_TYPE
	}

	return mimetype
}

// getTempFile returns the name of a new temp file.
func getTempFile() string {
	temp, _ := ioutil.TempFile("/tmp", "thumb")
	temp.Close()
	return temp.Name()
}

// printHelp() prints the command line help and exits.
func printHelp() {
	fmt.Printf("Thumbnailer Server v%s\n", thumbnailer.VERSION)
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("server -h <host> -p <port>")
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

// atoi converts a string to an integer.
func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		i = 0
	}

	return i
}
