package pkg2

import (
	"github.com/jae2274/goutils/terr"
	"github.com/jae2274/goutils/terr/test_pkg/pkg1/pkg2/pkg3"
)

func WrapExpected2() error {
	return terr.Wrap(pkg3.WrapExpected1())
}

func NewExpected2() error {
	return terr.Wrap(pkg3.NewExpected1())
}

func JustReturn2() error {
	return pkg3.NewExpected1()
}
