package types

import (
	"github.com/fasionchan/goutils/stl"
	"golang.org/x/exp/constraints"
)

type (
	Float32s = stl.Orderables[float32]
	Float64s = stl.Orderables[float64]

	Ints = stl.Orderables[int]
	Int8s = stl.Orderables[int8]
	Int16s = stl.Comparables[int16]
	Int32s = stl.Orderables[int32]
	Int64s = stl.Orderables[int64]

	Uints = stl.Orderables[uint]
	Uint8s = stl.Orderables[uint8]
	Uint16s = stl.Orderables[uint16]
	Uint32s = stl.Orderables[uint32]
	Uint64s = stl.Orderables[uint64]
)

type Number interface {
	constraints.Integer | constraints.Float
}

type ArithmeticProgression[T Number] struct {
	current T
	delta   T
}

func NewArithmeticProgression[T Number](start, delta T) *ArithmeticProgression[T] {
	return &ArithmeticProgression[T]{
		current: start - delta,
		delta:   delta,
	}
}

func (ap *ArithmeticProgression[T]) Current() T {
	if ap == nil {
		return 0
	}
	return ap.current
}

func (ap *ArithmeticProgression[T]) Next() T {
	if ap == nil {
		return 0
	}

	ap.current += ap.delta
	return ap.current
}

type SequenceNumber = ArithmeticProgression[int]

func NewSequenceNumber() *SequenceNumber {
	return NewSequenceNumberPro(1, 1)
}

func NewSequenceNumberPro(start, delta int) *SequenceNumber {
	return NewArithmeticProgression(start, delta)
}

type StringCounter = stl.Counter[string]

var NewStringCounter = stl.NewCounter[string]