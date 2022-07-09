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
			Pipe    bool `short:"p" long:"pipe" description:"Pipe files to stdout"`
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

	// This is terrible but not grokking a better way for mutually exclusive options
	if countBools(options.Commands.Pipe, options.Commands.Extract, options.Commands.List) != 1 {
		panic(errors.New("Please select one command at a time"))
	}
	command := List
	if options.Commands.Pipe {
		command = Pipe
	} else if options.Commands.Extract {
		command = Extract
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
				command(i, entry)
				break
			}
		}
	}
}

func List(n int, entry rwd.Entry) {
	fmt.Printf("%3d. %s (o=%d, l=%d)\n", n+1, entry.Filename, entry.Offset, entry.Length)
}
func Pipe(n int, entry rwd.Entry) {
	fmt.Printf("File: %s\n", entry.Filename)
	entry.WriteTo(os.Stdout)
}
func Extract(n int, entry rwd.Entry) {
	// TODO
}

func countBools(bools ...bool) int {
	count := 0
	for _, b := range bools {
		if b {
			count = count + 1
		}
	}
	return count
}
