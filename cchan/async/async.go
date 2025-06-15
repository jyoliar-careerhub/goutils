package async

import (
	"context"
)

type Result[T any] struct {
	Value T
	Err   error
}

func ExecAsync[T any](ctx context.Context, fn func(context.Context) (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		defer close(ch)

		value, err := fn(ctx)
		ch <- Result[T]{value, err}
	}()

	return ch
}

func ExecAsyncWithParam[T any, P any](ctx context.Context, param P, fn func(context.Context, P) (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		defer close(ch)

		value, err := fn(ctx, param)
		ch <- Result[T]{value, err}
	}()

	return ch
}
