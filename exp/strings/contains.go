package strings

import (
	"strings"
)

// ContainsAnyOf 返回给定字符串是否包含以下任何一个字符串。
func ContainsAnyOf(str string, args ...string) bool {
	for _, arg := range args {
		if strings.Contains(str, arg) {
			return true
		}
	}
	return false
}
