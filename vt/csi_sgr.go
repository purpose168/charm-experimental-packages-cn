package vt

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// handleSgr 处理 SGR（选择图形渲染）转义序列。
func (e *Emulator) handleSgr(params ansi.Params) {
	uv.ReadStyle(params, &e.scr.cur.Pen)
}
