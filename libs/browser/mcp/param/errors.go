package param

import (
	"fmt"
	"strings"
)

func errRequired(key string) error {
	return fmt.Errorf("%s is required", key)
}

func errMustBeOneOf(key string, allowed []string) error {
	return fmt.Errorf("%s must be one of: %s", key, strings.Join(allowed, ", "))
}

func errInvalidType(key, want string) error {
	return fmt.Errorf("%s must be a %s", key, want)
}
