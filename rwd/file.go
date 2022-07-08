package rwd

import (
	"encoding/binary"
	"io"
	"os"
	"strings"
)

type File struct {
	file *os.File
}

// Close closes the RWD File, rendering it unusable for I/O.
// Close will return an error if it has already been called.
func (r *File) Close() error {
	return r.file.Close()
}

// Header will read the header record from the RWD File.
// There is no caching and can be used to re-read an updated Header.
func (r *File) Header() (*Header, error) {
	_, err := r.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	header := Header{}
	err = binary.Read(r.file, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}

// Trailer will read the trailer record from the RWD File.
// There is no caching and can be used to re-read an updated Trailer.
func (r *File) Trailer() (*Trailer, error) {
	trailer := Trailer{}
	_, err := r.file.Seek(-int64(binary.Size(trailer)), io.SeekEnd)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r.file, binary.LittleEndian, &trailer)
	if err != nil {
		return nil, err
	}
	return &trailer, nil
}

// List will read a list of all Files from the RWD file.
func (r *File) List() (*[]Entry, error) {
	trailer, err := r.Trailer()
	if err != nil {
		return nil, err
	}

	_, err = r.file.Seek(-int64(binary.Size(trailer))-int64(trailer.DirectoryLength), io.SeekEnd)
	if err != nil {
		return nil, err
	}

	var numberOfFiles int32
	err = binary.Read(r.file, binary.LittleEndian, &numberOfFiles)
	if err != nil {
		return nil, err
	}

	files := []Entry{}
	for i := 0; i < int(numberOfFiles); i++ {
		var fileNameLength int16
		err = binary.Read(r.file, binary.LittleEndian, &fileNameLength)
		if err != nil {
			return nil, err
		}
		var fileNameRune = make([]int16, fileNameLength)
		err = binary.Read(r.file, binary.LittleEndian, &fileNameRune)
		if err != nil {
			return nil, err
		}
		sb := strings.Builder{}
		for _, ch := range fileNameRune {
			sb.WriteRune(rune(ch))
		}

		data := [6]int32{}
		err = binary.Read(r.file, binary.LittleEndian, &data)
		if err != nil {
			return nil, err
		}

		entry := Entry{
			Filename: sb.String(),
			Offset:   data[0],
			Length:   data[2],
		}
		files = append(files, entry)
	}
	return &files, err
}
