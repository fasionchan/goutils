package templatex

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

// ✅ TemplateFuncMap is a map of template functions.
type TemplateFuncMap template.FuncMap

func NewTemplateFuncMap() TemplateFuncMap {
	return TemplateFuncMap{}
}

func (m TemplateFuncMap) Len() int {
	return len(m)
}

func (m TemplateFuncMap) Empty() bool {
	return m.Len() == 0
}

func (m TemplateFuncMap) Native() template.FuncMap {
	return template.FuncMap(m)
}

func (m TemplateFuncMap) Dup() TemplateFuncMap {
	return stl.DupMap(m)
}

func (m TemplateFuncMap) Keys() types.Strings {
	return stl.MapKeys(m)
}

func (m TemplateFuncMap) Merge(other TemplateFuncMap) TemplateFuncMap {
	return stl.ConcatMapsTo(m, other)
}

func (m TemplateFuncMap) MergeManyX(others ...TemplateFuncMap) TemplateFuncMap {
	return stl.ConcatMapsTo(m, others...)
}

func (m TemplateFuncMap) SubMapByKeys(keys ...string) TemplateFuncMap {
	return stl.SubMapByKeys(m, keys...)
}

func (m TemplateFuncMap) With(key string, f any) TemplateFuncMap {
	m[key] = f
	return m
}

// WithDataContainer returns a DataContainer with the given name with the Set function registered.
func (m TemplateFuncMap) WithDataContainer(name string) *DataContainer {
	container := NewDataContainer()
	m["set"+name] = container.Set
	return container
}

// WithDataContainers returns a slice of DataContainers with the given names with the Set functions registered.
func (m TemplateFuncMap) WithDataContainers(names ...string) DataContainers {
	return stl.Map(names, m.WithDataContainer)
}

// WithDataContainersForMapping returns a mapping of DataContainers by string with the given names with the Set functions registered.
func (m TemplateFuncMap) WithDataContainersForMapping(names ...string) DataContainerMappingByString {
	return stl.BuildMap[DataContainerMappingByString](names, func(name string) (string, *DataContainer) {
		return name, m.WithDataContainer(name)
	})
}

// ParseTemplate parses the given text as a template and returns a SmartTemplate.
func (m TemplateFuncMap) ParseTemplate(name, text string) (*SmartTemplate, error) {
	return ParseSmartTemplateFromText(name, text, m.Native())
}

// ParseTemplateAndRenderToBuffer parses the given text as a template and renders it to a buffer.
func (m TemplateFuncMap) ParseTemplateAndRenderToBuffer(name, text string, data any) (*bytes.Buffer, error) {
	tpl, err := m.ParseTemplate(name, text)
	if err != nil {
		return nil, err
	}

	return tpl.RenderToBuffer(data)
}

// ParseTemplateAndRenderToBytes parses the given text as a template and renders it to a bytes.
func (m TemplateFuncMap) ParseTemplateAndRenderToBytes(name, text string, data any) ([]byte, error) {
	tpl, err := m.ParseTemplate(name, text)
	if err != nil {
		return nil, err
	}

	return tpl.RenderToBytes(data)
}

// ParseTemplateAndRenderToString parses the given text as a template and renders it to a string.
func (m TemplateFuncMap) ParseTemplateAndRenderToString(name, text string, data any) (string, error) {
	tpl, err := m.ParseTemplate(name, text)
	if err != nil {
		return "", err
	}

	return tpl.RenderToString(data)
}

func (m TemplateFuncMap) ParseTemplateAndRenderToTrimSpaceString(name, text string, data any) (string, error) {
	str, err := m.ParseTemplateAndRenderToString(name, text, data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}

type TemplateFuncMaps []TemplateFuncMap

func (maps TemplateFuncMaps) Merge() TemplateFuncMap {
	return stl.ConcatMaps(maps...)
}
