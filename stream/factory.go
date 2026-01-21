package stream

func Of[T any](items ...T) Stream[T] {
	return newStream(items)
}

func OfSlice[T any](items []T) Stream[T] {
	result := make([]T, len(items))
	copy(result, items)
	return newStream(result)
}

func Range(startInclusive, endExclusive int64) Stream[int64] {
	size := endExclusive - startInclusive
	if size <= 0 {
		return newStream([]int64{})
	}
	items := make([]int64, size)
	for i := int64(0); i < size; i++ {
		items[i] = startInclusive + i
	}
	return newStream(items)
}

func RangeClosed(startInclusive, endInclusive int64) Stream[int64] {
	size := endInclusive - startInclusive + 1
	if size <= 0 {
		return newStream([]int64{})
	}
	items := make([]int64, size)
	for i := int64(0); i < size; i++ {
		items[i] = startInclusive + i
	}
	return newStream(items)
}

func Empty[T any]() Stream[T] {
	return newStream([]T{})
}

func Generate[T any](supplier Supplier[T], count int) Stream[T] {
	items := make([]T, count)
	for i := 0; i < count; i++ {
		items[i] = supplier()
	}
	return newStream(items)
}

func Concat[T any](streams ...Stream[T]) Stream[T] {
	var allItems []T
	for _, s := range streams {
		allItems = append(allItems, s.ToSlice()...)
	}
	return newStream(allItems)
}
