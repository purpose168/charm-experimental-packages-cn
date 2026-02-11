// Package sixel 提供 sixel 图形格式功能。
package sixel

import (
	"fmt"
	"image/color"
	"io"

	"github.com/lucasb-eyer/go-colorful"
)

// ErrInvalidColor 在 Sixel 颜色无效时返回。
var ErrInvalidColor = fmt.Errorf("无效颜色")

// WriteColor 向写入器写入 Sixel 颜色。如果 pu 为 0，则忽略其余参数。
func WriteColor(w io.Writer, pc, pu, px, py, pz int) (int, error) {
	if pu <= 0 || pu > 2 {
		return fmt.Fprintf(w, "#%d", pc) //nolint:wrapcheck
	}

	return fmt.Fprintf(w, "#%d;%d;%d;%d;%d", pc, pu, px, py, pz) //nolint:wrapcheck
}

// ConvertChannel 将颜色通道从 color.Color 的 0xffff 范围转换为
// Sixel RGB 格式的 0-100 范围。
func ConvertChannel(c uint32) uint32 {
	// 我们加 328 是因为那大约是 sixel 0-100 颜色范围内的 0.5，我们试图
	// 四舍五入到最近的值
	return (c + 328) * 100 / 0xffff
}

// FromColor 从 color.Color 返回 Sixel 颜色。它将颜色通道转换为 0-100 范围。
func FromColor(c color.Color) Color {
	if c == nil {
		return Color{}
	}

	r, g, b, _ := c.RGBA()
	return Color{
		Pu: 2, // 始终使用 RGB 格式 "2"
		Px: int(ConvertChannel(r)),
		Py: int(ConvertChannel(g)),
		Pz: int(ConvertChannel(b)),
	}
}

// DecodeColor 从字节切片解码 Sixel 颜色。它返回颜色和读取的字节数。
func DecodeColor(data []byte) (c Color, n int) {
	if len(data) == 0 || data[0] != ColorIntroducer {
		return c, n
	}

	if len(data) < 2 { // 最小长度是 2：引导符和数字。
		return c, n
	}

	// 解析颜色编号和可选的颜色系统。
	pc := &c.Pc
	for n = 1; n < len(data); n++ {
		if data[n] == ';' {
			if pc == &c.Pc {
				pc = &c.Pu
			} else {
				n++
				break
			}
		} else if data[n] >= '0' && data[n] <= '9' {
			*pc = (*pc)*10 + int(data[n]-'0')
		} else {
			break
		}
	}

	// 解析颜色分量。
	ptr := &c.Px
	for ; n < len(data); n++ {
		if data[n] == ';' { //nolint:nestif
			if ptr == &c.Px {
				ptr = &c.Py
			} else if ptr == &c.Py {
				ptr = &c.Pz
			} else {
				n++
				break
			}
		} else if data[n] >= '0' && data[n] <= '9' {
			*ptr = (*ptr)*10 + int(data[n]-'0')
		} else {
			break
		}
	}

	return c, n
}

// Color 表示 Sixel 颜色。
type Color struct {
	// Pc 是颜色编号 (0-255)。
	Pc int
	// Pu 是可选的颜色系统
	//  - 0: 默认颜色映射
	//  - 1: HLS
	//  - 2: RGB
	Pu int
	// 颜色分量范围从 0-100 表示 RGB 值。对于 HLS 格式，Px（色相）分量范围从 0-360 度，
	// 而 L（亮度）和 S（饱和度）为 0-100。
	Px, Py, Pz int
}

// RGBA 实现 color.Color 接口。
func (c Color) RGBA() (r, g, b, a uint32) {
	switch c.Pu {
	case 1:
		return sixelHLS(c.Px, c.Py, c.Pz).RGBA()
	case 2:
		return sixelRGB(c.Px, c.Py, c.Pz).RGBA()
	default:
		return colorPalette[c.Pc].RGBA()
	}
}

// #define PALVAL(n,a,m) (((n) * (a) + ((m) / 2)) / (m))
func palval(n, a, m int) int {
	return (n*a + m/2) / m
}

func sixelRGB(r, g, b int) color.Color {
	return color.NRGBA{uint8(palval(r, 0xff, 100)), uint8(palval(g, 0xff, 100)), uint8(palval(b, 0xff, 100)), 0xFF} //nolint:gosec
}

func sixelHLS(h, l, s int) color.Color {
	return colorful.Hsl(float64(h), float64(s)/100.0, float64(l)/100.0).Clamped()
}
