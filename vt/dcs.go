package vt

import "github.com/purpose168/charm-experimental-packages-cn/ansi"

// handleDcs 处理 DCS（设备控制字符串）转义序列。
func (e *Emulator) handleDcs(cmd ansi.Cmd, params ansi.Params, data []byte) {
	e.flushGrapheme() // 在处理 DCS 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handleDcs(cmd, params, data) {
		e.logf("未处理的序列: DCS %q %q", paramsString(cmd, params), data)
	}
}

// handleApc 处理 APC（应用程序程序命令）转义序列。
func (e *Emulator) handleApc(data []byte) {
	e.flushGrapheme() // 在处理 APC 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handleApc(data) {
		e.logf("未处理的序列: APC %q", data)
	}
}

// handleSos 处理 SOS（字符串开始）转义序列。
func (e *Emulator) handleSos(data []byte) {
	e.flushGrapheme() // 在处理 SOS 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handleSos(data) {
		e.logf("未处理的序列: SOS %q", data)
	}
}

// handlePm 处理 PM（隐私消息）转义序列。
func (e *Emulator) handlePm(data []byte) {
	e.flushGrapheme() // 在处理 PM 序列之前，先刷新任何待处理的字形。
	if !e.handlers.handlePm(data) {
		e.logf("未处理的序列: PM %q", data)
	}
}
