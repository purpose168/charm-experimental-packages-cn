package cellbuf

import (
	"github.com/charmbracelet/colorprofile"
)

// ConvertStyle 根据给定的颜色配置文件转换样式。
// 参数：
// - s: 要转换的样式
// - p: 目标颜色配置文件
// 返回值：
// - 转换后的样式
func ConvertStyle(s Style, p colorprofile.Profile) Style {
	// 根据不同的颜色配置文件类型进行处理
	switch p { //nolint:exhaustive
	case colorprofile.TrueColor:
		// 真彩色模式不需要转换，直接返回原样式
		return s
	case colorprofile.Ascii:
		// ASCII模式下移除所有颜色和下划线
		s.Fg = nil // 前景色设为nil
		s.Bg = nil // 背景色设为nil
		s.Ul = nil // 下划线设为nil
	case colorprofile.NoTTY:
		// 非终端模式返回空样式
		return Style{}
	}

	// 对于其他颜色配置文件，转换各颜色属性
	if s.Fg != nil {
		s.Fg = p.Convert(s.Fg) // 转换前景色
	}
	if s.Bg != nil {
		s.Bg = p.Convert(s.Bg) // 转换背景色
	}
	if s.Ul != nil {
		s.Ul = p.Convert(s.Ul) // 转换下划线颜色
	}

	return s
}
