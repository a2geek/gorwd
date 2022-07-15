package rwd

import (
	"encoding/binary"
	"io"
	"os"
)

// Entry represents what amounts to a file in the archive.
type Entry struct {
	rwdFile *rwdFile

	Filename string
	Offset   int32
	Length   int32
}

// WriteTo will copy the bytes for this Entry. Note the logic to adjust for header.
func (e *Entry) WriteTo(writer io.Writer) (int64, error) {
	replacementFilename, changed := e.rwdFile.newFiles[e.Filename]
	if changed {
		return e.writeReplacementContentTo(writer, replacementFilename)
	} else {
		return e.writeExistingContentTo(writer)
	}
}
func (e *Entry) writeReplacementContentTo(writer io.Writer, replacementFilename string) (int64, error) {
	_, err := os.Stat(replacementFilename)
	if err != nil {
		return 0, err
	}

	src, err := os.Open(replacementFilename)
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
	e.rwdFile.newFiles[e.Filename] = filename
}
