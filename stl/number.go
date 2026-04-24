package stl

import "golang.org/x/exp/constraints"

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

type Addend interface {
	constraints.Integer | constraints.Float
}

func Sum[Datas ~[]Data, Data Addend](datas Datas, start Data) Data {
	for _, data := range datas {
		start += data
	}
	return start
}
