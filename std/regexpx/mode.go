package regexpx

import "regexp"

var (
	SpaceMode      = regexp.MustCompile(`\s`)
	SpaceOrEscMode = regexp.MustCompile(`\s|\\n|\\t|\\r|\\f|\\v`)
)

func PurgeStringSpace(s string) string {
	return SpaceMode.ReplaceAllString(s, "")
}

// PurgeStringSpaceOrEsc 使用正则清空所有空间（含转义字符）
func PurgeStringSpaceOrEsc(s string) string {
	return SpaceOrEscMode.ReplaceAllString(s, "")
}
