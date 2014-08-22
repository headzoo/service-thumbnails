package cli

import (
	"os"
	"strings"

	"github.com/dulo-tech/thumbnailer/cli/commands"
	"github.com/dulo-tech/thumbnailer/core"
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
		if !fileExists(file) {
			core.VPrintfError("The input file %q does not exist.", file)
			os.Exit(1)
		}
	}

	return files
}

// fileExists returns whether the given file exists.
func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
