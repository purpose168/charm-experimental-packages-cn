package vt

import uv "github.com/charmbracelet/ultraviolet"

// CursorStyle 表示光标样式。
type CursorStyle int

// 光标样式。
const (
	CursorBlock CursorStyle = iota     // 块光标
	CursorUnderline                     // 下划线光标
	CursorBar                           // 竖线光标
)

// Cursor 表示终端中的光标。
type Cursor struct {
	Pen  uv.Style  // 光标样式
	Link uv.Link   // 光标链接

	uv.Position    // 光标位置

	Style  CursorStyle // 光标样式类型
	Steady bool        // 不闪烁
	Hidden bool        // 隐藏
}
