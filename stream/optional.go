package stream

type Optional[T any] struct {
	value *T
}

func OfOptional[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

func OfOptionalPtr[T any](value *T) Optional[T] {
	return Optional[T]{value: value}
}

func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{value: nil}
}

func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

func (o Optional[T]) IsEmpty() bool {
	return o.value == nil
}

func (o Optional[T]) Get() T {
	if o.value == nil {
		var zero T
		return zero
	}
	return *o.value
}

func (o Optional[T]) OrElse(other T) T {
	if o.value == nil {
		return other
	}
	return *o.value
}

func (o Optional[T]) OrElseGet(supplier Supplier[T]) T {
	if o.value == nil {
		return supplier()
	}
	return *o.value
}

func (o Optional[T]) OrElsePanic(msg string) T {
	if o.value == nil {
		panic(msg)
	}
	return *o.value
}

func (o Optional[T]) IfPresent(consumer Consumer[T]) {
	if o.value != nil {
		consumer(*o.value)
	}
}

func (o Optional[T]) IfPresentOrElse(consumer Consumer[T], action func()) {
	if o.value != nil {
		consumer(*o.value)
	} else {
		action()
	}
}

func (o Optional[T]) Filter(predicate Predicate[T]) Optional[T] {
	if o.value == nil || !predicate(*o.value) {
		return EmptyOptional[T]()
	}
	return o
}
