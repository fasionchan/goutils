/*
 * Author: fasion
 * Created time: 2023-04-14 17:44:10
 * Last Modified by: fasion
 * Last Modified time: 2026-04-11 19:51:18
 */

package stl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReslice(t *testing.T) {
	bytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	a := bytes[:5]
	b := bytes[5:]
	fmt.Println(a[5:10], b)
	assert.Equal(t, a[5:10], b)
}

func TestCompiling(t *testing.T) {

}
