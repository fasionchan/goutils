/*
 * Author: fasion
 * Created time: 2025-05-15 08:59:45
 * Last Modified by: fasion
 * Last Modified time: 2025-05-27 11:10:02
 */

package logging

import (
	"context"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"go.uber.org/zap"
)

const (
	ContextKeyLoggerRef = "__LoggerRef__"
)

type LoggerRef struct {
	*zapLogger
}

func NewLoggerRef(logger *zap.Logger) *LoggerRef {
	return &LoggerRef{
		zapLogger: logger,
	}
}

func LoggerFromContextWithFallbacksX(ctx context.Context, named string, fallbacks ...*zap.Logger) (logger *zap.Logger) {
	loggerRef, _ := LoggerRefFromContextPro(ctx, true, false, named, fallbacks...)
	return loggerRef.GetLoggerWithFallbacksX(GetNopLogger())
}

func LoggerRefFromContextPro(ctx context.Context, create, wrapContext bool, named string, fallbacks ...*zap.Logger) (*LoggerRef, context.Context) {
	ref := LoggerRefFromContext(ctx)
	defer func() {
		ref.NamedOnce(named)
	}()

	if ref != nil {
		return ref, ctx
	}

	if !create {
		return nil, ctx
	}

	ref = NewLoggerRef(stl.FindFirstNotZero(fallbacks))
	if wrapContext {
		if ctx == nil {
			ctx = context.Background()
		}

		ctx = ContextWithLoggerRef(ctx, ref)
	}

	return ref, ctx
}

func LoggerRefFromContext(ctx context.Context) (ref *LoggerRef) {
	if ctx == nil {
		return nil
	}

	ref, _ = ctx.Value(ContextKeyLoggerRef).(*LoggerRef)
	return
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return ContextWithLoggerRef(ctx, NewLoggerRef(logger))
}

func ContextWithLoggerRef(ctx context.Context, ref *LoggerRef) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ContextKeyLoggerRef, ref)
}

func (ref *LoggerRef) GetLogger() *zap.Logger {
	if ref == nil {
		return nil
	}
	return ref.zapLogger
}

func (ref *LoggerRef) GetLoggerWithFallbacksX(loggers ...*zap.Logger) *zap.Logger {
	logger := ref.GetLogger()
	if logger != nil {
		return logger
	}

	return stl.FindFirstNotZero(loggers)
}

func (ref *LoggerRef) WrapContext(ctx context.Context) context.Context {
	return ContextWithLoggerRef(ctx, ref)
}

func (ref *LoggerRef) Named(name string) *LoggerRef {
	if ref == nil {
		return nil
	}

	ref.zapLogger = ref.zapLogger.Named(name)
	return ref
}

func (ref *LoggerRef) NamedOnce(name string) *LoggerRef {
	if ref == nil {
		return nil
	}

	if name == "" {
		return ref
	}

	current := ref.Name()
	if strings.HasSuffix(current, name) {
		clen, nlen := len(current), len(name)
		if clen == nlen {
			return ref
		} else if current[clen-nlen-1] == '.' {
			return ref
		}
	}

	ref.zapLogger = ref.zapLogger.Named(name)

	return ref
}

func (ref *LoggerRef) Reset() func() {
	current := *ref
	return func() {
		*ref = current
	}
}

func (ref *LoggerRef) With(fields ...zap.Field) *LoggerRef {
	if ref == nil {
		return nil
	}

	ref.zapLogger = ref.zapLogger.With(fields...)
	return ref
}
