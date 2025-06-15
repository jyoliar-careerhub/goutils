package async

type Result[T any] struct {
	Value T
	Err   error
}

func ExecAsync[T any](fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		defer close(ch)

		value, err := fn()
		ch <- Result[T]{value, err}
	}()

	return ch
}

func ExecAsyncWithParam[T any, P any](param P, fn func(P) (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		defer close(ch)

		value, err := fn(param)
		ch <- Result[T]{value, err}
	}()

	return ch
}
