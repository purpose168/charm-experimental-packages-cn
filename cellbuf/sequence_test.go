package cellbuf

import (
	"image/color"
	"testing"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// TestReadStyleColor 测试读取样式颜色的功能
func TestReadStyleColor(t *testing.T) {
	tests := []struct {
		name      string        // 测试用例名称
		params    []ansi.Param  // ANSI 参数
		wantN     int           // 期望读取的参数数量
		wantColor color.Color   // 期望的颜色
		wantNil   bool          // 期望颜色为 nil
	}{
		{
			name:    "无效 - 参数太少",
			params:  []ansi.Param{38},
			wantN:   0,
			wantNil: true,
		},
		{
			name:    "实现定义",
			params:  []ansi.Param{38, 0},
			wantN:   2,
			wantNil: true,
		},
		{
			name:      "透明",
			params:    []ansi.Param{38, 1},
			wantN:     2,
			wantColor: color.Transparent,
		},
		{
			name:      "RGB 分号分隔",
			params:    []ansi.Param{38, 2, 100, 150, 200},
			wantN:     5,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 255},
		},
		{
			name: "RGB 冒号分隔",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200,
			},
			wantN:     5,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 255},
		},
		{
			name: "带色彩空间的 RGB",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				1 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200,
			},
			wantN:     6,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 255},
		},
		// {
		// 	name:      "CMY 分号分隔",
		// 	params:    []ansi.Parameter{38, 3, 100, 150, 200},
		// 	wantN:     5,
		// 	wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 0},
		// },
		{
			name: "带色彩空间的 CMY",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				3 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200,
			},
			wantN:     6,
			wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 0},
		},
		// {
		// 	name: "CMY 冒号分隔",
		// 	params: []ansi.Parameter{
		// 		38 | parser.HasMoreFlag,
		// 		3 | parser.HasMoreFlag,
		// 		100 | parser.HasMoreFlag,
		// 		150 | parser.HasMoreFlag,
		// 		200,
		// 	},
		// 	wantN:     5,
		// 	wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 0},
		// },
		// {
		// 	name:      "CMYK 分号分隔",
		// 	params:    []ansi.Parameter{38, 4, 100, 150, 200, 50},
		// 	wantN:     6,
		// 	wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 50},
		// },
		{
			name: "带色彩空间的 CMYK",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				4 | parser.HasMoreFlag,
				1 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200 | parser.HasMoreFlag,
				50,
			},
			wantN:     7,
			wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 50},
		},
		// {
		// 	name: "CMYK 冒号分隔",
		// 	params: []ansi.Parameter{
		// 		38 | parser.HasMoreFlag,
		// 		4 | parser.HasMoreFlag,
		// 		100 | parser.HasMoreFlag,
		// 		150 | parser.HasMoreFlag,
		// 		200 | parser.HasMoreFlag,
		// 		50,
		// 	},
		// 	wantN:     6,
		// 	wantColor: color.CMYK{C: 100, M: 150, Y: 200, K: 50},
		// },
		{
			name:      "索引颜色分号分隔",
			params:    []ansi.Param{38, 5, 123},
			wantN:     3,
			wantColor: ansi.ExtendedColor(123),
		},
		{
			name: "索引颜色冒号分隔",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				5 | parser.HasMoreFlag,
				123,
			},
			wantN:     3,
			wantColor: ansi.ExtendedColor(123),
		},
		{
			name:    "无效的颜色类型",
			params:  []ansi.Param{38, 99},
			wantN:   0,
			wantNil: true,
		},
		{
			name: "带容差和色彩空间的 RGB",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				1 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200 | parser.HasMoreFlag,
				0 | parser.HasMoreFlag, // 容差值
				1,                      // 容差色彩空间
			},
			wantN:     8,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 255},
		},
		// 无效情况
		{
			name:    "空参数",
			params:  []ansi.Param{},
			wantN:   0,
			wantNil: true,
		},
		{
			name:    "单个参数",
			params:  []ansi.Param{38},
			wantN:   0,
			wantNil: true,
		},
		{
			name:    "nil 参数",
			params:  nil,
			wantN:   0,
			wantNil: true,
		},
		// 混合分隔符情况（应失败）
		{
			name: "RGB 混合分隔符",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2,                        // 分号
				100 | parser.HasMoreFlag, // 冒号
				150,                      // 分号
				200,
			},
			wantN:   0,
			wantNil: true,
		},
		{
			name: "CMYK 混合分隔符",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				4,                        // 分号
				100 | parser.HasMoreFlag, // 冒号
				150,                      // 分号
				200 | parser.HasMoreFlag, // 冒号
				50,
			},
			wantN:   0,
			wantNil: true,
		},
		// 边界情况
		{
			name: "RGB 最大值",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				255 | parser.HasMoreFlag,
				255 | parser.HasMoreFlag,
				255,
			},
			wantN:     5,
			wantColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name: "RGB 负值",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				-1 | parser.HasMoreFlag,
				-1 | parser.HasMoreFlag,
				-1,
			},
			wantN:   0,
			wantNil: true,
		},
		{
			name: "索引超出范围的索引颜色",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				5 | parser.HasMoreFlag,
				256, // 超出范围
			},
			wantN:     3,
			wantColor: ansi.ExtendedColor(0),
		},
		{
			name: "负索引的索引颜色",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				5 | parser.HasMoreFlag,
				-1,
			},
			wantN:   0,
			wantNil: true,
		},
		{
			name: "RGB 参数截断",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				2 | parser.HasMoreFlag,
				100 | parser.HasMoreFlag,
				150,
			},
			wantN:   0,
			wantNil: true,
		},
		{
			name: "CMYK 参数截断",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				4 | parser.HasMoreFlag,
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200,
			},
			wantN:   0,
			wantNil: true,
		},
		// RGBA (类型 6) 测试用例
		// {
		// 	name:      "RGBA 分号分隔",
		// 	params:    []Parameter{38, 6, 100, 150, 200, 128},
		// 	wantN:     6,
		// 	wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 128},
		// },
		// {
		// 	name: "RGBA 冒号分隔",
		// 	params: []ansi.Parameter{
		// 		38 | parser.HasMoreFlag,
		// 		6 | parser.HasMoreFlag,
		// 		100 | parser.HasMoreFlag,
		// 		150 | parser.HasMoreFlag,
		// 		200 | parser.HasMoreFlag,
		// 		128,
		// 	},
		// 	wantN:     6,
		// 	wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 128},
		// },
		{
			name: "带色彩空间的 RGBA",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				6 | parser.HasMoreFlag,
				1 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200 | parser.HasMoreFlag,
				128,
			},
			wantN:     7,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 128},
		},
		{
			name: "带容差和色彩空间的 RGBA",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				6 | parser.HasMoreFlag,
				1 | parser.HasMoreFlag, // 色彩空间 ID
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200 | parser.HasMoreFlag,
				128 | parser.HasMoreFlag,
				0 | parser.HasMoreFlag, // 容差值
				1,                      // 容差色彩空间
			},
			wantN:     9,
			wantColor: color.RGBA{R: 100, G: 150, B: 200, A: 128},
		},
		{
			name: "RGBA 最大值",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				6 | parser.HasMoreFlag,
				0 | parser.HasMoreFlag, // 色彩空间 ID
				255 | parser.HasMoreFlag,
				255 | parser.HasMoreFlag,
				255 | parser.HasMoreFlag,
				255,
			},
			wantN:     7,
			wantColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name: "RGBA 参数截断",
			params: []ansi.Param{
				38 | parser.HasMoreFlag,
				6 | parser.HasMoreFlag,
				100 | parser.HasMoreFlag,
				150 | parser.HasMoreFlag,
				200,
			},
			wantN:   0,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotColor color.Color
			gotN := ReadStyleColor(tt.params, &gotColor)
			if gotN != tt.wantN {
				t.Errorf("ReadColor() gotN = %v, want %v", gotN, tt.wantN)
			}
			if tt.wantNil {
				if gotColor != nil {
					t.Errorf("ReadColor() gotColor = %v, want nil", gotColor)
				}
				return
			}
			if gotColor != tt.wantColor {
				t.Errorf("ReadColor() gotColor = %v, want %v", gotColor, tt.wantColor)
			}
		})
	}
}
