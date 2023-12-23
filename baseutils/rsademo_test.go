/*
 * Author: fasion
 * Created time: 2023-01-12 15:01:49
 * Last Modified by: fasion
 * Last Modified time: 2023-01-14 17:19:44
 */

package baseutils

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestGenerateRsaPrivateKey(t *testing.T) {
	// 随机生成一对密钥
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 序列化公钥
	publicBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		panic(err)
	}

	// 编码公钥
	publicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	})
	fmt.Println(publicKey)

	// 编码私钥
	privateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key), // 序列化私钥
	})
	fmt.Println(privateKey)
}

func TestRsaEncryption(t *testing.T) {
	// 例子直接生成一对新的密钥进行实验
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 待加密明文
	plaintext := []byte("fasionchan.com")

	// 用公钥加密，得到密文
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &key.PublicKey, plaintext)
	if err != nil {
		panic(err)
	}

	// 用私钥解密，得到原来的明文
	text, err := rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
	if err != nil {
		panic(err)
	}

	if bytes.Compare(text, plaintext) != 0 {
		panic("error")
	}
}

func TestRsaSignature(t *testing.T) {
	// 例子直接生成一对新的密钥进行实验
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 待签名数据
	data := []byte("fasionchan.com")

	// 对数据求哈希
	dataHash := sha256.Sum256(data)

	// 使用私钥对哈希值进行加密，结果作为数据签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, dataHash[:])
	if err != nil {
		panic(err)
	}

	// 使用公钥对签名进行解密，与哈希值对比，验证签名是否真实
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA256, dataHash[:], signature); err != nil {
		panic(err)
	}
}
