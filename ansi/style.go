package ansi

import (
	"image/color"
	"strconv"
	"strings"
)

// ResetStyle 是一个 SGR（选择图形渲染）样式序列，用于重置所有属性。
// 参见：https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
const ResetStyle = "\x1b[m"

// Attr 是一个 SGR（选择图形渲染）样式属性。
type Attr = int

// Style 表示一个 ANSI SGR（选择图形渲染）样式。
type Style []string

// NewStyle 返回一个带有给定属性的新样式。属性是控制文本格式的 SGR（选择图形渲染）代码，如粗体、斜体、颜色等。
func NewStyle(attrs ...Attr) Style {
	if len(attrs) == 0 {
		return Style{}
	}
	s := make(Style, 0, len(attrs))
	for _, a := range attrs {
		attr, ok := attrStrings[a]
		if ok {
			s = append(s, attr)
		} else {
			if a < 0 {
				a = 0
			}
			s = append(s, strconv.Itoa(a))
		}
	}
	return s
}

// String 返回给定样式的 ANSI SGR（选择图形渲染）样式序列。
func (s Style) String() string {
	if len(s) == 0 {
		return ResetStyle
	}
	return "\x1b[" + strings.Join(s, ";") + "m"
}

// Styled 返回一个应用了给定样式的带样式字符串。样式在字符串开头应用，在结尾重置。
func (s Style) Styled(str string) string {
	if len(s) == 0 {
		return str
	}
	return s.String() + str + ResetStyle
}

// Reset 向样式添加重置样式属性。这会将所有格式属性重置为默认值。
func (s Style) Reset() Style {
	return append(s, attrReset)
}

// Bold 向样式添加粗体或正常强度的样式属性。
// 您可以使用 [Style.Normal] 重置为正常强度。
func (s Style) Bold() Style {
	return append(s, attrBold)
}

// Faint 向样式添加微弱或正常强度的样式属性。
// 您可以使用 [Style.Normal] 重置为正常强度。
func (s Style) Faint() Style {
	return append(s, attrFaint)
}

// Italic 向样式添加斜体或非斜体样式属性。
// 当 v 为 true 时，文本以斜体呈现。当为 false 时，斜体被禁用。
func (s Style) Italic(v bool) Style {
	if v {
		return append(s, attrItalic)
	}
	return append(s, attrNoItalic)
}

// Underline 向样式添加下划线或非下划线样式属性。
// 当 v 为 true 时，文本带下划线。当为 false 时，下划线被禁用。
func (s Style) Underline(v bool) Style {
	if v {
		return append(s, attrUnderline)
	}
	return append(s, attrNoUnderline)
}

// UnderlineStyle 向样式添加下划线样式属性。
// 支持多种下划线样式，包括单线、双线、波浪线、点线和虚线。
func (s Style) UnderlineStyle(u Underline) Style {
	switch u {
	case UnderlineNone:
		return s.Underline(false)
	case UnderlineSingle:
		return s.Underline(true)
	case UnderlineDouble:
		return append(s, underlineDouble)
	case UnderlineCurly:
		return append(s, underlineCurly)
	case UnderlineDotted:
		return append(s, underlineDotted)
	case UnderlineDashed:
		return append(s, underlineDashed)
	}
	return s
}

// Blink 向样式添加慢速闪烁或非闪烁样式属性。
// 当 v 为 true 时，文本缓慢闪烁（每分钟少于 150 次）。当为 false 时，闪烁被禁用。
func (s Style) Blink(v bool) Style {
	if v {
		return append(s, attrBlink)
	}
	return append(s, attrNoBlink)
}

// RapidBlink 向样式添加快速闪烁或非闪烁样式属性。
// 当 v 为 true 时，文本快速闪烁（每分钟 150 次以上）。当为 false 时，闪烁被禁用。
//
// 注意：这在终端模拟器中支持度不高。
func (s Style) RapidBlink(v bool) Style {
	if v {
		return append(s, attrRapidBlink)
	}
	return append(s, attrNoBlink)
}

