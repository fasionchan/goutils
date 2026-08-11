package httpx

import (
	"context"
	"fmt"
	"testing"

	"github.com/fasionchan/goutils/libs/datarender"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
	"github.com/stretchr/testify/assert"
)

func TestGenericHttpTemplate_Do(t *testing.T) {
	for _, testCase := range []struct {
		tpl      *GenericHttpTemplate[datarender.BooleanTemplate, datarender.BooleanRender, bool]
		data     types.SmartJsonMap
		expected bool
	}{
		{
			tpl: &GenericHttpTemplate[datarender.BooleanTemplate, datarender.BooleanRender, bool]{
				RequestTemplate: NewRequestTemplate(
					"POST",
					"https://httpbin.org/post",
					"",
					"{{ .Body }}"),
				ResponseTemplate: "{{ eq .ResponseResult.Body.json.key \"value\" }}",
			},
			data:     types.JsonMap{"Body": `{"key": "value"}`},
			expected: true,
		},
	} {

		client, err := NewClient("")
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}

		rawResult, result, err := testCase.tpl.Do(context.Background(), client, templatex.TemplateFuncs, testCase.data)
		if err != nil {
			t.Errorf("case failed: %s", err)
			return
		}

		fmt.Println(rawResult)

		assert.True(t, result)
	}
}
