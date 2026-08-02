package types

import (
	"github.com/fasionchan/goutils/stl"
	"golang.org/x/exp/constraints"
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
