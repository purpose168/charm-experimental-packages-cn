package vt

import (
	"io"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// Focus 如果启用了焦点事件模式，则向终端发送焦点事件。
// 这与 [Blur] 相反。
func (e *Emulator) Focus() {
	e.focus(true)
}

// Blur 如果启用了焦点事件模式，则向终端发送失焦事件。
// 这与 [Focus] 相反。
func (e *Emulator) Blur() {
	e.focus(false)
}

func (e *Emulator) focus(focus bool) {
	if mode, ok := e.modes[ansi.FocusEventMode]; ok && mode.IsSet() {
		if focus {
			_, _ = io.WriteString(e.pw, ansi.Focus)
		} else {
			_, _ = io.WriteString(e.pw, ansi.Blur)
		}
	}
}
