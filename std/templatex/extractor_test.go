package templatex

import (
	"testing"

	"github.com/fasionchan/goutils/baseutils"
)

var Data = baseutils.JsonMap{
	"a": baseutils.JsonMap{
		"b": baseutils.JsonMap{
			"c": 1,
		},
		"c": 2,
	},
	"b": 3,
}
var PathTemplateDataExtractorCases = []struct {
	data     any
	path     string
	expected any
}{
	{
		data:     1,
		path:     "",
		expected: 1,
	},
	{
		data:     1,
		path:     ".",
		expected: 1,
	},
	{
		data:     Data,
		path:     "b",
		expected: 3,
	},
	{
		data:     Data,
		path:     ".b",
		expected: 3,
	},
	{
		data:     Data,
		path:     "a.c",
		expected: 2,
	},
	{
		data:     Data,
		path:     ".a.c",
		expected: 2,
	},
	{
		data:     Data,
		path:     "a.b.c",
		expected: 1,
	},
	{
		data:     Data,
		path:     ".a.b.c",
		expected: 1,
	},
}

func TestExtracDataByPathTempalte(t *testing.T) {
	for _, item := range PathTemplateDataExtractorCases {
		result, err := ExtractDataByPathTempalte(item.data, item.path, nil)
		if err != nil {
			t.Error(err)
			continue
		}

		if result != item.expected {
			t.Fatalf("`%v` expected but got `%v`", item.expected, result)
		}
	}
}