package logging

import (
	"context"

	"go.uber.org/zap"
)

type LoggerFunc func(msg string, fields ...zap.Field)

func LogContextError(ctx context.Context, log LoggerFunc, msg string, fields ...zap.Field) error {
	err := ctx.Err()
	fields = Fields(fields).Dup().Append(zap.Error(err))
	log(msg, fields...)
	return err
}