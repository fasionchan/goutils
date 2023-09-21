/*
 * Author: fasion
 * Created time: 2023-09-21 11:14:51
 * Last Modified by: fasion
 * Last Modified time: 2023-09-21 11:33:21
 */

package stl

import "golang.org/x/exp/constraints"

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

type Addend interface {
	constraints.Integer | constraints.Float
}

func Sum[Datas ~[]Data, Data Addend](datas Datas, start Data) Data {
	for _, data := range datas {
		start += data
	}
	return start
}
