package reflectx

import "testing"

func TestLen(t *testing.T) {
	tests := []struct {
		data any
		want int
	}{
		{[]int{1, 2, 3}, 3},
		{map[string]int{"a": 1, "b": 2}, 2},
		{[]int{}, 0},
		{map[string]int{}, 0},
		{make(chan int, 10), 0},
		{"abc", 3},
	}
	for _, test := range tests {
		got, ok := Len(test.data)
		if !ok {
			t.Errorf("Len(%v) = false, want true", test.data)
		}
		if got != test.want {
			t.Errorf("Len(%v) = %d, want %d", test.data, got, test.want)
		}
	}
}