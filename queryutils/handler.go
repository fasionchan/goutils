/*
 * Author: fasion
 * Created time: 2024-06-08 01:00:31
 * Last Modified by: fasion
 * Last Modified time: 2024-06-08 12:36:56
 */

package queryutils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/stl"
)

type TypelessSetinHandler func(ctx context.Context, dataX any, setins []string) error

func (handler TypelessSetinHandler) Setin(ctx context.Context, dataX any, setins []string) error {
	return handler(ctx, dataX, setins)
}

func (handler TypelessSetinHandler) SetinOne(ctx context.Context, dataX any, setin string) error {
	return handler(ctx, dataX, []string{setin})
}

func (handler TypelessSetinHandler) SetinX(ctx context.Context, dataX any, setins ...string) error {
	return handler(ctx, dataX, setins)
}

type SubDataSetiner = TypelessSetinHandlerRegistry

var NewSubDataSetiner = NewTypelessSetinHandlerRegistry

type TypelessSetinHandlerRegistry TypelessSetinHandlerMappingByString

func NewTypelessSetinHandlerRegistry() TypelessSetinHandlerRegistry {
	return TypelessSetinHandlerRegistry{}
}

func (registry TypelessSetinHandlerRegistry) WithHandlerByIdent(handler TypelessSetinHandler, ident string) TypelessSetinHandlerRegistry {
	registry[ident] = handler
	return registry
}

func (registry TypelessSetinHandlerRegistry) WithHandlerByData(handler TypelessSetinHandler, data any) TypelessSetinHandlerRegistry {
	return registry.WithHandlerByIdent(handler, EssentialDataTypeIdent(data))
}

func (registry TypelessSetinHandlerRegistry) HandlerByData(data any) (TypelessSetinHandler, error) {
	ident := EssentialDataTypeIdent(data)
	handler, ok := registry[ident]
	if !ok {
		return nil, baseutils.NewGenericNotFoundError("TypelessSetinHandler", ident)
	}

	return handler, nil
}

func (registry TypelessSetinHandlerRegistry) setin(ctx context.Context, dataX any, setins []string) error {
	for _, setin := range setins {
		if err := registry.setinOne(ctx, dataX, setin); err != nil {
			return err
		}
	}
	return nil
}

func (registry TypelessSetinHandlerRegistry) setinOne(ctx context.Context, dataX any, setin string) error {
Start:
	setin = strings.TrimSpace(setin)

	if setin == "" {
		return nil
	}

	switch setin[0] {
	case '.':
		setin = setin[1:]
		goto Start
		// return registry.setinOne(ctx, dataX, setin[1:])
	case '-':
		return registry.setinOneWithDataHandler(ctx, dataX, setin[1:])
	case '(':
		// todo
		setins, err := parseSetinExpression(setin)
		if err != nil {
			return err
		}

		switch len(setins) {
		case 0:
			return nil
		case 1:
			setin = setins[0]
			goto Start
			// return registry.setinOne(ctx, dataX, setins[0])
		default:
			return registry.setin(ctx, dataX, setins)
		}
	}

	dotIndex := strings.Index(setin, ".")
	minusIndex := strings.Index(setin, "-")
	bracketIndex := strings.Index(setin, "(")

	indexes := stl.PurgeValue([]int{dotIndex, minusIndex, bracketIndex}, -1)
	if len(indexes) == 0 {
		return registry.setinOneWithDataHandler(ctx, dataX, setin)
	}

	index := stl.Min(indexes, 0)

	name := setin[:index]
	setin = setin[index:]

	if setin == "" {
		// todo visit data?
		return nil
	}

	subdata, err := GetSubData(dataX, name)
	if err != nil {
		return err
	}

	dataX = subdata
	goto Start
	// return registry.setinOne(ctx, subdata, setin)
}

func (registry TypelessSetinHandlerRegistry) setinOneWithDataHandler(ctx context.Context, dataX any, setin string) error {
	handler, err := registry.HandlerByData(dataX)
	if err != nil {
		return err
	}
	return handler(ctx, dataX, []string{setin})
}

type TypelessSetinHandlerMappingByString map[string]TypelessSetinHandler

func GetSubData(data any, name string) (any, error) {
	return GetSubDataFromReflectValue(reflect.ValueOf(data), name)
}

func GetSubDataFromReflectValue(value reflect.Value, name string) (any, error) {
	valueKind := value.Kind()

	if valueKind == reflect.Struct {
		field := value.FieldByName(name)
		if field.IsValid() {
			return field.Interface(), nil
		}
	}

	method := value.MethodByName(name)
	if method.IsValid() {
		methodType := method.Type()
		if methodType.NumIn() != 0 {
			return nil, baseutils.NewBadTypeError("TransformMethod", DataTypeIdentFromReflectType(methodType))
		}

		if methodType.NumOut() != 1 {
			return nil, baseutils.NewBadTypeError("TransformMethod", DataTypeIdentFromReflectType(methodType))
		}

		return method.Call(nil)[0].Interface(), nil
	}

	if valueKind == reflect.Pointer {
		if !value.IsNil() {
			return GetSubDataFromReflectValue(value.Elem(), name)
		}
	}

	return nil, baseutils.NewGenericNotFoundError(DataTypeIdentFromReflectType(value.Type()), name)
}

func EssentialDataType(data any) reflect.Type {
	type_ := reflect.TypeOf(data)

	for {
		switch type_.Kind() {
		case reflect.Array, reflect.Slice, reflect.Pointer:
			type_ = type_.Elem()
			continue
		default:
			return type_
		}
	}
}

func DataTypeIdent(data any) string {
	return DataTypeIdentFromReflectValue(reflect.ValueOf(data))
}

func DataTypeIdentFromReflectValue(value reflect.Value) string {
	return DataTypeIdentFromReflectType(value.Type())
}

func DataTypeIdentFromReflectType(type_ reflect.Type) string {
	name := type_.Name()
	if name == "" {
		return type_.String()
	}

	pkgPath := type_.PkgPath()
	if pkgPath == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", pkgPath, name)
}

func EssentialDataTypeIdent(data any) string {
	return DataTypeIdentFromReflectType(EssentialDataType(data))
}

var defaultSubDataSetiner = NewSubDataSetiner()

func GetDefaultSubDataSetiner() SubDataSetiner {
	return defaultSubDataSetiner
}
