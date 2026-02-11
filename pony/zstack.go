package pony

import uv "github.com/charmbracelet/ultraviolet"

// ZStack 表示一个分层堆栈容器，子元素会绘制在彼此之上。
// 堆栈中后面的子元素会绘制在前面子元素的上面。
type ZStack struct {
	BaseElement
	items             []Element
	width             SizeConstraint
	height            SizeConstraint
	alignment         string // 水平对齐方式：leading（左对齐）, center（居中）, trailing（右对齐）
	verticalAlignment string // 垂直对齐方式：top（顶部）, center（居中）, bottom（底部）
}

var _ Element = (*ZStack)(nil)

// NewZStack 创建一个新的分层堆栈。
func NewZStack(children ...Element) *ZStack {
	return &ZStack{
		items:             children,
		alignment:         AlignmentCenter,
		verticalAlignment: AlignmentCenter,
	}
}

// Alignment 设置水平对齐方式并返回 zstack 以支持链式调用。
func (z *ZStack) Alignment(alignment string) *ZStack {
	z.alignment = alignment
	return z
}

// VerticalAlignment 设置垂直对齐方式并返回 zstack 以支持链式调用。
func (z *ZStack) VerticalAlignment(alignment string) *ZStack {
	z.verticalAlignment = alignment
	return z
}

// Width 设置宽度约束并返回 zstack 以支持链式调用。
func (z *ZStack) Width(width SizeConstraint) *ZStack {
	z.width = width
	return z
}

// Height 设置高度约束并返回 zstack 以支持链式调用。
func (z *ZStack) Height(height SizeConstraint) *ZStack {
	z.height = height
	return z
}

// Draw 将分层堆栈渲染到屏幕上。
func (z *ZStack) Draw(scr uv.Screen, area uv.Rectangle) {
	z.SetBounds(area)

	if len(z.items) == 0 {
		return
	}

	// 首先布局所有子元素以获取它们的大小
	childConstraints := Constraints{
		MinWidth:  0,
		MaxWidth:  area.Dx(),
		MinHeight: 0,
		MaxHeight: area.Dy(),
	}

	childSizes := make([]Size, len(z.items))
	for i, child := range z.items {
		childSizes[i] = child.Layout(childConstraints)
	}

	// 按顺序绘制每个子元素（后面的子元素绘制在上面）
	for i, child := range z.items {
		// 定位元素自行处理布局和定位
		if _, isPositioned := child.(*Positioned); isPositioned {
			child.Draw(scr, area)
			continue
		}

		childSize := childSizes[i]

		// 使用 UV 布局辅助函数基于对齐方式计算子元素区域
		var childArea uv.Rectangle

		// 根据水平和垂直对齐方式确定定位
		switch {
		case z.alignment == AlignmentLeading && z.verticalAlignment == AlignmentTop:
			childArea = uv.TopLeftRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentCenter && z.verticalAlignment == AlignmentTop:
			childArea = uv.TopCenterRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentTrailing && z.verticalAlignment == AlignmentTop:
			childArea = uv.TopRightRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentLeading && z.verticalAlignment == AlignmentCenter:
			childArea = uv.LeftCenterRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentCenter && z.verticalAlignment == AlignmentCenter:
			childArea = uv.CenterRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentTrailing && z.verticalAlignment == AlignmentCenter:
			childArea = uv.RightCenterRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentLeading && z.verticalAlignment == AlignmentBottom:
			childArea = uv.BottomLeftRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentCenter && z.verticalAlignment == AlignmentBottom:
			childArea = uv.BottomCenterRect(area, childSize.Width, childSize.Height)
		case z.alignment == AlignmentTrailing && z.verticalAlignment == AlignmentBottom:
			childArea = uv.BottomRightRect(area, childSize.Width, childSize.Height)
		default:
			// 默认为左上角
			childArea = uv.TopLeftRect(area, childSize.Width, childSize.Height)
		}

		child.Draw(scr, childArea)
	}
}

// Layout 计算分层堆栈的总大小。
// ZStack 取所有子元素的最大宽度和高度。
func (z *ZStack) Layout(constraints Constraints) Size {
	if len(z.items) == 0 {
		return Size{Width: 0, Height: 0}
	}

	maxWidth := 0
	maxHeight := 0

	// 查找最大尺寸
	for _, child := range z.items {
		size := child.Layout(constraints)
		if size.Width > maxWidth {
			maxWidth = size.Width
		}
		if size.Height > maxHeight {
			maxHeight = size.Height
		}
	}

	result := Size{Width: maxWidth, Height: maxHeight}

	// 如果指定了宽度约束，则应用它
	if !z.width.IsAuto() {
		result.Width = z.width.Apply(constraints.MaxWidth, result.Width)
	}

	// 如果指定了高度约束，则应用它
	if !z.height.IsAuto() {
		result.Height = z.height.Apply(constraints.MaxHeight, result.Height)
	}

	return constraints.Constrain(result)
}

// Children 返回子元素。
func (z *ZStack) Children() []Element {
	return z.items
}
