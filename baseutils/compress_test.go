/*
 * Author: fasion
 * Created time: 2025-03-06 11:35:14
 * Last Modified by: fasion
 * Last Modified time: 2025-03-06 13:47:03
 */

package baseutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	original := "abc"
	compressed, err := Gzip([]byte(original))
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println("gzip | base64:", StdBase64EncodeToString(compressed))

	uncompressed, err := Gunzip(compressed)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println("gunzip:", string(uncompressed))
	assert.Equal(t, original, string(uncompressed))
}
