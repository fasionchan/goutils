/*
 * Author: fasion
 * Created time: 2023-05-15 15:59:37
 * Last Modified by: fasion
 * Last Modified time: 2024-08-07 10:43:28
 */

package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, NewStrings("a\tb c", "d e", "f g\th\ni\rj").Split(" ", "\t", "\n", "\r").Native(), []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"})
}

func TestCountStrings(t *testing.T) {
	fmt.Println(NewStrings("a", "b", "c", "a", "c", "a").Count())
}
