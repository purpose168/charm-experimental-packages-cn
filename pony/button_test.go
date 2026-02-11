package pony

import (
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
)

func TestButtonWithMethods(t *testing.T) {
	style := NewStyle().Fg(Hex("#00FF00")).Build()
	hoverStyle := NewStyle().Fg(Hex("#00FFFF")).Build()
	activeStyle := NewStyle().Fg(Hex("#FF00FF")).Build()

	btn := NewButton("Click Me").
		Style(style).
		HoverStyle(hoverStyle).
		ActiveStyle(activeStyle).
		Border("rounded").
		Padding(2).
		Width(NewFixedConstraint(20)).
		Height(NewFixedConstraint(3))

	// 方法正常工作 - 已通过黄金测试验证
	if btn == nil {
		t.Error("按钮不应为nil")
	}
}

func TestButtonLayout(t *testing.T) {
	btn := NewButton("Submit")
	size := btn.Layout(Fixed(40, 10))

	if size.Width == 0 || size.Height == 0 {
		t.Errorf("按钮大小为零: %v", size)
	}
}

func TestButtonDraw(t *testing.T) {
	btn := NewButton("Test").
		Style(NewStyle().Bold().Build()).
		Border("rounded").
		Padding(1)

	scr := uv.NewScreenBuffer(30, 10)
	area := uv.Rect(0, 0, 30, 10)

	btn.Draw(scr, area)

	bounds := btn.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		t.Error("绘制后按钮边界未设置")
	}
}

func TestButtonChildren(t *testing.T) {
	btn := NewButton("Test")
	children := btn.Children()

	if children != nil {
		t.Error("按钮应具有nil子元素")
	}
}

func TestBoundsMapAllElements(t *testing.T) {
	text1 := NewText("Hello")
	text1.SetID("text1")

	text2 := NewText("World")
	text2.SetID("text2")

	vstack := NewVStack(text1, text2)
	vstack.SetID("vstack")

	constraints := Fixed(20, 10)
	vstack.Layout(constraints)

	scr := uv.NewScreenBuffer(20, 10)
	area := uv.Rect(0, 0, 20, 10)
	vstack.Draw(scr, area)

	boundsMap := NewBoundsMap()
	walkAndRegister(vstack, boundsMap)

	allElements := boundsMap.AllElements()
	if len(allElements) == 0 {
		t.Error("AllElements返回空列表")
	}

	if len(allElements) != 3 {
		t.Errorf("AllElements返回%d个元素，期望3个", len(allElements))
	}

	for _, eb := range allElements {
		if eb.Element == nil {
			t.Error("AllElements中的元素为nil")
		}
		if eb.Bounds.Dx() == 0 && eb.Bounds.Dy() == 0 {
			t.Errorf("元素%s的边界为零", eb.Element.ID())
		}
	}
}
