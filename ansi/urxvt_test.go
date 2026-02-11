package ansi

import "testing"

// TestUrxvtExt 测试 URxvtExt 函数的功能
func TestUrxvtExt(t *testing.T) {
	tests := []struct {
		extension string   // 扩展名称
		params    []string // 参数列表
		expected  string   // 预期结果
	}{
		{
			extension: "foo",
			params:    []string{"bar", "baz"},
			expected:  "\x1b]777;foo;bar;baz\x07",
		},
		{
			extension: "test",
			params:    []string{},
			expected:  "\x1b]777;test;\x07",
		},
		{
			extension: "example",
			params:    []string{"param1"},
			expected:  "\x1b]777;example;param1\x07",
		},
		{
			extension: "notify",
			params:    []string{"message", "info"},
			expected:  "\x1b]777;notify;message;info\x07",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable to avoid loop variable reuse issue
		t.Run(tt.extension, func(t *testing.T) {
			result := URxvtExt(tt.extension, tt.params...)
			if result != tt.expected {
				t.Errorf("URxvtExt(%q, %v) = %q; want %q", tt.extension, tt.params, result, tt.expected)
			}
		})
	}
}
