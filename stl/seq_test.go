package stl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataSeq(t *testing.T) {
	seq := DataSeq(1, 2, 3)
	datas := NewSliceFromSeq(seq)
	assert.Equal(t, datas.Len(), 3)
	assert.Equal(t, datas[0], 1)
	assert.Equal(t, datas[1], 2)
	assert.Equal(t, datas[2], 3)
}

func TestEmptySeq(t *testing.T) {
	datas := NewSliceFromSeq(EmptySeq[int])
	assert.Equal(t, datas.Len(), 0)
}

func TestEmptySeq2(t *testing.T) {
}

func TestSingularDataSeq(t *testing.T) {
	datas := NewSliceFromSeq(SingularDataSeq(1))
	assert.Equal(t, datas.Len(), 1)
	assert.Equal(t, datas[0], 1)
}
