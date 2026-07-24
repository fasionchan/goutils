package stl

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/constraints"
)

type Writer[Datas ~[]Data, Data any] interface {
	Write(datas Datas) (n int, err error)
}

type WriteCloser[Datas ~[]Data, Data any] interface {
	Writer[Datas, Data]
	io.Closer
}

func NewWriteCloser[Datas ~[]Data, Data any](writer Writer[Datas, Data], closer io.Closer) WriteCloser[Datas, Data] {
	return struct {
		Writer[Datas, Data]
		io.Closer
	}{
		Writer: writer,
		Closer: closer,
	}
}

type NopWriter[Datas ~[]Data, Data any] struct {}

func NewNopWriter[Datas ~[]Data, Data any]() NopWriter[Datas, Data] {
	return NopWriter[Datas, Data]{}
}

func (nopWriter NopWriter[Datas, Data]) Write(datas Datas) (n int, err error) {
	return len(datas), nil
}

func NewNopCloseWriter[Datas ~[]Data, Data any](writer Writer[Datas, Data]) WriteCloser[Datas, Data] {
	return struct {
		Writer[Datas, Data]
		io.Closer
	}{
		Writer: writer,
		Closer: NopCloser{},
	}
}

type NopWriteCloser[Datas ~[]Data, Data any] struct {
	NopWriter[Datas, Data]
	NopCloser
}

func NewNopWriteCloser[Datas ~[]Data, Data any]() NopWriteCloser[Datas, Data] {
	return NopWriteCloser[Datas, Data]{}
}

type Writers[Datas ~[]Data, Data any] []Writer[Datas, Data]

func NewWriters[Datas ~[]Data, Data any](writers ...Writer[Datas, Data]) Writers[Datas, Data] {
	return writers
}

func MultiWriter[Datas ~[]Data, Data any](writers ...Writer[Datas, Data]) Writers[Datas, Data] {
	return writers
}

func (writers Writers[Datas, Data]) Append(others ...Writer[Datas, Data]) Writers[Datas, Data] {
	return append(writers, others...)
}

func (writers Writers[Datas, Data]) Native() []Writer[Datas, Data] {
	return writers
}

func (writers Writers[Datas, Data]) PurgeNil() Writers[Datas, Data] {
	return PurgeZero(writers)
}

func (writers Writers[Datas, Data]) Write(datas Datas) (n int, err error) {
	for _, writer := range writers {
		if _, err := writer.Write(datas); err != nil {
			return 0, err
		}
	}

	return len(datas), nil
}

func (writers Writers[Datas, Data]) Close1() error {
	var errs Errors = Map(writers, func(writer Writer[Datas, Data]) error {
		if closer, ok := writer.(io.Closer); ok {
			return closer.Close()
		}
		return nil
	})

	return errs.Simplify()
}

type UnaryWriteFunc[Datas ~[]Data, Data any] func(data Data) (err error)

func NewUnaryWriteFunc[Datas ~[]Data, Data any](write func(data Data) (err error)) UnaryWriteFunc[Datas, Data] {
	return write
}

func (write UnaryWriteFunc[Datas, Data]) Write(datas Datas) (n int, err error) {
	for i, data := range datas {
		if err := write(data); err != nil {
			return i, err
		}
	}

	return len(datas), nil
}

type WriteFunc[Datas any, Data any] func(datas Datas) (n int, err error)

func (write WriteFunc[Datas, Data]) Write(datas Datas) (n int, err error) {
	return write(datas)
}

type NilNopWriteFunc[Datas ~[]Data, Data any] func(datas Datas) (n int, err error)

func (write NilNopWriteFunc[Datas, Data]) Write(datas Datas) (n int, err error) {
	if write == nil {
		return len(datas), nil
	}

	return write(datas)
}

type CountWriter[Datas ~[]Data, Results ~[]Result, Data any, Result constraints.Integer] struct {
	count  Result
	output Writer[Results, Result]
}

