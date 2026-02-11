package cellbuf

import (
	"image"
)

// Position 表示一个 x, y 位置。
type Position = image.Point

// Pos 是 Position{X: x, Y: y} 的简写。
func Pos(x, y int) Position {
	return image.Pt(x, y)
}

// Rectangle 表示一个矩形。
type Rectangle = image.Rectangle

// Rect 是 Rectangle 的简写。
func Rect(x, y, w, h int) Rectangle {
	return image.Rect(x, y, x+w, y+h)
}
