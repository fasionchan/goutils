package stl

import "io"

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
