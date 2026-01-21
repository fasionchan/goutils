/*
 * Author: fasion
 * Created time: 2026-01-21 19:25:02
 * Last Modified by: fasion
 * Last Modified time: 2026-01-21 19:45:59
 */

package stl

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMapSeq(t *testing.T) {
	number := reflect.ValueOf(10)
	fmt.Println(number.Type().CanSeq())

	results := MapSeq(number.Seq(), func(data reflect.Value) int {
		return data.Interface().(int) * 2
	})
	fmt.Println(results)
}
