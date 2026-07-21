package stl

import (
	"io"
)

type Writer[Datas ~[]Data, Data any] interface {
	Write(datas Datas) (n int, err error)
}

type WriteCloser[Datas ~[]Data, Data any] interface {
	Writer[Datas, Data]
	io.Closer
}