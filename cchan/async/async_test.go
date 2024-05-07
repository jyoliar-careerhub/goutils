package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAsync(t *testing.T) {

	// given
	waitSecond3Return1 := func() (interface{}, error) {
		time.Sleep(3 * time.Second)
		return 1, nil
	}

	t.Run("run function asynchronously", func(t *testing.T) {

		// when
		start := time.Now()
		ExecAsync(waitSecond3Return1)
		returnAsync := time.Now()
		require.Less(t, returnAsync.Sub(start), 1*time.Millisecond)
	})
	t.Run("chan Result returns after finish func", func(t *testing.T) {
		// when
		start := time.Now()
		result := ExecAsync(waitSecond3Return1)

		resultValue := <-result
		end := time.Now()
		require.GreaterOrEqual(t, end.Sub(start), 3*time.Second)

		// then
		require.Equal(t, 1, resultValue.Value)
		require.Nil(t, resultValue.Err)

		select {
		case _, ok := <-result:
			require.False(t, ok)
		default:
			require.Fail(t, "result channel should be closed")
		}
	})
}
