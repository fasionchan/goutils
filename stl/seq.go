package stl

import "iter"

func DataSeq[Data any](datas ...Data) iter.Seq[Data] {
	return NewSlice(datas...).DataSeq()
}

func IndexDataSeq[Data any](datas ...Data) iter.Seq2[int, Data] {
	return NewSlice(datas...).IndexDataSeq()
}

func SingularDataSeq[Data any](data Data) iter.Seq[Data] {
	return func(yield func(Data) bool) {
		yield(data)
	}
}

func SingularIndexDataSeq[Data any](index int, data Data) iter.Seq2[int, Data] {
	return func(yield func(int, Data) bool) {
		yield(index, data)
	}
}

func SingularKeyValueSeq[K any, V any](key K, value V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(key, value)
	}
}

func EmptySeq[Data any](yield func(Data) bool) {
}

func EmptySeq2[K any, V any](yield func(K, V) bool) {
}

func MultiSeq[Data any](seqs ...iter.Seq[Data]) iter.Seq[Data] {
	return NewSeqs(seqs...).AsSeq()
}

func MultiSeq2[K any, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return NewSeq2s(seqs...).AsSeq2()
}

type Seqs[Data any] []iter.Seq[Data]

func NewSeqs[Data any](seqs ...iter.Seq[Data]) Seqs[Data] {
	return seqs
}

func (seqs Seqs[Data]) Seq(yield func(Data) bool) {
	for _, seq := range seqs {
		for data := range seq {
			if !yield(data) {
				return
			}
		}
	}
}

func (seqs Seqs[Data]) AsSeq() iter.Seq[Data] {
	if len(seqs) == 0 {
		return EmptySeq[Data]
	}

	if len(seqs) == 1 {
		return seqs[0]
	}

	return seqs.Seq
}

type Seq2s[K any, V any] []iter.Seq2[K, V]

func NewSeq2s[K any, V any](seqs ...iter.Seq2[K, V]) Seq2s[K, V] {
	return seqs
}

func (seqs Seq2s[K, V]) Seq2(yield func(K, V) bool) {
	for _, seq := range seqs {
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (seqs Seq2s[K, V]) AsSeq2() iter.Seq2[K, V] {
	if len(seqs) == 0 {
		return EmptySeq2[K, V]
	}

	if len(seqs) == 1 {
		return seqs[0]
	}

	return seqs.Seq2
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

func WriteSeq2Key[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq2[Data, any]) (int64, error) {
	var total int64
	for key := range seq {
		n, err := dst.Write(Datas{key})
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func WriteSeq2Value[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq2[any, Data]) (int64, error) {
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

func WriteSeq2DataError[Datas ~[]Data, Data any](dst Writer[Datas, Data], seq iter.Seq2[Data, error]) (int64, error) {
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
