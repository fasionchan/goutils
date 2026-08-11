/*
 * Author: fasion
 * Created time: 2025-10-14 10:12:36
 * Last Modified by: fasion
 * Last Modified time: 2025-10-16 13:38:56
 */

package httpx

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fasionchan/goutils/std/templatex"
	"github.com/stretchr/testify/assert"
)

func TestParseAndRenderHeader(t *testing.T) {
	for _, testCase := range []struct {
		text   string
		data   any
		header http.Header
	}{
		{
			text: "Content-Type: application/json\nKey: {{ .Value }}",
			data: map[string]string{
				"Value": "value",
			},
			header: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"value"},
			},
		},

		{
			text: HeaderTemplateTypeMagicText + "Content-Type: application/json\nKey: {{ .Value }}",
			data: map[string]string{
				"Value": "value",
			},
			header: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"value"},
			},
		},

		{
			text: HeaderTemplateTypePlain + "Content-Type: application/json\nKey: {{ .Value }}",
			data: map[string]string{
				"Value": "value",
			},
			header: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"{{ .Value }}"},
			},
		},

		{
			text: HeaderTemplateTypeMagicExtractor + `{{ $header := httpHeaderNew }}{{ httpHeaderSet $header "Content-Type" "application/json" }}{{ httpHeaderSet $header "Key" .Value }}{{ setHeader $header }}`,
			data: map[string]string{
				"Value": "value",
			},
			header: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"value"},
			},
		},

		{
			text: HeaderTemplateTypeMagicExtractor + `{{ httpHeaderNew | call (httpHeaderSetter "Content-Type" "application/json") | call (httpHeaderSetter "Key" .Value) | setHeader }}`,
			data: map[string]string{
				"Value": "value",
			},
			header: http.Header{
				"Content-Type": []string{"application/json"},
				"Key":          []string{"value"},
			},
		},
	} {
		header, err := ParseAndRenderHeaderTemplate(testCase.text, templatex.TemplateFuncs, testCase.data)
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}
		fmt.Println(header)

		// todo
		// assert.Equal(t, header, testCase.header)
	}
}

func TestHeader_ContentDisposition(t *testing.T) {
	for _, testCase := range []struct {
		header http.Header
		filename string
		params map[string]string
	}{
		{
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename=%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1%E6%88%AA%E5%9B%BE_46dc91c9-3817-4a5f-8d25-6411018cd39f.png`},
			},
			filename: "企业微信截图_46dc91c9-3817-4a5f-8d25-6411018cd39f.png",
		},
		{
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename*=UTF-8''test.txt`},
			},
			filename: "test.txt",
		},
		{
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename*=UTF-8'zh-CN'test.txt`},
			},
			filename: "test.txt",
		},
		{
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename*=UTF-8''%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1.txt`},
			},
			filename: "企业微信.txt",
		},
		{
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename*=UTF-8'zh-CN'%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1.txt`},
			},
			filename: "企业微信.txt",
		},
	} {
		filename, err := Header(testCase.header).ContentDispositionFilename()
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}

		fmt.Println(filename)

		assert.Equal(t, testCase.filename, filename)
	}
}