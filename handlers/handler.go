package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/dulo-tech/thumbnailer/core"
	"github.com/rakyll/magicmime"
)

// Misc constants.
const (
	// Max amount of memory to use when reading uploaded files.
	MAX_MEMORY = 1 * 1024 * 1024

	// Default mime type when an uploaded file type cannot be determined.
	DEFAULT_MIME_TYPE = "binary/octet-stream"

	// Default width for simple thumbnails.
	DEFAULT_WIDTH_SIMPLE = 0

	// Default width for sprite thumbnails.
	DEFAULT_WIDTH_SPRITE = 180
)

var (
	// numRequests counts the number of total requests handled by the http server.
	numRequests int = 0

	// numErrors counts the number of errors generated by the http server.
	numErrors int = 0

	// pulseIPWhiteList is a list of ip masks allowed to access the pulse end point.
	pulseIPWhiteList = []string{
		"127.*",
		"10.0.*",
		"192.168.*",
	}
)

// Upload stores the values of an uploaded file.
type Upload struct {
	// Name is the original name of an uploaded file.
	Name string

	// Size is the file byte size.
	Size int64

	// MimeType is the file mime type.
	MimeType string

	// Temp is the path to the temp copy of the uploaded file.
	Temp string
}

// Handler is the default HTTP handler.
type Handler struct {
	// opts is the command line options.
	opts *core.Options
}

// New creates and returns a new Handler instance.
func New(opts *core.Options) *Handler {
	return &Handler{
		opts: opts,
	}
}

// getFile returns the uploaded file.
func getFile(w http.ResponseWriter, r *http.Request) *Upload {
	files, _ := writeUploadedFiles(r)
	if len(files) > 1 {
		numErrors++
		w.WriteHeader(400)
		w.Write([]byte("Only a single file allowed."))
		return nil
	} else if len(files) == 0 {
		numErrors++
		w.WriteHeader(400)
		w.Write([]byte("No files uploaded."))
		return nil
	}

	core.VPrintf("Got upload %#v\n", files[0])
	return &files[0]
}

// writeUploadedFiles writes all uploaded files to the temp dir.
func writeUploadedFiles(r *http.Request) ([]Upload, error) {
	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		numErrors++
		return nil, err
	}

	files := []Upload{}
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			fin, err := fileHeader.Open()
			if err != nil {
				numErrors++
				return nil, err
			}

			fout, err := ioutil.TempFile("/tmp", "upload")
			if err != nil {
				numErrors++
				return nil, err
			}
			defer fout.Close()

			size, err := io.Copy(fout, fin)
			if err != nil {
				numErrors++
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
		numErrors++
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
		numErrors++
		return DEFAULT_MIME_TYPE
	}

	mimetype, err := mm.TypeByFile(file)
	if err != nil {
		numErrors++
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

// atoi converts a string to an integer.
func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		numErrors++
		i = 0
	}

	return i
}
