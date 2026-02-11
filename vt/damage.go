package vt

import uv "github.com/charmbracelet/ultraviolet"

// Damage 表示一个损坏的区域。
type Damage interface {
	// Bounds 返回损坏区域的边界。
	Bounds() uv.Rectangle
}

// CellDamage 表示一个损坏的单元格。
type CellDamage struct {
	X, Y  int
	Width int
}

// Bounds 返回损坏区域的边界。
func (d CellDamage) Bounds() uv.Rectangle {
	return uv.Rect(d.X, d.Y, d.Width, 1)
}

// RectDamage 表示一个损坏的矩形。
type RectDamage uv.Rectangle

// Bounds 返回损坏区域的边界。
func (d RectDamage) Bounds() uv.Rectangle {
	return uv.Rectangle(d)
}

// X 返回损坏区域的 x 坐标。
func (d RectDamage) X() int {
	return uv.Rectangle(d).Min.X
}

// Y 返回损坏区域的 y 坐标。
func (d RectDamage) Y() int {
	return uv.Rectangle(d).Min.Y
}

// Width 返回损坏区域的宽度。
func (d RectDamage) Width() int {
	return uv.Rectangle(d).Dx()
}

// Height 返回损坏区域的高度。
func (d RectDamage) Height() int {
	return uv.Rectangle(d).Dy()
}

// ScreenDamage 表示一个损坏的屏幕。
type ScreenDamage struct {
	Width, Height int
}

// Bounds 返回损坏区域的边界。
func (d ScreenDamage) Bounds() uv.Rectangle {
	return uv.Rect(0, 0, d.Width, d.Height)
}

// MoveDamage 表示一个移动的区域。
// 该区域从源位置移动到目标位置。
type MoveDamage struct {
	Src, Dst uv.Rectangle
}

// ScrollDamage 表示一个滚动的区域。
// 该区域按给定的增量滚动。
type ScrollDamage struct {
	uv.Rectangle
	Dx, Dy int
}
