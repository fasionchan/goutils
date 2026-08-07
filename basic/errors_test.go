package basic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func panicFunc() (err error) {
	defer RecoverPanic(&err, nil, nil, nil)
	panic("test")
}

func callPanicFunc() {
	fmt.Println("calling")
	fmt.Println("call result:", panicFunc())
	fmt.Println("called")
}

func TestPanicRecover(t *testing.T) {
	callPanicFunc()
}

func TestCatchPanic(t *testing.T) {
	err := CatchPanic(func () error {
		panic("testCatchPanic")
	})

	assert.Error(t, err)
	fmt.Println("catch result:", err)
}