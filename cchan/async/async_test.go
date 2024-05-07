package async

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testSample struct {
	name string
}

func TestAsync(t *testing.T) {

	// given
	waitSecond3Return1 := func() (int, error) {
		time.Sleep(3 * time.Second)
		return 1, nil
	}

	returnInt := func() (int, error) {
		return 1, nil
	}

	returnString := func() (string, error) {
		return "hello", nil
	}

	returnStruct := func() (testSample, error) {
		return testSample{name: "test"}, nil
	}

	returnError := func() (any, error) {
		return nil, fmt.Errorf("error")
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

	t.Run("return variable values", func(t *testing.T) {
		t.Run("return int", func(t *testing.T) {
			// when
			result := ExecAsync(returnInt)

			// then
			resultValue := <-result
			require.Equal(t, 1, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return string", func(t *testing.T) {
			// when
			result := ExecAsync(returnString)

			// then
			resultValue := <-result
			require.Equal(t, "hello", resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return struct", func(t *testing.T) {
			// when
			result := ExecAsync(returnStruct)

			// then
			resultValue := <-result
			require.Equal(t, testSample{name: "test"}, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return error", func(t *testing.T) {
			// when
			result := ExecAsync(returnError)

			// then
			resultValue := <-result
			require.Nil(t, resultValue.Value)
			require.Error(t, resultValue.Err)
		})
	})
}
