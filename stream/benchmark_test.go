package stream

import (
	"fmt"
	"testing"
)

func BenchmarkFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Filter(func(n int) bool {
			return n%2 == 0
		}).ToSlice()
	}
}

func BenchmarkMap(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Map(func(n int) int {
			return n * 2
		}).ToSlice()
	}
}

func BenchmarkFlatMap(b *testing.B) {
	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).FlatMap(func(n int) Stream[int] {
			return Of(n, n*10)
		}).ToSlice()
	}
}

func BenchmarkDistinct(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i % 100
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Distinct().ToSlice()
	}
}

func BenchmarkSorted(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = 1000 - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Sorted(func(a, b int) int {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		}).ToSlice()
	}
}

func BenchmarkLimit(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Limit(100).ToSlice()
	}
}

func BenchmarkSkip(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Skip(900).ToSlice()
	}
}

func BenchmarkChainedOperations(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).
			Filter(func(n int) bool {
				return n%2 == 0
			}).
			Map(func(n int) int {
				return n * 2
			}).
			Sorted(func(a, b int) int {
				return b - a
			}).
			Limit(100).
			ToSlice()
	}
}

func BenchmarkReduce(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Reduce(0, func(a, b int) int {
			return a + b
		})
	}
}

func BenchmarkCount(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Count()
	}
}

func BenchmarkCollectToSlice(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(ToSlice[int]())
	}
}

func BenchmarkCollectToMap(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	data := make([]Person, 100)
	for i := 0; i < 100; i++ {
		data[i] = Person{
			Name: fmt.Sprintf("Person%d", i),
			Age:  i,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(ToMap(func(p Person) string {
			return p.Name
		}, func(p Person) int {
			return p.Age
		}))
	}
}

func BenchmarkCollectGroupingBy(b *testing.B) {
	type Person struct {
		Name string
		Age  int
	}

	data := make([]Person, 100)
	for i := 0; i < 100; i++ {
		data[i] = Person{
			Name: fmt.Sprintf("Person%d", i),
			Age:  i % 10,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(GroupingBy(func(p Person) int {
			return p.Age
		}))
	}
}

func BenchmarkCollectCounting(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(Counting[int]())
	}
}

func BenchmarkCollectSumming(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(Summing(func(n int) int64 {
			return int64(n)
		}))
	}
}

func BenchmarkCollectAveraging(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(Averaging(func(n int) int64 {
			return int64(n)
		}))
	}
}

func BenchmarkCollectJoining(b *testing.B) {
	data := make([]string, 100)
	for i := 0; i < 100; i++ {
		data[i] = fmt.Sprintf("Item%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		OfSlice(data).Collect(Joining[string](","))
	}
}

func BenchmarkNativeFilter(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0)
		for _, n := range data {
			if n%2 == 0 {
				result = append(result, n)
			}
		}
		_ = result
	}
}

func BenchmarkNativeMap(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, len(data))
		for i, n := range data {
			result[i] = n * 2
		}
		_ = result
	}
}

func BenchmarkNativeReduce(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := 0
		for _, n := range data {
			result += n
		}
		_ = result
	}
}
