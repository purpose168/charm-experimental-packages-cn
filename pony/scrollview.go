package pony

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// ScrollView 表示一个可滚动容器。
// 内容可以大于视口，超出部分将被裁剪。
type ScrollView struct {
	BaseElement
	child Element

	// 滚动位置
	offsetX int
	offsetY int

	// 视口大小约束
	width  SizeConstraint
	height SizeConstraint

	// 滚动条选项
	showScrollbar  bool
	scrollbarColor color.Color

	// 滚动方向
	horizontal bool // 如果为 true，则水平滚动
	vertical   bool // 如果为 true，则垂直滚动（默认）
}

var _ Element = (*ScrollView)(nil)

// NewScrollView 创建一个新的可滚动视图。
func NewScrollView(child Element) *ScrollView {
	return &ScrollView{
		child:         child,
		vertical:      true, // 默认垂直滚动
		showScrollbar: true,
	}
}

// Offset 设置滚动偏移量并返回滚动视图以支持链式调用。
func (s *ScrollView) Offset(x, y int) *ScrollView {
	s.offsetX = x
	s.offsetY = y
	return s
}

// Vertical 启用/禁用垂直滚动。
func (s *ScrollView) Vertical(enabled bool) *ScrollView {
	s.vertical = enabled
	return s
}

// Horizontal 启用/禁用水平滚动。
func (s *ScrollView) Horizontal(enabled bool) *ScrollView {
	s.horizontal = enabled
	return s
}

// Scrollbar 启用/禁用滚动条。
func (s *ScrollView) Scrollbar(show bool) *ScrollView {
	s.showScrollbar = show
	return s
}

// ScrollbarColor 设置滚动条颜色。
func (s *ScrollView) ScrollbarColor(c color.Color) *ScrollView {
	s.scrollbarColor = c
	return s
}

// Width 设置宽度约束。
func (s *ScrollView) Width(width SizeConstraint) *ScrollView {
	s.width = width
	return s
}

// Height 设置高度约束。
func (s *ScrollView) Height(height SizeConstraint) *ScrollView {
	s.height = height
	return s
}

// ScrollUp 向上滚动指定的量。
func (s *ScrollView) ScrollUp(amount int) {
	s.offsetY = max(0, s.offsetY-amount)
}

// ScrollDown 向下滚动指定的量。
func (s *ScrollView) ScrollDown(amount int, contentHeight, viewportHeight int) {
	maxOffset := max(0, contentHeight-viewportHeight)
	s.offsetY = min(maxOffset, s.offsetY+amount)
}

// ScrollLeft 向左滚动指定的量。
func (s *ScrollView) ScrollLeft(amount int) {
	s.offsetX = max(0, s.offsetX-amount)
}

// ScrollRight 向右滚动指定的量。
func (s *ScrollView) ScrollRight(amount int, contentWidth, viewportWidth int) {
	maxOffset := max(0, contentWidth-viewportWidth)
	s.offsetX = min(maxOffset, s.offsetX+amount)
}

