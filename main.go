package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/a2geek/gorwd/rwd"
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
		fmt.Printf("%3d. %s (o=%d, l=%d)\n", i+1, entry.Filename, entry.Offset, entry.Length)
	}
}
