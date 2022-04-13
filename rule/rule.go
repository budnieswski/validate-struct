package rule

import "reflect"

type Interface interface {
	Validate(value reflect.Value) (bool, error)
}
