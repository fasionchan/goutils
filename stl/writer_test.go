package stl

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFunc(t *testing.T) {
	printer := NewPrinter[[]int](os.Stdout)
	printer.Write([]int{1, 2, 3})
}

func TestConverterWriter(t *testing.T) {
	printer := NewPrinter[[]int](os.Stdout)
	converter := NewConvertWriter[[]int](func(input int) (int, error) {
		return input * input, nil
	}, printer)

	converter.Write([]int{1, 2, 3})
}

func TestBuffer(t *testing.T) {
	buffer := NewBuffer[[]int]()
	buffer.Write([]int{1})
	buffer.Write([]int{2, 3})
	assert.Equal(t, []int{1, 2, 3}, buffer.Datas())
}

func TestSkipWriter(t *testing.T) {
	buffer := NewBuffer[[]int]()
	skipWriter := NewSkipWriter(3, buffer)
	skipWriter.Write([]int{1})
	skipWriter.Write([]int{2, 3, 4})
	skipWriter.Write([]int{5, 6})
	assert.Equal(t, []int{4, 5, 6}, buffer.Datas())
}

func TestLimitWriter(t *testing.T) {
	buffer := NewBuffer[[]int]()
	limitWriter := NewLimitWriter(3, buffer)
	limitWriter.Write([]int{1})
	limitWriter.Write([]int{2, 3, 4})
	limitWriter.Write([]int{5, 6})
	assert.Equal(t, []int{1, 2, 3}, buffer.Datas())
}

func TestWriterType(t *testing.T) {
	fmt.Println(reflect.TypeOf((*Writer[[]byte, byte])(nil)).Elem())
	fmt.Println(reflect.TypeOf((*io.Writer)(nil)).Elem())
}