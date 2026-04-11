package stl

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuffer_byteReadWrite(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	n, err := b.Write([]byte("hello"))
	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, 5, b.Len())
	require.Equal(t, []byte("hello"), []byte(b.Datas()))

	p := make([]byte, 2)
	n, err = b.Read(p)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []byte("he"), p)
	require.Equal(t, 3, b.Len())

	n, err = b.Read(p[:0])
	require.NoError(t, err)
	require.Equal(t, 0, n)
}

func TestBuffer_byteEOF(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	_, err := b.Read(make([]byte, 1))
	require.ErrorIs(t, err, io.EOF)
}

func TestBuffer_ReadFrom(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	src := bytes.NewReader([]byte("abcd"))
	n, err := b.ReadFrom(src)
	require.NoError(t, err)
	require.Equal(t, int64(4), n)
	require.Equal(t, []byte("abcd"), []byte(b.Datas()))
}

func TestBuffer_ReadnFrom(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	src := bytes.NewReader([]byte("abcdefghij"))
	n, err := b.ReadnFrom(src, 4)
	require.NoError(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, []byte("abcd"), []byte(b.Datas()))
	rest, _ := io.ReadAll(src)
	require.Equal(t, []byte("efghij"), rest)
}

func TestBuffer_ReadnFrom_nNonPositive(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	n, err := b.ReadnFrom(bytes.NewReader([]byte("x")), 0)
	require.NoError(t, err)
	require.Equal(t, 0, n)
	require.Equal(t, 0, b.Len())
	n, err = b.ReadnFrom(bytes.NewReader([]byte("x")), -1)
	require.NoError(t, err)
	require.Equal(t, 0, n)
}

func TestBuffer_ReadnFrom_EOFBeforeN(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	n, err := b.ReadnFrom(bytes.NewReader([]byte("ab")), 10)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []byte("ab"), []byte(b.Datas()))
}

func TestBuffer_ReadnFrom_appendAfterWrite(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]byte, byte]()
	_, err := b.Write([]byte("pre"))
	require.NoError(t, err)
	n, err := b.ReadnFrom(bytes.NewReader([]byte("XY")), 2)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []byte("preXY"), []byte(b.Datas()))
}

func TestBuffer_TruncateResetGrow(t *testing.T) {
	t.Parallel()
	b := NewBufferFrom[[]byte, byte]([]byte{1, 2, 3, 4, 5})
	require.Equal(t, 5, b.Len())
	b.Truncate(2)
	require.Equal(t, []byte{1, 2}, []byte(b.Datas()))
	b.Reset()
	require.Equal(t, 0, b.Len())
	b.Grow(100)
	require.GreaterOrEqual(t, b.Cap(), 100)
}

func TestBuffer_intSlice(t *testing.T) {
	t.Parallel()
	b := NewBuffer[[]int, int]()
	_, err := b.Write([]int{7, 8, 9})
	require.NoError(t, err)
	p := make([]int, 2)
	n, err := b.Read(p)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []int{7, 8}, p)
	require.Equal(t, 1, b.Len())
}
