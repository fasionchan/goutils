/*
 * Author: fasion
 * Created time: 2023-05-15 15:59:37
 * Last Modified by: fasion
 * Last Modified time: 2024-11-29 10:41:10
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

func TestReverse(t *testing.T) {
	fmt.Println(NewStrings().ReverseInplace())
	fmt.Println(NewStrings().ReverseDup())

	fmt.Println(NewStrings("a").ReverseInplace())
	fmt.Println(NewStrings("a").ReverseDup())

	fmt.Println(NewStrings("a", "b").ReverseInplace())
	fmt.Println(NewStrings("a", "b").ReverseDup())

	fmt.Println(NewStrings("a", "b", "c").ReverseInplace())
	fmt.Println(NewStrings("a", "b", "c").ReverseDup())
}

func TestCommaSeparatedValueLine(t *testing.T) {
	fmt.Println(CommaSeparatedValueRecord("a,b,c").Values())
	fmt.Println(CommaSeparatedValueRecord(`a,,b,c`).Values())
	fmt.Println(CommaSeparatedValueRecord(`"a,",b,c`).Values())
	fmt.Println(CommaSeparatedValueRecord(`a,b,c,"a,b`).Values())
	fmt.Println(CommaSeparatedValueRecord(`a,b,c,"a,b"`).Values())
}
