package rule

import (
	"fmt"
	"reflect"
)

type IPv6Rule struct {
	Interface
}

func (v IPv6Rule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	if !isIPv6(fmt.Sprintf("%v", value)) {
		return false, fmt.Errorf("must be a valid IPv6")
	}

	return true, nil
}
