package ansi

import (
	"strconv"
)

// SaveCursor (DECSC) 是一个保存当前光标位置的转义序列。
//
//	ESC 7
//
// 请参阅：https://vt100.net/docs/vt510-rm/DECSC.html
const (
	SaveCursor = "\x1b7"
	DECSC      = SaveCursor
)

// RestoreCursor (DECRC) 是一个恢复光标位置的转义序列。
//
//	ESC 8
//
// 请参阅：https://vt100.net/docs/vt510-rm/DECRC.html
const (
	RestoreCursor = "\x1b8"
	DECRC         = RestoreCursor
)

// RequestCursorPosition 是一个请求当前光标位置的转义序列。
//
//	CSI 6 n
//
// 终端将以 CSI 序列的形式报告光标位置，格式如下：
//
//	CSI Pl ; Pc R
//
// 其中 Pl 是行号，Pc 是列号。
// 请参阅：https://vt100.net/docs/vt510-rm/CPR.html
//
// 已弃用：请使用 [RequestCursorPositionReport] 代替。
const RequestCursorPosition = "\x1b[6n"

// RequestExtendedCursorPosition (DECXCPR) 是一个请求光标位置报告的序列，包括当前页码。
//
//	CSI ? 6 n
//
// 终端将以 CSI 序列的形式报告光标位置，格式如下：
//
//	CSI ? Pl ; Pc ; Pp R
//
// 其中 Pl 是行号，Pc 是列号，Pp 是页码。
// 请参阅：https://vt100.net/docs/vt510-rm/DECXCPR.html
//
// 已弃用：请使用 [RequestExtendedCursorPositionReport] 代替。
const RequestExtendedCursorPosition = "\x1b[?6n"

// CursorUp (CUU) 返回一个使光标上移 n 个单元格的序列。
//
//	CSI n A
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUU.html
func CursorUp(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "A"
}

// CUU 是 [CursorUp] 的别名。
func CUU(n int) string {
	return CursorUp(n)
}

// CUU1 是一个使光标上移一个单元格的序列。
const CUU1 = "\x1b[A"

// CursorUp1 是一个使光标上移一个单元格的序列。
//
// 这相当于 CursorUp(1)。
//
// 已弃用：请使用 [CUU1] 代替。
const CursorUp1 = "\x1b[A"

// CursorDown (CUD) 返回一个使光标下移 n 个单元格的序列。
//
//	CSI n B
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUD.html
func CursorDown(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "B"
}

// CUD 是 [CursorDown] 的别名。
func CUD(n int) string {
	return CursorDown(n)
}

// CUD1 是一个使光标下移一个单元格的序列。
const CUD1 = "\x1b[B"

// CursorDown1 是一个使光标下移一个单元格的序列。
//
// 这相当于 CursorDown(1)。
//
// 已弃用：请使用 [CUD1] 代替。
const CursorDown1 = "\x1b[B"

// CursorForward (CUF) 返回一个使光标右移 n 个单元格的序列。
//
// # CSI n C
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUF.html
func CursorForward(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "C"
}

// CUF 是 [CursorForward] 的别名。
func CUF(n int) string {
	return CursorForward(n)
}

// CUF1 是一个使光标右移一个单元格的序列。
const CUF1 = "\x1b[C"

// CursorRight (CUF) 返回一个使光标右移 n 个单元格的序列。
//
//	CSI n C
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUF.html
//
// 已弃用：请使用 [CursorForward] 代替。
func CursorRight(n int) string {
	return CursorForward(n)
}

// CursorRight1 是一个使光标右移一个单元格的序列。
//
// 这相当于 CursorRight(1)。
//
// 已弃用：请使用 [CUF1] 代替。
const CursorRight1 = CUF1

// CursorBackward (CUB) 返回一个使光标左移 n 个单元格的序列。
//
// # CSI n D
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUB.html
func CursorBackward(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "D"
}

// CUB 是 [CursorBackward] 的别名。
func CUB(n int) string {
	return CursorBackward(n)
}

// CUB1 是一个使光标左移一个单元格的序列。
const CUB1 = "\x1b[D"

