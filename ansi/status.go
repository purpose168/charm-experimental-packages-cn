package ansi

import (
	"strconv"
	"strings"
)

// StatusReport 表示终端状态报告。
type StatusReport interface {
	// StatusReport 返回状态报告标识符。
	StatusReport() int
}

// ANSIStatusReport 表示ANSI终端状态报告。
type ANSIStatusReport int //nolint:revive

// StatusReport 返回状态报告标识符。
func (s ANSIStatusReport) StatusReport() int {
	return int(s)
}

// DECStatusReport 表示DEC终端状态报告。
type DECStatusReport int

// StatusReport 返回状态报告标识符。
func (s DECStatusReport) StatusReport() int {
	return int(s)
}

// DeviceStatusReport (DSR) 是一个报告终端状态的控制序列。
// 终端会用DSR序列响应。
//
//	CSI Ps n
//	CSI ? Ps n
//
// 如果其中一个状态是[DECStatus]，序列将使用DEC格式。
//
// 另请参阅 https://vt100.net/docs/vt510-rm/DSR.html
func DeviceStatusReport(statues ...StatusReport) string {
	var dec bool
	list := make([]string, len(statues))
	seq := "\x1b["
	for i, status := range statues {
		list[i] = strconv.Itoa(status.StatusReport())
		switch status.(type) {
		case DECStatusReport:
			dec = true
		}
	}
	if dec {
		seq += "?"
	}
	return seq + strings.Join(list, ";") + "n"
}

// DSR 是[DeviceStatusReport]的别名。
func DSR(status StatusReport) string {
	return DeviceStatusReport(status)
}

// RequestCursorPositionReport 是一个请求当前光标位置的转义序列。
//
//	CSI 6 n
//
// 终端将以以下格式的CSI序列报告光标位置：
//
//	CSI Pl ; Pc R
//
// 其中Pl是行号，Pc是列号。
// 请参阅：https://vt100.net/docs/vt510-rm/CPR.html
const RequestCursorPositionReport = "\x1b[6n"

// RequestExtendedCursorPositionReport (DECXCPR) 是一个请求光标位置报告的序列，包括当前页码。
//
//	CSI ? 6 n
//
// 终端将以以下格式的CSI序列报告光标位置：
//
//	CSI ? Pl ; Pc ; Pp R
//
// 其中Pl是行号，Pc是列号，Pp是页码。
// 请参阅：https://vt100.net/docs/vt510-rm/DECXCPR.html
const RequestExtendedCursorPositionReport = "\x1b[?6n"

// RequestLightDarkReport 是一个请求终端报告其操作系统明暗颜色偏好的控制序列。
// 支持的终端应使用以下[LightDarkReport]序列响应：
//
//	CSI ? 997 ; 1 n   暗色模式
//	CSI ? 997 ; 2 n   亮色模式
//
// 请参阅：https://contour-terminal.org/vt-extensions/color-palette-update-notifications/
const RequestLightDarkReport = "\x1b[?996n"

// CursorPositionReport (CPR) 是一个报告光标位置的控制序列。
//
//	CSI Pl ; Pc R
//
// 其中Pl是行号，Pc是列号。
//
// 另请参阅 https://vt100.net/docs/vt510-rm/CPR.html
func CursorPositionReport(line, column int) string {
	if line < 1 {
		line = 1
	}
	if column < 1 {
		column = 1
	}
	return "\x1b[" + strconv.Itoa(line) + ";" + strconv.Itoa(column) + "R"
}

// CPR 是[CursorPositionReport]的别名。
func CPR(line, column int) string {
	return CursorPositionReport(line, column)
}

// ExtendedCursorPositionReport (DECXCPR) 是一个报告光标位置（可选包含页码）的控制序列。
//
//	CSI ? Pl ; Pc R
//	CSI ? Pl ; Pc ; Pv R
//
// 其中Pl是行号，Pc是列号，Pv是页码。
//
// 如果页码为零或负数，返回的序列将不包含页码。
//
// 另请参阅 https://vt100.net/docs/vt510-rm/DECXCPR.html
func ExtendedCursorPositionReport(line, column, page int) string {
	if line < 1 {
		line = 1
	}
	if column < 1 {
		column = 1
	}
	if page < 1 {
		return "\x1b[?" + strconv.Itoa(line) + ";" + strconv.Itoa(column) + "R"
	}
	return "\x1b[?" + strconv.Itoa(line) + ";" + strconv.Itoa(column) + ";" + strconv.Itoa(page) + "R"
}

// DECXCPR 是[ExtendedCursorPositionReport]的别名。
func DECXCPR(line, column, page int) string {
	return ExtendedCursorPositionReport(line, column, page)
}

// LightDarkReport 是一个报告终端操作系统明暗颜色偏好的控制序列。
//
//	CSI ? 997 ; 1 n   暗色模式
//	CSI ? 997 ; 2 n   亮色模式
//
// 请参阅：https://contour-terminal.org/vt-extensions/color-palette-update-notifications/
func LightDarkReport(dark bool) string {
	if dark {
		return "\x1b[?997;1n"
	}
	return "\x1b[?997;2n"
}
