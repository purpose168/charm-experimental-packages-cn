package vt

import (
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// handleEsc 处理 ESC 转义序列。
func (e *Emulator) handleEsc(cmd ansi.Cmd) {
	e.flushGrapheme() // 在处理 ESC 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handleEsc(int(cmd)) {
		var str string
		if inter := cmd.Intermediate(); inter != 0 {
			str += string(inter) + " "
		}
		if final := cmd.Final(); final != 0 {
			str += string(final)
		}
		e.logf("未处理的序列: ESC %q", str)
	}
}

// fullReset 执行完整的终端重置，如 [ansi.RIS] 中所述。
func (e *Emulator) fullReset() {
	e.scrs[0].Reset()
	e.scrs[1].Reset()
	e.resetTabStops()

	// XXX: 我们是否在这里重置所有模式？需要调查。
	e.resetModes()

	e.gl, e.gr = 0, 1
	e.gsingle = 0
	e.charsets = [4]CharSet{}
	e.atPhantom = false
	e.grapheme = e.grapheme[:0]
	e.lastChar = 0
	e.lastState = parser.GroundState
}
