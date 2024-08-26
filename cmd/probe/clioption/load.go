package clioption

import (
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type CLIOption struct {
	File f
}

type f struct {
	*os.File
	Size int64
}

func Load() (CLIOption, error) {
	// Options for key features
	filePath := pflag.StringP("file", "f", "", "file path")

	// Options for developer
	debugLogEnable := pflag.Bool("debug", false, "Enable debug logs")

	pflag.Parse()

	file, err := loadFile(filePath)
	if err != nil {
		return CLIOption{}, err
	}

	if *debugLogEnable {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return CLIOption{file}, nil
}
