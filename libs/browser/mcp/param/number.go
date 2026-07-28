package param

import (
	"strconv"
	"strings"
)

// OptionalBool reads an optional bool with a default when missing.
func OptionalBool(args map[string]any, key string, defaultValue bool) (bool, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return defaultValue, nil
	}
	switch v := raw.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, errInvalidType(key, "bool")
		}
		return b, nil
	case float64:
		return v != 0, nil
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	default:
		return false, errInvalidType(key, "bool")
	}
}

// OptionalInt reads an optional int with a default when missing.
func OptionalInt(args map[string]any, key string, defaultValue int) (int, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return defaultValue, nil
	}
	switch v := raw.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		i, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, errInvalidType(key, "number")
		}
		return i, nil
	default:
		return 0, errInvalidType(key, "number")
	}
}

// Has reports whether key is present and non-nil.
func Has(args map[string]any, key string) bool {
	raw, ok := args[key]
	return ok && raw != nil
}
