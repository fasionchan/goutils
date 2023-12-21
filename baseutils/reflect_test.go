/*
 * Author: fasion
 * Created time: 2023-12-21 14:32:53
 * Last Modified by: fasion
 * Last Modified time: 2023-12-21 17:25:48
 */

package baseutils

import (
	"fmt"
	"testing"
)

var nilInts *[]int
var nilStruct *struct{}

func TestMapListToAnysByReflect(t *testing.T) {

	fmt.Println(MapListToAnysByReflect([]int(nil)))
	fmt.Println()

	fmt.Println(MapListToAnysByReflect([]int{}))
	fmt.Println(MapListToAnysByReflect([]int{0}))
	fmt.Println(MapListToAnysByReflect([]int{0, 1, 2}))
	fmt.Println()

	fmt.Println(MapListToAnysByReflect([]string{}))
	fmt.Println(MapListToAnysByReflect([]string{""}))
	fmt.Println(MapListToAnysByReflect([]string{"abc"}))
	fmt.Println(MapListToAnysByReflect([]string{"abc", ""}))
	fmt.Println()

	fmt.Println(MapListToAnysByReflect(""))
	fmt.Println(MapListToAnysByReflect("a"))
	fmt.Println(MapListToAnysByReflect(" "))
	fmt.Println(MapListToAnysByReflect("abc"))
	fmt.Println()

	fmt.Println(MapListToAnysByReflect([0]int{}))
	fmt.Println(MapListToAnysByReflect([1]int{}))
	fmt.Println(MapListToAnysByReflect([1]int{1}))
	fmt.Println(MapListToAnysByReflect([2]int{1, 2}))
}

func TestIteratorDataToTypelessSliceByReflect(t *testing.T) {
	fmt.Println(IteratorDataToTypelessSliceByReflect(nilInts))
	fmt.Println(IteratorDataToTypelessSliceByReflect(nilStruct))

	fmt.Println(IteratorDataToTypelessSliceByReflect([]int(nil)))
	fmt.Println()

	fmt.Println(IteratorDataToTypelessSliceByReflect([]int{}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([]int{0}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([]int{0, 1, 2}))
	fmt.Println()

	fmt.Println(IteratorDataToTypelessSliceByReflect([]string{}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([]string{""}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([]string{"abc"}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([]string{"abc", ""}))
	fmt.Println()

	fmt.Println(IteratorDataToTypelessSliceByReflect(""))
	fmt.Println(IteratorDataToTypelessSliceByReflect("a"))
	fmt.Println(IteratorDataToTypelessSliceByReflect(" "))
	fmt.Println(IteratorDataToTypelessSliceByReflect("abc"))
	fmt.Println()

	fmt.Println(IteratorDataToTypelessSliceByReflect([0]int{}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([1]int{}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([1]int{1}))
	fmt.Println(IteratorDataToTypelessSliceByReflect([2]int{1, 2}))
	fmt.Println()

	fmt.Println(IteratorDataToTypelessSliceByReflect(map[string]int(nil)))
	fmt.Println(IteratorDataToTypelessSliceByReflect(map[string]int{}))
	fmt.Println(IteratorDataToTypelessSliceByReflect(map[string]int{
		"a": 1,
	}))
	fmt.Println(IteratorDataToTypelessSliceByReflect(map[string]int{
		"a": 1,
		"b": 1,
		"c": 1,
	}))
	fmt.Println()
}
