/*
 * Author: fasion
 * Created time: 2025-10-14 09:41:32
 * Last Modified by: fasion
 * Last Modified time: 2025-10-27 10:22:42
 */

package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/fasionchan/goutils/libs/datarender"
	"github.com/fasionchan/goutils/std/templatex"
)

type Request struct {
	Method    string      `bson:"Method" json:"Method"`
	Url       string      `bson:"Url" json:"Url"`
	Header    Header      `bson:"Header" json:"Header"`
	BodyBytes []byte      `bson:"BodyBytes" json:"BodyBytes"`
}

func (request *Request) NewHttpRequest(ctx context.Context) (*http.Request, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, request.Method, request.Url, bytes.NewReader(request.BodyBytes))
	if err != nil {
		return nil, err
	}

	if request.Header != nil {
		for key, values := range request.Header {
			for _, value := range values {
				httpRequest.Header.Add(key, value)
			}
		}
	}

	return httpRequest, nil
}

func (request *Request) NewHttpServerRequest(ctx context.Context) (*http.Request, error) {
	parsedUrl, err := url.Parse(request.Url)
	if err != nil {
		return nil, err
	}

	header := request.Header
	if header == nil {
		header = make(Header)
	}

	body := io.NopCloser(bytes.NewReader(request.BodyBytes))

	return &http.Request{
		Method: request.Method,
		URL:    parsedUrl,
		Header: header.Native(),
		Body:   body,
	}, nil
}

func (request *Request) Print() {
	fmt.Println("Method:", request.Method)
	fmt.Println("Url:", request.Url)
	fmt.Println("Header:", request.Header)
	fmt.Println("Body:", len(request.BodyBytes))
}

func (request *Request) WithRange(start, end int64) *Request {
	if request == nil {
		return nil
	}

	request.Header = request.Header.WithRange(start, end)
	return request
}

type RequestTemplate struct {
	Method string         `bson:"Method" json:"Method"`
	Url    UrlTemplate    `bson:"Url" json:"Url"`
	Header HeaderTemplate `bson:"Header" json:"Header"`
	Body   BodyTemplate   `bson:"Body" json:"Body"`

	UrlRender    UrlRender    `bson:"-" json:"-"`
	HeaderRender HeaderRender `bson:"-" json:"-"`
	BodyRender   BodyRender   `bson:"-" json:"-"`
}

func NewRequestBlank() *RequestTemplate {
	return &RequestTemplate{}
}

func NewRequestTemplate(method string, url UrlTemplate, header HeaderTemplate, body BodyTemplate) *RequestTemplate {
	return &RequestTemplate{
		Method: method,
		Url:    url,
		Header: header,
		Body:   body,
	}
}

func ParseRequestTemplate(method string, url UrlTemplate, header HeaderTemplate, body BodyTemplate, funcMap templatex.TemplateFuncMap) (*RequestTemplate, error) {
	tpl := NewRequestTemplate(method, url, header, body)
	if err := tpl.Parse(funcMap); err != nil {
		return nil, err
	}

	return tpl, nil
}

func (tpl *RequestTemplate) Parse(funcMap templatex.TemplateFuncMap) error {
	if err := tpl.ParseUrl(funcMap); err != nil {
		return err
	}

	if err := tpl.ParseHeader(funcMap); err != nil {
		return err
	}

	if err := tpl.ParseBody(funcMap); err != nil {
		return err
	}

	return nil
}

func (tpl *RequestTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, data any) (*Request, error) {
	if err := tpl.Parse(funcMap); err != nil {
		return nil, err
	}

	return tpl.Render(data)
}

func (tpl *RequestTemplate) ParseUrl(funcMap templatex.TemplateFuncMap) error {
	if tpl == nil {
		return nil
	}

	urlTpl := tpl.Url
	if urlTpl == "" {
		return nil
	}

	render, err := urlTpl.ParseForRender(funcMap, false)
	if err != nil {
		return err
	}

	tpl.UrlRender = render

	return nil
}

func (tpl *RequestTemplate) ParseHeader(funcMap templatex.TemplateFuncMap) error {
	if tpl == nil {
		return nil
	}

	headerTpl := tpl.Header
	if headerTpl == "" {
		return nil
	}

	render, err := headerTpl.ParseForRender(funcMap, true)
	if err != nil {
		return err
	}

	tpl.HeaderRender = render

	return nil
}

