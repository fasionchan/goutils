/*
 * Author: fasion
 * Created time: 2024-08-07 17:31:50
 * Last Modified by: fasion
 * Last Modified time: 2025-09-18 15:56:39
 */

package baseutils

import (
	"bytes"
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
