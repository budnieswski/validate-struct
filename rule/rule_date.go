package rule

import (
	"fmt"
	"reflect"
	"strings"
)

type DateRule struct {
	Interface
	Format string
}

func (v DateRule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	switch v.Format {
	case "dd-mm-yyyy":
		if !isDateDDMMYY(fmt.Sprintf("%v", value)) {
			return false, fmt.Errorf("must be a valid date format. e.g: %s", strings.Join(DateLayoutDDMMYY[:], ", "))
		}
	case "timestamp":
		if !isTimestamp(fmt.Sprintf("%v", value)) {
			return false, fmt.Errorf("must be a valid timestamp format. e.g: %s", strings.Join(DateLayoutTimestamp[:], ", "))
		}
	default:
		if !isDate(fmt.Sprintf("%v", value)) {
			return false, fmt.Errorf("must be a valid date format. e.g: %s", strings.Join(DateLayoutDDMMYY[:], ", "))
		}
	}

	return true, nil
}
