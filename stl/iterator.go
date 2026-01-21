/*
 * Author: fasion
 * Created time: 2023-12-21 13:36:40
 * Last Modified by: fasion
 * Last Modified time: 2026-01-21 19:36:32
 */

package stl

import "iter"

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

func MapSeq[Data any, Result any](seq iter.Seq[Data], mapper func(Data) Result) []Result {
	var results []Result
	for data := range seq {
		results = append(results, mapper(data))
	}
	return results
}
