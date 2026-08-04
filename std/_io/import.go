package _io

import "github.com/fasionchan/goutils/stl"

type (
	Closers = stl.Closers

	UnaryWriteFunc = stl.UnaryWriteFunc[[]byte, byte]
	WriteFunc = stl.WriteFunc[[]byte, byte]
	NilNopWriteFunc = stl.NilNopWriteFunc[[]byte, byte]
)

var (
	NewReadCloser = stl.NewReadCloser[[]byte, byte]

	NewUnaryWriteFunc = stl.NewUnaryWriteFunc[[]byte, byte]
	NewWriteFunc = stl.NewWriteFunc[[]byte, byte]
	NewNilNopWriteFunc = stl.NewNilNopWriteFunc[[]byte, byte]
)