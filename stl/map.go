/*
 * Author: fasion
 * Created time: 2022-11-19 17:43:35
 * Last Modified by: fasion
 * Last Modified time: 2023-12-13 10:58:42
 */

package stl

type KeyValuePair[Key any, Value any] struct {
	Key   Key
	Value Value
}

type KeyValuePairs[Key any, Value any] []KeyValuePair[Key, Value]

func MapKeyValuePairs[Map ~map[Key]Value, Key comparable, Value any](m Map) KeyValuePairs[Key, Value] {
	return MapMapToSlice[KeyValuePairs[Key, Value]](m, func(key Key, value Value, m Map) KeyValuePair[Key, Value] {
		return KeyValuePair[Key, Value]{
			Key:   key,
			Value: value,
		}
	})
}

type KeyValuePairPtrs[Key any, Value any] []*KeyValuePair[Key, Value]

func MapKeyValuePairPtrs[Map ~map[Key]Value, Key comparable, Value any](m Map) KeyValuePairPtrs[Key, Value] {
	return MapMapToSlice[KeyValuePairPtrs[Key, Value]](m, func(key Key, value Value, m Map) *KeyValuePair[Key, Value] {
		return &KeyValuePair[Key, Value]{
			Key:   key,
			Value: value,
		}
	})
}

func MapMap[Map ~map[Key]Value, Key comparable, Value any](m Map, mapper func(Key, Value, Map) (Key, Value)) Map {
	result := Map{}
	for key, value := range m {
		key, value = mapper(key, value, m)
		result[key] = value
	}
	return result
}

func MapMapToSlice[Slice ~[]SliceItem, Map ~map[Key]Value, SliceItem any, Key comparable, Value any](m Map, convert func(Key, Value, Map) SliceItem) Slice {
	result := make(Slice, 0, len(m))

	for key, value := range m {
		result = append(result, convert(key, value, m))
	}

	return result
}

func MapMapToSlicePro[Slice ~[]SliceItem, Map ~map[Key]Value, SliceItem any, Key comparable, Value any](m Map, convert func(Key, Value, Map) (SliceItem, bool, error)) (Slice, error) {
	result := make(Slice, 0, len(m))

	for key, value := range m {
		item, ok, err := convert(key, value, m)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}

func BuildMap[Datas ~[]Data, Map ~map[Key]Value, Data any, Key comparable, Value any](datas Datas, kv func(data Data) (Key, Value)) Map {
	result := Map{}
	for _, data := range datas {
		key, value := kv(data)
		result[key] = value
	}
	return result
}

func BuildMapPro[Datas ~[]Data, Map ~map[Key]Value, Data any, Key comparable, Value any](datas Datas, kv func(data Data) (Key, Value, bool, error)) (Map, error) {
	result := Map{}

	for _, data := range datas {
		key, value, ok, err := kv(data)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		result[key] = value
	}

	return result, nil
}

func MapKeys[Map ~map[Key]Value, Key comparable, Value any](m Map) []Key {
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[Key comparable, Value any, Map ~map[Key]Value](m Map) []Value {
	values := make([]Value, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func MapValuesByKeys[Key comparable, Value any, Map ~map[Key]Value](m Map, keys ...Key) []Value {
	values := make([]Value, 0, len(keys))
	for _, key := range keys {
		values = append(values, m[key])
	}
	return values
}

func DupMap[Key comparable, Value any, Map ~map[Key]Value](m Map) Map {
	dup := Map{}
	for key, value := range m {
		dup[key] = value
	}
	return dup
}

func ConcatMapInplace[Key comparable, Value any, Map ~map[Key]Value](m1, m2 Map) Map {
	for key, value := range m2 {
		m1[key] = value
	}
	return m1
}

func ConcatMap[Key comparable, Value any, Map ~map[Key]Value](m1, m2 Map) Map {
	return ConcatMap(DupMap(m1), m2)
}

func PopMap[Key comparable, Value any, Map ~map[Key]Value](m Map, key Key) (value Value, ok bool) {
	value, ok = m[key]
	if ok {
		delete(m, key)
	}
	return
}

func BatchDeleteMap[Key comparable, Keys ~[]Key, Value any, Map ~map[Key]Value](m Map, keys Keys) Map {
	for _, key := range keys {
		delete(m, key)
	}
	return m
}

func BatchDeleteMapFromAnother[Key comparable, Value any, Map ~map[Key]Value](m Map, keys Map) Map {
	for key := range keys {
		delete(m, key)
	}
	return m
}

func CacheMapValueWithInitializer[Key comparable, Value any, Map ~map[Key]Value](m Map, key Key, initializer func() Value) Value {
	cached, ok := m[key]
	if !ok {
		cached = initializer()
		m[key] = cached
	}
	return cached
}

func CacheMapValue[Key comparable, Value any, Map ~map[Key]Value](m Map, key Key, value Value) Value {
	cached, ok := m[key]
	if !ok {
		cached = value
		m[key] = cached
	}
	return cached
}

func SubMapByKeys[Key comparable, Value any, Map ~map[Key]Value](m Map, keys ...Key) Map {
	result := Map{}
	for _, key := range keys {
		if value, ok := m[key]; ok {
			result[key] = value
		}
	}
	return result
}

func MapValueGetter[Key comparable, Value any, Map ~map[Key]Value](m Map) func(Key) Value {
	return func(k Key) Value {
		return m[k]
	}
}

func MapValueGetterPro[Key comparable, Value any, Map ~map[Key]Value](m Map, keys ...Key) func(Key) (Value, bool) {
	return func(k Key) (v Value, ok bool) {
		v, ok = m[k]
		return
	}
}
