package reflectx

import (
	"reflect"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

type StructFields []reflect.StructField

func ParseStructFields(structType reflect.Type) StructFields {
	return stl.CountAndMap(0, structType.NumField(), 1, structType.Field)
}

func (fields StructFields) MappingByName() StructFieldMappingByString {
	return stl.BuildMap[StructFieldMappingByString](fields, func(field reflect.StructField) (string, reflect.StructField) {
		return field.Name, field
	})
}

func (fields StructFields) MappingTag(tag string) StructFieldMappingByString {
	return stl.MappingByKeys(fields, func(field reflect.StructField) types.Strings {
		return types.SplitToStrings(field.Tag.Get(tag), ",")
	})
}

type StructFieldMappingByString map[string]reflect.StructField