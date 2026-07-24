package templatex

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSmartTemplateWithTimeFormat(t *testing.T) {
	// Test the current SmartTemplate implementation
	// After fix: should use text/template and preserve '+'
	tplText := `{{ .Now.Format "20060102 15:04:05Z0700" }}`

	// Mock data with time
	data := map[string]interface{}{
		"Now": time.Date(2026, 3, 2, 15, 5, 48, 0, time.FixedZone("CST", 8*3600)),
	}

	smartTpl, err := ParseSmartTemplateFromText("timeformat", tplText, nil)
	assert.NoError(t, err)

	result, err := smartTpl.RenderToString(data)
	assert.NoError(t, err)

	t.Logf("SmartTemplate result: %s", result)

	// After switching to text/template, should NOT contain HTML entities
	assert.True(t, !bytes.Contains([]byte(result), []byte("&#43;")),
		"SmartTemplate should NOT escape '+' after switching to text/template")
	assert.True(t, bytes.Contains([]byte(result), []byte("+0800")),
		"SmartTemplate should preserve raw '+0800'")
}
