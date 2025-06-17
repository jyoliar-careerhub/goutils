package async

import (
	"context"
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
	waitSecond3Return1 := func(ctx context.Context) (int, error) {
		select {
		case <-time.After(3 * time.Second):
			return 1, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	returnInt := func(context.Context) (int, error) {
		return 1, nil
	}

	returnString := func(context.Context) (string, error) {
		return "hello", nil
	}

	returnStruct := func(context.Context) (testSample, error) {
		return testSample{name: "test"}, nil
	}

	returnError := func(context.Context) (any, error) {
		return nil, fmt.Errorf("error")
	}

	t.Run("run function asynchronously", func(t *testing.T) {

		// when
		start := time.Now()
		ExecAsync(context.Background(), waitSecond3Return1)
		returnAsync := time.Now()
		require.Less(t, returnAsync.Sub(start), 1*time.Millisecond)
	})
	t.Run("chan Result returns after finish func", func(t *testing.T) {
		// when
		start := time.Now()
		result := ExecAsync(context.Background(), waitSecond3Return1)

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
			result := ExecAsync(context.Background(), returnInt)

			// then
			resultValue := <-result
			require.Equal(t, 1, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return string", func(t *testing.T) {
			// when
			result := ExecAsync(context.Background(), returnString)

			// then
			resultValue := <-result
			require.Equal(t, "hello", resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return struct", func(t *testing.T) {
			// when
			result := ExecAsync(context.Background(), returnStruct)

			// then
			resultValue := <-result
			require.Equal(t, testSample{name: "test"}, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return error", func(t *testing.T) {
			// when
			result := ExecAsync(context.Background(), returnError)

			// then
			resultValue := <-result
			require.Nil(t, resultValue.Value)
			require.Error(t, resultValue.Err)
		})
	})

	t.Run("context cancelled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())

		// when
		result := ExecAsync(ctx, waitSecond3Return1)

		cancel() // cancel context before function finishes

		resultValue := <-result

		require.Less(t, time.Millisecond, 3*time.Second)
		require.ErrorIs(t, resultValue.Err, context.Canceled)

		select {
		case _, ok := <-result:
			require.False(t, ok)
		default:
			require.Fail(t, "result channel should be closed")
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		// given
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// when
		start := time.Now()
		result := ExecAsync(ctx, waitSecond3Return1)

		resultValue := <-result

		require.WithinDuration(t, start.Add(time.Second), time.Now(), 10*time.Millisecond)
		require.Equal(t, context.DeadlineExceeded, resultValue.Err)

		select {
		case _, ok := <-result:
			require.False(t, ok)
		default:
			require.Fail(t, "result channel should be closed")
		}
	})
}
func TestAsyncWithParam(t *testing.T) {

	// given
	waitSecond3Return1 := func(ctx context.Context, p int) (int, error) {
		select {
		case <-time.After(3 * time.Second):
			return p, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	returnInt := func(ctx context.Context, p int) (int, error) {
		return p, nil
	}

	returnString := func(ctx context.Context, p string) (string, error) {
		return p, nil
	}

	returnStruct := func(ctx context.Context, p string) (testSample, error) {
		return testSample{name: p}, nil
	}

	returnError := func(ctx context.Context, p string) (any, error) {
		return nil, fmt.Errorf(p)
	}

	t.Run("run function asynchronously", func(t *testing.T) {

		// when
		start := time.Now()
		ExecAsyncWithParam(context.Background(), 1, waitSecond3Return1)
		returnAsync := time.Now()
		require.Less(t, returnAsync.Sub(start), 1*time.Millisecond)
	})
	t.Run("chan Result returns after finish func", func(t *testing.T) {
		// when
		start := time.Now()
		result := ExecAsyncWithParam(context.Background(), 1, waitSecond3Return1)

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
			result := ExecAsyncWithParam(context.Background(), 1, returnInt)

			// then
			resultValue := <-result
			require.Equal(t, 1, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return string", func(t *testing.T) {
			// when
			result := ExecAsyncWithParam(context.Background(), "hello", returnString)

			// then
			resultValue := <-result
			require.Equal(t, "hello", resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return struct", func(t *testing.T) {
			// when
			result := ExecAsyncWithParam(context.Background(), "test", returnStruct)

			// then
			resultValue := <-result
			require.Equal(t, testSample{name: "test"}, resultValue.Value)
			require.Nil(t, resultValue.Err)
		})

		t.Run("return error", func(t *testing.T) {
			// when
			result := ExecAsyncWithParam(context.Background(), "errorTest", returnError)

			// then
			resultValue := <-result
			require.Nil(t, resultValue.Value)
			require.Error(t, resultValue.Err)
			require.Equal(t, "errorTest", resultValue.Err.Error())
		})
	})

	t.Run("context cancelled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())

		// when
		result := ExecAsyncWithParam(ctx, 1, waitSecond3Return1)

		cancel() // cancel context before function finishes

		resultValue := <-result

		require.Less(t, time.Millisecond, 3*time.Second)
		require.ErrorIs(t, resultValue.Err, context.Canceled)

		select {
		case _, ok := <-result:
			require.False(t, ok)
		default:
			require.Fail(t, "result channel should be closed")
		}
	})

	t.Run("context timeout", func(t *testing.T) {
		// given
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// when
		start := time.Now()
		result := ExecAsyncWithParam(ctx, 1, waitSecond3Return1)

		resultValue := <-result

		require.WithinDuration(t, start.Add(time.Second), time.Now(), 10*time.Millisecond)
		require.Equal(t, context.DeadlineExceeded, resultValue.Err)

		select {
		case _, ok := <-result:
			require.False(t, ok)
		default:
			require.Fail(t, "result channel should be closed")
		}
	})
}
