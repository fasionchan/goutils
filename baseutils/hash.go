/*
 * Author: fasion
 * Created time: 2023-01-12 15:04:53
 * Last Modified by: fasion
 * Last Modified time: 2023-01-12 15:10:31
 */

package baseutils

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

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
