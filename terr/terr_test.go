package terr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jae2274/goutils/terr"
	"github.com/jae2274/goutils/terr/test_pkg/pkg1"
	"github.com/jae2274/goutils/terr/test_pkg/pkg1/pkg2/pkg3"
	"github.com/stretchr/testify/require"
)

func TestErrStack(t *testing.T) {
	t.Run("Start Wrapping error", func(t *testing.T) {
		err := pkg1.WrapExpected4()

		traceErr, ok := err.(*terr.TraceError)

		require.True(t, ok)
		frames := traceErr.Frames()

		require.Equal(t, "expected1.go", frames[0].File)
		require.Equal(t, "expected2.go", frames[1].File)
		require.Equal(t, "expected3.go", frames[2].File)
		require.Equal(t, "expected4.go", frames[3].File)

		lines := []int{16, 9, 15, 6}
		require.Equal(t, lines[0], frames[0].Line)
		require.Equal(t, lines[1], frames[1].Line)
		require.Equal(t, lines[2], frames[2].Line)
		require.Equal(t, lines[3], frames[3].Line)

		require.True(t,
			strings.HasPrefix(err.Error(), fmt.Sprintf("ExampleError\tStackTrace: expected1.go:%d -> expected2.go:%d -> expected3.go:%d -> expected4.go:%d", lines[0], lines[1], lines[2], lines[3])),
		)

		_, ok = err.(*pkg3.ExampleError)
		require.False(t, ok)

		err = terr.UnWrap(err)
		_, ok = err.(*pkg3.ExampleError)
		require.True(t, ok)
	})

	t.Run("Start TraceError", func(t *testing.T) {
		err := pkg1.NewExpected4()

		traceErr, ok := err.(*terr.TraceError)

		require.True(t, ok)
		frames := traceErr.Frames()
		// require.Equal(t, 4, len(frames))

		require.Equal(t, "expected1.go", frames[0].File)
		require.Equal(t, "expected2.go", frames[1].File)
		require.Equal(t, "expected3.go", frames[2].File)
		require.Equal(t, "expected4.go", frames[3].File)

		lines := []int{24, 13, 19, 10}
		require.Equal(t, lines[0], frames[0].Line)
		require.Equal(t, lines[1], frames[1].Line)
		require.Equal(t, lines[2], frames[2].Line)
		require.Equal(t, lines[3], frames[3].Line)

		require.True(t,
			strings.HasPrefix(err.Error(), fmt.Sprintf("VariableError\tStackTrace: expected1.go:%d -> expected2.go:%d -> expected3.go:%d -> expected4.go:%d", lines[0], lines[1], lines[2], lines[3])),
		)

		require.Equal(t, pkg3.ErrVariable, terr.UnWrap(err))
	})

	// terr.Wrap의 메소드는 인자로 전달되는 err가 TraceError인지 확인하고, TraceError가 아니라면 TraceError로 감싸서 반환, TraceError라면 그대로 반환한다.
	// 따라서 terr.Wrap을 호출한 위의 두 테스트와 동일한 동작을 수행한다.
	t.Run("Without Wrapped, Just return error ", func(t *testing.T) {
		err := pkg1.Justreturn4()

		traceErr, ok := err.(*terr.TraceError)

		require.True(t, ok)
		frames := traceErr.Frames()
		// require.Equal(t, 4, len(frames))

		require.Equal(t, "expected1.go", frames[0].File)
		require.Equal(t, "expected2.go", frames[1].File)
		require.Equal(t, "expected3.go", frames[2].File)
		require.Equal(t, "expected4.go", frames[3].File)

		lines := []int{24, 17, 23, 14}
		require.Equal(t, lines[0], frames[0].Line)
		require.Equal(t, lines[1], frames[1].Line)
		require.Equal(t, lines[2], frames[2].Line)
		require.Equal(t, lines[3], frames[3].Line)

		require.True(t,
			strings.HasPrefix(err.Error(), fmt.Sprintf("VariableError\tStackTrace: expected1.go:%d -> expected2.go:%d -> expected3.go:%d -> expected4.go:%d", lines[0], lines[1], lines[2], lines[3])),
		)

		require.Equal(t, pkg3.ErrVariable, terr.UnWrap(err))
	})
}
