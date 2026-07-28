package param

import (
	"fmt"
	"strings"
)

// RequiredString reads a non-empty string argument.
func RequiredString(args map[string]any, key string) (string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return "", errRequired(key)
	}
	s, ok := raw.(string)
	if !ok {
		return "", errInvalidType(key, "string")
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errRequired(key)
	}
	return s, nil
}

// OptionalString reads an optional string; missing or empty yields "".
func OptionalString(args map[string]any, key string) (string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return "", nil
	}
	s, ok := raw.(string)
	if !ok {
		return "", errInvalidType(key, "string")
	}
	return strings.TrimSpace(s), nil
}

// StringEnum reads a required string that must be one of allowed values.
func StringEnum(args map[string]any, key string, allowed []string) (string, error) {
	s, err := RequiredString(args, key)
	if err != nil {
		return "", err
	}
	for _, a := range allowed {
		if s == a {
			return s, nil
		}
	}
	return "", errMustBeOneOf(key, allowed)
}

// OptionalStringEnum reads an optional enum string; missing yields defaultValue.
func OptionalStringEnum(args map[string]any, key string, allowed []string, defaultValue string) (string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return defaultValue, nil
	}
	s, ok := raw.(string)
	if !ok {
		return "", errInvalidType(key, "string")
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return defaultValue, nil
	}
	for _, a := range allowed {
		if s == a {
			return s, nil
		}
	}
	return "", errMustBeOneOf(key, allowed)
}

// RequiredStringSlice reads a non-empty string array.
func RequiredStringSlice(args map[string]any, key string) ([]string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return nil, errRequired(key)
	}
	items, err := toStringSlice(raw, key)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errRequired(key)
	}
	return items, nil
}

// OptionalStringSlice reads an optional string array; missing yields nil.
func OptionalStringSlice(args map[string]any, key string) ([]string, error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return nil, nil
	}
	return toStringSlice(raw, key)
}

func toStringSlice(raw any, key string) ([]string, error) {
	switch v := raw.(type) {
	case []string:
		out := make([]string, 0, len(v))
		for _, s := range v {
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, s)
			}
		}
		return out, nil
	case []any:
		out := make([]string, 0, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("%s[%d] must be a string", key, i)
			}
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, s)
			}
		}
		return out, nil
	default:
		return nil, errInvalidType(key, "string array")
	}
}
