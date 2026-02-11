package pony

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// Slot 表示动态内容的占位符。
// 通过 RenderWithSlots 传递的 Elements 会填充到 Slots 中。
type Slot struct {
	BaseElement
	Name    string
	element Element // 实际要渲染的元素（在渲染期间填充）
}

var _ Element = (*Slot)(nil)

// NewSlot 创建一个新的插槽元素。
func NewSlot(name string) *Slot {
	return &Slot{Name: name}
}

// Draw 如果插槽的元素存在，则渲染它。
func (s *Slot) Draw(scr uv.Screen, area uv.Rectangle) {
	s.SetBounds(area)

	if s.element != nil {
		s.element.Draw(scr, area)
	}
}

// Layout 如果插槽的元素存在，则计算其大小。
func (s *Slot) Layout(constraints Constraints) Size {
	if s.element != nil {
		return s.element.Layout(constraints)
	}
	return Size{Width: 0, Height: 0}
}

// Children 如果插槽的元素存在，则返回其元素子项。
func (s *Slot) Children() []Element {
	if s.element != nil {
		return []Element{s.element}
	}
	return nil
}

// setElement 为此插槽设置元素（在渲染期间内部使用）。
func (s *Slot) setElement(elem Element) {
	s.element = elem
}
