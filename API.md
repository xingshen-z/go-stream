# Go Stream API 文档

## 目录
- [概述](#概述)
- [版本信息](#版本信息)
- [依赖要求](#依赖要求)
- [核心类型](#核心类型)
- [Stream 接口](#stream-接口)
- [工厂函数](#工厂函数)
- [收集器 (Collectors)](#收集器-collectors)
- [Optional 类型](#optional-类型)
- [使用示例](#使用示例)
- [注意事项](#注意事项)

---

## 概述

`go-stream` 是一个基于 Go 1.18+ 泛型实现的流式处理库，提供了类似 Java Stream API 的功能。该库采用惰性求值模式，支持链式调用，提供了丰富的中间操作和终端操作。

### 主要特性
- **类型安全**: 使用 Go 泛型确保编译时类型检查
- **惰性求值**: 中间操作延迟执行，提高性能
- **链式调用**: 流畅的 API 设计，支持方法链
- **功能完整**: 支持过滤、映射、排序、聚合、收集等操作
- **可扩展**: 基于接口设计，易于扩展

---

## 版本信息

| 项目 | 版本 |
|------|------|
| Go 版本要求 | 1.21+ |
| 模块名称 | github.com/xingshen-z/go-stream |
| 当前版本 | 1.0.0 |

---

## 依赖要求

```go
module github.com/xingshen-z/go-stream

go 1.21
```

### 安装

```bash
go get github.com/xingshen-z/go-stream
```

### 导入

```go
import "github.com/xingshen-z/go-stream/stream"
```

---

## 核心类型

### 函数式接口类型

#### Predicate[T any]
```go
type Predicate[T any] func(T) bool
```
**描述**: 谓词函数，用于测试元素是否满足条件

**参数**:
- `T`: 元素类型

**返回值**:
- `bool`: 如果元素满足条件返回 true，否则返回 false

**示例**:
```go
isEven := func(n int) bool {
    return n%2 == 0
}
```

---

#### Function[T, R any]
```go
type Function[T, R any] func(T) R
```
**描述**: 函数接口，将一个类型转换为另一个类型

**类型参数**:
- `T`: 输入类型
- `R`: 输出类型

**返回值**:
- `R`: 转换后的结果

**示例**:
```go
double := func(n int) int {
    return n * 2
}
```

---

#### BiFunction[T, U, R any]
```go
type BiFunction[T, U, R any] func(T, U) R
```
**描述**: 双参数函数接口

**类型参数**:
- `T`: 第一个参数类型
- `U`: 第二个参数类型
- `R`: 返回值类型

---

#### Consumer[T any]
```go
type Consumer[T any] func(T)
```
**描述**: 消费者函数，接收一个参数但不返回值

**参数**:
- `T`: 元素类型

**示例**:
```go
print := func(n int) {
    fmt.Println(n)
}
```

---

#### BiConsumer[T, U any]
```go
type BiConsumer[T, U any] func(T, U)
```
**描述**: 双参数消费者函数

---

#### Supplier[T any]
```go
type Supplier[T any] func() T
```
**描述**: 供应者函数，不接受参数但返回一个值

**返回值**:
- `T`: 供应的值

**示例**:
```go
random := func() int {
    return rand.Intn(100)
}
```

---

#### BinaryOperator[T any]
```go
type BinaryOperator[T any] func(T, T) T
```
**描述**: 二元操作符，对两个相同类型的值进行操作

**参数**:
- `T`: 操作数类型

**返回值**:
- `T`: 操作结果

**示例**:
```go
sum := func(a, b int) int {
    return a + b
}
```

---

#### Comparator[T any]
```go
type Comparator[T any] func(T, T) int
```
**描述**: 比较器函数，用于比较两个元素

**参数**:
- `T`: 元素类型

**返回值**:
- `int`: 比较结果
  - 负数: 第一个元素小于第二个元素
  - 0: 两个元素相等
  - 正数: 第一个元素大于第二个元素

**示例**:
```go
ascending := func(a, b int) int {
    return a - b
}
```

---

## Stream 接口

### 接口定义

```go
type Stream[T any] interface {
    // 中间操作
    Filter(predicate Predicate[T]) Stream[T]
    Map(mapper Function[T, T]) Stream[T]
    FlatMap(mapper Function[T, Stream[T]]) Stream[T]
    Distinct() Stream[T]
    Sorted(comparator Comparator[T]) Stream[T]
    Limit(maxSize int64) Stream[T]
    Skip(n int64) Stream[T]
    Peek(consumer Consumer[T]) Stream[T]

    // 终端操作
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

---

### 中间操作

#### Filter
```go
Filter(predicate Predicate[T]) Stream[T]
```

**描述**: 过滤流中的元素，只保留满足谓词条件的元素

**参数**:
- `predicate`: 谓词函数，返回 true 的元素会被保留

**返回值**:
- `Stream[T]`: 过滤后的新流

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Filter(func(n int) bool { return n%2 == 0 }).
    ToSlice()
// 结果: [2, 4]
```

**注意事项**:
- 这是一个惰性操作，只有在终端操作时才会执行
- Stream 只能被消费一次，重复使用会 panic

---

#### Map
```go
Map(mapper Function[T, T]) Stream[T]
```

**描述**: 对流中的每个元素应用映射函数，转换元素

**参数**:
- `mapper`: 映射函数，将元素转换为同类型的另一个值

**返回值**:
- `Stream[T]`: 映射后的新流

**示例**:
```go
result := stream.Of(1, 2, 3).
    Map(func(n int) int { return n * 2 }).
    ToSlice()
// 结果: [2, 4, 6]
```

**注意事项**:
- 映射函数必须返回相同类型 T
- 如果需要类型转换，请使用 FlatMap 或先收集到切片再处理

---

#### FlatMap
```go
FlatMap(mapper Function[T, Stream[T]]) Stream[T]
```

**描述**: 将流中的每个元素映射为一个流，然后将所有流连接成一个流

**参数**:
- `mapper`: 映射函数，将元素转换为一个 Stream

**返回值**:
- `Stream[T]`: 扁平化后的新流

**示例**:
```go
result := stream.Of(1, 2, 3).
    FlatMap(func(n int) stream.Stream[int] {
        return stream.Of(n, n*10)
    }).
    ToSlice()
// 结果: [1, 10, 2, 20, 3, 30]
```

---

#### Distinct
```go
Distinct() Stream[T]
```

**描述**: 去除流中的重复元素

**返回值**:
- `Stream[T]`: 去重后的新流

**示例**:
```go
result := stream.Of(1, 2, 2, 3, 3, 3).
    Distinct().
    ToSlice()
// 结果: [1, 2, 3]
```

**注意事项**:
- 使用 map 来跟踪已见过的元素
- 元素类型 T 必须可以用作 map 的键（即可比较）

---

#### Sorted
```go
Sorted(comparator Comparator[T]) Stream[T]
```

**描述**: 根据比较器对流中的元素进行排序

**参数**:
- `comparator`: 比较器函数，定义排序规则

**返回值**:
- `Stream[T]`: 排序后的新流

**示例**:
```go
result := stream.Of(3, 1, 4, 1, 5).
    Sorted(func(a, b int) int { return a - b }).
    ToSlice()
// 结果: [1, 1, 3, 4, 5]
```

**注意事项**:
- 使用冒泡排序算法，对于大数据集可能效率不高
- 排序是稳定的，保持相等元素的相对顺序

---

#### Limit
```go
Limit(maxSize int64) Stream[T]
```

**描述**: 截断流，只保留前 maxSize 个元素

**参数**:
- `maxSize`: 最大元素数量，必须 >= 0

**返回值**:
- `Stream[T]`: 截断后的新流

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Limit(3).
    ToSlice()
// 结果: [1, 2, 3]
```

**注意事项**:
- 如果 maxSize 大于流的大小，返回整个流
- maxSize 为负数时行为未定义

---

#### Skip
```go
Skip(n int64) Stream[T]
```

**描述**: 跳过流的前 n 个元素

**参数**:
- `n`: 要跳过的元素数量，必须 >= 0

**返回值**:
- `Stream[T]`: 跳过后的新流

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Skip(2).
    ToSlice()
// 结果: [3, 4, 5]
```

**注意事项**:
- 如果 n 大于等于流的大小，返回空流
- n 为负数时行为未定义

---

#### Peek
```go
Peek(consumer Consumer[T]) Stream[T]
```

**描述**: 对流中的每个元素执行操作，但不改变流本身，主要用于调试

**参数**:
- `consumer`: 消费者函数，对每个元素执行的操作

**返回值**:
- `Stream[T]`: 原流（未被修改）

**示例**:
```go
result := stream.Of(1, 2, 3).
    Peek(func(n int) { fmt.Println("Processing:", n) }).
    Map(func(n int) int { return n * 2 }).
    ToSlice()
// 输出: Processing: 1, Processing: 2, Processing: 3
// 结果: [2, 4, 6]
```

**注意事项**:
- 主要用于调试和日志记录
- 不应修改元素的状态

---

### 终端操作

#### ForEach
```go
ForEach(consumer Consumer[T])
```

**描述**: 对流中的每个元素执行操作

**参数**:
- `consumer`: 消费者函数

**返回值**:
- 无

**示例**:
```go
stream.Of(1, 2, 3).ForEach(func(n int) {
    fmt.Println(n)
})
// 输出: 1, 2, 3
```

**注意事项**:
- 这是一个终端操作，会触发流的执行
- 执行后流被消费，不能再使用

---

#### Collect
```go
Collect(collector Collector[T, any, any]) any
```

**描述**: 使用收集器将流中的元素收集到某种数据结构中

**参数**:
- `collector`: 收集器，定义如何收集元素

**返回值**:
- `any`: 收集结果，需要类型断言获取具体类型

**示例**:
```go
result := stream.Of(1, 2, 3).
    Collect(stream.ToSlice[int]())
slice, ok := result.([]int)
// slice = [1, 2, 3], ok = true
```

**注意事项**:
- 返回值类型为 any，需要使用类型断言
- 常用收集器见 Collectors 章节

---

#### Reduce
```go
Reduce(identity T, accumulator BinaryOperator[T]) T
```

**描述**: 使用累加器将流中的元素归约为单个值

**参数**:
- `identity`: 初始值
- `accumulator`: 累加器函数，将当前结果和下一个元素合并

**返回值**:
- `T`: 归约结果

**示例**:
```go
sum := stream.Of(1, 2, 3, 4, 5).
    Reduce(0, func(a, b int) int { return a + b })
// sum = 15

product := stream.Of(1, 2, 3, 4).
    Reduce(1, func(a, b int) int { return a * b })
// product = 24
```

**注意事项**:
- 如果流为空，返回 identity
- 累加器函数应该是结合的（满足结合律）

---

#### Count
```go
Count() int64
```

**描述**: 返回流中的元素数量

**返回值**:
- `int64`: 元素数量

**示例**:
```go
count := stream.Of(1, 2, 3, 4, 5).Count()
// count = 5
```

---

#### AnyMatch
```go
AnyMatch(predicate Predicate[T]) bool
```

**描述**: 检查流中是否有任意元素满足谓词条件

**参数**:
- `predicate`: 谓词函数

**返回值**:
- `bool`: 如果有元素满足条件返回 true，否则返回 false

**示例**:
```go
hasEven := stream.Of(1, 3, 5, 7).
    AnyMatch(func(n int) bool { return n%2 == 0 })
// hasEven = false

hasPositive := stream.Of(-1, 0, 1).
    AnyMatch(func(n int) bool { return n > 0 })
// hasPositive = true
```

**注意事项**:
- 短路操作，找到匹配元素后立即返回

---

#### AllMatch
```go
AllMatch(predicate Predicate[T]) bool
```

**描述**: 检查流中是否所有元素都满足谓词条件

**参数**:
- `predicate`: 谓词函数

**返回值**:
- `bool`: 如果所有元素都满足条件返回 true，否则返回 false

**示例**:
```go
allPositive := stream.Of(1, 2, 3, 4, 5).
    AllMatch(func(n int) bool { return n > 0 })
// allPositive = true

allEven := stream.Of(2, 4, 6, 8).
    AllMatch(func(n int) bool { return n%2 == 0 })
// allEven = true
```

**注意事项**:
- 空流返回 true
- 短路操作，找到不匹配元素后立即返回

---

#### NoneMatch
```go
NoneMatch(predicate Predicate[T]) bool
```

**描述**: 检查流中是否没有元素满足谓词条件

**参数**:
- `predicate`: 谓词函数

**返回值**:
- `bool`: 如果没有元素满足条件返回 true，否则返回 false

**示例**:
```go
noneNegative := stream.Of(1, 2, 3, 4, 5).
    NoneMatch(func(n int) bool { return n < 0 })
// noneNegative = true

noneEven := stream.Of(1, 3, 5, 7).
    NoneMatch(func(n int) bool { return n%2 == 0 })
// noneEven = true
```

**注意事项**:
- 等价于 !AnyMatch(predicate)
- 空流返回 true

---

#### FindFirst
```go
FindFirst() Optional[T]
```

**描述**: 返回流中的第一个元素

**返回值**:
- `Optional[T]`: 包含第一个元素的 Optional，如果流为空则返回空 Optional

**示例**:
```go
first := stream.Of(3, 1, 4, 1, 5).FindFirst()
// first = Optional{value: &3}

emptyFirst := stream.Empty[int]().FindFirst()
// emptyFirst = Optional{value: nil}
```

**注意事项**:
- 返回 Optional 类型，安全地处理可能为空的情况
- 短路操作，找到第一个元素后立即返回

---

#### FindAny
```go
FindAny() Optional[T]
```

**描述**: 返回流中的任意一个元素

**返回值**:
- `Optional[T]`: 包含某个元素的 Optional，如果流为空则返回空 Optional

**示例**:
```go
any := stream.Of(1, 2, 3).FindAny()
// any 包含 1, 2, 或 3 中的一个

emptyAny := stream.Empty[int]().FindAny()
// emptyAny = Optional{value: nil}
```

**注意事项**:
- 当前实现等同于 FindFirst()
- 在并行流中可能返回不同的元素

---

#### ToSlice
```go
ToSlice() []T
```

**描述**: 将流中的元素收集到切片中

**返回值**:
- `[]T`: 包含所有元素的切片

**示例**:
```go
slice := stream.Of(1, 2, 3, 4, 5).
    Filter(func(n int) bool { return n%2 == 0 }).
    ToSlice()
// slice = [2, 4]
```

**注意事项**:
- 等价于 Collect(ToSlice[T]())
- 返回的是新切片，不会修改原始数据

---

## 工厂函数

### Of
```go
func Of[T any](items ...T) Stream[T]
```

**描述**: 从可变参数创建流

**参数**:
- `items`: 可变数量的元素

**返回值**:
- `Stream[T]`: 包含指定元素的新流

**示例**:
```go
s := stream.Of(1, 2, 3, 4, 5)
```

---

### OfSlice
```go
func OfSlice[T any](items []T) Stream[T]
```

**描述**: 从切片创建流

**参数**:
- `items`: 元素切片

**返回值**:
- `Stream[T]`: 包含切片元素的新流

**示例**:
```go
nums := []int{1, 2, 3, 4, 5}
s := stream.OfSlice(nums)
```

**注意事项**:
- 会复制切片，不会修改原始数据

---

### Range
```go
func Range(startInclusive, endExclusive int64) Stream[int64]
```

**描述**: 创建一个从 startInclusive（包含）到 endExclusive（不包含）的整数范围流

**参数**:
- `startInclusive`: 起始值（包含）
- `endExclusive`: 结束值（不包含）

**返回值**:
- `Stream[int64]`: 整数范围流

**示例**:
```go
s := stream.Range(1, 5)
// 流包含: 1, 2, 3, 4
```

**注意事项**:
- 如果 endExclusive <= startInclusive，返回空流

---

### RangeClosed
```go
func RangeClosed(startInclusive, endInclusive int64) Stream[int64]
```

**描述**: 创建一个从 startInclusive（包含）到 endInclusive（包含）的整数范围流

**参数**:
- `startInclusive`: 起始值（包含）
- `endInclusive`: 结束值（包含）

**返回值**:
- `Stream[int64]`: 整数范围流

**示例**:
```go
s := stream.RangeClosed(1, 5)
// 流包含: 1, 2, 3, 4, 5
```

**注意事项**:
- 如果 endInclusive < startInclusive，返回空流

---

### Empty
```go
func Empty[T any]() Stream[T]
```

**描述**: 创建一个空流

**返回值**:
- `Stream[T]`: 空流

**示例**:
```go
s := stream.Empty[int]()
```

---

### Generate
```go
func Generate[T any](supplier Supplier[T], count int) Stream[T]
```

**描述**: 使用供应者函数生成指定数量的元素创建流

**参数**:
- `supplier`: 供应者函数，每次调用生成一个元素
- `count`: 要生成的元素数量

**返回值**:
- `Stream[T]`: 包含生成元素的新流

**示例**:
```go
s := stream.Generate(func() int {
    return rand.Intn(100)
}, 10)
```

**注意事项**:
- count 必须为非负数
- 供应者函数会被调用 count 次

---

### Concat
```go
func Concat[T any](streams ...Stream[T]) Stream[T]
```

**描述**: 连接多个流，按顺序合并为一个流

**参数**:
- `streams`: 可变数量的流

**返回值**:
- `Stream[T]`: 合并后的新流

**示例**:
```go
s1 := stream.Of(1, 2, 3)
s2 := stream.Of(4, 5, 6)
s := stream.Concat(s1, s2)
// 流包含: 1, 2, 3, 4, 5, 6
```

**注意事项**:
- 会消费所有输入流
- 保持原始流的顺序

---

## 收集器 (Collectors)

### Collector 接口
```go
type Collector[T, A, R any] interface {
    Collect(items []T) any
}
```

---

### ToSlice
```go
func ToSlice[T any]() Collector[T, any, []T]
```

**描述**: 将元素收集到切片中

**返回值**:
- `Collector[T, any, []T]`: 切片收集器

**示例**:
```go
result := stream.Of(1, 2, 3).
    Collect(stream.ToSlice[int]())
slice := result.([]int)
// slice = [1, 2, 3]
```

---

### ToMap
```go
func ToMap[T any, K comparable, V any](
    keyMapper Function[T, K],
    valueMapper Function[T, V],
) Collector[T, any, map[K]V]
```

**描述**: 将元素收集到 map 中

**类型参数**:
- `K`: 键类型（必须可比较）
- `V`: 值类型

**参数**:
- `keyMapper`: 键映射函数
- `valueMapper`: 值映射函数

**返回值**:
- `Collector[T, any, map[K]V]`: Map 收集器

**示例**:
```go
type Person struct {
    ID   int
    Name string
}

people := []Person{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
}

result := stream.OfSlice(people).
    Collect(stream.ToMap(
        func(p Person) int { return p.ID },
        func(p Person) string { return p.Name },
    ))
m := result.(map[int]string)
// m = {1: "Alice", 2: "Bob"}
```

**注意事项**:
- 如果有重复的键，后面的值会覆盖前面的值

---

### GroupingBy
```go
func GroupingBy[T any, K comparable](
    keyMapper Function[T, K],
) Collector[T, any, map[K][]T]
```

**描述**: 按键对元素进行分组

**类型参数**:
- `K`: 键类型（必须可比较）

**参数**:
- `keyMapper`: 键映射函数

**返回值**:
- `Collector[T, any, map[K][]T]`: 分组收集器

**示例**:
```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
    {Name: "Charlie", Age: 25},
}

result := stream.OfSlice(people).
    Collect(stream.GroupingBy(func(p Person) int { return p.Age }))
groups := result.(map[int][]Person)
// groups = {25: [{Alice, 25}, {Charlie, 25}], 30: [{Bob, 30}]}
```

---

### Counting
```go
func Counting[T any]() Collector[T, any, int64]
```

**描述**: 统计元素数量

**返回值**:
- `Collector[T, any, int64]`: 计数收集器

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Counting[int]())
count := result.(int64)
// count = 5
```

---

### Summing
```go
func Summing[T any](mapper Function[T, int64]) Collector[T, any, int64]
```

**描述**: 计算元素的总和

**参数**:
- `mapper`: 将元素映射为 int64 的函数

**返回值**:
- `Collector[T, any, int64]`: 求和收集器

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Summing(func(n int) int64 { return int64(n) }))
sum := result.(int64)
// sum = 15
```

---

### Averaging
```go
func Averaging[T any](mapper Function[T, int64]) Collector[T, any, float64]
```

**描述**: 计算元素的平均值

**参数**:
- `mapper`: 将元素映射为 int64 的函数

**返回值**:
- `Collector[T, any, float64]`: 平均值收集器

**示例**:
```go
result := stream.Of(1, 2, 3, 4, 5).
    Collect(stream.Averaging(func(n int) int64 { return int64(n) }))
avg := result.(float64)
// avg = 3.0
```

**注意事项**:
- 如果流为空，返回 0.0

---

### Joining
```go
func Joining[T any](delimiter string) Collector[T, any, string]
```

**描述**: 将元素连接成字符串

**参数**:
- `delimiter`: 分隔符

**返回值**:
- `Collector[T, any, string]`: 连接收集器

**示例**:
```go
result := stream.Of("a", "b", "c").
    Collect(stream.Joining[string](","))
str := result.(string)
// str = "a,b,c"
```

---

### JoiningWithMapper
```go
func JoiningWithMapper[T any](
    mapper Function[T, string],
    delimiter string,
) Collector[T, any, string]
```

**描述**: 使用映射函数将元素连接成字符串

**参数**:
- `mapper`: 将元素映射为字符串的函数
- `delimiter`: 分隔符

**返回值**:
- `Collector[T, any, string]`: 连接收集器

**示例**:
```go
result := stream.Of(1, 2, 3).
    Collect(stream.JoiningWithMapper(
        func(n int) string { return strconv.Itoa(n) },
        ",",
    ))
str := result.(string)
// str = "1,2,3"
```

---

### JoiningWithPrefixSuffix
```go
func JoiningWithPrefixSuffix[T any](
    mapper Function[T, string],
    delimiter, prefix, suffix string,
) Collector[T, any, string]
```

**描述**: 使用映射函数将元素连接成字符串，并添加前缀和后缀

**参数**:
- `mapper`: 将元素映射为字符串的函数
- `delimiter`: 分隔符
- `prefix`: 前缀
- `suffix`: 后缀

**返回值**:
- `Collector[T, any, string]`: 连接收集器

**示例**:
```go
result := stream.Of(1, 2, 3).
    Collect(stream.JoiningWithPrefixSuffix(
        func(n int) string { return strconv.Itoa(n) },
        ",",
        "[",
        "]",
    ))
str := result.(string)
// str = "[1,2,3]"
```

---

### ToInt
```go
func ToInt[T any](mapper Function[T, int64]) Function[T, int64]
```

**描述**: 创建一个 int64 映射函数（辅助函数）

**参数**:
- `mapper`: 映射函数

**返回值**:
- `Function[T, int64]`: 映射函数

---

### ToString
```go
func ToString[T any](mapper Function[T, string]) Function[T, string]
```

**描述**: 创建一个 string 映射函数（辅助函数）

**参数**:
- `mapper`: 映射函数

**返回值**:
- `Function[T, string]`: 映射函数

---

### IntToString
```go
func IntToString(value int64) string
```

**描述**: 将 int64 转换为字符串（辅助函数）

**参数**:
- `value`: int64 值

**返回值**:
- `string`: 字符串表示

---

## Optional 类型

### 类型定义
```go
type Optional[T any] struct {
    value *T
}
```

**描述**: Optional 类型用于安全地处理可能为空的值，避免空指针异常

---

### OfOptional
```go
func OfOptional[T any](value T) Optional[T]
```

**描述**: 创建一个包含非空值的 Optional

**参数**:
- `value`: 非空值

**返回值**:
- `Optional[T]`: 包含指定值的 Optional

**示例**:
```go
opt := stream.OfOptional(42)
```

---

### OfOptionalPtr
```go
func OfOptionalPtr[T any](value *T) Optional[T]
```

**描述**: 从指针创建 Optional

**参数**:
- `value`: 指针，可以为 nil

**返回值**:
- `Optional[T]`: Optional 对象

**示例**:
```go
val := 42
opt := stream.OfOptionalPtr(&val)

emptyOpt := stream.OfOptionalPtr[int](nil)
```

---

### EmptyOptional
```go
func EmptyOptional[T any]() Optional[T]
```

**描述**: 创建一个空的 Optional

**返回值**:
- `Optional[T]`: 空 Optional

**示例**:
```go
opt := stream.EmptyOptional[int]()
```

---

### IsPresent
```go
func (o Optional[T]) IsPresent() bool
```

**描述**: 检查 Optional 是否包含值

**返回值**:
- `bool`: 如果包含值返回 true，否则返回 false

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.IsPresent()) // true

empty := stream.EmptyOptional[int]()
fmt.Println(empty.IsPresent()) // false
```

---

### IsEmpty
```go
func (o Optional[T]) IsEmpty() bool
```

**描述**: 检查 Optional 是否为空

**返回值**:
- `bool`: 如果为空返回 true，否则返回 false

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.IsEmpty()) // false

empty := stream.EmptyOptional[int]()
fmt.Println(empty.IsEmpty()) // true
```

---

### Get
```go
func (o Optional[T]) Get() T
```

**描述**: 获取 Optional 中的值，如果为空返回零值

**返回值**:
- `T`: Optional 中的值或零值

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.Get()) // 42

empty := stream.EmptyOptional[int]()
fmt.Println(empty.Get()) // 0
```

**注意事项**:
- 如果 Optional 为空，返回类型的零值
- 建议使用 OrElse 或 OrElsePanic 来处理空值情况

---

### OrElse
```go
func (o Optional[T]) OrElse(other T) T
```

**描述**: 获取 Optional 中的值，如果为空返回默认值

**参数**:
- `other`: 默认值

**返回值**:
- `T`: Optional 中的值或默认值

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.OrElse(0)) // 42

empty := stream.EmptyOptional[int]()
fmt.Println(empty.OrElse(0)) // 0
```

---

### OrElseGet
```go
func (o Optional[T]) OrElseGet(supplier Supplier[T]) T
```

**描述**: 获取 Optional 中的值，如果为空使用供应者函数生成默认值

**参数**:
- `supplier`: 供应者函数

**返回值**:
- `T`: Optional 中的值或供应者生成的值

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.OrElseGet(func() int { return 0 })) // 42

empty := stream.EmptyOptional[int]()
fmt.Println(empty.OrElseGet(func() int { return 0 })) // 0
```

---

### OrElsePanic
```go
func (o Optional[T]) OrElsePanic(msg string) T
```

**描述**: 获取 Optional 中的值，如果为空则 panic

**参数**:
- `msg`: panic 消息

**返回值**:
- `T`: Optional 中的值

**示例**:
```go
opt := stream.OfOptional(42)
fmt.Println(opt.OrElsePanic("Value is empty")) // 42

empty := stream.EmptyOptional[int]()
// fmt.Println(empty.OrElsePanic("Value is empty")) // panic: Value is empty
```

---

### IfPresent
```go
func (o Optional[T]) IfPresent(consumer Consumer[T])
```

**描述**: 如果 Optional 包含值，则执行消费者函数

**参数**:
- `consumer`: 消费者函数

**示例**:
```go
opt := stream.OfOptional(42)
opt.IfPresent(func(n int) {
    fmt.Println("Value:", n)
})
// 输出: Value: 42

empty := stream.EmptyOptional[int]()
empty.IfPresent(func(n int) {
    fmt.Println("Value:", n)
})
// 无输出
```

---

### IfPresentOrElse
```go
func (o Optional[T]) IfPresentOrElse(consumer Consumer[T], action func())
```

**描述**: 如果 Optional 包含值则执行消费者函数，否则执行空操作

**参数**:
- `consumer`: 消费者函数（当有值时执行）
- `action`: 空操作函数（当为空时执行）

**示例**:
```go
opt := stream.OfOptional(42)
opt.IfPresentOrElse(
    func(n int) { fmt.Println("Value:", n) },
    func() { fmt.Println("No value") },
)
// 输出: Value: 42

empty := stream.EmptyOptional[int]()
empty.IfPresentOrElse(
    func(n int) { fmt.Println("Value:", n) },
    func() { fmt.Println("No value") },
)
// 输出: No value
```

---

### Filter
```go
func (o Optional[T]) Filter(predicate Predicate[T]) Optional[T]
```

**描述**: 如果 Optional 包含值且值满足谓词，返回此 Optional，否则返回空 Optional

**参数**:
- `predicate`: 谓词函数

**返回值**:
- `Optional[T]`: 过滤后的 Optional

**示例**:
```go
opt := stream.OfOptional(42)
filtered := opt.Filter(func(n int) bool { return n > 0 })
fmt.Println(filtered.IsPresent()) // true

filtered2 := opt.Filter(func(n int) bool { return n < 0 })
fmt.Println(filtered2.IsPresent()) // false
```

---

## 使用示例

### 示例 1: 基本过滤和映射
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

func main() {
    result := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
        Filter(func(n int) bool { return n%2 == 0 }).
        Map(func(n int) int { return n * n }).
        ToSlice()
    
    fmt.Println(result)
    // 输出: [4 16 36 64 100]
}
```

---

### 示例 2: 排序和限制
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

func main() {
    result := stream.Of(5, 3, 8, 1, 9, 2, 7, 4, 6).
        Sorted(func(a, b int) int { return a - b }).
        Limit(5).
        ToSlice()
    
    fmt.Println(result)
    // 输出: [1 2 3 4 5]
}
```

---

### 示例 3: 分组和统计
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 25},
        {Name: "Bob", Age: 30},
        {Name: "Charlie", Age: 25},
        {Name: "David", Age: 30},
        {Name: "Eve", Age: 35},
    }
    
    result := stream.OfSlice(people).
        Collect(stream.GroupingBy(func(p Person) int { return p.Age }))
    
    groups := result.(map[int][]Person)
    
    for age, persons := range groups {
        fmt.Printf("Age %d: ", age)
        for _, p := range persons {
            fmt.Printf("%s ", p.Name)
        }
        fmt.Println()
    }
    // 输出:
    // Age 25: Alice Charlie
    // Age 30: Bob David
    // Age 35: Eve
}
```

---

### 示例 4: 聚合操作
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 求和
    sum := stream.OfSlice(numbers).
        Reduce(0, func(a, b int) int { return a + b })
    fmt.Println("Sum:", sum)
    
    // 求积
    product := stream.OfSlice(numbers).
        Reduce(1, func(a, b int) int { return a * b })
    fmt.Println("Product:", product)
    
    // 最大值
    max := stream.OfSlice(numbers).
        Reduce(numbers[0], func(a, b int) int {
            if a > b {
                return a
            }
            return b
        })
    fmt.Println("Max:", max)
}
```

---

### 示例 5: 匹配操作
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // 检查是否有偶数
    hasEven := stream.OfSlice(numbers).
        AnyMatch(func(n int) bool { return n%2 == 0 })
    fmt.Println("Has even:", hasEven)
    
    // 检查是否都是正数
    allPositive := stream.OfSlice(numbers).
        AllMatch(func(n int) bool { return n > 0 })
    fmt.Println("All positive:", allPositive)
    
    // 检查是否没有负数
    noneNegative := stream.OfSlice(numbers).
        NoneMatch(func(n int) bool { return n < 0 })
    fmt.Println("None negative:", noneNegative)
}
```

---

### 示例 6: 查找操作
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

func main() {
    numbers := []int{5, 3, 8, 1, 9, 2, 7, 4, 6}
    
    // 查找第一个偶数
    firstEven := stream.OfSlice(numbers).
        Filter(func(n int) bool { return n%2 == 0 }).
        FindFirst()
    
    firstEven.IfPresent(func(n int) {
        fmt.Println("First even:", n)
    })
    
    // 查找大于 10 的数
    greaterThan10 := stream.OfSlice(numbers).
        Filter(func(n int) bool { return n > 10 }).
        FindFirst()
    
    if greaterThan10.IsEmpty() {
        fmt.Println("No number greater than 10 found")
    }
}
```

---

### 示例 7: 范围和生成
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
    "math/rand"
    "time"
)

func main() {
    // 使用 Range
    rangeStream := stream.Range(1, 10).
        Filter(func(n int64) bool { return n%2 == 0 }).
        ToSlice()
    fmt.Println("Even numbers from 1 to 9:", rangeStream)
    
    // 使用 RangeClosed
    closedRange := stream.RangeClosed(1, 5).ToSlice()
    fmt.Println("Numbers from 1 to 5:", closedRange)
    
    // 使用 Generate
    rand.Seed(time.Now().UnixNano())
    randomNumbers := stream.Generate(func() int {
        return rand.Intn(100)
    }, 5).ToSlice()
    fmt.Println("Random numbers:", randomNumbers)
}
```

---

### 示例 8: 复杂链式操作
```go
package main

import (
    "fmt"
    "github.com/xingshen-z/go-stream/stream"
)

type Product struct {
    ID       int
    Name     string
    Category string
    Price    float64
}

func main() {
    products := []Product{
        {ID: 1, Name: "Laptop", Category: "Electronics", Price: 999.99},
        {ID: 2, Name: "Mouse", Category: "Electronics", Price: 29.99},
        {ID: 3, Name: "Keyboard", Category: "Electronics", Price: 79.99},
        {ID: 4, Name: "Monitor", Category: "Electronics", Price: 299.99},
        {ID: 5, Name: "Desk", Category: "Furniture", Price: 199.99},
        {ID: 6, Name: "Chair", Category: "Furniture", Price: 149.99},
    }
    
    // 找出电子产品中价格大于 50 的产品，按价格排序，取前 3 个
    result := stream.OfSlice(products).
        Filter(func(p Product) bool { return p.Category == "Electronics" }).
        Filter(func(p Product) bool { return p.Price > 50 }).
        Sorted(func(a, b Product) int {
            if a.Price < b.Price {
                return -1
            } else if a.Price > b.Price {
                return 1
            }
            return 0
        }).
        Limit(3).
        ToSlice()
    
    fmt.Println("Top 3 expensive electronics:")
    for _, p := range result {
        fmt.Printf("- %s: $%.2f\n", p.Name, p.Price)
    }
}
```

---

## 注意事项

### 1. 流的生命周期
- Stream 只能被消费一次，重复使用会 panic
- 终端操作会触发流的执行并消费流
- 中间操作是惰性的，只有在终端操作时才会执行

```go
s := stream.Of(1, 2, 3)
s.ToSlice()
s.ToSlice() // panic: stream has already been operated upon or closed
```

### 2. 类型安全
- Map 操作只能转换为相同类型
- 如果需要类型转换，使用 FlatMap 或先收集到切片

```go
// 错误：Map 不能改变类型
stream.Of(1, 2, 3).Map(func(n int) string { return strconv.Itoa(n) })

// 正确：使用 FlatMap
stream.Of(1, 2, 3).FlatMap(func(n int) stream.Stream[string] {
    return stream.Of(strconv.Itoa(n))
})

// 或先收集到切片再转换
nums := stream.Of(1, 2, 3).ToSlice()
strs := make([]string, len(nums))
for i, n := range nums {
    strs[i] = strconv.Itoa(n)
}
```

### 3. 收集器类型断言
- Collect 返回 `any` 类型，需要类型断言
- 建议使用类型断言的 ok 模式

```go
result := stream.Of(1, 2, 3).Collect(stream.ToSlice[int]())
slice, ok := result.([]int)
if !ok {
    // 处理错误
}
```

### 4. 性能考虑
- 对于简单操作，原生 Go 代码可能更快
- Stream 的优势在于代码可读性和链式调用
- 大数据集时考虑性能影响

```go
// 原生 Go（更快）
nums := []int{1, 2, 3, 4, 5}
var sum int
for _, n := range nums {
    if n%2 == 0 {
        sum += n
    }
}

// Stream（更易读）
sum := stream.OfSlice(nums).
    Filter(func(n int) bool { return n%2 == 0 }).
    Reduce(0, func(a, b int) int { return a + b })
```

### 5. 空流处理
- 空流的 Count 返回 0
- 空流的 AllMatch 返回 true
- 空流的 NoneMatch 返回 true
- 空流的 FindFirst/FindAny 返回空 Optional

### 6. Distinct 的限制
- 元素类型必须可用作 map 的键（即可比较）
- 对于自定义类型，确保实现了正确的相等比较

### 7. 排序性能
- 当前使用冒泡排序，对于大数据集可能效率不高
- 对于大数据集，建议先收集到切片再使用标准库排序

### 8. 并发安全
- Stream 不是并发安全的
- 不要在多个 goroutine 中同时使用同一个 Stream

---

## 错误处理

### 常见错误

#### 1. 流已被消费
```
panic: stream has already been operated upon or closed
```
**原因**: Stream 被多次消费
**解决**: 每次使用创建新的 Stream

#### 2. 类型断言失败
```
panic: interface conversion: interface {} is []int, not []string
```
**原因**: Collect 返回的类型与预期不符
**解决**: 使用 ok 模式进行类型断言

#### 3. 空指针访问
```
panic: runtime error: invalid memory address or nil pointer dereference
```
**原因**: Optional 为空时调用 Get()
**解决**: 使用 OrElse 或 OrElsePanic

---

## 最佳实践

1. **使用方法链**: 充分利用 Stream 的链式调用特性
2. **惰性求值**: 将过滤等操作放在前面，减少数据处理量
3. **使用 Optional**: 安全地处理可能为空的值
4. **类型断言**: 使用 ok 模式处理类型断言
5. **性能权衡**: 在性能和可读性之间找到平衡
6. **避免重复消费**: 每个 Stream 只使用一次
7. **合理使用收集器**: 选择合适的收集器简化代码

---

## 许可证

本项目采用 MIT 许可证。

---

## 联系方式

如有问题或建议，请通过以下方式联系：
- GitHub: https://github.com/xingshen-z/go-stream
- Issues: https://github.com/xingshen-z/go-stream/issues
