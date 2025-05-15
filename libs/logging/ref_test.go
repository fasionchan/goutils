/*
 * Author: fasion
 * Created time: 2025-05-15 09:21:53
 * Last Modified by: fasion
 * Last Modified time: 2025-05-15 15:32:58
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

	ref, ctx := LoggerRefFromContextPro(ctx, true, true, GetLogger().Named("top"))
	logger := ref.With(zap.String("TOP", "TOP")).
		GetLogger().
		With(zap.String("top", "top"))

	logger.Info("top")

	Upper(ctx)
	Last(ctx)
}

func Upper(ctx context.Context) {
	ref, ctx := LoggerRefFromContextPro(ctx, true, true, GetLogger().Named("upper"))
	logger := ref.With(zap.String("upper", "upper"))
	logger.Info("upper")

	Middle(ctx)
}

func Middle(ctx context.Context) {
	ref, ctx := LoggerRefFromContextPro(ctx, true, true, GetLogger().Named("middle"))
	defer ref.Reset()()

	ref.Info("middle")

	Lower(ctx)
}

func Lower(ctx context.Context) {
	LoggerFromContextWithFallbacksX(ctx, GetLogger().Named("lower")).
		With(zap.String("lower", "lower")).
		Info("lowner")
}

func Last(ctx context.Context) {
	LoggerFromContextWithFallbacksX(ctx, GetLogger().Named("last")).
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
