package binary

import (
	"encoding/binary"
	"io"
)

var BigEdian = &bigEndian{}

type bigEndian struct{}

func (*bigEndian) BytesI8(num int8) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(num))
	return buf[1:]
}

func (*bigEndian) BytesI16(num int16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(num))
	return buf
}

func (*bigEndian) BytesI32(num int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(num))
	return buf
}

func (*bigEndian) ReadI8(r io.Reader) (int8, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	num := binary.BigEndian.Uint16(buf)
	return int8(num), nil
}

func (*bigEndian) ReadI16(r io.Reader) (int16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	num := binary.BigEndian.Uint16(buf)
	return int16(num), nil
}

func (*bigEndian) ReadI32(r io.Reader) (int32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
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
