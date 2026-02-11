package ansi

import "testing"

func TestSelectGraphicRendition(t *testing.T) {
	tests := []struct {
		name string
		args []Attr
		want string
	}{
		{
			name: "无属性",
			args: []Attr{},
			want: "\x1b[m",
		},
		{
			name: "单个基本属性",
			args: []Attr{BoldAttr},
			want: "\x1b[1m",
		},
		{
			name: "多个基本属性",
			args: []Attr{BoldAttr, ItalicAttr, UnderlineAttr},
			want: "\x1b[1;3;4m",
		},
		{
			name: "前景色",
			args: []Attr{RedForegroundColorAttr, BoldAttr},
			want: "\x1b[31;1m",
		},
		{
			name: "背景色",
			args: []Attr{BlueBackgroundColorAttr, BoldAttr},
			want: "\x1b[44;1m",
		},
		{
			name: "亮色",
			args: []Attr{BrightRedForegroundColorAttr, BrightBlueBackgroundColorAttr},
			want: "\x1b[91;104m",
		},
		{
			name: "重置属性",
			args: []Attr{ResetAttr},
			want: "\x1b[0m",
		},
		{
			name: "负属性值",
			args: []Attr{-1},
			want: "\x1b[0m",
		},
		{
			name: "自定义属性值",
			args: []Attr{99},
			want: "\x1b[99m",
		},
		{
			name: "混合已知和自定义属性",
			args: []Attr{BoldAttr, 99, ItalicAttr},
			want: "\x1b[1;99;3m",
		},
		{
			name: "所有文本装饰",
			args: []Attr{
				BoldAttr,
				FaintAttr,
				ItalicAttr,
				UnderlineAttr,
				SlowBlinkAttr,
				ReverseAttr,
				ConcealAttr,
				StrikethroughAttr,
			},
			want: "\x1b[1;2;3;4;5;7;8;9m",
		},
		{
			name: "所有颜色重置属性",
			args: []Attr{
				DefaultForegroundColorAttr,
				DefaultBackgroundColorAttr,
				DefaultUnderlineColorAttr,
			},
			want: "\x1b[39;49;59m",
		},
		{
			name: "扩展颜色属性",
			args: []Attr{
				ExtendedForegroundColorAttr,
				ExtendedBackgroundColorAttr,
				ExtendedUnderlineColorAttr,
			},
			want: "\x1b[38;48;58m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectGraphicRendition(tt.args...); got != tt.want {
				t.Errorf("SelectGraphicRendition() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSGR(t *testing.T) {
	// Test that SGR is an alias for SelectGraphicRendition
	tests := []struct {
		name string
		args []Attr
	}{
		{
			name: "empty args",
			args: []Attr{},
		},
		{
			name: "single arg",
			args: []Attr{BoldAttr},
		},
		{
			name: "multiple args",
			args: []Attr{BoldAttr, RedForegroundColorAttr, BlueBackgroundColorAttr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SGR(tt.args...)
			want := SelectGraphicRendition(tt.args...)
			if got != want {
				t.Errorf("SGR() = %q, want %q", got, want)
			}
		})
	}
}
