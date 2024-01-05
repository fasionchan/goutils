/*
 * Author: fasion
 * Created time: 2023-01-18 15:44:58
 * Last Modified by: fasion
 * Last Modified time: 2024-01-05 11:32:09
 */

package stl

import (
	"fmt"
	"testing"
)

func TestFilterMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(FilterMap(m, func(key string, value int, m map[string]int) bool {
		return value < 3
	}))
}

func TestDemo(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(MapValuesByKeys(m, "a", "b"))
}

func TestMapKeyValuePairs(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(MapKeyValuePairs(m))
}

func TestMapKeyValuePairPtrs(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(MapKeyValuePairPtrs(m))
}

func TestMapMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(MapMap(m, func(key string, value int, m map[string]int) (string, int) {
		return key, value * 2
	}))
}

func TestMapKeysOfNil(t *testing.T) {
	fmt.Println(MapKeys[map[string]string](nil))
}

func TestConcatMaps(t *testing.T) {
	var empty map[string]int
	var a map[string]int = map[string]int{"a": 1, "b": 2}
	var b map[string]int = map[string]int{"c": 3, "b": 1}
	fmt.Println(ConcatMaps(empty))
	fmt.Println(ConcatMaps(empty, empty))
	fmt.Println(ConcatMaps(a, empty))
	fmt.Println(ConcatMaps(a, empty, b))
	fmt.Println(ConcatMaps(a, empty, b, a))
	fmt.Println(a, b)
}
