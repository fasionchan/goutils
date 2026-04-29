package types

import (
	"strconv"
	"strings"
)

func FormatFloat(value float64, precision int, trimZero bool) string {
	s := strconv.FormatFloat(value, 'f', precision, 64)
	if trimZero {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
