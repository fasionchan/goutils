/*
 * Author: fasion
 * Created time: 2023-03-24 08:46:11
 * Last Modified by: fasion
 * Last Modified time: 2024-06-25 15:20:08
 */

package stl

import "reflect"

func New[Data any]() (data Data) {
	return
}

func NewPtr[Ptr ~*Data, Data any]() Ptr {
	var data Data
	return &data
}

func Dup[Data any, Ptr ~*Data](ptr Ptr) Ptr {
	if ptr == nil {
		return nil
	}

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

func ReflectType[T any]() reflect.Type {
	var value T
	return reflect.TypeOf(value)
}
