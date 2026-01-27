package stream

type Stream[T any] interface {
	Filter(predicate Predicate[T]) Stream[T]
	Map(mapper Function[T, T]) Stream[T]
	FlatMap(mapper Function[T, Stream[T]]) Stream[T]
	Distinct() Stream[T]
	Sorted(comparator Comparator[T]) Stream[T]
	Limit(maxSize int64) Stream[T]
	Skip(n int64) Stream[T]
	Peek(consumer Consumer[T]) Stream[T]

	ForEach(consumer Consumer[T])
	Collect(collector Collector[T, any, any]) any
	Reduce(identity T, accumulator BinaryOperator[T]) T
	Count() int64
	AnyMatch(predicate Predicate[T]) bool
	AllMatch(predicate Predicate[T]) bool
	NoneMatch(predicate Predicate[T]) bool
	FindFirst() Optional[T]
	FindAny() Optional[T]
	ToSlice() []T

	execute() []T
}

type streamImpl[T any] struct {
	source     []T
	operations []func([]T) []T
	isConsumed bool
}

func newStream[T any](source []T) Stream[T] {
	return &streamImpl[T]{
		source:     source,
		operations: make([]func([]T) []T, 0),
	}
}

func (s *streamImpl[T]) Filter(predicate Predicate[T]) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		// 预分配切片容量，最坏情况与原切片大小相同
		result := make([]T, 0, len(items))
		for _, item := range items {
			if predicate(item) {
				result = append(result, item)
			}
		}
		return result
	})
	return s
}

func (s *streamImpl[T]) Map(mapper Function[T, T]) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		result := make([]T, len(items))
		for i, item := range items {
			result[i] = mapper(item)
		}
		return result
	})
	return s
}

func (s *streamImpl[T]) FlatMap(mapper Function[T, Stream[T]]) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		// 预分配切片容量，基于原切片大小估计
		result := make([]T, 0, len(items))
		for _, item := range items {
			stream := mapper(item)
			result = append(result, stream.ToSlice()...)
		}
		return result
	})
	return s
}

func (s *streamImpl[T]) Distinct() Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		seen := make(map[any]bool)
		result := make([]T, 0)
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
		return result
	})
	return s
}

func (s *streamImpl[T]) Sorted(comparator Comparator[T]) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		result := make([]T, len(items))
		copy(result, items)
		sortSlice(result, comparator)
		return result
	})
	return s
}

func sortSlice[T any](slice []T, comparator Comparator[T]) {
	quickSort(slice, 0, len(slice)-1, comparator)
}

func quickSort[T any](slice []T, low, high int, comparator Comparator[T]) {
	if low < high {
		pivotIndex := partition(slice, low, high, comparator)
		quickSort(slice, low, pivotIndex-1, comparator)
		quickSort(slice, pivotIndex+1, high, comparator)
	}
}

func partition[T any](slice []T, low, high int, comparator Comparator[T]) int {
	pivot := slice[high]
	i := low - 1

	for j := low; j < high; j++ {
		if comparator(slice[j], pivot) <= 0 {
			i++
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

	slice[i+1], slice[high] = slice[high], slice[i+1]
	return i + 1
}

func (s *streamImpl[T]) Limit(maxSize int64) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		if int64(len(items)) <= maxSize {
			return items
		}
		return items[:maxSize]
	})
	return s
}

func (s *streamImpl[T]) Skip(n int64) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		if int64(len(items)) <= n {
			return []T{}
		}
		return items[n:]
	})
	return s
}

func (s *streamImpl[T]) Peek(consumer Consumer[T]) Stream[T] {
	s.checkNotConsumed()
	s.operations = append(s.operations, func(items []T) []T {
		for _, item := range items {
			consumer(item)
		}
		return items
	})
	return s
}

func (s *streamImpl[T]) ForEach(consumer Consumer[T]) {
	items := s.execute()
	for _, item := range items {
		consumer(item)
	}
}

func (s *streamImpl[T]) Collect(collector Collector[T, any, any]) any {
	items := s.execute()
	return collector.Collect(items)
}

func (s *streamImpl[T]) Reduce(identity T, accumulator BinaryOperator[T]) T {
	items := s.execute()
	result := identity
	for _, item := range items {
		result = accumulator(result, item)
	}
	return result
}

func (s *streamImpl[T]) Count() int64 {
	return int64(len(s.execute()))
}

func (s *streamImpl[T]) AnyMatch(predicate Predicate[T]) bool {
	items := s.execute()
	for _, item := range items {
		if predicate(item) {
			return true
		}
	}
	return false
}

func (s *streamImpl[T]) AllMatch(predicate Predicate[T]) bool {
	items := s.execute()
	for _, item := range items {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func (s *streamImpl[T]) NoneMatch(predicate Predicate[T]) bool {
	return !s.AnyMatch(predicate)
}

func (s *streamImpl[T]) FindFirst() Optional[T] {
	items := s.execute()
	if len(items) == 0 {
		return EmptyOptional[T]()
	}
	return OfOptional(items[0])
}

func (s *streamImpl[T]) FindAny() Optional[T] {
	return s.FindFirst()
}

func (s *streamImpl[T]) ToSlice() []T {
	return s.execute()
}

func (s *streamImpl[T]) execute() []T {
	if s.isConsumed {
		panic("stream has already been operated upon or closed")
	}
	s.isConsumed = true

	result := make([]T, len(s.source))
	copy(result, s.source)

	for _, op := range s.operations {
		result = op(result)
	}

	return result
}

func (s *streamImpl[T]) checkNotConsumed() {
	if s.isConsumed {
		panic("stream has already been operated upon or closed")
	}
}
