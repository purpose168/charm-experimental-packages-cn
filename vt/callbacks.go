package vt

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// Callbacks 表示终端的一组回调函数。
type Callbacks struct {
	// Bell 回调。当设置时，此函数在接收到响铃字符时被调用。
	Bell func()

	// Title 回调。当设置时，此函数在终端标题更改时被调用。
	Title func(string)

	// IconName 回调。当设置时，此函数在终端图标名称更改时被调用。
	IconName func(string)

	// AltScreen 回调。当设置时，此函数在交替屏幕被激活或停用时被调用。
	AltScreen func(bool)

	// CursorPosition 回调。当设置时，此函数在光标位置更改时被调用。
	CursorPosition func(old, new uv.Position) //nolint:predeclared,revive

	// CursorVisibility 回调。当设置时，此函数在光标可见性更改时被调用。
	CursorVisibility func(visible bool)

	// CursorStyle 回调。当设置时，此函数在光标样式更改时被调用。
	CursorStyle func(style CursorStyle, blink bool)

	// CursorColor 回调。当设置时，此函数在光标颜色更改时被调用。nil 表示默认终端颜色。
	CursorColor func(color color.Color)

	// BackgroundColor 回调。当设置时，此函数在背景颜色更改时被调用。nil 表示默认终端颜色。
	BackgroundColor func(color color.Color)

	// ForegroundColor 回调。当设置时，此函数在前景颜色更改时被调用。nil 表示默认终端颜色。
	ForegroundColor func(color color.Color)

	// WorkingDirectory 回调。当设置时，此函数在当前工作目录更改时被调用。
	WorkingDirectory func(string)

	// EnableMode 回调。当设置时，此函数在模式启用时被调用。
	EnableMode func(mode ansi.Mode)

	// DisableMode 回调。当设置时，此函数在模式禁用时被调用。
	DisableMode func(mode ansi.Mode)
}
