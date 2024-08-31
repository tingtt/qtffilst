package tags

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"qtffilst/internal/binary"
	"qtffilst/qtff/tags/meta/ilst"
	"slices"
)

type Box struct {
	Id           string
	Level        int64
	Path         string
	DataPosition int64
	DataSize     int64
}

const (
	ROOT_LEVEL = 0
	ROOT_PATH  = ""
)

var ErrBreakWalk = errors.New("break walk")

func Walk(rs io.ReadSeeker, size int64) iter.Seq2[Box, error] {
	return func(yield func(Box, error) (_continue bool)) {
		acturlYield := func(t Box) (_continue bool) {
			return yield(t, nil)
		}
		err := walkBoxes(rs, size, 0, ROOT_LEVEL /* start from */, ROOT_PATH /* start from */, acturlYield)
		if err != nil && !errors.Is(err, ErrBreakWalk) {
			yield(Box{}, err)
		}
	}
}

func walkBoxes(rs io.ReadSeeker, parentEndsAt, offset, level int64, path string, yield func(Box) (_continue bool)) (err error) {
	startPosition, err := rs.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if startPosition >= parentEndsAt {
		return nil
	}

	boxSize, err := binary.BigEdian.ReadI32(rs)
	if err != nil {
		return err
	}
	endPosition := startPosition + int64(boxSize)

	boxNameBuf, err := binary.Read(rs, 4)
	if err != nil {
		return err
	}
	boxName := string(boxNameBuf)
	if boxNameBuf[0] == 0xA9 {
		boxName = "(c)" + boxName[1:]
	}

	_continue := yield(Box{
		Id:           boxName,
		Level:        level,
		Path:         path + "." + boxName,
		DataPosition: startPosition + 8,  /* add size (bytes) of fixed fields (size, name)) */
		DataSize:     int64(boxSize) - 8, /* add size (bytes) of fixed fields (size, name)) */
	})
	if !_continue {
		return ErrBreakWalk
	}

	if containableBox(boxName) {
		childOffset := startPosition + 8 /* add size (bytes) of fixed fields (size, name)) */
		if boxName == "meta" {
			childOffset += 4 /* bytes */
		}
		err = walkBoxes(rs, endPosition, childOffset, level+1, path+"."+boxName, yield)
		if err != nil {
			return err
		}
	}

	nextStartPosition, err := rs.Seek(endPosition, io.SeekStart)
	if err != nil {
		return err
	}
	return walkBoxes(rs, parentEndsAt, nextStartPosition, level, path, yield)
}

func containableBox(boxName string) bool {
	for id := range ilst.Ids() {
		if id == boxName {
			return true
		}
	}
	return slices.Contains([]string{"moov",
		"udta", "meta", "ilst",
		"trak", "mdia", "minf", "stbl",
	}, boxName)
}

type WritableBox struct {
	Box
	Write func([]byte) (size int64, err error)
}

func WritableWalk(rs io.ReadSeeker, size int64, dest io.Writer) iter.Seq2[WritableBox, error] {
	return func(yield func(WritableBox, error) (_continue bool)) {
		acturlYield := func(t WritableBox) (_continue bool) {
			return yield(t, nil)
		}
		err := walkCopyBoxes(rs, size, 0, ROOT_LEVEL /* start from */, ROOT_PATH /* start from */, dest, acturlYield)
		if err != nil && !errors.Is(err, ErrBreakWalk) {
			yield(WritableBox{}, err)
		}
	}
}

func walkCopyBoxes(rs io.ReadSeeker, parentEndsAt, offset, level int64, basePath string, dest io.Writer, yield func(box WritableBox) (_continue bool)) (err error) {
	startPosition, err := rs.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if startPosition >= parentEndsAt {
		return nil
	}

	boxSize, err := binary.BigEdian.ReadI32(rs)
	if err != nil {
		return err
	}
	endPosition := startPosition + int64(boxSize)

	boxNameBuf, err := binary.Read(rs, 4)
	if err != nil {
		return err
	}
	boxName := string(boxNameBuf)
	if boxNameBuf[0] == 0xA9 {
		boxName = "(c)" + boxName[1:]
	}

	currPath := basePath + "." + boxName

	if containableBox(boxName) {
		childOffset := startPosition + 8 /* add size (bytes) of fixed fields (size, name)) */
		childBuf := &bytes.Buffer{}
		if boxName == "meta" {
			childBuf.Write(bytes.Repeat([]byte{0x0}, 4))
			childOffset += 4 /* bytes */
		}
		slog.Debug(fmt.Sprintf("%-36s    ->", currPath))
		err = walkCopyBoxes(rs, endPosition, childOffset, level+1, currPath, childBuf, yield)
		if err != nil {
			return err
		}
		slog.Debug(fmt.Sprintf("%-36s    [] %8d -> %8d (%+d)\n",
			currPath, boxSize-8, childBuf.Len(), childBuf.Len()-int(boxSize-8),
		))
		err = writeBox(dest, boxNameBuf, childBuf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write box: %w", err)
		}
	} else /* writable data */ {
		dataPosition := startPosition + 8 /* add size (bytes) of fixed fields (size, name)) */
		dataSize := int64(boxSize) - 8    /* add size (bytes) of fixed fields (size, name)) */

		var (
			modified   bool          = false
			newDataBuf *bytes.Buffer = &bytes.Buffer{}
		)

		writer := func(data []byte) (size int64, err error) {
			if modified {
				return 0, fmt.Errorf("`%s` already written", currPath)
			}
			modified = true
			_, err = newDataBuf.Write(data)
			if err != nil {
				return 0, err
			}
			return int64(newDataBuf.Len()), nil
		}

		_continue := yield(WritableBox{
			Box{
				Id:           boxName,
				Level:        level,
				Path:         currPath,
				DataPosition: dataPosition,
				DataSize:     dataSize,
			},
			writer,
		})
		if !_continue {
			return ErrBreakWalk
		}
		if !modified {
			copy(rs, dataPosition, dataSize, newDataBuf)
		}

		slog.Debug(fmt.Sprintf("%-36s    *  %8d -> %8d (%+d)\n",
			currPath, dataSize, newDataBuf.Len(), newDataBuf.Len()-int(dataSize),
		))
		err = writeBox(dest, boxNameBuf, newDataBuf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write box: %w", err)
		}
	}

	nextStartPosition, err := rs.Seek(endPosition, io.SeekStart)
	if err != nil {
		return err
	}
	return walkCopyBoxes(rs, parentEndsAt, nextStartPosition, level, basePath, dest, yield)
}

func writeBox(dest io.Writer, name []byte, data []byte) error {
	_, err := dest.Write(binary.BigEdian.BytesI32(int32(len(data) + 8)))
	if err != nil {
		return err
	}
	if len(name) != 4 {
		return fmt.Errorf("invalid name length")
	}
	_, err = dest.Write(name)
	if err != nil {
		return err
	}
	_, err = dest.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func copy(rs io.ReadSeeker, position, size int64, w io.Writer) error {
	_, err := rs.Seek(position, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(rs, buf)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
