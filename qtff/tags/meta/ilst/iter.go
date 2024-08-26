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
