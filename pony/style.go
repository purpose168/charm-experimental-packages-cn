package pony

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// parseColor 将颜色字符串解析为 color.Color 类型。
// 支持：
//   - 命名颜色：red, blue, green 等
//   - 十六进制颜色：#FF0000, #f00
//   - RGB：rgb(255, 0, 0)
//   - ANSI 颜色：0-255
func parseColor(s string) (color.Color, error) {
	s = strings.TrimSpace(s)

	// 检查是否为十六进制颜色
	if strings.HasPrefix(s, "#") {
		return parseHexColor(s)
	}

	// 检查是否为 rgb() 格式
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") {
		return parseRGBColor(s)
	}

	// 检查是否为 ANSI 颜色代码 (0-255)
	if num, err := strconv.Atoi(s); err == nil && num >= 0 && num <= 255 {
		return ansiColor(num), nil
	}

	// 命名颜色
	return parseNamedColor(s)
}

// parseHexColor 解析十六进制颜色，如 #FF0000 或 #f00。
func parseHexColor(s string) (color.Color, error) {
	c, err := colorful.Hex(s)
	if err != nil {
		return nil, fmt.Errorf("无效的十六进制颜色: %w", err)
	}
	return c, nil
}

// parseRGBColor 解析 RGB 颜色，如 rgb(255, 0, 0)。
func parseRGBColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "rgb(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("rgb() 需要 3 个值")
	}

	var rgb [3]uint8
	for i, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || val < 0 || val > 255 {
			return nil, fmt.Errorf("无效的 rgb 值: %s", part)
		}
		rgb[i] = uint8(val)
	}

	return color.RGBA{R: rgb[0], G: rgb[1], B: rgb[2], A: 255}, nil
}

// parseNamedColor 解析命名颜色。
func parseNamedColor(s string) (color.Color, error) {
	s = strings.ToLower(s)

	// 基本 ANSI 颜色
	switch s {
	case "black":
		return color.RGBA{0, 0, 0, 255}, nil
	case "red":
		return color.RGBA{170, 0, 0, 255}, nil
	case "green":
		return color.RGBA{0, 170, 0, 255}, nil
	case "yellow":
		return color.RGBA{170, 85, 0, 255}, nil
	case "blue":
		return color.RGBA{0, 0, 170, 255}, nil
	case "magenta":
		return color.RGBA{170, 0, 170, 255}, nil
	case "cyan":
		return color.RGBA{0, 170, 170, 255}, nil
	case "white":
		return color.RGBA{170, 170, 170, 255}, nil

	// 明亮 ANSI 颜色
	case "bright-black", "gray", "grey":
		return color.RGBA{85, 85, 85, 255}, nil
	case "bright-red":
		return color.RGBA{255, 85, 85, 255}, nil
	case "bright-green":
		return color.RGBA{85, 255, 85, 255}, nil
	case "bright-yellow":
		return color.RGBA{255, 255, 85, 255}, nil
	case "bright-blue":
		return color.RGBA{85, 85, 255, 255}, nil
	case "bright-magenta":
		return color.RGBA{255, 85, 255, 255}, nil
	case "bright-cyan":
		return color.RGBA{85, 255, 255, 255}, nil
	case "bright-white":
		return color.RGBA{255, 255, 255, 255}, nil

	default:
		return nil, fmt.Errorf("未知颜色: %s", s)
	}
}

// ansiColor 返回 ANSI 256 色 palette 中的颜色。
func ansiColor(code int) color.Color {
	if code < 0 || code > 255 {
		return nil
	}

	// 0-7: 标准颜色
	if code < 8 {
		colors := []color.RGBA{
			{0, 0, 0, 255},       // 黑色
			{170, 0, 0, 255},     // 红色
			{0, 170, 0, 255},     // 绿色
			{170, 85, 0, 255},    // 黄色
			{0, 0, 170, 255},     // 蓝色
			{170, 0, 170, 255},   // 洋红色
			{0, 170, 170, 255},   // 青色
			{170, 170, 170, 255}, // 白色
		}
		return colors[code]
	}

	// 8-15: 明亮颜色
	if code < 16 {
		colors := []color.RGBA{
			{85, 85, 85, 255},    // 明亮黑色
			{255, 85, 85, 255},   // 明亮红色
			{85, 255, 85, 255},   // 明亮绿色
			{255, 255, 85, 255},  // 明亮黄色
			{85, 85, 255, 255},   // 明亮蓝色
			{255, 85, 255, 255},  // 明亮洋红色
			{85, 255, 255, 255},  // 明亮青色
			{255, 255, 255, 255}, // 明亮白色
		}
		return colors[code-8]
	}

	// 16-231: 216 色立方体
	if code < 232 {
		code -= 16
		r := (code / 36) * 51
		g := ((code % 36) / 6) * 51
		b := (code % 6) * 51
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255} //nolint:gosec // values bounded to 0-255
	}

	// 232-255: 灰度
	gray := 8 + (code-232)*10
	return color.RGBA{uint8(gray), uint8(gray), uint8(gray), 255} //nolint:gosec // value bounded to 0-255
}
