package ansi

import (
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

// HexColor 是一个可以作为十六进制字符串格式化的 [color.Color]。
type HexColor string

// RGBA 返回颜色的 RGBA 值。
func (h HexColor) RGBA() (r, g, b, a uint32) {
	hex := h.color()
	if hex == nil {
		return 0, 0, 0, 0
	}
	return hex.RGBA()
}

// Hex 返回颜色的十六进制表示。如果颜色无效，则返回空字符串。
func (h HexColor) Hex() string {
	hex := h.color()
	if hex == nil {
		return ""
	}
	return hex.Hex()
}

// String 将颜色作为十六进制字符串返回。如果颜色为 nil，则返回空字符串。
func (h HexColor) String() string {
	return h.Hex()
}

// color 返回 HexColor 的底层颜色。
func (h HexColor) color() *colorful.Color {
	hex, err := colorful.Hex(string(h))
	if err != nil {
		return nil
	}
	return &hex
}

// XRGBColor 是一个可以作为 XParseColor rgb: 字符串格式化的 [color.Color]。
//
// 参见: https://linux.die.net/man/3/xparsecolor
type XRGBColor struct {
	color.Color
}

// RGBA 返回颜色的 RGBA 值。
func (x XRGBColor) RGBA() (r, g, b, a uint32) {
	if x.Color == nil {
		return 0, 0, 0, 0
	}
	return x.Color.RGBA()
}

// String 将颜色作为 XParseColor rgb: 字符串返回。如果颜色为 nil，则返回空字符串。
func (x XRGBColor) String() string {
	if x.Color == nil {
		return ""
	}
	r, g, b, _ := x.Color.RGBA()
	// 获取低 8 位
	return fmt.Sprintf("rgb:%04x/%04x/%04x", r, g, b)
}

// XRGBAColor 是一个可以作为 XParseColor rgba: 字符串格式化的 [color.Color]。
//
// 参见: https://linux.die.net/man/3/xparsecolor
type XRGBAColor struct {
	color.Color
}

// RGBA 返回颜色的 RGBA 值。
func (x XRGBAColor) RGBA() (r, g, b, a uint32) {
	if x.Color == nil {
		return 0, 0, 0, 0
	}
	return x.Color.RGBA()
}

// String 将颜色作为 XParseColor rgba: 字符串返回。如果颜色为 nil，则返回空字符串。
func (x XRGBAColor) String() string {
	if x.Color == nil {
		return ""
	}
	r, g, b, a := x.RGBA()
	// 获取低 8 位
	return fmt.Sprintf("rgba:%04x/%04x/%04x/%04x", r, g, b, a)
}

// SetForegroundColor 返回一个设置默认终端前景色的序列。
//
//	OSC 10 ; color ST
//	OSC 10 ; color BEL
//
// 其中 color 是编码的颜色编号。大多数终端支持十六进制、
// XParseColor rgb: 和 rgba: 字符串。您可以使用 [HexColor]、[XRGBColor]
// 或 [XRGBAColor] 来格式化颜色。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func SetForegroundColor(s string) string {
	return "\x1b]10;" + s + "\x07"
}

// RequestForegroundColor 是一个请求当前默认终端前景色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const RequestForegroundColor = "\x1b]10;?\x07"

// ResetForegroundColor 是一个重置默认终端前景色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const ResetForegroundColor = "\x1b]110\x07"

// SetBackgroundColor 返回一个设置默认终端背景色的序列。
//
//	OSC 11 ; color ST
//	OSC 11 ; color BEL
//
// 其中 color 是编码的颜色编号。大多数终端支持十六进制、
// XParseColor rgb: 和 rgba: 字符串。您可以使用 [HexColor]、[XRGBColor]
// 或 [XRGBAColor] 来格式化颜色。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func SetBackgroundColor(s string) string {
	return "\x1b]11;" + s + "\x07"
}

// RequestBackgroundColor 是一个请求当前默认终端背景色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const RequestBackgroundColor = "\x1b]11;?\x07"

// ResetBackgroundColor 是一个重置默认终端背景色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const ResetBackgroundColor = "\x1b]111\x07"

// SetCursorColor 返回一个设置终端光标颜色的序列。
//
//	OSC 12 ; color ST
//	OSC 12 ; color BEL
//
// 其中 color 是编码的颜色编号。大多数终端支持十六进制、
// XParseColor rgb: 和 rgba: 字符串。您可以使用 [HexColor]、[XRGBColor]
// 或 [XRGBAColor] 来格式化颜色。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func SetCursorColor(s string) string {
	return "\x1b]12;" + s + "\x07"
}

// RequestCursorColor 是一个请求当前终端光标颜色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const RequestCursorColor = "\x1b]12;?\x07"

// ResetCursorColor 是一个重置终端光标颜色的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
const ResetCursorColor = "\x1b]112\x07"
