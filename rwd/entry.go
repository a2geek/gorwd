package rwd

import (
	"encoding/binary"
	"io"
)

// Entry represents what amounts to a file in the archive.
type Entry struct {
	rwdFile  *File
	Filename string
	Offset   int32
	Length   int32
}

// WriteTo will copy the bytes for this Entry. Note the logic to adjust for header.
func (e *Entry) WriteTo(writer io.Writer) (int64, error) {
	// Need to add Header to the offset; choosing to do it dynamically.
	header := Header{}
	headerBytes := binary.Size(header)

	var offset int
	offset = int(e.Offset)
	offset += headerBytes

	sr := io.NewSectionReader(e.rwdFile.file, int64(offset), int64(e.Length))
	return io.Copy(writer, sr)
}
