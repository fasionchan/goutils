/*
 * Author: fasion
 * Created time: 2026-01-25 23:36:36
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 00:07:26
 */

package types

import "github.com/fasionchan/goutils/stl"

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
