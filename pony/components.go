package pony

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// Badge 表示一个徽章元素（如 "NEW"、"BETA"、状态指示器等）。
type Badge struct {
	BaseElement
	text  string
	color color.Color
}

var _ Element = (*Badge)(nil)

// NewBadge 创建一个徽章组件。
func NewBadge(props Props, children []Element) Element {
	text := props.Get("text")
	if text == "" && len(children) > 0 {
		// 使用第一个子元素作为文本
		if t, ok := children[0].(*Text); ok {
			text = t.Content()
		}
	}

	var fgColor color.Color
	if colorStr := props.Get("foreground-color"); colorStr != "" {
		if c, err := parseColor(colorStr); err == nil {
			fgColor = c
		}
	}

	return &Badge{
		text:  text,
		color: fgColor,
	}
}

// Draw 渲染徽章。
func (b *Badge) Draw(scr uv.Screen, area uv.Rectangle) {
	b.SetBounds(area)

	if b.text == "" {
		return
	}

	// 徽章渲染为带样式的 [文本]
	content := "[" + b.text + "]"
	if b.color != nil {
		style := uv.Style{Fg: b.color}
		content = style.Styled(content)
	}

	styled := uv.NewStyledString(content)
	styled.Draw(scr, area)
}

// Layout 计算徽章大小。
func (b *Badge) Layout(constraints Constraints) Size {
	width := len(b.text) + 2 // 文本 + 括号
	return constraints.Constrain(Size{Width: width, Height: 1})
}

// Children 返回 nil。
func (b *Badge) Children() []Element {
	return nil
}

// ProgressView 表示一个进度条元素。
type ProgressView struct {
	BaseElement
	value int
	max   int
	width SizeConstraint
	color color.Color
	char  string
}

var _ Element = (*ProgressView)(nil)

// NewProgressView 创建一个进度条组件。
func NewProgressView(props Props, _ []Element) Element {
	value := parseIntAttr(props, "value", 0)
	maxValue := parseIntAttr(props, "max", 100)
	width := parseSizeConstraint(props.Get("width"))
	char := props.GetOr("char", "█")

	var fgColor color.Color
	if colorStr := props.Get("foreground-color"); colorStr != "" {
		if c, err := parseColor(colorStr); err == nil {
			fgColor = c
		}
	}

	return &ProgressView{
		value: value,
		max:   maxValue,
		width: width,
		color: fgColor,
		char:  char,
	}
}

// Draw 渲染进度条。
func (p *ProgressView) Draw(scr uv.Screen, area uv.Rectangle) {
	p.SetBounds(area)

	if area.Dx() == 0 {
		return
	}

	// 计算填充部分
	filled := 0
	if p.max > 0 {
		filled = min((area.Dx()*p.value)/p.max, area.Dx())
	}

	// 创建填充部分的单元格
	filledCell := uv.NewCell(scr.WidthMethod(), p.char)
	if filledCell != nil && p.color != nil {
		filledCell.Style = uv.Style{Fg: p.color}
	}

	// 创建空部分的单元格
	emptyCell := uv.NewCell(scr.WidthMethod(), "░")

	// 绘制进度条
	for x := 0; x < area.Dx(); x++ {
		if x < filled {
			scr.SetCell(area.Min.X+x, area.Min.Y, filledCell)
		} else {
			scr.SetCell(area.Min.X+x, area.Min.Y, emptyCell)
		}
	}
}

// Layout 计算进度条大小。
func (p *ProgressView) Layout(constraints Constraints) Size {
	// 未指定时的默认宽度
	width := 20

	// 如果指定了宽度约束，则应用
	if !p.width.IsAuto() {
		// 对于固定宽度，直接使用约束值
		width = p.width.Apply(constraints.MaxWidth, width)
	} else {
		// 对于自动宽度，使用可用宽度
		width = constraints.MaxWidth
	}

	return Size{Width: width, Height: 1}
}

// Children 返回 nil。
func (p *ProgressView) Children() []Element {
	return nil
}

// init 注册内置的自定义组件。
func init() {
	Register("badge", NewBadge)
	Register("progressview", NewProgressView)
}
