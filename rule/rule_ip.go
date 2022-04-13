package rule

import (
	"fmt"
	"reflect"
)

type IPRule struct {
	Interface
}

func (v IPRule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	if !isIP(fmt.Sprintf("%v", value)) {
		return false, fmt.Errorf("must be a valid IP address")
	}

	return true, nil
}
