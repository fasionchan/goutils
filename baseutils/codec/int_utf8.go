/*
 * Author: fasion
 * Created time: 2024-01-27 15:13:01
 * Last Modified by: fasion
 * Last Modified time: 2024-01-27 19:38:02
 */

package codec

import (
	"encoding/binary"
	"errors"

	"golang.org/x/exp/constraints"
)

func UnsignedToUtf8Bytes[Unsigned constraints.Unsigned](u Unsigned) []byte {
	return Uint64ToUtf8Bytes(uint64(u))
}

func Uint64ToUtf8Bytes(u uint64) []byte {
	var buffer [10]byte
	binary.BigEndian.PutUint64(buffer[2:], u)

	for bytes := 1; ; bytes++ {
		prefixBits := bytes
		prefix := (1 << prefixBits) - 2

		ebits := bytes*8 - prefixBits
		if u>>ebits > 0 {
			continue
		}

		result := buffer[10-bytes:]
		if bytes <= 8 {
			// 0xxxxxxx
			// 10xxxxxx XX
			// 110xxxxx XX XX
			// 1110xxxx XX XX XX
			// 11110xxx XX XX XX XX
			// 111110xx XX XX XX XX XX
			// 1111110x XX XX XX XX XX XX
			// 11111110 XX XX XX XX XX XX XX
			result[0] |= byte(prefix << (8 - prefixBits))
		} else if bytes == 9 {
			// 11111111 0xxxxxxx XX XX XX XX XX XX XX
			result[0] = 0xff
		} else {
			// 11111111 10xxxxxx XX XX XX XX XX XX XX XX
			result[0] = 0xff
			result[1] |= 0x80
		}

		return result
	}
}

var BadUtf8BytesError = errors.New("bad utf8 bytes")

func ParseUint64FromUtf8Bytes(raw []byte) (uint64, error) {
	nBytes, err := CountUtf8Bytes(raw)
	if err != nil {
		return 0, err
	}

	var value uint64
	if nBytes < 8 {
		var buffer [8]byte
		start := 8 - nBytes
		for i := 0; i < nBytes; i++ {
			buffer[start+i] = raw[i]
		}
		value = binary.BigEndian.Uint64(buffer[:])
	} else {
		value = binary.BigEndian.Uint64(raw[nBytes-8:])
	}

	eBits := (nBytes << 3) - nBytes
	if eBits < 64 {
		pBit := uint64(64 - eBits)
		value = value << pBit >> pBit
	}

	return value, nil
}

func CountUtf8Bytes(raw []byte) (int, error) {
	for maskBytes := 0; ; maskBytes++ { // the preciding 11111111 bytes
		baseBytes := (maskBytes << 3) + 1
		if len(raw) < baseBytes {
			return 0, BadUtf8BytesError
		}

		prefix := raw[maskBytes]
		for i, value := range []byte{
			0x80, // 10000000
			0xc0, // 11000000
			0xe0, // 11100000
			0xf0, // 11110000
			0xf8, // 11111000
			0xfc, // 11111100
			0xfe, // 11111110
		} {
			if prefix < value {
				return EnsureUtf8BytesCount(raw, baseBytes+i)
			}
		}

		if prefix == 0xfe {
			return EnsureUtf8BytesCount(raw, baseBytes+7)
		}
	}
}

func EnsureUtf8BytesCount(raw []byte, expected int) (int, error) {
	if len(raw) < expected {
		return 0, BadUtf8BytesError
	} else {
		return expected, nil
	}
}
