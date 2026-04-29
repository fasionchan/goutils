/*
 * Author: fasion
 * Created time: 2026-01-25 23:36:36
 * Last Modified by: fasion
 * Last Modified time: 2026-04-26 17:56:46
 */

package types

import (
	"github.com/fasionchan/goutils/stl"
)

type BytesBoundedBuffer = stl.BoundedBuffer[[]byte, byte]

func NewBytesBoundedBuffer(size int) *BytesBoundedBuffer {
	return stl.NewBoundedBuffer[[]byte](size)
}

type BytesTruncatedBuffer = stl.TruncatedBuffer[[]byte, byte]

func NewBytesTruncatedBuffer(size int) *BytesTruncatedBuffer {
	return stl.NewTruncatedBuffer[[]byte](size)
}

type BytesRingBuffer = stl.RingBuffer[[]byte, byte]

func NewBytesRingBuffer(size int) *BytesRingBuffer {
	return stl.NewRingBuffer[[]byte](size)
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

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	value := float64(b)
	i := 0
	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}

	return FormatFloat(value, 2, true) + units[i]
}
