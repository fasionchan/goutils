/*
 * Author: fasion
 * Created time: 2023-03-24 08:46:11
 * Last Modified by: fasion
 * Last Modified time: 2023-12-21 16:51:40
 */

package stl

func Dup[Data any, Ptr ~*Data](ptr Ptr) Ptr {
	dup := *ptr
	return &dup
}

func Echo[Data any](data Data) Data {
	return data
}

func Reference[Data any](data Data) *Data {
	return &data
}

func Dereference[Ptr ~*Data, Data any](ptr Ptr) Data {
	return *ptr
}

func ToTypeless[Data any](data Data) any {
	return data
}

func ToTypelessSlice[Datas ~[]Data, Data any](datas Datas) []any {
	return Map(datas, ToTypeless[Data])
}
