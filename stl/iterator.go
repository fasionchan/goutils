/*
 * Author: fasion
 * Created time: 2023-12-21 13:36:40
 * Last Modified by: fasion
 * Last Modified time: 2023-12-21 15:06:54
 */

package stl

type Iterator[Data any] interface {
	Len() int
	Data() Data
	Next() bool
}

func MapIterator[Data any, Result any](it Iterator[Data], mapper func(Data) Result) []Result {
	var result []Result
	if n := it.Len(); n > 0 {
		result = make([]Result, 0, n)
	}

	for it.Next() {
		result = append(result, mapper(it.Data()))
	}

	return result
}
