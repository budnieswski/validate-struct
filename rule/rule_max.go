package rule

import (
	"fmt"
	"reflect"
)

type MaxRule struct {
	Interface
	Max int
}

func (v MaxRule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) {
		return true, nil
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Map {
		return v.validateArray(value)
	}

	if len(fmt.Sprintf("%v", value)) > v.Max {
		return false, fmt.Errorf("must be at most %d characters long", v.Max)
	}

	return true, nil
}

func (v MaxRule) validateArray(value reflect.Value) (bool, error) {
	if value.Len() > v.Max {
		return false, fmt.Errorf("must have at most %d items", v.Max)
	}

	return true, nil
}
