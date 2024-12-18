package qtffilst

import "strings"

func ilstDataBox(box Box) bool {
	return strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.") &&
		!box.IsContainable &&
		box.Name == "data"
}

func ilstDataBoxName(path string) string {
	return strings.Split(path, ".")[5]
}
