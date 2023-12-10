package pkg1

import "goutils/terr"

func WrapExpected4() error {
	return terr.Wrap(WrapExpected3())
}

func NewExpected4() error {
	return terr.Wrap(NewExpected3())
}

func Justreturn4() error {
	return JustReturn3()
}
