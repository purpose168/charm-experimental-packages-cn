package pony

import uv "github.com/charmbracelet/ultraviolet"

// Positioned 表示一个绝对定位的元素。
// 该元素相对于其父元素定位在特定坐标处。
type Positioned struct {
	BaseElement
	child  Element
	x      int // X 位置（从左侧开始的单元格数）
	y      int // Y 位置（从顶部开始的单元格数）
	right  int // 距右边缘的距离（如果 >= 0，则覆盖 X）
	bottom int // 距下边缘的距离（如果 >= 0，则覆盖 Y）
	width  SizeConstraint
	height SizeConstraint
}

var _ Element = (*Positioned)(nil)

// NewPositioned 创建一个新的绝对定位元素。
func NewPositioned(child Element, x, y int) *Positioned {
	return &Positioned{
		child:  child,
		x:      x,
		y:      y,
		right:  -1,
		bottom: -1,
	}
}

// Right 设置右边缘距离并返回定位元素以支持链式调用。
// 当设置为 >= 0 时，这会覆盖 X 位置。
func (p *Positioned) Right(right int) *Positioned {
	p.right = right
	return p
}

// Bottom 设置下边缘距离并返回定位元素以支持链式调用。
// 当设置为 >= 0 时，这会覆盖 Y 位置。
func (p *Positioned) Bottom(bottom int) *Positioned {
	p.bottom = bottom
	return p
}

// Width 设置宽度约束并返回定位元素以支持链式调用。
func (p *Positioned) Width(width SizeConstraint) *Positioned {
	p.width = width
	return p
}

// Height 设置高度约束并返回定位元素以支持链式调用。
func (p *Positioned) Height(height SizeConstraint) *Positioned {
	p.height = height
	return p
}

// Draw 渲染定位元素。
func (p *Positioned) Draw(scr uv.Screen, area uv.Rectangle) {
	p.SetBounds(area)

	if p.child == nil {
		return
	}

	// 计算子元素大小
	constraints := Constraints{
		MinWidth:  0,
		MaxWidth:  area.Dx(),
		MinHeight: 0,
		MaxHeight: area.Dy(),
	}

	// 如果指定了宽度/高度约束，则应用它们
	if !p.width.IsAuto() {
		width := p.width.Apply(area.Dx(), area.Dx())
		constraints.MinWidth = width
		constraints.MaxWidth = width
	}

	if !p.height.IsAuto() {
		height := p.height.Apply(area.Dy(), area.Dy())
		constraints.MinHeight = height
		constraints.MaxHeight = height
	}

	childSize := p.child.Layout(constraints)

	// 根据定位约束计算位置
	var childArea uv.Rectangle

	// 使用 UV 布局辅助函数处理右/下定位
	if p.right >= 0 && p.bottom >= 0 {
		// 同时设置了右和下 - 从右下角定位
		childArea = uv.BottomRightRect(area, childSize.Width+p.right, childSize.Height+p.bottom)
		// 调整偏移量
		childArea.Min.X = childArea.Max.X - childSize.Width - p.right
		childArea.Max.X = childArea.Max.X - p.right
		childArea.Min.Y = childArea.Max.Y - childSize.Height - p.bottom
		childArea.Max.Y = childArea.Max.Y - p.bottom
	} else if p.right >= 0 {
		// 设置了右 - 从右边缘定位
		x := area.Max.X - p.right - childSize.Width
		y := area.Min.Y + p.y
		childArea = uv.Rect(x, y, childSize.Width, childSize.Height)
	} else if p.bottom >= 0 {
		// 设置了下 - 从下边缘定位
		x := area.Min.X + p.x
		y := area.Max.Y - p.bottom - childSize.Height
		childArea = uv.Rect(x, y, childSize.Width, childSize.Height)
	} else {
		// 标准 X/Y 定位（从左上角开始）
		x := area.Min.X + p.x
		y := area.Min.Y + p.y
		childArea = uv.Rect(x, y, childSize.Width, childSize.Height)
	}

	// 确保区域在父元素边界内
	childArea = childArea.Intersect(area)

	p.child.Draw(scr, childArea)
}

// Layout 计算定位元素的大小。
// 定位元素不影响父元素布局 - 它们返回 0 大小。
func (p *Positioned) Layout(_ Constraints) Size {
	// 定位元素被排除在正常流之外
	return Size{Width: 0, Height: 0}
}

// Children 返回子元素。
func (p *Positioned) Children() []Element {
	if p.child == nil {
		return nil
	}
	return []Element{p.child}
}
