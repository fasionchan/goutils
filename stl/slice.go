/*
 * Author: fasion
 * Created time: 2022-11-14 11:27:56
 * Last Modified by: fasion
 * Last Modified time: 2023-04-20 15:14:31
 */

package stl

import (
	"sort"

	"golang.org/x/exp/constraints"
)

func AnyMatch[Data any](datas []Data, test func(Data) bool) bool {
	for _, data := range datas {
		if test(data) {
			return true
		}
	}
	return false
}

func AllMatch[Data any](datas []Data, test func(Data) bool) bool {
	for _, data := range datas {
		if !test(data) {
			return false
		}
	}
	return true
}

func Equal[Data comparable](as []Data, bs []Data) bool {
	if len(as) != len(bs) {
		return false
	}

	if len(as) == 0 {
		return true
	}

	for i, a := range as {
		if a != bs[i] {
			return false
		}
	}

	return true
}

func EqualBySort[Data constraints.Ordered](as []Data, bs []Data) bool {
	if len(as) != len(bs) {
		return false
	}

	if len(as) == 0 {
		return true
	}

	less := func(a, b Data) bool {
		return a < b
	}

	bs = Sort(bs, less)
	for i, a := range Sort(as, less) {
		if a != bs[i] {
			return false
		}
	}

	return true
}

func EqualBySet[Data comparable](as []Data, bs []Data) bool {
	if len(as) != len(bs) {
		return false
	}

	if len(as) == 0 {
		return true
	}

	return NewSet(as...).Equal(NewSet(bs...))
}

func Compare[Data constraints.Ordered](as []Data, bs []Data) int {
	bn := len(bs)
	for i, a := range as {
		if i >= bn {
			return 1
		}

		b := bs[i]
		if a > b {
			return 1
		} else if a < b {
			return -1
		}
	}

	an := len(as)
	if an == bn {
		return 0
	} else {
		return -1
	}
}

func Find[Data any](datas []Data, test func(Data) bool) (Data, bool) {
	for _, data := range datas {
		if test(data) {
			return data, true
		}
	}

	var zero Data
	return zero, false
}

func IndexOf[Data comparable](datas []Data, target Data) int {
	for i, data := range datas {
		if data == target {
			return i
		}
	}

	return -1
}

func ForEach[Data any](datas []Data, handler func(data Data)) {
	for _, data := range datas {
		handler(data)
	}
}

func Filter[Data any, Datas ~[]Data](datas Datas, filter func(Data) bool) Datas {
	result := make(Datas, 0, len(datas))
	for _, data := range datas {
		if filter(data) {
			result = append(result, data)
		}
	}
	return result
}

func Purge[Data any, Datas ~[]Data](datas Datas, filter func(Data) bool) Datas {
	result := make(Datas, 0, len(datas))
	for _, data := range datas {
		if !filter(data) {
			result = append(result, data)
		}
	}
	return result
}

func PurgeValue[Data comparable, Datas ~[]Data](datas Datas, value Data) Datas {
	result := make(Datas, 0, len(datas))
	for _, data := range datas {
		if data != value {
			result = append(result, data)
		}
	}
	return result

}

func PurgeZero[Data comparable, Datas ~[]Data](datas Datas) Datas {
	var zero Data
	return PurgeValue(datas, zero)
}

func Map[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(Data) Result) []Result {
	results := make([]Result, 0, len(datas))
	for _, data := range datas {
		results = append(results, mapper(data))
	}
	return results
}

func MapArgs[Data any, Result any](mapper func(Data) Result, args ...Data) []Result {
	return Map(args, mapper)
}

func MapAndConcat[Data any, Datas ~[]Data, Result any, Results ~[]Result](datas Datas, mapper func(Data) Results) Results {
	slices := Map(datas, mapper)
	return ConcatSlices(slices...)
}

func Reduce[Data any, Datas ~[]Data, Result any](datas Datas, reducer func(Data, Result) Result, initial Result) (result Result) {
	result = initial
	for _, data := range datas {
		result = reducer(data, result)
	}
	return
}

func Sort[Data any, Datas ~[]Data](datas Datas, less func(a, b Data) bool) Datas {
	sort.Slice(datas, func(i, j int) bool {
		return less(datas[i], datas[j])
	})
	return datas
}

func SortFast[Data constraints.Ordered, Datas ~[]Data](datas Datas) Datas {
	sort.Slice(datas, func(i, j int) bool {
		return datas[i] < datas[j]
	})
	return datas
}

func Unique[Data comparable, Datas ~[]Data](datas Datas, equal func(Data, Data) bool) Datas {
	result := make(Datas, 0, len(datas))
	var last Data
	for i, data := range datas {
		if i == 0 || !equal(data, last) {
			result = append(result, data)
			last = data
		}
	}
	return result
}

func UniqueFast[Data comparable, Datas ~[]Data](datas Datas) Datas {
	result := make(Datas, 0, len(datas))
	var last Data
	for i, data := range datas {
		if i == 0 || data != last {
			result = append(result, data)
			last = data
		}
	}
	return result
}

