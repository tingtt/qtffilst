package ilst

import (
	"errors"

	"github.com/tingtt/iterutil"
)

func (il *ItemList) SetDecoded(id string, value []byte) error {
	for _, v := range iterutil.FilterKey(IterateFieldWriters(il), id) {
		return v.SetDecoded(value)
	}
	return errors.New("field not found")
}