// Draw 渲染可滚动视图。
func (s *ScrollView) Draw(scr uv.Screen, area uv.Rectangle) {
	s.SetBounds(area)

	if s.child == nil {
		return
	}

	// 计算视口大小
	viewportWidth := area.Dx()
	viewportHeight := area.Dy()

	// 如果显示滚动条，预留空间
	scrollbarWidth := 0
	scrollbarHeight := 0
	if s.showScrollbar {
		if s.vertical && !s.horizontal {
			scrollbarWidth = 1
			viewportWidth -= scrollbarWidth
		}
		if s.horizontal && !s.vertical {
			scrollbarHeight = 1
			viewportHeight -= scrollbarHeight
		}
		if s.horizontal && s.vertical {
			scrollbarWidth = 1
			scrollbarHeight = 1
			viewportWidth -= scrollbarWidth
			viewportHeight -= scrollbarHeight
		}
	}

	// 使用无界约束布局子元素以获取完整内容大小
	contentConstraints := Constraints{
		MinWidth:  0,
		MaxWidth:  1 << 30, // 非常大的数字
		MinHeight: 0,
		MaxHeight: 1 << 30,
	}
	contentSize := s.child.Layout(contentConstraints)

	// 为完整内容创建缓冲区
	contentBuffer := uv.NewScreenBuffer(contentSize.Width, contentSize.Height)
	contentArea := uv.Rect(0, 0, contentSize.Width, contentSize.Height)
	s.child.Draw(contentBuffer, contentArea)

	// 调整子元素边界到屏幕坐标（考虑视口位置和滚动偏移）
	s.adjustChildBounds(s.child, area.Min.X-s.offsetX, area.Min.Y-s.offsetY)

	// 将可见部分复制到屏幕（带偏移）
	for y := 0; y < viewportHeight; y++ {
		for x := 0; x < viewportWidth; x++ {
			// 内容缓冲区中的源位置（带偏移）
			srcX := x + s.offsetX
			srcY := y + s.offsetY

			// 屏幕上的目标位置
			dstX := area.Min.X + x
			dstY := area.Min.Y + y

			// 如果在边界内，复制单元格
			if srcY < contentSize.Height && srcX < contentSize.Width {
				cell := contentBuffer.CellAt(srcX, srcY)
				scr.SetCell(dstX, dstY, cell)
			}
		}
	}

	// 如果启用，绘制滚动条
	if s.showScrollbar {
		if s.vertical {
			s.drawVerticalScrollbar(scr, area, contentSize.Height, viewportHeight, scrollbarWidth)
		}
		if s.horizontal {
			s.drawHorizontalScrollbar(scr, area, contentSize.Width, viewportWidth, scrollbarHeight)
		}
	}
}

// drawVerticalScrollbar 绘制垂直滚动条。
func (s *ScrollView) drawVerticalScrollbar(scr uv.Screen, area uv.Rectangle, contentHeight, viewportHeight, scrollbarWidth int) {
	if contentHeight <= viewportHeight {
		return // 不需要滚动条
	}

	scrollbarX := area.Max.X - scrollbarWidth
	scrollbarStart := area.Min.Y
	scrollbarEnd := area.Max.Y
	trackHeight := scrollbarEnd - scrollbarStart

	// 计算滚动条滑块大小
	thumbHeight := max(1, (viewportHeight*trackHeight)/contentHeight)

	// 计算滚动条滑块位置
	// scrollableRange 是我们可以滚动的距离
	scrollableRange := contentHeight - viewportHeight
	// trackRange 是滑块可以移动的距离
	trackRange := trackHeight - thumbHeight

	// 按比例定位滑块
	thumbPos := scrollbarStart
	if scrollableRange > 0 {
		thumbPos = scrollbarStart + (s.offsetY*trackRange)/scrollableRange
	}

	// 确保滑块保持在边界内（处理舍入边缘情况）
	if thumbPos+thumbHeight > scrollbarEnd {
		thumbPos = scrollbarEnd - thumbHeight
	}
	if thumbPos < scrollbarStart {
		thumbPos = scrollbarStart
	}

	// 创建滚动条单元格
	trackCell := uv.NewCell(scr.WidthMethod(), "░")
	thumbCell := uv.NewCell(scr.WidthMethod(), "█")
	if thumbCell != nil && s.scrollbarColor != nil {
		thumbCell.Style = uv.Style{Fg: s.scrollbarColor}
	}

	// 绘制滚动条
	for y := scrollbarStart; y < scrollbarEnd; y++ {
		if y >= thumbPos && y < thumbPos+thumbHeight {
			scr.SetCell(scrollbarX, y, thumbCell)
		} else {
			scr.SetCell(scrollbarX, y, trackCell)
		}
	}
}

