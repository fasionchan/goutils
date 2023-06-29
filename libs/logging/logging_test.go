/*
 * Author: fasion
 * Created time: 2023-06-29 10:12:44
 * Last Modified by: fasion
 * Last Modified time: 2023-06-29 10:52:36
 */

package logging

import (
	"testing"

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

func TestCompile(t *testing.T) {
}
