package datarender

import (
	"testing"

	"github.com/fasionchan/goutils/std/_testing"
	"github.com/stretchr/testify/assert"
)

type TemplateMagicTestCase struct {
	name string
	text string
	plainMagic string
	expected *TemplateMagic
}

func (testCase *TemplateMagicTestCase) GetName() string {
	return testCase.name
}

func (testCase *TemplateMagicTestCase) Run(t *testing.T) {
	magic, err := ParseTemplateMagic(testCase.text, testCase.plainMagic)
	if err != nil {
		t.Errorf("case failed: %s", err)
		return
	}
	assert.Equal(t, testCase.expected, magic)
}

func TestParseTemplateMagic(t *testing.T) {
	_testing.TypedRunNamedTestCases(t, []*TemplateMagicTestCase{
		{
			name: "plain",
			text: TemplateTypeMagicPlain+"abc",
			plainMagic: TemplateTypeMagicPlain,
			expected: &TemplateMagic{
				Type: TemplateTypePlain,
				Attributes: make(map[string]string),
			},
		},
		{
			name: "text",
			text: TemplateTypeMagicText + `{[( format=json)]} */}}\n{"key": "value}`,
			expected: &TemplateMagic{
				Type: TemplateTypeText,
				Attributes: map[string]string{
					"format": "json",
				},
			},
		},
		{
			name: "extractor",
			text: TemplateTypeMagicExtractor + `{[( format=json )]} */}}\n{"key": "value}`,
			expected: &TemplateMagic{
				Type: TemplateTypeExtractor,
				Attributes: map[string]string{
					"format": "json",
				},
			},
		},
		{
			name: "lua",
			text: TemplateTypeMagicLua + `{[( format=json )]}\n{"key": "value}`,
			expected: &TemplateMagic{
				Type: TemplateTypeLua,
				Attributes: map[string]string{
					"format": "json",
				},
			},
		},
	})
}