package stl

import (
	"context"
	"errors"
	"time"

	"github.com/fasionchan/goutils/basic"
)

var (
	ErrChanPipeCanceled = errors.New("ChanPipeCanceled")
)

type ChanPtr[Data any] = *Chan[Data]

type Chan[Data any] chan Data

func (c Chan[Data]) PushPro(ctx context.Context, data Data, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- data:
		return nil
	case <-time.After(timeout):
		return basic.NewTimeoutError("Chan.PushPro")
	}
}

func (c Chan[Data]) Push(data Data) Chan[Data] {
	c <- data
	return c
}

func (c Chan[Data]) Pull() Data {
	return <-c
}

func PushDataToChanX[DataChan ~chan Data, Data any](dataChan DataChan, datas ...Data) DataChan {
	for _, data := range datas {
		dataChan <- data
	}
	return dataChan
}

func NewChanFromDatasX[DataChan ~chan Data, Data any](datas ...Data) DataChan {
	return PushDataToChanX(make(DataChan, len(datas)), datas...)
}

type ChanPipePtr[Data any] = *ChanPipe[Data]

type ChanPipe[Data any] struct {
	pipe   chan Data
	cancel chan struct{}
}

func NewChanPipe[Data any]() *ChanPipe[Data] {
	return &ChanPipe[Data]{
		pipe:   make(chan Data),
		cancel: make(chan struct{}),
	}
}

func NewBufferedChanPipe[Data any](cap int) *ChanPipe[Data] {
	return &ChanPipe[Data]{
		pipe:   make(chan Data, cap),
		cancel: make(chan struct{}),
	}
}

func (pipe *ChanPipe[Data]) Cancel() {
	close(pipe.cancel)
}

func (pipe *ChanPipe[Data]) Canceled() <-chan struct{} {
	return pipe.cancel
}

func (pipe *ChanPipe[Data]) IsCanceled() bool {
	select {
	case <-pipe.cancel:
		return true
	default:
		return false
	}
}

func (pipe *ChanPipe[Data]) Reader() <-chan Data {
	return pipe.pipe
}

func (pipe *ChanPipe[Data]) Writer() chan<- Data {
	return pipe.pipe
}

func (pipe *ChanPipe[Data]) Close() {
	close(pipe.pipe)
}

func (pipe *ChanPipe[Data]) Pull() (data Data) {
	return <-pipe.pipe
}

func (pipe *ChanPipe[Data]) PullWithCtx(ctx context.Context) (data Data, err error) {
	if ctx == nil {
		return <-pipe.pipe, nil
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	case data = <-pipe.pipe:
		return
	}
}

func (pipe *ChanPipe[Data]) Push(data Data) error {
	select {
	case pipe.pipe <- data:
		return nil
	case <-pipe.cancel:
		return ErrChanPipeCanceled
	}
}

func (pipe *ChanPipe[Data]) PushWithCtx(ctx context.Context, data Data) error {
	if ctx == nil {
		return pipe.Push(data)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case pipe.pipe <- data:
		return nil
	case <-pipe.cancel:
		return ErrChanPipeCanceled
	}
}

type ChanPipes[Data any] []*ChanPipe[Data]

func NewChanPipes[Data any](pipes ...*ChanPipe[Data]) ChanPipes[Data] {
	return pipes
}

func (pipes ChanPipes[Data]) Close() {
	ForEach(pipes, ChanPipePtr[Data].Close)
}