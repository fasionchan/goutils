/*
 * Author: fasion
 * Created time: 2023-01-12 15:22:39
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:39:53
 */

package baseutils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestGenerateRsaKey(t *testing.T) {
	publicKey, privateKey, err := GenerateRsaKey(2048)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(publicKey)
	fmt.Println(privateKey)
}

func TestRsaSigning(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Error(err)
		return
	}

	pk := NewPrivateKeyCipher(key)

	for _, data := range plaintexts {
		sign, err := pk.Sign(data)
		if err != nil {
			t.Error(err)
			continue
		}

		if err := pk.Verify(data, sign); err != nil {
			t.Error(err)
			continue
		}

	}
}

func TestRsaEncrypting(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Error(err)
		return
	}

	pk := NewPrivateKeyCipher(key)

	for _, data := range plaintexts {
		ciphertext, err := pk.Encrypt(data)
		if err != nil {
			t.Error(err)
			continue
		}

		plaintext, err := pk.Decrypt(ciphertext)
		if err != nil {
			t.Error(err)
			continue
		}

		if bytes.Compare(plaintext, data) != 0 {
			t.Errorf("case failed: %s", data)
			t.Fail()
			continue
		}
	}
}
