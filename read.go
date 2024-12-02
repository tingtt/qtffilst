package qtffilst

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"log/slog"
	"strings"

	"github.com/tingtt/iterutil"
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

	for box, err := range WalkSupportedBox(r.f, r.size) {
		if err != nil {
			return ilst.ItemList{}, err
		}

		buf := &bytes.Buffer{}
		err = copy(r.f, box.DataPosition, box.DataSize, buf)
		if err != nil {
			return ilst.ItemList{}, err
		}

		ilstBoxName := ilstDataBoxName(box.Path)
		err = itemList.SetDecoded(ilstBoxName, buf.Bytes())
		if err != nil {
			return ilst.ItemList{}, fmt.Errorf("%w (id: %s)", err, box.Name)
		}

		if /* binary data */ strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.covr") {
			slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB) binary data (skip display)\n", box.Path, box.DataPosition, box.DataSize))
			continue
		}
		slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB) \"%+v\"\n", box.Path, box.DataPosition, box.DataSize, buf.Bytes()))
	}

	return itemList, nil
}

func WalkSupportedBox(rs io.ReadSeeker, size int64) iter.Seq2[Box, error] {
	matchSupporedBox := func(v Box) bool {
		if ilstDataBox(v) {
			return true
		}
		slog.Debug(fmt.Sprintf("box: %-36s (%v, %vB)\n", v.Path, v.DataPosition, v.DataSize))
		return false
	}

	return iterutil.FilterKeyFunc(Walk(rs, size), matchSupporedBox)
}
