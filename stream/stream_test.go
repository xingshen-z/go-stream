package stream

import (
	"testing"
)

func TestOf(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.ToSlice()

	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}

	expected := []int{1, 2, 3, 4, 5}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestOfSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	s := OfSlice(input)
	result := s.ToSlice()

	if len(result) != 5 {
		t.Errorf("Expected length 5, got %d", len(result))
	}

	for i, v := range result {
		if v != input[i] {
			t.Errorf("Expected %d at index %d, got %d", input[i], i, v)
		}
	}
}

func TestEmpty(t *testing.T) {
	s := Empty[int]()
	result := s.ToSlice()

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(result))
	}
}

func TestRange(t *testing.T) {
	s := Range(0, 5)
	result := s.ToSlice()

	expected := []int64{0, 1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestRangeClosed(t *testing.T) {
	s := RangeClosed(1, 5)
	result := s.ToSlice()

	expected := []int64{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFilter(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6)
	result := s.Filter(func(n int) bool {
		return n%2 == 0
	}).ToSlice()

	expected := []int{2, 4, 6}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestMap(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.Map(func(n int) int {
		return n * 2
	}).ToSlice()

	expected := []int{2, 4, 6, 8, 10}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestFlatMap(t *testing.T) {
	s := Of(1, 2, 3)
	result := s.FlatMap(func(n int) Stream[int] {
		return Of(n, n*10)
	}).ToSlice()

	expected := []int{1, 10, 2, 20, 3, 30}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestDistinct(t *testing.T) {
	s := Of(1, 2, 2, 3, 3, 3, 4, 5)
	result := s.Distinct().ToSlice()

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestSorted(t *testing.T) {
	s := Of(5, 3, 1, 4, 2)
	result := s.Sorted(func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}).ToSlice()

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestLimit(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	result := s.Limit(5).ToSlice()

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestSkip(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	result := s.Skip(5).ToSlice()

	expected := []int{6, 7, 8, 9, 10}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestPeek(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	peeked := []int{}
	result := s.Peek(func(n int) {
		peeked = append(peeked, n)
	}).ToSlice()

	if len(peeked) != 5 {
		t.Errorf("Expected 5 peeked elements, got %d", len(peeked))
	}

	for i, v := range peeked {
		if v != result[i] {
			t.Errorf("Peeked value %d at index %d doesn't match result %d", v, i, result[i])
		}
	}
}

func TestForEach(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := []int{}
	s.ForEach(func(n int) {
		result = append(result, n)
	})

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestReduce(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.Reduce(0, func(a, b int) int {
		return a + b
	})

	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
}

func TestCount(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	count := s.Count()

	if count != 5 {
		t.Errorf("Expected 5, got %d", count)
	}
}

func TestAnyMatch(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.AnyMatch(func(n int) bool {
		return n > 3
	})

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestAllMatch(t *testing.T) {
	s := Of(2, 4, 6, 8, 10)
	result := s.AllMatch(func(n int) bool {
		return n%2 == 0
	})

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestNoneMatch(t *testing.T) {
	s := Of(1, 3, 5, 7, 9)
	result := s.NoneMatch(func(n int) bool {
		return n%2 == 0
	})

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestFindFirst(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result := s.FindFirst()

	if !result.IsPresent() {
		t.Errorf("Expected present value, got empty")
	}

	if result.Get() != 1 {
		t.Errorf("Expected 1, got %d", result.Get())
	}
}

func TestFindFirstEmpty(t *testing.T) {
	s := Empty[int]()
	result := s.FindFirst()

	if result.IsPresent() {
		t.Errorf("Expected empty value, got present")
	}
}

func TestCollectToSlice(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result, ok := s.Collect(ToSlice[int]()).([]int)

	if !ok {
		t.Errorf("Expected []int type")
	}

	expected := []int{1, 2, 3, 4, 5}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestCollectToMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	s := Of(
		Person{"Alice", 25},
		Person{"Bob", 30},
		Person{"Charlie", 35},
	)

	result, ok := s.Collect(ToMap(func(p Person) string {
		return p.Name
	}, func(p Person) int {
		return p.Age
	})).(map[string]int)

	if !ok {
		t.Errorf("Expected map[string]int type")
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(result))
	}

	if result["Alice"] != 25 {
		t.Errorf("Expected Alice's age 25, got %d", result["Alice"])
	}
}

func TestCollectGroupingBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	s := Of(
		Person{"Alice", 25},
		Person{"Bob", 30},
		Person{"Charlie", 25},
		Person{"David", 30},
	)

	result, ok := s.Collect(GroupingBy(func(p Person) int {
		return p.Age
	})).(map[int][]Person)

	if !ok {
		t.Errorf("Expected map[int][]Person type")
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}

	if len(result[25]) != 2 {
		t.Errorf("Expected 2 people with age 25, got %d", len(result[25]))
	}
}

func TestCollectCounting(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result, ok := s.Collect(Counting[int]()).(int64)

	if !ok {
		t.Errorf("Expected int64 type")
	}

	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestCollectSumming(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result, ok := s.Collect(Summing(func(n int) int64 {
		return int64(n)
	})).(int64)

	if !ok {
		t.Errorf("Expected int64 type")
	}

	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
}

func TestCollectAveraging(t *testing.T) {
	s := Of(1, 2, 3, 4, 5)
	result, ok := s.Collect(Averaging(func(n int) int64 {
		return int64(n)
	})).(float64)

	if !ok {
		t.Errorf("Expected float64 type")
	}

	expected := 3.0
	if result != expected {
		t.Errorf("Expected %f, got %f", expected, result)
	}
}

func TestCollectJoining(t *testing.T) {
	s := Of("a", "b", "c", "d", "e")
	result, ok := s.Collect(Joining[string](",")).(string)

	if !ok {
		t.Errorf("Expected string type")
	}

	expected := "a,b,c,d,e"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestChainedOperations(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	result := s.
		Filter(func(n int) bool {
			return n%2 == 0
		}).
		Map(func(n int) int {
			return n * 2
		}).
		Sorted(func(a, b int) int {
			return b - a
		}).
		Limit(3).
		ToSlice()

	expected := []int{20, 16, 12}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestOptionalIsPresent(t *testing.T) {
	opt := OfOptional(42)
	if !opt.IsPresent() {
		t.Errorf("Expected present, got empty")
	}

	emptyOpt := EmptyOptional[int]()
	if emptyOpt.IsPresent() {
		t.Errorf("Expected empty, got present")
	}
}

func TestOptionalGet(t *testing.T) {
	opt := OfOptional(42)
	if opt.Get() != 42 {
		t.Errorf("Expected 42, got %d", opt.Get())
	}
}

func TestOptionalOrElse(t *testing.T) {
	opt := OfOptional(42)
	if opt.OrElse(0) != 42 {
		t.Errorf("Expected 42, got %d", opt.OrElse(0))
	}

	emptyOpt := EmptyOptional[int]()
	if emptyOpt.OrElse(0) != 0 {
		t.Errorf("Expected 0, got %d", emptyOpt.OrElse(0))
	}
}

func TestOptionalOrElseGet(t *testing.T) {
	opt := OfOptional(42)
	if opt.OrElseGet(func() int { return 0 }) != 42 {
		t.Errorf("Expected 42, got %d", opt.OrElseGet(func() int { return 0 }))
	}

	emptyOpt := EmptyOptional[int]()
	if emptyOpt.OrElseGet(func() int { return 100 }) != 100 {
		t.Errorf("Expected 100, got %d", emptyOpt.OrElseGet(func() int { return 100 }))
	}
}

func TestOptionalIfPresent(t *testing.T) {
	opt := OfOptional(42)
	called := false
	opt.IfPresent(func(n int) {
		called = true
		if n != 42 {
			t.Errorf("Expected 42, got %d", n)
		}
	})

	if !called {
		t.Errorf("Expected consumer to be called")
	}

	emptyOpt := EmptyOptional[int]()
	called = false
	emptyOpt.IfPresent(func(n int) {
		called = true
	})

	if called {
		t.Errorf("Expected consumer not to be called")
	}
}

func TestOptionalFilter(t *testing.T) {
	opt := OfOptional(42)
	result := opt.Filter(func(n int) bool {
		return n > 40
	})

	if !result.IsPresent() {
		t.Errorf("Expected present, got empty")
	}

	result = opt.Filter(func(n int) bool {
		return n < 40
	})

	if result.IsPresent() {
		t.Errorf("Expected empty, got present")
	}
}
