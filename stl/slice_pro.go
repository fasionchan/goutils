/*
 * Author: fasion
 * Created time: 2022-11-19 17:24:48
 * Last Modified by: fasion
 * Last Modified time: 2023-12-21 11:51:21
 */

package stl

func AnyMatchPro[Data any, Datas ~[]Data](datas Datas, test func(int, Data, Datas) bool) bool {
	for i, data := range datas {
		if test(i, data, datas) {
			return true
		}
	}
	return false
}

func AllMatchPro[Data any, Datas ~[]Data](datas Datas, test func(int, Data, Datas) bool) bool {
	for i, data := range datas {
		if !test(i, data, datas) {
			return false
		}
	}
	return true
}

func ForEachPro[Data any, Datas ~[]Data](datas Datas, handler func(i int, data Data, datas Datas)) {
	for i, data := range datas {
		handler(i, data, datas)
	}
}

func FilterPro[Data any, Datas ~[]Data](datas Datas, filter func(int, Data, Datas) bool) Datas {
	result := make(Datas, 0, len(datas))
	for i, data := range datas {
		if filter(i, data, datas) {
			result = append(result, data)
		}
	}
	return result
}

func MapPro[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(int, Data, Datas) Result) []Result {
	results := make([]Result, 0, len(datas))
	for i, data := range datas {
		results = append(results, mapper(i, data, datas))
	}
	return results
}

func MapProArgs[Data any, Datas ~[]Data, Result any](mapper func(int, Data, Datas) Result, datas ...Data) []Result {
	return MapPro(Datas(datas), mapper)
}

func MapWithErrorPro[Datas ~[]Data, Result any, Data any](datas Datas, stopWhenError bool, mapper func(int, Data, Datas) (Result, error)) (results []Result, errs Errors) {
	// 分配空间
	results = make([]Result, 0, len(datas))
	errs = make(Errors, 0, len(datas))

	for i, data := range datas {
		result, err := mapper(i, data, datas)

		results = append(results, result)
		errs = append(errs, err)

		if err != nil && stopWhenError {
			return
		}
	}

	return
}

func ReducePro[Data any, Datas ~[]Data, Result any](datas Datas, reducer func(Result, Data, int, Datas) Result, initial Result) (result Result) {
	result = initial
	for i, data := range datas {
		result = reducer(result, data, i, datas)
	}
	return
}

func MappingByKeyPro[Data any, Datas ~[]Data, Key comparable](datas Datas, key func(int, Data, Datas) Key) map[Key]Data {
	m := map[Key]Data{}
	for i, data := range datas {
		m[key(i, data, datas)] = data
	}
	return m
}

func InstancesToPointers[Data any, Instances ~[]Data](instances Instances) []*Data {
	return MapPro(instances, func(i int, _ Data, datas Instances) *Data {
		return &datas[i]
	})
}
