/*
 * @Author: fasion
 * @Created time: Do not edit
 * @Last Modified by: fasion
 * @Last Modified time: Do not edit
 */
/*
 * Author: fasion
 * Created time: 2024-08-07 17:31:50
 * Last Modified by: fasion
 * Last Modified time: 2025-09-18 15:56:39
 */

package baseutils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMd5Hasher(t *testing.T) {
	hasher := GetMd5Hasher()

	for _, testCase := range []struct {
		data string
		sum  string
	}{
		{
			data: "",
			sum:  "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			data: "a",
			sum:  "0cc175b9c0f1b6a831c399e269772661",
		},
		{
			data: "abc",
			sum:  "900150983cd24fb0d6963f7d28e17f72",
		},
	} {
		sum := hasher.SumStos(testCase.data)
		assert.Equal(t, testCase.sum, sum, "data: %s", testCase.data)

		sum2 := hasher.SumRtos(bytes.NewBufferString(testCase.data))
		assert.Equal(t, testCase.sum, sum2, "data: %s", testCase.data)
	}

	for _, testCase := range []struct {
		datas []string
		sum   string
	}{
		{
			datas: nil,
			sum:   "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			datas: []string{},
			sum:   "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			datas: []string{"a"},
			sum:   "0cc175b9c0f1b6a831c399e269772661",
		},
		{
			datas: []string{"a", "b", "c"},
			sum:   "900150983cd24fb0d6963f7d28e17f72",
		},
	} {
		sum := hasher.SumSstos(testCase.datas)
		assert.Equal(t, testCase.sum, sum, "datas: %v", testCase.datas)
	}
}

func TestHashSum(t *testing.T) {
	for _, testCase := range []struct {
		hash func() hash.Hash
		datas string
		sum string
	}{
		{
			hash: md5.New,
			datas: "hello",
			sum: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			hash: md5.New,
			datas: "",
			sum: "d41d8cd98f00b204e9800998ecf8427e",
		},
	} {
		sum1 := HashSum(testCase.hash(), []byte(testCase.datas))
		assert.Equal(t, testCase.sum, fmt.Sprintf("%x", sum1), "hash: %s, datas: %v", testCase.hash, testCase.datas)

		sum2 := HashSumBtob(testCase.hash(), []byte(testCase.datas))
		assert.Equal(t, testCase.sum, fmt.Sprintf("%x", sum2), "hash: %s, datas: %v", testCase.hash, testCase.datas)

		sum3 := HashSumBtos(testCase.hash(), []byte(testCase.datas))
		assert.Equal(t, testCase.sum, sum3, "hash: %s, datas: %v", testCase.hash, testCase.datas)

		sum4 := HashSumStob(testCase.hash(), testCase.datas)
		assert.Equal(t, testCase.sum, fmt.Sprintf("%x", sum4), "hash: %s, datas: %v", testCase.hash, testCase.datas)

		sum5 := HashSumStos(testCase.hash(), testCase.datas)
		assert.Equal(t, testCase.sum, sum5, "hash: %s, datas: %v", testCase.hash, testCase.datas)

		sum := HashSumStos(testCase.hash(), testCase.datas)
		assert.Equal(t, testCase.sum, sum, "hash: %s, datas: %v", testCase.hash, testCase.datas)
	}
}