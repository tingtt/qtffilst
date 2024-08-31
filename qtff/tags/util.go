package tags

import "strings"

func ilstDataBox(box Box) bool {
	return strings.HasPrefix(box.Path, ".moov.udta.meta.ilst.") &&
		box.Id == "data"
}

func ilstDataBoxName(path string) string {
	return strings.Split(path, ".")[5]
}
