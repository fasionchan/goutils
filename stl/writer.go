package stl

import (
	"io"
	"iter"
)

type Writer[Datas ~[]Data, Data any] interface {
	Write(datas Datas) (n int, err error)
}

type WriteCloser[Datas ~[]Data, Data any] interface {
	Writer[Datas, Data]
	io.Closer
}

func WriteSeq[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq[Data]) (int64, error) {
	var total int64
	for data := range seq {
		n, err := dst.Write(Datas{data})
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func WriteSeq2[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq2[any, Data]) (int64, error) {
	var total int64
	for _, data := range seq {
		n, err := dst.Write(Datas{data})
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func WriteDataErrorSeq[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq2[Data, error]) (int64, error) {
	var total int64
	for data, err := range seq {
		if err != nil {
			return total, err
		}

		n, err := dst.Write(Datas{data})
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
