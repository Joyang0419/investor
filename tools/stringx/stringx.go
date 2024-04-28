package stringx

import (
	"strings"
)

func CheckEmptyStr(s string) bool {
	return len(s) == 0
}

func CheckNotEmptyStr(s string) bool {
	return len(s) != 0
}

func CheckHasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
