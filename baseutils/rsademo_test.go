/*
 * Author: fasion
 * Created time: 2023-01-12 15:01:49
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:08:41
 */

package baseutils

import (
	"crypto/rand"
	"crypto/rsa"
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
