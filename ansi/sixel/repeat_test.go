package sixel

import (
	"bytes"
	"testing"
)

func TestWriteRepeat(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		char     byte
		expected string
		wantErr  bool
	}{
		{
			name:     "basic repeat",
			count:    3,
			char:     'A',
			expected: "!3A",
		},
		{
			name:     "single digit",
			count:    5,
			char:     '#',
			expected: "!5#",
		},
		{
			name:     "multiple digits",
			count:    123,
			char:     'x',
			expected: "!123x",
		},
		{
			name:     "zero count",
			count:    0,
			char:     'B',
			expected: "!0B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			n, err := WriteRepeat(buf, tt.count, tt.char)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteRepeat() 错误 = %v，期望错误 %v", err, tt.wantErr)
				return
			}
			if got := buf.String(); got != tt.expected {
				t.Errorf("WriteRepeat() = %v，期望 %v", got, tt.expected)
			}
			if n != len(tt.expected) {
				t.Errorf("WriteRepeat() 返回长度 = %v，期望 %v", n, len(tt.expected))
			}
		})
	}
}

func TestDecodeRepeat(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		wantRepeat  Repeat
		wantN       int
		description string
	}{
		{
			name:        "basic repeat",
			input:       []byte("!3A"),
			wantRepeat:  Repeat{Count: 3, Char: 'A'},
			wantN:       3,
			description: "简单的单位数重复",
		},
		{
			name:        "multiple digits",
			input:       []byte("!123x"),
			wantRepeat:  Repeat{Count: 123, Char: 'x'},
			wantN:       5,
			description: "多位数重复",
		},
		{
			name:        "empty input",
			input:       []byte{},
			wantRepeat:  Repeat{},
			wantN:       0,
			description: "空输入应返回零值",
		},
		{
			name:        "invalid introducer",
			input:       []byte("X3A"),
			wantRepeat:  Repeat{},
			wantN:       0,
			description: "没有正确引导符的输入",
		},
		{
			name:        "incomplete sequence",
			input:       []byte("!3"),
			wantRepeat:  Repeat{},
			wantN:       0,
			description: "没有字符的不完整序列",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepeat, gotN := DecodeRepeat(tt.input)
			if gotRepeat != tt.wantRepeat {
				t.Errorf("DecodeRepeat() 得到重复 = %v，期望 %v", gotRepeat, tt.wantRepeat)
			}
			if gotN != tt.wantN {
				t.Errorf("DecodeRepeat() 得到 N = %v，期望 %v", gotN, tt.wantN)
			}
		})
	}
}

func TestRepeat_String(t *testing.T) {
	tests := []struct {
		name     string
		repeat   Repeat
		expected string
	}{
		{
			name:     "basic repeat",
			repeat:   Repeat{Count: 3, Char: 'A'},
			expected: "!3A",
		},
		{
			name:     "multiple digits",
			repeat:   Repeat{Count: 123, Char: 'x'},
			expected: "!123x",
		},
		{
			name:     "zero count",
			repeat:   Repeat{Count: 0, Char: 'B'},
			expected: "!0B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.repeat.String(); got != tt.expected {
				t.Errorf("Repeat.String() = %v，期望 %v", got, tt.expected)
			}
		})
	}
}
