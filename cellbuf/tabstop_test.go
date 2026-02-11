package cellbuf

import (
	"testing"
)

// TestTabStops 测试制表位功能
func TestTabStops(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		interval int
		checks   []struct {
			col      int
			expected bool
		}
	}{
		{
			name:     "默认间隔为8",
			width:    24,
			interval: DefaultTabInterval,
			checks: []struct {
				col      int
				expected bool
			}{
				{0, true},   // 第一个制表位
				{7, false},  // 不是制表位
				{8, true},   // 第二个制表位
				{15, false}, // 不是制表位
				{16, true},  // 第三个制表位
				{23, false}, // 不是制表位
			},
		},
		{
			name:     "自定义间隔为4",
			width:    16,
			interval: 4,
			checks: []struct {
				col      int
				expected bool
			}{
				{0, true},   // 第一个制表位
				{3, false},  // 不是制表位
				{4, true},   // 第二个制表位
				{7, false},  // 不是制表位
				{8, true},   // 第三个制表位
				{12, true},  // 第四个制表位
				{15, false}, // 不是制表位
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTabStops(tt.width, tt.interval)

			// 测试初始制表位
			for _, check := range tt.checks {
				if got := ts.IsStop(check.col); got != check.expected {
					t.Errorf("IsStop(%d) = %v, want %v", check.col, got, check.expected)
				}
			}

			// 测试设置自定义制表位
			customCol := tt.interval + 1
			ts.Set(customCol)
			if !ts.IsStop(customCol) {
				t.Errorf("After Set(%d), IsStop(%d) = false, want true", customCol, customCol)
			}

			// 测试重置制表位
			regularStop := tt.interval
			ts.Reset(regularStop)
			if ts.IsStop(regularStop) {
				t.Errorf("After Reset(%d), IsStop(%d) = true, want false", regularStop, regularStop)
			}
		})
	}
}

// TestTabStopsNavigation 测试制表位导航功能
func TestTabStopsNavigation(t *testing.T) {
	ts := NewTabStops(24, DefaultTabInterval)

	tests := []struct {
		name     string
		col      int
		wantNext int
		wantPrev int
	}{
		{
			name:     "从列0开始",
			col:      0,
			wantNext: 8,
			wantPrev: 0,
		},
		{
			name:     "从列4开始",
			col:      4,
			wantNext: 8,
			wantPrev: 0,
		},
		{
			name:     "从列8开始",
			col:      8,
			wantNext: 16,
			wantPrev: 0,
		},
		{
			name:     "从列20开始",
			col:      20,
			wantNext: 23,
			wantPrev: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ts.Next(tt.col); got != tt.wantNext {
				t.Errorf("Next(%d) = %v, want %v", tt.col, got, tt.wantNext)
			}
			if got := ts.Prev(tt.col); got != tt.wantPrev {
				t.Errorf("Prev(%d) = %v, want %v", tt.col, got, tt.wantPrev)
			}
		})
	}
}

// TestTabStopsClear 测试清除所有制表位
func TestTabStopsClear(t *testing.T) {
	ts := NewTabStops(24, DefaultTabInterval)

	// 验证初始状态
	if !ts.IsStop(0) || !ts.IsStop(8) || !ts.IsStop(16) {
		t.Error("初始制表位设置不正确")
	}

	// 清除所有制表位
	ts.Clear()

	// 验证所有制表位已清除
	for i := range 24 {
		if ts.IsStop(i) {
			t.Errorf("Clear()后列%d的制表位仍然设置", i)
		}
	}
}

// TestTabStopsResize 测试调整缓冲区大小时的制表位行为
func TestTabStopsResize(t *testing.T) {
	tests := []struct {
		name        string
		initialSize int
		newSize     int
		checks      []struct {
			col      int
			expected bool
		}
	}{
		{
			name:        "增大缓冲区",
			initialSize: 16,
			newSize:     24,
			checks: []struct {
				col      int
				expected bool
			}{
				{0, true},   // 原始制表位
				{8, true},   // 原始制表位
				{16, true},  // 新制表位
				{23, false}, // 不是制表位
			},
		},
		{
			name:        "相同大小 - 无变化",
			initialSize: 16,
			newSize:     16,
			checks: []struct {
				col      int
				expected bool
			}{
				{0, true},   // 原始制表位
				{8, true},   // 原始制表位
				{15, false}, // 不是制表位
			},
		},
		{
			name:        "使用自定义间隔调整大小",
			initialSize: 8,
			newSize:     16,
			checks: []struct {
				col      int
				expected bool
			}{
				{0, true},   // 第一个制表位
				{4, true},   // 第二个制表位
				{8, true},   // 第三个制表位
				{12, true},  // 第四个制表位
				{15, false}, // 不是制表位
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts *TabStops
			if tt.name == "使用自定义间隔调整大小" {
				ts = NewTabStops(tt.initialSize, 4) // 自定义间隔为4
			} else {
				ts = DefaultTabStops(tt.initialSize)
			}

			// 验证初始状态
			if ts.width != tt.initialSize {
				t.Errorf("初始宽度 = %d, 期望 %d", ts.width, tt.initialSize)
			}

			// 执行调整大小
			ts.Resize(tt.newSize)

			// 验证新大小
			if ts.width != tt.newSize {
				t.Errorf("调整大小后，宽度 = %d, 期望 %d", ts.width, tt.newSize)
			}

			// 检查调整大小后的制表位
			for _, check := range tt.checks {
				if got := ts.IsStop(check.col); got != check.expected {
					t.Errorf("调整大小后，IsStop(%d) = %v, 期望 %v",
						check.col, got, check.expected)
				}
			}

			// 验证stops切片长度正确
			expectedStopsLen := (tt.newSize + (ts.interval - 1)) / ts.interval
			if len(ts.stops) != expectedStopsLen {
				t.Errorf("stops切片长度 = %d, 期望 %d",
					len(ts.stops), expectedStopsLen)
			}
		})
	}
}

// TestTabStopsResizeEdgeCases 测试调整大小时的边界情况
func TestTabStopsResizeEdgeCases(t *testing.T) {
	t.Run("调整到零大小", func(t *testing.T) {
		ts := DefaultTabStops(8)
		ts.Resize(0)

		if ts.width != 0 {
			t.Errorf("宽度 = %d, 期望 0", ts.width)
		}

		// 验证没有制表位可访问
		if ts.IsStop(0) {
			t.Error("对于零宽度，IsStop(0) 应返回 false")
		}
	})

	t.Run("调整到非常大的宽度", func(t *testing.T) {
		ts := DefaultTabStops(8)
		largeWidth := 1000
		ts.Resize(largeWidth)

		// 检查较高位置的一些制表位
		checks := []struct {
			col      int
			expected bool
		}{
			{992, true},  // 8的倍数
			{999, false}, // 不是制表位
		}

		for _, check := range checks {
			if got := ts.IsStop(check.col); got != check.expected {
				t.Errorf("IsStop(%d) = %v, 期望 %v",
					check.col, got, check.expected)
			}
		}
	})

	t.Run("多次调整大小", func(t *testing.T) {
		ts := DefaultTabStops(8)

		// 执行多次调整大小
		sizes := []int{16, 8, 24, 4}
		for _, size := range sizes {
			ts.Resize(size)

			// 验证每次调整大小后的基本属性
			if ts.width != size {
				t.Errorf("宽度 = %d, 期望 %d", ts.width, size)
			}

			// 检查第一个制表位始终设置
			if !ts.IsStop(0) {
				t.Errorf("调整大小到 %d 后，IsStop(0) = false, 期望 true", size)
			}
		}
	})
}
