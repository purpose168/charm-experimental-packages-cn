package ansi

import (
	"strconv"
	"strings"
)

// EraseDisplay（ED）清除整个显示或显示的部分内容。屏幕是终端显示的可见部分，不包括回滚缓冲区。
// 可能的值：
//
// 默认值为 0。
//
//	 0: 清除从光标到屏幕末尾的内容。
//	 1: 清除从光标到屏幕开始的内容。
//	 2: 清除整个屏幕（并在 DOS 上将光标移到左上角）。
//	 3: 清除整个显示，删除回滚缓冲区中保存的所有行（xterm）。
//
//	CSI <n> J
//
// 参考：https://vt100.net/docs/vt510-rm/ED.html
func EraseDisplay(n int) string {
	var s string
	if n > 0 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "J"
}

// ED 是 [EraseDisplay] 的别名。
func ED(n int) string {
	return EraseDisplay(n)
}

// EraseDisplay 常量。
// 这些是 EraseDisplay 函数的可能值。
const (
	EraseScreenBelow   = "\x1b[J"
	EraseScreenAbove   = "\x1b[1J"
	EraseEntireScreen  = "\x1b[2J"
	EraseEntireDisplay = "\x1b[3J"
)

// EraseLine（EL）清除当前行或行的部分内容。可能的值：
//
//	0: 清除从光标到行尾的内容。
//	1: 清除从光标到行首的内容。
//	2: 清除整行内容。
//
// 光标位置不受影响。
//
//	CSI <n> K
//
// 参考：https://vt100.net/docs/vt510-rm/EL.html
func EraseLine(n int) string {
	var s string
	if n > 0 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "K"
}

// EL 是 [EraseLine] 的别名。
func EL(n int) string {
	return EraseLine(n)
}

// EraseLine 常量。
// 这些是 EraseLine 函数的可能值。
const (
	EraseLineRight  = "\x1b[K"
	EraseLineLeft   = "\x1b[1K"
	EraseEntireLine = "\x1b[2K"
)

// ScrollUp（SU）向上滚动屏幕 n 行。新行添加到屏幕底部。
//
//	CSI Pn S
//
// 参考：https://vt100.net/docs/vt510-rm/SU.html
func ScrollUp(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "S"
}

// PanDown 是 [ScrollUp] 的别名。
func PanDown(n int) string {
	return ScrollUp(n)
}

// SU 是 [ScrollUp] 的别名。
func SU(n int) string {
	return ScrollUp(n)
}

// ScrollDown（SD）向下滚动屏幕 n 行。新行添加到屏幕顶部。
//
//	CSI Pn T
//
// 参考：https://vt100.net/docs/vt510-rm/SD.html
func ScrollDown(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "T"
}

// PanUp 是 [ScrollDown] 的别名。
func PanUp(n int) string {
	return ScrollDown(n)
}

// SD 是 [ScrollDown] 的别名。
func SD(n int) string {
	return ScrollDown(n)
}

// InsertLine（IL）在当前光标位置插入 n 个空行。现有行向下移动。
//
//	CSI Pn L
//
// 参考：https://vt100.net/docs/vt510-rm/IL.html
func InsertLine(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "L"
}

// IL 是 [InsertLine] 的别名。
func IL(n int) string {
	return InsertLine(n)
}

// DeleteLine（DL）删除当前光标位置的 n 行。现有行向上移动。
//
//	CSI Pn M
//
// 参考：https://vt100.net/docs/vt510-rm/DL.html
func DeleteLine(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "M"
}

// DL 是 [DeleteLine] 的别名。
func DL(n int) string {
	return DeleteLine(n)
}

// SetTopBottomMargins（DECSTBM）设置滚动区域的顶部和底部边距。默认值是整个屏幕。
//
// 默认值为 1 和屏幕底部。
//
//	CSI Pt ; Pb r
//
// 参考：https://vt100.net/docs/vt510-rm/DECSTBM.html
func SetTopBottomMargins(top, bot int) string {
	var t, b string
	if top > 0 {
		t = strconv.Itoa(top)
	}
	if bot > 0 {
		b = strconv.Itoa(bot)
	}
	return "\x1b[" + t + ";" + b + "r"
}

// DECSTBM 是 [SetTopBottomMargins] 的别名。
func DECSTBM(top, bot int) string {
	return SetTopBottomMargins(top, bot)
}

// SetLeftRightMargins（DECSLRM）设置滚动区域的左侧和右侧边距。
//
// 默认值为 1 和屏幕右侧。
//
//	CSI Pl ; Pr s
//
// 参考：https://vt100.net/docs/vt510-rm/DECSLRM.html
func SetLeftRightMargins(left, right int) string {
	var l, r string
	if left > 0 {
		l = strconv.Itoa(left)
	}
	if right > 0 {
		r = strconv.Itoa(right)
	}
	return "\x1b[" + l + ";" + r + "s"
}

// DECSLRM 是 [SetLeftRightMargins] 的别名。
func DECSLRM(left, right int) string {
	return SetLeftRightMargins(left, right)
}

