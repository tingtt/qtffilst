package ilst

import (
	"reflect"
	"strconv"
	"strings"
)

type decoder struct {
	targetField reflect.Value
}

func (d decoder) Decode(str string) ([]byte, error) {
	switch d.targetField.Interface().(type) {
	case *internationalText:
		return NewInternationalText(str).Bytes()
	case *Genre:
		// TODO: implement decode to Genre
		panic("unimplemented: decode to Genre")
	case *BoolWithHeader0x15_0:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, err
		}
		return (&BoolWithHeader0x15_0{b}).Bytes(), nil
	case *Int16WithHeader0x15_0:
		i, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return nil, err
		}
		return (&Int16WithHeader0x15_0{int16(i)}).Bytes(), nil
	case *TrackNumber:
		number, total, err := decodeSlashedStr(str)
		if err != nil {
			return nil, err
		}
		if number == 0 {
			number = 1
		}
		if total == 0 {
			total = number
		}
		return (&TrackNumber{int16(number), int16(total)}).Bytes()
	case *DiskNumber:
		number, total, err := decodeSlashedStr(str)
		if err != nil {
			return nil, err
		}
		if number == 0 {
			number = 1
		}
		if total == 0 {
			total = number
		}
		return (&TrackNumber{int16(number), int16(total)}).Bytes()
	default:
		panic("unsupported item type")
	}
}

func decodeSlashedStr(str string) (number, total int, err error) {
	splittedStr := strings.Split(str, "/")

	numberStr := splittedStr[0]
	number, err = strconv.Atoi(numberStr)
	if err != nil {
		return 0, 0, err
	}

	if len(splittedStr) < 2 {
		return number, 0, nil
	}

	totalStr := splittedStr[1]
	total, err = strconv.Atoi(totalStr)
	if err != nil {
		return 0, 0, err
	}

	return number, total, nil
}
