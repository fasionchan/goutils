/*
 * Author: fasion
 * Created time: 2025-05-15 09:00:37
 * Last Modified by: fasion
 * Last Modified time: 2025-06-13 14:24:29
 */

package logging

import (
	"github.com/fasionchan/goutils/stl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger = zap.Logger

var nopLogger = zap.NewNop()

func GetNopLogger() *zapLogger {
	return nopLogger
}

func RecoverPanicAndLog(logger *zapLogger, level zapcore.Level, msg string, fields ...zap.Field) {
	if r := recover(); r != nil {
		logger.Log(level, msg,
			make(stl.Slice[zap.Field], 0, len(fields)+1).
				Append(zap.Any("Panic", r)).
				Append(fields...).
				Native()...,
		)
	}
}
