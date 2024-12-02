package clioption

import (
	"log/slog"
	"os"

	"github.com/spf13/pflag"
	"github.com/tingtt/qtffilst/ilst"
)

type CLIOption struct {
	File          f
	Dest          *os.File
	TmpDest       *os.File
	TmpDest2      *os.File
	KeepTmpFile   bool
	ItemList      *ilst.ItemList
	DeleteItemIds []string
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
	changeDatas := pflag.StringSliceP("data", "d", nil, "Write QTFF ItemList tag.\n\tformat: <id>=<value>")
	removeIds := pflag.StringSliceP("rm", "r", nil, "")

	// Options for developer
	debugLogEnable := pflag.Bool("debug", false, "Enable debug logs")

	pflag.Parse()

	file, err := loadFile(filePath)
	if err != nil {
		return CLIOption{}, err
	}
	dest, tmpDest, tmpDest2, err := createDestFile(destPath, tmpDestPath)
	if err != nil {
		return CLIOption{}, err
	}

	itemList, deleteIds, err := loadChanges(*changeDatas, *removeIds)
	if err != nil {
		return CLIOption{}, err
	}

	if *debugLogEnable {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return CLIOption{file, dest, tmpDest, tmpDest2, *keepTmpFile, itemList, deleteIds}, nil
}
