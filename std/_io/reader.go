package _io

import (
	"io"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/std/iox"
)

func TeeReader(r io.Reader, ws ...io.Writer) io.Reader {
	writers := NewWriters(ws).PurgeNil()
	if writers.Empty() {
		return r
	}

	return io.TeeReader(r, writers.MultiWriter())
}


func NewTeeReadCloser(readCloser io.ReadCloser, closeWriters bool, writers ...io.Writer) io.ReadCloser {
	return NewTeeReadCloserBuilder(readCloser, closeWriters, writers...).Build()
}

type TeeReadCloserBuilder struct {
	readCloser io.ReadCloser

	writers      IoWriters
	closeWriters bool

	writerClosers WriteClosers

	closers baseutils.Closers
}

func NewBlankTeeReadCloserBuilder() *TeeReadCloserBuilder {
	return &TeeReadCloserBuilder{}
}

func NewTeeReadCloserBuilder(readCloser io.ReadCloser, closeWriters bool, writers ...io.Writer) *TeeReadCloserBuilder {
	return &TeeReadCloserBuilder{
		readCloser:   readCloser,
		closeWriters: closeWriters,
		writers:      writers,
	}
}

func (builder *TeeReadCloserBuilder) Build() io.ReadCloser {
	if builder == nil {
		return nil
	}

	writers := NewIoWriters().Concat(
		builder.writers,
		builder.writerClosers.AsWriters(),
	)

	closers := baseutils.NewClosers(
		builder.readCloser,
	).Concat(
		builder.closers,
		builder.writerClosers.AsClosers(),
	)
	if builder.closeWriters {
		closers = closers.Concat(
			builder.writers.ValidClosers(),
		)
	}

	var reader io.Reader = builder.readCloser
	if writer := writers.PurgeNil().MultiWriter(); writer != nil {
		reader = io.TeeReader(builder.readCloser, writer)
	}

	return iox.NewReadCloser(reader, closers.PurgeNil().Simplify())
}

func (builder *TeeReadCloserBuilder) WithCloseWriters(closeWriters bool) *TeeReadCloserBuilder {
	if builder == nil {
		return nil
	}

	builder.closeWriters = closeWriters
	return builder
}

func (builder *TeeReadCloserBuilder) WithClosers(closers baseutils.Closers) *TeeReadCloserBuilder {
	if builder == nil {
		return nil
	}

	builder.closers = closers
	return builder
}

func (builder *TeeReadCloserBuilder) WithClosersX(closers ...io.Closer) *TeeReadCloserBuilder {
	return builder.WithClosers(closers)
}

func (builder *TeeReadCloserBuilder) WithReadCloser(readCloser io.ReadCloser) *TeeReadCloserBuilder {
	if builder == nil {
		return nil
	}

	builder.readCloser = readCloser
	return builder
}

func (builder *TeeReadCloserBuilder) WithWriters(writers IoWriters) *TeeReadCloserBuilder {
	if builder == nil {
		return nil
	}

	builder.writers = writers
	return builder

}

func (builder *TeeReadCloserBuilder) WithWritersX(writers ...io.Writer) *TeeReadCloserBuilder {
	return builder.WithWriters(writers)
}

func (builder *TeeReadCloserBuilder) WithWriterClosers(writerClosers WriteClosers) *TeeReadCloserBuilder {
	if builder == nil {
		return nil
	}

	builder.writerClosers = writerClosers
	return builder
}

func (builder *TeeReadCloserBuilder) WithWriterClosersX(writerClosers ...io.WriteCloser) *TeeReadCloserBuilder {
	return builder.WithWriterClosers(writerClosers)
}
