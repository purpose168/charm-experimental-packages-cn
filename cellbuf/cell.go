package cellbuf

import (
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

var (
	// BlankCell 是一个包含单个空格的单元格，宽度为 1，无样式或链接。
	BlankCell = Cell{Rune: ' ', Width: 1}

	// EmptyCell 只是一个空单元格，用于比较和作为宽单元格的占位符。
	EmptyCell = Cell{}
)

// Cell 表示终端屏幕中的单个单元格。
type Cell struct {
	// 单元格的样式。Nil 样式表示无样式。零值会打印重置序列。
	Style Style

	// Link 是单元格的超链接。
	Link Link

	// Comb 是单元格的组合符文。如果单元格是单个符文或
	// 是宽单元格的一部分的零宽度单元格，则为 nil。
	Comb []rune

	// Width 是字形簇的等宽宽度。
	Width int

	// Rune 是单元格的主符文。如果单元格是宽单元格的一部分，则为零。
	Rune rune
}

// Append 向单元格追加符文而不更改宽度。这在我们想要使用单元格存储
// 转义序列或其他不影响单元格宽度的符文时非常有用。
func (c *Cell) Append(r ...rune) {
	for i, r := range r {
		if i == 0 && c.Rune == 0 {
			c.Rune = r
			continue
		}
		c.Comb = append(c.Comb, r)
	}
}

// String 返回单元格的字符串内容，不包括任何样式、链接和转义序列。
func (c Cell) String() string {
	if c.Rune == 0 {
		return ""
	}
	if len(c.Comb) == 0 {
		return string(c.Rune)
	}
	return string(append([]rune{c.Rune}, c.Comb...))
}

// Equal 返回单元格是否等于另一个单元格。
func (c *Cell) Equal(o *Cell) bool {
	return o != nil &&
		c.Width == o.Width &&
		c.Rune == o.Rune &&
		runesEqual(c.Comb, o.Comb) &&
		c.Style.Equal(&o.Style) &&
		c.Link.Equal(&o.Link)
}

// Empty 返回单元格是否为空单元格。空单元格是宽度为 0、符文为 0 且无组合符文的单元格。
func (c Cell) Empty() bool {
	return c.Width == 0 &&
		c.Rune == 0 &&
		len(c.Comb) == 0
}

// Reset 将单元格重置为默认状态零值。
func (c *Cell) Reset() {
	c.Rune = 0
	c.Comb = nil
	c.Width = 0
	c.Style.Reset()
	c.Link.Reset()
}

// Clear 返回单元格是否仅包含不影响空格字符外观的属性。
func (c *Cell) Clear() bool {
	return c.Rune == ' ' && len(c.Comb) == 0 && c.Width == 1 && c.Style.Clear() && c.Link.Empty()
}

// Clone 返回单元格的副本。
func (c *Cell) Clone() (n *Cell) {
	n = new(Cell)
	*n = *c
	return n
}

// Blank 通过将符文设置为空格、comb 设置为 nil 并将宽度设置为 1 来使单元格成为空白单元格。
func (c *Cell) Blank() *Cell {
	c.Rune = ' '
	c.Comb = nil
	c.Width = 1
	return c
}

// Link 表示终端屏幕中的超链接。
type Link struct {
	URL    string
	Params string
}

// String 返回超链接的字符串表示形式。
func (h Link) String() string {
	return h.URL
}

// Reset 将超链接重置为默认状态零值。
func (h *Link) Reset() {
	h.URL = ""
	h.Params = ""
}

// Equal 返回超链接是否等于另一个超链接。
func (h *Link) Equal(o *Link) bool {
	return o != nil && h.URL == o.URL && h.Params == o.Params
}

// Empty 返回超链接是否为空。
func (h Link) Empty() bool {
	return h.URL == "" && h.Params == ""
}

// AttrMask 是可更改文本外观的文本属性的位掩码。
// 这些属性可以组合以创建不同的样式。
type AttrMask uint8

// 这些是可用的文本属性，可以组合以创建不同的样式。
const (
	BoldAttr AttrMask = 1 << iota
	FaintAttr
	ItalicAttr
	SlowBlinkAttr
	RapidBlinkAttr
	ReverseAttr
	ConcealAttr
	StrikethroughAttr

	ResetAttr AttrMask = 0
)

// Contains 返回属性掩码是否包含该属性。
func (a AttrMask) Contains(attr AttrMask) bool {
	return a&attr == attr
}

// UnderlineStyle 是文本使用的下划线样式。
type UnderlineStyle = ansi.UnderlineStyle

// 这些是可用的下划线样式。
const (
	NoUnderline     = ansi.UnderlineStyleNone
	SingleUnderline = ansi.UnderlineStyleSingle
	DoubleUnderline = ansi.UnderlineStyleDouble
	CurlyUnderline  = ansi.UnderlineStyleCurly
	DottedUnderline = ansi.UnderlineStyleDotted
	DashedUnderline = ansi.UnderlineStyleDashed
)

// Style 表示单元格的样式。
type Style struct {
	Fg      ansi.Color
	Bg      ansi.Color
	Ul      ansi.Color
	Attrs   AttrMask
	UlStyle UnderlineStyle
}

// Sequence 返回设置样式的 ANSI 序列。
func (s Style) Sequence() string {
	if s.Empty() {
		return ansi.ResetStyle
	}

	var b ansi.Style

	if s.Attrs != 0 { //nolint:nestif
		if s.Attrs&BoldAttr != 0 {
			b = b.Bold()
		}
		if s.Attrs&FaintAttr != 0 {
			b = b.Faint()
		}
		if s.Attrs&ItalicAttr != 0 {
			b = b.Italic(true)
		}
		if s.Attrs&SlowBlinkAttr != 0 {
			b = b.Blink(true)
		}
		if s.Attrs&RapidBlinkAttr != 0 {
			b = b.RapidBlink(true)
		}
		if s.Attrs&ReverseAttr != 0 {
			b = b.Reverse(true)
		}
		if s.Attrs&ConcealAttr != 0 {
			b = b.Conceal(true)
		}
		if s.Attrs&StrikethroughAttr != 0 {
			b = b.Strikethrough(true)
		}
	}
	if s.UlStyle != NoUnderline {
		switch u := s.UlStyle; u {
		case NoUnderline:
			b = b.Underline(false)
		default:
			b = b.Underline(true)
			b = b.UnderlineStyle(u)
		}
	}
	if s.Fg != nil {
		b = b.ForegroundColor(s.Fg)
	}
	if s.Bg != nil {
		b = b.BackgroundColor(s.Bg)
	}
	if s.Ul != nil {
		b = b.UnderlineColor(s.Ul)
	}

	return b.String()
}

// DiffSequence 返回将样式设置为与另一种样式的差异的 ANSI 序列。
func (s Style) DiffSequence(o Style) string {
	if o.Empty() {
		return s.Sequence()
	}

	var b ansi.Style

	if !colorEqual(s.Fg, o.Fg) {
		b = b.ForegroundColor(s.Fg)
	}

	if !colorEqual(s.Bg, o.Bg) {
		b = b.BackgroundColor(s.Bg)
	}

	if !colorEqual(s.Ul, o.Ul) {
		b = b.UnderlineColor(s.Ul)
	}

	var (
		noBlink  bool
		isNormal bool
	)

	if s.Attrs != o.Attrs { //nolint:nestif
		if s.Attrs&BoldAttr != o.Attrs&BoldAttr {
			if s.Attrs&BoldAttr != 0 {
				b = b.Bold()
			} else if !isNormal {
				isNormal = true
				b = b.Normal()
			}
		}
		if s.Attrs&FaintAttr != o.Attrs&FaintAttr {
			if s.Attrs&FaintAttr != 0 {
				b = b.Faint()
			} else if !isNormal {
				b = b.Normal()
			}
		}
		if s.Attrs&ItalicAttr != o.Attrs&ItalicAttr {
			b = b.Italic(s.Attrs&ItalicAttr != 0)
		}
		if s.Attrs&SlowBlinkAttr != o.Attrs&SlowBlinkAttr {
			if s.Attrs&SlowBlinkAttr != 0 {
				b = b.Blink(true)
			} else if !noBlink {
				noBlink = true
				b = b.Blink(false)
			}
		}
		if s.Attrs&RapidBlinkAttr != o.Attrs&RapidBlinkAttr {
			if s.Attrs&RapidBlinkAttr != 0 {
				b = b.RapidBlink(true)
			} else if !noBlink {
				b = b.Blink(false)
			}
		}
		if s.Attrs&ReverseAttr != o.Attrs&ReverseAttr {
			b = b.Reverse(s.Attrs&ReverseAttr != 0)
		}
		if s.Attrs&ConcealAttr != o.Attrs&ConcealAttr {
			b = b.Conceal(s.Attrs&ConcealAttr != 0)
		}
		if s.Attrs&StrikethroughAttr != o.Attrs&StrikethroughAttr {
			b = b.Strikethrough(s.Attrs&StrikethroughAttr != 0)
		}
	}

	if s.UlStyle != o.UlStyle {
		b = b.UnderlineStyle(s.UlStyle)
	}

	return b.String()
}

// Equal 如果样式等于另一种样式，则返回 true。
func (s *Style) Equal(o *Style) bool {
	return s.Attrs == o.Attrs &&
		s.UlStyle == o.UlStyle &&
		colorEqual(s.Fg, o.Fg) &&
		colorEqual(s.Bg, o.Bg) &&
		colorEqual(s.Ul, o.Ul)
}

func colorEqual(c, o ansi.Color) bool {
	if c == nil && o == nil {
		return true
	}
	if c == nil || o == nil {
		return false
	}
	cr, cg, cb, ca := c.RGBA()
	or, og, ob, oa := o.RGBA()
	return cr == or && cg == og && cb == ob && ca == oa
}

// Bold 设置粗体属性。
func (s *Style) Bold(v bool) *Style {
	if v {
		s.Attrs |= BoldAttr
	} else {
		s.Attrs &^= BoldAttr
	}
	return s
}

// Faint 设置 faint 属性。
func (s *Style) Faint(v bool) *Style {
	if v {
		s.Attrs |= FaintAttr
	} else {
		s.Attrs &^= FaintAttr
	}
	return s
}

// Italic 设置斜体属性。
func (s *Style) Italic(v bool) *Style {
	if v {
		s.Attrs |= ItalicAttr
	} else {
		s.Attrs &^= ItalicAttr
	}
	return s
}

// SlowBlink 设置慢闪烁属性。
func (s *Style) SlowBlink(v bool) *Style {
	if v {
		s.Attrs |= SlowBlinkAttr
	} else {
		s.Attrs &^= SlowBlinkAttr
	}
	return s
}

// RapidBlink 设置快闪烁属性。
func (s *Style) RapidBlink(v bool) *Style {
	if v {
		s.Attrs |= RapidBlinkAttr
	} else {
		s.Attrs &^= RapidBlinkAttr
	}
	return s
}

// Reverse 设置反转属性。
func (s *Style) Reverse(v bool) *Style {
	if v {
		s.Attrs |= ReverseAttr
	} else {
		s.Attrs &^= ReverseAttr
	}
	return s
}

// Conceal 设置隐藏属性。
func (s *Style) Conceal(v bool) *Style {
	if v {
		s.Attrs |= ConcealAttr
	} else {
		s.Attrs &^= ConcealAttr
	}
	return s
}

// Strikethrough 设置删除线属性。
func (s *Style) Strikethrough(v bool) *Style {
	if v {
		s.Attrs |= StrikethroughAttr
	} else {
		s.Attrs &^= StrikethroughAttr
	}
	return s
}

// UnderlineStyle 设置下划线样式。
func (s *Style) UnderlineStyle(style UnderlineStyle) *Style {
	s.UlStyle = style
	return s
}

// Underline 设置下划线属性。
// 这是 [UnderlineStyle] 的语法糖。
func (s *Style) Underline(v bool) *Style {
	if v {
		return s.UnderlineStyle(SingleUnderline)
	}
	return s.UnderlineStyle(NoUnderline)
}

// Foreground 设置前景色。
func (s *Style) Foreground(c ansi.Color) *Style {
	s.Fg = c
	return s
}

// Background 设置背景色。
func (s *Style) Background(c ansi.Color) *Style {
	s.Bg = c
	return s
}

// UnderlineColor 设置下划线颜色。
func (s *Style) UnderlineColor(c ansi.Color) *Style {
	s.Ul = c
	return s
}

// Reset 将样式重置为默认值。
func (s *Style) Reset() *Style {
	s.Fg = nil
	s.Bg = nil
	s.Ul = nil
	s.Attrs = ResetAttr
	s.UlStyle = NoUnderline
	return s
}

// Empty 如果样式为空，则返回 true。
func (s *Style) Empty() bool {
	return s.Fg == nil && s.Bg == nil && s.Ul == nil && s.Attrs == ResetAttr && s.UlStyle == NoUnderline
}

// Clear 返回样式是否仅包含不影响空格字符外观的属性。
func (s *Style) Clear() bool {
	return s.UlStyle == NoUnderline &&
		s.Attrs&^(BoldAttr|FaintAttr|ItalicAttr|SlowBlinkAttr|RapidBlinkAttr) == 0 &&
		s.Fg == nil &&
		s.Bg == nil &&
		s.Ul == nil
}

func runesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, r := range a {
		if r != b[i] {
			return false
		}
	}
	return true
}
