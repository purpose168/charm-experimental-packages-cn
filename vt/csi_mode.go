package vt

import (
	"io"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// handleMode 处理模式设置参数，根据 set 标志决定是设置还是重置模式。
// isAnsi 标志指示是否为 ANSI 模式。
func (e *Emulator) handleMode(params ansi.Params, set, isAnsi bool) {
	for _, p := range params {
		param := p.Param(-1)
		if param == -1 {
			// 缺少参数，忽略
			continue
		}

		var mode ansi.Mode = ansi.DECMode(param)
		if isAnsi {
			mode = ansi.ANSIMode(param)
		}

		setting := e.modes[mode]
		if setting == ansi.ModePermanentlyReset || setting == ansi.ModePermanentlySet {
			// 永久设置的模式被忽略。
			continue
		}

		setting = ansi.ModeReset
		if set {
			setting = ansi.ModeSet
		}

		e.setMode(mode, setting)
	}
}

// setAltScreenMode 设置交替屏幕模式。
func (e *Emulator) setAltScreenMode(on bool) {
	if (on && e.scr == &e.scrs[1]) || (!on && e.scr == &e.scrs[0]) {
		// 已经在交替屏幕模式或正常屏幕，不执行任何操作。
		return
	}
	if on {
		e.scr = &e.scrs[1]
		e.scrs[1].cur = e.scrs[0].cur
		e.scr.Clear()
		e.scr.buf.Touched = nil
		e.setCursor(0, 0)
	} else {
		e.scr = &e.scrs[0]
	}
	if e.cb.AltScreen != nil {
		e.cb.AltScreen(on)
	}
	if e.cb.CursorVisibility != nil {
		e.cb.CursorVisibility(!e.scr.cur.Hidden)
	}
}

// saveCursor 保存光标位置。
func (e *Emulator) saveCursor() {
	e.scr.SaveCursor()
}

// restoreCursor 恢复光标位置。
func (e *Emulator) restoreCursor() {
	e.scr.RestoreCursor()
}

// setMode 将模式设置为给定值。
func (e *Emulator) setMode(mode ansi.Mode, setting ansi.ModeSetting) {
	e.logf("setting mode %T(%v) to %v", mode, mode, setting)
	e.modes[mode] = setting
	switch mode {
	case ansi.TextCursorEnableMode:
		e.scr.setCursorHidden(!setting.IsSet())
	case ansi.AltScreenMode:
		e.setAltScreenMode(setting.IsSet())
	case ansi.SaveCursorMode:
		if setting.IsSet() {
			e.saveCursor()
		} else {
			e.restoreCursor()
		}
	case ansi.AltScreenSaveCursorMode: // 交替屏幕保存光标 (1047 & 1048)
		// 保存主屏幕光标位置
		// 切换到交替屏幕
		// 不支持滚动回退
		if setting.IsSet() {
			e.saveCursor()
		}
		e.setAltScreenMode(setting.IsSet())
	case ansi.InBandResizeMode:
		if setting.IsSet() {
			_, _ = io.WriteString(e.pw, ansi.InBandResize(e.Height(), e.Width(), 0, 0))
		}
	}
	if setting.IsSet() {
		if e.cb.EnableMode != nil {
			e.cb.EnableMode(mode)
		}
	} else if setting.IsReset() {
		if e.cb.DisableMode != nil {
			e.cb.DisableMode(mode)
		}
	}
}

// isModeSet 如果模式已设置，则返回 true。
func (e *Emulator) isModeSet(mode ansi.Mode) bool {
	m, ok := e.modes[mode]
	return ok && m.IsSet()
}
