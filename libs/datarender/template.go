/*
 * Author: fasion
 * Created time: 2025-10-16 11:41:37
 * Last Modified by: fasion
 * Last Modified time: 2026-08-02 19:56:51
 */

package datarender

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/std/_time"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

const (
	DataParserTypeJson    = "json"
	DataParserTypeYaml    = "yaml"
	DataParserTypeCsv     = "csv"
	DataParserTypeDefault = "default"
)

// 数据文本型模板
type DataTextTemplate[Data any] struct {
	text string
	*templatex.SmartTemplate
	dataParser func(text string) (Data, error)
}

func ParseDataTextTemplate[Data any](text string, funcMap templatex.TemplateFuncMap, dataParser func(text string) (Data, error)) (*DataTextTemplate[Data], error) {
	tpl, err := funcMap.ParseTemplate("", text)
	if err != nil {
		return nil, err
	}

	return &DataTextTemplate[Data]{
		text:          text,
		SmartTemplate: tpl,
		dataParser:    dataParser,
	}, nil
}

func (tpl *DataTextTemplate[Data]) Serialize(data any) (result Data, err error) {
	text, err := tpl.RenderToString(data)
	if err != nil {
		return
	}

	return tpl.dataParser(text)
}

func (tpl *DataTextTemplate[Data]) Render(data any) (Data, error) {
	return tpl.Serialize(data)
}

// 数据提取型模板
type DataExtractorTemplate[Data any] struct {
	text string
	*templatex.TemplateDataExtractor
	dataCast func(any) (Data, error)
}

func ParseDataExtractorTemplate[Data any](text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (Data, error)) (*DataExtractorTemplate[Data], error) {
	extractor, err := templatex.NewTemplateDataExtractorCustom(text, funcMap, concurrent, "", containerName)
	if err != nil {
		return nil, err
	}

	return &DataExtractorTemplate[Data]{
		text:                  text,
		TemplateDataExtractor: extractor,
		dataCast:              dataCast,
	}, nil
}

func (tpl *DataExtractorTemplate[Data]) Extract(data any) (result Data, err error) {
	anyResult, err := tpl.TemplateDataExtractor.Extract(data)
	if err != nil {
		return
	}

	if dataCast := tpl.dataCast; dataCast != nil {
		return dataCast(anyResult)
	}

	resultType := reflect.TypeOf(result)

	value := reflect.ValueOf(anyResult)
	for {
		if value.Type() == resultType {
			return value.Interface().(Data), nil
		}

		if value.CanConvert(resultType) {
			return value.Convert(resultType).Interface().(Data), nil
		}

		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				err = baseutils.NewNilError(fmt.Sprintf("pointer is nil: %s", value.Type().String()))
				return
			}

			value = value.Elem()
		} else {
			err = baseutils.NewBlankBadTypeError().WithExpectedReflectType(resultType).WithGivenReflectType(value.Type())
			return
		}
	}
}

func (tpl *DataExtractorTemplate[Data]) Render(data any) (result Data, err error) {
	return tpl.Extract(data)
}

// 数据渲染器
type DataRender[Data any] interface {
	Render(data any) (Data, error)
}

// 数据渲染函数
type DataRenderFunc[Data any] func(data any) (Data, error)

func (fn DataRenderFunc[Data]) Render(data any) (Data, error) {
	return fn(data)
}

