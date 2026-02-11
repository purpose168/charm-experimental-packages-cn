package cellbuf

import (
	"testing"
)

func TestNewCell(t *testing.T) {
	tests := []struct {
		name     string
		mainRune rune
		combRune []rune
		want     *Cell
	}{{
		name:     "简单 ASCII 字符",
		mainRune: 'a',
		want:     &Cell{Rune: 'a', Width: 1},
	},
		{
			name:     "宽字符",
			mainRune: '世',
			want:     &Cell{Rune: '世', Width: 2},
		},
		{
			name:     "组合字符",
			mainRune: 'e',
			combRune: []rune{'́'}, // 重音符号
			want:     &Cell{Rune: 'e', Comb: []rune{'́'}, Width: 1},
		},
		{
			name:     "零宽度字符",
			mainRune: '\u200B', // 零宽度空格
			want:     &Cell{Rune: '\u200B', Width: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCell(tt.mainRune, tt.combRune...)
			if got.Rune != tt.want.Rune {
				t.Errorf("NewCell().Rune = %v, want %v", got.Rune, tt.want.Rune)
			}
			if got.Width != tt.want.Width {
				t.Errorf("NewCell().Width = %v, want %v", got.Width, tt.want.Width)
			}
			if len(got.Comb) != len(tt.want.Comb) {
				t.Errorf("NewCell().Comb length = %v, want %v", len(got.Comb), len(tt.want.Comb))
			}
		})
	}
}

func TestNewCellString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want *Cell
	}{{
		name: "空字符串",
		str:  "",
		want: &Cell{Width: 0},
	},
		{
			name: "简单 ASCII 字符",
			str:  "a",
			want: &Cell{Rune: 'a', Width: 1},
		},
		{
			name: "组合字符",
			str:  "é", // 带重音符号的 e
			want: &Cell{Rune: 'é', Width: 1},
		},
		{
			name: "宽字符",
			str:  "世",
			want: &Cell{Rune: '世', Width: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCellString(tt.str)
			if got.Rune != tt.want.Rune {
				t.Errorf("NewCellString().Rune = %v, want %v", got.Rune, tt.want.Rune)
			}
			if got.Width != tt.want.Width {
				t.Errorf("NewCellString().Width = %v, want %v", got.Width, tt.want.Width)
			}
		})
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		name      string
		line      Line
		wantStr   string
		wantLen   int
		wantWidth int
	}{{
		name:      "空行",
		line:      Line{},
		wantStr:   "",
		wantLen:   0,
		wantWidth: 0,
	},
		{
			name:      "简单行",
			line:      Line{NewCell('a'), NewCell('b'), NewCell('c')},
			wantStr:   "abc",
			wantLen:   3,
			wantWidth: 3,
		},
		{
			name:      "包含 nil 单元格的行",
			line:      Line{nil, NewCell('a'), nil},
			wantStr:   " a",
			wantLen:   3,
			wantWidth: 3,
		},
		{
			name:      "包含宽字符的行",
			line:      Line{NewCell('世'), NewCell('界')},
			wantStr:   "世界",
			wantLen:   2,
			wantWidth: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.String(); got != tt.wantStr {
				t.Errorf("Line.String() = %q, want %q", got, tt.wantStr)
			}
			if got := tt.line.Len(); got != tt.wantLen {
				t.Errorf("Line.Len() = %v, want %v", got, tt.wantLen)
			}
			if got := tt.line.Width(); got != tt.wantWidth {
				t.Errorf("Line.Width() = %v, want %v", got, tt.wantWidth)
			}
		})
	}
}

func TestBuffer(t *testing.T) {
	t.Run("创建和调整大小", func(t *testing.T) {
		b := NewBuffer(3, 2)
		if b.Width() != 3 {
			t.Errorf("Buffer width = %d, want 3", b.Width())
		}
		if b.Height() != 2 {
			t.Errorf("Buffer height = %d, want 2", b.Height())
		}

		b.Resize(4, 3)
		if b.Width() != 4 {
			t.Errorf("调整大小后，缓冲区宽度 = %d, 期望 4", b.Width())
		}
		if b.Height() != 3 {
			t.Errorf("调整大小后，缓冲区高度 = %d, 期望 3", b.Height())
		}
	})

	t.Run("单元格操作", func(t *testing.T) {
		b := NewBuffer(3, 3)
		cell := NewCell('A')

		b.SetCell(1, 1, cell)
		got := b.Cell(1, 1)
		if got.Rune != 'A' {
			t.Errorf("设置单元格后，得到的字符为 %c, 期望 A", got.Rune)
		}
	})

	t.Run("清除操作", func(t *testing.T) {
		b := NewBuffer(2, 2)
		b.SetCell(0, 0, NewCell('A'))
		b.SetCell(1, 0, NewCell('B'))
		b.Clear()

		// if b.Cell(0, 0) != nil {
		// TODO: 我们应该返回 nil 而不是 BlankCell 吗？nil 表示默认单元格。

		if !b.Cell(0, 0).Equal(&BlankCell) {
			t.Error("清除后，单元格应为 nil")
		}
	})

	t.Run("插入行", func(t *testing.T) {
		b := NewBuffer(3, 3)
		b.SetCell(0, 0, NewCell('A'))
		b.SetCell(0, 1, NewCell('B'))

		b.InsertLine(1, 1, nil)
		got := b.Cell(0, 1)
		if !got.Equal(&BlankCell) {
			t.Error("插入行后，插入的行应为空")
		}
	})

	t.Run("删除行", func(t *testing.T) {
		b := NewBuffer(3, 3)
		b.SetCell(0, 0, NewCell('A'))
		b.SetCell(0, 1, NewCell('B'))

		b.DeleteLine(0, 1, nil)
		got := b.Cell(0, 0)
		if !got.Equal(NewCell('B')) {
			t.Error("删除行后，第一行应为空")
		}
	})
}

func TestBufferBounds(t *testing.T) {
	b := NewBuffer(4, 3)
	bounds := b.Bounds()

	if bounds.Min.X != 0 || bounds.Min.Y != 0 {
		t.Errorf("Buffer bounds min = (%d,%d), want (0,0)", bounds.Min.X, bounds.Min.Y)
	}
	if bounds.Max.X != 4 || bounds.Max.Y != 3 {
		t.Errorf("Buffer bounds max = (%d,%d), want (4,3)", bounds.Max.X, bounds.Max.Y)
	}
}