// Reverse 向样式添加反转或非反转样式属性。
// 当 v 为 true 时，前景色和背景色交换。当为 false 时，反转视频被禁用。
func (s Style) Reverse(v bool) Style {
	if v {
		return append(s, attrReverse)
	}
	return append(s, attrNoReverse)
}

// Conceal 向样式添加隐藏或非隐藏样式属性。
// 当 v 为 true 时，文本被隐藏/掩盖。当为 false 时，隐藏被禁用。
func (s Style) Conceal(v bool) Style {
	if v {
		return append(s, attrConceal)
	}
	return append(s, attrNoConceal)
}

// Strikethrough 向样式添加删除线或非删除线样式属性。
// 当 v 为 true 时，文本中间会显示一条水平线。当为 false 时，删除线被禁用。
func (s Style) Strikethrough(v bool) Style {
	if v {
		return append(s, attrStrikethrough)
	}
	return append(s, attrNoStrikethrough)
}

// Normal 向样式添加正常强度样式属性。这会重置 [Style.Bold] 和 [Style.Faint] 属性。
func (s Style) Normal() Style {
	return append(s, attrNormalIntensity)
}

// NoItalic 向样式添加非斜体样式属性。
//
// 已弃用：请使用 [Style.Italic](false) 代替。
func (s Style) NoItalic() Style {
	return append(s, attrNoItalic)
}

// NoUnderline 向样式添加非下划线样式属性。
//
// 已弃用：请使用 [Style.Underline](false) 代替。
func (s Style) NoUnderline() Style {
	return append(s, attrNoUnderline)
}

// NoBlink 向样式添加非闪烁样式属性。
//
// 已弃用：请使用 [Style.Blink](false) 或 [Style.RapidBlink](false) 代替。
func (s Style) NoBlink() Style {
	return append(s, attrNoBlink)
}

// NoReverse 向样式添加非反转样式属性。
//
// 已弃用：请使用 [Style.Reverse](false) 代替。
func (s Style) NoReverse() Style {
	return append(s, attrNoReverse)
}

// NoConceal 向样式添加非隐藏样式属性。
//
// 已弃用：请使用 [Style.Conceal](false) 代替。
func (s Style) NoConceal() Style {
	return append(s, attrNoConceal)
}

// NoStrikethrough 向样式添加非删除线样式属性。
//
// 已弃用：请使用 [Style.Strikethrough](false) 代替。
func (s Style) NoStrikethrough() Style {
	return append(s, attrNoStrikethrough)
}

// DefaultForegroundColor 向样式添加默认前景色样式属性。
//
// 已弃用：请使用 [Style.ForegroundColor](nil) 代替。
func (s Style) DefaultForegroundColor() Style {
	return append(s, attrDefaultForegroundColor)
}

// DefaultBackgroundColor 向样式添加默认背景色样式属性。
//
// 已弃用：请使用 [Style.BackgroundColor](nil) 代替。
func (s Style) DefaultBackgroundColor() Style {
	return append(s, attrDefaultBackgroundColor)
}

// DefaultUnderlineColor 向样式添加默认下划线颜色样式属性。
//
// 已弃用：请使用 [Style.UnderlineColor](nil) 代替。
func (s Style) DefaultUnderlineColor() Style {
	return append(s, attrDefaultUnderlineColor)
}

// ForegroundColor 向样式添加前景色样式属性。
// 如果 c 为 nil，则使用默认前景色。支持 [BasicColor]、[IndexedColor]（256 色）和 [color.Color]（24 位 RGB）。
func (s Style) ForegroundColor(c Color) Style {
	if c == nil {
		return append(s, attrDefaultForegroundColor)
	}
	return append(s, foregroundColorString(c))
}

