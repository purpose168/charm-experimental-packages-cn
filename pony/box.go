package pony

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// Box 表示一个带有可选边框的容器。
type Box struct {
	BaseElement
	child        Element
	border       string // normal, rounded, thick, double, hidden, none
	borderColor  color.Color
	width        SizeConstraint
	height       SizeConstraint
	padding      int
	margin       int // 所有边的边距
	marginTop    int
	marginRight  int
	marginBottom int
	marginLeft   int
}

var _ Element = (*Box)(nil)

// NewBox 创建一个新的 box 元素。
func NewBox(child Element) *Box {
	return &Box{
		child:  child,
		border: BorderNone,
	}
}

// Border 设置边框样式并返回 box 以支持链式调用。
func (b *Box) Border(border string) *Box {
	b.border = border
	return b
}

// BorderColor 设置边框颜色并返回 box 以支持链式调用。
func (b *Box) BorderColor(c color.Color) *Box {
	b.borderColor = c
	return b
}

// Padding 设置内边距并返回 box 以支持链式调用。
func (b *Box) Padding(padding int) *Box {
	b.padding = padding
	return b
}

// Margin 设置所有边的边距并返回 box 以支持链式调用。
func (b *Box) Margin(margin int) *Box {
	b.margin = margin
	return b
}

// MarginTop 设置顶部边距并返回 box 以支持链式调用。
func (b *Box) MarginTop(margin int) *Box {
	b.marginTop = margin
	return b
}

// MarginRight 设置右侧边距并返回 box 以支持链式调用。
func (b *Box) MarginRight(margin int) *Box {
	b.marginRight = margin
	return b
}

// MarginBottom 设置底部边距并返回 box 以支持链式调用。
func (b *Box) MarginBottom(margin int) *Box {
	b.marginBottom = margin
	return b
}

// MarginLeft 设置左侧边距并返回 box 以支持链式调用。
func (b *Box) MarginLeft(margin int) *Box {
	b.marginLeft = margin
	return b
}

// Width 设置宽度约束并返回 box 以支持链式调用。
func (b *Box) Width(width SizeConstraint) *Box {
	b.width = width
	return b
}

// Height 设置高度约束并返回 box 以支持链式调用。
func (b *Box) Height(height SizeConstraint) *Box {
	b.height = height
	return b
}

// Draw 将 box 渲染到屏幕上。
func (b *Box) Draw(scr uv.Screen, area uv.Rectangle) {
	b.SetBounds(area)

	// 应用边距（在绘制前缩小区域）
	marginTop := b.marginTop
	if marginTop == 0 {
		marginTop = b.margin
	}
	marginRight := b.marginRight
	if marginRight == 0 {
		marginRight = b.margin
	}
	marginBottom := b.marginBottom
	if marginBottom == 0 {
		marginBottom = b.margin
	}
	marginLeft := b.marginLeft
	if marginLeft == 0 {
		marginLeft = b.margin
	}

	marginH := marginLeft + marginRight
	marginV := marginTop + marginBottom

	if area.Dx() > marginH && area.Dy() > marginV {
		area = uv.Rect(
			area.Min.X+marginLeft,
			area.Min.Y+marginTop,
			area.Dx()-marginH,
			area.Dy()-marginV,
		)
	}

	// 如果指定了边框，则绘制边框
	if b.border != "" && b.border != BorderNone {
		var uvBorder uv.Border
		switch b.border {
		case BorderNormal:
			uvBorder = uv.NormalBorder()
		case BorderRounded:
			uvBorder = uv.RoundedBorder()
		case BorderThick:
			uvBorder = uv.ThickBorder()
		case BorderDouble:
			uvBorder = uv.DoubleBorder()
		case BorderHidden:
			uvBorder = uv.HiddenBorder()
		default:
			uvBorder = uv.NormalBorder()
		}

		// 如果指定了边框颜色，则应用
		if b.borderColor != nil {
			uvBorder = uvBorder.Style(uv.Style{Fg: b.borderColor})
		}

		uvBorder.Draw(scr, area)

		// 为子内容缩小区域（为边框留出空间）
		if area.Dx() > 2 && area.Dy() > 2 {
			area = uv.Rect(area.Min.X+1, area.Min.Y+1, area.Dx()-2, area.Dy()-2)
		}
	}

	// 应用内边距
	if b.padding > 0 {
		padH := b.padding * 2 // 左侧 + 右侧
		padV := b.padding * 2 // 顶部 + 底部
		if area.Dx() > padH && area.Dy() > padV {
			area = uv.Rect(
				area.Min.X+b.padding,
				area.Min.Y+b.padding,
				area.Dx()-padH,
				area.Dy()-padV,
			)
		}
	}

	// 如果存在子元素，则绘制子元素
	if b.child != nil {
		b.child.Draw(scr, area)
	}
}

// Layout 计算 box 的大小。
func (b *Box) Layout(constraints Constraints) Size {
	// 考虑边距
	marginTop := b.marginTop
	if marginTop == 0 {
		marginTop = b.margin
	}
	marginRight := b.marginRight
	if marginRight == 0 {
		marginRight = b.margin
	}
	marginBottom := b.marginBottom
	if marginBottom == 0 {
		marginBottom = b.margin
	}
	marginLeft := b.marginLeft
	if marginLeft == 0 {
		marginLeft = b.margin
	}

	marginWidth := marginLeft + marginRight
	marginHeight := marginTop + marginBottom

	// 考虑边框
	borderWidth := 0
	borderHeight := 0
	if b.border != "" && b.border != BorderNone {
		borderWidth = 2
		borderHeight = 2
	}

	// 考虑内边距
	paddingWidth := b.padding * 2
	paddingHeight := b.padding * 2

	totalReduction := marginWidth + borderWidth + paddingWidth
	totalReductionH := marginHeight + borderHeight + paddingHeight

	childConstraints := Constraints{
		MinWidth:  max(0, constraints.MinWidth-totalReduction),
		MaxWidth:  max(0, constraints.MaxWidth-totalReduction),
		MinHeight: max(0, constraints.MinHeight-totalReductionH),
		MaxHeight: max(0, constraints.MaxHeight-totalReductionH),
	}

	var childSize Size
	if b.child != nil {
		childSize = b.child.Layout(childConstraints)
	}

	totalSize := Size{
		Width:  childSize.Width + marginWidth + borderWidth + paddingWidth,
		Height: childSize.Height + marginHeight + borderHeight + paddingHeight,
	}

	// 如果指定了宽度约束，则应用
	if !b.width.IsAuto() {
		totalSize.Width = b.width.Apply(constraints.MaxWidth, totalSize.Width)
	}

	// 如果指定了高度约束，则应用
	if !b.height.IsAuto() {
		totalSize.Height = b.height.Apply(constraints.MaxHeight, totalSize.Height)
	}

	return constraints.Constrain(totalSize)
}

// Children 返回子元素。
func (b *Box) Children() []Element {
	if b.child == nil {
		return nil
	}
	return []Element{b.child}
}
