package vt

import (
	"fmt"
	"io"
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// handleCsi 处理 CSI（控制序列引导器）指令。
// 在处理 CSI 序列之前，先刷新任何待处理的字形。
func (e *Emulator) handleCsi(cmd ansi.Cmd, params ansi.Params) {
	e.flushGrapheme() // 在处理 CSI 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handleCsi(cmd, params) {
		e.logf("未处理的序列: CSI %q", paramsString(cmd, params))
	}
}

// handleRequestMode 处理模式请求。
// 参数 params 包含请求的模式参数，isAnsi 指示是否为 ANSI 模式。
func (e *Emulator) handleRequestMode(params ansi.Params, isAnsi bool) {
	n, _, ok := params.Param(0, 0)
	if !ok || n == 0 {
		return
	}

	var mode ansi.Mode = ansi.DECMode(n)
	if isAnsi {
		mode = ansi.ANSIMode(n)
	}

	setting := e.modes[mode]
	_, _ = io.WriteString(e.pw, ansi.ReportMode(mode, setting))
}

// paramsString 将 CSI 命令和参数转换为字符串表示。
// 用于日志记录和调试未处理的序列。
func paramsString(cmd ansi.Cmd, params ansi.Params) string {
	var s strings.Builder
	if mark := cmd.Prefix(); mark != 0 {
		s.WriteByte(mark)
	}
	params.ForEach(-1, func(i, p int, more bool) {
		s.WriteString(fmt.Sprintf("%d", p))
		if i < len(params)-1 {
			if more {
				s.WriteByte(':')
			} else {
				s.WriteByte(';')
			}
		}
	})
	if inter := cmd.Intermediate(); inter != 0 {
		s.WriteByte(inter)
	}
	if final := cmd.Final(); final != 0 {
		s.WriteByte(final)
	}
	return s.String()
}
