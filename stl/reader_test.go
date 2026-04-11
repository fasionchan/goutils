package stl

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func readAllByteReader(r Reader[[]byte, byte]) ([]byte, error) {
	var out []byte
	buf := make([]byte, 64)
	for {
		n, err := r.Read(buf)
		out = append(out, buf[:n]...)
		if err != nil {
			if err == io.EOF {
				return out, nil
			}
			return out, err
		}
	}
}

type sliceReader struct {
	b []byte
}

func (r *sliceReader) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

func TestMultiReader_concat(t *testing.T) {
	t.Parallel()
	a := &sliceReader{b: []byte("ab")}
	b := &sliceReader{b: []byte("cd")}
	mr := MultiReader[[]byte, byte](a, b)
	out, err := readAllByteReader(mr)
	require.NoError(t, err)
	require.Equal(t, []byte("abcd"), out)
}

func TestMultiReader_emptyReaders(t *testing.T) {
	t.Parallel()
	mr := MultiReader[[]byte, byte]()
	_, err := mr.Read(make([]byte, 1))
	require.ErrorIs(t, err, io.EOF)
}

func TestMultiReader_nestedFlatten(t *testing.T) {
	t.Parallel()
	inner := MultiReader[[]byte, byte](
		&sliceReader{b: []byte("1")},
		&sliceReader{b: []byte("2")},
	)
	mr := MultiReader[[]byte, byte](inner, &sliceReader{b: []byte("3")})
	out, err := readAllByteReader(mr)
	require.NoError(t, err)
	require.Equal(t, []byte("123"), out)
}

func TestMultiReader_errorPropagates(t *testing.T) {
	t.Parallel()
	errBoom := io.ErrUnexpectedEOF
	bad := errReader{err: errBoom}
	mr := MultiReader[[]byte, byte](&sliceReader{b: []byte("x")}, bad)
	buf := make([]byte, 10)
	n, err := mr.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	n, err = mr.Read(buf)
	require.ErrorIs(t, err, errBoom)
	require.Equal(t, 0, n)
}

type errReader struct {
	err error
}

func (e errReader) Read(p []byte) (int, error) {
	return 0, e.err
}

func TestPeekReader_firstPeek(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("xyz")})
	n, err := pr.Peek(2)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, 2, pr.Len())
	require.Equal(t, []byte("xy"), []byte(pr.Datas()))
}

func TestPeekReader_beforePeekLenAndDatas(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("a")})
	require.Equal(t, 0, pr.Len())
	require.Nil(t, pr.Datas())
}

func TestPeekReader_readWithoutPeek(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("ok")})
	out, err := readAllByteReader(pr)
	require.NoError(t, err)
	require.Equal(t, []byte("ok"), out)
}

func TestPeekReader_peekThenReadAll(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("abcdef")})
	_, err := pr.Peek(3)
	require.NoError(t, err)
	require.Equal(t, []byte("abc"), []byte(pr.Datas()))
	out, err := readAllByteReader(pr)
	require.NoError(t, err)
	require.Equal(t, []byte("abcdef"), out)
}

func TestPeekReader_twoPeeks(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("abcdef")})
	n1, err := pr.Peek(3)
	require.NoError(t, err)
	require.Equal(t, 3, n1)
	n2, err := pr.Peek(2)
	require.NoError(t, err)
	require.Equal(t, 2, n2)
	require.Equal(t, 5, pr.Len(), "两次 Peek 应把共 5 字节放入窥视缓冲")
	out, err := readAllByteReader(pr)
	require.NoError(t, err)
	require.Equal(t, []byte("abcdef"), out)
}

func TestPeekReader_peekNonPositive(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("ab")})
	n, err := pr.Peek(0)
	require.NoError(t, err)
	require.Equal(t, 0, n)
	n, err = pr.Peek(-1)
	require.NoError(t, err)
	require.Equal(t, 0, n)
	// n<=0 不向底层拉数据，Reader 仍为底层；再 Peek(1) 才挂 MultiReader
	n, err = pr.Peek(1)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	out, err := readAllByteReader(pr)
	require.NoError(t, err)
	require.Equal(t, []byte("ab"), out)
}

func TestPeekReader_peekEOFPartial(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]byte, byte](&sliceReader{b: []byte("x")})
	n, err := pr.Peek(10)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, []byte("x"), []byte(pr.Datas()))
	out, err := readAllByteReader(pr)
	require.NoError(t, err)
	require.Equal(t, []byte("x"), out)
	require.Equal(t, 0, pr.Len())
	require.Empty(t, pr.Datas())
}

func TestPeekReader_errorPropagates(t *testing.T) {
	t.Parallel()
	errBoom := io.ErrUnexpectedEOF
	pr := NewPeekReader[[]byte, byte](errReader{err: errBoom})
	_, err := pr.Peek(1)
	require.ErrorIs(t, err, errBoom)
}

func TestPeekReader_intSlice(t *testing.T) {
	t.Parallel()
	pr := NewPeekReader[[]int, int](&sliceIntReader{vals: []int{1, 2, 3, 4}})
	n, err := pr.Peek(2)
	require.NoError(t, err)
	require.Equal(t, 2, n)
	require.Equal(t, []int{1, 2}, []int(pr.Datas()))
	out := make([]int, 0, 8)
	buf := make([]int, 4)
	for {
		n, err := pr.Read(buf)
		out = append(out, buf[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
	}
	require.Equal(t, []int{1, 2, 3, 4}, out)
}

type sliceIntReader struct {
	vals []int
}

func (r *sliceIntReader) Read(p []int) (int, error) {
	if len(r.vals) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.vals)
	r.vals = r.vals[n:]
	return n, nil
}
