/*
 * Author: fasion
 * Created time: 2023-03-24 08:46:11
 * Last Modified by: fasion
 * Last Modified time: 2025-01-02 11:05:13
 */

package stl

import (
	"context"
	"fmt"
	"reflect"
)

type ContextValues[Key comparable, Value any] map[Key]Value

func NewContextValues[Key comparable, Value any]() ContextValues[Key, Value] {
	return ContextValues[Key, Value]{}
}

func (values ContextValues[Key, Value]) ApplyTo(ctx context.Context) context.Context {
	if values.Empty() {
		return ctx
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for key, value := range values {
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx
}

func (values ContextValues[Key, Value]) Empty() bool {
	return values.Len() == 0
}

func (values ContextValues[Key, Value]) Len() int {
	return len(values)
}

func (values ContextValues[Key, Value]) With(key Key, value Value) ContextValues[Key, Value] {
	if values != nil {
		values[key] = value
	}
	return values
}

func LookupContextValue[Value any, Key any](ctx context.Context, key Key) (value Value, ok bool) {
	if ctx != nil {
		value, ok = ctx.Value(key).(Value)
	}
	return
}

func MustLookupContextValue[Value any, Key any](ctx context.Context, key Key) (value Value) {
	value, ok := LookupContextValue[Value](ctx, key)
	if !ok {
		panic(fmt.Sprintf("context value not found: %v", key))
	}
	return
}

func GetContextValue[Value any, Key any](ctx context.Context, key Key) (value Value) {
	value, _ = LookupContextValue[Value](ctx, key)
	return
}

func New[Data any]() (data Data) {
	return
}

func NewAsAny[T any]() any {
	return new(T)
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

func TypeAsserter[Dst any, Src any](src Src) (dst Dst, ok bool) {
	dst, ok = any(src).(Dst)
	return
}

func ReflectType[T any]() reflect.Type {
	var value T
	return reflect.TypeOf(value)
}

func Zero[Data any]() (_ Data) {
	return
}
