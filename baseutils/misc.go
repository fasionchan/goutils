package baseutils

import "strings"

func ParseBoolDefaultFalse(s string) bool {
	return ParseBool(s, false)
}

func ParseBoolDefaultTrue(s string) bool {
	return ParseBool(s, true)
}

func ParseBool(s string, defaultValue bool) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true":
		return true
	case "false":
		return false
	case "":
		return defaultValue
	default:
		return false
	}
}