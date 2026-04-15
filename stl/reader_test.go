package stl

import (
	"errors"
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

// readOneThenErr 第一次 Read 返回 1 字节与错误，用于触发 ResumeReader 续读。
type readOneThenErr struct {
	b   []byte
	err error
}

func (r *readOneThenErr) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = r.b[0]
	r.b = r.b[1:]
	if r.err != nil {
		return 1, r.err
	}
	return 1, nil
}

type closeSpyReader struct {
	*sliceReader
	closed *int
}

func (s *closeSpyReader) Close() error {
	if s.closed != nil {
		*s.closed++
	}
	return nil
}

// immediateErrCloser 首次 Read 即返回错误，用于验证 ResumeReader 会 Close。
type immediateErrCloser struct {
	err    error
	closed *int
}

func (i *immediateErrCloser) Read(p []byte) (int, error) {
	return 0, i.err
}

func (i *immediateErrCloser) Close() error {
	if i.closed != nil {
		*i.closed++
	}
	return nil
}

func TestResumeReader_recoverSameReadCall(t *testing.T) {
	t.Parallel()
	full := []byte("hello")
	errBoom := errors.New("boom")
	var opens int
	factory := func(offset int) (Reader[[]byte, byte], error) {
		opens++
		if opens == 1 {
			require.Equal(t, 0, offset)
			return &readOneThenErr{b: full, err: errBoom}, nil
		}
		require.Equal(t, 1, offset)
		return &sliceReader{b: append([]byte(nil), full[1:]...)}, nil
	}
	rr := NewResumeReader[[]byte, byte](factory, 2)
	out, err := readAllByteReader(rr)
	require.NoError(t, err)
	require.Equal(t, full, out)
	require.Equal(t, 2, opens)
}

func TestResumeReader_closeOnRecoverableError(t *testing.T) {
	t.Parallel()
	errBoom := errors.New("boom")
	var closes int
	var opens int
	factory := func(offset int) (Reader[[]byte, byte], error) {
		opens++
		if opens == 1 {
			require.Equal(t, 0, offset)
			return &immediateErrCloser{err: errBoom, closed: &closes}, nil
		}
		require.Equal(t, 0, offset)
		return &sliceReader{b: []byte("ok")}, nil
	}
	rr := NewResumeReader[[]byte, byte](factory, 2)
	out, err := readAllByteReader(rr)
	require.NoError(t, err)
	require.Equal(t, []byte("ok"), out)
	require.Equal(t, 1, closes)
	require.Equal(t, 2, opens)
}

func TestResumeReader_closeOnEOF(t *testing.T) {
	t.Parallel()
	var closes int
	rr := NewResumeReader[[]byte, byte](func(offset int) (Reader[[]byte, byte], error) {
		return &closeSpyReader{
			sliceReader: &sliceReader{b: []byte("z")},
			closed:      &closes,
		}, nil
	}, 2)
	out, err := readAllByteReader(rr)
	require.NoError(t, err)
	require.Equal(t, []byte("z"), out)
	require.Equal(t, 1, closes)
}

func TestResumeReader_factoryError(t *testing.T) {
	t.Parallel()
	errOpen := errors.New("open failed")
	rr := NewResumeReader[[]byte, byte](func(offset int) (Reader[[]byte, byte], error) {
		return nil, errOpen
	}, 2)
	_, err := rr.Read(make([]byte, 1))
	require.ErrorIs(t, err, errOpen)
}

func TestResumeReader_factoryNilReader(t *testing.T) {
	t.Parallel()
	rr := NewResumeReader[[]byte, byte](func(offset int) (Reader[[]byte, byte], error) {
		return nil, nil
	}, 2)
	_, err := rr.Read(make([]byte, 1))
	require.Error(t, err)
}
