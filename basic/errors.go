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

func PanicRecover(panicError *error, onPanicErrors ...func(*PanicError)) error {
	if exception := recover(); exception != nil {
		var stackBuffer [1024000]byte
		n := runtime.Stack(stackBuffer[:], false)
		stackTrace := string(stackBuffer[:n])

		err := NewPanicError(exception, stackTrace)
		if panicError != nil {
			*panicError = err
		}

		for _, onPanicError := range onPanicErrors {
			if onPanicError != nil {
				onPanicError(err)
			}
		}

		return err
	}

	return nil
}