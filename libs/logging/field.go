package logging

import (
	"github.com/fasionchan/goutils/stl"
	"go.uber.org/zap"
)

type Fields []zap.Field

func NewFields(fields ...zap.Field) Fields {
	return Fields(fields)
}

func (fields Fields) Dup() Fields {
	return stl.DupSlice(fields)
}

func (fields Fields) Append(more ...zap.Field) Fields {
	return append(fields, more...)
}