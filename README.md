# Go Stream API

A powerful and expressive stream processing library for Go, inspired by Java Stream API. This library provides a fluent, chainable API for processing collections of data with operations like filter, map, reduce, collect, and more.

## Features

- **Type-safe**: Built with Go 1.18+ generics for compile-time type safety
- **Fluent API**: Chain multiple operations together for readable code
- **Comprehensive operations**: Filter, Map, FlatMap, Distinct, Sorted, Limit, Skip, Peek, Reduce, Collect, and more
- **Rich collectors**: ToSlice, ToMap, GroupingBy, Counting, Summing, Averaging, Joining
- **Optional type**: Safe handling of potentially null values
- **Well-tested**: Comprehensive unit tests and benchmarks

## Installation

```bash
go get github.com/xingshen/go-stream
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/xingshen/go-stream/stream"
)

func main() {
    // Filter even numbers
    result := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
        Filter(func(n int) bool {
            return n%2 == 0
        }).
        ToSlice()
    
    fmt.Println(result) // [2 4 6 8 10]
}
```

## Core Concepts

### Creating Streams

```go
// From variadic arguments
s := stream.Of(1, 2, 3, 4, 5)

// From a slice
data := []int{1, 2, 3, 4, 5}
s := stream.OfSlice(data)

// Range of numbers
s := stream.Range(0, 10) // 0-9
s := stream.RangeClosed(1, 10) // 1-10

// Empty stream
s := stream.Empty[int]()

// Generate values
s := stream.Generate(func() int { return 42 }, 5)

// Concatenate streams
s := stream.Concat(stream.Of(1, 2), stream.Of(3, 4))
```

### Intermediate Operations

