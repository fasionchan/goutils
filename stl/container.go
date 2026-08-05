package stl

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

type Container[T any] interface {
	Len() int
	IsEmpty() bool
	Push(T)
	Pop() T
	Peek() T
	Clear()
}

type Stack[T any] []T

func NewStack[T any](capacity int) *Stack[T] {
	return AddrOf(make(Stack[T], 0, capacity))
}

func NewStackAsContainer[T any](capacity int) Container[T] {
	return NewStack[T](capacity)
}

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *Stack[T]) TryPeek() (data T, ok bool) {
	n := len(*s)
	if n == 0 {
		return
	}

	return (*s)[n-1], true
}

func (s *Stack[T]) Peek() T {
	data, _ := s.TryPeek()
	return data
}

func (s *Stack[T]) TryPop() (data T, ok bool) {
	n := len(*s)
	if n == 0 {
		return
	}

	data = (*s)[n-1]
	*s = (*s)[:n-1]
	ok = true

	return
}

func (s *Stack[T]) Pop() T {
	data, _ := s.TryPop()
	return data
}

func (s *Stack[T]) Clear() {
	*s = (*s)[:0]
}

type HeapIface[T any] struct {
	datas []T
	less  func(i, j T) bool
}

func NewHeapIface[T any](less func(a, b T) bool, capacity int) *HeapIface[T] {
	return &HeapIface[T]{
		datas: make([]T, 0, capacity),
		less:  less,
	}
}

func (h *HeapIface[T]) Len() int {
	return len(h.datas)
}

func (h HeapIface[T]) Less(i, j int) bool {
	return h.less(h.datas[i], h.datas[j])
}

func (h *HeapIface[T]) Swap(i, j int) {
	h.datas[i], h.datas[j] = h.datas[j], h.datas[i]
}

func (h *HeapIface[T]) Push(x any) {
	h.datas = append(h.datas, x.(T))
}

func (h *HeapIface[T]) Pop() any {
	n := len(h.datas)
	if n == 0 {
		return nil
	}

	data := h.datas[n-1]
	h.datas = h.datas[:n-1]
	return data
}

type Heap[T any] HeapIface[T]

func NewHeap[T any](less func(a, b T) bool, capacity int) *Heap[T] {
	return (*Heap[T])(NewHeapIface(less, capacity))
}

func NewHeapAsContainer[T any](less func(a, b T) bool, capacity int) Container[T] {
	return NewHeap[T](less, capacity)
}

func NewMinHeap[T constraints.Ordered](capacity int) Container[T] {
	return NewHeapAsContainer(Less[T], capacity)
}

func NewMinHeapAsContainer[T constraints.Ordered](capacity int) Container[T] {
	return NewMinHeap[T](capacity)
}

func NewMaxHeap[T constraints.Ordered](capacity int) Container[T] {
	return NewHeapAsContainer(Greater[T], capacity)
}

func NewMaxHeapAsContainer[T constraints.Ordered](capacity int) Container[T] {
	return NewMaxHeap[T](capacity)
}

func (h *Heap[T]) Iface() *HeapIface[T] {
	return (*HeapIface[T])(h)
}

func (h *Heap[T]) Len() int {
	return h.Iface().Len()
}

func (h *Heap[T]) IsEmpty() bool {
	return h.Iface().Len() == 0
}

func (h *Heap[T]) TryPeek() (data T, ok bool) {
	if h.IsEmpty() {
		return
	}

	data = h.datas[0]
	ok = true

	return
}

func (h *Heap[T]) Peek() T {
	data, _ := h.TryPeek()
	return data
}

func (h *Heap[T]) Push(x T) {
	heap.Push(h.Iface(), x)
}

func (h *Heap[T]) TryPop() (data T, ok bool) {
	if h.IsEmpty() {
		return
	}

	data = heap.Pop(h.Iface()).(T)
	ok = true

	return
}

func (h *Heap[T]) Pop() T {
	data, _ := h.TryPop()
	return data
}

func (h *Heap[T]) Clear() {
	h.datas = h.datas[:0]
}
