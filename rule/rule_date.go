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

	// v.Format = "dd-mm-yyyy"
	// value = reflect.ValueOf("01-04-2022")
	// value = "01/04/2022"

	// v.Format = "timestamp"
	// value = "2022-04-06T17:13:22Z"
	// value = "2022-04-06T17:13:22+03:00"
	// value = "2022-04-06T17:13:22-03:00"
	// value = "2022-04-06T17:13:22.508Z"
	// value = "2022-04-06T17:13:22.508+03:00"
	// value = "2022-04-06T17:13:22.508-03:00"

	// value = "1999-04-01"

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
