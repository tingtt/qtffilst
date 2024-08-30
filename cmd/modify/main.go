package main

import (
	"log/slog"
	"os"
	"qtffilst/cmd/modify/clioption"
	"qtffilst/qtff/tags"
	"qtffilst/qtff/tags/meta/ilst"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}
}

func run() error {
	cliOption, err := clioption.Load()
	if err != nil {
		return err
	}

	r, err := tags.ParseReadWriter(cliOption.File.File)
	if err != nil {
		return err
	}

	err = r.Write(cliOption.Dest, cliOption.TmpDest, ilst.ItemList{}, nil)
	if err != nil {
		return err
	}

	if !cliOption.KeepTmpFile {
		os.Remove(cliOption.TmpDest.Name())
	}

	return nil
}
