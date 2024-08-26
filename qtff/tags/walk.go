package tags

import (
	"errors"
	"io"
	"iter"
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
	return slices.Contains([]string{"moov", "udta", "meta", "ilst"}, boxName)
}