// 解析数据渲染模板
func ParseDataRenderTemplate[Data any](
	text string, // 模板文本
	funcMap templatex.TemplateFuncMap, // 模板函数映射
	plainMagic string, // 纯文本魔术字符串（表示不使用模板）
	dataParser func(text string) (Data, error), // 数据解析函数（数据文本型模板使用）
	dataContainerName string, // 数据容器名称（数据提取型模板使用）
	concurrent bool, // 是否支持并发执行（数据提取型模板使用，控制是否需要加锁）
	dataCast func(any) (Data, error), // 数据转换函数（数据提取型模板使用）
) (DataRender[Data], error) {
	dataParserMapping := NewDataParserFuncMappingByString[Data]().WithDefault(dataParser)
	return ParseDataRenderTemplatePro(text, funcMap, plainMagic, dataParserMapping, dataContainerName, concurrent, dataCast)

	// switch {
	// case strings.HasPrefix(text, plainMagic):
	// 	return DataRenderFunc[Data](func(data any) (Data, error) {
	// 		return dataParser(strings.TrimPrefix(text, plainMagic))
	// 	}), nil
	// case strings.HasPrefix(text, TemplateTypeMagicLua):
	// 	return ParseLuaTextRenderTemplate(text, dataParser)
	// case strings.HasPrefix(text, TemplateTypeMagicExtractor):
	// 	return ParseDataExtractorTemplate(text, funcMap, concurrent, dataContainerName, dataCast)
	// default:
	// 	return ParseDataTextTemplate(text, funcMap, dataParser)
	// }
}

func ParseDataRenderTemplatePro[Data any](
	text string, // 模板文本
	funcMap templatex.TemplateFuncMap, // 模板函数映射
	plainMagic string, // 纯文本魔术字符串（表示不使用模板）
	dataParserMapping DataParserFuncMappingByString[Data], // 数据解析函数映射
	dataContainerName string, // 数据容器名称（数据提取型模板使用）
	concurrent bool, // 是否支持并发执行（数据提取型模板使用，控制是否需要加锁）
	dataCast func(any) (Data, error), // 数据转换函数（数据提取型模板使用）
) (DataRender[Data], error) {
	magic, err := ParseTemplateMagic(text, plainMagic)
	if err != nil {
		return nil, err
	}

	switch magic.Type {
	case TemplateTypePlain:
		return DataRenderFunc[Data](func(data any) (Data, error) {
			return dataParserMapping.Parse(magic.GetFormat(), strings.TrimPrefix(text, plainMagic))
		}), nil
	case TemplateTypeLua:
		return ParseLuaTextRenderTemplate(text, dataParserMapping.GetParser(magic.GetFormat()))
	case TemplateTypeExtractor:
		return ParseDataExtractorTemplate(text, funcMap, concurrent, dataContainerName, dataCast)
	case TemplateTypeText:
		return ParseDataTextTemplate(text, funcMap, dataParserMapping.GetParser(magic.GetFormat()))
	default:
		return nil, baseutils.NewNotImplementedError("TemplateType:" + strconv.Itoa(magic.Type))
	}
}

type DataTemplate[Data any] interface {
	Native() string
	ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (DataRender[Data], error)
	ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any) (Data, error)
}

type DataTemplateHandler[Data any] struct {
	tpl    DataTemplate[Data]
	render DataRender[Data]
	result Data

	funcMap       templatex.TemplateFuncMap
	containerName string
}

func NewDataTemplateHandler[Data any](tpl DataTemplate[Data], funcMap templatex.TemplateFuncMap, containerName string) *DataTemplateHandler[Data] {
	return &DataTemplateHandler[Data]{
		tpl:           tpl,
		funcMap:       funcMap,
		containerName: containerName,
	}
}

func (handler *DataTemplateHandler[Data]) AsTypelessDataTemplateHandler() TypelessDataTemplateHandler {
	return handler
}

func (handler *DataTemplateHandler[Data]) ParseAndRender(data any) (err error) {
	_, err = handler.ParseAndRenderForData(data)
	return
}

func (handler *DataTemplateHandler[Data]) ParseAndRenderForData(data any) (Data, error) {
	var zero Data
	if err := handler.Parse(false); err != nil {
		return zero, err
	}

	return handler.RenderForData(data)
}

func (handler *DataTemplateHandler[Data]) Parse(concurrent bool) (err error) {
	_, err = handler.ParseForRender(concurrent)
	return
}

