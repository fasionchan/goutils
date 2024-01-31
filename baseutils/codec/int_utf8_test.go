/*
 * Author: fasion
 * Created time: 2024-01-27 15:13:38
 * Last Modified by: fasion
 * Last Modified time: 2024-01-27 19:42:19
 */

package codec

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testUintUtf8(t *testing.T, u uint64) {
	utf8 := Uint64ToUtf8Bytes(u)
	parsed, err := ParseUint64FromUtf8Bytes(utf8)
	if err != nil {
		t.Fatal(parsed)
		return
	}

	assert.Equal(t, parsed, u)
}

func TestXxx(t *testing.T) {
	testUintUtf8(t, 0)
	for bytes := 1; bytes <= 10; bytes += 1 {
		next := uint64(1 << (bytes*8 - bytes))
		max := next - 1
		testUintUtf8(t, max)
		testUintUtf8(t, next)
		testUintUtf8(t, next+1)
	}

	for i := 0; i < 64; i++ {
		testUintUtf8(t, 1<<uint64(i))
		testUintUtf8(t, (1<<uint64(i+1))-1)
	}

	for i := 0; i < 10000; i++ {
		testUintUtf8(t, rand.Uint64())
	}
}
