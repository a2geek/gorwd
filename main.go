package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/a2geek/gorwd/rwd"
	"github.com/gobwas/glob"
	flags "github.com/jessevdk/go-flags"
)

func main() {
	var options struct {
		Commands struct {
			List    bool `short:"l" long:"list" description:"List contents of file"`
			Info    bool `short:"i" long:"info" description:"Display details of file structures"`
			Extract bool `short:"x" long:"extract" description:"Extract files"`
			Pipe    bool `short:"p" long:"pipe" description:"Pipe files to stdout"`
			Update  bool `short:"u" long:"update" description:"Update file in archive"`
		} `group:"Commands" required:"true"`

		Filename  string `short:"f" long:"file" env:"GORWD_FILENAME" description:"File to process"`
		Directory string `short:"d" long:"dir" description:"Read or write files to this directory (valid for extract only)"`

		Args struct {
			Glob []string `description:"glob patterns to match (ex: *.ttf)" positional-arg-name:"file(s)"`
		} `positional-args:"yes"`
	}

	parser := flags.NewParser(&options, flags.Default)
	parser.LongDescription = `
	This is a little utility to help manipulate RWD archive files.

	The end goal is to get a utility available for those of us in Linux installing games 
	like Kohan II using Proton for Windows emulation. Some games do not work without a 
	little bit of additional configuration.

	See https://github.com/a2geek/gorwd for more details.
	`

	_, err := parser.Parse()
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		panic(err)
	}

	// This is terrible but not grokking a better way for mutually exclusive options
	if countBools(options.Commands.Pipe, options.Commands.Extract, options.Commands.List, options.Commands.Update, options.Commands.Info) != 1 {
		panic(errors.New("Please select one command at a time"))
	}
	command := List
	if options.Commands.Pipe {
		command = Pipe
	} else if options.Commands.Extract {
		command = Extract
	} else if options.Commands.Update {
		command = Update
	} else if options.Commands.Info {
		// Do nothing; we do not want to run details
		//command = Info
	}

	// This adjusts so we can use 'os.Chdir' later on
	absolutePath, err := filepath.Abs(options.Filename)
	if err != nil {
		panic(err)
	}
	options.Filename = absolutePath

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

	if len(options.Directory) > 0 {
		err := os.Chdir(options.Directory)
		if err != nil {
			panic(err)
		}
	}

	f, err := rwd.New(options.Filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if options.Commands.Info {
		err = Info(f)
		if err != nil {
			panic(err)
		}
		return
	}

	entries, err := f.List()
	if err != nil {
		panic(err)
	}

	for i, entry := range entries {
		for _, glob := range globs {
			if glob.Match(entry.Filename) {
				command(i, entry)
				break
			}
		}
	}

	if options.Commands.Update {
		err = f.Save()
		if err != nil {
			panic(err)
		}
	}
}

func List(n int, entry *rwd.Entry) {
	fmt.Printf("%3d. %s (o=%d, l=%d)\n", n+1, entry.Filename, entry.Offset, entry.Length)
}
func Pipe(n int, entry *rwd.Entry) {
	fmt.Printf("File: %s\n", entry.Filename)
	entry.WriteTo(os.Stdout)
}
func Extract(n int, entry *rwd.Entry) {
	path, _ := filepath.Split(entry.Filename)

	if len(path) > 0 {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(entry.Filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = entry.WriteTo(f)
	if err != nil {
		panic(err)
	}
}
func Update(n int, entry *rwd.Entry) {
	_, err := os.Stat(entry.Filename)
	if err != nil {
		fmt.Printf("WARNING: File '%s' does not exist.", entry.Filename)
		return
	}

	fmt.Printf("Updating file '%s' in archive.\n", entry.Filename)
	entry.ReplaceWithFile(entry.Filename)
}

func Info(f rwd.File) error {
	header, err := f.Header()
	if err != nil {
		return err
	}

	fmt.Printf("*** HEADER ***\n")
	fmt.Printf("Magic Bytes: %08x (%s)\n", header.Magic, byteToString(header.Magic[:]))
	for n, value := range header.Value {
		fmt.Printf("Value %d:     %08x (%d)\n", n, value, value)
	}
	fmt.Printf("Name Length: %04x (%d)\n", header.NameLength, header.NameLength)
	fmt.Printf("Name:        %s\n", wcharToString(header.Name[:], header.NameLength))
	fmt.Printf("Unknown      %08x (%d)\n", header.Unknown, header.Unknown)
	fmt.Println()

	trailer, err := f.Trailer()
	if err != nil {
		return err
	}

	fmt.Printf("*** TRAILER ***\n")
	printSection("Header", trailer.Header)
	printSection("Files", trailer.Files)
	printSection("Footer", trailer.Footer)

	return nil
}
func printSection(name string, section rwd.Section) {
	fmt.Printf("--> Section    %s\n", name)
	fmt.Printf("    Name:      %s\n", wcharzToString(section.NameZ[:]))
	fmt.Printf("    Offset:    %08x (%d)\n", section.Offset, section.Offset)
	fmt.Printf("    Unknown1:  %08x (%d)\n", section.Unknown1, section.Unknown1)
	fmt.Printf("    Length:    %08x (%d)\n", section.Length, section.Length)
	fmt.Printf("    Unknown3:  %08x (%d)\n", section.Unknown3, section.Unknown3)
	fmt.Printf("    Unknown4:  %08x (%d)\n", section.Unknown4, section.Unknown4)
	fmt.Printf("    Unknown5:  %08x (%d)\n", section.Unknown5, section.Unknown5)
	fmt.Printf("    Alt. Len.: %08x (%d)\n", section.AlternateLength, section.AlternateLength)
	fmt.Printf("    Unknown6:  %08x (%d)\n", section.Unknown7, section.Unknown7)
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

func wcharToString(data []uint16, len uint16) string {
	sb := strings.Builder{}
	for n, ch := range data {
		if n >= int(len) {
			break
		}
		sb.WriteRune(rune(ch))
	}
	return sb.String()
}

func wcharzToString(data []uint16) string {
	sb := strings.Builder{}
	for _, ch := range data {
		if ch == 0 {
			break
		}
		sb.WriteRune(rune(ch))
	}
	return sb.String()
}

func byteToString(data []byte) string {
	sb := strings.Builder{}
	for _, ch := range data {
		sb.WriteRune(rune(ch))
	}
	return sb.String()
}
