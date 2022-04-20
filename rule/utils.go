package rule

import (
	"net"
	"reflect"
	"strings"
	"time"
)

var (
	DateLayoutTimestamp = [3]string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.999999999-07:00",
	}
	DateLayout = [1]string{
		"2006-01-02",
	}
	DateLayoutDDMMYY = [2]string{
		"02-01-2006",
		"02/01/2006",
	}
)

func isEmpty(x reflect.Value) bool {
	switch x.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String, reflect.Array:
		return x.Len() == 0
	case reflect.Map, reflect.Slice:
		return x.Len() == 0 || x.IsNil()
	case reflect.Bool:
		return !x.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return x.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return x.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return x.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return x.IsNil()
	}

	return reflect.DeepEqual(x.Interface(), reflect.Zero(x.Type()).Interface()) || strings.HasPrefix(x.String(), "0001-01-01")
}

// checks if is either IP version 4 or 6
func isIP(str string) bool {
	return net.ParseIP(str) != nil
}

func isIPv4(str string) bool {
	return net.ParseIP(str) != nil && strings.Contains(str, ".")
}

func isIPv6(str string) bool {
	return net.ParseIP(str) != nil && strings.Contains(str, ":")
}

func isDate(date string) bool {
	return dateValidate(date, DateLayout[:])
}

func isDateDDMMYY(date string) bool {
	return dateValidate(date, DateLayoutDDMMYY[:])
}

func isTimestamp(date string) bool {
	return dateValidate(date, DateLayoutTimestamp[:])
}

func dateValidate(date string, layout []string) bool {
	var failed bool

	for _, l := range layout {
		if _, err := time.Parse(l, date); err == nil {
			failed = false
			break
		}

		failed = true
	}

	return !failed
}