func (counter *CountWriter[Datas, Results, Data, Result]) Close() error {
	if output := counter.output; output != nil {
		// defer output.Close()
		if _, err := output.Write([]Result{counter.count}); err != nil {
			return err
		}
	}

	return nil
}

func (counter *CountWriter[Datas, Results, Data, Result]) GetCount() Result {
	return counter.count
}

func (counter *CountWriter[Datas, Results, Data, Result]) Write(datas Datas) (n int, err error) {
	counter.count += Result(len(datas))
	return len(datas), nil
}

type LimitWriter[Datas ~[]Data, Data any] struct {
	limit  int
	writer Writer[Datas, Data]
}

func NewLimitWriter[Datas ~[]Data, Data any](limit int, writer Writer[Datas, Data]) *LimitWriter[Datas, Data] {
	return &LimitWriter[Datas, Data]{
		limit:  limit,
		writer: writer,
	}
}

func (writer *LimitWriter[Datas, Data]) Write(datas Datas) (int, error) {
	n := len(datas)

	limit := writer.limit
	if limit < 1 {
		return n, nil
	}

	if n > limit {
		datas = datas[:limit]
	}

	writtens, err := writer.writer.Write(datas)
	writer.limit -= writtens
	if err != nil {
		return writtens, err
	}

	return n, nil
}

type SkipWriter[Datas ~[]Data, Data any] struct {
	skip   int
	writer Writer[Datas, Data]
}

func NewSkipWriter[Datas ~[]Data, Data any](skip int, writer Writer[Datas, Data]) *SkipWriter[Datas, Data] {
	return &SkipWriter[Datas, Data]{
		skip:   skip,
		writer: writer,
	}
}

func (writer *SkipWriter[Datas, Data]) Write(datas Datas) (int, error) {
	skip := writer.skip
	n := len(datas)
	if n <= skip {
		writer.skip -= n
		return n, nil
	}

	datas = datas[skip:]
	writer.skip = 0

	writtens, err := writer.writer.Write(datas)
	return writtens + skip, err
}

type MapWriter[Inputs ~[]Input, Outputs ~[]Output, Input any, Output any] struct {
	mapper func(input Input) Output
	output Writer[Outputs, Output]
	NilNopCloseFunc
}

func NewMapWriter[Inputs ~[]Input, Outputs ~[]Output, Input any, Output any](mapper func(input Input) Output, output Writer[Outputs, Output]) *MapWriter[Inputs, Outputs, Input, Output] {
	var close NilNopCloseFunc
	closer, ok := output.(io.Closer)
	if ok {
		close = closer.Close
	}

	return &MapWriter[Inputs, Outputs, Input, Output]{
		mapper:          mapper,
		output:          output,
		NilNopCloseFunc: close,
	}
}

func (mapper *MapWriter[Inputs, Outputs, Input, Output]) Write(datas Inputs) (n int, err error) {
	results := Map(datas, mapper.mapper)
	return mapper.output.Write(results)
}

type WriteBuffer[Datas ~[]Data, Data any] struct {
	datas Datas
}

func NewWriteBuffer[Datas ~[]Data, Data any]() *WriteBuffer[Datas, Data] {
	return &WriteBuffer[Datas, Data]{}
}

func (buffer *WriteBuffer[Datas, Data]) Datas() Datas {
	if buffer == nil {
		return nil
	}

	return buffer.datas
}

func (buffer *WriteBuffer[Datas, Data]) Write(datas Datas) (n int, err error) {
	buffer.datas = append(buffer.datas, datas...)
	return len(datas), nil
}

type reduceWriteBuffer[Datas ~[]Data, Data any, Result any] struct {
	result Result
	reduce func(result Result, data Data) Result
}

