/*
 * Author: fasion
 * Created time: 2023-01-12 15:04:53
 * Last Modified by: fasion
 * Last Modified time: 2025-05-12 14:26:26
 */

package baseutils

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/fasionchan/goutils/stl"
)

var (
	md5Hasher    Hasher = md5.New
	sha1Hasher   Hasher = sha1.New
	sha224Hasher Hasher = sha256.New224
	sha256Hasher Hasher = sha256.New
	sha512Hasher Hasher = sha512.New
)

// todo: where to move?
func StringToBytes(s string) []byte {
	return []byte(s)
}

func GetMd5Hasher() Hasher {
	return md5Hasher
}

func GetSha1Hasher() Hasher {
	return sha1Hasher
}

func GetSha224Hasher() Hasher {
	return sha224Hasher
}

func GetSha256Hasher() Hasher {
	return sha256Hasher
}

func GetSha512Hasher() Hasher {
	return sha512Hasher
}

type Hasher func() hash.Hash

func (hasher Hasher) New() hash.Hash {
	return hasher()
}

func (hasher Hasher) NewHashPro() *HashPro {
	return NewHashPro(hasher())
}

func (hasher Hasher) Sum(data []byte) []byte {
	h := hasher()
	h.Write(data)
	return h.Sum(nil)
}

func (hasher Hasher) SumBstob(datas [][]byte) []byte {
	h := hasher()

	for _, data := range datas {
		h.Write(data)
	}

	return h.Sum(nil)
}

func (hasher Hasher) SumBstobX(datas ...[]byte) []byte {
	return hasher.SumBstob(datas)
}

func (hasher Hasher) SumBstos(datas [][]byte) string {
	return fmt.Sprintf("%x", hasher.SumBstob(datas))
}

func (hasher Hasher) SumBstosX(datas ...[]byte) string {
	return hasher.SumBstos(datas)
}

func (hasher Hasher) SumBtob(data []byte) []byte {
	return hasher.Sum(data)
}

func (hasher Hasher) SumBtos(data []byte) string {
	return fmt.Sprintf("%x", hasher.Sum(data))
}

func (hasher Hasher) SumSstob(datas []string) []byte {
	return hasher.SumBstob(stl.Map(datas, StringToBytes))
}

func (hasher Hasher) SumSstobX(datas ...string) []byte {
	return hasher.SumBstob(stl.Map(datas, StringToBytes))
}

func (hasher Hasher) SumSstos(datas []string) string {
	return hasher.SumBstos(stl.Map(datas, StringToBytes))
}

func (hasher Hasher) SumSstosX(datas ...string) string {
	return hasher.SumBstos(stl.Map(datas, StringToBytes))
}

func (hasher Hasher) SumStob(data string) []byte {
	return hasher.Sum([]byte(data))
}

func (hasher Hasher) SumStos(data string) string {
	return hasher.SumBtos([]byte(data))
}

func HashSum(hash hash.Hash, data []byte) []byte {
	return hash.Sum(data)
}

func HashSumBtob(hash hash.Hash, data []byte) []byte {
	return hash.Sum(data)
}

func HashSumBtos(hash hash.Hash, data []byte) string {
	return fmt.Sprintf("%x", hash.Sum(data))
}

func HashSumStob(hash hash.Hash, data string) []byte {
	return hash.Sum([]byte(data))
}

func HashSumStos(hash hash.Hash, data string) string {
	return HashSumBtos(hash, []byte(data))
}

type HashPro struct {
	hash.Hash
}

func NewHashPro(hash hash.Hash) *HashPro {
	return &HashPro{
		Hash: hash,
	}
}

func (hash *HashPro) Native() hash.Hash {
	if hash == nil {
		return nil
	}
	return hash.Hash
}

func (hash *HashPro) SumBtob(data []byte) []byte {
	return HashSumBtob(hash.Hash, data)
}

func (hash *HashPro) SumBtos(data []byte) string {
	return HashSumBtos(hash.Hash, data)
}

func (hash *HashPro) SumStob(data string) []byte {
	return HashSumStob(hash.Hash, data)
}

func (hash *HashPro) SumStos(data string) string {
	return HashSumBtos(hash, []byte(data))
}

var hasherCreators = map[crypto.Hash]func() hash.Hash{
	crypto.MD5:    md5.New,
	crypto.SHA1:   sha1.New,
	crypto.SHA224: sha256.New224,
	crypto.SHA256: sha256.New,
	crypto.SHA512: sha512.New,
}

func NewHash(hash crypto.Hash) (hash.Hash, error) {
	creator, ok := hasherCreators[hash]
	if !ok {
		return nil, nil
	}

	return creator(), nil
}

func HashData(hash crypto.Hash, datas ...[]byte) ([]byte, error) {
	hasher, err := NewHash(hash)
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		hasher.Write(data)
	}

	return hasher.Sum(nil), nil
}

func HashDataSilently(hash crypto.Hash, datas ...[]byte) []byte {
	hashed, _ := HashData(hash, datas...)
	return hashed
}
