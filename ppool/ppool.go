package ppool

import "context"

type token struct{}
type Pool[T any] struct {
	sem  chan token
	idle chan T

	newFunc func() (T, error)
}

func NewPool[T any](size int, newFunc func() (T, error)) *Pool[T] {
	return &Pool[T]{
		sem:     make(chan token, size),
		idle:    make(chan T, size),
		newFunc: newFunc,
	}
}

func (p *Pool[T]) Release(value T) {
	p.idle <- value
}

func (p *Pool[T]) Hijack() {
	<-p.sem
}

func (p *Pool[T]) Acquire(ctx context.Context) (T, error) {
	select {
	case value := <-p.idle:
		return value, nil
	case p.sem <- token{}:
		newV, err := p.newFunc()
		if err != nil {
			<-p.sem
		}

		return newV, err
	case <-ctx.Done():
		return *new(T), ctx.Err()
	}
}
