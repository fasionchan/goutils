/*
 * Author: fasion
 * Created time: 2023-01-18 15:44:58
 * Last Modified by: fasion
 * Last Modified time: 2023-10-23 15:30:57
 */

package stl

import (
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Println(MapValuesByKeys(m, "a", "b"))
}

func TestMapKeysOfNil(t *testing.T) {
	fmt.Println(MapKeys[map[string]string](nil))
}
