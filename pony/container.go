package pony

import uv "github.com/charmbracelet/ultraviolet"

// VStack 表示一个垂直堆栈容器。
type VStack struct {
	BaseElement
	items     []Element
	spacing   int
	width     SizeConstraint
	height    SizeConstraint
	alignment string // leading, center, trailing (子元素的水平对齐方式)
}

var _ Element = (*VStack)(nil)

// NewVStack 创建一个新的垂直堆栈。
func NewVStack(children ...Element) *VStack {
	return &VStack{items: children}
}

// Spacing 设置子元素之间的间距并返回 vstack 以支持链式调用。
func (v *VStack) Spacing(spacing int) *VStack {
	v.spacing = spacing
	return v
}

// Alignment 设置子元素的水平对齐方式并返回 vstack 以支持链式调用。
func (v *VStack) Alignment(alignment string) *VStack {
	v.alignment = alignment
	return v
}

// Width 设置宽度约束并返回 vstack 以支持链式调用。
func (v *VStack) Width(width SizeConstraint) *VStack {
	v.width = width
	return v
}

// Height 设置高度约束并返回 vstack 以支持链式调用。
func (v *VStack) Height(height SizeConstraint) *VStack {
	v.height = height
	return v
}

// calculateChildSizes 对 VStack 的子元素执行两步布局。
// 第一步：布局固定大小的子元素，第二步：根据 flex-grow 分配空间给弹性子元素。
func (v *VStack) calculateChildSizes(constraints Constraints) []Size {
	childSizes := make([]Size, len(v.items))
	if len(v.items) == 0 {
		return childSizes
	}

	// 第一步：布局固定大小的子元素并计算弹性项
	fixedHeight := 0
	totalFlexGrow := 0

	for i, child := range v.items {
		flexGrow := GetFlexGrow(child)

		if flexGrow > 0 {
			// 弹性项 - 将在第二步中设置大小
			totalFlexGrow += flexGrow
			childSizes[i] = Size{Width: 0, Height: 0}
		} else {
			// 固定大小项 - 现在布局
			childConstraints := Constraints{
				MinWidth:  0,
				MaxWidth:  constraints.MaxWidth,
				MinHeight: 0,
				MaxHeight: constraints.MaxHeight - fixedHeight,
			}

			size := child.Layout(childConstraints)
			childSizes[i] = size
			fixedHeight += size.Height
		}

		if i < len(v.items)-1 {
			fixedHeight += v.spacing
		}
	}

	// 第二步：根据 flex-grow 在弹性项之间分配剩余空间
	if totalFlexGrow > 0 {
		remainingHeight := constraints.MaxHeight - fixedHeight
		if remainingHeight > 0 {
			for i, child := range v.items {
				flexGrow := GetFlexGrow(child)
				if flexGrow > 0 {
					// 按 flex-grow 比例分配空间
					flexHeight := (remainingHeight * flexGrow) / totalFlexGrow

					// 使用分配的空间布局弹性子元素
					childConstraints := Constraints{
						MinWidth:  0,
						MaxWidth:  constraints.MaxWidth,
						MinHeight: flexHeight,
						MaxHeight: flexHeight,
					}

					childSizes[i] = child.Layout(childConstraints)
				}
			}
		}
	}

	return childSizes
}

// Draw 将垂直堆栈渲染到屏幕上。
func (v *VStack) Draw(scr uv.Screen, area uv.Rectangle) {
	v.SetBounds(area)

	if len(v.items) == 0 {
		return
	}

	// 使用两步布局计算子元素大小
	childSizes := v.calculateChildSizes(Constraints{
		MinWidth:  0,
		MaxWidth:  area.Dx(),
		MinHeight: 0,
		MaxHeight: area.Dy(),
	})

	// 使用计算的大小绘制所有子元素
	y := area.Min.Y

	for i, child := range v.items {
		if y >= area.Max.Y {
			break
		}

		childSize := childSizes[i]

		// 根据水平对齐方式计算 x 位置
		var x int
		switch v.alignment {
		case AlignmentCenter:
			if childSize.Width < area.Dx() {
				x = area.Min.X + (area.Dx()-childSize.Width)/2
			} else {
				x = area.Min.X
			}
		case AlignmentTrailing:
			if childSize.Width < area.Dx() {
				x = area.Max.X - childSize.Width
			} else {
				x = area.Min.X
			}
		default: // AlignmentLeading
			x = area.Min.X
		}

		childArea := uv.Rect(x, y, childSize.Width, childSize.Height)
		// 裁剪到父容器边界
		childArea = childArea.Intersect(area)

		child.Draw(scr, childArea)

		y += childSize.Height
		if i < len(v.items)-1 {
			y += v.spacing
		}
	}
}

// Layout 计算垂直堆栈的总大小。
func (v *VStack) Layout(constraints Constraints) Size {
	if len(v.items) == 0 {
		return Size{Width: 0, Height: 0}
	}

	// 使用两步布局计算子元素大小
	childSizes := v.calculateChildSizes(constraints)

	// 计算总大小
	totalHeight := 0
	maxWidth := 0

	for i, size := range childSizes {
		totalHeight += size.Height
		if size.Width > maxWidth {
			maxWidth = size.Width
		}

		if i < len(v.items)-1 {
			totalHeight += v.spacing
		}
	}

	result := Size{Width: maxWidth, Height: totalHeight}

	if !v.width.IsAuto() {
		result.Width = v.width.Apply(constraints.MaxWidth, result.Width)
	}

	if !v.height.IsAuto() {
		result.Height = v.height.Apply(constraints.MaxHeight, result.Height)
	}

	return constraints.Constrain(result)
}

