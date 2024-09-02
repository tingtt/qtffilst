package ilst

import (
	"iter"
	"reflect"
)

func Ids() iter.Seq[string] {
	return func(yield func(id string) (_continue bool)) {
		rv := reflect.ValueOf(ItemList{})
		rt := rv.Type()

		for i := range make([]interface{}, rt.NumField()) {
			f := rt.Field(i)
			_continue := yield(f.Tag.Get("id"))
			if !_continue {
				break
			}
		}
	}
}

type writableValue struct {
	id     string
	value  any
	set    Setter
	remove Remover
}
type Setter func(buf []byte) error
type Remover func()

func IdWriters(ilst *ItemList) iter.Seq2[string, writableValue] {
	return func(yield func(id string, v writableValue) (_continue bool)) {
		rv := reflect.ValueOf(ilst).Elem()
		rt := rv.Type()

		for i := range make([]interface{}, rt.NumField()) {
			id := rt.Field(i).Tag.Get("id")
			_continue := yield(id, newWritableValue(id, rv.Field(i)))
			if !_continue {
				break
			}
		}
	}
}

func newWritableValue(id string, field reflect.Value) writableValue {
	v := field.Interface()

	set := func(buf []byte) (err error) {
		switch v.(type) {
		case *InternationalText:
			err = setField(field, decodeInternationalText, buf)
		case *Genre:
			err = setField(field, decodeGenre, buf)
		case *BoolWithHeader0x15_0:
			err = setField(field, decodeBoolWithHeader0x15_0, buf)
		case *Int16WithHeader0x15_0:
			err = setField(field, decodeInt16WithHeader0x15_0, buf)
		case *TrackNumber:
			err = setField(field, decodeTrackNumber, buf)
		case *DiskNumber:
			err = setField(field, decodeDiskNumber, buf)
		default:
			panic("unsupported item type")
		}
		return err
	}
	remove := func() { field.Set(reflect.ValueOf(nil)) }

	return writableValue{id, v, set, remove}
}

func setField[T any](field reflect.Value, decode func(data []byte) (T, error), data []byte) error {
	var v T
	v, err := decode(data)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(&v))
	return nil
}
