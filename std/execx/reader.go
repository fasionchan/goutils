/*
 * Author: fasion
 * Created time: 2026-01-26 22:11:06
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 23:05:46
 */

package execx

import (
	"io"
	"unsafe"

	"github.com/fasionchan/goutils/stl"
)

type TruncateStrategy string

const (
	TruncateHead TruncateStrategy = "head"
	TruncateTail TruncateStrategy = "tail"
)

// ReadOptions 输出读取选项
type ReadOptions struct {
	Limit    int              // 输出长度限制，0 表示不限制
	Strategy TruncateStrategy // 截断策略，head 或 tail
	ReadRest bool             // 是否读取剩余数据，默认不读取
}

type ReadOption func(*ReadOptions)

func WithLimit(limit int) ReadOption {
	return func(opts *ReadOptions) {
		opts.Limit = limit
	}
}

func WithStrategy(strategy TruncateStrategy) ReadOption {
	return func(opts *ReadOptions) {
		opts.Strategy = strategy
	}
}

func WithReadRest(readRest bool) ReadOption {
	return func(opts *ReadOptions) {
		opts.ReadRest = readRest
	}
}

// ReadResult 输出读取结果
type ReadResult[Datas ~[]Data, Data any] struct {
	Data      Datas // 读取的数据
	Total     int   // 总读取的字节数，包括截断后的数据
	Truncated bool  // 是否被截断
	Error     error // 读取错误
}

// Apply 应用选项
func (opts *ReadOptions) Apply(options ...ReadOption) {
	for _, opt := range options {
		opt(opts)
	}
}

type Buffer[Datas ~[]Data, Data any] interface {
	Write(datas Datas) (n int, err error)
	Datas() Datas
	TotalWritten() int64
	IsTruncated() bool
	Reset()
}

func NewReadBuffer[Datas ~[]Data, Data any](opts ReadOptions) Buffer[Datas, Data] {
	if opts.Strategy == TruncateTail {
		return stl.NewRingBuffer[Datas](opts.Limit)
	}

	if opts.ReadRest {
		return stl.NewTruncatedBuffer[Datas](opts.Limit)
	}

	return stl.NewBoundedBuffer[Datas](opts.Limit)
}

type Reader[Datas ~[]Data, Data any] interface {
	Read(p Datas) (n int, err error)
}

type Writer[Datas ~[]Data, Data any] interface {
	Write(datas Datas) (n int, err error)
}

func Copy[Datas ~[]Data, Data any](writer Writer[Datas, Data], reader Reader[Datas, Data]) (n int64, err error) {
	return CopyBuffer(writer, reader, nil)
}

// CopyBuffer 从 reader 复制数据到 writer，使用 buf 作为中间缓冲区
// 返回总共复制的字节数和错误
// 如果 buf 为 nil 或长度为 0，会根据 Data 类型大小动态创建一个约 32KB 的缓冲区
func CopyBuffer[Datas ~[]Data, Data any](writer Writer[Datas, Data], reader Reader[Datas, Data], buf Datas) (n int64, err error) {
	if buf == nil || len(buf) == 0 {
		// 如果没有提供缓冲区，根据 Data 类型大小动态计算缓冲区长度
		// 目标：约 32KB (32*1024 字节)内存占用
		const targetSize = 32 * 1024

		var zero Data
		dataSize := int(unsafe.Sizeof(zero))
		if dataSize <= 0 {
			// 如果类型大小为 0 或负数（不应该发生），使用默认值
			dataSize = 1
		}

		// 计算需要多少个元素才能达到约 32KB
		bufLen := targetSize / dataSize
		if bufLen <= 0 {
			// 如果单个元素就超过 32KB，至少分配 1 个元素
			bufLen = 1
		}

		buf = make(Datas, bufLen)
	}

	for {
		// 从 reader 读取数据到缓冲区
		nr, er := reader.Read(buf)
		if nr > 0 {
			// 将读取的数据写入 writer
			nw, ew := writer.Write(buf[:nr])
			if nw > 0 {
				n += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return n, err
}

// Read 读取输出流，支持限制和截断策略
func Read[Datas ~[]Data, Data any](reader Reader[Datas, Data], opts ...ReadOption) ReadResult[Datas, Data] {
	readOpts := ReadOptions{
		Limit:    0,
		Strategy: TruncateHead,
		ReadRest: false,
	}
	readOpts.Apply(opts...)

	buffer := NewReadBuffer[Datas, Data](readOpts)
	n, err := Copy(buffer, reader)
	return ReadResult[Datas, Data]{
		Data:      buffer.Datas(),
		Total:     int(n),
		Truncated: buffer.IsTruncated(),
		Error:     err,
	}
}

type BytesReadResult = ReadResult[[]byte, byte]
