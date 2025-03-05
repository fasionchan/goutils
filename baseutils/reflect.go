/*
 * Author: fasion
 * Created time: 2023-12-21 14:28:34
 * Last Modified by: fasion
 * Last Modified time: 2025-03-05 15:06:26
 */

package baseutils

import (
	"reflect"

	"fmt"

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

// ReflectMap traverses elements in datas, calls the specified method on each element,
// and returns the results as a slice.
//
// Parameters:
// - datas: a slice or array of data to process
// - methodName: the name of the method to call on each element
// - methodOutValueIndex: the index of the value returned by the method
// - methodOutErrorIndex: the index of the error returned by the method
// - args: additional arguments to pass to the method (can be empty)
// - stopWhenError: whether to stop processing when an error occurs
//
// Returns:
// - values: a slice containing the values returned by the method calls
// - error: the first error encountered, or nil if no errors occurred
func ReflectMap(datas any, methodName string, methodOutValueIndex, methodOutErrorIndex int, args []reflect.Value, stopWhenError bool) (values any, errs Errors, err error) {
	datasValue := reflect.ValueOf(datas)

	// Dereference pointer if necessary
	datasValue, _ = DereferenceReflectPointer(datasValue)

	// Check if datas is a slice or array
	switch datasValue.Kind() {
	case reflect.Array, reflect.Slice:
		// Valid types
	default:
		return nil, nil, NewBlankBadTypeError().WithExpected("slice/array").WithGivenReflectType(datasValue.Type())
	}

	dataType := datasValue.Type().Elem()
	methodType, ok := dataType.MethodByName(methodName)
	if !ok {
		return nil, nil, NewGenericNotFoundError(methodName, dataType.String())
	}

	funcType := methodType.Func.Type()
	funcNumOut := funcType.NumOut()
	if methodOutValueIndex >= funcNumOut || methodOutErrorIndex >= funcNumOut {
		return nil, nil, fmt.Errorf("method %s on type %s has wrong return values: %d", methodName, dataType.String(), funcNumOut)
	}

	// Check if the error type is correct
	if methodOutErrorIndex >= 0 {
		errorType := funcType.Out(methodOutErrorIndex)
		if errorType != stl.ReflectType[error]() {
			return nil, nil, NewBlankBadTypeError().WithExpected("error").WithGivenReflectType(errorType)
		}
	}

	length := datasValue.Len()

	// Initialize result slice
	var valuesValue reflect.Value
	if methodOutValueIndex >= 0 {
		valueType := funcType.Out(methodOutValueIndex)
		valuesValue = reflect.MakeSlice(reflect.SliceOf(valueType), 0, length)
	}

	args = append(make([]reflect.Value, 1, len(args)+1), args...)

	// Process each element
	for i := 0; i < length; i++ {
		elem := datasValue.Index(i)

		// Call the method with provided arguments
		args[0] = elem
		results := methodType.Func.Call(args)

		if methodOutValueIndex >= 0 {
			// Get the value
			value := results[methodOutValueIndex]
			valuesValue = reflect.Append(valuesValue, value)
		}

		if methodOutErrorIndex >= 0 {
			err = results[methodOutErrorIndex].Interface().(error)
			errs = append(errs, err)
			if err != nil {
				break
			}
		}
	}

	values = valuesValue.Interface()

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
