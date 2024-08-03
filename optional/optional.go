package optional

// Optional은 제네릭을 사용하여 다양한 타입 T에 대해 동작할 수 있는 구조체입니다.
type Optional[T any] struct {
	value *T
}

// NewOptional은 주어진 값을 포함하는 Optional을 생성합니다.
func NewOptional[T any](value *T) Optional[T] {
	return Optional[T]{value: value}
}

func NewOptionalPtr[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

// IsPresent는 Optional에 값이 설정되어 있는지 여부를 반환합니다.
func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o Optional[T]) GetPtr() *T {
	return o.value
}

func (o Optional[T]) Get() T {
	if o.value == nil {
		return *new(T)
	}
	return *o.value
}

// OrElsePtr는 Optional에 값이 있으면 그 값을,
// 없으면 기본값을 반환합니다.
func (o Optional[T]) OrElsePtr(defaultValue T) *T {
	if o.value == nil {
		return &defaultValue
	}
	return o.value
}

func (o Optional[T]) OrElse(defaultValue T) T {
	if o.value == nil {
		return defaultValue
	}
	return *o.value
}
