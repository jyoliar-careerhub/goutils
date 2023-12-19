package pkg1

import "github.com/jae2274/goutils/terr"

func WrapExpected4() error {
	return terr.Wrap(WrapExpected3())
}

func NewExpected4() error {
	return terr.Wrap(NewExpected3())
}

func Justreturn4() error {
	return JustReturn3()
}
