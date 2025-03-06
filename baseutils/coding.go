/*
 * Author: fasion
 * Created time: 2023-01-11 14:53:49
 * Last Modified by: fasion
 * Last Modified time: 2025-03-06 13:47:27
 */

package baseutils

import (
	"bytes"
	"encoding/base64"
)

func Pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddings := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, paddings...)
}

func Pkcs5Unpadding(data []byte) []byte {
	bytes := len(data)
	padding := int(data[bytes-1])
	return data[:bytes-padding]
}

func Base64Encode(enc *base64.Encoding, data []byte) []byte {
	result := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(result, data)
	return result
}

func Base64Decode(enc *base64.Encoding, data []byte) ([]byte, error) {
	result := make([]byte, enc.DecodedLen(len(data)))
	if n, err := enc.Decode(result, data); err != nil {
		return nil, err
	} else {
		return result[:n], nil
	}
}

func StdBase64Encode(data []byte) []byte {
	return Base64Encode(base64.StdEncoding, data)
}

func StdBase64Decode(data []byte) ([]byte, error) {
	return Base64Decode(base64.StdEncoding, data)
}

func UrlBase64Encode(data []byte) []byte {
	return Base64Encode(base64.URLEncoding, data)
}

func UrlBase64Decode(data []byte) ([]byte, error) {
	return Base64Decode(base64.URLEncoding, data)
}

var StdBase64EncodeToString = base64.StdEncoding.EncodeToString
var StdBase64DecodeString = base64.StdEncoding.DecodeString

var UrlBase64EncodeToString = base64.URLEncoding.EncodeToString
var UrlBase64DecodeString = base64.URLEncoding.DecodeString
