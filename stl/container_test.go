package stl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackBasic(t *testing.T) {
	s := NewStack[int](2)
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Len())

	s.Push(1)
	s.Push(2)
	s.Push(3)
	assert.False(t, s.IsEmpty())
	assert.Equal(t, 3, s.Len())
	assert.Equal(t, 3, s.Peek())

	assert.Equal(t, 3, s.Pop())
	assert.Equal(t, 2, s.Pop())
	assert.Equal(t, 1, s.Peek())
	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestStackTryMethodsAndClear(t *testing.T) {
	s := NewStack[string](0)

	data, ok := s.TryPeek()
	assert.False(t, ok)
	assert.Equal(t, "", data)

	data, ok = s.TryPop()
	assert.False(t, ok)
	assert.Equal(t, "", data)
	assert.Equal(t, "", s.Peek())
	assert.Equal(t, "", s.Pop())

	s.Push("a")
	s.Push("b")
	data, ok = s.TryPeek()
	assert.True(t, ok)
	assert.Equal(t, "b", data)

	data, ok = s.TryPop()
	assert.True(t, ok)
	assert.Equal(t, "b", data)
	assert.Equal(t, 1, s.Len())

	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Len())
}

func TestStackAsContainer(t *testing.T) {
	var c Container[int] = NewStackAsContainer[int](1)
	c.Push(10)
	c.Push(20)
	assert.Equal(t, 20, c.Peek())
	assert.Equal(t, 20, c.Pop())
	assert.Equal(t, 10, c.Pop())
	assert.True(t, c.IsEmpty())
	c.Clear()
}

func TestMinHeapOrder(t *testing.T) {
	h := NewMinHeap[int](0)
	for _, v := range []int{5, 1, 4, 2, 3} {
		h.Push(v)
	}

	assert.Equal(t, 5, h.Len())
	assert.Equal(t, 1, h.Peek())

	got := make([]int, 0, 5)
	for !h.IsEmpty() {
		got = append(got, h.Pop())
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, got)
}

func TestMaxHeapOrder(t *testing.T) {
	h := NewMaxHeap[int](4)
	for _, v := range []int{5, 1, 4, 2, 3} {
		h.Push(v)
	}

	assert.Equal(t, 5, h.Peek())

	got := make([]int, 0, 5)
	for !h.IsEmpty() {
		got = append(got, h.Pop())
	}
	assert.Equal(t, []int{5, 4, 3, 2, 1}, got)
}

func TestHeapTryMethodsAndClear(t *testing.T) {
	h := NewHeap[int](Less[int], 0)

	data, ok := h.TryPeek()
	assert.False(t, ok)
	assert.Equal(t, 0, data)

	data, ok = h.TryPop()
	assert.False(t, ok)
	assert.Equal(t, 0, data)
	assert.Equal(t, 0, h.Peek())
	assert.Equal(t, 0, h.Pop())

	h.Push(2)
	h.Push(1)
	data, ok = h.TryPeek()
	assert.True(t, ok)
	assert.Equal(t, 1, data)

	data, ok = h.TryPop()
	assert.True(t, ok)
	assert.Equal(t, 1, data)
	assert.Equal(t, 1, h.Len())

	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Len())
}

func TestHeapAsContainerFactories(t *testing.T) {
	min := NewMinHeapAsContainer[string](0)
	max := NewMaxHeapAsContainer[string](0)
	custom := NewHeapAsContainer(func(a, b int) bool { return a < b }, 0)

	min.Push("b")
	min.Push("a")
	assert.Equal(t, "a", min.Pop())
	assert.Equal(t, "b", min.Pop())

	max.Push("a")
	max.Push("b")
	assert.Equal(t, "b", max.Pop())
	assert.Equal(t, "a", max.Pop())

	custom.Push(3)
	custom.Push(1)
	assert.Equal(t, 1, custom.Pop())
	assert.Equal(t, 3, custom.Pop())
}