func UniqueBySet[Data comparable, Datas ~[]Data](datas Datas) Datas {
	return Datas(NewSet(datas...).Slice())
}

func SortUniqueFast[Data constraints.Ordered, Datas ~[]Data](datas Datas) Datas {
	return UniqueFast(SortFast(datas))
}

func DupSlice[Data any, Slice ~[]Data](slice Slice) Slice {
	return append(make(Slice, 0, len(slice)), slice...)
}

func ConcatSlices[Data any, Slice ~[]Data](slices ...Slice) Slice {
	return ConcatSlicesTo(nil, slices...)
}

func GetSliceElemPointers[Data any, Datas ~[]Data](datas Datas) []*Data {
	return Map(datas, func(data Data) *Data { return &data })
}

func ConcatSlicesTo[Data any, Slice ~[]Data](slice Slice, slices ...Slice) Slice {
	return Reduce(slices, func(current Slice, result Slice) Slice {
		return append(result, current...)
	}, slice)
}

func MappingByKey[Data any, Datas ~[]Data, Key comparable](datas Datas, key func(Data) Key) map[Key]Data {
	m := map[Key]Data{}
	for _, data := range datas {
		m[key(data)] = data
	}
	return m
}

func MappingByKeys[Data any, Datas ~[]Data, Key comparable, Keys ~[]Key](datas Datas, keys func(Data) Keys) map[Key]Data {
	m := map[Key]Data{}
	for _, data := range datas {
		for _, key := range UniqueBySet(keys(data)) {
			m[key] = data
		}
	}
	return m
}

func SliceMappingByKey[Data any, Datas ~[]Data, Key comparable](datas Datas, key func(Data) Key) map[Key]Datas {
	m := map[Key]Datas{}
	for _, data := range datas {
		k := key(data)
		m[k] = append(m[k], data)
	}
	return m
}

func SliceMappingByKeys[Data any, Datas ~[]Data, Key comparable, Keys ~[]Key](datas Datas, keys func(Data) Keys) map[Key]Datas {
	m := map[Key]Datas{}
	for _, data := range datas {
		for _, key := range UniqueBySet(keys(data)) {
			m[key] = append(m[key], data)
		}
	}
	return m
}

func FillSliceToCap[Data any, Datas ~[]Data](datas Datas, g func(int) Data) Datas {
	n := cap(datas)
	for i := len(datas); i < n; i++ {
		datas = append(datas, g(i))
	}
	return datas
}

func ReadChan[Data any](c chan Data, n int) []Data {
	datas := make([]Data, 0, n)
	for i := 0; i < n; i++ {
		datas = append(datas, <-c)
	}
	return datas
}

func ReadChanAll[Data any, Datas ~[]Data](c chan Data) (datas Datas) {
	for data := range c {
		datas = append(datas, data)
	}
	return
}

type Slice[Data any] []Data

func NewSlice[Data any](datas ...Data) Slice[Data] {
	return datas
}

func (slice Slice[Data]) Native() []Data {
	return slice
}

func (slice Slice[Data]) NotNilSlice() Slice[Data] {
	if slice == nil {
		return Slice[Data]{}
	}
	return slice
}

func (slice Slice[Data]) AnyMatch(f func(Data) bool) bool {
	return AnyMatch(slice, f)
}

func (slice Slice[Data]) AllMatch(f func(Data) bool) bool {
	return AllMatch(slice, f)
}

func (slice Slice[Data]) ForEachPro(f func(int, Data, Slice[Data])) {
	ForEachPro(slice, f)
}

func (slice Slice[Data]) Map(f func(Data) Data) Slice[Data] {
	return Map(slice, f)
}

func (slice Slice[Data]) Filter(f func(Data) bool) Slice[Data] {
	return Filter(slice, f)
}

func (slice Slice[Data]) Dup() Slice[Data] {
	return DupSlice(slice)
}

type ComparableSlice[Data comparable] []Data

func NewComparableSlice[Data comparable](datas ...Data) ComparableSlice[Data] {
	return datas
}

func (slice ComparableSlice[Data]) Native() []Data {
	return slice
}

func (slice ComparableSlice[Data]) NotNilSlice() ComparableSlice[Data] {
	if slice == nil {
		return ComparableSlice[Data]{}
	}
	return slice
}

func (slice ComparableSlice[Data]) AnyMatch(f func(Data) bool) bool {
	return AnyMatch(slice, f)
}

func (slice ComparableSlice[Data]) AllMatch(f func(Data) bool) bool {
	return AllMatch(slice, f)
}

func (slice ComparableSlice[Data]) ForEachPro(f func(int, Data, ComparableSlice[Data])) {
	ForEachPro(slice, f)
}

func (slice ComparableSlice[Data]) Map(f func(Data) Data) Slice[Data] {
	return Map(slice, f)
}

func (slice ComparableSlice[Data]) Filter(f func(Data) bool) ComparableSlice[Data] {
	return Filter(slice, f)
}

func (slice ComparableSlice[Data]) Dup() ComparableSlice[Data] {
	return DupSlice(slice)
}
