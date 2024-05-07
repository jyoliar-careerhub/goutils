package async

type Result[T any] struct {
	Value T
	Err   error
}

func ExecAsync[T any](fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T])

	go func() {
		defer close(ch)

		value, err := fn()
		ch <- Result[T]{value, err}
	}()

	return ch
}
