package tags

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"qtffilst/internal/binary"
	"qtffilst/qtff/tags/meta/ilst"
	"strings"
)

type Writer interface {
	Write(dest, tmpDest *os.File, tags ilst.ItemList, deleteIds []string) error
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

func (r *readWriter) Write(dest, tmpDest *os.File, tags ilst.ItemList, deleteIds []string) error {
	ilstSizeDiff := int32(0)

	// Modify `.moov.udta.meta.ilst`
	for box, err := range WritableWalk(r.f, r.size, tmpDest) {
		if err != nil {
			return err
		}

		if !strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.") || box.Id != "data" {
			continue
		}

		if box.Path == ".moov.udta.meta.ilst.(c)nam.data" {
			buf := &bytes.Buffer{}
			err = copy(r.f, box.DataPosition, box.DataSize, buf)
			if err != nil {
				return err
			}

			data := append(buf.Bytes(), []byte(" modified")...)
			size, err := box.Set(data)
			if err != nil {
				return err
			}

			ilstSizeDiff += int32(size - box.DataSize)
		}
	}

	stat, err := tmpDest.Stat()
	if err != nil {
		return err
	}

	// Modify `.moov.trak.mdia.minf.stbl.stco`
	for box, err := range WritableWalk(tmpDest, stat.Size(), dest) {
		if err != nil {
			return err
		}

		if box.Path == ".moov.trak.mdia.minf.stbl.stco" {
			// Data format
			// https://developer.apple.com/documentation/quicktime-file-format/chunk_offset_atom
			buf := &bytes.Buffer{}

			{ // copy fixed fields (version, flags, number of entries)
				err = copy(tmpDest, box.DataPosition, 8, buf)
				if err != nil {
					return err
				}
			}

			var entryCount int32
			{ // get entry count
				entryCountBuf := &bytes.Buffer{}
				err = copy(tmpDest, box.DataPosition+4, 4, entryCountBuf)
				if err != nil {
					return err
				}
				entryCount, err = binary.BigEdian.ReadI32(entryCountBuf)
				if err != nil {
					return err
				}
			}

			positionOfChunkOffsetTable := box.DataPosition + 8
			tmpDest.Seek(positionOfChunkOffsetTable, io.SeekStart) // Seek to chunk offset table
			for range entryCount {
				offset, err := binary.BigEdian.ReadI32(tmpDest)
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
			_, err = box.Set(buf.Bytes())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
