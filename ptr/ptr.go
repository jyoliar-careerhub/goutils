package ptr

import "reflect"

func P[T any](value T) *T {
	return &value
}

func IsNil(target any) bool { //주로 generic type의 파라미터를 전달받는 함수에서 해당 파라미터가 nil인지 확인할 때 사용한다.
	if target == nil {
		return true
	}

	value := reflect.ValueOf(target)
	return value.Kind() == reflect.Ptr && value.IsNil()
}
