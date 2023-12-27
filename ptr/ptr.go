package ptr

import "reflect"

func P[T any](value T) *T {
	return &value
}

func IsNil(target any) bool {
	if target == nil {
		return true
	}

	value := reflect.ValueOf(target)
	return value.Kind() == reflect.Ptr && value.IsNil()
}
