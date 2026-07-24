package stl

import (
	"io"
	"sync"
)

type NopCloser struct {}

func (nopCloser NopCloser) Close() error {
	return nil
}

type CloseFunc func() error

func (close CloseFunc) Close() error {
	return close()
}

type NilNopCloseFunc func() error

func (close NilNopCloseFunc) Close() error {
	if close == nil {
		return nil
	}

	return close()
}

type PipeParser[
	SrcDatas ~[]SrcData,
	ResultDatas ~[]ResultData,
	SrcData any,
	ResultData any,
] struct {
	result ResultData
	output Writer[ResultDatas, ResultData]

	Writer[SrcDatas, SrcData]
	io.Closer

	wg sync.WaitGroup
}

func NewPipeParser[
	SrcDatas ~[]SrcData,
	ResultDatas ~[]ResultData,
	SrcData any,
	ResultData any,
](
	pipe func() (ReadCloser[SrcDatas, SrcData], WriteCloser[SrcDatas, SrcData]),
	output Writer[ResultDatas, ResultData],
	readParse func (reader Reader[SrcDatas, SrcData]) (ResultData, error),
) (*PipeParser[SrcDatas, ResultDatas, SrcData, ResultData]) {
	reader, writer := pipe()

	parser := &PipeParser[SrcDatas, ResultDatas, SrcData, ResultData]{
		output: output,
		Writer: writer,
		Closer: writer,
	}

	parser.wg.Add(1)
	go func() {
		defer parser.wg.Done()
		defer reader.Close()

		var err error
		defer func() {
			parser.Writer = WriteFunc[SrcDatas, SrcData](func(data SrcDatas) (int, error) {
				if err != nil {
					return 0, err
				}

				// ignore extra data
				return len(data), nil
			})
		}()

		parser.result, err = readParse(reader)
		if err == nil {
			_, err = parser.output.Write(ResultDatas{parser.result})
		}
	}()

	return parser
}

func (parser *PipeParser[SrcDatas, ResultDatas, SrcData, ResultData]) CloseAndJoin() {
	parser.Close()
	parser.Join()
}

func (parser *PipeParser[SrcDatas, ResultDatas, SrcData, ResultData]) Join() {
	if parser == nil {
		return
	}
	parser.wg.Wait()
}


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
