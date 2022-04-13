package rule

import (
	"reflect"
	"regexp"
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
	regexIP = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
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

func isIP(ip string) bool {
	return regexIP.MatchString(ip)
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
