package rwd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var MAGIC = []byte{0x54, 0x47, 0x43, 0x4b}

// New will construct an rwd.File and validate the type of file.
func New(filename string) (File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = CheckMagic(f)
	if err != nil {
		return nil, err
	}
	rwd := rwdFile{
		file:     f,
		newFiles: make(map[string]string),
	}
	return &rwd, nil
}

// CheckMagic will validate the "magic" bytes for the file signature.
func CheckMagic(f *os.File) error {
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	var data [4]byte
	err = binary.Read(f, binary.LittleEndian, &data)
	if err != nil {
		return err
	}
	if bytes.Compare(data[:], MAGIC) != 0 {
		return fmt.Errorf("Unexpected magic bytes: %02x%02x%02x%02x",
			data[0], data[1], data[2], data[3])
	}
	return nil
}
