package pkg3

import (
	"errors"
	"goutils/terr"
)

type ExampleError struct{}

func (e *ExampleError) Error() string {
	return "ExampleError"
}

func WrapExpected1() error {
	return terr.Wrap(&ExampleError{})
}

var (
	ErrVariable = errors.New("VariableError") //다음과 같이 traceError를 초기화시켜 재사용할 경우 애플리케이션이 로드되는 시점의 stacktrace가 저장되어 버린다.
)

func NewExpected1() error {
	return terr.Wrap(ErrVariable)
}
