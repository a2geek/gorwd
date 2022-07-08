package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/a2geek/gorwd/rwd"
	"github.com/gobwas/glob"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Commands struct {
			List    bool `short:"l" long:"list" description:"List contents of file"`
			Extract bool `short:"x" long:"extract" description:"Extract files"`
		} `group:"Commands" required:"true"`

		Filename string `short:"f" long:"file" env:"GORWD_FILENAME" description:"File to process"`

		Args struct {
			Glob []string `description:"glob patterns to match (ex: *.ttf)"`
		} `positional-args:"yes"`
	}

	_, err := flags.Parse(&options)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		panic(err)
	}

	if !options.Commands.List {
		panic(errors.New("Only support list at this time"))
	}

	var globs []glob.Glob
	for _, globPattern := range options.Args.Glob {
		glob, err := glob.Compile(globPattern)
		if err != nil {
			panic(err)
		}
		globs = append(globs, glob)
	}
	if len(globs) == 0 {
		defaultGlob := glob.MustCompile("*")
		globs = append(globs, defaultGlob)
	}

	f, err := rwd.New(options.Filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	entries, err := f.List()
	if err != nil {
		panic(err)
	}

	for i, entry := range *entries {
		for _, glob := range globs {
			if glob.Match(entry.Filename) {
				fmt.Printf("%3d. %s (o=%d, l=%d)\n", i+1, entry.Filename, entry.Offset, entry.Length)
			}
		}
	}
}
