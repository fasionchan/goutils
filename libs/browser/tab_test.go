package browser

import (
	"net/url"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/stretchr/testify/require"
)

func TestTabGetHtmlsRequest(t *testing.T) {
	var request TabGetHtmlsParams

	formDecoder := form.NewDecoder()
	formDecoder.SetTagName("query")
	formDecoder.SetMode(form.ModeExplicit)

	if err := formDecoder.Decode(&request, url.Values{
		"target":     {"target"},
		"targetType": {"targetType"},
	}); err != nil {
		t.Fatalf("Failed to decode request: %v", err)
	}

	require.Equal(t, request.Target, "target")
	require.Equal(t, request.TargetType, "targetType")
}