// Children 返回子元素。
func (v *VStack) Children() []Element {
	return v.items
}

// HStack 表示一个水平堆栈容器。
type HStack struct {
	BaseElement
	items     []Element
	spacing   int
	width     SizeConstraint
	height    SizeConstraint
	alignment string // top, center, bottom (子元素的垂直对齐方式)
}

var _ Element = (*HStack)(nil)

// NewHStack 创建一个新的水平堆栈。
func NewHStack(children ...Element) *HStack {
	return &HStack{items: children}
}

// Spacing 设置子元素之间的间距并返回 hstack 以支持链式调用。
func (h *HStack) Spacing(spacing int) *HStack {
	h.spacing = spacing
	return h
}

// Alignment 设置子元素的垂直对齐方式并返回 hstack 以支持链式调用。
func (h *HStack) Alignment(alignment string) *HStack {
	h.alignment = alignment
	return h
}

// Width 设置宽度约束并返回 hstack 以支持链式调用。
func (h *HStack) Width(width SizeConstraint) *HStack {
	h.width = width
	return h
}

// Height 设置高度约束并返回 hstack 以支持链式调用。
func (h *HStack) Height(height SizeConstraint) *HStack {
	h.height = height
	return h
}

// calculateChildSizes 对 HStack 的子元素执行两步布局。
// 第一步：布局固定大小的子元素，第二步：根据 flex-grow 分配空间给弹性子元素。
func (h *HStack) calculateChildSizes(constraints Constraints) []Size {
	childSizes := make([]Size, len(h.items))
	if len(h.items) == 0 {
		return childSizes
	}

	// 第一步：布局固定大小的子元素并计算弹性项
	fixedWidth := 0
	totalFlexGrow := 0

	for i, child := range h.items {
		flexGrow := GetFlexGrow(child)

		if flexGrow > 0 {
			// 弹性项 - 将在第二步中设置大小
			totalFlexGrow += flexGrow
			childSizes[i] = Size{Width: 0, Height: 0}
		} else {
			// 固定大小项 - 现在布局
			childConstraints := Constraints{
				MinWidth:  0,
				MaxWidth:  constraints.MaxWidth - fixedWidth,
				MinHeight: 0,
				MaxHeight: constraints.MaxHeight,
			}

			size := child.Layout(childConstraints)
			childSizes[i] = size
			fixedWidth += size.Width
		}

		if i < len(h.items)-1 {
			fixedWidth += h.spacing
		}
	}

	// 第二步：根据 flex-grow 在弹性项之间分配剩余空间
	if totalFlexGrow > 0 {
		remainingWidth := constraints.MaxWidth - fixedWidth
		if remainingWidth > 0 {
			for i, child := range h.items {
				flexGrow := GetFlexGrow(child)
				if flexGrow > 0 {
					// 按 flex-grow 比例分配空间
					flexWidth := (remainingWidth * flexGrow) / totalFlexGrow

					// 使用分配的空间布局弹性子元素
					childConstraints := Constraints{
						MinWidth:  flexWidth,
						MaxWidth:  flexWidth,
						MinHeight: 0,
						MaxHeight: constraints.MaxHeight,
					}

					childSizes[i] = child.Layout(childConstraints)
				}
			}
		}
	}

	return childSizes
}

// Draw 将水平堆栈渲染到屏幕上。
func (h *HStack) Draw(scr uv.Screen, area uv.Rectangle) {
	h.SetBounds(area)

	if len(h.items) == 0 {
		return
	}

	// 使用两步布局计算子元素大小
	childSizes := h.calculateChildSizes(Constraints{
		MinWidth:  0,
		MaxWidth:  area.Dx(),
		MinHeight: 0,
		MaxHeight: area.Dy(),
	})

	// 绘制所有子元素
	x := area.Min.X

	for i, child := range h.items {
		if x >= area.Max.X {
			break
		}

		childSize := childSizes[i]

		// 根据垂直对齐方式计算 y 位置
		var y int
		switch h.alignment {
		case AlignmentCenter:
			if childSize.Height < area.Dy() {
				y = area.Min.Y + (area.Dy()-childSize.Height)/2
			} else {
				y = area.Min.Y
			}
		case AlignmentBottom:
			if childSize.Height < area.Dy() {
				y = area.Max.Y - childSize.Height
			} else {
				y = area.Min.Y
			}
		default: // AlignmentTop
			y = area.Min.X
		}

		childArea := uv.Rect(x, y, childSize.Width, childSize.Height)
		// 裁剪到父容器边界
		childArea = childArea.Intersect(area)

		child.Draw(scr, childArea)

		x += childSize.Width
		if i < len(h.items)-1 {
			x += h.spacing
		}
	}
}

// Layout 计算水平堆栈的总大小。
func (h *HStack) Layout(constraints Constraints) Size {
	if len(h.items) == 0 {
		return Size{Width: 0, Height: 0}
	}

	// 使用两步布局计算子元素大小
	childSizes := h.calculateChildSizes(constraints)

	// 计算总大小
	totalWidth := 0
	maxHeight := 0

	for i, size := range childSizes {
		totalWidth += size.Width
		if size.Height > maxHeight {
			maxHeight = size.Height
		}

		if i < len(h.items)-1 {
			totalWidth += h.spacing
		}
	}

	result := Size{Width: totalWidth, Height: maxHeight}

	if !h.width.IsAuto() {
		result.Width = h.width.Apply(constraints.MaxWidth, result.Width)
	}

	if !h.height.IsAuto() {
		result.Height = h.height.Apply(constraints.MaxHeight, result.Height)
	}

	return constraints.Constrain(result)
}

// Children 返回子元素。
func (h *HStack) Children() []Element {
	return h.items
}
