package main

import (
	"fmt"

	"github.com/xingshen-z/go-stream/stream"
)

type Person struct {
	Name string
	Age  int
	City string
}

func main() {
	fmt.Println("=== Go Stream API Examples ===")

	example1_BasicOperations()
	example2_FilterAndMap()
	example3_Collectors()
	example4_ReduceAndAggregate()
	example5_MatchingAndFinding()
	example6_ChainedOperations()
	example7_Optional()
	example8_RangeAndGenerate()
}

func example1_BasicOperations() {
	fmt.Println("Example 1: Basic Operations")
	fmt.Println("-----------------------------")

	numbers := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	filtered := numbers.Filter(func(n int) bool {
		return n%2 == 0
	}).ToSlice()

	fmt.Printf("Even numbers: %v\n", filtered)

	mapped := stream.Of(1, 2, 3, 4, 5).Map(func(n int) int {
		return n * 2
	}).ToSlice()

	fmt.Printf("Doubled numbers: %v\n", mapped)

	fmt.Println()
}

func example2_FilterAndMap() {
	fmt.Println("Example 2: Filter and Map")
	fmt.Println("--------------------------")

	people := []Person{
		{"Alice", 25, "New York"},
		{"Bob", 30, "London"},
		{"Charlie", 35, "Paris"},
		{"David", 28, "New York"},
	}

	result := stream.OfSlice(people).
		Filter(func(p Person) bool {
			return p.Age >= 30
		}).
		Map(func(p Person) Person {
			return Person{
				Name: p.Name,
				Age:  p.Age + 1,
				City: p.City,
			}
		}).
		ToSlice()

	fmt.Printf("People aged 30+ with age incremented: %v\n", result)
	fmt.Println()
}

func example3_Collectors() {
	fmt.Println("Example 3: Collectors")
	fmt.Println("---------------------")

	people := []Person{
		{"Alice", 25, "New York"},
		{"Bob", 30, "London"},
		{"Charlie", 25, "Paris"},
		{"David", 30, "New York"},
	}

	toSlice := stream.OfSlice(people).Collect(stream.ToSlice[Person]()).([]Person)
	fmt.Printf("Collected to slice: %v\n", toSlice)

	toMap := stream.OfSlice(people).Collect(stream.ToMap(
		func(p Person) string { return p.Name },
		func(p Person) int { return p.Age },
	)).(map[string]int)
	fmt.Printf("Collected to map: %v\n", toMap)

	groupingBy := stream.OfSlice(people).Collect(stream.GroupingBy(func(p Person) int {
		return p.Age
	})).(map[int][]Person)
	fmt.Printf("Grouped by age: %v\n", groupingBy)

	counting := stream.OfSlice(people).Collect(stream.Counting[Person]()).(int64)
	fmt.Printf("Count: %d\n", counting)

	summing := stream.OfSlice(people).Collect(stream.Summing(func(p Person) int64 {
		return int64(p.Age)
	})).(int64)
	fmt.Printf("Sum of ages: %d\n", summing)

	averaging := stream.OfSlice(people).Collect(stream.Averaging(func(p Person) int64 {
		return int64(p.Age)
	})).(float64)
	fmt.Printf("Average age: %.2f\n", averaging)

	joining := stream.Of("Alice", "Bob", "Charlie").Collect(stream.Joining[string](", ")).(string)
	fmt.Printf("Joined names: %s\n", joining)

	fmt.Println()
}

func example4_ReduceAndAggregate() {
	fmt.Println("Example 4: Reduce and Aggregate")
	fmt.Println("--------------------------------")

	numbers := stream.Of(1, 2, 3, 4, 5)

	sum := numbers.Reduce(0, func(a, b int) int {
		return a + b
	})
	fmt.Printf("Sum: %d\n", sum)

	product := stream.Of(1, 2, 3, 4, 5).Reduce(1, func(a, b int) int {
		return a * b
	})
	fmt.Printf("Product: %d\n", product)

	count := stream.Of(1, 2, 3, 4, 5).Count()
	fmt.Printf("Count: %d\n", count)

	fmt.Println()
}

