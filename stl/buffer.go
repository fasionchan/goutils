/*
 * Author: fasion
 * Created time: 2026-01-25 23:06:46
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 23:07:47
 */

package stl

import (
	"bufio"
)

type BoundedBuffer[Datas ~[]Data, Data any] struct {
	buffer Datas
	size   int
}

func NewBoundedBuffer[Datas ~[]Data, Data any](size int) *BoundedBuffer[Datas, Data] {
	return &BoundedBuffer[Datas, Data]{
		buffer: nil,
		size:   size,
	}
}

func (lb *BoundedBuffer[Datas, Data]) Datas() Datas {
	return lb.buffer
}

func (lb *BoundedBuffer[Datas, Data]) IsFull() bool {
	return len(lb.buffer) >= lb.size
}

func (lb *BoundedBuffer[Datas, Data]) Write(datas Datas) (n int, err error) {
	if len(datas) == 0 {
		return 0, nil
	}

	if lb.size <= 0 {
		return 0, bufio.ErrBufferFull
	}

	freeBytes := lb.size - len(lb.buffer)
	if freeBytes <= 0 {
		return 0, bufio.ErrBufferFull
	}

	if n := len(datas); n <= freeBytes {
		lb.buffer = append(lb.buffer, datas...)
		return n, nil
	}

	lb.buffer = append(lb.buffer, datas[:freeBytes]...)
	return freeBytes, nil
}

func (lb *BoundedBuffer[Datas, Data]) TotalWritten() int64 {
	return int64(len(lb.buffer))
}

func (lb *BoundedBuffer[Datas, Data]) IsTruncated() bool {
	return false
}

func (lb *BoundedBuffer[Datas, Data]) Reset() {
	lb.buffer = lb.buffer[:0]
}

type TruncatedBuffer[Datas ~[]Data, Data any] struct {
	BoundedBuffer[Datas, Data]
	totalWritten int64
}

func NewTruncatedBuffer[Datas ~[]Data, Data any](size int) *TruncatedBuffer[Datas, Data] {
	return &TruncatedBuffer[Datas, Data]{
		BoundedBuffer: *NewBoundedBuffer[Datas, Data](size),
		totalWritten:  0,
	}
}

func (tb *TruncatedBuffer[Datas, Data]) Write(datas Datas) (int, error) {
	totalWritten := 0
	var err error

	for len(datas) > 0 {
		var written int
		written, err = tb.BoundedBuffer.Write(datas)

		totalWritten += written
		datas = datas[written:]

		if err == nil {
			continue
		} else if err == bufio.ErrBufferFull {
			totalWritten += len(datas)
			break
		} else {
			break
		}
	}

	tb.totalWritten += int64(totalWritten)

	return totalWritten, err
}

func (tb *TruncatedBuffer[Datas, Data]) TotalWritten() int64 {
	return tb.totalWritten
}

func (tb *TruncatedBuffer[Datas, Data]) IsTruncated() bool {
	return tb.totalWritten > int64(tb.size)
}

func (tb *TruncatedBuffer[Datas, Data]) Reset() {
	tb.totalWritten = 0
	tb.BoundedBuffer.Reset()
}

// CircularBuffer 循环缓冲区，实现 io.Writer 接口
// 保留最后 size 字节的数据，新数据会覆盖旧数据
// 采用动态扩容策略，按写入量逐步增长，避免一次性分配大内存
type RingBuffer[Datas ~[]Data, Data any] struct {
	BoundedBuffer[Datas, Data]
	writePos     int
	totalWritten int64
}

// NewRingBuffer 创建新的循环缓冲区
// 缓冲区采用动态扩容策略，初始为空，随写入量增长
func NewRingBuffer[Datas ~[]Data, Data any](size int) *RingBuffer[Datas, Data] {
	if size <= 0 {
		size = 0
	}
	return &RingBuffer[Datas, Data]{
		BoundedBuffer: *NewBoundedBuffer[Datas, Data](size),
		writePos:      0,
		totalWritten:  0,
	}
}

func (rb *RingBuffer[Datas, Data]) Write(datas Datas) (n int, err error) {
	if len(datas) == 0 {
		return 0, nil
	}

	if rb.size <= 0 {
		return len(datas), nil
	}

	if rb.BoundedBuffer.IsFull() {
		return rb.writeByRing(datas)
	}

	written, err := rb.BoundedBuffer.Write(datas)
	rb.totalWritten += int64(written)
	rb.writePos = (rb.writePos + written) % rb.size
	if err != nil {
		return written, err
	}

	left := len(datas) - written
	if left <= 0 {
		return written, nil
	}

	written2, err := rb.writeByRing(datas[written:])
	return written + written2, err
}

func (rb *RingBuffer[Datas, Data]) writeByRing(datas Datas) (int, error) {
	n := len(datas)
	rb.totalWritten += int64(n)

	// 如果数据长度大于等于缓冲区大小，整个缓冲区都被覆盖
	if n >= rb.size {
		// 先计算数据写入后，新的写入位置
		rb.writePos = (rb.writePos + n) % rb.size

		// 将数据末尾部分直接拷贝进缓冲区
		copy(rb.buffer[rb.writePos:], datas[n-rb.size:n-rb.writePos])
		copy(rb.buffer[:rb.writePos], datas[n-rb.writePos:])

		return n, nil
	}

	// 如果数据长度小于缓冲区大小，直接将数据拷贝进缓冲区
	for len(datas) > 0 {
		written := copy(rb.buffer[rb.writePos:], datas)
		if written == 0 {
			panic("writeByCircular: written == 0")
		}

		rb.writePos = (rb.writePos + written) % rb.size
		datas = datas[written:]
	}

	return n, nil
}

// Datas 获取缓冲区中的数据（最后写入的字节）
func (rb *RingBuffer[Datas, Data]) Datas() Datas {
	if rb.size == 0 || rb.totalWritten == 0 {
		return nil
	}

	// 如果缓冲区还未分配，返回 nil
	if rb.buffer == nil {
		return nil
	}

	// 如果写入的数据少于缓冲区大小，返回实际数据
	if rb.totalWritten < int64(rb.size) {
		result := make(Datas, rb.writePos)
		copy(result, rb.buffer[:rb.writePos])
		return result
	}

	// 如果写入的数据超过缓冲区大小，返回最后 size 字节
	// 此时 buffer 的容量应该是 size（循环模式）
	// writePos 指向下一个写入位置，所以最后写入的数据是从 writePos 往前 size 字节
	result := make(Datas, rb.size)
	// 从 writePos 开始复制到 buffer 末尾
	copy(result, rb.buffer[rb.writePos:])
	// 从 buffer 开头复制到 writePos
	copy(result[rb.size-rb.writePos:], rb.buffer[:rb.writePos])

	return result
}

// TotalWritten 返回总共写入的字节数
func (rb *RingBuffer[Datas, Data]) TotalWritten() int64 {
	return rb.totalWritten
}

// IsTruncated 返回是否被截断（写入的数据超过缓冲区大小）
func (rb *RingBuffer[Datas, Data]) IsTruncated() bool {
	return rb.totalWritten > int64(rb.size)
}

// Reset 重置缓冲区
func (rb *RingBuffer[Datas, Data]) Reset() {
	rb.writePos = 0
	rb.totalWritten = 0

	// 可以选择释放缓冲区以节省内存，但为了性能通常保留
	// rb.buffer = nil
}
