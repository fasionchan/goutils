/*
 * Author: fasion
 * Created time: 2022-11-14 11:27:56
 * Last Modified by: fasion
 * Last Modified time: 2023-12-15 18:17:54
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

func Backmost[Datas ~[]Data, Data any](datas Datas, before func(a, b Data) bool) (result Data, index int) {
	result, _ = BackmostPro(datas, before)
	return
}

func BackmostPro[Datas ~[]Data, Data any](datas Datas, before func(a, b Data) bool) (result Data, index int) {
	index = -1
	for i, data := range datas {
		if index == 0 || !before(data, result) {
			result = data
			index = i
		}
	}
	return
}

func SliceEqual[Datas ~[]Data, Data comparable](as Datas, bs Datas) bool {
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

func FirstOneOrZero[Datas ~[]Data, Data any](datas Datas) (data Data) {
	if len(datas) > 0 {
		data = datas[0]
	}
	return
}

func LastOneOrZero[Datas ~[]Data, Data any](datas Datas) (data Data) {
	if len(datas) > 0 {
		data = datas[len(datas)-1]
	}
	return
}

func FindFirst[Datas ~[]Data, Data any](datas Datas, test func(Data) bool) (Data, bool) {
	for _, data := range datas {
		if test(data) {
			return data, true
		}
	}

	var zero Data
	return zero, false
}

func FindFirstOrDefault[Data any](datas []Data, test func(Data) bool, defaultData Data) Data {
	if data, ok := FindFirst(datas, test); ok {
		return data
	} else {
		return defaultData
	}
}

func FindFirstOrZero[Data any](datas []Data, test func(Data) bool) Data {
	data, _ := FindFirst(datas, test)
	return data
}

func FindFirstNotZero[Data comparable](datas []Data) Data {
	var zero Data
	return FindFirstOrZero(datas, func(data Data) bool {
		return data != zero
	})
}

func FindLast[Data any](datas []Data, test func(Data) bool) (Data, bool) {
	for i := len(datas) - 1; i >= 0; i-- {
		data := datas[i]
		if test(data) {
			return data, true
		}

	}

	var zero Data
	return zero, false
}

func FindLastOrDefault[Data any](datas []Data, test func(Data) bool, defaultData Data) Data {
	if data, ok := FindLast(datas, test); ok {
		return data
	} else {
		return defaultData
	}
}

func FindLastOrZero[Data any](datas []Data, test func(Data) bool) Data {
	data, _ := FindLast(datas, test)
	return data
}

func Index[Data any](datas []Data, i int) (data Data, ok bool) {
	ok = i >= 0 && i < len(datas)
	if ok {
		data = datas[i]
	}
	return
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

func Headmost[Datas ~[]Data, Data any](datas Datas, before func(a, b Data) bool) (result Data) {
	result, _ = HeadmostPro(datas, before)
	return
}

func HeadmostPro[Datas ~[]Data, Data any](datas Datas, before func(a, b Data) bool) (result Data, index int) {
	index = -1
	for i, data := range datas {
		if i == 0 || before(data, result) {
			result = data
			index = i
		}
	}
	return
}

func JoinSlices[Slice ~[]Data, Data any](sep Slice, slices ...Slice) Slice {
	return JoinSlicesTo(nil, sep, false, slices...)
}

func JoinSlicesTo[Slice ~[]Data, Data any](slice Slice, sep Slice, keepFirstSep bool, slices ...Slice) Slice {
	keepSep := keepFirstSep
	return Reduce(slices, func(result Slice, current Slice) Slice {
		if keepSep {
			result = append(result, sep...)
		} else {
			keepSep = true
		}
		return append(result, current...)
	}, slice)
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

// func ConvertSliceItemType[Sources ~[]Source, Destinations ~[]Destination, Source any, Destination any](sources Sources) Destinations {
// 	return Map(sources, func(source Source) Destination {
// 		return source
// 	})
// }

func MapAndConcatWithError[Datas ~[]Data, Results ~[]Result, Data any, Result any](datas Datas, stopWhenError bool, mapper func(Data) (Results, error)) (Results, error) {
	slices, errs := MapWithError(datas, stopWhenError, mapper)
	if err := errs.FirstError(); err != nil {
		return nil, err
	}

	return ConcatSlices(slices...), nil
}

func MapAndJoinWithError[Datas ~[]Data, Results ~[]Result, Data any, Result any](datas Datas, sep Results, stopWhenError bool, mapper func(Data) (Results, error)) (Results, error) {
	slices, errs := MapWithError(datas, stopWhenError, mapper)
	if err := errs.FirstError(); err != nil {
		return nil, err
	}

	return JoinSlices(sep, slices...), nil
}

func MapWithError[Datas ~[]Data, Result any, Data any](datas Datas, stopWhenError bool, mapper func(Data) (Result, error)) (results []Result, errs Errors) {
	// 分配空间
	results = make([]Result, 0, len(datas))
	errs = make(Errors, 0, len(datas))

	for _, data := range datas {
		result, err := mapper(data)

		results = append(results, result)
		errs = append(errs, err)

		if err != nil && stopWhenError {
			return
		}
	}

	return
}

func Reduce[Data any, Datas ~[]Data, Result any](datas Datas, reducer func(Result, Data) Result, initial Result) (result Result) {
	result = initial
	for _, data := range datas {
		result = reducer(result, data)
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

func UniqueSorteds[Data comparable, Datas ~[]Data](datas Datas) Datas {
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

func UniqueByKeySet[Datas ~[]Data, Data any, Key comparable](datas Datas, key func(Data) Key) Datas {
	set := NewSet[Key]()
	return Filter(datas, func(data Data) bool {
		_key := key(data)
		exists := set.Contain(_key)
		if !exists {
			set.Push(_key)
		}
		return !exists
	})
}

func SortUniqueFast[Data constraints.Ordered, Datas ~[]Data](datas Datas) Datas {
	return UniqueSorteds(SortFast(datas))
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
	return Reduce(slices, func(result Slice, current Slice) Slice {
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

func (slice Slice[Data]) AnyMatch(f func(Data) bool) bool {
	return AnyMatch(slice, f)
}

func (slice Slice[Data]) AllMatch(f func(Data) bool) bool {
	return AllMatch(slice, f)
}

func (slice Slice[Data]) Dup() Slice[Data] {
	return DupSlice(slice)
}

func (slice Slice[Data]) Filter(f func(Data) bool) Slice[Data] {
	return Filter(slice, f)
}

func (slice Slice[Data]) FilterPro(f func(int, Data, Slice[Data]) bool) Slice[Data] {
	return FilterPro(slice, f)
}

func (slice Slice[Data]) ForEach(f func(Data)) {
	ForEach(slice, f)
}

func (slice Slice[Data]) ForEachPro(f func(int, Data, Slice[Data])) {
	ForEachPro(slice, f)
}

func (slice Slice[Data]) Map(f func(Data) Data) Slice[Data] {
	return Map(slice, f)
}

func (slice Slice[Data]) MapPro(f func(int, Data, Slice[Data]) Data) Slice[Data] {
	return MapPro(slice, f)
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

func (slice ComparableSlice[Data]) Contain(data Data) bool {
	return IndexOf(slice, data) >= 0
}

func (slice ComparableSlice[Data]) ContainAll(datas ...Data) bool {
	return AllMatch(datas, slice.Contain)
}

func (slice ComparableSlice[Data]) ContainAny(datas ...Data) bool {
	return AnyMatch(datas, slice.Contain)
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
