package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/user"
	"reflect"
	"strings"
	"text/template"

	"github.com/dulo-tech/service-thumbnails/cli"
	"github.com/dulo-tech/service-thumbnails/core"
	"github.com/dulo-tech/service-thumbnails/http"
	"path"
	"strconv"
)

func main() {
	config()
	if core.Opts.PrintHelp {
		executeHelpTemplate()
	}
	if core.Opts.Mode == "cli" {
		if core.Opts.InFile == "" || core.Opts.OutFile == "" || core.Opts.ThumbType == "" {
			core.VPrintfError("Missing -i, -o, or -t.")
			executeHelpTemplate()
		}
		if core.Opts.ThumbType != "sprite" && core.Opts.ThumbType != "simple" {
			core.VPrintfError("Invalid thumbnail type.")
			executeHelpTemplate()
		}
		cli.Go()
	} else if core.Opts.Mode == "http" {
		http.Go()
	} else {
		core.VPrintfError("Invalid mode.")
		executeHelpTemplate()
	}
}

// config parses command line arguments and reads from configuration files.
//
// First attempts to read from a configuration file specified at the command
// line. Then tries to read from .service-thumbnails.conf in the user's home
// directory. Then tries reading from /etc/service-thumbnails.conf. Finally
// parses the command line arguments. The command line arguments override
// configuration file values.
func config() {
	confCli := ""
	confHome := ""
	confEtc := "/etc/service-thumbnails.conf"
	u, err := user.Current()
	if err == nil {
		confHome = path.Join(u.HomeDir, "/.service-thumbnails.conf")
	}

	set := flag.NewFlagSet("conf", flag.ContinueOnError)
	set.StringVar(
		&confCli,
		"conf",
		"",
		"Path to configuration file.")
	set.Parse(os.Args[1:])
	if confCli != "" {
		readConfigFile(confCli, core.Opts)
	} else if core.FileExists(confHome) {
		readConfigFile(confHome, core.Opts)
	} else if core.FileExists(confEtc) {
		readConfigFile(confEtc, core.Opts)
	}

	flag.String(
		"conf",
		"",
		"Path to configuration file.")
	flag.StringVar(
		&core.Opts.Mode,
		"m",
		core.Opts.Mode,
		"Running mode, either 'cli' or 'http'. Defaults to 'cli'.")
	flag.BoolVar(
		&core.Opts.PrintHelp,
		"help",
		core.Opts.PrintHelp,
		"Display command help.")
	flag.BoolVar(
		&core.Opts.Quiet,
		"q",
		core.Opts.Quiet,
		"Run in quiet mode.")
	flag.IntVar(
		&core.Opts.SkipSeconds,
		"s",
		core.Opts.SkipSeconds,
		"Skip this number of seconds into the video before thumbnailing.")
	flag.IntVar(
		&core.Opts.Count,
		"c",
		core.Opts.Count,
		"Number of thumbs to generate in a sprite. 30 is the default.")
	flag.IntVar(
		&core.Opts.Width,
		"w",
		core.Opts.Width,
		"The thumbnail width. Overrides the built in defaults.")
	flag.StringVar(
		&core.Opts.ThumbType,
		"t",
		core.Opts.ThumbType,
		"The type of thumbnail to generate. 'simple' is the default.")
	flag.StringVar(
		&core.Opts.InFile,
		"i",
		core.Opts.InFile,
		"The input video file. Separate multiple files with a comma.")
	flag.StringVar(
		&core.Opts.OutFile,
		"o",
		core.Opts.OutFile,
		"The output image file.")
	flag.StringVar(
		&core.Opts.Host,
		"h",
		core.Opts.Host,
		"The host name to listen on.")
	flag.IntVar(
		&core.Opts.Port,
		"p",
		core.Opts.Port,
		"The port to listen on.")
	flag.Parse()
}

// ExecuteHelpTemplate() prints the command line help using the given template and exits.
func executeHelpTemplate() {
	buff := bytes.Buffer{}
	flag.VisitAll(func(f *flag.Flag) {
		buff.WriteString(fmt.Sprintf("\t-%-8s%s\n", f.Name, f.Usage))
	})

	data := struct {
		Version string
		Flags   string
	}{
		core.AppVersion,
		buff.String(),
	}

	t, _ := template.New("help").Parse(helpTemplate)
	t.Execute(os.Stdout, data)

	os.Exit(1)
}

// readConfigFile sets the values in opts by reading a configuration file.
func readConfigFile(file string, opts *core.Options) {
	of := reflect.ValueOf(opts)
	st := of.Elem()

	fin, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	line := 1
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		text := strings.Trim(scanner.Text(), " \t\r\n")
		if !strings.HasPrefix(text, "#") && text != "" {
			parts := strings.SplitN(text, "=", 2)
			if len(parts) != 2 {
				panic(fmt.Sprintf("Invalid configuration at line %d: %q", line, text))
			}

			field := st.FieldByName(parts[0])
			if field.IsValid() && field.CanSet() {
				if field.Kind() == reflect.Int {
					x, err := strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						panic(fmt.Sprintf("Invalid configuration at line %d. Expecting integer: %q", line, text))
					}
					if !field.OverflowInt(x) {
						field.SetInt(x)
					}
				} else if field.Kind() == reflect.String {
					field.SetString(parts[1])
				} else if field.Kind() == reflect.Bool {
					x := strings.ToLower(parts[1])
					if x == "true" || x == "yes" || x == "1" {
						field.SetBool(true)
					} else {
						field.SetBool(false)
					}
				} else {
					panic(fmt.Sprintf("Invalid configuration at line %d: %q", line, text))
				}
			}
		}

		line++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// helpTemplate is the template used for displaying command line help.
const helpTemplate = `Thumbnailer v{{.Version}} - Used to generate thumbnails from videos.

The app can run as a command line app or as an http server. Use the -m
switch to change modes. See usage and examples of each below.

Options can be set using a configuration file by using the -conf switch.

	service-thumbnails -conf thumbnails.conf

When the -conf switch isn't used the app will try to read from
.service-thumbnails.conf in the user's directory. Finally the app will try to
read from /etc/service-thumbnails.conf.

See the example thumbnails.conf for a description of each configuration value.


OPTIONS:

{{.Flags}}
CLI USAGE:
	thumbnailer -t <sprite|simple> -i <video> -o <image>

	<sprite|simple> determines the type of thumbnail being generated. Either
	a sprite or a simple thumbnail. Simple is the default when not specified.

	<video> is one or more source videos. Separate multiple videos with commas.

	<image> may contain the place holders {name} and {type} which correspond
	to the name of the source video (without file extension) and the type of
	of thumbnail. One of 'sprite' or 'simple'. The <image> may also contain
	the verb %d which will be replaced with the file number. See the fmt package
	for more information on verbs.

CLI EXAMPLES:

	thumbnailer -t sprite -i source.mp4 -o thumb.jpg
	thumbnailer -i source1.mp4,source2.mp4 -o out%02d.jpg
	thumbnailer -t sprite -i source.mp4 -o thumb{name}{type}.jpg

HTTP USAGE:
	thumbnailer -m http -h <host> -p <port>

HTTP EXAMPLES:
	thumbnailer -m http -h 127.0.0.1 -p 3366
`
