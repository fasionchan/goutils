package param

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TestArgsEmpty(t *testing.T) {
	req := mcp.CallToolRequest{}
	args := Args(req)
	require.NotNil(t, args)
	require.Empty(t, args)
}

func TestRequiredString(t *testing.T) {
	_, err := RequiredString(map[string]any{}, "tabId")
	require.Error(t, err)
	require.Contains(t, err.Error(), "tabId is required")

	_, err = RequiredString(map[string]any{"tabId": "  "}, "tabId")
	require.Error(t, err)

	got, err := RequiredString(map[string]any{"tabId": "abc"}, "tabId")
	require.NoError(t, err)
	require.Equal(t, "abc", got)
}

func TestOptionalString(t *testing.T) {
	got, err := OptionalString(map[string]any{}, "url")
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestStringEnum(t *testing.T) {
	allowed := []string{"ref", "css-selector", "xpath"}
	_, err := StringEnum(map[string]any{"targetType": "foo"}, "targetType", allowed)
	require.Error(t, err)
	require.Contains(t, err.Error(), "targetType must be one of")

	got, err := StringEnum(map[string]any{"targetType": "ref"}, "targetType", allowed)
	require.NoError(t, err)
	require.Equal(t, "ref", got)
}

func TestOptionalBoolDefault(t *testing.T) {
	got, err := OptionalBool(map[string]any{}, "doubleClick", false)
	require.NoError(t, err)
	require.False(t, got)

	got, err = OptionalBool(map[string]any{"doubleClick": true}, "doubleClick", false)
	require.NoError(t, err)
	require.True(t, got)
}

func TestOptionalInt(t *testing.T) {
	got, err := OptionalInt(map[string]any{}, "count", 1)
	require.NoError(t, err)
	require.Equal(t, 1, got)

	got, err = OptionalInt(map[string]any{"count": float64(2)}, "count", 1)
	require.NoError(t, err)
	require.Equal(t, 2, got)
}

func TestRequiredStringSlice(t *testing.T) {
	_, err := RequiredStringSlice(map[string]any{}, "paths")
	require.Error(t, err)
	require.Contains(t, err.Error(), "paths")

	_, err = RequiredStringSlice(map[string]any{"paths": []any{}}, "paths")
	require.Error(t, err)

	got, err := RequiredStringSlice(map[string]any{"values": []any{"a", "b"}}, "values")
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b"}, got)
}
