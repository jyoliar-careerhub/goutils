package ptr

func P[T any](value T) *T {
	return &value
}
