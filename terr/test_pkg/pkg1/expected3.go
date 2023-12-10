package pkg1

import (
	"goutils/terr"
	"goutils/terr/test_pkg/pkg1/pkg2"
)

func WrapExpected3() error {
	//Just make
	//Just make
	//Just make
	//Just make
	//Just make
	//Just make Line 15
	return terr.Wrap(pkg2.WrapExpected2())
}

func NewExpected3() error {
	return terr.Wrap(pkg2.NewExpected2())
}

func JustReturn3() error {
	return pkg2.JustReturn2()
}
