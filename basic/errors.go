package basic

import (
	"fmt"
	"runtime"
)

func ErrorStringOrZero(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func MarkError(err error, mark string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", mark, err)
}

type TimeoutError struct {
	hint string
}

func NewTimeoutError(hint string) *TimeoutError {
	return &TimeoutError{hint: hint}
}

func (err *TimeoutError) Error() string {
	return fmt.Sprintf("%s timeout", err.hint)
}

func (err *TimeoutError) Timeout() bool {
	return true
}

type PanicError struct {
	exception any
	stack     string
}

func NewPanicError(exception any, stack string) *PanicError {
	return &PanicError{
		exception: exception,
		stack:     stack,
	}
}

func (err *PanicError) Error() string {
	if err == nil {
		return "panic: PanicError=nil"
	}
	return fmt.Sprintf("panic: exception=%v || stack=\n%s", err.exception, err.stack)
}

func (err *PanicError) GetException() any {
	if err == nil {
		return nil
	}
	return err.exception
}

func (err *PanicError) GetStack() string {
	if err == nil {
		return ""
	}
	return err.stack
}

func RecoverPanic(panicError *error, onPanicErrors ...func(*PanicError)) {
	if exception := recover(); exception != nil {
		err := NewPanicError(exception, GetStackTrace())
		if panicError != nil {
			*panicError = err
		}

		for _, onPanicError := range onPanicErrors {
			if onPanicError != nil {
				onPanicError(err)
			}
		}

		return
	}
}

func CatchPanic(fn func() error, onPanicErrors ...func(*PanicError)) (err error) {
	defer RecoverPanic(&err, onPanicErrors...)
	return fn()
}

type GetStackTraceConfig struct {
	Limit int
	All   bool
}

type GetStackTraceOption = func(*GetStackTraceConfig)

func GetStackTraceWithLimit(limit int) GetStackTraceOption {
	return func(config *GetStackTraceConfig) {
		config.Limit = limit
	}
}

func GetStackTraceAll(all bool) GetStackTraceOption {
	return func(config *GetStackTraceConfig) {
		config.All = all
	}
}

func GetStackTrace(options ...GetStackTraceOption) string {
	config := &GetStackTraceConfig{
		Limit: 102400,
	}

	for _, option := range options {
		option(config)
	}

	buffer := make([]byte, config.Limit)
	n := runtime.Stack(buffer, config.All)
	return string(buffer[:n])
}