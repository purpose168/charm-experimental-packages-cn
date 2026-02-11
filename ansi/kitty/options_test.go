package kitty

import (
	"reflect"
	"slices"
	"sort"
	"testing"
)

func TestOptions_Options(t *testing.T) {
	tests := []struct {
		name     string
		options  Options
		expected []string
	}{
		{
			name:     "默认选项",
			options:  Options{},
			expected: []string{}, // 默认值不生成选项
		},
		{
			name: "基本传输选项",
			options: Options{
				Format: PNG,
				ID:     1,
				Action: TransmitAndPut,
			},
			expected: []string{
				"f=100",
				"i=1",
				"a=T",
			},
		},
		{
			name: "显示选项",
			options: Options{
				X:      100,
				Y:      200,
				Z:      3,
				Width:  400,
				Height: 300,
			},
			expected: []string{
				"x=100",
				"y=200",
				"z=3",
				"w=400",
				"h=300",
			},
		},
		{
			name: "压缩和分块",
			options: Options{
				Compression: Zlib,
				Chunk:       true,
				Size:        1024,
			},
			expected: []string{
				"S=1024",
				"o=z",
			},
		},
		{
			name: "删除选项",
			options: Options{
				Delete:          DeleteID,
				DeleteResources: true,
			},
			expected: []string{
				"d=I", // 由于 DeleteResources 为 true 而转为大写
			},
		},
		{
			name: "虚拟放置",
			options: Options{
				VirtualPlacement:  true,
				ParentID:          5,
				ParentPlacementID: 2,
			},
			expected: []string{
				"U=1",
				"P=5",
				"Q=2",
			},
		},
		{
			name: "单元格定位",
			options: Options{
				OffsetX: 10,
				OffsetY: 20,
				Columns: 80,
				Rows:    24,
			},
			expected: []string{
				"X=10",
				"Y=20",
				"c=80",
				"r=24",
			},
		},
		{
			name: "传输详情",
			options: Options{
				Transmission: File,
				File:         "/tmp/image.png",
				Offset:       100,
				Number:       2,
				PlacementID:  3,
			},
			expected: []string{
				"p=3",
				"I=2",
				"t=f",
				"O=100",
			},
		},
		{
			name: "安静模式和格式",
			options: Options{
				Quite:  2,
				Format: RGB,
			},
			expected: []string{
				"f=24",
				"q=2",
			},
		},
		{
			name: "全零值",
			options: Options{
				Format: 0,
				Action: 0,
				Delete: 0,
			},
			expected: []string{}, // 应使用默认值且不生成选项
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.options.Options()

			// Sort both slices to ensure consistent comparison
			sortStrings(got)
			sortStrings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Options.Options() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestOptions_Validation(t *testing.T) {
	tests := []struct {
		name    string
		options Options
		check   func([]string) bool
	}{
		{
			name: "格式验证",
			options: Options{
				Format: 999, // 无效格式
			},
			check: func(opts []string) bool {
				// 即使无效也应该输出格式
				return containsOption(opts, "f=999")
			},
		},
		{
			name: "带资源的删除",
			options: Options{
				Delete:          DeleteID,
				DeleteResources: true,
			},
			check: func(opts []string) bool {
				// 当 DeleteResources 为 true 时应为大写
				return containsOption(opts, "d=I")
			},
		},
		{
			name: "带文件的传输",
			options: Options{
				File: "/tmp/test.png",
			},
			check: func(opts []string) bool {
				return containsOption(opts, "t=f")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.options.Options()
			if !tt.check(got) {
				t.Errorf("Options validation failed for %s: %v", tt.name, got)
			}
		})
	}
}

func TestOptions_String(t *testing.T) {
	tests := []struct {
		name string
		o    Options
		want string
	}{
		{
			name: "empty options",
			o:    Options{},
			want: "",
		},
		{
			name: "full options",
			o: Options{
				Action:            'A',
				Quite:             'Q',
				Compression:       'C',
				Transmission:      'T',
				Delete:            'd',
				DeleteResources:   true,
				ID:                123,
				PlacementID:       456,
				Number:            789,
				Format:            1,
				ImageWidth:        800,
				ImageHeight:       600,
				Size:              1024,
				Offset:            10,
				Chunk:             true,
				X:                 100,
				Y:                 200,
				Z:                 300,
				Width:             400,
				Height:            500,
				OffsetX:           50,
				OffsetY:           60,
				Columns:           4,
				Rows:              3,
				VirtualPlacement:  true,
				ParentID:          999,
				ParentPlacementID: 888,
			},
			want: "f=1,q=81,i=123,p=456,I=789,s=800,v=600,t=T,S=1024,O=10,U=1,P=999,Q=888,x=100,y=200,z=300,w=400,h=500,X=50,Y=60,c=4,r=3,d=D,a=A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("Options.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		o       Options
		want    []byte
		wantErr bool
	}{
		{
			name: "序列化空选项",
			o:    Options{},
			want: []byte(""),
		},
		{
			name: "带值序列化",
			o: Options{
				Action:          'A',
				ID:              123,
				Width:           400,
				Height:          500,
				Quite:           2,
				DoNotMoveCursor: true,
			},
			want: []byte("q=2,i=123,C=1,w=400,h=500,a=A"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Options.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.MarshalText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOptions_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    Options
		wantErr bool
	}{
		{
			name: "反序列化空",
			text: []byte(""),
			want: Options{},
		},
		{
			name: "反序列化基本选项",
			text: []byte("a=A,i=123,w=400,h=500"),
			want: Options{
				Action: 'A',
				ID:     123,
				Width:  400,
				Height: 500,
			},
		},
		{
			name: "反序列化带无效数字",
			text: []byte("i=abc"),
			want: Options{},
		},
		{
			name: "反序列化带删除资源",
			text: []byte("d=D"),
			want: Options{
				Delete:          'd',
				DeleteResources: true,
			},
		},
		{
			name: "反序列化带布尔分块",
			text: []byte("m=1"),
			want: Options{
				Chunk: true,
			},
		},
		{
			name: "反序列化带虚拟放置",
			text: []byte("U=1"),
			want: Options{
				VirtualPlacement: true,
			},
		},
		{
			name: "反序列化带无效格式",
			text: []byte("invalid=format"),
			want: Options{},
		},
		{
			name: "反序列化带缺失值",
			text: []byte("a="),
			want: Options{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o Options
			err := o.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Options.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("Options.UnmarshalText() = %+v, want %+v", o, tt.want)
			}
		})
	}
}

// Helper functions

func sortStrings(s []string) {
	sort.Strings(s)
}

func containsOption(opts []string, target string) bool {
	return slices.Contains(opts, target)
}
