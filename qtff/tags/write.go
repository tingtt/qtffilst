package tags

import (
	"os"
	"qtffilst/qtff/tags/meta/ilst"
)

type Writer interface {
	Write(tags ilst.ItemList) error
}

type ReadWriter interface {
	Reader
	Writer
}

func Open(trackPath string) (ReadWriter, error) {
	return open(trackPath)
}

func open(trackPath string) (*readWriter, error) {
	f, err := os.Open(trackPath)
	if err != nil {
		return nil, err
	}
	return parse(f)
}

func ParseReadWriter(f *os.File) (ReadWriter, error) {
	return parse(f)
}

func parse(f *os.File) (*readWriter, error) {
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	return &readWriter{
		reader: reader{
			f:    f,
			size: stat.Size(),
		},
		f: f,
	}, nil
}

type readWriter struct {
	reader
	f *os.File
}

func (r *readWriter) Read() (ilst.ItemList, error) {
	return r.reader.Read()
}

func (r *readWriter) Write(tags ilst.ItemList) error {
	panic("unimplemented")
}
