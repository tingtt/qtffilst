package clioption

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/tingtt/iterutil"
	"github.com/tingtt/qtffilst/ilst"
)

func loadChanges(changeDatas, removeIds []string) (itemList *ilst.ItemList, deleteIds []string, err error) {
	newItemList := new(ilst.ItemList)

	for _, changeDataStr := range changeDatas {
		id, value, err := decodeChangeData(changeDataStr)
		if err != nil {
			return nil, nil, fmt.Errorf("CLI option `--data`,`-d` %w", err)
		}

		for _, v := range iterutil.FilterKey(ilst.IterateFieldWriters(newItemList), id) {
			decodedValue, err := v.GetDecorder().Decode(value)
			if err != nil {
				return nil, nil, err
			}
			err = newItemList.SetDecoded(id, decodedValue)
			if err != nil {
				return nil, nil, err
			}
			break
		}
	}

	for _, id := range removeIds {
		if !validItemListBoxId(id) {
			return nil, nil, fmt.Errorf("CLI option `--rm`,`-r` invalid ItemList id")
		}
	}

	return newItemList, removeIds, nil
}

func decodeChangeData(str string) (id, value string, err error) {
	l := strings.Split(str, "=")
	if len(l) != 2 {
		return "", "", errors.New("invalid format")
	}
	return l[0], strings.Trim(l[1], "\""), nil
}

func validItemListBoxId(id string) bool {
	return len(slices.Collect(iterutil.Filter(ilst.Ids(), id))) == 1
}
