package pony

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// Divider 表示一条水平或垂直线。
type Divider struct {
	BaseElement
	vertical bool
	char     string
	color    color.Color
}

var _ Element = (*Divider)(nil)

// NewDivider 创建一个新的分隔线元素。
func NewDivider() *Divider {
	return &Divider{}
}

// NewVerticalDivider 创建一个新的垂直分隔线。
func NewVerticalDivider() *Divider {
	return &Divider{vertical: true}
}

// ForegroundColor 设置颜色并返回分隔线以支持链式调用。
func (d *Divider) ForegroundColor(c color.Color) *Divider {
	d.color = c
	return d
}

// Char 设置字符并返回分隔线以支持链式调用。
func (d *Divider) Char(char string) *Divider {
	d.char = char
	return d
}

// Draw 将分隔线渲染到屏幕上。
func (d *Divider) Draw(scr uv.Screen, area uv.Rectangle) {
	d.SetBounds(area)

	char := d.char
	if char == "" {
		if d.vertical {
			char = "│"
		} else {
			char = "─"
		}
	}

	cell := uv.NewCell(scr.WidthMethod(), char)
	if cell != nil && d.color != nil {
		cell.Style = uv.Style{Fg: d.color}
	}

	if d.vertical {
		for y := area.Min.Y; y < area.Max.Y; y++ {
			scr.SetCell(area.Min.X, y, cell)
		}
	} else {
		for x := area.Min.X; x < area.Max.X; x++ {
			scr.SetCell(x, area.Min.Y, cell)
		}
	}
}

// Layout 计算分隔线的大小。
func (d *Divider) Layout(constraints Constraints) Size {
	if d.vertical {
		return Size{Width: 1, Height: constraints.MaxHeight}
	}
	return Size{Width: constraints.MaxWidth, Height: 1}
}

// Children 为分隔线返回 nil。
func (d *Divider) Children() []Element {
	return nil
}