Intermediate operations return a new Stream and are lazy (they don't execute until a terminal operation is called).

#### Filter

```go
stream.Of(1, 2, 3, 4, 5, 6).
    Filter(func(n int) bool {
        return n%2 == 0
    }).
    ToSlice() // [2 4 6]
```

#### Map

```go
stream.Of(1, 2, 3, 4, 5).
    Map(func(n int) int {
        return n * 2
    }).
    ToSlice() // [2 4 6 8 10]
```

#### FlatMap

```go
stream.Of(1, 2, 3).
    FlatMap(func(n int) stream.Stream[int] {
        return stream.Of(n, n*10)
    }).
    ToSlice() // [1 10 2 20 3 30]
```

#### Distinct

```go
stream.Of(1, 2, 2, 3, 3, 3, 4).
    Distinct().
    ToSlice() // [1 2 3 4]
```

#### Sorted

```go
stream.Of(5, 3, 1, 4, 2).
    Sorted(func(a, b int) int {
        if a < b {
            return -1
        } else if a > b {
            return 1
        }
        return 0
    }).
    ToSlice() // [1 2 3 4 5]
```

#### Limit and Skip

```go
stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
    Skip(5).
    Limit(3).
    ToSlice() // [6 7 8]
```

#### Peek

```go
stream.Of(1, 2, 3).
    Peek(func(n int) {
        fmt.Println("Processing:", n)
    }).
    ToSlice()
```

### Terminal Operations

Terminal operations consume the stream and produce a result.

#### ForEach

```go
stream.Of(1, 2, 3, 4, 5).
    ForEach(func(n int) {
        fmt.Println(n)
    })
```

#### Collect

```go
// Collect to slice
result := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.ToSlice[int]()).([]int)

// Collect to map
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 25},
    {"Bob", 30},
}

nameToAge := stream.OfSlice(people).
    Collect(stream.ToMap(
        func(p Person) string { return p.Name },
        func(p Person) int { return p.Age },
    )).(map[string]int)

// Grouping
groupedByAge := stream.OfSlice(people).
    Collect(stream.GroupingBy(func(p Person) int {
        return p.Age
    })).(map[int][]Person)

// Counting
count := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Counting[int]()).(int64)

// Summing
sum := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Summing(func(n int) int64 {
        return int64(n)
    })).(int64)

// Averaging
avg := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Averaging(func(n int) int64 {
        return int64(n)
    })).(float64)

// Joining
joined := stream.Of("a", "b", "c").
    Collect(stream.Joining[string](", ")).(string)
```

#### Reduce

```go
// Sum
sum := stream.Of(1, 2, 3, 4, 5).
    Reduce(0, func(a, b int) int {
        return a + b
    }) // 15

// Product
product := stream.Of(1, 2, 3, 4, 5).
    Reduce(1, func(a, b int) int {
        return a * b
    }) // 120
```

#### Count

```go
count := stream.Of(1, 2, 3, 4, 5).Count() // 5
```

#### Matching Operations

```go
// AnyMatch
anyMatch := stream.Of(1, 2, 3, 4, 5).
    AnyMatch(func(n int) bool {
        return n > 3
    }) // true

// AllMatch
allMatch := stream.Of(2, 4, 6, 8, 10).
    AllMatch(func(n int) bool {
        return n%2 == 0
    }) // true

// NoneMatch
noneMatch := stream.Of(1, 3, 5, 7, 9).
    NoneMatch(func(n int) bool {
        return n%2 == 0
    }) // true
```

#### Finding Operations

```go
// FindFirst
first := stream.Of(1, 2, 3, 4, 5).
    FindFirst()
first.IfPresent(func(n int) {
    fmt.Println("First element:", n)
})

// FindAny (same as FindFirst for sequential streams)
any := stream.Of(1, 2, 3, 4, 5).
    FindAny()
```

#### ToSlice

```go
result := stream.Of(1, 2, 3, 4, 5).ToSlice() // [1 2 3 4 5]
```

### Optional Type

The Optional type provides a way to handle potentially null values safely.

```go
// Create Optional
opt := stream.OfOptional(42)
emptyOpt := stream.EmptyOptional[int]()

// Check if present
if opt.IsPresent() {
    fmt.Println("Value exists:", opt.Get())
}

// Get with default
value := emptyOpt.OrElse(0)

// Get with supplier
value := emptyOpt.OrElseGet(func() int {
    return 100
})

// Conditional action
opt.IfPresent(func(n int) {
    fmt.Println("Processing:", n)
})

// Filter Optional
filtered := opt.Filter(func(n int) bool {
    return n > 40
})
```

## Chained Operations

You can chain multiple operations together:

```go
result := stream.Range(1, 101).
    Filter(func(n int64) bool {
        return n%2 == 0
    }).
    Map(func(n int64) int64 {
        return n * n
    }).
    Sorted(func(a, b int64) int {
        if a < b {
            return 1
        } else if a > b {
            return -1
        }
        return 0
    }).
    Limit(5).
    ToSlice()
```

## Performance

The library includes comprehensive benchmarks. Here are some sample results on Apple M1:

```
BenchmarkFilter-8                 208956              4870 ns/op
BenchmarkMap-8                    284678              3564 ns/op
BenchmarkReduce-8                 427936              2505 ns/op
BenchmarkCount-8                  698742              1668 ns/op
BenchmarkChainedOperations-8        5965            201650 ns/op
```

For comparison with native Go operations:

```
BenchmarkNativeFilter-8            637639              1892 ns/op
BenchmarkNativeMap-8              1000000              1078 ns/op
BenchmarkNativeReduce-8           3320910               352.9 ns/op
```

While the stream API has some overhead compared to native operations, it provides significant benefits in code readability and maintainability for complex data processing tasks.

## API Reference

### Stream Interface

```go
type Stream[T any] interface {
    // Intermediate operations
    Filter(predicate Predicate[T]) Stream[T]
    Map(mapper Function[T, T]) Stream[T]
    FlatMap(mapper Function[T, Stream[T]]) Stream[T]
    Distinct() Stream[T]
    Sorted(comparator Comparator[T]) Stream[T]
    Limit(maxSize int64) Stream[T]
    Skip(n int64) Stream[T]
    Peek(consumer Consumer[T]) Stream[T]

    // Terminal operations
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
}
```

### Type Definitions

```go
type Predicate[T any] func(T) bool
type Function[T, R any] func(T) R
type Consumer[T any] func(T)
type BinaryOperator[T any] func(T, T) T
type Comparator[T any] func(T, T) int
```

### Collectors

- `ToSlice[T any]()`: Collect to slice
- `ToMap[T any, K comparable, V any](keyMapper, valueMapper)`: Collect to map
- `GroupingBy[T any, K comparable](keyMapper)`: Group by key
- `Counting[T any]()`: Count elements
- `Summing[T any](mapper)`: Sum mapped values
- `Averaging[T any](mapper)`: Average mapped values
- `Joining[T any](delimiter)`: Join as string

## Examples

See the [examples](examples/examples.go) directory for comprehensive usage examples covering all major features.

## Testing

Run the test suite:

```bash
go test ./stream/... -v
```

Run benchmarks:

```bash
go test ./stream/... -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Acknowledgments

Inspired by Java Stream API and designed to provide similar functionality in Go with idiomatic Go patterns.
