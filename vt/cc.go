package vt

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// handleControl 处理控制字符。
func (e *Emulator) handleControl(r byte) {
	e.flushGrapheme() // 在处理控制代码之前，先刷新任何待处理的字形。
	if !e.handleCc(r) {
		e.logf("未处理的序列: ControlCode %q", r)
	}
}

// linefeed 与 [index] 相同，但会尊重 [ansi.LNM] 模式。
func (e *Emulator) linefeed() {
	e.index()
	if e.isModeSet(ansi.LineFeedNewLineMode) {
		e.carriageReturn()
	}
}

// index 将光标向下移动一行，必要时向上滚动。这总是会重置幻影状态，即待换行状态。
func (e *Emulator) index() {
	x, y := e.scr.CursorPosition()
	scroll := e.scr.ScrollRegion()
	// XXX: 当我们添加滚动回退时处理它。
	if y == scroll.Max.Y-1 && x >= scroll.Min.X && x < scroll.Max.X {
		e.scr.ScrollUp(1)
	} else if y < scroll.Max.Y-1 || !uv.Pos(x, y).In(scroll) {
		e.scr.moveCursor(0, 1)
	}
	e.atPhantom = false
}

// horizontalTabSet 在当前光标位置设置水平制表位。
func (e *Emulator) horizontalTabSet() {
	x, _ := e.scr.CursorPosition()
	e.tabstops.Set(x)
}

// reverseIndex 将光标向上移动一行，或向下滚动。这不会重置幻影状态，即待换行状态。
func (e *Emulator) reverseIndex() {
	x, y := e.scr.CursorPosition()
	scroll := e.scr.ScrollRegion()
	if y == scroll.Min.Y && x >= scroll.Min.X && x < scroll.Max.X {
		e.scr.ScrollDown(1)
	} else {
		e.scr.moveCursor(0, -1)
	}
}

// backspace 如果可能，将光标向后移动一个单元格。
func (e *Emulator) backspace() {
	// 这类似于 [ansi.CUB]
	e.moveCursor(-1, 0)
}