func NewReduceWriteBuffer[Datas ~[]Data, Data any, Result any](reduce func(result Result, data Data) Result) *reduceWriteBuffer[Datas, Data, Result] {
	return &reduceWriteBuffer[Datas, Data, Result]{
		reduce: reduce,
	}
}

func (buffer *reduceWriteBuffer[Datas, Data, Result]) Result() (result Result) {
	if buffer != nil {
		result = buffer.result
	}
	return
}

func (buffer *reduceWriteBuffer[Datas, Data, Result]) Write(datas Datas) (n int, err error) {
	if len(datas) > 0 {
		buffer.result = Reduce(datas, buffer.reduce, buffer.result)
	}

	return len(datas), nil
}

type ReduceWriter[Datas ~[]Data, Results ~[]Result, Data any, Result any] struct {
	*reduceWriteBuffer[Datas, Data, Result]
	output Writer[Results, Result]
}

func NewReduceWriter[Datas ~[]Data, Results ~[]Result, Data any, Result any](reduce func(result Result, data Data) Result, output Writer[Results, Result]) *ReduceWriter[Datas, Results, Data, Result] {
	return &ReduceWriter[Datas, Results, Data, Result]{
		reduceWriteBuffer: NewReduceWriteBuffer[Datas](reduce),
		output:       output,
	}
}

func (writer *ReduceWriter[Datas, Results, Data, Result]) Close(datas Datas) (err error) {
	if output := writer.output; output != nil {
		if _, err := output.Write(Results{writer.Result()}); err != nil {
			return err
		}
	}

	return nil
}

func NewFirstWriteBuffer[Datas ~[]Data, Data any]() *reduceWriteBuffer[Datas, Data, Data] {
	first := true
	return NewReduceWriteBuffer[Datas](func(result Data, data Data) Data {
		if first {
			first = false
			return data
		}
		return result
	})
}

func NewLastWriteBuffer[Datas ~[]Data, Data any]() *reduceWriteBuffer[Datas, Data, Data] {
	return NewReduceWriteBuffer[Datas](func(result Data, data Data) Data {
		return data
	})
}

type ConvertWriter[Inputs ~[]Input, Outputs ~[]Output, Input any, Output any] struct {
	convert func(input Input) (Output, error)
	output  Writer[Outputs, Output]
	NilNopCloseFunc
}

func NewConvertWriter[Inputs ~[]Input, Outputs ~[]Output, Input any, Output any](convert func(input Input) (Output, error), output Writer[Outputs, Output]) *ConvertWriter[Inputs, Outputs, Input, Output] {
	var close NilNopCloseFunc
	closer, ok := output.(io.Closer)
	if ok {
		close = closer.Close
	}

	return &ConvertWriter[Inputs, Outputs, Input, Output]{
		convert:         convert,
		output:          output,
		NilNopCloseFunc: close,
	}
}

func (converter *ConvertWriter[Inputs, Outputs, Input, Output]) Write(datas Inputs) (n int, err error) {
	results, errs := MapWithError(datas, true, converter.convert)
	if err := errs.Simplify(); err != nil {
		return 0, err
	}

	if output := converter.output; output != nil {
		return output.Write(results)
	}

	return len(datas), nil
}

type ProcessWriter[Datas ~[]Data, Data any] struct {
	process func(data Data) (err error)
}

func NewProcessWriter[Datas ~[]Data, Data any](process func(data Data) (err error)) *ProcessWriter[Datas, Data] {
	return &ProcessWriter[Datas, Data]{
		process: process,
	}
}

func (processor *ProcessWriter[Datas, Data]) Write(datas Datas) (n int, err error) {
	for i, data := range datas {
		if err := processor.process(data); err != nil {
			return i, err
		}
	}

	return len(datas), nil
}

func NewPrinter[Datas ~[]Data, Data any](fout *os.File) *ProcessWriter[Datas, Data] {
	return NewProcessWriter[Datas](func(data Data) (err error) {
		_, err = fmt.Fprintln(fout, data)
		return
	})
}