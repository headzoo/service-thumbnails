package main

import (
	"github.com/dulo-tech/thumbnailer/cli"
	"github.com/dulo-tech/thumbnailer/core"
	"github.com/dulo-tech/thumbnailer/http"
)

func main() {
	core.Init()
	switch core.Opts.Mode {
	case "cli":
		cli.Go()
	case "http":
		http.Go()
	}
}
