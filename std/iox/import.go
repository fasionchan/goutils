package iox

import "github.com/fasionchan/goutils/stl"

type (
	Closers = stl.Closers

	Writers = stl.Writers[[]byte, byte]
	WriteClosers = stl.WriteClosers[[]byte, byte]

	UnaryWriteFunc = stl.UnaryWriteFunc[[]byte, byte]
	WriteFunc = stl.WriteFunc[[]byte, byte]
	NilNopWriteFunc = stl.NilNopWriteFunc[[]byte, byte]
)

var (
	NewReadCloser = stl.NewReadCloser[[]byte, byte]

	NewWriteCloser = stl.NewWriteCloser[[]byte, byte]
	NewWriters = stl.NewWriters[[]byte, byte]

	NewUnaryWriteFunc = stl.NewUnaryWriteFunc[[]byte, byte]
	NewWriteFunc = stl.NewWriteFunc[[]byte, byte]
	NewNilNopWriteFunc = stl.NewNilNopWriteFunc[[]byte, byte]
)