// CursorLeft (CUB) 返回一个使光标左移 n 个单元格的序列。
//
//	CSI n D
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUB.html
//
// 已弃用：请使用 [CursorBackward] 代替。
func CursorLeft(n int) string {
	return CursorBackward(n)
}

// CursorLeft1 是一个使光标左移一个单元格的序列。
//
// 这相当于 CursorLeft(1)。
//
// 已弃用：请使用 [CUB1] 代替。
const CursorLeft1 = CUB1

// CursorNextLine (CNL) 返回一个使光标移动到下一行开头 n 次的序列。
//
//	CSI n E
//
// 请参阅：https://vt100.net/docs/vt510-rm/CNL.html
func CursorNextLine(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "E"
}

// CNL 是 [CursorNextLine] 的别名。
func CNL(n int) string {
	return CursorNextLine(n)
}

// CursorPreviousLine (CPL) 返回一个使光标移动到上一行开头 n 次的序列。
//
//	CSI n F
//
// 请参阅：https://vt100.net/docs/vt510-rm/CPL.html
func CursorPreviousLine(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "F"
}

// CPL 是 [CursorPreviousLine] 的别名。
func CPL(n int) string {
	return CursorPreviousLine(n)
}

// CursorHorizontalAbsolute (CHA) 返回一个使光标移动到指定列的序列。
//
// 默认值为 1。
//
//	CSI n G
//
// 请参阅：https://vt100.net/docs/vt510-rm/CHA.html
func CursorHorizontalAbsolute(col int) string {
	var s string
	if col > 0 {
		s = strconv.Itoa(col)
	}
	return "\x1b[" + s + "G"
}

// CHA 是 [CursorHorizontalAbsolute] 的别名。
func CHA(col int) string {
	return CursorHorizontalAbsolute(col)
}

// CursorPosition (CUP) 返回一个将光标设置到指定行和列的序列。
//
// 默认值为 1,1。
//
//	CSI n ; m H
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUP.html
func CursorPosition(col, row int) string {
	if row <= 1 && col <= 1 {
		return CursorHomePosition
	}

	var r, c string
	if row > 0 {
		r = strconv.Itoa(row)
	}
	if col > 0 {
		c = strconv.Itoa(col)
	}
	return "\x1b[" + r + ";" + c + "H"
}

// CUP 是 [CursorPosition] 的别名。
func CUP(col, row int) string {
	return CursorPosition(col, row)
}

// CursorHomePosition 是一个使光标移动到滚动区域左上角的序列。
//
// 这相当于 [CursorPosition](1, 1)。
const CursorHomePosition = "\x1b[H"

// SetCursorPosition (CUP) 返回一个将光标设置到指定行和列的序列。
//
//	CSI n ; m H
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUP.html
//
// 已弃用：请使用 [CursorPosition] 代替。
func SetCursorPosition(col, row int) string {
	if row <= 0 && col <= 0 {
		return HomeCursorPosition
	}

	var r, c string
	if row > 0 {
		r = strconv.Itoa(row)
	}
	if col > 0 {
		c = strconv.Itoa(col)
	}
	return "\x1b[" + r + ";" + c + "H"
}

// HomeCursorPosition 是一个使光标移动到滚动区域左上角的序列。这相当于 `SetCursorPosition(1, 1)`。
//
// 已弃用：请使用 [CursorHomePosition] 代替。
const HomeCursorPosition = CursorHomePosition

// MoveCursor (CUP) 返回一个将光标设置到指定行和列的序列。
//
//	CSI n ; m H
//
// 请参阅：https://vt100.net/docs/vt510-rm/CUP.html
//
// 已弃用：请使用 [CursorPosition] 代替。
func MoveCursor(col, row int) string {
	return SetCursorPosition(col, row)
}

// CursorOrigin 是一个使光标移动到显示区域左上角的序列。这相当于 `SetCursorPosition(1, 1)`。
//
// 已弃用：请使用 [CursorHomePosition] 代替。
const CursorOrigin = "\x1b[1;1H"

// MoveCursorOrigin 是一个使光标移动到显示区域左上角的序列。这相当于 `SetCursorPosition(1, 1)`。
//
// 已弃用：请使用 [CursorHomePosition] 代替。
const MoveCursorOrigin = CursorOrigin

