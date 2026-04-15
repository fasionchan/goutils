package stl

import (
	"errors"
	"io"
)

// closer 与 io.Closer 等价，避免强制工厂返回类型带 Close。
type closer interface {
	Close() error
}

func closeReaderIfCloser[Datas ~[]Data, Data any](r Reader[Datas, Data]) {
	if r == nil {
		return
	}
	if c, ok := any(r).(closer); ok {
		_ = c.Close()
	}
}

type Reader[Datas ~[]Data, Data any] interface {
	Read(p []Data) (n int, err error)
}

type multiReader[Datas ~[]Data, Data any] struct {
	readers []Reader[Datas, Data]
}

func (mr *multiReader[Datas, Data]) Read(p []Data) (n int, err error) {
	for len(mr.readers) > 0 {
		// 展平嵌套的 multiReader（与 io.multiReader 一致，见 issue 13558）
		if len(mr.readers) == 1 {
			if r, ok := mr.readers[0].(*multiReader[Datas, Data]); ok {
				mr.readers = r.readers
				continue
			}
		}
		n, err = mr.readers[0].Read(p)
		if err == io.EOF {
			mr.readers[0] = NoOpIo[Datas, Data]{}
			mr.readers = mr.readers[1:]
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}

// MultiReader 返回按顺序串联多个 Reader 的逻辑 Reader；全部读完后再 Read 返回 io.EOF。
// 任一子 Reader 返回非 nil 且非 io.EOF 的错误时，该错误会向上返回。
func MultiReader[Datas ~[]Data, Data any](readers ...Reader[Datas, Data]) Reader[Datas, Data] {
	r := make([]Reader[Datas, Data], len(readers))
	copy(r, readers)
	return &multiReader[Datas, Data]{readers: r}
}

type PeekReader[Datas ~[]Data, Data any] struct {
	original Reader[Datas, Data]  // 原始 Reader
	buffer   *Buffer[Datas, Data] // 预读缓冲
	reader   Reader[Datas, Data]  // 当前读取器
}

// NewPeekReader 包装 r：未 Peek 时 Read 直接走 r；Peek 后 Read 先消费预读缓冲再读底层。
func NewPeekReader[Datas ~[]Data, Data any](r Reader[Datas, Data]) *PeekReader[Datas, Data] {
	return &PeekReader[Datas, Data]{
		original: r,
		buffer:   NewBuffer[Datas, Data](),
		reader:   r,
	}
}

// Peek 预读 n 个数据，读到新数据时，将原始 Reader 包装为 MultiReader(buffer, original)，后续 Read 时先消费预读缓冲再读底层。
func (r *PeekReader[Datas, Data]) Peek(n int) (int, error) {
	nr, err := r.buffer.ReadnFrom(r.original, n)
	if nr > 0 {
		if r.reader == r.original {
			r.reader = MultiReader[Datas, Data](r.buffer, r.original)
		}
	}

	return nr, err
}

// 读取数据，如果预读缓冲为空，则使用原始 Reader 读取
func (r *PeekReader[Datas, Data]) Read(p []Data) (n int, err error) {
	n, err = r.reader.Read(p)
	if r.buffer.Len() == 0 {
		r.reader = r.original
	}
	return
}

// Datas 返回当前预读缓冲中的未读数据；尚未 Peek 时为 nil。
func (r *PeekReader[Datas, Data]) Datas() Datas {
	if r.buffer == nil {
		return nil
	}
	return r.buffer.Datas()
}

// Len 返回预读缓冲未读长度；尚未 Peek 时为 0。
func (r *PeekReader[Datas, Data]) Len() int {
	if r.buffer == nil {
		return 0
	}
	return r.buffer.Len()
}

// ResumeReader 出错续读：底层 Read 返回非 nil 且非 io.EOF 时，若底层实现 closer 则先关闭，
// 再以当前已成功交付给调用方的元素个数为 offset 调用 factory 创建新 Reader，并继续本次 Read（填满 p）。
// offset 语义：从逻辑流起点算起，此前各次 Read 已累计返回给调用方的 Data 元素个数。
type ResumeReader[Datas ~[]Data, Data any] struct {
	factory      func(offset int) (Reader[Datas, Data], error)
	cur          Reader[Datas, Data]
	read         int // 已累计成功 Read 出的元素个数，作为下次 factory 的 offset
	triesPerRead int
}

// NewResumeReader 使用 factory 延迟打开底层；factory(offset) 应返回从逻辑流第 offset 个元素起的 Reader。
func NewResumeReader[Datas ~[]Data, Data any](factory func(offset int) (Reader[Datas, Data], error), triesPerRead int) *ResumeReader[Datas, Data] {
	if triesPerRead <= 0 {
		triesPerRead = 1
	}
	return &ResumeReader[Datas, Data]{factory: factory, triesPerRead: triesPerRead}
}

func (rr *ResumeReader[Datas, Data]) openCurrent() error {
	r, err := rr.factory(rr.read)
	if err != nil {
		return err
	}

	if r == nil {
		return errors.New("stl: ResumeReader factory returned nil reader")
	}

	rr.cur = r

	return nil
}

// Read 实现 Reader；底层非 EOF 错误时会关闭并换 reader 后重试，直到成功、EOF 或 factory 失败。
func (rr *ResumeReader[Datas, Data]) Read(p []Data) (total int, err error) {
	if len(p) == 0 {
		return
	}

	for tries := rr.triesPerRead; tries > 0 && len(p) > 0; tries-- {
		if rr.cur == nil {
			if err = rr.openCurrent(); err != nil {
				return
			}
		}

		var n int
		n, err = rr.cur.Read(p)
		if n > 0 {
			rr.read += n
			p = p[n:]
			total += n
		}

		if err == nil {
			return
		}

		closeReaderIfCloser[Datas, Data](rr.cur)
		rr.cur = nil

		if err == io.EOF {
			return
		}
	}

	return total, nil
}
