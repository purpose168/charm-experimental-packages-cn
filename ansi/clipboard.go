package ansi

import "encoding/base64"

// 剪贴板名称。
const (
	SystemClipboard  = 'c'
	PrimaryClipboard = 'p'
)

// SetClipboard 返回一个用于操作剪贴板的序列。
//
//	OSC 52 ; Pc ; Pd ST
//	OSC 52 ; Pc ; Pd BEL
//
// 其中 Pc 是剪贴板名称，Pd 是 base64 编码的数据。
// 空数据或无效的 base64 数据将重置剪贴板。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func SetClipboard(c byte, d string) string {
	if d != "" {
		d = base64.StdEncoding.EncodeToString([]byte(d))
	}
	return "\x1b]52;" + string(c) + ";" + d + "\x07"
}

// SetSystemClipboard 返回一个用于设置系统剪贴板的序列。
//
// 这等同于 SetClipboard(SystemClipboard, d)。
func SetSystemClipboard(d string) string {
	return SetClipboard(SystemClipboard, d)
}

// SetPrimaryClipboard 返回一个用于设置主剪贴板的序列。
//
// 这等同于 SetClipboard(PrimaryClipboard, d)。
func SetPrimaryClipboard(d string) string {
	return SetClipboard(PrimaryClipboard, d)
}

// ResetClipboard 返回一个用于重置剪贴板的序列。
//
// 这等同于 SetClipboard(c, "")。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func ResetClipboard(c byte) string {
	return SetClipboard(c, "")
}

// ResetSystemClipboard 是一个用于重置系统剪贴板的序列。
//
// 这等同于 ResetClipboard(SystemClipboard)。
const ResetSystemClipboard = "\x1b]52;c;\x07"

// ResetPrimaryClipboard 是一个用于重置主剪贴板的序列。
//
// 这等同于 ResetClipboard(PrimaryClipboard)。
const ResetPrimaryClipboard = "\x1b]52;p;\x07"

// RequestClipboard 返回一个用于请求剪贴板的序列。
//
// 参见: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func RequestClipboard(c byte) string {
	return "\x1b]52;" + string(c) + ";?\x07"
}

// RequestSystemClipboard 是一个用于请求系统剪贴板的序列。
//
// 这等同于 RequestClipboard(SystemClipboard)。
const RequestSystemClipboard = "\x1b]52;c;?\x07"

// RequestPrimaryClipboard 是一个用于请求主剪贴板的序列。
//
// 这等同于 RequestClipboard(PrimaryClipboard)。
const RequestPrimaryClipboard = "\x1b]52;p;?\x07"
