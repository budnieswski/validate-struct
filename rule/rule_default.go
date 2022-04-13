package rule

import "reflect"

type DefaultRule struct {
	Interface
}

func (v DefaultRule) Validate(value reflect.Value) (bool, error) {
	return true, nil
}
