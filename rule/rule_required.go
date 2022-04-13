package rule

import (
	"fmt"
	"reflect"
)

type RequiredRule struct {
	Interface
}

func (v RequiredRule) Validate(value reflect.Value) (bool, error) {
	if !isEmpty(value) {
		return true, nil
	}

	return false, fmt.Errorf("the field is required")
}
