/*
 * Author: fasion
 * Created time: 2023-01-12 15:12:56
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:32:06
 */

package baseutils

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var plaintexts = [][]byte{
	[]byte("fasionchan"),
	[]byte("小菜学编程"),
	[]byte("fasionchan.com"),
	[]byte("coding-fans"),
	[]byte("https://fasionchan.com"),
	[]byte("Python源码剖析"),
}

func TestAes(t *testing.T) {
	rand.Seed(time.Now().UnixMicro())
	cipher, err := NewAesCipher([]byte("abc"), 256)
	if err != nil {
		t.Error(err)
		return
	}

	for _, text := range plaintexts {
		ciphertext, err := cipher.Encrypt(text, true)
		if err != nil {
			t.Error(err)
			continue
		}

		plaintext, err := cipher.Decrypt(ciphertext, true)
		if err != nil {
			t.Error(err)
			continue
		}

		if bytes.Compare(plaintext, text) != 0 {
			t.Errorf(fmt.Sprintf("case failed: %s", text))
			t.Fail()
			continue
		}
	}
}
