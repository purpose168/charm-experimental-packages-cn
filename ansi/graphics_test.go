package ansi

import (
	"testing"
)

func TestKittyGraphics(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		opts    []string
		want    string
	}{
		{
			name:    "空负载无选项",
			payload: []byte{},
			opts:    nil,
			want:    "\x1b_G\x1b\\",
		},
		{
			name:    "带负载无选项",
			payload: []byte("test"),
			opts:    nil,
			want:    "\x1b_G;test\x1b\\",
		},
		{
			name:    "带负载和选项",
			payload: []byte("test"),
			opts:    []string{"a=t", "f=100"},
			want:    "\x1b_Ga=t,f=100;test\x1b\\",
		},
		{
			name:    "多选项无负载",
			payload: []byte{},
			opts:    []string{"q=2", "C=1", "f=24"},
			want:    "\x1b_Gq=2,C=1,f=24\x1b\\",
		},
		{
			name:    "负载中有特殊字符",
			payload: []byte("\x1b_G"),
			opts:    []string{"a=t"},
			want:    "\x1b_Ga=t;\x1b_G\x1b\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KittyGraphics(tt.payload, tt.opts...)
			if got != tt.want {
				t.Errorf("KittyGraphics() = %q，期望 %q", got, tt.want)
			}
		})
	}
}
