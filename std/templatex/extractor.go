package templatex

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fasionchan/goutils/types"
)


const (
	ExtractingContainerName = "Result"
)

type TemplateDataExtractor struct {
	template   *SmartTemplate
	mutex      *sync.Mutex
	container  *DataContainer               // 第一个数据（兼容单例场景）
	containers DataContainerMappingByString // 数据映射表（多例场景）
}

func NewTemplateDataExtractorCustom(tplText string, funcMap TemplateFuncMap, concurrent bool, name string, containerNames ...string) (*TemplateDataExtractor, error) {
	funcMap = TemplateFuncs.Dup().Merge(funcMap)

	// 容器名去空、去重
	containerNames = types.Strings(containerNames).PurgeZero().UniqueBySet()
	if len(containerNames) == 0 {
		// 未定义则采用默认名字
		containerNames = []string{ExtractingContainerName}
	}

	// 初始化容器
	containers := funcMap.WithDataContainersForMapping(containerNames...)

	// 解析模板
	tpl, err := ParseSmartTemplateFromText(name, tplText, funcMap.Native())
	if err != nil {
		return nil, err
	}

	// 分配锁
	var mutex *sync.Mutex
	if concurrent {
		mutex = &sync.Mutex{}
	}

	return &TemplateDataExtractor{
		template:   tpl,
		container:  containers[containerNames[0]], // 以第一个容器为单例
		containers: containers,                    // 多例
		mutex:      mutex,
	}, nil
}

func NewTemplateDataExtractor(tplText string, funcMap TemplateFuncMap, concurrent bool) (*TemplateDataExtractor, error) {
	return NewTemplateDataExtractorCustom(tplText, funcMap, concurrent, "", "")
}

func NewTemplateDataExtractorByPathCustom(path string, funcMap TemplateFuncMap, concurrent bool, name string, containerName string) (*TemplateDataExtractor, error) {
	if path == "" {
		path = "."
	}
	if !strings.HasPrefix(path, ".") {
		path = "." + path
	}
	if containerName == "" {
		containerName = ExtractingContainerName
	}
	return NewTemplateDataExtractorByExpressionCustom(path, "", funcMap, concurrent, name, containerName)
}

func NewTemplateDataExtractorByPath(path string, funcMap TemplateFuncMap, concurrent bool) (*TemplateDataExtractor, error) {
	return NewTemplateDataExtractorByPathCustom(path, funcMap, concurrent, "", "")
}

func (extractor *TemplateDataExtractor) Extract(data any) (any, error) {
	if err := extractor.ClearAndRender(data); err != nil {
		return nil, err
	}

	return extractor.container.Get(), nil
}

func (extractor *TemplateDataExtractor) ExtractMany(data any) (map[string]any, error) {
	if err := extractor.ClearAndRender(data); err != nil {
		return nil, err
	}

	return extractor.containers.GetAllDatas(), nil
}

func (extractor *TemplateDataExtractor) ClearAndRender(data any) error {
	if mutex := extractor.mutex; mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	extractor.containers.Clear()
	return extractor.template.RenderToDiscard(data)
}

func ExtractDataByTemplateCustom(data any, tplText string, funcMap TemplateFuncMap, tplName string, containerName string) (any, error) {
	extractor, err := NewTemplateDataExtractorCustom(tplText, funcMap, false, tplName, containerName)
	if err != nil {
		return nil, err
	}

	return extractor.Extract(data)
}

func ExtractDatasByTemplateCustom(data any, tplText string, funcMap TemplateFuncMap, tplName string, containerNames ...string) (map[string]any, error) {
	extractor, err := NewTemplateDataExtractorCustom(tplText, funcMap, false, tplName, containerNames...)
	if err != nil {
		return nil, err
	}

	return extractor.ExtractMany(data)
}

func ExtractDataByTemplate(data any, tplText string, funcMap TemplateFuncMap) (any, error) {
	return ExtractDataByTemplateCustom(data, tplText, funcMap, "", "")
}

func ExtractDataByPathTempalteCustom(data any, path string, funcMap TemplateFuncMap, tplName string, containerName string) (any, error) {
	if path == "" || path == "." {
		return data, nil
	}

	extractor, err := NewTemplateDataExtractorByPathCustom(path, funcMap, false, tplName, containerName)
	if err != nil {
		return nil, err
	}

	return extractor.Extract(data)
}

func ExtractDataByPathTempalte(data any, path string, funcMap TemplateFuncMap) (any, error) {
	return ExtractDataByPathTempalteCustom(data, path, funcMap, "", "")
}

func ExtractDataByExpressionTempalte(data any, expression, pretreatment string, funcMap TemplateFuncMap) (any, error) {
	return ExtractDataByExpressionTempalteCustom(data, expression, pretreatment, funcMap, "", "")
}

func ExtractDataByExpressionTempalteCustom(data any, expression, pretreatment string, funcMap TemplateFuncMap, tplName string, containerName string) (any, error) {
	extractor, err := NewTemplateDataExtractorByExpressionCustom(expression, pretreatment, funcMap, false, tplName, containerName)
	if err != nil {
		return nil, err
	}

	return extractor.Extract(data)
}

func NewTemplateDataExtractorByExpressionCustom(expression, pretreatment string, funcMap TemplateFuncMap, concurrent bool, name string, containerName string) (*TemplateDataExtractor, error) {
	if containerName == "" {
		containerName = ExtractingContainerName
	}

	if expression == "" {
		expression = "."
	}

	setExpression := fmt.Sprintf("{{ set%s %s }}", containerName, expression)

	return NewTemplateDataExtractorCustom(fmt.Sprintf("%s %s", pretreatment, setExpression), funcMap, concurrent, name, containerName)
}

type DataExtractorTemplate struct {
	PretreatmentTemplate string `bson:"PretreatmentTemplate" json:"PretreatmentTemplate"`
	ExpressionTemplate   string `bson:"ExpressionTemplate" json:"ExpressionTemplate"`
}

func (template *DataExtractorTemplate) NewDataExtractor(funcMap TemplateFuncMap, concurrent bool, name string, containerName string) (*TemplateDataExtractor, error) {
	return NewTemplateDataExtractorByExpressionCustom(template.ExpressionTemplate, template.PretreatmentTemplate, funcMap, concurrent, name, containerName)
}

func (template *DataExtractorTemplate) Extract(data any, funcMap TemplateFuncMap) (any, error) {
	if template == nil {
		return data, nil
	}

	if template.ExpressionTemplate == "" && template.PretreatmentTemplate == "" {
		return data, nil
	}

	extractor, err := template.NewDataExtractor(funcMap, false, "", "")
	if err != nil {
		return nil, err
	}

	return extractor.Extract(data)
}