package cli

import (
	"os"
	"strings"

	"github.com/dulo-tech/service-thumbnails/cli/commands"
	"github.com/dulo-tech/service-thumbnails/core"
)

func Go() {
	router := commands.NewRouter(splitFiles(core.Opts.InFile), core.Opts.OutFile)
	router.Command("simple", commands.NewSimple())
	router.Command("sprite", commands.NewSprite())
	err := router.Route(core.Opts.ThumbType)
	if err != nil {
		panic(err)
	}
}

// splitFiles converts a comma separated list of files into an array of file names.
func splitFiles(inFiles string) []string {
	files := strings.Split(inFiles, ",")
	for i, f := range files {
		files[i] = strings.Trim(f, " ")
	}

	for _, file := range files {
		if !core.FileExists(file) {
			core.VErrorf("The input file %q does not exist.", file)
			os.Exit(1)
		}
	}

	return files
}
