/*
 * 参考 Go 标准库 bytes.Buffer（buf/off 模型、grow/ReadFrom 等），提供基于切片类型 Datas ~[]Data 的泛型缓冲区。
 */

package stl

import (
	"errors"
	"io"
)

const smallBufferSize = 64

// MinRead 与 bytes.Buffer 一致：ReadFrom 单次至少向 r 预留这么多容量，减少扩容次数。
const MinRead = 512

const maxInt = int(^uint(0) >> 1)

// ErrTooLarge 在无法为缓冲区分配足够内存时由 grow 触发 panic 的值（与 bytes.Buffer 语义一致）。
var ErrTooLarge = errors.New("stl.Buffer: too large")

var errNegativeRead = errors.New("stl.Buffer: reader returned negative count from Read")

// Buffer 为可变长度切片缓冲，语义对齐标准库 bytes.Buffer，元素类型为 Data，底层切片类型为 Datas（须 ~[]Data）。
// 零值为空缓冲，可直接使用。
type Buffer[Datas ~[]Data, Data any] struct {
	buf Datas // 内容为 buf[off : len(buf)] 中未读部分
	off int   // 读从 &buf[off] 起，写从 len(buf) 起
}

func (b *Buffer[Datas, Data]) empty() bool { return len(b.buf) <= b.off }

// Datas 返回长度等于 Len() 的切片，指向当前未读数据；下一次读写/Reset/Truncate 前有效。
func (b *Buffer[Datas, Data]) Datas() Datas { return b.buf[b.off:] }

// Len 返回未读元素个数。
func (b *Buffer[Datas, Data]) Len() int { return len(b.buf) - b.off }

// Cap 返回底层切片总容量。
func (b *Buffer[Datas, Data]) Cap() int { return cap(b.buf) }

// Available 返回 buf 尾部尚未使用的容量（cap-len）。
func (b *Buffer[Datas, Data]) Available() int { return cap(b.buf) - len(b.buf) }

// AvailableBuffer 返回 len 为 0、指向 buf 空闲尾部的切片，供追加后立即 Write 使用（与 bytes.Buffer.AvailableBuffer 一致）。
func (b *Buffer[Datas, Data]) AvailableBuffer() Datas { return b.buf[len(b.buf):] }

// Truncate 保留未读的前 n 个元素，仍复用已分配空间；n 非法时 panic。
func (b *Buffer[Datas, Data]) Truncate(n int) {
	if n == 0 {
		b.Reset()
		return
	}
	if n < 0 || n > b.Len() {
		panic("stl.Buffer: truncation out of range")
	}
	b.buf = b.buf[:b.off+n]
}

// Reset 清空缓冲，保留底层数组供后续写入复用。
func (b *Buffer[Datas, Data]) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

func (b *Buffer[Datas, Data]) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow 保证至少还能写入 n 个元素，返回新写入应开始的下标。
func (b *Buffer[Datas, Data]) grow(n int) int {
	m := b.Len()
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make(Datas, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		b.buf = bufferGrowSlice(b.buf[b.off:], b.off+n)
	}
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

// Grow 将容量至少再扩大 n；n 为负时 panic。
func (b *Buffer[Datas, Data]) Grow(n int) {
	if n < 0 {
		panic("stl.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

// Write 追加 datas，必要时扩容；过大时 panic(ErrTooLarge)。
func (b *Buffer[Datas, Data]) Write(datas Datas) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(datas))
	if !ok {
		m = b.grow(len(datas))
	}
	return copy(b.buf[m:], datas), nil
}

// Read 从缓冲读出数据到 p；无数据可读且 len(p)>0 时返回 io.EOF。
func (b *Buffer[Datas, Data]) Read(p []Data) (n int, err error) {
	if b.empty() {
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, []Data(b.buf[b.off:]))
	b.off += n
	return n, nil
}

// Next 返回接下来至多 n 个未读元素组成的切片，并推进读指针（与 bytes.Buffer.Next 一致）。
func (b *Buffer[Datas, Data]) Next(n int) Datas {
	m := b.Len()
	if n > m {
		n = m
	}
	data := b.buf[b.off : b.off+n]
	b.off += n
	return data
}

// ReadFrom 从 r 读到 EOF 并追加到缓冲；返回值 n 为读到的元素个数（与 io.ReaderFrom 一致用 int64）。
func (b *Buffer[Datas, Data]) ReadFrom(r Reader[Datas, Data]) (n int64, err error) {
	for {
		i := b.grow(MinRead)
		b.buf = b.buf[:i]
		tail := []Data(b.buf[i:cap(b.buf)])
		m, e := r.Read(tail)
		if m < 0 {
			panic(errNegativeRead)
		}
		b.buf = b.buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil
		}
		if e != nil {
			return n, e
		}
	}
}

// ReadnFrom 从 r 最多读取 n 个元素并追加到缓冲；用于 Peek 等场景。
func (b *Buffer[Datas, Data]) ReadnFrom(r Reader[Datas, Data], n int) (read int, err error) {
	if n <= 0 {
		return 0, nil
	}
	for read < n {
		need := n - read
		growN := MinRead
		if need > growN {
			growN = need
		}
		i := b.grow(growN)
		b.buf = b.buf[:i]
		tail := []Data(b.buf[i:cap(b.buf)])
		if len(tail) > need {
			tail = tail[:need]
		}
		m, e := r.Read(tail)
		if m < 0 {
			panic(errNegativeRead)
		}
		b.buf = b.buf[:i+m]
		read += m
		if e == io.EOF {
			return read, nil
		}
		if e != nil {
			return read, e
		}
		if m == 0 {
			return read, errors.New("stl.Buffer.ReadnFrom: Read returned 0 with nil error")
		}
	}
	return read, nil
}

// NewBuffer 创建空缓冲。
func NewBuffer[Datas ~[]Data, Data any]() *Buffer[Datas, Data] {
	return &Buffer[Datas, Data]{}
}

// NewBufferFrom 使用已有切片作为初始内容（取得 buf 所有权，调用后勿再使用传入的 buf）。
func NewBufferFrom[Datas ~[]Data, Data any](buf Datas) *Buffer[Datas, Data] {
	return &Buffer[Datas, Data]{buf: buf}
}

func bufferGrowSlice[Datas ~[]Data, Data any](b Datas, n int) Datas {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	c := len(b) + n
	if c < 2*cap(b) {
		c = 2 * cap(b)
	}
	b2 := make(Datas, c)
	i := copy(b2, b)
	return b2[:i]
}
