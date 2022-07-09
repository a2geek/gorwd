package rwd

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
func (r *File) List() ([]*Entry, error) {
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

	files := []*Entry{}
	for i := 0; i < int(numberOfFiles); i++ {
		entry, err := r.readEntry()
		if err != nil {
			return nil, err
		}
		files = append(files, entry)
	}
	return files, err
}

func (r *File) readEntry() (*Entry, error) {
	var fileNameLength int16
	err := binary.Read(r.file, binary.LittleEndian, &fileNameLength)
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
		rwdFile:  r,
		Filename: sb.String(),
		Offset:   data[0],
		Length:   data[2],
	}
	return &entry, nil
}

// Save will save this (modified) RWD archive back to disk.
func (r *File) Save() error {
	dir, _ := filepath.Split(r.file.Name())
	newf, err := ioutil.TempFile(dir, "rwd-")
	if err != nil {
		return err
	}
	defer newf.Close()

	header, err := r.Header()
	if err != nil {
		return err
	}
	headerLength := binary.Size(header)

	entries, err := r.List()
	if err != nil {
		return err
	}

	trailer, err := r.Trailer()
	if err != nil {
		return err
	}

	err = binary.Write(newf, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		offset, err := newf.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		length, err := entry.WriteTo(newf)
		if err != nil {
			return err
		}

		// Don't want to change settings until _after_ content is saved
		entry.Offset = int32(int(offset) - headerLength)
		entry.Length = int32(length)
	}

	var numberOfFiles int32 = int32(len(entries))
	err = binary.Write(newf, binary.LittleEndian, &numberOfFiles)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err = r.writeEntry(newf, entry)
		if err != nil {
			return err
		}
	}

	err = binary.Write(newf, binary.LittleEndian, trailer)
	if err != nil {
		return err
	}

	err = r.Close()
	if err != nil {
		return err
	}

	err = os.Rename(r.file.Name(), r.file.Name()+".bak")
	if err != nil {
		return err
	}

	err = newf.Close()
	if err != nil {
		return err
	}

	return os.Rename(newf.Name(), r.file.Name())
}

func (r *File) writeEntry(f *os.File, entry *Entry) error {
	var fileNameLength int16 = int16(len(entry.Filename))
	err := binary.Write(f, binary.LittleEndian, &fileNameLength)
	if err != nil {
		return err
	}
	var fileNameRune = make([]int16, fileNameLength)
	for i, ch := range entry.Filename {
		fileNameRune[i] = int16(ch)
	}
	err = binary.Write(f, binary.LittleEndian, &fileNameRune)
	if err != nil {
		return err
	}

	data := [6]int32{}
	data[0] = entry.Offset
	data[2] = entry.Length
	return binary.Write(f, binary.LittleEndian, &data)
}
