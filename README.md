# What is this?

This is a little utility to help manipulate RWD archive files.

The end goal is to get a utility available for those of us in Linux installing games like Kohan II using Proton for Windows emulation. Some games do not work without a little bit of additional configuration.

## Usage

```
$ gorwd --help
Usage:
  gorwd [OPTIONS] [file(s)...]

This is a little utility to help manipulate RWD archive files.

The end goal is to get a utility available for those of us in Linux installing games
like Kohan II using Proton for Windows emulation. Some games do not work without a
little bit of additional configuration.

See https://github.com/a2geek/gorwd for more details.


Application Options:
  -f, --file=    File to process [$GORWD_FILENAME]
  -d, --dir=     Read or write files to this directory (valid for extract only)

Commands:
  -l, --list     List contents of file
  -i, --info     Display details of file structures
  -x, --extract  Extract files
  -p, --pipe     Pipe files to stdout
  -u, --update   Update file in archive

Help Options:
  -h, --help     Show this help message

Arguments:
  file(s):       glob patterns to match (ex: *.ttf)

```

### List example

```
$ gorwd -lf ~/.steam/steam/steamapps/common/Kohan\ II/Warchest/Warchest.rwd 
  1. Fonts/font_medium.tgi (o=2343, l=1468)
  2. Fonts/font_tiny.tgi (o=5254, l=1596)
  3. Fonts/font_large.tgi (o=595, l=1748)
  4. Fonts/font_small.tgi (o=3811, l=1443)
  5. AVars_version.tgi (o=0, l=595)
```

### Detailed Usage

* [Kohan II](docs/KOHAN-II.md) setup

## File Structure

This is an updated file structure from what was originally posted. Note that some names were altered based on text in the trailer section.

```
header (30 bytes)
    byte[4]         magic bytes ('TGCK')
    uint32[3]       values (2, 3, 2; uncertain of significance)
    uint16          length of following string
    wchar[4]        name(?); padded with zeroes
    uint32          unknown value (hash code? timestamp?)

file content

footer (aka file directory)
    uint32          number of files
    file entry      one for each file
        uint16      length of file name (in unicode characters)
        wchar[]     unicode file name (relative to Data directory)
        uint16      offset to file content from end of header
        uint16      unknown
        uint16      file length
        uint16      unknown
        uint16      unknown
        uint16      unknown

trailer (288 bytes)
    sections (3 of them - header, files, footer)
      wchar[32]     name of section, zero terminated
      uint32        offset of section from beginning of file
      uint32        unknown
      uint32        length of section
      uint32        unknown
      uint32        unknown
      uint32        unknown
      uint32        alternate length of section
      uint32        unknown

char is 8 bits; wchar_t is 16 bits.
int is 32 bits; short is 8 bits.
```

## Resources

* [Kohan II](https://www.protondb.com/app/97130) entry on ProtonDB. The primary solution requires the user to also have a Windows machine available.
* [Note.txt](Note.txt) file. This was in the archive posted to ProtonDB and has enough details to write code and get Kohan II running!
