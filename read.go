package qtffilst

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/tingtt/qtffilst/ilst"

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

		if /* not supporting data */ !ilstDataBox(box) {
			slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB)\n", box.Path, box.DataPosition, box.DataSize))
			continue
		}

		buf := &bytes.Buffer{}
		err = copy(r.f, box.DataPosition, box.DataSize, buf)
		if err != nil {
			return ilst.ItemList{}, err
		}

		ilstBoxName := ilstDataBoxName(box.Path)
		err = itemList.Set(ilstBoxName, buf.Bytes())
		if err != nil {
			return ilst.ItemList{}, fmt.Errorf("%w (id: %s)", err, box.Id)
		}

		if /* binary data */ strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.covr") {
			slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB) binary data (skip display)\n", box.Path, box.DataPosition, box.DataSize))
			continue
		}
		slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB) \"%+v\"\n", box.Path, box.DataPosition, box.DataSize, buf.Bytes()))
	}

	return itemList, nil
}
