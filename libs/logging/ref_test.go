/*
 * Author: fasion
 * Created time: 2025-05-15 09:21:53
 * Last Modified by: fasion
 * Last Modified time: 2025-05-27 11:11:42
 */

package logging

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"gotest.tools/assert"
)

func TestLoggerRef(t *testing.T) {
	ctx := context.Background()

	ref, ctx := LoggerRefFromContextPro(ctx, true, true, "top", GetLogger())
	logger := ref.With(zap.String("TOP", "TOP")).
		GetLogger().
		With(zap.String("top", "top"))

	logger.Info("top")

	Upper(ctx)
	Last(ctx)
}

func Upper(ctx context.Context) {
	ref, ctx := LoggerRefFromContextPro(ctx, true, true, "upper", GetLogger())
	logger := ref.With(zap.String("upper", "upper"))
	logger.Info("upper")

	Middle(ctx)
}

func Middle(ctx context.Context) {
	ref, ctx := LoggerRefFromContextPro(ctx, true, true, "middle", GetLogger())
	defer ref.Reset()()

	ref.Info("middle")

	Lower(ctx)
}

func Lower(ctx context.Context) {
	LoggerFromContextWithFallbacksX(ctx, "lower", GetLogger()).
		With(zap.String("lower", "lower")).
		Info("lowner")
}

func Last(ctx context.Context) {
	LoggerFromContextWithFallbacksX(ctx, "last", GetLogger()).
		With(zap.String("last", "last")).
		Info("last")
}

func TestLoggerRefNamedOnce(t *testing.T) {
	ref := NewLoggerRef(GetLogger())
	ref.NamedOnce("bar")
	assert.Equal(t, ref.Name(), "bar")

	ref.NamedOnce("bar")
	assert.Equal(t, ref.Name(), "bar")

	ref.NamedOnce("foo")
	assert.Equal(t, ref.Name(), "bar.foo")

	ref.NamedOnce("foo")
	assert.Equal(t, ref.Name(), "bar.foo")

	ref.NamedOnce("bar")
	assert.Equal(t, ref.Name(), "bar.foo.bar")

	ref.NamedOnce("bar")
	assert.Equal(t, ref.Name(), "bar.foo.bar")
}