// BackgroundColor 向样式添加背景色样式属性。
// 如果 c 为 nil，则使用默认背景色。支持 [BasicColor]、[IndexedColor]（256 色）和 [color.Color]（24 位 RGB）。
func (s Style) BackgroundColor(c Color) Style {
	if c == nil {
		return append(s, attrDefaultBackgroundColor)
	}
	return append(s, backgroundColorString(c))
}

// UnderlineColor 向样式添加下划线颜色样式属性。
// 如果 c 为 nil，则使用默认下划线颜色。支持 [BasicColor]、[IndexedColor]（256 色）和 [color.Color]（24 位 RGB）。
func (s Style) UnderlineColor(c Color) Style {
	if c == nil {
		return append(s, attrDefaultUnderlineColor)
	}
	return append(s, underlineColorString(c))
}

// Underline 表示一个 ANSI SGR（选择图形渲染）下划线样式。
type Underline = byte

// UnderlineStyle 表示一个 ANSI SGR（选择图形渲染）下划线样式。
//
// 已弃用：请使用 [Underline] 代替。
type UnderlineStyle = byte

const (
	underlineDouble = "4:2"
	underlineCurly  = "4:3"
	underlineDotted = "4:4"
	underlineDashed = "4:5"
)

// 下划线样式常量。
const (
	UnderlineNone Underline = iota
	UnderlineSingle
	UnderlineDouble
	UnderlineCurly
	UnderlineDotted
	UnderlineDashed
)

// 下划线样式常量。
//
// 已弃用：请使用 [UnderlineNone]、[UnderlineSingle] 等代替。
const (
	NoUnderlineStyle Underline = iota
	SingleUnderlineStyle
	DoubleUnderlineStyle
	CurlyUnderlineStyle
	DottedUnderlineStyle
	DashedUnderlineStyle
)

// 下划线样式常量。
//
// 已弃用：请使用 [UnderlineNone]、[UnderlineSingle] 等代替。
const (
	UnderlineStyleNone Underline = iota
	UnderlineStyleSingle
	UnderlineStyleDouble
	UnderlineStyleCurly
	UnderlineStyleDotted
	UnderlineStyleDashed
)

// SGR（选择图形渲染）样式属性。
// 参见：https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
const (
	AttrReset                        Attr = 0
	AttrBold                         Attr = 1
	AttrFaint                        Attr = 2
	AttrItalic                       Attr = 3
	AttrUnderline                    Attr = 4
	AttrBlink                        Attr = 5
	AttrRapidBlink                   Attr = 6
	AttrReverse                      Attr = 7
	AttrConceal                      Attr = 8
	AttrStrikethrough                Attr = 9
	AttrNormalIntensity              Attr = 22
	AttrNoItalic                     Attr = 23
	AttrNoUnderline                  Attr = 24
	AttrNoBlink                      Attr = 25
	AttrNoReverse                    Attr = 27
	AttrNoConceal                    Attr = 28
	AttrNoStrikethrough              Attr = 29
	AttrBlackForegroundColor         Attr = 30
	AttrRedForegroundColor           Attr = 31
	AttrGreenForegroundColor         Attr = 32
	AttrYellowForegroundColor        Attr = 33
	AttrBlueForegroundColor          Attr = 34
	AttrMagentaForegroundColor       Attr = 35
	AttrCyanForegroundColor          Attr = 36
	AttrWhiteForegroundColor         Attr = 37
	AttrExtendedForegroundColor      Attr = 38
	AttrDefaultForegroundColor       Attr = 39
	AttrBlackBackgroundColor         Attr = 40
	AttrRedBackgroundColor           Attr = 41
	AttrGreenBackgroundColor         Attr = 42
	AttrYellowBackgroundColor        Attr = 43
	AttrBlueBackgroundColor          Attr = 44
	AttrMagentaBackgroundColor       Attr = 45
	AttrCyanBackgroundColor          Attr = 46
	AttrWhiteBackgroundColor         Attr = 47
	AttrExtendedBackgroundColor      Attr = 48
	AttrDefaultBackgroundColor       Attr = 49
	AttrExtendedUnderlineColor       Attr = 58
	AttrDefaultUnderlineColor        Attr = 59
	AttrBrightBlackForegroundColor   Attr = 90
	AttrBrightRedForegroundColor     Attr = 91
	AttrBrightGreenForegroundColor   Attr = 92
	AttrBrightYellowForegroundColor  Attr = 93
	AttrBrightBlueForegroundColor    Attr = 94
	AttrBrightMagentaForegroundColor Attr = 95
	AttrBrightCyanForegroundColor    Attr = 96
	AttrBrightWhiteForegroundColor   Attr = 97
	AttrBrightBlackBackgroundColor   Attr = 100
	AttrBrightRedBackgroundColor     Attr = 101
	AttrBrightGreenBackgroundColor   Attr = 102
	AttrBrightYellowBackgroundColor  Attr = 103
	AttrBrightBlueBackgroundColor    Attr = 104
	AttrBrightMagentaBackgroundColor Attr = 105
	AttrBrightCyanBackgroundColor    Attr = 106
	AttrBrightWhiteBackgroundColor   Attr = 107

	AttrRGBColorIntroducer      Attr = 2
	AttrExtendedColorIntroducer Attr = 5
)

