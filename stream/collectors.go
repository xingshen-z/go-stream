package stream

import (
	"fmt"
	"strconv"
)

type Collector[T, A, R any] interface {
	Collect(items []T) any
}

type sliceCollector[T any] struct{}

func (c sliceCollector[T]) Collect(items []T) any {
	result := make([]T, len(items))
	copy(result, items)
	return result
}

func ToSlice[T any]() Collector[T, any, []T] {
	return sliceCollector[T]{}
}

type mapCollector[T any, K comparable, V any] struct {
	keyMapper   Function[T, K]
	valueMapper Function[T, V]
}

func (c mapCollector[T, K, V]) Collect(items []T) any {
	result := make(map[K] V)
	for _, item := range items {
		key := c.keyMapper(item)
		value := c.valueMapper(item)
		result[key] = value
	}
	return result
}

func ToMap[T any, K comparable, V any](keyMapper Function[T, K], valueMapper Function[T, V]) Collector[T, any, map[K]V] {
	return mapCollector[T, K, V]{
		keyMapper:   keyMapper,
		valueMapper: valueMapper,
	}
}

type groupingCollector[T any, K comparable] struct {
	keyMapper Function[T, K]
}

func (c groupingCollector[T, K]) Collect(items []T) any {
	result := make(map[K][]T)
	for _, item := range items {
		key := c.keyMapper(item)
		result[key] = append(result[key], item)
	}
	return result
}

func GroupingBy[T any, K comparable](keyMapper Function[T, K]) Collector[T, any, map[K][]T] {
	return groupingCollector[T, K]{
		keyMapper: keyMapper,
	}
}

type countingCollector[T any] struct{}

func (c countingCollector[T]) Collect(items []T) any {
	return int64(len(items))
}

func Counting[T any]() Collector[T, any, int64] {
	return countingCollector[T]{}
}

type summingCollector[T any] struct {
	mapper Function[T, int64]
}

func (c summingCollector[T]) Collect(items []T) any {
	var sum int64
	for _, item := range items {
		sum += c.mapper(item)
	}
	return sum
}

func Summing[T any](mapper Function[T, int64]) Collector[T, any, int64] {
	return summingCollector[T]{
		mapper: mapper,
	}
}

type averagingCollector[T any] struct {
	mapper Function[T, int64]
}

func (c averagingCollector[T]) Collect(items []T) any {
	if len(items) == 0 {
		return 0.0
	}
	var sum int64
	for _, item := range items {
		sum += c.mapper(item)
	}
	return float64(sum) / float64(len(items))
}

func Averaging[T any](mapper Function[T, int64]) Collector[T, any, float64] {
	return averagingCollector[T]{
		mapper: mapper,
	}
}

type joiningCollector[T any] struct {
	mapper    Function[T, string]
	delimiter string
	prefix    string
	suffix    string
}

func (c joiningCollector[T]) Collect(items []T) any {
	var builder string
	builder += c.prefix
	for i, item := range items {
		if i > 0 {
			builder += c.delimiter
		}
		builder += c.mapper(item)
	}
	builder += c.suffix
	return builder
}

func Joining[T any](delimiter string) Collector[T, any, string] {
	return JoiningWithMapper[T](func(t T) string {
		return fmt.Sprintf("%v", t)
	}, delimiter)
}

func JoiningWithMapper[T any](mapper Function[T, string], delimiter string) Collector[T, any, string] {
	return joiningCollector[T]{
		mapper:    mapper,
		delimiter: delimiter,
	}
}

func JoiningWithPrefixSuffix[T any](mapper Function[T, string], delimiter, prefix, suffix string) Collector[T, any, string] {
	return joiningCollector[T]{
		mapper:    mapper,
		delimiter: delimiter,
		prefix:    prefix,
		suffix:    suffix,
	}
}

func ToInt[T any](mapper Function[T, int64]) Function[T, int64] {
	return mapper
}

func ToString[T any](mapper Function[T, string]) Function[T, string] {
	return mapper
}

func IntToString(value int64) string {
	return strconv.FormatInt(value, 10)
}
