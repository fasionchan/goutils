package _io

import (
	"io"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/stl"
)

func NewNopCloseWriter(writer io.Writer) io.WriteCloser {
	return stl.NewNopCloseWriter(writer)
}

type Writers []io.Writer

func NewWriters(writers ...io.Writer) Writers {
	return Writers(writers)
}

func (writers Writers) Append(others ...io.Writer) Writers {
	return append(writers, others...)
}

func (writers Writers) Concat(others ...Writers) Writers {
	return stl.ConcatSlicesTo(writers, others...)
}

func (writers Writers) Empty() bool {
	return len(writers) == 0
}

func (writers Writers) MultiWriter() io.Writer {
	if len(writers) == 0 {
		return nil
	} else if len(writers) == 1 {
		return writers[0]
	} else {
		return io.MultiWriter(writers...)
	}
}

func (writers Writers) PurgeNil() Writers {
	return stl.PurgeZero(writers)
}

func (writers Writers) Closers() baseutils.Closers {
	return stl.Map(writers, func(writer io.Writer) io.Closer {
		if closer, ok := writer.(io.Closer); ok {
			return closer
		}
		return nil
	})
}

func (writers Writers) ValidClosers() baseutils.Closers {
	return writers.Closers().PurgeNil()
}

type WriteClosers []io.WriteCloser

func NewWriteClosers(writeClosers ...io.WriteCloser) WriteClosers {
	return WriteClosers(writeClosers)
}

func (writeClosers WriteClosers) AsWriters() Writers {
	return stl.Map(writeClosers, func(writer io.WriteCloser) io.Writer {
		return writer
	})
}

func (writeClosers WriteClosers) AsClosers() baseutils.Closers {
	return stl.Map(writeClosers, func(writer io.WriteCloser) io.Closer {
		return writer
	})
}