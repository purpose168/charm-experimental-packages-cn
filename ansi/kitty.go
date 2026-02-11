package ansi

import "strconv"

// Kitty 键盘协议渐进增强标志。
// 参见: https://sw.kovidgoyal.net/kitty/keyboard-protocol/#progressive-enhancement
const (
	KittyDisambiguateEscapeCodes = 1 << iota
	KittyReportEventTypes
	KittyReportAlternateKeys
	KittyReportAllKeysAsEscapeCodes
	KittyReportAssociatedKeys

	KittyAllFlags = KittyDisambiguateEscapeCodes | KittyReportEventTypes |
		KittyReportAlternateKeys | KittyReportAllKeysAsEscapeCodes | KittyReportAssociatedKeys
)

// RequestKittyKeyboard 是一个请求启用终端 Kitty 键盘协议标志的序列。
//
// 参见: https://sw.kovidgoyal.net/kitty/keyboard-protocol/
const RequestKittyKeyboard = "\x1b[?u"

// KittyKeyboard 返回一个从终端请求键盘增强功能的序列。
// flags 参数是一个位掩码，表示 Kitty 键盘协议的标志。
// mode 指定如何解释这些标志。
//
// 标志掩码的可能的值：
//
//	1:  消除转义码歧义
//	2:  报告事件类型
//	4:  报告备用键
//	8:  将所有键报告为转义码
//	16: 报告关联文本
//
// 模式的可能的值：
//
//	1: 设置给定标志并清除所有其他标志
//	2: 设置给定标志并保持现有标志不变
//	3: 清除给定标志并保持现有标志不变
//
// 参见 https://sw.kovidgoyal.net/kitty/keyboard-protocol/#progressive-enhancement
func KittyKeyboard(flags, mode int) string {
	return "\x1b[=" + strconv.Itoa(flags) + ";" + strconv.Itoa(mode) + "u"
}

// PushKittyKeyboard 返回一个将给定标志压入终端 Kitty 键盘堆栈的序列。
//
// 标志掩码的可能的值：
//
//	0:  禁用所有功能
//	1:  消除转义码歧义
//	2:  报告事件类型
//	4:  报告备用键
//	8:  将所有键报告为转义码
//	16: 报告关联文本
//
//	CSI > flags u
//
// 参见 https://sw.kovidgoyal.net/kitty/keyboard-protocol/#progressive-enhancement
func PushKittyKeyboard(flags int) string {
	var f string
	if flags > 0 {
		f = strconv.Itoa(flags)
	}

	return "\x1b[>" + f + "u"
}

// DisableKittyKeyboard 是一个将零压入终端 Kitty 键盘堆栈以禁用该协议的序列。
//
// 这等同于 PushKittyKeyboard(0)。
const DisableKittyKeyboard = "\x1b[>u"

// PopKittyKeyboard 返回一个从终端 Kitty 键盘堆栈弹出一个或多个标志的序列。
//
//	CSI < flags u
//
// 参见 https://sw.kovidgoyal.net/kitty/keyboard-protocol/#progressive-enhancement
func PopKittyKeyboard(n int) string {
	var num string
	if n > 0 {
		num = strconv.Itoa(n)
	}

	return "\x1b[<" + num + "u"
}
