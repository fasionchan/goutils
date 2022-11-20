/*
 * Author: fasion
 * Created time: 2022-11-14 11:27:56
 * Last Modified by: fasion
 * Last Modified time: 2022-11-19 21:43:57
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

func DupSlice[Data any, Slice ~[]Data](slice Slice) Slice {
	return append(make(Slice, 0, len(slice)), slice...)
}

func ConcatSlices[Data any, Slice ~[]Data](slices ...Slice) Slice {
	return ConcatSlicesTo(nil, slices...)
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

func SliceMappingByKey[Data any, Datas ~[]Data, Key comparable](datas Datas, key func(Data) Key) map[Key]Datas {
	m := map[Key]Datas{}
	for _, data := range datas {
		k := key(data)
		m[k] = append(m[k], data)
	}
	return m
}
