package reflectx

import (
	"fmt"
	"reflect"
)

// extractValues 为核心实现：按 PathItems 从 reflect.Value 中提取，返回 Values。
func extractValues(v reflect.Value, path PathItems) (Values, error) {
	if len(path) == 0 {
		return NewValues(NewValue(&v)), nil
	}

	// 解引用指针
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	item := path[0]
	rest := path[1:]

	switch v.Kind() {
	case reflect.Struct:
		if item.Type != PathItemTypeAttr {
			return nil, fmt.Errorf("invalid path item(not attr): %s", item.String())
		}

		f := v.FieldByName(item.Name)
		if !f.IsValid() {
			return nil, nil
		}

		return extractValues(f, rest)
	case reflect.Map:
		if item.Type != PathItemTypeAttr {
			return nil, fmt.Errorf("invalid path item(not attr): %s", item.String())
		}

		k := reflect.ValueOf(item.Name)
		r := v.MapIndex(k)
		if !r.IsValid() {
			return nil, nil
		}

		return extractValues(r, rest)
	case reflect.Slice, reflect.Array:
		switch item.Type {
		case PathItemTypeIndex:
			values := NewValue(&v).IndexForValues(item.Index)
			return values.ExtractPath(rest)
		case PathItemTypeAttr:
			values := ItemValuesFromIndexable(v)
			return values.ExtractPath(path)
		default:
			return nil, fmt.Errorf("invalid path item(not index or attr): %s", item.String())
		}
	default:
		return nil, nil
	}
}

// ExtractValues 从 reflect.Value 按字符串路径提取，返回 Values。
// 路径格式见 ParsePath。
func ExtractValues(value reflect.Value, path string) (Values, error) {
	items, err := ParsePathExpand(path)
	if err != nil {
		return nil, err
	}
	return ExtractValuesPath(value, items)
}

// ExtractValuesPath 从 reflect.Value 按已解析的 PathItems 提取，返回 Values。
func ExtractValuesPath(value reflect.Value, path PathItems) (Values, error) {
	if len(path) == 0 {
		return nil, nil
	}
	return extractValues(value, path)
}

// Extract 从任意值按字符串路径提取，返回 []any。
func Extract(value any, path string) ([]any, error) {
	items, err := ParsePathExpand(path)
	if err != nil {
		return nil, err
	}
	return ExtractPath(value, items)
}

// ExtractPath 从任意值按已解析的 PathItems 提取，返回 []any。
func ExtractPath(value any, path PathItems) ([]any, error) {
	v := reflect.ValueOf(value)
	vals, err := ExtractValuesPath(v, path)
	if err != nil {
		return nil, err
	}

	return vals.Valids().Interfacables().Interfaces(), nil
}