func (handler *DataTemplateHandler[Data]) ParseForRender(concurrent bool) (DataRender[Data], error) {
	if handler == nil {
		return nil, baseutils.NewNilError("DataTemplateHandler.ParseForRender(handler)")
	}

	render, err := handler.tpl.ParseForRender(handler.funcMap, handler.containerName, concurrent)
	if err != nil {
		return nil, err
	}

	handler.render = render
	return render, nil
}

func (handler *DataTemplateHandler[Data]) Render(data any) error {
	_, err := handler.RenderForData(data)
	return err
}

func (handler *DataTemplateHandler[Data]) RenderForData(data any) (Data, error) {
	var zero Data

	if handler == nil {
		return zero, baseutils.NewNilError("DataTemplateHandler.Render(handler)")
	}

	render := handler.render
	if render == nil {
		return zero, baseutils.NewNilError("DataTemplateHandler.Render(handler.render)")
	}

	result, err := render.Render(data)
	if err != nil {
		return zero, err
	}

	handler.result = result
	return result, nil
}

type TypelessDataTemplateHandler interface {
	Parse(concurrent bool) error
	Render(data any) error
}

type TypelessDataTemplateHandlers []TypelessDataTemplateHandler

func NewTypelessDataTemplateHandlers(handlers ...TypelessDataTemplateHandler) TypelessDataTemplateHandlers {
	return handlers
}

func (handlers TypelessDataTemplateHandlers) Parse(concurrent bool) error {
	return stl.BatchProcessUntilFirstError(handlers, func(handler TypelessDataTemplateHandler) error {
		return handler.Parse(concurrent)
	})
}

func (handlers TypelessDataTemplateHandlers) Render(data any) error {
	return stl.BatchProcessUntilFirstError(handlers, func(handler TypelessDataTemplateHandler) error {
		return handler.Render(data)
	})
}

type GenericTemplatedData[DataRender_ DataRender[Data], Data any] interface {
	Native() string
	ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (DataRender_, error)
	ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any) (Data, error)
}

type StringTemplate string

func (tpl StringTemplate) Native() string {
	return string(tpl)
}

func (tpl StringTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (StringRender, error) {
	return ParseStringRenderTemplate(tpl.Native(), funcMap, containerName, concurrent)
}

func (tpl StringTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any) (string, error) {
	return ParseAndRenderStringTemplate(tpl.Native(), funcMap, containerName, data)
}

func (tpl StringTemplate) ParseAndRenderPurgeSpace(funcMap templatex.TemplateFuncMap, containerName string, data any) (string, error) {
	result, err := ParseAndRenderStringTemplate(tpl.Native(), funcMap, containerName, data)
	if err != nil {
		return "", err
	}

	return PurgeStringSpace(result), nil
}

var spaceReg = regexp.MustCompile(`\s|\\n|\\t|\\r|\\f|\\v`)

// PurgeStrSpace 使用正则清空所有代表空格的字符
func PurgeStringSpace(s string) string {
	return spaceReg.ReplaceAllString(s, "")
}

// 字符串渲染器
type StringRender = DataRender[string]

func ParseStringRenderTemplate(text string, funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (StringRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, parseStringText, containerName, concurrent, nil)
}

func ParseAndRenderStringTemplate(text string, funcMap templatex.TemplateFuncMap, containerName string, data any) (string, error) {
	render, err := ParseStringRenderTemplate(text, funcMap, containerName, false)
	if err != nil {
		return "", err
	}
	return render.Render(data)
}

// 字符串文本型模板
type StringTextTemplate = DataTextTemplate[string]

func ParseStringTextTemplate(text string, funcMap templatex.TemplateFuncMap) (*StringTextTemplate, error) {
	return ParseDataTextTemplate(text, funcMap, parseStringText)
}

func parseStringText(text string) (string, error) {
	return text, nil
}

// 字符串提取型模板
type StringExtractorTemplate = DataExtractorTemplate[string]

func ParseStringExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (string, error)) (*StringExtractorTemplate, error) {
	return ParseDataExtractorTemplate(text, funcMap, concurrent, containerName, dataCast)
}

type StringSeparator string

