package optional

import (
	"testing"

	"github.com/jae2274/goutils/ptr"
	"github.com/stretchr/testify/require"
)

func testOptional[T any](t *testing.T, value *T, defaultValue T, zeroValue T) {
	t.Run("if exists", func(t *testing.T) {
		testExist(t, value, defaultValue, zeroValue)
	})

	t.Run("if not exists", func(t *testing.T) {
		testNotExist(t, nil, defaultValue, zeroValue)
	})
}

func testExist[T any](t *testing.T, value *T, defaultValue T, zeroValue T) {
	opt := NewOptional(value)
	require.True(t, opt.IsPresent())

	result := opt.Get()
	require.Equal(t, *value, result)

	resultPtr := opt.GetPtr()
	require.Equal(t, *value, *resultPtr)

	require.Equal(t, *value, *opt.OrElsePtr(defaultValue))
	require.Equal(t, *value, opt.OrElse(defaultValue))
}

func testNotExist[T any](t *testing.T, value *T, defaultValue T, zeroValue T) {
	opt := NewOptional(value)
	require.False(t, opt.IsPresent())

	result := opt.Get()
	require.Equal(t, zeroValue, result)

	resultPtr := opt.GetPtr()
	require.Nil(t, resultPtr)

	require.Equal(t, defaultValue, *opt.OrElsePtr(defaultValue))
	require.Equal(t, defaultValue, opt.OrElse(defaultValue))
}

func TestOptional(t *testing.T) {

	testOptional(t, ptr.P(42), 5, 0)
	testOptional(t, ptr.P("Hello, Optional!"), "Default", "")

	type User struct {
		Name string
		Age  int
	}
	testOptional(t, ptr.P(User{Name: "Alice", Age: 25}), User{Name: "Bob", Age: 30}, User{})
}
