package types

import (
	"github.com/fasionchan/goutils/stl"
	"golang.org/x/exp/constraints"
)

type (
	Float32s = stl.Comparables[float32]
	Float64s = stl.Comparables[float64]

	Ints = stl.Comparables[int]
	Int8s = stl.Comparables[int8]
	Int16s = stl.Comparables[int16]
	Int32s = stl.Comparables[int32]
	Int64s = stl.Comparables[int64]

	Uints = stl.Comparables[uint]
	Uint8s = stl.Comparables[uint8]
	Uint16s = stl.Comparables[uint16]
	Uint32s = stl.Comparables[uint32]
	Uint64s = stl.Comparables[uint64]
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