package basic

import (
	"fmt"
	"testing"
)

func panicFunc() (err error) {
	defer PanicRecover(&err, nil, nil, nil)
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