func example5_MatchingAndFinding() {
	fmt.Println("Example 5: Matching and Finding")
	fmt.Println("--------------------------------")

	numbers := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	anyMatch := numbers.AnyMatch(func(n int) bool {
		return n > 5
	})
	fmt.Printf("Any number > 5: %v\n", anyMatch)

	allMatch := stream.Of(2, 4, 6, 8, 10).AllMatch(func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("All numbers are even: %v\n", allMatch)

	noneMatch := stream.Of(1, 3, 5, 7, 9).NoneMatch(func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("None of the numbers are even: %v\n", noneMatch)

	findFirst := stream.Of(1, 2, 3, 4, 5).FindFirst()
	fmt.Printf("Find first: ")
	findFirst.IfPresent(func(n int) {
		fmt.Printf("%d\n", n)
	})

	findFirstEmpty := stream.Empty[int]().FindFirst()
	fmt.Printf("Find first from empty: present=%v\n", findFirstEmpty.IsPresent())

	fmt.Println()
}

func example6_ChainedOperations() {
	fmt.Println("Example 6: Chained Operations")
	fmt.Println("------------------------------")

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

	fmt.Printf("Top 5 squared even numbers from 1-100: %v\n", result)

	people := []Person{
		{"Alice", 25, "New York"},
		{"Bob", 30, "London"},
		{"Charlie", 35, "Paris"},
		{"David", 28, "New York"},
		{"Eve", 32, "London"},
	}

	filteredPeople := stream.OfSlice(people).
		Filter(func(p Person) bool {
			return p.City == "New York" || p.City == "London"
		}).
		Sorted(func(a, b Person) int {
			if a.Age < b.Age {
				return -1
			} else if a.Age > b.Age {
				return 1
			}
			return 0
		}).
		ToSlice()

	formattedPeople := make([]string, len(filteredPeople))
	for i, p := range filteredPeople {
		formattedPeople[i] = fmt.Sprintf("%s (%d)", p.Name, p.Age)
	}
	fmt.Printf("People from NY/London sorted by age: %v\n", formattedPeople)

	fmt.Println()
}

func example7_Optional() {
	fmt.Println("Example 7: Optional")
	fmt.Println("-------------------")

	value := stream.OfOptional(42)
	fmt.Printf("Optional value present: %v\n", value.IsPresent())
	fmt.Printf("Optional value: %d\n", value.Get())

	empty := stream.EmptyOptional[int]()
	fmt.Printf("Empty optional present: %v\n", empty.IsPresent())
	fmt.Printf("Empty optional or else: %d\n", empty.OrElse(0))

	value2 := stream.OfOptional(100)
	value2.IfPresent(func(n int) {
		fmt.Printf("Value exists: %d\n", n)
	})

	value3 := stream.OfOptional(50)
	filtered := value3.Filter(func(n int) bool {
		return n > 40
	})
	fmt.Printf("Filtered optional present: %v\n", filtered.IsPresent())

	value4 := stream.OfOptional(30)
	filtered2 := value4.Filter(func(n int) bool {
		return n > 40
	})
	fmt.Printf("Filtered optional (no match) present: %v\n", filtered2.IsPresent())

	fmt.Println()
}

func example8_RangeAndGenerate() {
	fmt.Println("Example 8: Range and Generate")
	fmt.Println("-----------------------------")

	rangeStream := stream.Range(1, 11)
	fmt.Printf("Range 1-10: %v\n", rangeStream.ToSlice())

	rangeClosed := stream.RangeClosed(1, 10)
	fmt.Printf("RangeClosed 1-10: %v\n", rangeClosed.ToSlice())

	generated := stream.Generate(func() int {
		return 42
	}, 5)
	fmt.Printf("Generated 5 times: %v\n", generated.ToSlice())

	concatenated := stream.Concat(
		stream.Of(1, 2, 3),
		stream.Of(4, 5, 6),
		stream.Of(7, 8, 9),
	)
	fmt.Printf("Concatenated streams: %v\n", concatenated.ToSlice())

	fmt.Println()
}
