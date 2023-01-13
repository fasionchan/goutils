/*
 * Author: fasion
 * Created time: 2023-01-11 10:36:39
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:13:46
 */

package baseutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
)

type AesCipher struct {
	block      cipher.Block
	blockBytes int
}

func NewAesCipher(key []byte, bits int) (*AesCipher, error) {
	// 对密钥算哈希，这样密钥就不用固定长度
	hash := sha256.Sum256(key)
	key = hash[:bits>>3]

	// 以哈希值为密钥，初始化加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AesCipher{
		block:      block,
		blockBytes: block.BlockSize(),
	}, nil
}

func (a *AesCipher) Encrypt(plaintext []byte, b64 bool) ([]byte, error) {
	// 对明文进行填充
	plaintext = Pkcs5Padding(plaintext, a.blockBytes)

	// 随机生成初始向量
	iv := make([]byte, a.blockBytes)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// 为密文分配空间（为拼接IV预留空间），并完成加密
	ciphertext := make([]byte, len(plaintext), len(plaintext)+len(iv))
	cipher.NewCBCEncrypter(a.block, iv).CryptBlocks(ciphertext, plaintext)

	// 将IV拼接在密文后面，两者合在一起作为加密结果
	result := append(ciphertext, iv...)

	// 由于结果是无序的二进制字节序列，可以做BASE64编码（可选）
	if b64 {
		result = Base64Encode(base64.StdEncoding, result)
	}

	return result, nil
}

func (a *AesCipher) Decrypt(cipherdata []byte, b64 bool) ([]byte, error) {
	// 如果做了BASE64编码，需要先解密
	if b64 {
		var err error
		cipherdata, err = Base64Decode(base64.StdEncoding, cipherdata)
		if err != nil {
			return nil, err
		}
	}

	// 计算密文长度
	textBytes := len(cipherdata) - a.blockBytes

	// 切出密文和IV
	iv := cipherdata[textBytes:]
	ciphertext := cipherdata[:textBytes]

	// 为明文分配空间并解密
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(a.block, iv).CryptBlocks(plaintext, ciphertext)

	// 取出填充字符
	return Pkcs5Unpadding(plaintext), nil
}