// CursorHorizontalForwardTab (CHT) 返回一个使光标移动到下一个制表位 n 次的序列。
//
// 默认值为 1。
//
//	CSI n I
//
// 请参阅：https://vt100.net/docs/vt510-rm/CHT.html
func CursorHorizontalForwardTab(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "I"
}

// CHT 是 [CursorHorizontalForwardTab] 的别名。
func CHT(n int) string {
	return CursorHorizontalForwardTab(n)
}

// EraseCharacter (ECH) 返回一个从屏幕上擦除 n 个字符的序列。这不会影响其他单元格属性。
//
// 默认值为 1。
//
//	CSI n X
//
// 请参阅：https://vt100.net/docs/vt510-rm/ECH.html
func EraseCharacter(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "X"
}

// ECH 是 [EraseCharacter] 的别名。
func ECH(n int) string {
	return EraseCharacter(n)
}

// CursorBackwardTab (CBT) 返回一个使光标移动到上一个制表位 n 次的序列。
//
// 默认值为 1。
//
//	CSI n Z
//
// 请参阅：https://vt100.net/docs/vt510-rm/CBT.html
func CursorBackwardTab(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "Z"
}

// CBT 是 [CursorBackwardTab] 的别名。
func CBT(n int) string {
	return CursorBackwardTab(n)
}

// VerticalPositionAbsolute (VPA) 返回一个使光标移动到指定行的序列。
//
// 默认值为 1。
//
//	CSI n d
//
// 请参阅：https://vt100.net/docs/vt510-rm/VPA.html
func VerticalPositionAbsolute(row int) string {
	var s string
	if row > 0 {
		s = strconv.Itoa(row)
	}
	return "\x1b[" + s + "d"
}

// VPA 是 [VerticalPositionAbsolute] 的别名。
func VPA(row int) string {
	return VerticalPositionAbsolute(row)
}

// VerticalPositionRelative (VPR) 返回一个使光标相对于当前位置下移 n 行的序列。
//
// 默认值为 1。
//
//	CSI n e
//
// 请参阅：https://vt100.net/docs/vt510-rm/VPR.html
func VerticalPositionRelative(n int) string {
	var s string
	if n > 1 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "e"
}

// VPR 是 [VerticalPositionRelative] 的别名。
func VPR(n int) string {
	return VerticalPositionRelative(n)
}

// HorizontalVerticalPosition (HVP) 返回一个使光标移动到指定行和列的序列。
//
// 默认值为 1,1。
//
//	CSI n ; m f
//
// 这与 [CursorPosition] 效果相同。
//
// 请参阅：https://vt100.net/docs/vt510-rm/HVP.html
func HorizontalVerticalPosition(col, row int) string {
	var r, c string
	if row > 0 {
		r = strconv.Itoa(row)
	}
	if col > 0 {
		c = strconv.Itoa(col)
	}
	return "\x1b[" + r + ";" + c + "f"
}

// HVP 是 [HorizontalVerticalPosition] 的别名。
func HVP(col, row int) string {
	return HorizontalVerticalPosition(col, row)
}

// HorizontalVerticalHomePosition 是一个使光标移动到滚动区域左上角的序列。这相当于
// `HorizontalVerticalPosition(1, 1)`。
const HorizontalVerticalHomePosition = "\x1b[f"

// SaveCurrentCursorPosition (SCOSC) 是一个用于 SCO 控制台模式的保存当前光标位置的序列。
//
//	CSI s
//
// 这与 [DECSC] 类似，只是不保存光标所在的页码。
//
// 请参阅：https://vt100.net/docs/vt510-rm/SCOSC.html
const (
	SaveCurrentCursorPosition = "\x1b[s"
	SCOSC                     = SaveCurrentCursorPosition
)

// SaveCursorPosition (SCP 或 SCOSC) 是一个保存光标位置的序列。
//
//	CSI s
//
// 这与 Save 类似，只是不保存光标所在的页码。
//
// 请参阅：https://vt100.net/docs/vt510-rm/SCOSC.html
//
// 已弃用：请使用 [SaveCurrentCursorPosition] 代替。
const SaveCursorPosition = "\x1b[s"

