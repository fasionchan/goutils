package templatex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"text/template"

	"github.com/fasionchan/goutils/baseutils/netutils"
	"github.com/fasionchan/goutils/stl"
	"gopkg.in/yaml.v3"
)

const (
	ContextValueTemplateData  = "TemplateData"

	TemplateValueNoValue = "<no value>"
)


func WrapContextWithTemplateData(ctx context.Context, data any) context.Context {
	return context.WithValue(ctx, ContextValueTemplateData, data)
}

func TemplateDataFromContext(ctx context.Context) any {
	return ctx.Value(ContextValueTemplateData)
}

// ✅ SmartTemplate is a template that can be rendered to a buffer, bytes, string, etc.
type SmartTemplate template.Template

func MustSmartTemplate(tpl *SmartTemplate, err error) *SmartTemplate {
	if err != nil {
		panic(err)
	}
	return tpl
}

func ParseSmartTemplateFromText(name, text string, funcMap template.FuncMap) (*SmartTemplate, error) {
	tpl := template.New(name)
	if funcMap != nil {
		tpl = tpl.Funcs(funcMap)
	}

	tpl, err := tpl.Parse(text)
	if err != nil {
		return nil, err
	}

	return (*SmartTemplate)(tpl), nil
}

func MustParseSmartTemplateFromText(name, text string, funcMap template.FuncMap) *SmartTemplate {
	return MustSmartTemplate(ParseSmartTemplateFromText(name, text, funcMap))
}

func (tpl *SmartTemplate) Native() *template.Template {
	return (*template.Template)(tpl)
}

func (tpl *SmartTemplate) WithFunc(funcMap template.FuncMap) *SmartTemplate {
	return (*SmartTemplate)(tpl.Native().Funcs(funcMap))
}

func (tpl *SmartTemplate) RenderToBuffer(data any) (buffer *bytes.Buffer, err error) {
	buffer = bytes.NewBuffer(nil)
	err = tpl.Native().Execute(buffer, data)
	return
}

func (tpl *SmartTemplate) RenderToBytes(data any) ([]byte, error) {
	buffer, err := tpl.RenderToBuffer(data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (tpl *SmartTemplate) RenderToDiscard(data any) error {
	return tpl.Native().Execute(io.Discard, data)
}

func (tpl *SmartTemplate) RenderToInt(data any) (int, error) {
	str, err := tpl.RenderToString(data)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (tpl *SmartTemplate) RenderToJson(data, result any) error {
	bytes, err := tpl.RenderToBytes(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}

func (tpl *SmartTemplate) RenderToString(data any) (string, error) {
	buffer, err := tpl.RenderToBuffer(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (tpl *SmartTemplate) RenderToYaml(data, result any) error {
	bytes, err := tpl.RenderToBytes(data)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, result)
}

type SmartTemplates = []*SmartTemplate

type SmartTemplateMappingByString map[string]*SmartTemplate

func NewSmartTemplateMappingByString() SmartTemplateMappingByString {
	return SmartTemplateMappingByString{}
}

func ParseAndRenderNamedTemplates(templates map[string]string, funcMap TemplateFuncMap, data interface{}) (map[string]string, error) {
	result, _, err := stl.MapMapPro(templates, func(templateName string, templateValue string, _, _ map[string]string) (string, string, bool, error) {
		value, err := funcMap.ParseTemplateAndRenderToString(templateName, templateValue, data)
		if err != nil {
			return "", "", false, fmt.Errorf("RenderTemplateFailed, templateName=%s, err:%w", templateName, err)
		}
		return templateName, value, true, nil
	})
	return result, err
}

func ParseAndRenderNamedTemplatesAsUrlValues(templates map[string]string, funcMap TemplateFuncMap, data interface{}) (url.Values, error) {
	values, err := ParseAndRenderNamedTemplates(templates, funcMap, data)
	if err != nil {
		return nil, err
	}
	return netutils.NewUrlValuesFromMap(values), nil
}
