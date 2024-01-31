/*
 * Author: fasion
 * Created time: 2022-12-25 18:41:34
 * Last Modified by: fasion
 * Last Modified time: 2024-01-17 10:46:02
 */

package baseutils

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	value := reflect.ValueOf(NotImplementedError.Error)
	fmt.Println(value)
	fmt.Println(value.Kind())
	fmt.Println(value.Type())
	fmt.Println(value.Type().Name())
	fmt.Println(value.Type().PkgPath())
}
