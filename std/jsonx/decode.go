package jsonx

import (
	"encoding/json"
	"io"

	"github.com/fasionchan/goutils/stl"
)

func NewPipeDecoder[
	Datas ~[]Data,
	Data any,
] (output stl.Writer[Datas, Data]) *stl.PipeParser[[]byte, Datas, byte, Data] {
	return stl.NewPipeParser[[]byte](
		func () (stl.ReadCloser[[]byte, byte], stl.WriteCloser[[]byte, byte]) {
			return io.Pipe()
		},
		output,
		func(reader stl.Reader[[]byte, byte]) (data Data, err error) {
			err = json.NewDecoder(reader).Decode(&data)
			return
		},
	)
}