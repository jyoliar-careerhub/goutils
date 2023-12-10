package enum

import (
	"encoding/json"
	"reflect"
)

type Values interface {
	Values() []string
}

type Enum[T Values] string

func (e Enum[T]) MarshalText() (text []byte, err error) {

	enums := *new(T)
	for _, enumValue := range enums.Values() {
		if enumValue == string(e) {
			return []byte(e), nil
		}
	}

	return nil, &json.UnsupportedValueError{
		Value: reflect.ValueOf(e),
		Str:   string(e),
	}
}

func (e *Enum[T]) UnmarshalText(text []byte) error {
	enums := *new(T)
	parseText := string(text)

	for _, enumValue := range enums.Values() {
		if enumValue == parseText {
			*e = Enum[T](enumValue)
			return nil
		}
	}

	return &json.UnsupportedValueError{
		Value: reflect.ValueOf(e),
		Str:   string(text),
	}
}
