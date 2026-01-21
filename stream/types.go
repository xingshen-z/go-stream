package stream

type Predicate[T any] func(T) bool

type Function[T, R any] func(T) R

type BiFunction[T, U, R any] func(T, U) R

type Consumer[T any] func(T)

type BiConsumer[T, U any] func(T, U)

type Supplier[T any] func() T

type BinaryOperator[T any] func(T, T) T

type Comparator[T any] func(T, T) int
