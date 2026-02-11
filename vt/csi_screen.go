package vt

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// eraseCharacter 从光标位置开始擦除 n 个字符。它不会移动光标。
// 这相当于 [ansi.ECH]。
func (e *Emulator) eraseCharacter(n int) {
	if n <= 0 {
		n = 1
	}
	x, y := e.scr.CursorPosition()
	rect := uv.Rect(x, y, n, 1)
	e.scr.FillArea(e.scr.blankCell(), rect)
	e.atPhantom = false
	// ECH 不会移动光标。
}
