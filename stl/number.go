package stl

import (
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Added = Number

type Addend interface {
	constraints.Integer | constraints.Float
}

type Numbers[Data Number] []Data

func (numbers Numbers[Data]) Append(others ...Data) Numbers[Data] {
	return append(numbers, others...)
}

func (numbers Numbers[Data]) Avg() Data {
	return Avg(numbers)
}

func (numbers Numbers[Data]) Concat(others ...Numbers[Data]) Numbers[Data] {
	return ConcatSlicesTo(numbers, others...)
}

func (numbers Numbers[Data]) Dup() Numbers[Data] {
	return DupSlice(numbers)
}

func (numbers Numbers[Data]) Empty() bool {
	return numbers.Len() == 0
}

func (numbers Numbers[Data]) Len() int {
	return len(numbers)
}

func (numbers Numbers[Data]) Max() Data {
	return Max(numbers, 0)
}

func (numbers Numbers[Data]) Negtives() Numbers[Data] {
	return Filter(numbers, IsNegativeNumber[Data])
}

func (numbers Numbers[Data]) Min() Data {
	return Min(numbers, 0)
}

func (numbers Numbers[Data]) positionFor(target Data, cond func(data, target Data) bool) int {
	return sort.Search(numbers.Len(), func(i int) bool { return cond(numbers[i], target) })
}

func (numbers Numbers[Data]) AscPositionFor(target Data) int {
	return numbers.positionFor(target, NotLess[Data])
}

func (numbers Numbers[Data]) DescPositionFor(target Data) int {
	return numbers.positionFor(target, NotGreater[Data])
}

func (numbers Numbers[Data]) AscRatioFor(target Data) float64 {
	return float64(numbers.AscPositionFor(target)+1) / float64(numbers.Len()+1)
}

func (numbers Numbers[Data]) DescRatioFor(target Data) float64 {
	return float64(numbers.DescPositionFor(target)+1) / float64(numbers.Len()+1)
}

func (numbers Numbers[Data]) AscPercentFor(target Data) float64 {
	return numbers.AscRatioFor(target) * 100
}

func (numbers Numbers[Data]) DescPercentFor(target Data) float64 {
	return numbers.DescRatioFor(target) * 100
}

func (numbers Numbers[Data]) Positives() Numbers[Data] {
	return Filter(numbers, IsPositiveNumber[Data])
}

func (numbers Numbers[Data]) PurgeZero() Numbers[Data] {
	return PurgeZero(numbers)
}

func (numbers Numbers[Data]) Sort() Numbers[Data] {
	return numbers.SortPro(Less[Data])
}

func (numbers Numbers[Data]) SortAsc() Numbers[Data] {
	return numbers.SortPro(Less[Data])
}

func (numbers Numbers[Data]) SortDesc() Numbers[Data] {
	return numbers.SortPro(Greater[Data])
}

func (numbers Numbers[Data]) SortPro(less func(a, b Data) bool) Numbers[Data] {
	return Sort(numbers, Less[Data])
}

func (numbers Numbers[Data]) Sum() Data {
	return Sum(numbers, 0)
}

func Avg[
	Datas ~[]Data,
	Data Number,
](datas Datas) Data {
	if len(datas) == 0 {
		return 0
	}

	return Sum(datas, 0) / Data(len(datas))
}

func Less[Data constraints.Ordered](a, b Data) bool {
	return a < b
}

func Greater[Data constraints.Ordered](a, b Data) bool {
	return a > b
}

func NotLess[Data constraints.Ordered](a, b Data) bool {
	return a >= b
}

func NotGreater[Data constraints.Ordered](a, b Data) bool {
	return a <= b
}

func IsNegativeNumber[Data Number](number Data) bool {
	return number < 0
}

func IsPositiveNumber[Data Number](number Data) bool {
	return number > 0
}

func AlignFloor[Int constraints.Integer](value Int, base Int) Int {
	return value / base * base
}

func AlignCeil[Int constraints.Integer](value Int, base Int) Int {
	return (value + base - 1) / base * base
}

func Max[Datas ~[]Data, Data constraints.Ordered](datas Datas, _default Data) (result Data) {
	if len(datas) == 0 {
		return _default
	}

	result = datas[0]
	for _, data := range datas {
		if data > result {
			result = data
		}
	}

	return
}

func Min[Datas ~[]Data, Data constraints.Ordered](datas Datas, _default Data) (result Data) {
	if len(datas) == 0 {
		return _default
	}

	result = datas[0]
	for _, data := range datas {
		if data < result {
			result = data
		}
	}

	return
}

func Sum[Datas ~[]Data, Data Addend](datas Datas, start Data) Data {
	for _, data := range datas {
		start += data
	}
	return start
}

func CalculatePercentage[R Number, T Number](part, total T, precision int) R {
	if total == 0 {
		return 0
	}
	if precision < 0 {
		precision = 0
	}
	return R(RoundFloat64(float64(part)/float64(total)*100, precision))
}

func RoundFloat32(value float32, precision int) float32 {
	result := RoundFloat64(float64(value), precision)
	return float32(result)
}

func RoundFloat64(value float64, precision int) float64 {
	base := math.Pow10(precision)
	return math.Round(value*base) / base
}

type Comparables[T comparable] Slice[T]

func (comparables Comparables[T]) Slice() Slice[T] {
	return Slice[T](comparables)
}

func (comparables Comparables[T]) Filter(value T) Comparables[T] {
	return FilterValue(comparables, value)
}

func (comparables Comparables[T]) Purge(value T) Comparables[T] {
	return PurgeValue(comparables, value)
}

func (comparables Comparables[T]) PurgeZero() Comparables[T] {
	return PurgeZero(comparables)
}

type Orderables[T constraints.Ordered] Comparables[T]

func (orderables Orderables[T]) Native() []T {
	return []T(orderables)
}

func (orderables Orderables[T]) Slice() Slice[T] {
	return (Slice[T])(orderables)
}

func (orderables Orderables[T]) Comparables() Comparables[T] {
	return (Comparables[T])(orderables)
}

func (orderables Orderables[T]) Compare(other Orderables[T]) int {
	return Compare(orderables, other)
}

func (slice Orderables[T]) Less(other Orderables[T]) bool {
	return slice.Compare(other) < 0
}

func (slice Orderables[T]) Greater(other Orderables[T]) bool {
	return slice.Compare(other) > 0
}

func (slice Orderables[T]) Equal(other Orderables[T]) bool {
	return slice.Compare(other) == 0
}