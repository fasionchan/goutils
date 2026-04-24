/*
 * Author: fasion
 * Created time: 2026-01-25 23:36:36
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 00:07:26
 */

package types

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

type BytesBoundedBuffer = stl.BoundedBuffer[[]byte, byte]

func NewBytesBoundedBuffer(size int) *BytesBoundedBuffer {
	return stl.NewBoundedBuffer[[]byte, byte](size)
}

type BytesTruncatedBuffer = stl.TruncatedBuffer[[]byte, byte]

func NewBytesTruncatedBuffer(size int) *BytesTruncatedBuffer {
	return stl.NewTruncatedBuffer[[]byte, byte](size)
}

type BytesRingBuffer = stl.RingBuffer[[]byte, byte]

func NewBytesRingBuffer(size int) *BytesRingBuffer {
	return stl.NewRingBuffer[[]byte, byte](size)
}

type ByteSize int

func (b ByteSize) Native() int {
	return int(b)
}
func (b ByteSize) String() string {
	const unit = 1024
	if b < 1 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	value := float64(b)
	i := 0
	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}

	return fmt.Sprintf("%.2f%s", value, units[i])
}
