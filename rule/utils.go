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

	if x.Kind() == reflect.Invalid {
		return true
	}

	switch x.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return x.Len() == 0
	}

	return reflect.DeepEqual(x, reflect.Zero(x.Type())) || strings.HasPrefix(x.String(), "0001-01-01")
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