// SetScrollingRegion（DECSTBM）设置滚动区域的顶部和底部边距。默认值是整个屏幕。
//
//	CSI <top> ; <bottom> r
//
// 参考：https://vt100.net/docs/vt510-rm/DECSTBM.html
//
// 已弃用：请使用 [SetTopBottomMargins] 代替。
func SetScrollingRegion(t, b int) string {
	if t < 0 {
		t = 0
	}
	if b < 0 {
		b = 0
	}
	return "\x1b[" + strconv.Itoa(t) + ";" + strconv.Itoa(b) + "r"
}

// InsertCharacter（ICH）在当前光标位置插入 n 个空字符。现有字符向右移动。移动超过右侧边距的字符将丢失。ICH 在滚动边距外无效。
//
// 默认值为 1。
//
//	CSI Pn @
//
// 参考：https://vt100.net/docs/vt510-rm/ICH.html
func InsertCharacter(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "@"
}

// ICH 是 [InsertCharacter] 的别名。
func ICH(n int) string {
	return InsertCharacter(n)
}

// DeleteCharacter（DCH）删除当前光标位置的 n 个字符。
// 当字符被删除时，剩余字符向左移动，光标保持在相同位置。
//
// 默认值为 1。
//
//	CSI Pn P
//
// 参考：https://vt100.net/docs/vt510-rm/DCH.html
func DeleteCharacter(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "P"
}

// DCH 是 [DeleteCharacter] 的别名。
func DCH(n int) string {
	return DeleteCharacter(n)
}

// SetTabEvery8Columns（DECST8C）在每 8 列设置制表位。
//
//	CSI ? 5 W
//
// 参考：https://vt100.net/docs/vt510-rm/DECST8C.html
const (
	SetTabEvery8Columns = "\x1b[?5W"
	DECST8C             = SetTabEvery8Columns
)

// HorizontalTabSet（HTS）在当前光标列设置水平制表位。
//
// 这等同于 [HTS]。
//
//	ESC H
//
// 参考：https://vt100.net/docs/vt510-rm/HTS.html
const HorizontalTabSet = "\x1bH"

// TabClear（TBC）清除制表位。
//
// 默认值为 0。
//
// 可能的值：
// 0: 清除当前列的制表位。（默认）
// 3: 清除所有制表位。
//
//	CSI Pn g
//
// 参考：https://vt100.net/docs/vt510-rm/TBC.html
func TabClear(n int) string {
	var s string
	if n > 0 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "g"
}

// TBC 是 [TabClear] 的别名。
func TBC(n int) string {
	return TabClear(n)
}

// RequestPresentationStateReport（DECRQPSR）请求终端发送显示状态报告。这包括光标信息 [DECCIR] 和制表位 [DECTABSR] 报告。
//
// 默认值为 0。
//
// 可能的值：
// 0: 错误，请求被忽略。
// 1: 光标信息报告 [DECCIR]。
// 2: 制表位报告 [DECTABSR]。
//
//	CSI Ps $ w
//
// 参考：https://vt100.net/docs/vt510-rm/DECRQPSR.html
func RequestPresentationStateReport(n int) string {
	var s string
	if n > 0 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "$w"
}

// DECRQPSR 是 [RequestPresentationStateReport] 的别名。
func DECRQPSR(n int) string {
	return RequestPresentationStateReport(n)
}

// TabStopReport（DECTABSR）是对制表位报告请求的响应。
// 它报告终端中设置的制表位。
//
// 响应是由斜杠（/）字符分隔的制表位列表。
//
//	DCS 2 $ u D ... D ST
//
// 其中 D 是表示制表位的十进制数字。
//
// 参考：https://vt100.net/docs/vt510-rm/DECTABSR.html
func TabStopReport(stops ...int) string {
	var s []string //nolint:prealloc
	for _, v := range stops {
		s = append(s, strconv.Itoa(v))
	}
	return "\x1bP2$u" + strings.Join(s, "/") + "\x1b\\"
}

// DECTABSR 是 [TabStopReport] 的别名。
func DECTABSR(stops ...int) string {
	return TabStopReport(stops...)
}

// CursorInformationReport（DECCIR）是对光标信息报告请求的响应。它报告光标位置、视觉属性和字符保护属性。它还报告原点模式 [DECOM] 的状态和当前活动字符集。
//
// 响应是由分号（;）字符分隔的值列表。
//
//	DCS 1 $ u D ... D ST
//
// 其中 D 是表示值的十进制数字。
//
// 参考：https://vt100.net/docs/vt510-rm/DECCIR.html
func CursorInformationReport(values ...int) string {
	var s []string //nolint:prealloc
	for _, v := range values {
		s = append(s, strconv.Itoa(v))
	}
	return "\x1bP1$u" + strings.Join(s, ";") + "\x1b\\"
}

// DECCIR 是 [CursorInformationReport] 的别名。
func DECCIR(values ...int) string {
	return CursorInformationReport(values...)
}

// RepeatPreviousCharacter（REP）重复前一个字符 n 次。
// 这与键入相同字符 n 次相同。
//
// 默认值为 1。
//
//	CSI Pn b
//
// 参考：ECMA-48 § 8.3.103。
func RepeatPreviousCharacter(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "b"
}

// REP 是 [RepeatPreviousCharacter] 的别名。
func REP(n int) string {
	return RepeatPreviousCharacter(n)
}
