/*
 * Author: fasion
 * Created time: 2022-11-19 17:24:48
 * Last Modified by: fasion
 * Last Modified time: 2022-11-19 17:35:21
 */

package stl

func AnyMatchCustom[Data any, Datas ~[]Data](datas Datas, test func(int, Data, Datas) bool) bool {
	for i, data := range datas {
		if test(i, data, datas) {
			return true
		}
	}
	return false
}

func AllMatchCustom[Data any, Datas ~[]Data](datas Datas, test func(int, Data, Datas) bool) bool {
	for i, data := range datas {
		if !test(i, data, datas) {
			return false
		}
	}
	return true
}

func ForEachCustom[Data any, Datas ~[]Data](datas Datas, handler func(i int, data Data, datas Datas)) {
	for i, data := range datas {
		handler(i, data, datas)
	}
}

func FilterCustom[Data any, Datas ~[]Data](datas Datas, filter func(int, Data, Datas) bool) Datas {
	result := make(Datas, 0, len(datas))
	for i, data := range datas {
		if filter(i, data, datas) {
			result = append(result, data)
		}
	}
	return result
}

func MapCustom[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(int, Data, Datas) Result) []Result {
	results := make([]Result, 0, len(datas))
	for i, data := range datas {
		results = append(results, mapper(i, data, datas))
	}
	return results
}

func MapCustomArgs[Data any, Datas ~[]Data, Result any](mapper func(int, Data, Datas) Result, datas ...Data) []Result {
	return MapCustom(Datas(datas), mapper)
}

func ReduceCustom[Data any, Datas ~[]Data, Result any](datas Datas, reducer func(int, Data, Datas, Result) Result, initial Result) (result Result) {
	result = initial
	for i, data := range datas {
		result = reducer(i, data, datas, result)
	}
	return
}

func MappingByKeyCustom[Data any, Datas ~[]Data, Key comparable](datas Datas, key func(int, Data, Datas) Key) map[Key]Data {
	m := map[Key]Data{}
	for i, data := range datas {
		m[key(i, data, datas)] = data
	}
	return m
}

func InstancesToPointers[Data any, Instances ~[]Data](instances Instances) []*Data {
	return MapCustom(instances, func(i int, _ Data, datas Instances) *Data {
		return &datas[i]
	})
}