// drawHorizontalScrollbar 绘制水平滚动条。
func (s *ScrollView) drawHorizontalScrollbar(scr uv.Screen, area uv.Rectangle, contentWidth, viewportWidth, scrollbarHeight int) {
	if contentWidth <= viewportWidth {
		return
	}

	scrollbarY := area.Max.Y - scrollbarHeight
	scrollbarStart := area.Min.X
	scrollbarEnd := area.Max.X
	trackWidth := scrollbarEnd - scrollbarStart

	// 计算滚动条滑块大小
	thumbWidth := max(1, (viewportWidth*trackWidth)/contentWidth)

	// 计算滚动条滑块位置
	scrollableRange := contentWidth - viewportWidth
	trackRange := trackWidth - thumbWidth

	thumbPos := scrollbarStart
	if scrollableRange > 0 {
		thumbPos = scrollbarStart + (s.offsetX*trackRange)/scrollableRange
	}

	// 确保滑块保持在边界内（处理舍入边缘情况）
	if thumbPos+thumbWidth > scrollbarEnd {
		thumbPos = scrollbarEnd - thumbWidth
	}
	if thumbPos < scrollbarStart {
		thumbPos = scrollbarStart
	}

	trackCell := uv.NewCell(scr.WidthMethod(), "░")
	thumbCell := uv.NewCell(scr.WidthMethod(), "█")
	if thumbCell != nil && s.scrollbarColor != nil {
		thumbCell.Style = uv.Style{Fg: s.scrollbarColor}
	}

	for x := scrollbarStart; x < scrollbarEnd; x++ {
		if x >= thumbPos && x < thumbPos+thumbWidth {
			scr.SetCell(x, scrollbarY, thumbCell)
		} else {
			scr.SetCell(x, scrollbarY, trackCell)
		}
	}
}

// Layout 计算滚动视图大小。
func (s *ScrollView) Layout(constraints Constraints) Size {
	// 从最大可用空间开始
	viewportWidth := constraints.MaxWidth
	viewportHeight := constraints.MaxHeight

	// 如果指定了宽度/高度约束，则应用
	if !s.width.IsAuto() {
		viewportWidth = s.width.Apply(constraints.MaxWidth, constraints.MaxWidth)
	}

	if !s.height.IsAuto() {
		viewportHeight = s.height.Apply(constraints.MaxHeight, constraints.MaxHeight)
	}

	// 约束最终大小
	return Size{
		Width:  min(viewportWidth, constraints.MaxWidth),
		Height: min(viewportHeight, constraints.MaxHeight),
	}
}

// Children 返回子元素。
func (s *ScrollView) Children() []Element {
	if s.child == nil {
		return nil
	}
	return []Element{s.child}
}

// ContentSize 返回内容的完整大小。
func (s *ScrollView) ContentSize() Size {
	if s.child == nil {
		return Size{Width: 0, Height: 0}
	}

	// 使用无界约束布局以获取完整大小
	unbounded := Constraints{
		MinWidth:  0,
		MaxWidth:  1 << 30,
		MinHeight: 0,
		MaxHeight: 1 << 30,
	}

	return s.child.Layout(unbounded)
}

// adjustChildBounds 递归调整所有子元素的边界
// 以考虑滚动视图的视口位置和滚动偏移。
// 这确保了滚动视图内元素的点击测试能够正确工作。
func (s *ScrollView) adjustChildBounds(elem Element, offsetX, offsetY int) {
	if elem == nil {
		return
	}

	// 获取当前边界（相对于内容缓冲区的 0,0）
	bounds := elem.Bounds()

	// 转换为屏幕坐标
	newBounds := uv.Rect(
		bounds.Min.X+offsetX,
		bounds.Min.Y+offsetY,
		bounds.Dx(),
		bounds.Dy(),
	)

	// 更新元素的边界
	elem.SetBounds(newBounds)

	// 递归调整子元素
	for _, child := range elem.Children() {
		if child != nil {
			s.adjustChildBounds(child, offsetX, offsetY)
		}
	}
}
