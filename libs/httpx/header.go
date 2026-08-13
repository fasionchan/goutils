/*
 * Author: fasion
 * Created time: 2025-10-14 09:54:44
 * Last Modified by: fasion
 * Last Modified time: 2025-10-16 13:22:18
 */

package httpx

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/emersion/go-message"
	"github.com/fasionchan/goutils/libs/datarender"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

const (
	HeaderTemplateTypePlain          = ": HeaderTemplateTypePlain" // 纯文本
	HeaderTemplateTypeMagicText      = datarender.TemplateTypeMagicText
	HeaderTemplateTypeMagicExtractor = datarender.TemplateTypeMagicExtractor
)

type Header http.Header

func NewHeader(header http.Header) Header {
	return Header(header)
}

func (header Header) Native() http.Header {
	return http.Header(header)
}

func (header Header) MessageHeader() message.Header {
	return message.HeaderFromMap(header)
}

func (header Header) MessageHeaderPtr() *message.Header {
	_header := header.MessageHeader()
	return &_header
}

func (header Header) ContentDisposition() (string, map[string]string, error) {
	return header.MessageHeaderPtr().ContentDisposition()
}

func (header Header) ContentDispositionFilename() (string, error) {
	_, params, err := header.ContentDisposition()
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(params["filename"])
}

func (header Header) ParseContentRange() (unit string, start, end int64, total int64, err error) {
	name := "Content-Range"
	return ParseContentRange(header.Native().Get(name))
}

func (header Header) WithRange(start, end int64) Header {
	if header == nil {
		return nil
	}

	header.Native().Set("Range", FormatRange("bytes", start, end))
	return header
}

func MergeHeader(header, other http.Header, fn func (header http.Header, name string, value string)) http.Header {
	for name, values := range other {
		for _, value := range values {
			fn(header, name, value)
		}
	}
	return header
}

func MergeHeaderByAdd(header, other http.Header) http.Header {
	return MergeHeader(header, other, http.Header.Add)
}

func MergeHeaderBySet(header, other http.Header) http.Header {
	return MergeHeader(header, other, http.Header.Set)
}

func MergeHeaders(header http.Header, fn func (header http.Header, name string, value string), others ...http.Header) http.Header {
	return ReduceUnary(others, MergeHeader, header, fn)
}

func MergeHeadersByAdd(header http.Header, others ...http.Header) http.Header {
	return MergeHeaders(header, http.Header.Add, others...)
}

func MergeHeadersBySet(header http.Header, others ...http.Header) http.Header {
	return MergeHeaders(header, http.Header.Set, others...)
}

func ReduceUnary[Result any, Arg any, Datas ~[]Data, Data any](datas Datas, reducer func(Result, Data, Arg) Result, initial Result, arg Arg) (result Result) {
	result = initial
	for _, data := range datas {
		result = reducer(result, data, arg)
	}
	return
}

// 文本型解析函数
func ParseHeaderText(text string) (http.Header, error) {
	return HeaderText(text).Parse(), nil
}

// 模板
type HeaderTemplate string

func (str HeaderTemplate) Native() string {
	return string(str)
}

func (str HeaderTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, concurrent bool) (HeaderRender, error) {
	return ParseHeaderRenderTemplate(str.Native(), funcMap, concurrent)
}

func (str HeaderTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, data any) (http.Header, error) {
	if str == "" {
		return nil, nil
	}

	return ParseAndRenderHeaderTemplate(str.Native(), funcMap, data)
}

// 渲染器
type HeaderRender = datarender.DataRender[http.Header]

func ParseHeaderRenderTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (HeaderRender, error) {
	return datarender.ParseDataRenderTemplate(text, funcMap, HeaderTemplateTypePlain, ParseHeaderText, "Header", concurrent, nil)
}

func ParseAndRenderHeaderTemplate(text string, funcMap templatex.TemplateFuncMap, data any) (http.Header, error) {
	render, err := ParseHeaderRenderTemplate(text, funcMap, false)
	if err != nil {
		return nil, err
	}
	return render.Render(data)
}

// 文本型模板
type HeaderTextTemplate = datarender.DataTextTemplate[http.Header]

func ParseHeaderTextTemplate(text string, funcMap templatex.TemplateFuncMap) (*HeaderTextTemplate, error) {
	return datarender.ParseDataTextTemplate[http.Header](text, funcMap, nil)
}

// 提取型模板
type HeaderExtractorTemplate = datarender.DataExtractorTemplate[http.Header]

func ParseHeaderExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (*HeaderExtractorTemplate, error) {
	return datarender.ParseDataExtractorTemplate[http.Header](text, funcMap, concurrent, "Header", nil)
}

type HeaderText string

func (t HeaderText) Parse() http.Header {
	header := make(http.Header)

	for _, line := range strings.Split(string(t), "\n") {
		fields := strings.SplitN(line, ":", 2)
		if len(fields) != 2 {
			continue
		}

		name := strings.TrimSpace(fields[0])
		if name == "" {
			continue
		}

		header[name] = append(header[name], strings.TrimSpace(fields[1]))
	}

	return header
}

type UserAgent string

func (ua UserAgent) Native() string {
	return string(ua)
}

func (ua UserAgent) Contains(keyword string) bool {
	return strings.Contains(ua.Native(), keyword)
}

func (ua UserAgent) ContainsAny(keywords types.Strings) bool {
	return stl.AnyMatch(keywords, ua.Contains)
}

func (ua UserAgent) ContainsAnyX(keywords ...string) bool {
	return ua.ContainsAny(keywords)
}

func (ua UserAgent) IsPc() bool {
	return ua.ContainsAnyX(
		"windows",
		"macintosh",
		"x86",
		"x86_64",
		"x64",
	)
}

func (ua UserAgent) IsWxWork() bool {
	return ua.ContainsAnyX(
		"wxwork",
	)
}

func (ua UserAgent) IsWeixin() bool {
	return ua.ContainsAnyX(
		"weixin",
	)
}

func (ua UserAgent) IsWxWorkOrWeixin() bool {
	return ua.IsWxWork() || ua.IsWeixin()
}

func (ua UserAgent) IsWindows() bool {
	return ua.ContainsAnyX(
		"windows",
	)
}

func (ua UserAgent) IsMac() bool {
	return ua.ContainsAnyX(
		"macintosh",
	)
}
