package input

import (
	"fmt"
	"image/color"
	"math"
)

// ForegroundColorEvent 表示前景色事件。当终端使用 [ansi.RequestForegroundColor] 请求终端前景色时，会发出此事件。
type ForegroundColorEvent struct{ color.Color }

// String 返回颜色的十六进制表示。
func (e ForegroundColorEvent) String() string {
	return colorToHex(e.Color)
}

// IsDark 返回颜色是否为深色。
func (e ForegroundColorEvent) IsDark() bool {
	return isDarkColor(e.Color)
}

// BackgroundColorEvent 表示背景色事件。当终端使用 [ansi.RequestBackgroundColor] 请求终端背景色时，会发出此事件。
type BackgroundColorEvent struct{ color.Color }

// String 返回颜色的十六进制表示。
func (e BackgroundColorEvent) String() string {
	return colorToHex(e)
}

// IsDark 返回颜色是否为深色。
func (e BackgroundColorEvent) IsDark() bool {
	return isDarkColor(e.Color)
}

// CursorColorEvent 表示光标颜色变化事件。当程序使用 [ansi.RequestCursorColor] 请求终端光标颜色时，会发出此事件。
type CursorColorEvent struct{ color.Color }

// String 返回颜色的十六进制表示。
func (e CursorColorEvent) String() string {
	return colorToHex(e)
}

// IsDark 返回颜色是否为深色。
func (e CursorColorEvent) IsDark() bool {
	return isDarkColor(e)
}

type shiftable interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

func shift[T shiftable](x T) T {
	if x > 0xff {
		x >>= 8
	}
	return x
}

func colorToHex(c color.Color) string {
	if c == nil {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", shift(r), shift(g), shift(b))
}

func getMaxMin(a, b, c float64) (ma, mi float64) {
	if a > b {
		ma = a
		mi = b
	} else {
		ma = b
		mi = a
	}
	if c > ma {
		ma = c
	} else if c < mi {
		mi = c
	}
	return ma, mi
}

func round(x float64) float64 {
	return math.Round(x*1000) / 1000
}

// rgbToHSL 将 RGB 三元组转换为 HSL 三元组。
func rgbToHSL(r, g, b uint8) (h, s, l float64) {
	// 将 uint32 预乘值转换为 uint8
	// r,g,b 值除以 255，将范围从 0..255 更改为 0..1：
	Rnot := float64(r) / 255
	Gnot := float64(g) / 255
	Bnot := float64(b) / 255
	Cmax, Cmin := getMaxMin(Rnot, Gnot, Bnot)
	Δ := Cmax - Cmin
	// 亮度计算：
	l = (Cmax + Cmin) / 2
	// 色相和饱和度计算：
	if Δ == 0 {
		h = 0
		s = 0
	} else {
		switch Cmax {
		case Rnot:
			h = 60 * (math.Mod((Gnot-Bnot)/Δ, 6))
		case Gnot:
			h = 60 * (((Bnot - Rnot) / Δ) + 2)
		case Bnot:
			h = 60 * (((Rnot - Gnot) / Δ) + 4)
		}
		if h < 0 {
			h += 360
		}

		s = Δ / (1 - math.Abs((2*l)-1))
	}

	return h, round(s), round(l)
}

// isDarkColor 返回给定颜色是否为深色。
func isDarkColor(c color.Color) bool {
	if c == nil {
		return true
	}

	r, g, b, _ := c.RGBA()
	_, _, l := rgbToHSL(uint8(r>>8), uint8(g>>8), uint8(b>>8)) //nolint:gosec
	return l < 0.5
}
