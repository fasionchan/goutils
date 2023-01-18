/*
 * Author: fasion
 * Created time: 2023-01-18 15:44:58
 * Last Modified by: fasion
 * Last Modified time: 2023-01-18 15:46:38
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
