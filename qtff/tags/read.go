package tags

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"qtffilst/qtff/tags/meta/ilst"
	"strings"

	"gitlab.com/osaki-lab/iowrapper"
)

type Reader interface {
	Read() (ilst.ItemList, error)
}

func NewReader(f fs.File) (Reader, error) {
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}

	return &reader{
		f:    iowrapper.NewSeeker(f, iowrapper.MaxBufferSize(uint64(stat.Size()))),
		size: stat.Size(),
	}, nil
}

type reader struct {
	f    io.ReadSeeker
	size int64
}

func (r *reader) Read() (ilst.ItemList, error) {
	itemList := ilst.ItemList{}

	for box, err := range Walk(r.f, r.size) {
		if err != nil {
			return ilst.ItemList{}, err
		}

		if /* not supporting data */ !strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.") ||
			box.Id != "data" {
			continue
		}

		_, err := r.f.Seek(box.DataPosition, io.SeekStart)
		if err != nil {
			return ilst.ItemList{}, err
		}

		buf := make([]byte, box.DataSize)
		_, err = io.ReadFull(r.f, buf)
		if err != nil {
			return ilst.ItemList{}, err
		}

		err = itemList.Set(strings.Split(box.Path, ".")[5], buf)
		if err != nil {
			return ilst.ItemList{}, fmt.Errorf("%w (id: %s)", err, box.Id)
		}

		if /* binary data */ strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.covr") {
			slog.Debug(fmt.Sprintf("box: %v (%v, %vB) binary data (skip display)\n", box.Path, box.DataPosition, box.DataSize))
			continue
		}
		slog.Debug(fmt.Sprintf("box: %v (%v, %vB) \"%s\"\n", box.Path, box.DataPosition, box.DataSize, string(buf)))
	}

	return itemList, nil
}
