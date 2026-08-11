/*
 * Author: fasion
 * Created time: 2025-10-16 14:24:52
 * Last Modified by: fasion
 * Last Modified time: 2026-07-05 10:59:50
 */

package httpx

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestTemplate_ParseAndRender(t *testing.T) {
	for _, testCase := range []struct {
		tpl      *RequestTemplate
		data     any
		expected *Request
	}{
		{
			tpl:  NewRequestTemplate("GET", "https://example.com?value={{ .Value }}", "Token: {{ .Token }}", "{{ .Body }}"),
			data: types.JsonMap{"Token": "token", "Value": "value", "Body": "body"},
			expected: &Request{
				Method: "GET",
				Url:    "https://example.com?value=value",
				Header: Header{
					"Token": []string{"token"},
				},
				BodyBytes: []byte("body"),
			},
		},
	} {
		request, err := testCase.tpl.ParseAndRender(templatex.TemplateFuncs, testCase.data)
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}
		fmt.Println(request)

		assert.Equal(t, testCase.expected.Method, request.Method)
		assert.Equal(t, testCase.expected.Url, request.Url)
		assert.Equal(t, string(testCase.expected.BodyBytes), string(request.BodyBytes))
	}
}
