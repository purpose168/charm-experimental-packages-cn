package pony

import uv "github.com/charmbracelet/ultraviolet"

// Spacer 表示可以增长以填充可用空间的空白区域。
type Spacer struct {
	BaseElement
	fixedSize int // 固定大小，0 表示灵活大小
}

var _ Element = (*Spacer)(nil)

// NewSpacer 创建一个新的间隔元素。
func NewSpacer() *Spacer {
	return &Spacer{}
}

// NewFixedSpacer 创建一个具有固定大小的新间隔。
func NewFixedSpacer(size int) *Spacer {
	return &Spacer{fixedSize: size}
}

// FixedSize 设置大小并返回间隔以支持链式调用。
func (s *Spacer) FixedSize(size int) *Spacer {
	s.fixedSize = size
	return s
}

// Draw 渲染间隔（无需绘制任何内容）。
func (s *Spacer) Draw(_ uv.Screen, area uv.Rectangle) {
	s.SetBounds(area)
	// 间隔是不可见的
}

// Layout 计算间隔大小。
func (s *Spacer) Layout(constraints Constraints) Size {
	if s.fixedSize > 0 {
		return constraints.Constrain(Size{Width: s.fixedSize, Height: s.fixedSize})
	}
	// 灵活间隔 - 占用所有可用空间
	return Size{Width: constraints.MaxWidth, Height: constraints.MaxHeight}
}

// Children 对间隔返回 nil。
func (s *Spacer) Children() []Element {
	return nil
}
