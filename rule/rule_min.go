package rule

import (
	"fmt"
	"reflect"
)

type MinRule struct {
	Interface
	Min int
}

func (v MinRule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Map {
		return v.validateArray(value)
	}

	if len(fmt.Sprintf("%v", value)) < v.Min {
		return false, fmt.Errorf("must be at least %d characters long", v.Min)
	}

	return true, nil
}

func (v MinRule) validateArray(value reflect.Value) (bool, error) {
	if value.Len() < v.Min {
		return false, fmt.Errorf("must have at least %d items", v.Min)
	}

	return true, nil
}
