package ilst

import (
	"errors"

	"github.com/tingtt/iterutil"
)

func (il *ItemList) Set(id string, value []byte) error {
	for _, v := range iterutil.FilterKey(IdWriters(il), id) {
		return v.set(value)
	}
	return errors.New("field not found")
}