// RestoreCurrentCursorPosition (SCORC) 是一个用于 SCO 控制台模式的恢复当前光标位置的序列。
//
//	CSI u
//
// 这与 [DECRC] 类似，只是不恢复光标保存时的页码。
//
// 请参阅：https://vt100.net/docs/vt510-rm/SCORC.html
const (
	RestoreCurrentCursorPosition = "\x1b[u"
	SCORC                        = RestoreCurrentCursorPosition
)

// RestoreCursorPosition (RCP 或 SCORC) 是一个恢复光标位置的序列。
//
//	CSI u
//
// 这与 Restore 类似，只是光标会停留在保存时的同一页。
//
// 请参阅：https://vt100.net/docs/vt510-rm/SCORC.html
//
// 已弃用：请使用 [RestoreCurrentCursorPosition] 代替。
const RestoreCursorPosition = "\x1b[u"

// SetCursorStyle (DECSCUSR) 返回一个更改光标样式的序列。
//
// 默认值为 1。
//
//	CSI Ps SP q
//
// 其中 Ps 是光标样式：
//
//	0: 闪烁块
//	1: 闪烁块（默认）
//	2: 稳定块
//	3: 闪烁下划线
//	4: 稳定下划线
//	5: 闪烁竖线（xterm）
//	6: 稳定竖线（xterm）
//
// 请参阅：https://vt100.net/docs/vt510-rm/DECSCUSR.html
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h4-Functions-using-CSI-_-ordered-by-the-final-character-lparen-s-rparen:CSI-Ps-SP-q.1D81
func SetCursorStyle(style int) string {
	if style < 0 {
		style = 0
	}
	return "\x1b[" + strconv.Itoa(style) + " q"
}

// DECSCUSR 是 [SetCursorStyle] 的别名。
func DECSCUSR(style int) string {
	return SetCursorStyle(style)
}

// SetPointerShape 返回一个更改鼠标指针光标形状的序列。使用 "default" 表示默认指针形状。
//
//	OSC 22 ; Pt ST
//	OSC 22 ; Pt BEL
//
// 其中 Pt 是指针形状名称。该名称可以是操作系统能够理解的任何内容。一些常见名称包括：
//
//   - copy
//   - crosshair
//   - default
//   - ew-resize
//   - n-resize
//   - text
//   - wait
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Operating-System-Commands
func SetPointerShape(shape string) string {
	return "\x1b]22;" + shape + "\x07"
}

// ReverseIndex (RI) 是一个使光标在同一列中上移一行的转义序列。如果光标在顶部边缘，屏幕会向下滚动。
//
// 这与 [RI] 效果相同。
const ReverseIndex = "\x1bM"

// HorizontalPositionAbsolute (HPA) 返回一个使光标移动到指定列的序列。这与 [CUP] 效果相同。
//
// 默认值为 1。
//
//	CSI n \`
//
// 请参阅：https://vt100.net/docs/vt510-rm/HPA.html
func HorizontalPositionAbsolute(col int) string {
	var s string
	if col > 0 {
		s = strconv.Itoa(col)
	}
	return "\x1b[" + s + "`"
}

// HPA 是 [HorizontalPositionAbsolute] 的别名。
func HPA(col int) string {
	return HorizontalPositionAbsolute(col)
}

// HorizontalPositionRelative (HPR) 返回一个使光标相对于当前位置右移 n 列的序列。这与 [CUP] 效果相同。
//
// 默认值为 1。
//
//	CSI n a
//
// 请参阅：https://vt100.net/docs/vt510-rm/HPR.html
func HorizontalPositionRelative(n int) string {
	var s string
	if n > 0 {
		s = strconv.Itoa(n)
	}
	return "\x1b[" + s + "a"
}

// HPR 是 [HorizontalPositionRelative] 的别名。
func HPR(n int) string {
	return HorizontalPositionRelative(n)
}

// Index (IND) 是一个使光标在同一列中下移一行的转义序列。如果光标在底部边缘，屏幕会向上滚动。
// 这与 [IND] 效果相同。
const Index = "\x1bD"
