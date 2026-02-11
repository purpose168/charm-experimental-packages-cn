package ansi

import "testing"

// TestNotify 测试Notify函数

func TestNotify(t *testing.T) {
	tests := []struct {
		name string // 测试用例名称
		s    string // 输入字符串
		want string // 期望结果
	}{
		{
			name: "基本",
			s:    "Hello, World!",
			want: "\x1b]9;Hello, World!\x07",
		},
		{
			name: "空字符串",
			s:    "",
			want: "\x1b]9;\x07",
		},
		{
			name: "特殊字符",
			s:    "Line1\nLine2\tTabbed",
			want: "\x1b]9;Line1\nLine2\tTabbed\x07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Notify(tt.s); got != tt.want {
				t.Errorf("Notify() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestDesktopNotification 测试DesktopNotification函数

func TestDesktopNotification(t *testing.T) {
	tests := []struct {
		name     string   // 测试用例名称
		payload  string   // 通知内容
		metadata []string // 元数据
		want     string   // 期望结果
	}{
		{
			name:     "基本",
			payload:  "Task Completed",
			metadata: []string{},
			want:     "\x1b]99;;Task Completed\x07",
		},
		{
			name:     "带元数据",
			payload:  "New Message",
			metadata: []string{"i=1", "a=focus"},
			want:     "\x1b]99;i=1:a=focus;New Message\x07",
		},
		{
			name:     "空内容",
			payload:  "",
			metadata: []string{"i=2"},
			want:     "\x1b]99;i=2;\x07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DesktopNotification(tt.payload, tt.metadata...); got != tt.want {
				t.Errorf("DesktopNotification() = %q, want %q", got, tt.want)
			}
		})
	}
}
