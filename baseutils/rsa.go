/*
 * Author: fasion
 * Created time: 2023-01-11 16:58:25
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:40:49
 */

package baseutils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

type RsaSigner = func(io.Reader, *rsa.PrivateKey, crypto.Hash, []byte) ([]byte, error)
type RsaVerifier = func(*rsa.PublicKey, crypto.Hash, []byte, []byte) error
type RsaDecrypter = func(io.Reader, *rsa.PrivateKey, []byte) ([]byte, error)
type RsaEncrypter = func(io.Reader, *rsa.PublicKey, []byte) ([]byte, error)

type PublicKeyCipher struct {
	publicKey *rsa.PublicKey
}

func (c *PublicKeyCipher) VerifyPro(data, sign []byte, verifier RsaVerifier, hash crypto.Hash) error {
	return verifier(c.publicKey, hash, HashDataSilently(hash, data), sign)
}

func (c *PublicKeyCipher) EncryptPro(plaintext []byte, encrypter RsaEncrypter, random io.Reader) ([]byte, error) {
	return encrypter(random, c.publicKey, plaintext)
}

func (c *PublicKeyCipher) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, plaintext)
}

func (c *PublicKeyCipher) Verify(data, sign []byte) error {
	return c.VerifyPro(data, sign, rsa.VerifyPKCS1v15, crypto.SHA256)
}

func NewPublicKeyCipher(key *rsa.PublicKey) *PublicKeyCipher {
	return &PublicKeyCipher{
		publicKey: key,
	}
}

type PrivateKeyCipher struct {
	PublicKeyCipher
	privateKey *rsa.PrivateKey
}

func NewPrivateKeyCipher(key *rsa.PrivateKey) *PrivateKeyCipher {
	return &PrivateKeyCipher{
		PublicKeyCipher: PublicKeyCipher{
			publicKey: &key.PublicKey,
		},
		privateKey: key,
	}
}

func (c *PrivateKeyCipher) DecryptPro(ciphertext []byte, decrypter RsaDecrypter, random io.Reader) ([]byte, error) {
	return decrypter(random, c.privateKey, ciphertext)
}

func (c *PrivateKeyCipher) Decrypt(ciphertext []byte) ([]byte, error) {
	return c.DecryptPro(ciphertext, rsa.DecryptPKCS1v15, rand.Reader)
}

func (c *PrivateKeyCipher) SignPro(data []byte, signer RsaSigner, random io.Reader, hash crypto.Hash) ([]byte, error) {
	return signer(random, c.privateKey, hash, HashDataSilently(hash, data))
}

func (c *PrivateKeyCipher) Sign(data []byte) ([]byte, error) {
	return c.SignPro(data, rsa.SignPKCS1v15, rand.Reader, crypto.SHA256)
}

func GenerateRsaKey(bits int) (publicKey, privateKey []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}

	publicBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return
	}

	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	})

	privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return
}
