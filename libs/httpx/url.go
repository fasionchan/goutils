/*
 * Author: fasion
 * Created time: 2025-10-16 11:43:33
 * Last Modified by: fasion
 * Last Modified time: 2025-10-16 13:39:16
 */

package httpx

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/fasionchan/goutils/libs/datarender"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/stl"
)

const (
	UrlTemplateTypePlain          = "\t\r\n"
	UrlTemplateTypeMagicText      = datarender.TemplateTypeMagicText
	UrlTemplateTypeMagicExtractor = datarender.TemplateTypeMagicExtractor
)

// 文本解析函数
func parseUrlText(text string) (string, error) {
	return strings.TrimSpace(text), nil
}

// 模板
type UrlTemplate string

func (str UrlTemplate) Native() string {
	return string(str)
}

func (str UrlTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, concurrent bool) (UrlRender, error) {
	return ParseUrlRenderTemplate(str.Native(), funcMap)
}

func (str UrlTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, data any) (string, error) {
	return ParseAndRenderUrlTemplate(str.Native(), funcMap, data)
}

// 渲染器
type UrlRender = datarender.DataRender[string]

func ParseUrlRenderTemplate(text string, funcMap templatex.TemplateFuncMap) (UrlRender, error) {
	return datarender.ParseDataRenderTemplate(text, funcMap, UrlTemplateTypePlain, parseUrlText, "Url", false, nil)
}

func ParseAndRenderUrlTemplate(text string, funcMap templatex.TemplateFuncMap, data any) (string, error) {
	render, err := ParseUrlRenderTemplate(text, funcMap)
	if err != nil {
		return "", err
	}
	return render.Render(data)
}

// 文本型模板
type UrlTextTemplate = datarender.DataTextTemplate[string]

func ParseUrlTextTemplate(text string, funcMap templatex.TemplateFuncMap) (*UrlTextTemplate, error) {
	return datarender.ParseDataTextTemplate(text, funcMap, func(text string) (string, error) {
		return strings.TrimSpace(text), nil
	})
}

// 提取型模板
type UrlExtractorTemplate = datarender.DataExtractorTemplate[string]

func ParseUrlExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (*UrlExtractorTemplate, error) {
	return datarender.ParseDataExtractorTemplate[string](text, funcMap, concurrent, "Url", nil)
}

// todo
type SingularValues map[string]string

func (values SingularValues) Dup() SingularValues {
	return stl.DupMap(values)
}

func (values SingularValues) UrlValues() url.Values {
	return stl.MapMap[url.Values](values, func(key string, value string, _ SingularUrlValues) (string, []string) {
		return key, []string{value}
	})
}

func (values SingularValues) HttpHeader() http.Header {
	return stl.MapMap[http.Header](values, func(key string, value string, _ SingularValues) (string, []string) {
		return key, []string{value}
	})
}

type SingularUrlValues = SingularValues

func SingularizeUrlValues(values url.Values) SingularUrlValues {
	return stl.MapMap[SingularUrlValues](values, func(key string, value []string, _ url.Values) (string, string) {
		return key, stl.LastOneOrZero(value)
	})
}

type SingularHeader = SingularValues

func SingularizeHeader(header http.Header) SingularHeader {
	return stl.MapMap[SingularHeader](header, func(key string, value []string, _ http.Header) (string, string) {
		return key, stl.LastOneOrZero(value)
	})
}