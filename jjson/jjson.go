package jjson

import (
	"encoding/json"
	"io"
)

func Unmarshal[T any](data []byte) (*T, error) {
	value := new(T)
	err := json.Unmarshal(data, value)
	return value, err
}

func UnmarshalReader[T any](reader io.Reader) (*T, error) {
	value := new(T)
	err := json.NewDecoder(reader).Decode(value)
	return value, err
}
