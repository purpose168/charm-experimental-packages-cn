package pony

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// Button 表示一个可点击的按钮元素。
type Button struct {
	BaseElement
	text        string
	style       uv.Style
	hoverStyle  uv.Style
	activeStyle uv.Style
	border      string
	padding     int
	width       SizeConstraint
	height      SizeConstraint
}

var _ Element = (*Button)(nil)

// NewButton 创建一个新的按钮元素。
func NewButton(text string) *Button {
	return &Button{
		text:    text,
		border:  BorderRounded,
		padding: 1,
	}
}

// Style 设置按钮样式并返回按钮以支持链式调用。
func (b *Button) Style(style uv.Style) *Button {
	b.style = style
	return b
}

// HoverStyle 设置悬停样式并返回按钮以支持链式调用。
func (b *Button) HoverStyle(style uv.Style) *Button {
	b.hoverStyle = style
	return b
}

// ActiveStyle 设置激活（按下）样式并返回按钮以支持链式调用。
func (b *Button) ActiveStyle(style uv.Style) *Button {
	b.activeStyle = style
	return b
}

// Border 设置边框类型并返回按钮以支持链式调用。
func (b *Button) Border(border string) *Button {
	b.border = border
	return b
}

// Padding 设置内边距并返回按钮以支持链式调用。
func (b *Button) Padding(padding int) *Button {
	b.padding = padding
	return b
}

// Width 设置宽度约束并返回按钮以支持链式调用。
func (b *Button) Width(width SizeConstraint) *Button {
	b.width = width
	return b
}

// Height 设置高度约束并返回按钮以支持链式调用。
func (b *Button) Height(height SizeConstraint) *Button {
	b.height = height
	return b
}

// Draw 将按钮渲染到屏幕上。
func (b *Button) Draw(scr uv.Screen, area uv.Rectangle) {
	b.SetBounds(area)

	// 创建文本元素
	textElem := NewText(b.text).Alignment(AlignmentCenter)
	if !b.style.IsZero() {
		// 应用样式到文本内容
		if b.style.Fg != nil {
			textElem = textElem.ForegroundColor(b.style.Fg)
		}
		if b.style.Attrs&uv.AttrBold != 0 {
			textElem = textElem.Bold()
		}
		if b.style.Attrs&uv.AttrItalic != 0 {
			textElem = textElem.Italic()
		}
	}

	// 用带边框和内边距的盒子包裹
	box := NewBox(textElem).
		Border(b.border).
		Padding(b.padding)

	if !b.style.IsZero() && b.style.Fg != nil {
		box = box.BorderColor(b.style.Fg)
	}

	box.Draw(scr, area)
}

// Layout 计算按钮大小。
func (b *Button) Layout(constraints Constraints) Size {
	// 创建用于计算大小的文本元素
	textElem := NewText(b.text)
	textSize := textElem.Layout(Unbounded())

	// 添加内边距和边框
	borderSize := 2
	if b.border == BorderNone || b.border == BorderHidden {
		borderSize = 0
	}

	paddingSize := b.padding * 2

	width := textSize.Width + borderSize + paddingSize
	height := textSize.Height + borderSize + paddingSize

	result := Size{Width: width, Height: height}

	// 如果指定了宽度约束，则应用
	if !b.width.IsAuto() {
		result.Width = b.width.Apply(constraints.MaxWidth, result.Width)
	}

	// 如果指定了高度约束，则应用
	if !b.height.IsAuto() {
		result.Height = b.height.Apply(constraints.MaxHeight, result.Height)
	}

	return constraints.Constrain(result)
}

// Children 为按钮返回 nil。
func (b *Button) Children() []Element {
	return nil
}

// NewButtonFromProps 从属性创建按钮（用于解析器）。
func NewButtonFromProps(props Props, children []Element) Element {
	text := props.Get("text")
	if text == "" && len(children) > 0 {
		if t, ok := children[0].(*Text); ok {
			text = t.Content()
		}
	}

	btn := NewButton(text)

	// 解析按钮文本/边框的前景色
	if fgColor := props.Get("foreground-color"); fgColor != "" {
		if c, err := parseColor(fgColor); err == nil {
			style := uv.Style{Fg: c}
			if props.Get("font-weight") == FontWeightBold {
				style.Attrs |= uv.AttrBold
			}
			btn = btn.Style(style)
		}
	}

	if border := props.Get("border"); border != "" {
		btn = btn.Border(border)
	}

	if padding := parseIntAttr(props, "padding", 0); padding > 0 {
		btn = btn.Padding(padding)
	}

	if width := props.Get("width"); width != "" {
		btn = btn.Width(parseSizeConstraint(width))
	}

	if height := props.Get("height"); height != "" {
		btn = btn.Height(parseSizeConstraint(height))
	}

	return btn
}

func init() {
	Register("button", NewButtonFromProps)
}
