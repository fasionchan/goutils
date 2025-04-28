/*
 * Author: fasion
 * Created time: 2024-12-28 23:24:22
 * Last Modified by: fasion
 * Last Modified time: 2025-04-28 08:52:00
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
