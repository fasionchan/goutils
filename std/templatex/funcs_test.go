package templatex

import (
	"testing"
	"time"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/fasionchan/goutils/std/_time"
	"github.com/stretchr/testify/assert"
)

type TemplateFuncsTestCase struct {
	_testing.TestCaseName
	data     any
	template string
	expected string
	wantErr  bool
}

func (testCase TemplateFuncsTestCase) Run(t *testing.T) {
	result, err := TemplateFuncs.ParseTemplateAndRenderToString(testCase.GetName(), testCase.template, testCase.data)
	if testCase.wantErr {
		assert.Error(t, err)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, result)
}

func TestTemplateFuncs(t *testing.T) {
	hms := _time.Hour + _time.Minute + _time.Second
	hs := _time.Hour + _time.Second
	nativeHms := time.Hour + time.Minute + time.Second
	nativeHs := time.Hour + time.Second

	_testing.TypedRunNamedTestCases(t, []TemplateFuncsTestCase{
		{
			TestCaseName: "now",
			template:     `{{ eq ((now).Format "2006-01-02") ((today).Format "2006-01-02") }}`,
			expected:     "true",
		},
		{
			TestCaseName: "yesterday",
			template:     `{{ eq ((yesterday).AddDate 0 0 1) (today) }}`,
			expected:     "true",
		},
		{
			TestCaseName: "today",
			template:     `{{ eq ((today).AddDate 0 0 1) (tomorrow) }}`,
			expected:     "true",
		},
		{
			TestCaseName: "tomorrow",
			template:     `{{ eq ((tomorrow).AddDate 0 0 -1) (today) }}`,
			expected:     "true",
		},
		{
			TestCaseName: "jsonDecode",
			data: map[string]any{
				"s": `{"name":"John","age":30}`,
			},
			template: `{{ index (jsonDecode .s) "name" }} {{ index (jsonDecode .s) "age" }}`,
			expected: "John 30",
		},
		{
			TestCaseName: "jsonDecodeBytes",
			data: map[string]any{
				"s": []byte(`{"ok":true}`),
			},
			template: `{{ index (jsonDecode .s) "ok" }}`,
			expected: "true",
		},
		{
			TestCaseName: "jsonEncode",
			data: map[string]any{
				"v": map[string]any{"a": 1},
			},
			template: `{{ jsonEncode .v }}`,
			expected: `{"a":1}`,
		},
		{
			TestCaseName: "jsonDecodeThenEncode",
			data: map[string]any{
				"s": `{"a":1}`,
			},
			template: `{{ jsonEncode (jsonDecode .s) }}`,
			expected: `{"a":1}`,
		},
		{
			TestCaseName: "jsonDecodeInvalid",
			data: map[string]any{
				"s": "not json",
			},
			template: `{{ jsonDecode .s }}`,
			wantErr:  true,
		},
		{
			TestCaseName: "durationFromNative",
			data: map[string]any{
				"d": nativeHms,
			},
			template: `{{ durationLocaleString (durationFromNative .d) "zh" }}`,
			expected: "1时1分",
		},
		{
			TestCaseName: "durationLocaleString",
			data: map[string]any{
				"d": hms,
			},
			template: `{{ durationLocaleString .d "zh" }}`,
			expected: "1时1分",
		},
		{
			TestCaseName: "durationLocaleStringEn",
			data: map[string]any{
				"d": hms,
			},
			template: `{{ durationLocaleString .d "en" }}`,
			expected: "1hr1min",
		},
		{
			TestCaseName: "durationLocaleStringPro",
			data: map[string]any{
				"d":         hms,
				"purgeHead": true,
				"purgeTail": true,
			},
			template: `{{ durationLocaleStringPro .d "zh" .purgeHead .purgeTail 3 }}`,
			expected: "1时1分1秒",
		},
		{
			TestCaseName: "durationLocaleStringProKeepZero",
			data: map[string]any{
				"d":         hs,
				"purgeHead": true,
				"purgeTail": false,
			},
			template: `{{ durationLocaleStringPro .d "zh" .purgeHead .purgeTail 6 }}`,
			expected: "1时0分1秒",
		},
		{
			TestCaseName: "nativeDurationLocaleString",
			data: map[string]any{
				"d": nativeHms,
			},
			template: `{{ nativeDurationLocaleString .d "zh" }}`,
			expected: "1时1分",
		},
		{
			TestCaseName: "nativeDurationLocaleStringEn",
			data: map[string]any{
				"d": nativeHms,
			},
			template: `{{ nativeDurationLocaleString .d "en" }}`,
			expected: "1hr1min",
		},
		{
			TestCaseName: "nativeDurationLocaleStringPro",
			data: map[string]any{
				"d":         nativeHs,
				"purgeHead": true,
				"purgeTail": false,
			},
			template: `{{ nativeDurationLocaleStringPro .d "zh" .purgeHead .purgeTail 6 }}`,
			expected: "1时0分1秒",
		},
	})
}