// SGR（选择图形渲染）样式属性。
//
// 已弃用：请使用 Attr* 常量代替。
const (
	ResetAttr                        = AttrReset
	BoldAttr                         = AttrBold
	FaintAttr                        = AttrFaint
	ItalicAttr                       = AttrItalic
	UnderlineAttr                    = AttrUnderline
	SlowBlinkAttr                    = AttrBlink
	RapidBlinkAttr                   = AttrRapidBlink
	ReverseAttr                      = AttrReverse
	ConcealAttr                      = AttrConceal
	StrikethroughAttr                = AttrStrikethrough
	NormalIntensityAttr              = AttrNormalIntensity
	NoItalicAttr                     = AttrNoItalic
	NoUnderlineAttr                  = AttrNoUnderline
	NoBlinkAttr                      = AttrNoBlink
	NoReverseAttr                    = AttrNoReverse
	NoConcealAttr                    = AttrNoConceal
	NoStrikethroughAttr              = AttrNoStrikethrough
	BlackForegroundColorAttr         = AttrBlackForegroundColor
	RedForegroundColorAttr           = AttrRedForegroundColor
	GreenForegroundColorAttr         = AttrGreenForegroundColor
	YellowForegroundColorAttr        = AttrYellowForegroundColor
	BlueForegroundColorAttr          = AttrBlueForegroundColor
	MagentaForegroundColorAttr       = AttrMagentaForegroundColor
	CyanForegroundColorAttr          = AttrCyanForegroundColor
	WhiteForegroundColorAttr         = AttrWhiteForegroundColor
	ExtendedForegroundColorAttr      = AttrExtendedForegroundColor
	DefaultForegroundColorAttr       = AttrDefaultForegroundColor
	BlackBackgroundColorAttr         = AttrBlackBackgroundColor
	RedBackgroundColorAttr           = AttrRedBackgroundColor
	GreenBackgroundColorAttr         = AttrGreenBackgroundColor
	YellowBackgroundColorAttr        = AttrYellowBackgroundColor
	BlueBackgroundColorAttr          = AttrBlueBackgroundColor
	MagentaBackgroundColorAttr       = AttrMagentaBackgroundColor
	CyanBackgroundColorAttr          = AttrCyanBackgroundColor
	WhiteBackgroundColorAttr         = AttrWhiteBackgroundColor
	ExtendedBackgroundColorAttr      = AttrExtendedBackgroundColor
	DefaultBackgroundColorAttr       = AttrDefaultBackgroundColor
	ExtendedUnderlineColorAttr       = AttrExtendedUnderlineColor
	DefaultUnderlineColorAttr        = AttrDefaultUnderlineColor
	BrightBlackForegroundColorAttr   = AttrBrightBlackForegroundColor
	BrightRedForegroundColorAttr     = AttrBrightRedForegroundColor
	BrightGreenForegroundColorAttr   = AttrBrightGreenForegroundColor
	BrightYellowForegroundColorAttr  = AttrBrightYellowForegroundColor
	BrightBlueForegroundColorAttr    = AttrBrightBlueForegroundColor
	BrightMagentaForegroundColorAttr = AttrBrightMagentaForegroundColor
	BrightCyanForegroundColorAttr    = AttrBrightCyanForegroundColor
	BrightWhiteForegroundColorAttr   = AttrBrightWhiteForegroundColor
	BrightBlackBackgroundColorAttr   = AttrBrightBlackBackgroundColor
	BrightRedBackgroundColorAttr     = AttrBrightRedBackgroundColor
	BrightGreenBackgroundColorAttr   = AttrBrightGreenBackgroundColor
	BrightYellowBackgroundColorAttr  = AttrBrightYellowBackgroundColor
	BrightBlueBackgroundColorAttr    = AttrBrightBlueBackgroundColor
	BrightMagentaBackgroundColorAttr = AttrBrightMagentaBackgroundColor
	BrightCyanBackgroundColorAttr    = AttrBrightCyanBackgroundColor
	BrightWhiteBackgroundColorAttr   = AttrBrightWhiteBackgroundColor
	RGBColorIntroducerAttr           = AttrRGBColorIntroducer
	ExtendedColorIntroducerAttr      = AttrExtendedColorIntroducer
)

