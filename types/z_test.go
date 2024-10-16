/*
 * Author: fasion
 * Created time: 2023-04-14 17:48:20
 * Last Modified by: fasion
 * Last Modified time: 2024-10-16 09:51:24
 */

package types

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/stl"
)

func TestSequenceNumber(t *testing.T) {
	no := NewSequenceNumber().Next
	for i := 0; i < 5; i++ {
		fmt.Println(no())
	}
}

func TestXxx(t *testing.T) {
	var strs stl.Slice[string]
	fmt.Println(Strings(strs))
}
