package ansi

import "strconv"

// KeyModifierOptions (XTMODKEYS) 设置/重置 xterm 键修饰符选项。
//
// 默认值为 0。
//
//	CSI > Pp m
//	CSI > Pp ; Pv m
//
// 如果省略 Pv，则资源将重置为其初始值。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
func KeyModifierOptions(p int, vs ...int) string {
	var pp, pv string
	if p > 0 {
		pp = strconv.Itoa(p)
	}

	if len(vs) == 0 {
		return "\x1b[>" + strconv.Itoa(p) + "m"
	}

	v := vs[0]
	if v > 0 {
		pv = strconv.Itoa(v)
		return "\x1b[>" + pp + ";" + pv + "m"
	}

	return "\x1b[>" + pp + "m"
}

// XTMODKEYS 是 [KeyModifierOptions] 的别名
func XTMODKEYS(p int, vs ...int) string {
	return KeyModifierOptions(p, vs...)
}

// SetKeyModifierOptions 设置 xterm 键修饰符选项。
// 这是 [KeyModifierOptions] 的别名
func SetKeyModifierOptions(pp int, pv int) string {
	return KeyModifierOptions(pp, pv)
}

// ResetKeyModifierOptions 重置 xterm 键修饰符选项。
// 这是 [KeyModifierOptions] 的别名
func ResetKeyModifierOptions(pp int) string {
	return KeyModifierOptions(pp)
}

// QueryKeyModifierOptions (XTQMODKEYS) 请求 xterm 键修饰符选项。
//
// 默认值为 0。
//
//	CSI ? Pp m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
func QueryKeyModifierOptions(pp int) string {
	var p string
	if pp > 0 {
		p = strconv.Itoa(pp)
	}
	return "\x1b[?" + p + "m"
}

// XTQMODKEYS 是 [QueryKeyModifierOptions] 的别名
func XTQMODKEYS(pp int) string {
	return QueryKeyModifierOptions(pp)
}

// Modify Other Keys (modifyOtherKeys) 是 xterm 的一个功能，允许终端修改某些键的行为，
// 使其在按下时发送不同的转义序列。
//
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
const (
	// SetModifyOtherKeys1 设置 modifyOtherKeys 模式 1
	SetModifyOtherKeys1 = "\x1b[>4;1m"
	// SetModifyOtherKeys2 设置 modifyOtherKeys 模式 2
	SetModifyOtherKeys2 = "\x1b[>4;2m"
	// ResetModifyOtherKeys 重置 modifyOtherKeys 模式
	ResetModifyOtherKeys = "\x1b[>4m"
	// QueryModifyOtherKeys 查询 modifyOtherKeys 模式
	QueryModifyOtherKeys = "\x1b[?4m"
)

// ModifyOtherKeys 返回一个设置 XTerm modifyOtherKeys 模式的序列。
// mode 参数指定要设置的模式：
//
//	0: 禁用 modifyOtherKeys 模式
//	1: 启用 modifyOtherKeys 模式 1
//	2: 启用 modifyOtherKeys 模式 2
//
//	CSI > 4 ; mode m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
//
// 已废弃：请改用 [SetModifyOtherKeys1] 或 [SetModifyOtherKeys2]
func ModifyOtherKeys(mode int) string {
	return "\x1b[>4;" + strconv.Itoa(mode) + "m"
}

// DisableModifyOtherKeys 禁用 modifyOtherKeys 模式
//
//	CSI > 4 ; 0 m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
//
// 已废弃：请改用 [ResetModifyOtherKeys]
const DisableModifyOtherKeys = "\x1b[>4;0m"

// EnableModifyOtherKeys1 启用 modifyOtherKeys 模式 1
//
//	CSI > 4 ; 1 m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
//
// 已废弃：请改用 [SetModifyOtherKeys1]
const EnableModifyOtherKeys1 = "\x1b[>4;1m"

// EnableModifyOtherKeys2 启用 modifyOtherKeys 模式 2
//
//	CSI > 4 ; 2 m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
//
// 已废弃：请改用 [SetModifyOtherKeys2]
const EnableModifyOtherKeys2 = "\x1b[>4;2m"

// RequestModifyOtherKeys 请求 modifyOtherKeys 模式
//
//	CSI ? 4  m
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
// 参见：https://invisible-island.net/xterm/manpage/xterm.html#VT100-Widget-Resources:modifyOtherKeys
//
// 已废弃：请改用 [QueryModifyOtherKeys]
const RequestModifyOtherKeys = "\x1b[?4m"