const (
	attrReset                        = "0"
	attrBold                         = "1"
	attrFaint                        = "2"
	attrItalic                       = "3"
	attrUnderline                    = "4"
	attrBlink                        = "5"
	attrRapidBlink                   = "6"
	attrReverse                      = "7"
	attrConceal                      = "8"
	attrStrikethrough                = "9"
	attrNormalIntensity              = "22"
	attrNoItalic                     = "23"
	attrNoUnderline                  = "24"
	attrNoBlink                      = "25"
	attrNoReverse                    = "27"
	attrNoConceal                    = "28"
	attrNoStrikethrough              = "29"
	attrBlackForegroundColor         = "30"
	attrRedForegroundColor           = "31"
	attrGreenForegroundColor         = "32"
	attrYellowForegroundColor        = "33"
	attrBlueForegroundColor          = "34"
	attrMagentaForegroundColor       = "35"
	attrCyanForegroundColor          = "36"
	attrWhiteForegroundColor         = "37"
	attrExtendedForegroundColor      = "38"
	attrDefaultForegroundColor       = "39"
	attrBlackBackgroundColor         = "40"
	attrRedBackgroundColor           = "41"
	attrGreenBackgroundColor         = "42"
	attrYellowBackgroundColor        = "43"
	attrBlueBackgroundColor          = "44"
	attrMagentaBackgroundColor       = "45"
	attrCyanBackgroundColor          = "46"
	attrWhiteBackgroundColor         = "47"
	attrExtendedBackgroundColor      = "48"
	attrDefaultBackgroundColor       = "49"
	attrExtendedUnderlineColor       = "58"
	attrDefaultUnderlineColor        = "59"
	attrBrightBlackForegroundColor   = "90"
	attrBrightRedForegroundColor     = "91"
	attrBrightGreenForegroundColor   = "92"
	attrBrightYellowForegroundColor  = "93"
	attrBrightBlueForegroundColor    = "94"
	attrBrightMagentaForegroundColor = "95"
	attrBrightCyanForegroundColor    = "96"
	attrBrightWhiteForegroundColor   = "97"
	attrBrightBlackBackgroundColor   = "100"
	attrBrightRedBackgroundColor     = "101"
	attrBrightGreenBackgroundColor   = "102"
	attrBrightYellowBackgroundColor  = "103"
	attrBrightBlueBackgroundColor    = "104"
	attrBrightMagentaBackgroundColor = "105"
	attrBrightCyanBackgroundColor    = "106"
	attrBrightWhiteBackgroundColor   = "107"
)

