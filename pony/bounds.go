package pony

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// BoundsMap 跟踪所有渲染元素及其位置，用于点击测试。
// 创建后不可变，并发读取安全。
type BoundsMap struct {
	elements   map[string]elementBounds
	byPosition []elementBounds // 按 z-index 排序（最后一个 = 在顶部）
}

type elementBounds struct {
	id     string
	elem   Element
	bounds uv.Rectangle
}

// NewBoundsMap 创建一个新的空边界映射。
func NewBoundsMap() *BoundsMap {
	return &BoundsMap{
		elements: make(map[string]elementBounds),
	}
}

// Register 记录元素及其渲染边界。
// 应在渲染过程中调用。
func (bm *BoundsMap) Register(elem Element, bounds uv.Rectangle) {
	eb := elementBounds{
		id:     elem.ID(),
		elem:   elem,
		bounds: bounds,
	}
	bm.elements[elem.ID()] = eb
	bm.byPosition = append(bm.byPosition, eb)
}

// HitTest 返回给定屏幕坐标处的最顶层元素。
// 当多个元素在某一点重叠时，优先选择具有显式设置的 ID 的元素，而不是自动生成的 ID (elem_*)。
//
// 这种行为对于交互组件至关重要：当你点击组件的渲染区域时，你希望获取组件的 ID，而不是其子元素的 ID。
// 例如，点击 Input 组件中的任何位置都应该返回 Input 的 ID，而不是其内部的 Text 或 Box 子元素的 ID。
//
// 要实现这一点，请在组件返回的根元素上设置组件的 ID：
//
//	func (i *Input) Render() pony.Element {
//	    vstack := pony.NewVStack(...)
//	    vstack.SetID(i.ID())  // 传递组件 ID
//	    return vstack
//	}
//
// 如果在该位置未找到元素，则返回 nil。
func (bm *BoundsMap) HitTest(x, y int) Element {
	var bestMatch Element
	var bestMatchHasExplicitID bool

	// 从末尾搜索（最后绘制 = 在顶部）
	for i := len(bm.byPosition) - 1; i >= 0; i-- {
		eb := bm.byPosition[i]
		if pointInRect(x, y, eb.bounds) {
			// 检查此元素是否有显式 ID（非自动生成）
			hasExplicitID := !strings.HasPrefix(eb.id, "elem_")

			// 第一个匹配项或更好的匹配项（优先选择显式 ID）
			if bestMatch == nil || (!bestMatchHasExplicitID && hasExplicitID) {
				bestMatch = eb.elem
				bestMatchHasExplicitID = hasExplicitID
			}

			// 如果找到具有显式 ID 的元素，那就是我们的最佳匹配项
			if hasExplicitID {
				return bestMatch
			}
		}
	}

	return bestMatch
}

// HitTestAll 返回给定屏幕坐标处的所有元素，
// 从上到下排序（第一个元素在视觉上位于顶部）。
//
// 这对于带有可点击子元素的嵌套交互组件（如滚动视图）非常有用，
// 你需要同时知道被点击的子元素和其父容器。
//
// 使用示例：
//
//	hits := boundsMap.HitTestAll(x, y)
//	for _, elem := range hits {
//	    switch elem.ID() {
//	    case "list-item-5":
//	        // 处理项目点击
//	    case "main-scroll-view":
//	        // 同时跟踪我们在滚动视图中
//	    }
//	}
//
// 如果在该位置未找到元素，则返回空切片。
func (bm *BoundsMap) HitTestAll(x, y int) []Element {
	var hits []Element

	// 从末尾搜索（最后绘制 = 在顶部）
	for i := len(bm.byPosition) - 1; i >= 0; i-- {
		eb := bm.byPosition[i]
		if pointInRect(x, y, eb.bounds) {
			hits = append(hits, eb.elem)
		}
	}

	return hits
}

// HitTestWithContainer 返回顶部元素和第一个带有显式 ID 的父容器。
// 这对于带有可点击项目的滚动视图非常有用，你希望同时知道点击了什么以及它在哪个容器中。
//
// "顶部"元素是坐标处视觉上最顶层的元素。
// "容器"是命中堆栈中（顶部之后）第一个具有显式 ID（非自动生成的 "elem_" 前缀）的元素。
//
// 使用示例：
//
//	top, container := boundsMap.HitTestWithContainer(x, y)
//	if top != nil {
//	    handleClick(top.ID())
//	}
//	if container != nil && container.ID() == "scroll-view" {
//	    // 我们知道我们点击了滚动视图内部
//	}
//
// 如果在该位置未找到元素，则返回 (nil, nil)。
func (bm *BoundsMap) HitTestWithContainer(x, y int) (top Element, container Element) {
	hits := bm.HitTestAll(x, y)
	if len(hits) == 0 {
		return nil, nil
	}

	top = hits[0]

	// 找到第一个容器（具有显式 ID 且不是顶部的元素）
	for i := 1; i < len(hits); i++ {
		if !strings.HasPrefix(hits[i].ID(), "elem_") {
			container = hits[i]
			break
		}
	}

	return top, container
}

// GetByID 通过 ID 检索元素。
func (bm *BoundsMap) GetByID(id string) (Element, bool) {
	eb, ok := bm.elements[id]
	return eb.elem, ok
}

// GetBounds 通过 ID 返回元素的渲染边界。
func (bm *BoundsMap) GetBounds(id string) (uv.Rectangle, bool) {
	eb, ok := bm.elements[id]
	return eb.bounds, ok
}

// AllElements 返回所有注册的元素及其边界。
func (bm *BoundsMap) AllElements() []ElementWithBounds {
	result := make([]ElementWithBounds, 0, len(bm.byPosition))
	for _, eb := range bm.byPosition {
		result = append(result, ElementWithBounds{
			Element: eb.elem,
			Bounds:  eb.bounds,
		})
	}
	return result
}

// ElementWithBounds 将元素与其渲染边界配对。
type ElementWithBounds struct {
	Element Element
	Bounds  uv.Rectangle
}

// pointInRect 检查点是否在矩形内。
func pointInRect(x, y int, rect uv.Rectangle) bool {
	return x >= rect.Min.X && x < rect.Max.X &&
		y >= rect.Min.Y && y < rect.Max.Y
}

// BaseElement 为所有元素提供通用功能。
// 元素应嵌入此结构以获得 ID 和边界跟踪。
type BaseElement struct {
	id     string
	bounds uv.Rectangle
}

// ID 返回元素的标识符。
// 如果未显式设置 ID，则返回基于指针的 ID。
func (b *BaseElement) ID() string {
	if b.id == "" {
		return fmt.Sprintf("elem_%p", b)
	}
	return b.id
}

// SetID 设置元素的标识符。
func (b *BaseElement) SetID(id string) {
	b.id = id
}

// Bounds 返回元素最后渲染的边界。
func (b *BaseElement) Bounds() uv.Rectangle {
	return b.bounds
}

// SetBounds 记录元素的渲染边界。
// 应在 Draw() 开始时调用。
func (b *BaseElement) SetBounds(bounds uv.Rectangle) {
	b.bounds = bounds
}

// walkAndRegister 递归遍历元素树并注册所有元素。
func walkAndRegister(elem Element, bm *BoundsMap) {
	bm.Register(elem, elem.Bounds())

	for _, child := range elem.Children() {
		if child != nil {
			walkAndRegister(child, bm)
		}
	}
}
