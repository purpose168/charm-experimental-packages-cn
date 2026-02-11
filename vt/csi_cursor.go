package vt

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// nextTab 将光标移动到下一个制表位，重复 n 次。这会尊重水平滚动区域。
// 执行与 [ansi.CHT] 相同的功能。
func (e *Emulator) nextTab(n int) {
	x, y := e.scr.CursorPosition()
	scroll := e.scr.ScrollRegion()
	for range n {
		ts := e.tabstops.Next(x)
		if ts < x {
			break
		}
		x = ts
	}

	if x >= scroll.Max.X {
		x = min(scroll.Max.X-1, e.Width()-1)
	}

	// 注意：我们使用 t.scr.setCursor 因为我们不想重置幻影状态。
	e.scr.setCursor(x, y, false)
}

// prevTab 将光标移动到上一个制表位，重复 n 次。当启用原点模式时，
// 这会尊重水平滚动区域。如果光标会移动过最左有效列，光标会停留在最左有效列并完成操作。
func (e *Emulator) prevTab(n int) {
	x, _ := e.scr.CursorPosition()
	leftmargin := 0
	scroll := e.scr.ScrollRegion()
	if e.isModeSet(ansi.DECOM) {
		leftmargin = scroll.Min.X
	}

	for range n {
		ts := e.tabstops.Prev(x)
		if ts > x {
			break
		}
		x = ts
	}

	if x < leftmargin {
		x = leftmargin
	}

	// 注意：我们使用 t.scr.setCursorX 因为我们不想重置幻影状态。
	e.scr.setCursorX(x, false)
}

// moveCursor 按给定的 x 和 y 增量移动光标。如果光标处于幻影状态，
// 状态将重置，光标回到屏幕中。
func (e *Emulator) moveCursor(dx, dy int) {
	e.scr.moveCursor(dx, dy)
	e.atPhantom = false
}

// setCursor 设置光标位置。这会重置幻影状态。
func (e *Emulator) setCursor(x, y int) {
	e.scr.setCursor(x, y, false)
	e.atPhantom = false
}

// setCursorPosition 设置光标位置。这会尊重 [ansi.DECOM]（原点模式）。
// 执行与 [ansi.CUP] 相同的功能。
func (e *Emulator) setCursorPosition(x, y int) {
	mode, ok := e.modes[ansi.DECOM]
	margins := ok && mode.IsSet()
	e.scr.setCursor(x, y, margins)
	e.atPhantom = false
}

// carriageReturn 将光标移动到最左列。如果设置了 [ansi.DECOM]，
// 光标会设置到左边界。如果没有，且光标在左边界或其右侧，
// 光标会设置到左边界。否则，光标会设置到屏幕的最左列。
// 执行与 [ansi.CR] 相同的功能。
func (e *Emulator) carriageReturn() {
	mode, ok := e.modes[ansi.DECOM]
	margins := ok && mode.IsSet()
	x, y := e.scr.CursorPosition()
	if margins {
		e.scr.setCursor(0, y, true)
	} else if region := e.scr.ScrollRegion(); uv.Pos(x, y).In(region) {
		e.scr.setCursor(region.Min.X, y, false)
	} else {
		e.scr.setCursor(0, y, false)
	}
	e.atPhantom = false
}

// repeatPreviousCharacter 重复前一个字符 n 次。这相当于键入相同的字符 n 次。
// 执行与 [ansi.REP] 相同的功能。
func (e *Emulator) repeatPreviousCharacter(n int) {
	if e.lastChar == 0 {
		return
	}
	for range n {
		e.handlePrint(e.lastChar)
	}
}