// foregroundColorString 返回给定前景色的样式 SGR 属性。
// 参见：https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
func foregroundColorString(c Color) string {
	switch c := c.(type) {
	case nil:
		return attrDefaultForegroundColor
	case BasicColor:
		// 3-bit or 4-bit ANSI foreground
		// "3<n>" or "9<n>" where n is the color number from 0 to 7
		switch c {
		case Black:
			return attrBlackForegroundColor
		case Red:
			return attrRedForegroundColor
		case Green:
			return attrGreenForegroundColor
		case Yellow:
			return attrYellowForegroundColor
		case Blue:
			return attrBlueForegroundColor
		case Magenta:
			return attrMagentaForegroundColor
		case Cyan:
			return attrCyanForegroundColor
		case White:
			return attrWhiteForegroundColor
		case BrightBlack:
			return attrBrightBlackForegroundColor
		case BrightRed:
			return attrBrightRedForegroundColor
		case BrightGreen:
			return attrBrightGreenForegroundColor
		case BrightYellow:
			return attrBrightYellowForegroundColor
		case BrightBlue:
			return attrBrightBlueForegroundColor
		case BrightMagenta:
			return attrBrightMagentaForegroundColor
		case BrightCyan:
			return attrBrightCyanForegroundColor
		case BrightWhite:
			return attrBrightWhiteForegroundColor
		}
	case ExtendedColor:
		// 256-color ANSI foreground
		// "38;5;<n>"
		return "38;5;" + strconv.FormatUint(uint64(c), 10)
	case TrueColor, color.Color:
		// 24-bit "true color" foreground
		// "38;2;<r>;<g>;<b>"
		r, g, b, _ := c.RGBA()
		return "38;2;" +
			strconv.FormatUint(uint64(shift(r)), 10) + ";" +
			strconv.FormatUint(uint64(shift(g)), 10) + ";" +
			strconv.FormatUint(uint64(shift(b)), 10)
	}
	return attrDefaultForegroundColor
}

// backgroundColorString 返回给定背景色的样式 SGR 属性。
// 参见：https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
func backgroundColorString(c Color) string {
	switch c := c.(type) {
	case nil:
		return attrDefaultBackgroundColor
	case BasicColor:
		// 3-bit or 4-bit ANSI foreground
		// "4<n>" or "10<n>" where n is the color number from 0 to 7
		switch c {
		case Black:
			return attrBlackBackgroundColor
		case Red:
			return attrRedBackgroundColor
		case Green:
			return attrGreenBackgroundColor
		case Yellow:
			return attrYellowBackgroundColor
		case Blue:
			return attrBlueBackgroundColor
		case Magenta:
			return attrMagentaBackgroundColor
		case Cyan:
			return attrCyanBackgroundColor
		case White:
			return attrWhiteBackgroundColor
		case BrightBlack:
			return attrBrightBlackBackgroundColor
		case BrightRed:
			return attrBrightRedBackgroundColor
		case BrightGreen:
			return attrBrightGreenBackgroundColor
		case BrightYellow:
			return attrBrightYellowBackgroundColor
		case BrightBlue:
			return attrBrightBlueBackgroundColor
		case BrightMagenta:
			return attrBrightMagentaBackgroundColor
		case BrightCyan:
			return attrBrightCyanBackgroundColor
		case BrightWhite:
			return attrBrightWhiteBackgroundColor
		}
	case ExtendedColor:
		// 256-color ANSI foreground
		// "48;5;<n>"
		return "48;5;" + strconv.FormatUint(uint64(c), 10)
	case TrueColor, color.Color:
		// 24-bit "true color" foreground
		// "38;2;<r>;<g>;<b>"
		r, g, b, _ := c.RGBA()
		return "48;2;" +
			strconv.FormatUint(uint64(shift(r)), 10) + ";" +
			strconv.FormatUint(uint64(shift(g)), 10) + ";" +
			strconv.FormatUint(uint64(shift(b)), 10)
	}
	return attrDefaultBackgroundColor
}

