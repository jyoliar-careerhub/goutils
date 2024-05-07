package async

type Result[T any] struct {
	Value T
	Err   error
}

func ExecAsync(fn func() (interface{}, error)) <-chan Result[interface{}] {
	ch := make(chan Result[interface{}])

	go func() {
		defer close(ch)

		value, err := fn()
		ch <- Result[interface{}]{value, err}
	}()

	return ch
}
