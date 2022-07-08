package rwd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

var MAGIC = []byte{0x54, 0x47, 0x43, 0x4b}

func New(filename string) (*File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	err = CheckMagic(f)
	if err != nil {
		return nil, err
	}
	return &File{file: f}, nil
}

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
