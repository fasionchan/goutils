/*
 * Author: fasion
 * Created time: 2025-10-16 11:48:57
 * Last Modified by: fasion
 * Last Modified time: 2025-10-16 13:38:30
 */

package httpx

import (
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/std/templatex"
	"github.com/stretchr/testify/assert"
)

func TestParseAndRenderUrlTemplate(t *testing.T) {

	for _, testCase := range []struct {
		text     string
		data     any
		expected string
	}{
		{
			text: "https://example.com?key={{.Value}}",
			data: map[string]string{
				"Value": "value",
			},
			expected: "https://example.com?key=value",
		},
		{
			text: UrlTemplateTypeMagicText + "https://example.com?key={{.Value}}",
			data: map[string]string{
				"Value": "value",
			},
			expected: "https://example.com?key=value",
		},
		{
			text: UrlTemplateTypePlain + "https://example.com?key={{ .Value }}",
			data: map[string]string{
				"Value": "value",
			},
			expected: "https://example.com?key={{ .Value }}",
		},
		{
			text: UrlTemplateTypeMagicExtractor + "{{ printf \"https://example.com?key=%s\" .Value | setUrl }}",
			data: map[string]string{
				"Value": "value",
			},
			expected: "https://example.com?key=value",
		},
	} {
		url, err := ParseAndRenderUrlTemplate(testCase.text, templatex.TemplateFuncs, testCase.data)
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}
		fmt.Println(url)

		assert.Equal(t, testCase.expected, url)
	}
}
