package stl

import "io"

// NoOpIo 为空的 IO 操作，用于占位。
type NoOpIo[Datas ~[]Data, Data any] struct{}

func (NoOpIo[Datas, Data]) Read(p []Data) (int, error) {
	return 0, io.EOF
}

func (NoOpIo[Datas, Data]) Write(p []Data) (int, error) {
	return len(p), nil
}

func (NoOpIo[Datas, Data]) Close() error {
	return nil
}

// ✅ 关闭一个对象
// stl需要调用，因此不能放在std/iox包中
func Close(x any) error {
	if x == nil {
		return nil
	}

	if closer, ok := x.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func CloseQuietly(x any) {
	Close(x)
}
