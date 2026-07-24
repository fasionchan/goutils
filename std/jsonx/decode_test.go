package jsonx

import (
	"os"
	"testing"

	"github.com/fasionchan/goutils/stl"
)

func TestNewPipeDecoder(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	printer := stl.NewPrinter[[]*Person](os.Stdout)

	decoder := NewPipeDecoder[[]*Person](printer)
	decoder.Write([]byte(`{"name":"John","age":30}`))
	decoder.Close()
	decoder.Join()
}