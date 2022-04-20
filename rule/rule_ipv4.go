package rule

import (
	"fmt"
	"reflect"
)

type IPv4Rule struct {
	Interface
}

func (v IPv4Rule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	if !isIPv4(fmt.Sprintf("%v", value)) {
		return false, fmt.Errorf("must be a valid IPv4")
	}

	return true, nil
}
