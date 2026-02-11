// colors 包提供颜色定义和工具函数
package colors

import "github.com/purpose168/lipgloss-cn"

//nolint:revive
var (
	// WhiteBright 亮白色 - 在亮色和暗色模式下都保持一致的白色
	WhiteBright = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}

	// Normal 普通黑色 - 亮色模式下为深黑色，暗色模式下为亮灰色
	Normal = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#dddddd"}
	// NormalDim 普通暗淡色 - 普通黑色的暗淡版本
	NormalDim = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}

	// Gray 灰色 - 中性灰色
	Gray = lipgloss.AdaptiveColor{Light: "#909090", Dark: "#626262"}
	// GrayMid 中灰色 - 介于亮灰和暗灰之间的灰色
	GrayMid = lipgloss.AdaptiveColor{Light: "#B2B2B2", Dark: "#4A4A4A"}
	// GrayDark 深灰色 - 较暗的灰色
	GrayDark = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#222222"}
	// GrayBright 亮灰色 - 较亮的灰色
	GrayBright = lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"}
	// GrayBrightDim 亮灰色暗淡版 - 亮灰色的暗淡版本
	GrayBrightDim = lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"}

	// Indigo 靛蓝色 - 鲜明的靛蓝色
	Indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	// IndigoDim 靛蓝色暗淡版 - 靛蓝的暗淡版本
	IndigoDim = lipgloss.AdaptiveColor{Light: "#9498FF", Dark: "#494690"}
	// IndigoSubtle 柔和靛蓝色 - 较柔和的靛蓝色
	IndigoSubtle = lipgloss.AdaptiveColor{Light: "#7D79F6", Dark: "#514DC1"}
	// IndigoSubtleDim 柔和靛蓝色暗淡版 - 柔和靛蓝的暗淡版本
	IndigoSubtleDim = lipgloss.AdaptiveColor{Light: "#BBBDFF", Dark: "#383584"}

	// YellowGreen 黄绿色 - 鲜明的黄绿色
	YellowGreen = lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#ECFD65"}
	// YellowGreenDull 黄绿色暗淡版 - 黄绿色的暗淡版本
	YellowGreenDull = lipgloss.AdaptiveColor{Light: "#6BCB94", Dark: "#9BA92F"}

	// Fuchsia 紫红色 - 鲜明的紫红色
	Fuschia = lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}
	// FuchsiaDim 紫红色暗淡版 - 紫红色的暗淡版本
	FuchsiaDim = lipgloss.AdaptiveColor{Light: "#F1A8FF", Dark: "#99519E"}
	// FuchsiaDull 暗淡紫红色 - 暗淡的紫红色
	FuchsiaDull = lipgloss.AdaptiveColor{Dark: "#AD58B4", Light: "#F793FF"}
	// FuchsiaDullDim 暗淡紫红色暗淡版 - 暗淡紫红色的更暗淡版本
	FuchsiaDullDim = lipgloss.AdaptiveColor{Light: "#F6C9FF", Dark: "#6B3A6F"}

	// Green 绿色 - 鲜明的绿色
	Green = lipgloss.Color("#04B575")
	// GreenDim 绿色暗淡版 - 绿色的暗淡版本
	GreenDim = lipgloss.AdaptiveColor{Light: "#72D2B0", Dark: "#0B5137"}

	// Red 红色 - 鲜明的红色
	Red = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	// RedDull 红色暗淡版 - 红色的暗淡版本
	RedDull = lipgloss.AdaptiveColor{Light: "#FF6F91", Dark: "#C74665"}
)
