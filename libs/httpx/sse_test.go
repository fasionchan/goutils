/*
 * Author: fasion
 * Created time: 2024-10-25 15:34:03
 * Last Modified by: fasion
 * Last Modified time: 2025-07-04 19:47:58
 */

package httpx

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/fasionchan/goutils/stl"
)

func TestServerSentEventReader(t *testing.T) {
	reader := NewServerSentEventReader(bytes.NewReader([]byte(`data: hello
data: hello2

data: hello3


`)), "\n")
	messages, err := reader.ReadAll()
	if err != nil {
		t.Fatal(messages)
		return
	}

	fmt.Println(messages)
	fmt.Println(messages.EventDatas(""))
}

func TestServerSentEventParser(t *testing.T) {
	printer := stl.NewPrinter[ServerSentEventMessages](os.Stdout)
	parser := NewServerSentEventParser(printer).WithEndOfLine("\n")

	parser.Write([]byte("data: hell01\n"))
	parser.Write([]byte("data: hello2\n\n"))
	parser.Write([]byte("data: hello3\n\n"))
}
