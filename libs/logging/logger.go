/*
 * Author: fasion
 * Created time: 2025-05-15 09:00:37
 * Last Modified by: fasion
 * Last Modified time: 2025-05-15 09:21:37
 */

package logging

import "go.uber.org/zap"

type zapLogger = zap.Logger

var nopLogger = zap.NewNop()

func GetNopLogger() *zapLogger {
	return nopLogger
}