// underlineColorString 返回给定下划线颜色的样式 SGR 属性。
// 参见：https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters
func underlineColorString(c Color) string {
	switch c := c.(type) {
	case nil:
		return attrDefaultUnderlineColor
	// NOTE: we can't use 3-bit and 4-bit ANSI color codes with underline
	// color, use 256-color instead.
	//
	// 256-color ANSI underline color
	// "58;5;<n>"
	case BasicColor:
		return "58;5;" + strconv.FormatUint(uint64(c), 10)
	case ExtendedColor:
		return "58;5;" + strconv.FormatUint(uint64(c), 10)
	case TrueColor, color.Color:
		// 24-bit "true color" foreground
		// "38;2;<r>;<g>;<b>"
		r, g, b, _ := c.RGBA()
		return "58;2;" +
			strconv.FormatUint(uint64(shift(r)), 10) + ";" +
			strconv.FormatUint(uint64(shift(g)), 10) + ";" +
			strconv.FormatUint(uint64(shift(b)), 10)
	}
	return attrDefaultUnderlineColor
}

// ReadStyleColor 从参数切片中解码颜色。它返回读取的参数数量和颜色。此函数用于按照 ITU T.416 标准读取 SGR 颜色参数。
//
// 它支持读取以下颜色类型：
//   - 0: 实现定义
//   - 1: 透明
//   - 2: RGB 直接颜色
//   - 3: CMY 直接颜色
//   - 4: CMYK 直接颜色
//   - 5: 索引颜色
//   - 6: RGBA 直接颜色（WezTerm 扩展）
//
// 参数可以用分号 (;) 或冒号 (:) 分隔。不允许混合使用分隔符。
//
// 规范支持定义颜色空间 ID、颜色容差 value 和容差颜色空间 ID。然而，这些值对返回的颜色没有影响，将被忽略。
//
// 此实现对规范做了一些修改：
//  1. 支持使用分号 (;) 分隔的传统颜色值（针对 RGB 和索引颜色）
//  2. 支持忽略和省略颜色空间 ID（第二个参数）（针对 RGB 颜色）
//  3. 支持忽略和省略第 6 个参数（针对 RGB 和 CMY 颜色）
//  4. 支持读取 RGBA 颜色
func ReadStyleColor(params Params, co *color.Color) int {
	if len(params) < 2 { // Need at least SGR type and color type
		return 0
	}

	// First parameter indicates one of 38, 48, or 58 (foreground, background, or underline)
	s := params[0]
	p := params[1]
	colorType := p.Param(0)
	n := 2

	paramsfn := func() (p1, p2, p3, p4 int) {
		// Where should we start reading the color?
		switch {
		case s.HasMore() && p.HasMore() && len(params) > 8 && params[2].HasMore() && params[3].HasMore() && params[4].HasMore() && params[5].HasMore() && params[6].HasMore() && params[7].HasMore():
			// We have color space id, a 6th parameter, a tolerance value, and a tolerance color space
			n += 7
			return params[3].Param(0), params[4].Param(0), params[5].Param(0), params[6].Param(0)
		case s.HasMore() && p.HasMore() && len(params) > 7 && params[2].HasMore() && params[3].HasMore() && params[4].HasMore() && params[5].HasMore() && params[6].HasMore():
			// We have color space id, a 6th parameter, and a tolerance value
			n += 6
			return params[3].Param(0), params[4].Param(0), params[5].Param(0), params[6].Param(0)
		case s.HasMore() && p.HasMore() && len(params) > 6 && params[2].HasMore() && params[3].HasMore() && params[4].HasMore() && params[5].HasMore():
			// We have color space id and a 6th parameter
			// 48 : 4 : : 1 : 2 : 3 :4
			n += 5
			return params[3].Param(0), params[4].Param(0), params[5].Param(0), params[6].Param(0)
		case s.HasMore() && p.HasMore() && len(params) > 5 && params[2].HasMore() && params[3].HasMore() && params[4].HasMore() && !params[5].HasMore():
			// We have color space
			// 48 : 3 : : 1 : 2 : 3
			n += 4
			return params[3].Param(0), params[4].Param(0), params[5].Param(0), -1
		case s.HasMore() && p.HasMore() && p.Param(0) == 2 && params[2].HasMore() && params[3].HasMore() && !params[4].HasMore():
			// We have color values separated by colons (:)
			// 48 : 2 : 1 : 2 : 3
			fallthrough
		case !s.HasMore() && !p.HasMore() && p.Param(0) == 2 && !params[2].HasMore() && !params[3].HasMore() && !params[4].HasMore():
			// Support legacy color values separated by semicolons (;)
			// 48 ; 2 ; 1 ; 2 ; 3
			n += 3
			return params[2].Param(0), params[3].Param(0), params[4].Param(0), -1
		}
		// Ambiguous SGR color
		return -1, -1, -1, -1
	}

	switch colorType {
	case 0: // implementation defined
		return 2
	case 1: // transparent
		*co = color.Transparent
		return 2
	case 2: // RGB direct color
		if len(params) < 5 {
			return 0
		}

		r, g, b, _ := paramsfn()
		if r == -1 || g == -1 || b == -1 {
			return 0
		}

		*co = color.RGBA{
			R: uint8(r), //nolint:gosec
			G: uint8(g), //nolint:gosec
			B: uint8(b), //nolint:gosec
			A: 0xff,
		}
		return n

	case 3: // CMY direct color
		if len(params) < 5 {
			return 0
		}

		c, m, y, _ := paramsfn()
		if c == -1 || m == -1 || y == -1 {
			return 0
		}

		*co = color.CMYK{
			C: uint8(c), //nolint:gosec
			M: uint8(m), //nolint:gosec
			Y: uint8(y), //nolint:gosec
			K: 0,
		}
		return n

	case 4: // CMYK direct color
		if len(params) < 6 {
			return 0
		}

		c, m, y, k := paramsfn()
		if c == -1 || m == -1 || y == -1 || k == -1 {
			return 0
		}

		*co = color.CMYK{
			C: uint8(c), //nolint:gosec
			M: uint8(m), //nolint:gosec
			Y: uint8(y), //nolint:gosec
			K: uint8(k), //nolint:gosec
		}
		return n

	case 5: // indexed color
		if len(params) < 3 {
			return 0
		}
		switch {
		case s.HasMore() && p.HasMore() && !params[2].HasMore():
			// Colon separated indexed color
			// 38 : 5 : 234
		case !s.HasMore() && !p.HasMore() && !params[2].HasMore():
			// Legacy semicolon indexed color
			// 38 ; 5 ; 234
		default:
			return 0
		}
		*co = ExtendedColor(params[2].Param(0)) //nolint:gosec
		return 3

	case 6: // RGBA direct color
		if len(params) < 6 {
			return 0
		}

		r, g, b, a := paramsfn()
		if r == -1 || g == -1 || b == -1 || a == -1 {
			return 0
		}

		*co = color.RGBA{
			R: uint8(r), //nolint:gosec
			G: uint8(g), //nolint:gosec
			B: uint8(b), //nolint:gosec
			A: uint8(a), //nolint:gosec
		}
		return n

	default:
		return 0
	}
}
