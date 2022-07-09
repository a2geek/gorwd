package rwd

import (
	"encoding/binary"
	"io"
	"os"
)

// Entry represents what amounts to a file in the archive.
type Entry struct {
	rwdFile             *File
	replacementFilename string

	Filename string
	Offset   int32
	Length   int32
}

// WriteTo will copy the bytes for this Entry. Note the logic to adjust for header.
func (e *Entry) WriteTo(writer io.Writer) (int64, error) {
	if len(e.replacementFilename) > 0 {
		return e.writeReplacementContentTo(writer)
	} else {
		return e.writeExistingContentTo(writer)
	}
}
func (e *Entry) writeReplacementContentTo(writer io.Writer) (int64, error) {
	_, err := os.Stat(e.replacementFilename)
	if err != nil {
		return 0, err
	}

	src, err := os.Open(e.replacementFilename)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	return io.Copy(writer, src)
}
func (e *Entry) writeExistingContentTo(writer io.Writer) (int64, error) {
	// Need to add Header to the offset; choosing to do it dynamically.
	header := Header{}
	headerBytes := binary.Size(header)

	var offset int
	offset = int(e.Offset)
	offset += headerBytes

	sr := io.NewSectionReader(e.rwdFile.file, int64(offset), int64(e.Length))
	return io.Copy(writer, sr)
}

// ReplaceWithFile will flag this entry for replacement with contents from a filesystem file.
func (e *Entry) ReplaceWithFile(filename string) {
	e.replacementFilename = filename
}
