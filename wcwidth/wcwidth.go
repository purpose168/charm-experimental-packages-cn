// Package wcwidth 提供用于计算字符宽度的工具。
package wcwidth

import (
	"github.com/mattn/go-runewidth"
)

// RuneWidth 返回字符的固定宽度。
//
// 已弃用：这现在是 go-runewidth 的包装器。直接使用 go-runewidth。
func RuneWidth(r rune) int {
	return runewidth.RuneWidth(r)
}

// StringWidth 返回字符串的固定宽度。
//
// 已弃用：这现在是 go-runewidth 的包装器。直接使用 go-runewidth。
func StringWidth(s string) (n int) {
	return runewidth.StringWidth(s)
}
