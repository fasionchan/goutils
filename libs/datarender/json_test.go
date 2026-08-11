package datarender

import (
	"testing"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/types"
	"github.com/stretchr/testify/assert"
)

func TestSmartJsonMapTemplate_ParseAndRender(t *testing.T) {
	tpl := SmartJsonMapTemplate(`{"time": {{ today | jsonEncode }}}`)
	render, err := tpl.ParseForRender(templatex.TemplateFuncs, false)
	if err != nil {
		t.Errorf("case failed: %s", err)
		return
	}

	result, err := render.Render(nil)
	if err != nil {
		t.Errorf("case failed: %s", err)
		return
	}

	assert.Equal(t, result, types.SmartJsonMap{
		"time": MustJsonUnmarshal[string](MustJsonMarshalRaw(baseutils.Today())),
	})
}
