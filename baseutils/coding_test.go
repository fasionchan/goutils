/*
 * Author: fasion
 * Created time: 2025-03-06 15:26:28
 * Last Modified by: fasion
 * Last Modified time: 2025-03-06 15:33:00
 */

package baseutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64(t *testing.T) {
	for i := byte(0); i < 64; i++ {
		encoded := StdBase64EncodeToString([]byte{i << 2})
		decoded, err := StdBase64DecodeString(encoded)
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Equal(t, []byte{i << 2}, decoded)
		fmt.Println(encoded)
	}

	for i := byte(0); i < 64; i++ {
		encoded := UrlBase64EncodeToString([]byte{i << 2})
		decoded, err := UrlBase64DecodeString(encoded)
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Equal(t, []byte{i << 2}, decoded)
		fmt.Println(encoded)
	}
}
