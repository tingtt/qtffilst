package clioption

import (
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type CLIOption struct {
	File        f
	Dest        *os.File
	TmpDest     *os.File
	KeepTmpFile bool
}

type f struct {
	*os.File
	Size int64
}

func Load() (CLIOption, error) {
	// Options for key features
	filePath := pflag.StringP("file", "f", "", "src file path")
	destPath := pflag.StringP("out", "o", "", "dest file path")
	tmpDestPath := pflag.String("tmp", "", "tmp dest file path")
	keepTmpFile := pflag.Bool("keep", false, "keep tmp dest file")

	// Options for developer
	debugLogEnable := pflag.Bool("debug", false, "Enable debug logs")

	pflag.Parse()

	file, err := loadFile(filePath)
	if err != nil {
		return CLIOption{}, err
	}
	dest, tmpDest, err := createDestFile(destPath, tmpDestPath)
	if err != nil {
		return CLIOption{}, err
	}

	if *debugLogEnable {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return CLIOption{file, dest, tmpDest, *keepTmpFile}, nil
}
