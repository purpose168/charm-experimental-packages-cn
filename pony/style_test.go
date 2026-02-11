package pony

import (
	"image/color"
	"testing"

	"github.com/purpose168/charm-experimental-packages-cn/exp/golden"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(color.Color) bool
	}{
		{
			name:    "命名颜色红色",
			input:   "red",
			wantErr: false,
			check: func(c color.Color) bool {
				return c != nil
			},
		},
		{
			name:    "十六进制颜色",
			input:   "#FF0000",
			wantErr: false,
			check: func(c color.Color) bool {
				return c != nil
			},
		},
		{
			name:    "短十六进制颜色",
			input:   "#f00",
			wantErr: false,
			check: func(c color.Color) bool {
				return c != nil
			},
		},
		{
			name:    "RGB颜色",
			input:   "rgb(255, 0, 0)",
			wantErr: false,
			check: func(c color.Color) bool {
				if c == nil {
					return false
				}
				r, g, b, _ := c.RGBA()
				return r > g && r > b
			},
		},
		{
			name:    "ANSI颜色代码",
			input:   "196",
			wantErr: false,
			check: func(c color.Color) bool {
				return c != nil
			},
		},
		{
			name:    "明亮颜色",
			input:   "bright-red",
			wantErr: false,
			check: func(c color.Color) bool {
				return c != nil
			},
		},
		{
			name:    "无效的十六进制颜色",
			input:   "#GGGGGG",
			wantErr: true,
		},
		{
			name:    "无效的RGB颜色",
			input:   "rgb(300, 0, 0)",
			wantErr: true,
		},
		{
			name:    "未知的命名颜色",
			input:   "notacolor",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := parseColor(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil && !tt.check(c) {
				t.Errorf("parseColor() color check failed for input %q", tt.input)
			}
		})
	}
}

func TestRenderWithStyle(t *testing.T) {
	const markup = `<text font-weight="bold" foreground-color="red">Styled Text</text>`

	tmpl, err := Parse[any](markup)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output := tmpl.Render(nil, 80, 24)
	golden.RequireEqual(t, output)
}

func TestRenderBoxWithBorderStyle(t *testing.T) {
	const markup = `<box border="rounded" border-color="cyan"><text>Content</text></box>`

	tmpl, err := Parse[any](markup)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output := tmpl.Render(nil, 80, 24)
	golden.RequireEqual(t, output)
}

// 测试命名颜色覆盖范围。
func TestNamedColorsCoverage(t *testing.T) {
	colors := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"gray", "grey", "bright-black",
		"bright-red", "bright-green", "bright-yellow",
		"bright-blue", "bright-magenta", "bright-cyan", "bright-white",
	}

	for _, colorName := range colors {
		t.Run(colorName, func(t *testing.T) {
			c, err := parseColor(colorName)
			if err != nil {
				t.Errorf("parseColor(%q) error = %v", colorName, err)
			}
			if c == nil {
				t.Errorf("parseColor(%q) returned nil", colorName)
			}
		})
	}
}

// 测试ANSI颜色覆盖范围。
func TestAnsiColorsCoverage(t *testing.T) {
	// 测试不同的ANSI颜色范围
	testCodes := []int{0, 7, 8, 15, 16, 100, 231, 232, 240, 255}

	for _, code := range testCodes {
		c := ansiColor(code)
		if c == nil {
			t.Errorf("ansiColor(%d) returned nil", code)
		}
	}

	// 测试超出范围的情况
	if ansiColor(-1) != nil {
		t.Error("ansiColor(-1) should return nil")
	}
	if ansiColor(256) != nil {
		t.Error("ansiColor(256) should return nil")
	}
}
