package rule

import (
	"fmt"
	"reflect"
	"strings"
)

type ENumRule struct {
	Interface
	ValidValues []string
}

func (v ENumRule) Validate(value reflect.Value) (bool, error) {
	if isEmpty(value) || len(v.ValidValues) < 1 {
		return true, nil
	}

	for _, vv := range v.ValidValues {
		if strings.Compare(vv, fmt.Sprintf("%v", value)) == 0 {
			return true, nil
		}
	}

	return false, fmt.Errorf("must be one of these: %s", strings.Join(v.ValidValues, ", "))
}
