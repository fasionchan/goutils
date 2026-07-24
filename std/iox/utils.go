package iox

import (
	"github.com/fasionchan/goutils/stl"
)

var (
	Close        = stl.Close
	CloseQuietly = stl.CloseQuietly
	NewNopCloseWriter = stl.NewNopCloseWriter[[]byte, byte]
)
