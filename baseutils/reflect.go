/*
 * Author: fasion
 * Created time: 2023-12-21 14:28:34
 * Last Modified by: fasion
 * Last Modified time: 2023-12-21 17:12:23
 */

package baseutils

import (
	"reflect"

	"github.com/fasionchan/goutils/stl"
)

func DereferenceReflectPointer(ptr reflect.Value) (result reflect.Value, ok bool) {
	for {
		switch ptr.Kind() {
		case reflect.Ptr:
			if ptr.IsNil() {
				return
			}

			ptr = ptr.Elem()
		default:
			return ptr, true
		}
	}
}

func DereferenceReflectPointerWithError(ptr reflect.Value) (result reflect.Value, err error) {
	result, ok := DereferenceReflectPointer(ptr)
	if !ok {
		err = NewBadTypeError("nonpointer", "").WithGivenReflectType(ptr.Type())
	}
	return
}

func MapListToAnysByReflect(datas any) ([]any, error) {
	return MapListByReflect(datas, stl.Echo[any])
}

func MapListByReflect[T any](datas any, mapper func(data any) T) ([]T, error) {
	datasValue := reflect.ValueOf(datas)

	datasValue, err := DereferenceReflectPointerWithError(datasValue)
	if err != nil {
		return nil, err
	}

	switch datasValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
	default:
		return nil, NewBadTypeError("slice/array/string", "").WithGivenReflectType(datasValue.Type())
	}

	return MapListByReflectValue(datasValue, mapper), nil
}

func MapListToAnysByReflectValue(datas reflect.Value) []any {
	return MapListByReflectValue(datas, stl.Echo[any])
}

func MapListByReflectValue[T any](datas reflect.Value, mapper func(data any) T) []T {
	return stl.MapIterator[any](NewSliceIteratorByReflectValue(datas), mapper)
}

func MapKeyValuePairPtrsByReflect(m any) (stl.KeyValuePairPtrs[any, any], error) {
	return MapMapWithReflect(m, stl.Reference[stl.KeyValuePair[any, any]])
}

func MapKeyValuePairPtrsByReflectValue(m reflect.Value) stl.KeyValuePairPtrs[any, any] {
	return MapMapWithReflectValue(m, stl.Reference[stl.KeyValuePair[any, any]])
}

func MapKeyValuePairsByReflect(m any) (stl.KeyValuePairs[any, any], error) {
	return MapMapWithReflect(m, stl.Echo[stl.KeyValuePair[any, any]])
}

func MapKeyValuePairsByReflectValue(m reflect.Value) stl.KeyValuePairs[any, any] {
	return MapMapWithReflectValue(m, stl.Echo[stl.KeyValuePair[any, any]])
}

func MapMapWithReflect[T any](m any, mapper func(stl.KeyValuePair[any, any]) T) ([]T, error) {
	mValue := reflect.ValueOf(m)

	mValue, err := DereferenceReflectPointerWithError(mValue)
	if err != nil {
		return nil, err
	}

	if mValue.Kind() != reflect.Map {
		return nil, NewBlankBadTypeError().WithExpected("map").WithGivenReflectType(mValue.Type())
	}

	return MapMapWithReflectValue(mValue, mapper), nil
}

func MapMapWithReflectValue[T any](m reflect.Value, mapper func(stl.KeyValuePair[any, any]) T) []T {
	return stl.MapIterator[stl.KeyValuePair[any, any]](NewMapIteratorByReflect(m), mapper)
}

func IteratorDataToTypelessSliceByReflect(data any) ([]any, error) {
	return IteratorDataToTypelessSliceByReflectValue(reflect.ValueOf(data))
}

func IteratorDataToTypelessSliceByReflectValue(data reflect.Value) ([]any, error) {
	data, err := DereferenceReflectPointerWithError(data)
	if err != nil {
		return nil, err
	}

	switch data.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		return MapListToAnysByReflectValue(data), nil
	case reflect.Map:
		return MapKeyValuePairsByReflectValue(data).ToTypelessSlice(), nil
	default:
		return nil, NewBlankBadTypeError().WithExpected("slice/array/string/map").WithGivenReflectType(data.Type())
	}
}

type ListIteratorByReflect struct {
	slice reflect.Value
	n     int
	i     int
}

func NewSliceIteratorByReflectValue(slice reflect.Value) *ListIteratorByReflect {
	return &ListIteratorByReflect{
		slice: slice,
		n:     slice.Len(),
		i:     -1,
	}
}

func (it *ListIteratorByReflect) Len() int {
	return it.n
}

func (it *ListIteratorByReflect) Next() bool {
	it.i++
	return it.i < it.n
}

func (it *ListIteratorByReflect) Data() any {
	return it.slice.Index(it.i).Interface()
}

type MapIteratorByReflect struct {
	m    reflect.Value
	n    int
	iter *reflect.MapIter
}

func NewMapIteratorByReflect(m reflect.Value) *MapIteratorByReflect {
	return &MapIteratorByReflect{
		m:    m,
		n:    m.Len(),
		iter: m.MapRange(),
	}
}

func (it *MapIteratorByReflect) Len() int {
	return it.n
}

func (it *MapIteratorByReflect) Next() bool {
	return it.iter.Next()
}

func (it *MapIteratorByReflect) Data() stl.KeyValuePair[any, any] {
	return stl.KeyValuePair[any, any]{
		Key:   it.iter.Key().Interface(),
		Value: it.iter.Value().Interface(),
	}
}
