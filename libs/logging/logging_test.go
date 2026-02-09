/*
 * Author: fasion
 * Created time: 2023-06-29 10:12:44
 * Last Modified by: fasion
 * Last Modified time: 2025-06-13 14:22:17
 */

package logging

import (
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogging(t *testing.T) {
	GetLogger().Debug("debug should not shown")

	SetLoggerLevel(zapcore.DebugLevel)
	GetLogger().Debug("debug should shown")

	creator := NewLoggerCreator()
	logger := creator.NewLogger()
	logger.Debug("debug should not shown")

	creator.WithLevel(zapcore.DebugLevel)
	creator.NewLogger().Debug("debug should shown")
	logger.Debug("debug should not shown")
}

func TestLengthLimit(t *testing.T) {
	container := NewLoggerCreator().NewLoggerContainer()
	container.DynamicEncoder.WithEntryLengthLimit(200)

	logger := container.GetLogger()
	logger = logger.With(zap.String("test", "test"))

	logger.Error("TestLengthLimit", zap.Any(
		"As", strings.Repeat("a", 1024),
	))
}

func TestDefaultLengthLimit(t *testing.T) {
	logger := GetLogger()
	logger.Error("TestDefaultLengthLimit", zap.Any(
		"As", strings.Repeat("a", 10240),
	))
}

func TestDuplidatedWith(t *testing.T) {
	container := NewLoggerCreator().NewLoggerContainer()
	container.DynamicEncoder.WithEntryLengthLimit(100)

	logger := container.GetLogger()
	for range 100 {
		logger = logger.With(zap.String("test", "test"))
	}

	logger.Info("done")
}

func TestLogger(t *testing.T) {
	logger := GetLogger()
	logger.Info("hello")

	named := logger.Named("bar")
	logger.Info("has bar?")
	named.Info("has bar!")

	named = named.Named("foo")
	named.Info("has bar.foo!")
	fmt.Println(named.Name())
}

func TestArgsX(t *testing.T) {
	datas := []int{1, 2, 3}
	dups := func(numbers ...int) []int {
		return numbers
	}(datas...)

	datas[0] = 0

	fmt.Println(datas)
	fmt.Println(dups)
}

func TestCompile(t *testing.T) {
}
