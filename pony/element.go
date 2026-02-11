package pony

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// Element 表示 TUI 中可渲染的组件。元素实现了 UV Drawable 接口，并且可以组合成树状结构。
type Element interface {
	uv.Drawable

	// Layout 计算元素在给定约束下的期望大小。
	// 返回元素实际占据的大小。
	Layout(constraints Constraints) Size

	// Children 返回容器类型的子元素。
	// 对于叶元素返回 nil。
	Children() []Element

	// ID 返回此元素的唯一标识符。
	// 用于命中测试和事件处理。
	ID() string

	// SetID 设置元素的标识符。
	SetID(id string)

	// Bounds 返回元素上次渲染的屏幕坐标。
	// 在 Draw() 期间更新，用于鼠标命中测试。
	Bounds() uv.Rectangle

	// SetBounds 记录元素的渲染边界。
	// 应在 Draw() 开始时调用。
	SetBounds(bounds uv.Rectangle)
}

// Size 表示终端单元格中的尺寸。
type Size struct {
	Width  int
	Height int
}

// Constraints 定义布局计算的尺寸约束。
type Constraints struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

// Constrain 返回满足约束的尺寸。
func (c Constraints) Constrain(size Size) Size {
	w := size.Width
	h := size.Height

	if w < c.MinWidth {
		w = c.MinWidth
	}
	if w > c.MaxWidth {
		w = c.MaxWidth
	}
	if h < c.MinHeight {
		h = c.MinHeight
	}
	if h > c.MaxHeight {
		h = c.MaxHeight
	}

	return Size{Width: w, Height: h}
}

// Unbounded 返回无限制的约束。
func Unbounded() Constraints {
	return Constraints{
		MinWidth:  0,
		MaxWidth:  1<<31 - 1,
		MinHeight: 0,
		MaxHeight: 1<<31 - 1,
	}
}

// Fixed 返回固定尺寸的约束。
func Fixed(width, height int) Constraints {
	return Constraints{
		MinWidth:  width,
		MaxWidth:  width,
		MinHeight: height,
		MaxHeight: height,
	}
}

// Constraint 表示可以应用的尺寸约束。
type Constraint interface {
	// Apply 将约束应用于给定的可用空间。
	Apply(available int) int
}

// FixedConstraint 表示单元格中的固定尺寸。
type FixedConstraint int

// Apply 返回固定尺寸，限制在可用空间内。
func (f FixedConstraint) Apply(available int) int {
	if int(f) > available {
		return available
	}
	if f < 0 {
		return 0
	}
	return int(f)
}

// PercentConstraint 表示可用空间的百分比 (0-100)。
type PercentConstraint int

// Apply 返回可用空间的百分比。
func (p PercentConstraint) Apply(available int) int {
	if p < 0 {
		return 0
	}
	if p > 100 {
		return available
	}
	return available * int(p) / 100
}

// AutoConstraint 表示基于内容的尺寸调整。
type AutoConstraint struct{}

// Apply 返回可用空间（将基于内容计算）。
func (a AutoConstraint) Apply(available int) int {
	return available
}

// Props 是传递给元素的属性映射。
type Props map[string]string

// Get 返回属性值，如果未找到则返回空字符串。
func (p Props) Get(key string) string {
	if p == nil {
		return ""
	}
	return p[key]
}

// GetOr 返回属性值，如果未找到则返回默认值。
func (p Props) GetOr(key, defaultValue string) string {
	if p == nil {
		return defaultValue
	}
	if v, ok := p[key]; ok {
		return v
	}
	return defaultValue
}

// Has 检查属性是否存在。
func (p Props) Has(key string) bool {
	if p == nil {
		return false
	}
	_, ok := p[key]
	return ok
}