func (sep StringSeparator) Split(s string) types.Strings {
	return types.SplitToStrings(s, string(sep))
}

func (sep StringSeparator) SplitWithError(s string) (types.Strings, error) {
	return sep.Split(s), nil
}

func (sep StringSeparator) SplitTpe(s string) (types.Strings, error) {
	return sep.Split(s).TrimSpace().PurgeZero(), nil
}

type StringsTemplate string

func (tpl StringsTemplate) Native() string {
	return string(tpl)
}

func (tpl StringsTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool, parseStringsText func(text string) (types.Strings, error), dataCast func(any) (types.Strings, error)) (StringsRender, error) {
	return ParseDataRenderTemplate(tpl.Native(), funcMap, TemplateTypeMagicPlain, parseStringsText, containerName, concurrent, dataCast)
}

func (tpl StringsTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any, parseStringsText func(text string) (types.Strings, error), dataCast func(any) (types.Strings, error)) (types.Strings, error) {
	return ParseAndRenderStringsTemplate(tpl.Native(), funcMap, parseStringsText, dataCast, containerName, data)
}

type StringsRender = DataRender[types.Strings]

func ParseStringsRenderTemplate(text string, funcMap templatex.TemplateFuncMap, parseStringsText func(text string) (types.Strings, error), dataCast func(any) (types.Strings, error), containerName string, concurrent bool) (StringsRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, parseStringsText, containerName, concurrent, dataCast)
}
func ParseAndRenderStringsTemplate(text string, funcMap templatex.TemplateFuncMap, parseStringsText func(text string) (types.Strings, error), dataCast func(any) (types.Strings, error), containerName string, data any) (types.Strings, error) {
	render, err := ParseStringsRenderTemplate(text, funcMap, parseStringsText, dataCast, containerName, false)
	if err != nil {
		return nil, err
	}
	return render.Render(data)
}

type StringsTextTemplate = DataTextTemplate[types.Strings]

func ParseStringsTextTemplate(text string, funcMap templatex.TemplateFuncMap, parseStringsText func(text string) (types.Strings, error)) (*StringsTextTemplate, error) {
	return ParseDataTextTemplate(text, funcMap, parseStringsText)
}

type StringsExtractorTemplate = DataExtractorTemplate[types.Strings]

func ParseStringsExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (types.Strings, error)) (*StringsExtractorTemplate, error) {
	return ParseDataExtractorTemplate(text, funcMap, concurrent, containerName, dataCast)
}

type SeparatedStringsTemplate string

func (tpl SeparatedStringsTemplate) Native() string {
	return string(tpl)
}

func (tpl SeparatedStringsTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, seps Separators, containerName string, concurrent bool) (StringsRender, error) {
	return ParseStringsRenderTemplate(tpl.Native(), funcMap, NewSeparatedStringsParseFunc(seps), nil, containerName, concurrent)
}

func (tpl SeparatedStringsTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, seps Separators, containerName string, data any) (types.Strings, error) {
	return ParseAndRenderStringsTemplate(tpl.Native(), funcMap, NewSeparatedStringsParseFunc(seps), nil, containerName, data)
}

type Separators types.Strings

func NewSeparators(seps ...string) Separators {
	return Separators(seps)
}

func (seps Separators) SplitText(text string) types.Strings {
	return types.NewStrings(text).Split(seps...)
}

func (seps Separators) SplitTexts(texts types.Strings) types.Strings {
	return texts.Split(seps...)
}

func NewSeparatedStringsParseFuncX(seps ...string) func(text string) (types.Strings, error) {
	return NewSeparatedStringsParseFunc(Separators(seps))
}

func NewSeparatedStringsParseFunc(seps Separators) func(text string) (types.Strings, error) {
	return func(text string) (types.Strings, error) {
		return seps.SplitText(text), nil
	}
}

type TimeTemplate string

func (tpl TimeTemplate) Native() string {
	return string(tpl)
}

