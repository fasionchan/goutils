package iox

import "github.com/fasionchan/goutils/stl"

type (
	Closers = stl.Closers

	Writers = stl.Writers[[]byte, byte]
	WriteClosers = stl.WriteClosers[[]byte, byte]
)