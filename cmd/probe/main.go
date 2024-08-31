package main

import (
	"fmt"
	"iter"
	"log/slog"
	"os"
	"qtffilst/cmd/probe/clioption"
	"qtffilst/qtff/tags"
	"qtffilst/qtff/tags/meta/ilst"
	"reflect"
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

	r, err := tags.NewReader(cliOption.File)
	if err != nil {
		return err
	}

	tag, err := r.Read()
	if err != nil {
		return err
	}

	fmt.Println("---")
	for f := range iterateIDs(&tag) {
		if f.tag.Get("id") == "covr" {
			fmt.Printf("%s: binary data (skip display)\n", f.tag.Get("id"))
		} else {
			v := f.value.Elem()
			if v.IsValid() {
				fmt.Printf("%s: %+v\n", f.tag.Get("id"), v)
			}
		}
	}

	return nil
}

type field struct {
	name  string
	tag   reflect.StructTag
	value reflect.Value
}

func iterateIDs(itemList *ilst.ItemList) iter.Seq[field] {
	return func(yield func(f field) bool) {
		rv := reflect.ValueOf(itemList).Elem()
		rt := rv.Type()

		for i := range make([]interface{}, rt.NumField()) {
			f := rt.Field(i)
			_continue := yield(field{f.Name, f.Tag, rv.FieldByName(f.Name)})
			if !_continue {
				break
			}
		}
	}
}
