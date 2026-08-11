package reflectx

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/fasionchan/goutils/std/jsonx"
	"github.com/fasionchan/goutils/stl"
)

var (
	IntType   = reflect.TypeOf(int(0))
	Int8Type  = reflect.TypeOf(int8(0))
	Int16Type = reflect.TypeOf(int16(0))
	Int32Type = reflect.TypeOf(int32(0))
	Int64Type = reflect.TypeOf(int64(0))

	UintType   = reflect.TypeOf(uint(0))
	Uint8Type  = reflect.TypeOf(uint8(0))
	Uint16Type = reflect.TypeOf(uint16(0))
	Uint32Type = reflect.TypeOf(uint32(0))
	Uint64Type = reflect.TypeOf(uint64(0))

	Float32Type = reflect.TypeOf(float32(0))
	Float64Type = reflect.TypeOf(float64(0))

	StringType = reflect.TypeOf("")
	BytesType = reflect.TypeOf([]byte{})

	ContextType  = TypeOf[context.Context]()
	TimeType     = TypeOf[time.Time]()
	DurationType = TypeOf[time.Duration]()

	TypeMapping = map[string]reflect.Type{
		"int":   IntType,
		"int8":  Int8Type,
		"int16": Int16Type,
		"int32": Int32Type,
		"int64": Int64Type,

		"uint":   UintType,
		"uint8":  Uint8Type,
		"uint16": Uint16Type,
		"uint32": Uint32Type,
		"uint64": Uint64Type,

		"float32": Float32Type,
		"float64": Float64Type,

		"context.Context": ContextType,
		"time.Time":       TimeType,
		"time.Duration":   DurationType,
	}
)

func TypeName(v any) string {
	return reflect.TypeOf(v).Name()
}

func TypeNameOrString(v any) string {
	rt := reflect.TypeOf(v)
	if name := rt.Name(); name != "" {
		return name
	}
	return rt.String()
}

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

type Types []reflect.Type

func (types Types) Len() int {
	return len(types)
}

func (types Types) ResolveValues(mapping ValueMappingByType) []reflect.Value {
	return types.ResolveValuesTo(mapping, make([]reflect.Value, types.Len()))
}

func (types Types) ResolveValuesTo(mapping ValueMappingByType, values []reflect.Value) []reflect.Value {
	n := min(
		len(types),
		len(values),
	)

	for i := range n {
		if value, ok := mapping[types[i]]; ok {
			values[i] = value
		} else {
			return values[:i]
		}
	}

	return values[:n]
}

func (types Types) ParseJsonRawMessages(args jsonx.JsonRawMessages) ([]reflect.Value, error) {
	return types.ParseJsonRawMessagesTo(args, make([]reflect.Value, types.Len()))
}

func (types Types) ParseJsonRawMessagesTo(msgs jsonx.JsonRawMessages, values []reflect.Value) ([]reflect.Value, error) {
	n := min(
		len(types),
		len(msgs),
		len(values),
	)

	for i := range n {
		value := reflect.New(types[i])
		if err := json.Unmarshal(msgs[i], value.Interface()); err != nil {
			return nil, err
		}
		values[i] = value.Elem()
	}

	return values[:n], nil
}

func (types Types) Zeros() []reflect.Value {
	return stl.Map(types, reflect.Zero)
}

func (types Types) ZerosTo(values []reflect.Value) []reflect.Value {
	n := min(
		len(types),
		len(values),
	)

	for i := range n {
		values[i] = reflect.Zero(types[i])
	}

	return values[:n]
}
