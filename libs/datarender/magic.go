package datarender

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/fasionchan/goutils/stl"
)

const (
	TemplateTypePlain = iota
	TemplateTypeText
	TemplateTypeExtractor
	TemplateTypeLua

	TemplateTypeMagicPlain         = "\t\r\n"
	TemplateTypeMagicPlainDisabled = "__TemplateTypePlainDisabled__"        // 禁用纯文本模板
	TemplateTypeMagicText          = "{{/* TemplateTypeMagicText"      // 数据文本型模板魔法字符串
	TemplateTypeMagicExtractor     = "{{/* TemplateTypeMagicExtractor" // 数据提取型模板魔法字符串
	TemplateTypeMagicLua           = "-- TemplateTypeMagicLua"       		// Lua 脚本型模板魔法字符串
)

var TemplateMagicAttributePattern = regexp.MustCompile(`(?s)\{\[\(\s*(.*?)\s*\)\]\}`)

type TemplateMagic struct {
	Type int
	Attributes map[string]string
}

func ParseTemplateMagic(text string, plainMagic string) (*TemplateMagic, error) {
	var attributes map[string]string

	matches := TemplateMagicAttributePattern.FindStringSubmatch(strings.SplitN(text, "\n", 2)[0])
	if len(matches) == 2 {
		query, err := url.ParseQuery(matches[1])
		if err != nil {
			return nil, err
		}

		attributes = stl.MapMap[map[string]string](query, func(key string, values []string, _ url.Values) (string, string) {
			return key, stl.LastOneOrZero(values)
		})
	} else {
		attributes = make(map[string]string)
	}

	switch {
	case plainMagic != "" && strings.HasPrefix(text, plainMagic):
		return &TemplateMagic{
			Type: TemplateTypePlain,
			Attributes: attributes,
		}, nil
	case strings.HasPrefix(text, TemplateTypeMagicLua):
		return &TemplateMagic{
			Type: TemplateTypeLua,
			Attributes: attributes,
		}, nil
	case strings.HasPrefix(text, TemplateTypeMagicExtractor):
		return &TemplateMagic{
			Type: TemplateTypeExtractor,
			Attributes: attributes,
		}, nil
	default:
		return &TemplateMagic{
			Type: TemplateTypeText,
			Attributes: attributes,
		}, nil
	}
}

func (magic *TemplateMagic) GetFormat() string {
	if magic == nil {
		return ""
	}

	return magic.Attributes["format"]
}