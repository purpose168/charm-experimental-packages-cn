package ansi

import (
	"testing"
)

// TestModeSetting_Methods 测试ModeSetting的方法

func TestModeSetting_Methods(t *testing.T) {
	tests := []struct {
		name     string      // 测试用例名称
		mode     ModeSetting // 模式设置
		notRecog bool        // 是否未识别
		isSet    bool        // 是否设置
		isReset  bool        // 是否重置
		permSet  bool        // 是否永久设置
		permRst  bool        // 是否永久重置
	}{
		{
			name:     "模式未识别",
			mode:     ModeNotRecognized,
			notRecog: true,
			isSet:    false,
			isReset:  false,
			permSet:  false,
			permRst:  false,
		},
		{
			name:     "模式设置",
			mode:     ModeSet,
			notRecog: false,
			isSet:    true,
			isReset:  false,
			permSet:  false,
			permRst:  false,
		},
		{
			name:     "模式重置",
			mode:     ModeReset,
			notRecog: false,
			isSet:    false,
			isReset:  true,
			permSet:  false,
			permRst:  false,
		},
		{
			name:     "模式永久设置",
			mode:     ModePermanentlySet,
			notRecog: false,
			isSet:    true,
			isReset:  false,
			permSet:  true,
			permRst:  false,
		},
		{
			name:     "模式永久重置",
			mode:     ModePermanentlyReset,
			notRecog: false,
			isSet:    false,
			isReset:  true,
			permSet:  false,
			permRst:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.IsNotRecognized(); got != tt.notRecog {
				t.Errorf("IsNotRecognized() = %v, want %v", got, tt.notRecog)
			}
			if got := tt.mode.IsSet(); got != tt.isSet {
				t.Errorf("IsSet() = %v, want %v", got, tt.isSet)
			}
			if got := tt.mode.IsReset(); got != tt.isReset {
				t.Errorf("IsReset() = %v, want %v", got, tt.isReset)
			}
			if got := tt.mode.IsPermanentlySet(); got != tt.permSet {
				t.Errorf("IsPermanentlySet() = %v, want %v", got, tt.permSet)
			}
			if got := tt.mode.IsPermanentlyReset(); got != tt.permRst {
				t.Errorf("IsPermanentlyReset() = %v, want %v", got, tt.permRst)
			}
		})
	}
}

// TestSetMode 测试SetMode函数

func TestSetMode(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		modes    []Mode // 模式列表
		expected string // 期望结果
	}{
		{
			name:     "空模式",
			modes:    []Mode{},
			expected: "",
		},
		{
			name:     "单个ANSI模式",
			modes:    []Mode{ModeKeyboardAction},
			expected: "\x1b[2h",
		},
		{
			name:     "单个DEC模式",
			modes:    []Mode{ModeCursorKeys},
			expected: "\x1b[?1h",
		},
		{
			name:     "多个ANSI模式",
			modes:    []Mode{ModeKeyboardAction, ModeInsertReplace},
			expected: "\x1b[2;4h",
		},
		{
			name:     "多个DEC模式",
			modes:    []Mode{ModeCursorKeys, ModeAutoWrap},
			expected: "\x1b[?1;7h",
		},
		{
			name:     "混合ANSI和DEC模式",
			modes:    []Mode{ModeKeyboardAction, ModeCursorKeys},
			expected: "\x1b[2h\x1b[?1h",
		},
		{
			name:     "多个混合ANSI和DEC模式",
			modes:    []Mode{ModeKeyboardAction, ModeInsertReplace, ModeCursorKeys, ModeAutoWrap},
			expected: "\x1b[2;4h\x1b[?1;7h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetMode(tt.modes...); got != tt.expected {
				t.Errorf("SetMode() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestResetMode 测试ResetMode函数

func TestResetMode(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		modes    []Mode // 模式列表
		expected string // 期望结果
	}{
		{
			name:     "空模式",
			modes:    []Mode{},
			expected: "",
		},
		{
			name:     "单个ANSI模式",
			modes:    []Mode{ModeKeyboardAction},
			expected: "\x1b[2l",
		},
		{
			name:     "单个DEC模式",
			modes:    []Mode{ModeCursorKeys},
			expected: "\x1b[?1l",
		},
		{
			name:     "多个ANSI模式",
			modes:    []Mode{ModeKeyboardAction, ModeInsertReplace},
			expected: "\x1b[2;4l",
		},
		{
			name:     "多个DEC模式",
			modes:    []Mode{ModeCursorKeys, ModeAutoWrap},
			expected: "\x1b[?1;7l",
		},
		{
			name:     "混合ANSI和DEC模式",
			modes:    []Mode{ModeKeyboardAction, ModeCursorKeys},
			expected: "\x1b[2l\x1b[?1l",
		},
		{
			name:     "多个混合ANSI和DEC模式",
			modes:    []Mode{ModeKeyboardAction, ModeInsertReplace, ModeCursorKeys, ModeAutoWrap},
			expected: "\x1b[2;4l\x1b[?1;7l",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResetMode(tt.modes...); got != tt.expected {
				t.Errorf("ResetMode() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestRequestMode 测试RequestMode函数

func TestRequestMode(t *testing.T) {
	tests := []struct {
		name     string // 测试用例名称
		mode     Mode   // 模式
		expected string // 期望结果
	}{
		{
			name:     "ANSI模式",
			mode:     ModeKeyboardAction,
			expected: "\x1b[2$p",
		},
		{
			name:     "DEC mode",
			mode:     ModeCursorKeys,
			expected: "\x1b[?1$p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RequestMode(tt.mode); got != tt.expected {
				t.Errorf("RequestMode() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestReportMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     Mode
		value    ModeSetting
		expected string
	}{
		{
			name:     "ANSI mode not recognized",
			mode:     ModeKeyboardAction,
			value:    ModeNotRecognized,
			expected: "\x1b[2;0$y",
		},
		{
			name:     "DEC mode set",
			mode:     ModeCursorKeys,
			value:    ModeSet,
			expected: "\x1b[?1;1$y",
		},
		{
			name:     "ANSI mode reset",
			mode:     ModeInsertReplace,
			value:    ModeReset,
			expected: "\x1b[4;2$y",
		},
		{
			name:     "DEC mode permanently set",
			mode:     ModeAutoWrap,
			value:    ModePermanentlySet,
			expected: "\x1b[?7;3$y",
		},
		{
			name:     "ANSI mode permanently reset",
			mode:     ModeSendReceive,
			value:    ModePermanentlyReset,
			expected: "\x1b[12;4$y",
		},
		{
			name:     "Invalid mode setting defaults to not recognized",
			mode:     ModeKeyboardAction,
			value:    5,
			expected: "\x1b[2;0$y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReportMode(tt.mode, tt.value); got != tt.expected {
				t.Errorf("ReportMode() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestModeImplementations(t *testing.T) {
	tests := []struct {
		name     string
		mode     Mode
		expected int
	}{
		{
			name:     "ANSIMode",
			mode:     ANSIMode(42),
			expected: 42,
		},
		{
			name:     "DECMode",
			mode:     DECMode(99),
			expected: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.Mode(); got != tt.expected {
				t.Errorf("Mode() = %v, want %v", got, tt.expected)
			}
		})
	}
}
