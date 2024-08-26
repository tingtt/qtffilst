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
