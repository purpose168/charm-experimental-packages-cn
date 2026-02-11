package ansi

import (
	"fmt"
	"image/color"
)

// SetPalette 设置给定索引的调色板颜色。索引是0到15之间的16色索引。颜色是24位RGB颜色。
//
//	OSC P n rrggbb BEL
//
// 其中n是十六进制的颜色索引（0-f），rrggbb是十六进制格式的颜色（例如，ff0000表示红色）。
//
// 此序列特定于Linux控制台，可能不适用于其他终端模拟器。
//
// 请参阅 https://man7.org/linux/man-pages/man4/console_codes.4.html
func SetPalette(i int, c color.Color) string {
	if c == nil || i < 0 || i > 15 {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("\x1b]P%x%02x%02x%02x\x07", i, r>>8, g>>8, b>>8)
}

// ResetPalette 将颜色调色板重置为默认值。
//
// 此序列特定于Linux控制台，可能不适用于其他终端模拟器。
//
// 请参阅 https://man7.org/linux/man-pages/man4/console_codes.4.html
const ResetPalette = "\x1b]R\x07"
