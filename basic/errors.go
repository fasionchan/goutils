package basic

import "fmt"


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