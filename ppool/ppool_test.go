package ppool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	require "github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	t.Run("Acquire and Release", testPoolAcquireAndRelease)
	t.Run("Acquire with Timeout", testPoolAcquireWithTimeout)
	t.Run("NewFunc Error", testPoolNewFuncError)
	t.Run("Hijack", testPoolHijack)
	t.Run("Limit Pool Size", testPoolLimitSize)

}

func testPoolAcquireAndRelease(t *testing.T) {
	var counter int32

	newFunc := func() (int, error) {
		return int(atomic.AddInt32(&counter, 1)), nil
	}

	pool := NewPool[int](2, newFunc)

	ctx := context.Background()

	// Acquire two values (pool size is 2)
	v1, err := pool.Acquire(ctx)
	require.NoError(t, err)

	v2, err := pool.Acquire(ctx)
	require.NoError(t, err)

	require.NotEqual(t, v1, v2, "expected different values")

	// Release one value back to pool
	pool.Release(v1)

	// Acquire again, should reuse the released value
	v3, err := pool.Acquire(ctx)
	require.NoError(t, err)

	require.Equal(t, v1, v3, "expected reused value to match released value")
}

func testPoolAcquireWithTimeout(t *testing.T) {
	newFunc := func() (int, error) {
		return 1, nil
	}
	pool := NewPool[int](1, newFunc)

	ctx := context.Background()

	// Fill the pool
	_, err := pool.Acquire(ctx)
	require.NoError(t, err)

	// Now attempt to acquire with short timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = pool.Acquire(timeoutCtx)
	require.Error(t, err, "expected error due to timeout")
}

func testPoolNewFuncError(t *testing.T) {
	errExpected := errors.New("creation error")

	pool := NewPool[int](1, func() (int, error) {
		return 0, errExpected
	})

	ctx := context.Background()

	_, err := pool.Acquire(ctx)
	require.ErrorIs(t, err, errExpected, "expected error from newFunc")
}

func testPoolHijack(t *testing.T) {
	newFunc := func() (int, error) {
		return 1, nil
	}

	pool := NewPool[int](1, newFunc)

	ctx := context.Background()

	// Acquire to fill the semaphore
	_, err := pool.Acquire(ctx)
	require.NoError(t, err, "expected no error when acquiring from pool")

	done := make(chan int)

	go func() {
		pool.Hijack()
		value, err := pool.Acquire(ctx)
		require.NoError(t, err, "expected no error when hijacking and acquiring from pool")
		done <- value
	}()

	select {
	case <-done:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Error("Hijack did not unblock as expected")
	}
}

func testPoolLimitSize(t *testing.T) {
	newFunc := func() (int, error) {
		return 1, nil
	}

	pool := NewPool[int](2, newFunc)

	ctx := context.Background()

	// Acquire two values (pool size is 2)
	_, err := pool.Acquire(ctx)
	require.NoError(t, err)

	_, err = pool.Acquire(ctx)
	require.NoError(t, err)

	// Attempt to acquire a third value should block or return an error
	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	_, err = pool.Acquire(timeoutCtx)
	require.Error(t, err, "expected error due to pool limit")
}
