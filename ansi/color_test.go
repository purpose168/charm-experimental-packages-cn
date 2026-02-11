package ansi

import (
	"image/color"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

// TestRGBAToHex 测试 RGBA 到十六进制颜色的转换
func TestRGBAToHex(t *testing.T) {
	cases := []struct {
		r, g, b, a uint32 // RGBA 颜色分量
		want       uint32 // 预期的十六进制颜色值
	}{
		{0, 0, 255, 0xffff, 0x0000ff},
		{255, 255, 255, 0xffff, 0xffffff},
		{255, 0, 0, 0xffff, 0xffff0000},
	}

	for _, c := range cases {
		gotR, gotG, gotB, _ := TrueColor(c.want).RGBA()
		gotR /= 256
		gotG /= 256
		gotB /= 256
		if gotR != c.r || gotG != c.g || gotB != c.b {
			t.Errorf("RGBA() of TrueColor(%06x): got (%v, %v, %v), want (%v, %v, %v)",
				c.want, gotR, gotG, gotB, c.r, c.g, c.b)
		}
	}
}

// TestColorToHexString 测试颜色到十六进制字符串的转换
func TestColorToHexString(t *testing.T) {
	cases := []struct {
		color color.Color // 颜色
		want  string      // 预期的十六进制字符串
	}{
		{TrueColor(0x0000ff), "#0000ff"},
		{TrueColor(0xffffff), "#ffffff"},
		{TrueColor(0xff0000), "#ff0000"},
	}

	for _, c := range cases {
		got := colorToHexString(c.color)
		if got != c.want {
			t.Errorf("colorToHexString(%v): got %v, want %v", c.color, got, c.want)
		}
	}
}

// TestAnsiToRGB 测试 ANSI 颜色到 RGB 颜色的转换
func TestAnsiToRGB(t *testing.T) {
	cases := []struct {
		ansi    byte   // ANSI 颜色索引
		r, g, b uint32 // 预期的 RGB 颜色分量
	}{
		{0, 0, 0, 0},         // 黑色
		{1, 128, 0, 0},       // 红色
		{255, 238, 238, 238}, // 最高ANSI颜色（灰度）
	}

	for _, c := range cases {
		gotR, gotG, gotB, _ := ansiToRGB(c.ansi).RGBA()
		// 我们需要将值下移到8位
		gotR >>= 8
		gotR &= 0xff
		gotG >>= 8
		gotG &= 0xff
		gotB >>= 8
		gotB &= 0xff
		if gotR != c.r || gotG != c.g || gotB != c.b {
			t.Errorf("ansiToRGB(%v): got (%v, %v, %v), want (%v, %v, %v)",
				c.ansi, gotR, gotG, gotB, c.r, c.g, c.b)
		}
	}
}

// TestHexToRGB 测试十六进制颜色到 RGB 颜色的转换
func TestHexToRGB(t *testing.T) {
	cases := []struct {
		hex     uint32 // 十六进制颜色值
		r, g, b uint32 // 预期的 RGB 颜色分量
	}{
		{0x0000FF, 0, 0, 255},     // 蓝色
		{0xFFFFFF, 255, 255, 255}, // 白色
		{0xFF0000, 255, 0, 0},     // 红色
	}

	for _, c := range cases {
		gotR, gotG, gotB := hexToRGB(c.hex)
		if gotR != c.r || gotG != c.g || gotB != c.b {
			t.Errorf("hexToRGB(%v): got (%v, %v, %v), want (%v, %v, %v)",
				c.hex, gotR, gotG, gotB, c.r, c.g, c.b)
		}
	}
}

// TestHexTo256 测试十六进制颜色到 256 色索引的转换
func TestHexTo256(t *testing.T) {
	testCases := map[string]struct {
		input          colorful.Color // 输入颜色
		expectedHex    string         // 预期的十六进制字符串
		expectedOutput IndexedColor   // 预期的 256 色索引
	}{
		"白色": {
			input:          colorful.Color{R: 1, G: 1, B: 1},
			expectedHex:    "#ffffff",
			expectedOutput: 231,
		},
		"米白色": {
			input:          colorful.Color{R: 0.9333, G: 0.9333, B: 0.933},
			expectedHex:    "#eeeeee",
			expectedOutput: 255,
		},
		"比米白色稍亮": {
			input:          colorful.Color{R: 0.95, G: 0.95, B: 0.95},
			expectedHex:    "#f2f2f2",
			expectedOutput: 255,
		},
		"红色": {
			input:          colorful.Color{R: 1, G: 0, B: 0},
			expectedHex:    "#ff0000",
			expectedOutput: 196,
		},
		"银箔色": {
			input:          colorful.Color{R: 0.6863, G: 0.6863, B: 0.6863},
			expectedHex:    "#afafaf",
			expectedOutput: 145,
		},
		"银杯色": {
			input:          colorful.Color{R: 0.698, G: 0.698, B: 0.698},
			expectedHex:    "#b2b2b2",
			expectedOutput: 249,
		},
		"更接近银箔色": {
			input:          colorful.Color{R: 0.692, G: 0.692, B: 0.692},
			expectedHex:    "#b0b0b0",
			expectedOutput: 145,
		},
		"更接近银杯色": {
			input:          colorful.Color{R: 0.694, G: 0.694, B: 0.694},
			expectedHex:    "#b1b1b1",
			expectedOutput: 249,
		},
		"灰色": {
			input:          colorful.Color{R: 0.5, G: 0.5, B: 0.5},
			expectedHex:    "#808080",
			expectedOutput: 244,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			// hex := fmt.Sprintf("#%02x%02x%02x", uint8(testCase.input.R*255), uint8(testCase.input.G*255), uint8(testCase.input.B*255))
			output := Convert256(testCase.input)
			if testCase.input.Hex() != testCase.expectedHex {
				t.Errorf("Expected %+v to map to %s, but instead received %s", testCase.input, testCase.expectedHex, testCase.input.Hex())
			}
			if output != testCase.expectedOutput {
				t.Errorf("Expected truecolor %+v to map to 256 color %d, but instead received %d", testCase.input, testCase.expectedOutput, output)
			}
		})
	}
}
