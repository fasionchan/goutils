package reflectx

import (
	"reflect"

	"github.com/fasionchan/goutils/stl"
)

var (
	inValidValue reflect.Value
)

type ValuePtr = *Value

type Value reflect.Value

func NewValue(value *reflect.Value) *Value {
	return (*Value)(value)
}

func AddrOfValue(value reflect.Value) *Value {
	return NewValue(&value)
}

func (value *Value) CanInterface() bool {
	return value.NativeValue().CanInterface()
}

func (value *Value) IsValid() bool {
	return value.NativeValue().IsValid()
}

func (value *Value) Interface() any {
	return value.NativeValue().Interface()
}

func (value *Value) Native() *reflect.Value {
	return (*reflect.Value)(value)
}

func (value *Value) NativeValue() reflect.Value {
	if value == nil {
		return inValidValue
	}

	return *value.Native()
}

func (value *Value) Extract(path PathItems) (Values, error) {
	if path.Empty() {
		return NewValues(value), nil
	}

	return extractValues(value.NativeValue(), path)
}

func (value *Value) Index(index *DataIndex) *Value {
	if value == nil {
		return nil
	}

	if index == nil {
		return nil
	}

	return value.index(index)
}

func (value *Value) IndexForValues(index *DataIndex) Values {
	if value == nil {
		return nil
	}

	if index == nil {
		return nil
	}

	value = value.index(index)
	switch value.Native().Kind() {
	case reflect.Slice, reflect.Array:
		return ItemValuesFromIndexable(*value.Native())
	default:
		return NewValues(value)
	}
}

func (value *Value) index(index *DataIndex) *Value {
	if ranges := index.ranges; ranges != nil {
		return value.slice(ranges[0], ranges[1])
	} else {
		return value.uniIndex(index.i)
	}
}

func (value *Value) UniIndex(index int) *Value {
	if value == nil {
		return nil
	}

	return value.uniIndex(index)
}

func (value *Value) uniIndex(index int) *Value {
	if index < 0 {
		index += value.Native().Len()
	}

	result := value.Native().Index(index)
	return (*Value)(&result)
}

func (value *Value) Slice(start, end *int) *Value {
	if value == nil {
		return nil
	}

	return value.slice(start, end)
}

func (value *Value) slice(start, end *int) *Value {
	istart, iend := 0, value.Native().Len()
	if start != nil {
		istart = *start
	}
	if end != nil {
		iend = *end
	}

	result := value.Native().Slice(istart, iend)
	return (*Value)(&result)
}

type Values []*Value

func NewValues(values ...*Value) Values {
	return Values(values)
}

func ValuesOf(values ...any) Values {
	return ValueInstancesOf(values...).Values()
}

func (values Values) Append(others ...*Value) Values {
	return append(values, others...)
}

func (values Values) AppendNatives(others ...reflect.Value) Values {
	for _, other := range others {
		values = values.Append(NewValue(&other))
	}
	return values
}

func ItemValuesFromIndexable(value reflect.Value) Values {
	n := value.Len()
	items := make(Values, 0, n)
	for i := range n {
		items = items.AppendNatives(value.Index(i))
	}
	return items
}


func (values Values) Empty() bool {
	return len(values) == 0
}

func (values Values) Len() int {
	return len(values)
}

func (values Values) Valids() Values {
	return stl.Filter(values, ValuePtr.IsValid)
}

func (values Values) Interfacables() Values {
	return stl.Filter(values, ValuePtr.CanInterface)
}

func (values Values) Interfaces() []any {
	return stl.Map(values, ValuePtr.Interface)
}

func (values Values) Extract(path string) (Values, error) {
	items, err := ParsePathExpand(path)
	if err != nil {
		return nil, err
	}
	return values.ExtractPath(items)
}

func (values Values) ExtractPath(path PathItems) (Values, error) {
	if values.Empty() {
		return nil, nil
	}

	if path.Empty() {
		return values, nil
	}

	if values.Len() == 1 {
		return values[0].Extract(path)
	}

	return stl.MapAndConcatWithError(values, true, func(value ValuePtr) (Values, error) {
		return value.Extract(path)
	})
}

type ValueInstances []reflect.Value 

func (instances ValueInstances) Values() Values {
	return stl.Map(instances, AddrOfValue)
}

func ValueInstancesOf(args ...any) ValueInstances {
	return stl.Map(args, reflect.ValueOf)
}

type ValueMappingByType = map[reflect.Type]reflect.Value