func (tpl TimeTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, parseTimeText func(text string) (time.Time, error), containerName string, concurrent bool) (TimeRender, error) {
	return ParseTimeRenderTemplate(tpl.Native(), funcMap, parseTimeText, containerName, concurrent)
}

func (tpl TimeTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, parseTimeText func(text string) (time.Time, error), containerName string, data any) (time.Time, error) {
	return ParseAndRenderTimeTemplate(tpl.Native(), funcMap, parseTimeText, containerName, data)
}

// 时间渲染器
type TimeRender = DataRender[time.Time]

func ParseTimeRenderTemplate(text string, funcMap templatex.TemplateFuncMap, parseTimeText func(text string) (time.Time, error), containerName string, concurrent bool) (TimeRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, parseTimeText, containerName, concurrent, nil)
}

func ParseAndRenderTimeTemplate(text string, funcMap templatex.TemplateFuncMap, parseTimeText func(text string) (time.Time, error), containerName string, data any) (time.Time, error) {
	render, err := ParseTimeRenderTemplate(text, funcMap, parseTimeText, containerName, false)
	if err != nil {
		return time.Time{}, err
	}
	return render.Render(data)
}

// 时间文本型模板
type TimeTextTemplate = DataTextTemplate[time.Time]

func ParseTimeTextTemplate(text string, funcMap templatex.TemplateFuncMap, parseTimeText func(text string) (time.Time, error)) (*TimeTextTemplate, error) {
	return ParseDataTextTemplate(text, funcMap, parseTimeText)
}

// 时间提取型模板
type TimeExtractorTemplate = DataExtractorTemplate[time.Time]

func ParseTimeExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (time.Time, error)) (*TimeExtractorTemplate, error) {
	return ParseDataExtractorTemplate(text, funcMap, concurrent, containerName, dataCast)
}

// 布尔模板
type BooleanTemplate string

func (tpl BooleanTemplate) Native() string {
	return string(tpl)
}

func (tpl BooleanTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (BooleanRender, error) {
	return ParseBooleanRenderTemplate(tpl.Native(), funcMap, containerName, concurrent)
}

func (tpl BooleanTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any) (bool, error) {
	return ParseAndRenderBooleanTemplate(tpl.Native(), funcMap, containerName, data)
}

// 布尔渲染器
type BooleanRender = DataRender[bool]

func ParseBooleanRenderTemplate(text string, funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (BooleanRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, parseBooleanText, "", false, nil)
}

func parseBooleanText(text string) (bool, error) {
	return strings.TrimSpace(text) == "true", nil
}

func ParseAndRenderBooleanTemplate(text string, funcMap templatex.TemplateFuncMap, containerName string, data any) (bool, error) {
	render, err := ParseBooleanRenderTemplate(text, funcMap, containerName, false)
	if err != nil {
		return false, err
	}
	return render.Render(data)
}

// 布尔文本型模板
type BooleanTextTemplate = DataTextTemplate[bool]

func ParseBooleanTextTemplate(text string, funcMap templatex.TemplateFuncMap, parseBooleanText func(text string) (bool, error)) (*BooleanTextTemplate, error) {
	return ParseDataTextTemplate(text, funcMap, parseBooleanText)
}

// 布尔提取型模板
type BooleanExtractorTemplate = DataExtractorTemplate[bool]

func ParseBooleanExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (bool, error)) (*BooleanExtractorTemplate, error) {
	return ParseDataExtractorTemplate(text, funcMap, concurrent, containerName, dataCast)
}

type TimeFormat string

func (format TimeFormat) Native() string {
	return string(format)
}

func (format TimeFormat) EnsureNonzero(zeroPlaceholder TimeFormat) TimeFormat {
	if format == "" {
		return zeroPlaceholder
	}
	return format
}

func (format TimeFormat) NonzeroOrDefault() TimeFormat {
	return format.EnsureNonzero(_time.DefaultTimeFormat)
}

func (format TimeFormat) ParseTime(text string) (time.Time, error) {
	return time.ParseInLocation(format.Native(), text, time.Local)
}

