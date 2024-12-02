package qtffilst

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"maps"
	"os"

	"github.com/tingtt/iterutil"
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
	modifyItemIds := maps.Collect(ilst.Values(&newItemList))
	for _, deleteId := range deleteIds {
		modifyItemIds[deleteId] = nil
	}

	if len(modifyItemIds) == 0 && len(deleteIds) == 0 {
		_, err := r.f.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.Copy(dest, r.f)
		return err
	}

	ilstSizeDiff := int32(0)
	oldItemList := ilst.ItemList{}
	lastLoadedIlstBoxName := ""

	// Modify `.moov.udta.meta.ilst`
	for box, err := range WalkSupportedWritabelBox(WritableWalk(r.f, r.size, tmpDest)) {
		if err != nil {
			return err
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

		err = oldItemList.SetDecoded(ilstBoxName, buf.Bytes())
		if err != nil {
			return err
		}

		if /* expect remove */ v, exists := modifyItemIds[ilstBoxName]; exists && v == nil {
			// Remove matched box
			_, err = box.Write(nil)
			if err != nil {
				return err
			}
			slog.Info("remove", slog.String("id", ilstBoxName), slog.String("diff", fmt.Sprintf("%+d", -box.DataSize)))

			ilstSizeDiff -= box.DataSize /* current box size */ + 16 /* parent box header size */
			continue
		}

		// Modify matched box
		boxNameMatcher := func(v ilst.EncodedValue) bool { return v.Id == ilstBoxName }
		for value, err := range iterutil.FilterKeyFunc(ilst.EncodedValues(&newItemList), boxNameMatcher) {
			if err != nil {
				return err
			}

			if _, exists := modifyItemIds[ilstBoxName]; !exists {
				continue
			}
			delete(modifyItemIds, ilstBoxName)
			size, err := box.Write(value.Bytes)
			if err != nil {
				return err
			}
			slog.Info("modify", slog.String("id", ilstBoxName), slog.String("diff", fmt.Sprintf("%+d", size-box.DataSize)))

			ilstSizeDiff += size - box.DataSize
			break
		}
	}

	stat, err := tmpDest.Stat()
	if err != nil {
		return err
	}

	if len(modifyItemIds) == 0 {
		_, err := tmpDest.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.Copy(tmpDest2, tmpDest)
		if err != nil {
			return err
		}
	} else {
		// Create remaining items `.moov.udta.meta.ilst`
		matchLastItemListBox := func(box WritableBox) bool {
			return box.Path == fmt.Sprintf(".moov.udta.meta.ilst.%s", lastLoadedIlstBoxName)
		}
		for box, err := range iterutil.FilterKeyFunc(WritableWalk(tmpDest, stat.Size(), tmpDest2), matchLastItemListBox) {
			if err != nil {
				return err
			}
			if box.InsertNewBox == nil {
				panic(fmt.Sprintf("box writer is nil (path: %s)", box.Path))
			}

			for value, err := range ilst.EncodedValues(&newItemList) {
				if err != nil {
					return err
				}

				if _, exists := modifyItemIds[value.Id]; !exists {
					continue
				}
				delete(modifyItemIds, value.Id)
				buf := &bytes.Buffer{}
				err = writeBox(buf, "data", value.Bytes)
				if err != nil {
					return err
				}
				size, err := box.InsertNewBox(value.Id, buf.Bytes())
				if err != nil {
					return err
				}
				slog.Info("append", slog.String("id", value.Id), slog.String("diff", fmt.Sprintf("%+d", size)))

				ilstSizeDiff += size
			}
		}
	}

	stat2, err := tmpDest2.Stat()
	if err != nil {
		return err
	}

	if ilstSizeDiff == 0 {
		slog.Debug("skip modification of chunk offset because .moov.udta.meta.ilst has no size changes", slog.String("diff", fmt.Sprintf("%+d", ilstSizeDiff)))
		_, err := tmpDest2.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.Copy(dest, tmpDest2)
		return err
	}

	mdatFoundBeforeIlst, err := mdatBoxIsBeforeIlst(Walk(tmpDest2, stat2.Size()))
	if err != nil {
		return err
	}
	if mdatFoundBeforeIlst {
		slog.Debug("skip modification of chunk offset because .mdat exists before .moov.udta.meta.ilst", slog.String("diff", fmt.Sprintf("%+d", ilstSizeDiff)))
		tmpDest2.Seek(0, io.SeekStart)
		_, err := io.Copy(dest, tmpDest2)
		return err
	}

	// Modify `.moov.trak.mdia.minf.stbl.stco`
	slog.Info("modify chunk offsets", slog.String("diff", fmt.Sprintf("%+d", ilstSizeDiff)))
	matchSampleTableChunkOffsetBox := func(box WritableBox) bool {
		return box.Path == ".moov.trak.mdia.minf.stbl.stco"
	}
	for box, err := range iterutil.FilterKeyFunc(WritableWalk(tmpDest2, stat2.Size(), dest), matchSampleTableChunkOffsetBox) {
		if err != nil {
			return err
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

func WalkSupportedWritabelBox(rw iter.Seq2[WritableBox, error]) iter.Seq2[WritableBox, error] {
	matchSupporedBox := func(v WritableBox) bool {
		return ilstDataBox(v.Box)
	}
	return iterutil.FilterKeyFunc(rw, matchSupporedBox)
}

var (
	ErrIlstBoxDoesNotExist = errors.New(".moov.udta.meta.ilst does not exists")
)

func mdatBoxIsBeforeIlst(seq iter.Seq2[Box, error]) (bool, error) {
	mdatFound := false
	for box, err := range seq {
		if err != nil {
			return false, err
		}

		switch box.Path {
		case ".mdat":
			mdatFound = true
		case ".moov.udta.meta.ilst":
			return mdatFound, nil
		}
	}
	return false, errors.New(".moov.udta.meta.ilst does not exists")
}
