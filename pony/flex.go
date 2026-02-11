package pony

import uv "github.com/charmbracelet/ultraviolet"

// Flex 表示支持 flex-grow 和 flex-shrink 的元素包装器。
// 使用它可以使元素在 VStack 或 HStack 中具有灵活性。
type Flex struct {
	BaseElement
	child  Element
	grow   int // flex-grow: 相对于兄弟元素的增长比例（默认 0 = 不增长）
	shrink int // flex-shrink: 相对于兄弟元素的收缩比例（默认 1）
	basis  int // flex-basis: 弹性计算前的初始大小（默认 0 = 自动）
}

var _ Element = (*Flex)(nil)

// NewFlex 创建一个新的弹性包装器。
func NewFlex(child Element) *Flex {
	return &Flex{
		child:  child,
		grow:   0,
		shrink: 1,
		basis:  0,
	}
}

// Grow 设置 flex-grow 并返回弹性包装器以支持链式调用。
func (f *Flex) Grow(grow int) *Flex {
	f.grow = grow
	return f
}

// Shrink 设置 flex-shrink 并返回弹性包装器以支持链式调用。
func (f *Flex) Shrink(shrink int) *Flex {
	f.shrink = shrink
	return f
}

// Basis 设置 flex-basis 并返回弹性包装器以支持链式调用。
func (f *Flex) Basis(basis int) *Flex {
	f.basis = basis
	return f
}

// Draw 渲染弹性子元素。
func (f *Flex) Draw(scr uv.Screen, area uv.Rectangle) {
	f.SetBounds(area)

	if f.child != nil {
		f.child.Draw(scr, area)
	}
}

// Layout 计算弹性子元素的大小。
func (f *Flex) Layout(constraints Constraints) Size {
	if f.child == nil {
		return Size{Width: 0, Height: 0}
	}

	// 如果设置了 basis，则将其用作初始大小
	if f.basis > 0 {
		// 创建以 basis 为首选大小的约束
		flexConstraints := constraints
		flexConstraints.MinWidth = min(f.basis, constraints.MaxWidth)
		flexConstraints.MinHeight = min(f.basis, constraints.MaxHeight)
		return f.child.Layout(flexConstraints)
	}

	return f.child.Layout(constraints)
}

// Children 返回子元素。
func (f *Flex) Children() []Element {
	if f.child == nil {
		return nil
	}
	return []Element{f.child}
}

// GetFlexGrow 返回元素的 flex-grow 值。
// 如果元素不是 Flex 包装器，则返回 0。
func GetFlexGrow(elem Element) int {
	if flex, ok := elem.(*Flex); ok {
		return flex.grow
	}
	// 检查是否为弹性间隔符
	if spacer, ok := elem.(*Spacer); ok && spacer.fixedSize == 0 {
		return 1
	}
	return 0
}

// GetFlexShrink 返回元素的 flex-shrink 值。
// 如果元素不是 Flex 包装器，则返回 1。
func GetFlexShrink(elem Element) int {
	if flex, ok := elem.(*Flex); ok {
		return flex.shrink
	}
	return 1
}

// GetFlexBasis 返回元素的 flex-basis 值。
// 如果元素不是 Flex 包装器，则返回 0（自动）。
func GetFlexBasis(elem Element) int {
	if flex, ok := elem.(*Flex); ok {
		return flex.basis
	}
	return 0
}

// IsFlexible 如果元素可以增长，则返回 true。
func IsFlexible(elem Element) bool {
	return GetFlexGrow(elem) > 0
}
