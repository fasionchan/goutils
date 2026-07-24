package templatex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataContainerRender(t *testing.T) {
	text := `
		{{ setString .Data.String }}
		{{ if .Data.Bool  }}
			{{ setNumber .Data.Number }}
		{{ end }}
	`

	dataExtractor, err := NewTemplateDataExtractorCustom(text, NewTemplateFuncMap(), true, "tpl", "String", "Number")
	if err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"Data": map[string]any{
			"String": "字符串",
			"Number": 1,
			"Bool":   true,
		},
	}

	mapping, err := dataExtractor.ExtractMany(data)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, "字符串", mapping["String"])
	assert.Equal(t, 1, mapping["Number"])
}