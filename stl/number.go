package stl

import (
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

func (numbers Numbers[Data]) positionFor(target Data, cond func (data, target Data) bool) int {
	return sort.Search(numbers.Len(), func(i int) bool { return cond(numbers[i], target) })
}

func (numbers Numbers[Data]) AscPositionFor(target Data) int {
	return numbers.positionFor(target, NotLess[Data])
}

func (numbers Numbers[Data]) DescPositionFor(target Data) int {
	return numbers.positionFor(target, NotGreater[Data])
}

func (numbers Numbers[Data]) AscRatioFor(target Data) float64 {
	return float64(numbers.AscPositionFor(target) + 1) / float64(numbers.Len() + 1)
}

func (numbers Numbers[Data]) DescRatioFor(target Data) float64 {
	return float64(numbers.DescPositionFor(target) + 1) / float64(numbers.Len() + 1)
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