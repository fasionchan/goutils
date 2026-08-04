package reflectx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	tests := []struct {
		data any
		ok bool
		want int
	}{
		{[]int{1, 2, 3}, true, 3},
		{map[string]int{"a": 1, "b": 2}, true, 2},
		{[]int{}, true, 0},
		{map[string]int{}, true, 0},
		{make(chan int, 10), true, 0},
		{"abc", true, 3},

		{nil, false, 0},
		{1, false, 0},
		{(*[]byte)(nil), false, 0},
		{struct{}{}, false, 0},
	}
	for _, test := range tests {
		got, ok := Len(test.data)
		assert.Equal(t, test.ok, ok)
		assert.Equal(t, test.want, got)
	}
}