func (tpl *RequestTemplate) ParseBody(funcMap templatex.TemplateFuncMap) error {
	if tpl == nil {
		return nil
	}

	bodyTpl := tpl.Body
	if bodyTpl == "" {
		return nil
	}

	render, err := bodyTpl.ParseForRender(funcMap, false)
	if err != nil {
		return err
	}

	tpl.BodyRender = render

	return nil
}

func (tpl *RequestTemplate) Print() {
	if tpl == nil {
		return
	}

	fmt.Println("Method:", tpl.Method)
	fmt.Println("Url:", tpl.Url)
	fmt.Println("Header:", tpl.Header)
	fmt.Println("Body:", len(tpl.Body))
}

func (tpl *RequestTemplate) Render(data any) (*Request, error) {
	url, err := tpl.RenderUrl(data)
	if err != nil {
		return nil, err
	}

	header, err := tpl.RenderHeader(data)
	if err != nil {
		return nil, err
	}

	body, err := tpl.RenderBody(data)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method:    tpl.Method,
		Url:       url,
		Header:    Header(header),
		BodyBytes: body,
	}, nil
}

func (tpl *RequestTemplate) RenderUrl(data any) (string, error) {
	if tpl == nil {
		return "", nil
	}

	render := tpl.UrlRender
	if render == nil {
		return "", nil
	}

	return render.Render(data)
}

func (tpl *RequestTemplate) RenderHeader(data any) (http.Header, error) {
	if tpl == nil {
		return nil, nil
	}

	render := tpl.HeaderRender
	if render == nil {
		return nil, nil
	}

	return render.Render(data)
}

func (tpl *RequestTemplate) RenderBody(data any) ([]byte, error) {
	if tpl == nil {
		return nil, nil
	}

	render := tpl.BodyRender
	if render == nil {
		return nil, nil
	}

	body, err := render.Render(data)
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}

func (tpl *RequestTemplate) WithMethod(method string) *RequestTemplate {
	if tpl == nil {
		return nil
	}

	tpl.Method = method
	return tpl
}

func (tpl *RequestTemplate) WithUrl(url UrlTemplate) *RequestTemplate {
	if tpl == nil {
		return nil
	}

	tpl.Url = url
	return tpl
}

func (tpl *RequestTemplate) WithHeader(header HeaderTemplate) *RequestTemplate {
	if tpl == nil {
		return nil
	}

	tpl.Header = header
	return tpl
}

func (tpl *RequestTemplate) WithBodyBytes(body BodyTemplate) *RequestTemplate {
	if tpl == nil {
		return nil
	}

	tpl.Body = body
	return tpl
}

const (
	BodyTemplateTypePlain = "\t\r\n"
)

func ParseBodyText(text string) (string, error) {
	return text, nil
}

// 模板
type BodyTemplate string

func (tpl BodyTemplate) Native() string {
	return string(tpl)
}

func (tpl BodyTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, concurrent bool) (BodyRender, error) {
	return ParseBodyRenderTemplate(tpl.Native(), funcMap, concurrent)
}

func (tpl BodyTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, data any) (string, error) {
	return ParseAndRenderBodyTemplate(tpl.Native(), funcMap, data)
}

// 渲染器
type BodyRender = datarender.StringRender

func ParseBodyRenderTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (BodyRender, error) {
	return datarender.ParseStringRenderTemplate(text, funcMap, "Body", concurrent)
}

func ParseAndRenderBodyTemplate(text string, funcMap templatex.TemplateFuncMap, data any) (string, error) {
	return datarender.ParseAndRenderStringTemplate(text, funcMap, "Body", data)
}

// 文本型模板
type BodyTextTemplate = datarender.StringTextTemplate

func ParseBodyTextTemplate(text string, funcMap templatex.TemplateFuncMap) (*BodyTextTemplate, error) {
	return datarender.ParseDataTextTemplate(text, funcMap, ParseBodyText)
}

// 提取型模板
type BodyExtractorTemplate = datarender.StringExtractorTemplate

func ParseBodyExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool) (*BodyExtractorTemplate, error) {
	return datarender.ParseDataExtractorTemplate[string](text, funcMap, concurrent, "Body", nil)
}