func (format TimeFormat) ParseTimeOrZero(text string) (time.Time, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return time.Time{}, nil
	}
	return format.ParseTime(text)
}

func (format TimeFormat) ParseRenderedTimeOrZero(text string) (time.Time, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return time.Time{}, nil
	}

	if text == templatex.TemplateValueNoValue {
		return time.Time{}, nil
	}

	return format.ParseTime(text)
}

type IntTemplate string

func (tpl IntTemplate) Native() string {
	return string(tpl)
}

func (tpl IntTemplate) ParseForRender(funcMap templatex.TemplateFuncMap, containerName string, concurrent bool) (IntRender, error) {
	return ParseIntRenderTemplate(tpl.Native(), funcMap, IntFormat("").ParseRenderedOrZero, containerName, concurrent)
}

func (tpl IntTemplate) ParseAndRender(funcMap templatex.TemplateFuncMap, containerName string, data any) (int, error) {
	render, err := ParseIntRenderTemplate(tpl.Native(), funcMap, IntFormat("").ParseRenderedOrZero, containerName, false)
	if err != nil {
		return 0, err
	}
	return render.Render(data)
}

type IntRender = DataRender[int]

func ParseIntRenderTemplate(text string, funcMap templatex.TemplateFuncMap, parseIntText func(text string) (int, error), containerName string, concurrent bool) (IntRender, error) {
	return ParseDataRenderTemplate(text, funcMap, TemplateTypeMagicPlain, parseIntText, containerName, concurrent, nil)
}

func ParseAndRenderIntTemplate(text string, funcMap templatex.TemplateFuncMap, parseIntText func(text string) (int, error), containerName string, data any) (int, error) {
	render, err := ParseIntRenderTemplate(text, funcMap, parseIntText, containerName, false)
	if err != nil {
		return 0, err
	}
	return render.Render(data)
}

type IntTextTemplate = DataTextTemplate[int]

func ParseIntTextTemplate(text string, funcMap templatex.TemplateFuncMap, parseIntText func(text string) (int, error)) (*IntTextTemplate, error) {
	return ParseDataTextTemplate(text, funcMap, parseIntText)
}

type IntExtractorTemplate = DataExtractorTemplate[int]

func ParseIntExtractorTemplate(text string, funcMap templatex.TemplateFuncMap, concurrent bool, containerName string, dataCast func(any) (int, error)) (*IntExtractorTemplate, error) {
	return ParseDataExtractorTemplate(text, funcMap, concurrent, containerName, dataCast)
}

type IntFormat = GenericJsonFormat[int]

type GenericJsonFormat[Data any] string

func (format GenericJsonFormat[Data]) Native() string {
	return string(format)
}

func (format GenericJsonFormat[Data]) Parse(text string) (data Data, err error) {
	err = json.Unmarshal([]byte(text), &data)
	return
}

func (format GenericJsonFormat[Data]) ParseRenderedOrZero(text string) (data Data, err error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	if text == templatex.TemplateValueNoValue {
		return
	}

	data, err = format.Parse(text)
	return
}

func GenericCastDatas[Datas ~[]DataPtr, DataPtr ~*Data, Data any](value any) (Datas, error) {
	if value == nil {
		return nil, nil
	}

	if records, ok := value.(Datas); ok {
		return records, nil
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() == reflect.Slice {
		records := make(Datas, 0, valueReflect.Len())
		for i := 0; i < valueReflect.Len(); i++ {
			item := valueReflect.Index(i).Interface()
			record, ok := item.(DataPtr)
			if !ok {
				if rec, ok := item.(Data); ok {
					copied := rec
					record = &copied
				} else {
					return nil, fmt.Errorf("extractor result element must be %s, got %T", stl.ReflectType[DataPtr](), item)
				}
			}
			records = append(records, record)
		}
		return records, nil
	}

	return nil, fmt.Errorf("extractor result must be %s or slice, got %T", stl.ReflectType[Datas](), value)
}
