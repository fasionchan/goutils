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
