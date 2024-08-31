package qtffilst

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/tingtt/qtffilst/ilst"
	"github.com/tingtt/qtffilst/internal/binary"
)

type Writer interface {
	Write(dest, tmpDest, tmpDest2 *os.File, tags ilst.ItemList, deleteIds []string) error
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

func (r *readWriter) Write(dest, tmpDest, tmpDest2 *os.File, newItemList ilst.ItemList, deleteIds []string) error {
	// TODO: remove test data
	newItemList.TitleC = ilst.NewInternationalText("modified title")
	newItemList.ReleaseDate = ilst.NewInternationalText(time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC).Format(time.DateOnly))
	modifyItemIds := map[string]any{"(c)nam": nil, "rldt": nil}
	// TODO: remove test data

	ilstSizeDiff := int32(0)
	oldItemList := ilst.ItemList{}
	lastLoadedIlstBoxName := ""

	// Modify `.moov.udta.meta.ilst`
	for box, err := range WritableWalk(r.f, r.size, tmpDest) {
		if err != nil {
			return err
		}

		if /* not supporting data */ !ilstDataBox(box.Box) {
			continue
		}
		if box.Write == nil {
			panic(fmt.Sprintf("box writer is nil (path: %s)", box.Path))
		}

		buf := &bytes.Buffer{}
		err = copy(r.f, box.DataPosition, box.DataSize, buf)
		if err != nil {
			return err
		}

		ilstBoxName := ilstDataBoxName(box.Path)
		lastLoadedIlstBoxName = ilstBoxName

		err = oldItemList.Set(ilstBoxName, buf.Bytes())
		if err != nil {
			return err
		}

		// TODO: write new data to matched box
		if ilstBoxName == "(c)nam" {
			delete(modifyItemIds, ilstBoxName)
			newTitleCBuf, err := newItemList.TitleC.Bytes()
			if err != nil {
				return err
			}
			size, err := box.Write(newTitleCBuf)
			if err != nil {
				return err
			}

			ilstSizeDiff += size - box.DataSize
		}
		if ilstBoxName == "rldt" {
			delete(modifyItemIds, ilstBoxName)
			dataBuf, err := newItemList.ReleaseDate.Bytes()
			if err != nil {
				return err
			}
			size, err := box.Write(dataBuf)
			if err != nil {
				return err
			}

			ilstSizeDiff += size - box.DataSize
		}
	}

	stat, err := tmpDest.Stat()
	if err != nil {
		return err
	}

	// Create remaining items `.moov.udta.meta.ilst`
	for box, err := range WritableWalk(tmpDest, stat.Size(), tmpDest2) {
		if err != nil {
			return err
		}

		if /* not supporting data */ !strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.") || !box.IsContainable {
			continue
		}
		if box.InsertNewBox == nil {
			panic(fmt.Sprintf("box writer is nil (path: %s)", box.Path))
		}

		if box.Name == lastLoadedIlstBoxName {
			dataBuf, err := newItemList.ReleaseDate.Bytes()
			if err != nil {
				return err
			}
			buf := &bytes.Buffer{}
			err = writeBox(buf, []byte("data"), dataBuf)
			if err != nil {
				return err
			}
			size, err := box.InsertNewBox("rldt", buf.Bytes())
			if err != nil {
				return err
			}
			ilstSizeDiff += size
		}
	}

	stat2, err := tmpDest2.Stat()
	if err != nil {
		return err
	}

	// Modify `.moov.trak.mdia.minf.stbl.stco`
	for box, err := range WritableWalk(tmpDest2, stat2.Size(), dest) {
		if err != nil {
			return err
		}

		if box.Path != ".moov.trak.mdia.minf.stbl.stco" {
			continue
		}

		// Data format
		// https://developer.apple.com/documentation/quicktime-file-format/chunk_offset_atom
		buf := &bytes.Buffer{}

		{ // copy fixed fields (version, flags, number of entries)
			err = copy(tmpDest2, box.DataPosition, 8, buf)
			if err != nil {
				return err
			}
		}

		var entryCount int32
		{ // get entry count
			_, err = tmpDest2.Seek(box.DataPosition+4, io.SeekStart)
			if err != nil {
				return err
			}
			entryCount, err = binary.BigEdian.ReadI32(tmpDest2)
			if err != nil {
				return err
			}
		}

		// create new chunk offset table
		positionOfChunkOffsetTable := box.DataPosition + 8
		tmpDest2.Seek(positionOfChunkOffsetTable, io.SeekStart) // Seek to chunk offset table
		for range entryCount {
			offset, err := binary.BigEdian.ReadI32(tmpDest2)
			if err != nil {
				return err
			}
			slog.Debug(fmt.Sprintf("offset: %8d -> %8d (%+d)\n", offset, offset+ilstSizeDiff, ilstSizeDiff))
			newOffset := binary.BigEdian.BytesI32(offset + ilstSizeDiff)
			_, err = buf.Write(newOffset)
			if err != nil {
				return err
			}
		}

		// write new `.moov.trak.mdia.minf.stbl.stco`
		_, err = box.Write(buf.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
