package stringx

import (
	"strings"
)

func IsEmptyStr(s string) bool {
	return len(s) == 0
}

func IsNotEmptyStr(s string) bool {
	return len(s) != 0
}

func IsHasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}
