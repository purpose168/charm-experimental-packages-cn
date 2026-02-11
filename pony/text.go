package pony

import (
	"image/color"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// Text 表示一个文本元素。
type Text struct {
	BaseElement
	content   string
	style     uv.Style
	wrap      bool
	alignment string // leading, center, trailing
}

var _ Element = (*Text)(nil)

// NewText 创建一个新的文本元素。
func NewText(content string) *Text {
	return &Text{content: content}
}

// Bold 将文本设置为粗体并返回文本以支持链式调用。
func (t *Text) Bold() *Text {
	t.style.Attrs |= uv.AttrBold
	return t
}

// Italic 将文本设置为斜体并返回文本以支持链式调用。
func (t *Text) Italic() *Text {
	t.style.Attrs |= uv.AttrItalic
	return t
}

// Underline 为文本添加下划线并返回文本以支持链式调用。
func (t *Text) Underline() *Text {
	t.style.Underline = uv.UnderlineSingle
	return t
}

// Strikethrough 为文本添加删除线并返回文本以支持链式调用。
func (t *Text) Strikethrough() *Text {
	t.style.Attrs |= uv.AttrStrikethrough
	return t
}

// Faint 将文本设置为淡色/暗淡并返回文本以支持链式调用。
func (t *Text) Faint() *Text {
	t.style.Attrs |= uv.AttrFaint
	return t
}

// ForegroundColor 设置文本前景色并返回文本以支持链式调用。
func (t *Text) ForegroundColor(c color.Color) *Text {
	t.style.Fg = c
	return t
}

// BackgroundColor 设置文本背景色并返回文本以支持链式调用。
func (t *Text) BackgroundColor(c color.Color) *Text {
	t.style.Bg = c
	return t
}

// Alignment 设置文本对齐方式并返回文本以支持链式调用。
func (t *Text) Alignment(alignment string) *Text {
	t.alignment = alignment
	return t
}

// Wrap 启用文本换行并返回文本以支持链式调用。
func (t *Text) Wrap(wrap bool) *Text {
	t.wrap = wrap
	return t
}

// Content 返回文本内容（供外部访问）。
func (t *Text) Content() string {
	return t.content
}

// Draw 将文本渲染到屏幕上。
func (t *Text) Draw(scr uv.Screen, area uv.Rectangle) {
	t.SetBounds(area)

	if t.content == "" {
		return
	}

	// 如果指定了样式，则应用到内容
	content := t.content
	if !t.style.IsZero() {
		content = t.style.Styled(content)
	}

	// 处理对齐
	if t.alignment != "" && t.alignment != AlignmentLeading {
		content = t.alignText(content, area.Dx())
	}

	// 创建带样式的字符串
	styled := uv.NewStyledString(content)
	styled.Wrap = t.wrap

	styled.Draw(scr, area)
}

// alignText 在给定宽度内对齐文本。
func (t *Text) alignText(content string, width int) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// 剥离 ANSI 代码以获取实际文本宽度
		plainText := ansi.Strip(line)
		textWidth := ansi.StringWidth(plainText)

		if textWidth >= width {
			result = append(result, line)
			continue
		}

		padding := width - textWidth

		switch t.alignment {
		case AlignmentCenter:
			leftPad := padding / 2
			rightPad := padding - leftPad
			aligned := strings.Repeat(" ", leftPad) + line + strings.Repeat(" ", rightPad)
			result = append(result, aligned)

		case AlignmentTrailing:
			aligned := strings.Repeat(" ", padding) + line
			result = append(result, aligned)

		default:
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// Layout 计算文本大小。
func (t *Text) Layout(constraints Constraints) Size {
	if t.content == "" {
		return Size{Width: 0, Height: 0}
	}

	// 计算尺寸
	lines := strings.Split(t.content, "\n")
	height := len(lines)

	width := 0
	for _, line := range lines {
		// 使用支持 ANSI 的宽度计算
		lineWidth := ansi.StringWidth(line)
		if lineWidth > width {
			width = lineWidth
		}
	}

	// 如果启用了换行，则应用
	if t.wrap && width > constraints.MaxWidth {
		width = constraints.MaxWidth
		totalChars := len(t.content)
		height = (totalChars + width - 1) / width
	}

	return constraints.Constrain(Size{Width: width, Height: height})
}

// Children 对于文本元素返回 nil。
func (t *Text) Children() []Element {
	return nil
}
