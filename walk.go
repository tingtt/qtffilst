package qtffilst

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"slices"
	"strings"

	"github.com/tingtt/qtffilst/ilst"
	"github.com/tingtt/qtffilst/internal/binary"
)

type Box struct {
	Name          string
	Level         int8
	Path          string
	DataPosition  int64
	DataSize      int32
	IsContainable bool
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

func walkBoxes(rs io.ReadSeeker, parentEndsAt, offset int64, level int8, path string, yield func(Box) (_continue bool)) (err error) {
	startPosition, err := rs.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if startPosition >= parentEndsAt {
		return nil
	}

	boxSize, boxName, err := readBoxHeader(rs)
	if err != nil {
		return err
	}
	endPosition := startPosition + int64(boxSize)

	_continue := yield(Box{
		Name:          boxName,
		Level:         level,
		Path:          path + "." + boxName,
		DataPosition:  startPosition + 8, /* add size (bytes) of fixed fields (size, name)) */
		DataSize:      boxSize - 8,       /* add size (bytes) of fixed fields (size, name)) */
		IsContainable: false,
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
		_continue := yield(Box{
			Name:          boxName,
			Level:         level,
			Path:          path + "." + boxName,
			DataPosition:  startPosition + 8, /* add size (bytes) of fixed fields (size, name)) */
			DataSize:      boxSize - 8,       /* add size (bytes) of fixed fields (size, name)) */
			IsContainable: true,
		})
		if !_continue {
			return ErrBreakWalk
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
	Write        func([]byte) (size int32, err error)
	InsertNewBox func(name string, data []byte) (size int32, err error)
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

func walkCopyBoxes(rs io.ReadSeeker, parentEndsAt, offset int64, level int8, basePath string, dest io.Writer, yield func(box WritableBox) (_continue bool)) (err error) {
	startPosition, err := rs.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if startPosition >= parentEndsAt {
		return nil
	}

	boxSize, boxName, err := readBoxHeader(rs)
	if err != nil {
		return err
	}
	endPosition := startPosition + int64(boxSize)

	box := Box{
		Name:          boxName,
		Level:         level,
		Path:          basePath + "." + boxName,
		DataPosition:  startPosition + 8, /* add size (bytes) of fixed fields (size, name)) */
		DataSize:      boxSize - 8,       /* add size (bytes) of fixed fields (size, name)) */
		IsContainable: containableBox(boxName),
	}

	if box.IsContainable {
		childOffset := box.DataPosition
		childBuf := &bytes.Buffer{}
		if box.Name == "meta" {
			childBuf.Write(bytes.Repeat([]byte{0x0}, 4))
			childOffset += 4 /* bytes */
		}
		slog.Debug(fmt.Sprintf("%-36s    ->", box.Path))
		err = walkCopyBoxes(rs, endPosition, childOffset, level+1, box.Path, childBuf, yield)
		if err != nil {
			return err
		}
		slog.Debug(fmt.Sprintf("%-36s    [] %8d -> %8d (%+d)\n",
			box.Path, box.DataSize, childBuf.Len(), int32(childBuf.Len())-box.DataSize,
		))

		var insertBoxes map[string][]byte = map[string][]byte{}
		nextBoxWriter := func(name string, data []byte) (size int32, err error) {
			insertBoxes[name] = data
			boxLengthWillWrite := int32(len(data) + 4 + 4)
			return boxLengthWillWrite, nil
		}
		_continue := yield(WritableBox{box,
			nil, // containable box does not support modify content
			nextBoxWriter,
		})
		if !_continue {
			return ErrBreakWalk
		}
		if childBuf.Len() != 0 {
			err = writeBox(dest, box.Name, childBuf.Bytes())
			if err != nil {
				return fmt.Errorf("failed to write box: %w", err)
			}
		}
		for name, data := range insertBoxes {
			insertBoxPath := basePath + "." + name
			insertBoxLength := len(data) + 8
			slog.Debug(fmt.Sprintf("%-36s    +  %8d -> %8d (%+d)\n", insertBoxPath, 0, insertBoxLength, insertBoxLength))
			err = writeBox(dest, name, data)
			if err != nil {
				return fmt.Errorf("failed to write box: %w (%s)", err, insertBoxPath)
			}
		}
	} else /* writable data */ {
		var (
			modified   bool          = false
			newDataBuf *bytes.Buffer = &bytes.Buffer{}
		)

		writer := func(data []byte) (size int32, err error) {
			if modified {
				return 0, fmt.Errorf("`%s` already written", box.Path)
			}
			modified = true
			_, err = newDataBuf.Write(data)
			if err != nil {
				return 0, err
			}
			return int32(newDataBuf.Len()), nil
		}

		_continue := yield(WritableBox{box,
			writer,
			nil, // `data` box does not support to insert next box
		})
		if !_continue {
			return ErrBreakWalk
		}
		if !modified {
			copy(rs, box.DataPosition, box.DataSize, newDataBuf)
		}

		slog.Debug(fmt.Sprintf("%-36s    *  %8d -> %8d (%+d)\n",
			box.Path, box.DataSize, newDataBuf.Len(), int32(newDataBuf.Len())-box.DataSize,
		))
		if /* box data not removed */ newDataBuf.Len() != 0 {
			err = writeBox(dest, box.Name, newDataBuf.Bytes())
			if err != nil {
				return fmt.Errorf("failed to write box: %w (%s)", err, box.Path)
			}
		}
	}

	nextStartPosition, err := rs.Seek(endPosition, io.SeekStart)
	if err != nil {
		return err
	}
	return walkCopyBoxes(rs, parentEndsAt, nextStartPosition, level, basePath, dest, yield)
}

func readBoxHeader(rs io.ReadSeeker) (size int32, name string, err error) {
	size, err = binary.BigEdian.ReadI32(rs)
	if err != nil {
		return 0, "", err
	}

	boxNameBuf, err := binary.Read(rs, 4)
	if err != nil {
		return 0, "", err
	}
	name = string(boxNameBuf)
	if boxNameBuf[0] == 0xA9 {
		name = "(c)" + name[1:]
	}

	return size, name, nil
}

func writeBox(dest io.Writer, name string, data []byte) error {
	_, err := dest.Write(binary.BigEdian.BytesI32(int32(len(data) + 8)))
	if err != nil {
		return err
	}
	if strings.HasPrefix(name, "(c)") {
		name = string([]byte{0xA9}) + name[3:]
	}
	if len(name) != 4 {
		return fmt.Errorf("invalid name length (\"%s\")", name)
	}
	_, err = dest.Write([]byte(name))
	if err != nil {
		return err
	}
	_, err = dest.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func copy(rs io.ReadSeeker, position int64, size int32, w io.Writer) error {
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
