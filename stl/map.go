/*
 * Author: fasion
 * Created time: 2022-11-19 17:43:35
 * Last Modified by: fasion
 * Last Modified time: 2022-11-19 20:37:45
 */

package stl

func MapKeys[Key comparable, Value any, Map ~map[Key]Value](m Map) []Key {
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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
