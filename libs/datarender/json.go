package datarender

import (
	"encoding/json"

	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
)

// 模板
type SmartJsonMapTemplate string

func (str SmartJsonMapTemplate) Native() string {
	return string(str)
}

func (str SmartJsonMapTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, concurrent bool) (SmartJsonMapRender, error) {
	return ParseSmartJsonMapRenderTemplate(str.Native(), funcMap, concurrent)
}

func (str SmartJsonMapTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, data any) (types.SmartJsonMap, error) {
	return ParseAndRenderSmartJsonMapTemplate(str.Native(), funcMap, data)
}

type SmartJsonMapRender = DataRender[types.SmartJsonMap]

func ParseSmartJsonMapRenderTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (SmartJsonMapRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, JsonUnmarshalString[types.SmartJsonMap], "", concurrent, nil)
}

func ParseAndRenderSmartJsonMapTemplate(text string, funcMap templatex.TemplateFuncMap, data any) (types.SmartJsonMap, error) {
	render, err := ParseSmartJsonMapRenderTemplate(text, funcMap, false)
	if err != nil {
		return nil, err
	}
	return render.Render(data)
}

type SmartJsonMapTextTemplate = DataTextTemplate[types.SmartJsonMap]

func ParseSmartJsonMapTextTemplate(text string, funcMap templatex.TemplateFuncMap) (*SmartJsonMapTextTemplate, error) {
	return ParseDataTextTemplate[types.SmartJsonMap](text, funcMap, nil)
}

type SmartJsonMapExtractorTemplate = DataExtractorTemplate[types.SmartJsonMap]

func ParseSmartJsonMapExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (*SmartJsonMapExtractorTemplate, error) {
	return ParseDataExtractorTemplate[types.SmartJsonMap](text, funcMap, concurrent, "", nil)
}

func JsonMarshalString(v any) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func MustJsonMarshalString(v any) string {
	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func JsonMarshalRaw(v any) (json.RawMessage, error) {
	return json.Marshal(v)
}

func MustJsonMarshalRaw(v any) json.RawMessage {
	json, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return json
}

func JsonUnmarshal[T any](data []byte) (result T, err error) {
	err = json.Unmarshal(data, &result)
	return
}

func MustJsonUnmarshal[T any](data []byte) T {
	result, err := JsonUnmarshal[T](data)
	if err != nil {
		panic(err)
	}
	return result
}

func JsonUnmarshalString[T any](data string) (result T, err error) {
	return JsonUnmarshal[T]([]byte(data))
}

func MustJsonUnmarshalString[T any](data string) T {
	result, err := JsonUnmarshalString[T](data)
	if err != nil {
		panic(err)
	}
	return result
}
