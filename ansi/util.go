package ansi

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// colorToHexString 返回颜色的十六进制字符串表示
func colorToHexString(c color.Color) string { //nolint:unused
	if c == nil {
		return "" // 如果颜色为nil，返回空字符串
	}
	shift := func(v uint32) uint32 {
		if v > 0xff {
			return v >> 8 // 如果值大于0xff，右移8位
		}
		return v // 否则直接返回
	}
	r, g, b, _ := c.RGBA()                       // 获取颜色的RGBA值
	r, g, b = shift(r), shift(g), shift(b)       // 转换为8位值
	return fmt.Sprintf("#%02x%02x%02x", r, g, b) // 格式化为十六进制字符串
}

// rgbToHex 将红、绿、蓝值转换为十六进制值
//
//	hex := rgbToHex(0, 0, 255) // 0x0000FF
func rgbToHex(r, g, b uint32) uint32 { //nolint:unused
	return r<<16 + g<<8 + b // 将红色左移16位，绿色左移8位，蓝色直接相加
}

// shiftable 是一个接口类型，定义了可以移位的无符号整数类型
type shiftable interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

// shift 将值移位，如果值大于0xff则右移8位
func shift[T shiftable](x T) T {
	if x > 0xff {
		x >>= 8 // 如果值大于0xff，右移8位
	}
	return x // 返回移位后的值
}

// XParseColor 是一个辅助函数，将字符串解析为color.Color
// 它提供了与Xlib中的XParseColor函数类似的接口
// 支持以下格式：
//
//   - #RGB
//   - #RRGGBB
//   - rgb:RRRR/GGGG/BBBB
//   - rgba:RRRR/GGGG/BBBB/AAAA
//
// 如果字符串不是有效的颜色，返回nil
//
// 请参阅：https://linux.die.net/man/3/xparsecolor
func XParseColor(s string) color.Color {
	switch {
	case strings.HasPrefix(s, "#"): // 处理#RGB或#RRGGBB格式
		c, err := colorful.Hex(s)
		if err != nil {
			return nil // 如果解析失败，返回nil
		}

		return c // 返回解析后的颜色
	case strings.HasPrefix(s, "rgb:"): // 处理rgb:RRRR/GGGG/BBBB格式
		parts := strings.Split(s[4:], "/") // 分割字符串
		if len(parts) != 3 {
			return nil // 如果分割后的部分不是3个，返回nil
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32) // 解析红色分量
		g, _ := strconv.ParseUint(parts[1], 16, 32) // 解析绿色分量
		b, _ := strconv.ParseUint(parts[2], 16, 32) // 解析蓝色分量

		// 转换为8位值并返回RGBA颜色
		return color.RGBA{uint8(shift(r)), uint8(shift(g)), uint8(shift(b)), 255} //nolint:gosec
	case strings.HasPrefix(s, "rgba:"): // 处理rgba:RRRR/GGGG/BBBB/AAAA格式
		parts := strings.Split(s[5:], "/") // 分割字符串
		if len(parts) != 4 {
			return nil // 如果分割后的部分不是4个，返回nil
		}

		r, _ := strconv.ParseUint(parts[0], 16, 32) // 解析红色分量
		g, _ := strconv.ParseUint(parts[1], 16, 32) // 解析绿色分量
		b, _ := strconv.ParseUint(parts[2], 16, 32) // 解析蓝色分量
		a, _ := strconv.ParseUint(parts[3], 16, 32) // 解析alpha分量

		// 转换为8位值并返回RGBA颜色
		return color.RGBA{uint8(shift(r)), uint8(shift(g)), uint8(shift(b)), uint8(shift(a))} //nolint:gosec
	}
	return nil // 如果格式不匹配，返回nil
}
