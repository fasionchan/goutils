/*
 * Author: fasion
 * Created time: 2024-12-28 23:24:22
 * Last Modified by: fasion
 * Last Modified time: 2026-03-17 20:48:47
 */

package stl

// var ErrTimeout = errors.New("timeout")

// type Chan[Data any] chan Data

// func (c Chan[Data]) PushPro(ctx context.Context, data Data, timeout time.Duration) error {
// 	select {
// 	case <-ctx.Done():
// 		return nil
// 	case c <- data:
// 		return nil
// 	case <-time.After(timeout):
// 		return nil
// 	}
// }

// func (c Chan[Data]) Push(data Data) Chan[Data] {
// 	c <- data
// 	return c
// }

// func (c Chan[Data]) Pull() Data {
// 	return <-c
// }

func PushDataToChanX[DataChan ~chan Data, Data any](dataChan DataChan, datas ...Data) DataChan {
	for _, data := range datas {
		dataChan <- data
	}
	return dataChan
}

func NewChanFromDatasX[DataChan ~chan Data, Data any](datas ...Data) DataChan {
	return PushDataToChanX(make(DataChan, len(datas)), datas...)
}

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
