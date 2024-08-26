package binary

import (
	"encoding/binary"
	"io"
)

var BigEdian = &bigEndian{}

type bigEndian struct{}

func (*bigEndian) ReadI16(r io.Reader) (int16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return -1, err
	}
	num := binary.BigEndian.Uint16(buf)
	return int16(num), nil
}

func (*bigEndian) ReadI32(r io.Reader) (int32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return -1, err
	}
	num := binary.BigEndian.Uint32(buf)
	return int32(num), nil
}

func Read(r io.Reader, bytes uint) ([]byte, error) {
	buf := make([]byte, bytes)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
