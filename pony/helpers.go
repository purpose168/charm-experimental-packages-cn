package pony

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// StyleBuilder 提供了一个流畅的 API 用于构建样式。
type StyleBuilder struct {
	style uv.Style
}

// NewStyle 创建一个新的样式构建器。
func NewStyle() *StyleBuilder {
	return &StyleBuilder{}
}

// Fg 设置前景色。
func (sb *StyleBuilder) Fg(c color.Color) *StyleBuilder {
	sb.style.Fg = c
	return sb
}

// Bg 设置背景色。
func (sb *StyleBuilder) Bg(c color.Color) *StyleBuilder {
	sb.style.Bg = c
	return sb
}

// UnderlineColor 设置下划线颜色。
func (sb *StyleBuilder) UnderlineColor(c color.Color) *StyleBuilder {
	sb.style.UnderlineColor = c
	return sb
}

// Bold 使文本加粗。
func (sb *StyleBuilder) Bold() *StyleBuilder {
	sb.style.Attrs |= uv.AttrBold
	return sb
}

// Faint 使文本变淡/变暗。
func (sb *StyleBuilder) Faint() *StyleBuilder {
	sb.style.Attrs |= uv.AttrFaint
	return sb
}

// Italic 使文本变为斜体。
func (sb *StyleBuilder) Italic() *StyleBuilder {
	sb.style.Attrs |= uv.AttrItalic
	return sb
}

// Underline 设置单下划线。
func (sb *StyleBuilder) Underline() *StyleBuilder {
	sb.style.Underline = uv.UnderlineSingle
	return sb
}

// UnderlineStyle 设置下划线样式。
func (sb *StyleBuilder) UnderlineStyle(style uv.Underline) *StyleBuilder {
	sb.style.Underline = style
	return sb
}

// Blink 使文本闪烁。
func (sb *StyleBuilder) Blink() *StyleBuilder {
	sb.style.Attrs |= uv.AttrBlink
	return sb
}

// Reverse 反转前景色和背景色。
func (sb *StyleBuilder) Reverse() *StyleBuilder {
	sb.style.Attrs |= uv.AttrReverse
	return sb
}

// Strikethrough 添加删除线。
func (sb *StyleBuilder) Strikethrough() *StyleBuilder {
	sb.style.Attrs |= uv.AttrStrikethrough
	return sb
}

// Build 返回构建好的样式。
func (sb *StyleBuilder) Build() uv.Style {
	return sb.style
}

// 颜色辅助函数

// Hex 从十六进制字符串创建颜色。
// 如果无效会 panic - 使用 HexSafe 进行错误处理。
func Hex(s string) color.Color {
	c, err := parseHexColor(s)
	if err != nil {
		panic(err)
	}
	return c
}

// HexSafe 从十六进制字符串创建颜色并进行错误处理。
func HexSafe(s string) (color.Color, error) {
	return parseHexColor(s)
}

// RGB 从 RGB 值创建颜色。
func RGB(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// 通用布局辅助函数

// Panel 创建带有边框和内边距的盒子。
func Panel(child Element, border string, padding int) *Box {
	return NewBox(child).
		Border(border).
		Padding(padding)
}

// PanelWithMargin 创建带有边框、内边距和外边距的盒子。
func PanelWithMargin(child Element, border string, padding, margin int) *Box {
	return NewBox(child).
		Border(border).
		Padding(padding).
		Margin(margin)
}

// Card 创建带有标题和内容的卡片。
func Card(title string, titleColor, borderColor color.Color, children ...Element) Element {
	titleText := NewText(title)
	if titleColor != nil {
		titleText = titleText.ForegroundColor(titleColor).Bold()
	}

	box := NewBox(
		NewVStack(
			titleText,
			NewDivider(),
			NewVStack(children...),
		),
	).Border("rounded").Padding(1)

	if borderColor != nil {
		box = box.BorderColor(borderColor)
	}

	return box
}

// Section 创建带有标题和内容的部分。
func Section(header string, headerColor color.Color, children ...Element) Element {
	headerText := NewText(header)
	if headerColor != nil {
		headerText = headerText.ForegroundColor(headerColor).Bold()
	}

	items := []Element{headerText}
	items = append(items, children...)
	return NewVStack(items...)
}

// Separated 在每个子元素之间添加分隔线。
func Separated(children ...Element) Element {
	if len(children) == 0 {
		return NewVStack()
	}

	items := make([]Element, 0, len(children)*2-1)
	for i, child := range children {
		items = append(items, child)
		if i < len(children)-1 {
			items = append(items, NewDivider())
		}
	}

	return NewVStack(items...)
}

// Overlay 创建一个 Z 轴堆叠，子元素层叠在一起。
func Overlay(children ...Element) Element {
	return NewZStack(children...)
}

// FlexGrow 创建带有指定增长值的弹性布局包装器。
func FlexGrow(child Element, grow int) *Flex {
	return NewFlex(child).Grow(grow)
}

// Position 创建绝对定位的元素。
func Position(child Element, x, y int) *Positioned {
	return NewPositioned(child, x, y)
}

// PositionRight 创建相对于右边缘定位的元素。
func PositionRight(child Element, right, y int) *Positioned {
	return NewPositioned(child, 0, y).Right(right)
}

// PositionBottom 创建相对于底部边缘定位的元素。
func PositionBottom(child Element, x, bottom int) *Positioned {
	return NewPositioned(child, x, 0).Bottom(bottom)
}

// PositionCorner 创建定位在角落的元素。
func PositionCorner(child Element, right, bottom int) *Positioned {
	return NewPositioned(child, 0, 0).Right(right).Bottom(bottom)
}
