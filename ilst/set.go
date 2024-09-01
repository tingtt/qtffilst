package ilst

func (il *ItemList) Set(id string, value []byte) (err error) {
	for _, v := range IdWriters(il, WithSomeId(id)) {
		err := v.set(value)
		if err != nil {
			return err
		}
	}
	return err
}
