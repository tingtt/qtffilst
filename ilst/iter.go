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

func Values(ilst *ItemList) iter.Seq2[string, any] {
	return func(yield func(string, any) (_continue bool)) {
		rv := reflect.ValueOf(ilst).Elem()
		rt := rv.Type()

		for i := range make([]interface{}, rt.NumField()) {
			id := rt.Field(i).Tag.Get("id")
			value := rv.Field(i).Interface()

			if reflect.ValueOf(value).IsNil() {
				continue
			}
			_continue := yield(id, value)
			if !_continue {
				break
			}
		}
	}
}

type EncodedValue struct {
	Id    string
	Bytes []byte
}

func EncodedValues(ilst *ItemList) iter.Seq2[EncodedValue, error] {
	return func(yield func(EncodedValue, error) (_continue bool)) {
		rv := reflect.ValueOf(ilst).Elem()
		rt := rv.Type()

		for i := range make([]interface{}, rt.NumField()) {
			id := rt.Field(i).Tag.Get("id")
			value := rv.Field(i).Interface()
			buf, err := encodeFieldValue(value)
			if buf == nil {
				continue
			}
			_continue := yield(EncodedValue{id, buf}, err)
			if !_continue {
				break
			}
		}
	}
}

func encodeFieldValue(value any) ([]byte, error) {
	switch v := value.(type) {
	case *internationalText:
		if v == nil {
			return nil, nil
		}
		return v.Bytes()
	case *Genre:
		if v == nil {
			return nil, nil
		}
		return v.Bytes()
	case *BoolWithHeader0x15_0:
		if v == nil {
			return nil, nil
		}
		return v.Bytes(), nil
	case *Int16WithHeader0x15_0:
		if v == nil {
			return nil, nil
		}
		return v.Bytes(), nil
	case *TrackNumber:
		if v == nil {
			return nil, nil
		}
		return v.Bytes()
	case *DiskNumber:
		if v == nil {
			return nil, nil
		}
		return v.Bytes()
	default:
		panic("unsupported item type")
	}
}

type writableValue struct {
	id    string
	field reflect.Value
}

func (w writableValue) SetDecoded(buf []byte) (err error) {
	switch w.field.Interface().(type) {
	case *internationalText:
		err = setField(w.field, decodeInternationalText, buf)
	case *Genre:
		err = setField(w.field, decodeGenre, buf)
	case *BoolWithHeader0x15_0:
		err = setField(w.field, decodeBoolWithHeader0x15_0, buf)
	case *Int16WithHeader0x15_0:
		err = setField(w.field, decodeInt16WithHeader0x15_0, buf)
	case *TrackNumber:
		err = setField(w.field, decodeTrackNumber, buf)
	case *DiskNumber:
		err = setField(w.field, decodeDiskNumber, buf)
	default:
		panic("unsupported item type")
	}
	return err
}

func (w writableValue) Remove() {
	w.field.Set(reflect.ValueOf(nil))
}

func (w writableValue) GetDecorder() decoder {
	return decoder{w.field}
}

func IterateFieldWriters(ilst *ItemList) iter.Seq2[string, writableValue] {
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
	return writableValue{id, field}
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
