package pony

import (
	"testing"

	"github.com/purpose168/charm-experimental-packages-cn/exp/golden"
)

func TestSizeConstraint(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		available int
		content   int
		want      int
	}{
		{
			name:      "auto uses content",
			input:     "auto",
			available: 100,
			content:   50,
			want:      50,
		},
		{
			name:      "auto limited by available",
			input:     "auto",
			available: 30,
			content:   50,
			want:      30,
		},
		{
			name:      "50 percent",
			input:     "50%",
			available: 100,
			content:   10,
			want:      50,
		},
		{
			name:      "100 percent",
			input:     "100%",
			available: 80,
			content:   10,
			want:      80,
		},
		{
			name:      "fixed size",
			input:     "20",
			available: 100,
			content:   10,
			want:      20,
		},
		{
			name:      "fixed size larger than available",
			input:     "150",
			available: 100,
			content:   10,
			want:      100,
		},
		{
			name:      "empty string is auto",
			input:     "",
			available: 100,
			content:   30,
			want:      30,
		},
		{
			name:      "max takes available",
			input:     "max",
			available: 100,
			content:   30,
			want:      100,
		},
		{
			name:      "min takes content",
			input:     "min",
			available: 100,
			content:   30,
			want:      30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := parseSizeConstraint(tt.input)
			got := sc.Apply(tt.available, tt.content)
			if got != tt.want {
				t.Errorf("Apply(%d, %d) = %d, 期望 %d (constraint: %s)",
					tt.available, tt.content, got, tt.want, sc.String())
			}
		})
	}
}

func TestBoxWithWidth(t *testing.T) {
	const markup = `
<hstack>
	<box border="normal" width="30%">
		<text>30%</text>
	</box>
	<box border="normal" width="70%">
		<text>70%</text>
	</box>
</hstack>
`

	tmpl, err := Parse[any](markup)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output := tmpl.Render(nil, 100, 10)
	golden.RequireEqual(t, output)
}

func TestZeroValueIsAuto(t *testing.T) {
	// 测试零值 SizeConstraint 是否表现得像 auto
	sc := SizeConstraint{}

	if !sc.IsAuto() {
		t.Error("零值 SizeConstraint 应该是 auto")
	}

	got := sc.Apply(100, 50)
	if got != 50 {
		t.Errorf("零值 Apply(100, 50) = %d, 期望 50 (内容大小)", got)
	}
}

// 测试尺寸约束方法。
func TestSizeConstraintMethods(t *testing.T) {
	// IsFixed
	sc := NewFixedConstraint(10)
	if !sc.IsFixed() {
		t.Error("IsFixed 对于固定约束应该返回 true")
	}

	// IsPercent
	sc = NewPercentConstraint(50)
	if !sc.IsPercent() {
		t.Error("IsPercent 对于百分比约束应该返回 true")
	}

	// String - fixed
	sc = NewFixedConstraint(10)
	if sc.String() != "10" {
		t.Errorf("String() = %s, 期望 10", sc.String())
	}

	// String - percent
	sc = NewPercentConstraint(50)
	if sc.String() != "50%" {
		t.Errorf("String() = %s, 期望 50%%", sc.String())
	}

	// String - auto
	sc = parseSizeConstraint("auto")
	if sc.String() != "auto" {
		t.Errorf("String() = %s, 期望 auto", sc.String())
	}

	// String - min
	sc = parseSizeConstraint("min")
	if sc.String() != "min" {
		t.Errorf("String() = %s, 期望 min", sc.String())
	}

	// String - max
	sc = parseSizeConstraint("max")
	if sc.String() != "max" {
		t.Errorf("String() = %s, 期望 max", sc.String())
	}
}

// 测试 VStack 链式方法。
func TestVStackWithMethods(t *testing.T) {
	vstack := NewVStack(NewText("a"), NewText("b"))

	result := vstack.
		Spacing(2).
		Alignment(AlignmentCenter).
		Width(NewFixedConstraint(10)).
		Height(NewFixedConstraint(5))

	if result == nil {
		t.Error("方法链式调用应该返回 vstack")
	}
}

// 测试 HStack 构造函数和链式方法。
func TestHStackConstructor(t *testing.T) {
	hstack := NewHStack(NewText("a"), NewText("b"))
	if hstack == nil {
		t.Fatal("NewHStack 返回了 nil")
	}
	if len(hstack.Children()) != 2 {
		t.Error("NewHStack 项目未设置")
	}

	hstack.Spacing(2)
	hstack.Alignment(AlignmentCenter)
	hstack.Width(NewFixedConstraint(10))
	hstack.Height(NewFixedConstraint(5))

	// 测试 Children
	children := hstack.Children()
	if len(children) != 2 {
		t.Error("HStack Children 应该返回项目")
	}
}
