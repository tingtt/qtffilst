package clioption

import (
	"errors"
	"os"
)

func loadFile(filePath *string) (f, error) {
	if *filePath == "" {
		return f{}, errors.New("CLI option `-f` cannot be empty")
	}
	file, err := os.Open(*filePath)
	if err != nil {
		return f{}, err
	}
	stat, err := file.Stat()
	if err != nil {
		return f{}, err
	}
	return f{file, stat.Size()}, nil
}

func createDestFile(destFilePath, tmpDestFilePath *string) (dest *os.File, tmpDest, tmpDest2 *os.File, err error) {
	if *destFilePath == "" {
		return nil, nil, nil, errors.New("CLI option `--out`,`-o` cannot be empty")
	}
	if *destFilePath == *tmpDestFilePath {
		return nil, nil, nil, errors.New("CLI option `--out` and `--tmp` cannot be same")
	}
	if *tmpDestFilePath == "" {
		*tmpDestFilePath = *destFilePath + ".tmp"
	}

	dest, err = os.Create(*destFilePath)
	if err != nil {
		return nil, nil, nil, err
	}
	tmpDest, err = os.Create(*tmpDestFilePath + "1")
	if err != nil {
		return nil, nil, nil, err
	}
	tmpDest2, err = os.Create(*tmpDestFilePath + "2")
	if err != nil {
		return nil, nil, nil, err
	}
	return dest, tmpDest, tmpDest2, nil
}
