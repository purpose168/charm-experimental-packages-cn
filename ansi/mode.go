package ansi

import (
	"strconv"
	"strings"
)

// ModeSetting 表示模式设置。
type ModeSetting byte

// ModeSetting 常量。
const (
	ModeNotRecognized ModeSetting = iota
	ModeSet
	ModeReset
	ModePermanentlySet
	ModePermanentlyReset
)

// IsNotRecognized 如果模式未被识别，则返回 true。
func (m ModeSetting) IsNotRecognized() bool {
	return m == ModeNotRecognized
}

// IsSet 如果模式已设置或永久设置，则返回 true。
func (m ModeSetting) IsSet() bool {
	return m == ModeSet || m == ModePermanentlySet
}

// IsReset 如果模式已重置或永久重置，则返回 true。
func (m ModeSetting) IsReset() bool {
	return m == ModeReset || m == ModePermanentlyReset
}

// IsPermanentlySet 如果模式被永久设置，则返回 true。
func (m ModeSetting) IsPermanentlySet() bool {
	return m == ModePermanentlySet
}

// IsPermanentlyReset 如果模式被永久重置，则返回 true。
func (m ModeSetting) IsPermanentlyReset() bool {
	return m == ModePermanentlyReset
}

// Mode 表示终端模式的接口。
// 模式可以被设置、重置和请求。
type Mode interface {
	Mode() int
}

// SetMode (SM 或 DECSET) 返回设置模式的序列。
// 模式参数是要设置的模式列表。
//
// 如果其中一个模式是 [DECMode]，该函数将返回两个转义序列。
//
// ANSI 格式：
//
//	CSI Pd ; ... ; Pd h
//
// DEC 格式：
//
//	CSI ? Pd ; ... ; Pd h
//
// 参见：https://vt100.net/docs/vt510-rm/SM.html
func SetMode(modes ...Mode) string {
	return setMode(false, modes...)
}

// SM 是 [SetMode] 的别名。
func SM(modes ...Mode) string {
	return SetMode(modes...)
}

// DECSET 是 [SetMode] 的别名。
func DECSET(modes ...Mode) string {
	return SetMode(modes...)
}

// ResetMode (RM 或 DECRST) 返回重置模式的序列。
// 模式参数是要重置的模式列表。
//
// 如果其中一个模式是 [DECMode]，该函数将返回两个转义序列。
//
// ANSI 格式：
//
//	CSI Pd ; ... ; Pd l
//
// DEC 格式：
//
//	CSI ? Pd ; ... ; Pd l
//
// 参见：https://vt100.net/docs/vt510-rm/RM.html
func ResetMode(modes ...Mode) string {
	return setMode(true, modes...)
}

// RM 是 [ResetMode] 的别名。
func RM(modes ...Mode) string {
	return ResetMode(modes...)
}

// DECRST 是 [ResetMode] 的别名。
func DECRST(modes ...Mode) string {
	return ResetMode(modes...)
}

func setMode(reset bool, modes ...Mode) (s string) {
	if len(modes) == 0 {
		return s
	}

	cmd := "h"
	if reset {
		cmd = "l"
	}

	seq := "\x1b["
	if len(modes) == 1 {
		switch modes[0].(type) {
		case DECMode:
			seq += "?"
		}
		return seq + strconv.Itoa(modes[0].Mode()) + cmd
	}

	dec := make([]string, 0, len(modes)/2)
	ansi := make([]string, 0, len(modes)/2)
	for _, m := range modes {
		switch m.(type) {
		case DECMode:
			dec = append(dec, strconv.Itoa(m.Mode()))
		case ANSIMode:
			ansi = append(ansi, strconv.Itoa(m.Mode()))
		}
	}

	if len(ansi) > 0 {
		s += seq + strings.Join(ansi, ";") + cmd
	}
	if len(dec) > 0 {
		s += seq + "?" + strings.Join(dec, ";") + cmd
	}
	return s
}

// RequestMode (DECRQM) 返回向终端请求模式的序列。
// 终端用报告模式函数 [DECRPM] 响应。
//
// ANSI 格式：
//
//	CSI Pa $ p
//
// DEC 格式：
//
//	CSI ? Pa $ p
//
// 参见：https://vt100.net/docs/vt510-rm/DECRQM.html
func RequestMode(m Mode) string {
	seq := "\x1b["
	switch m.(type) {
	case DECMode:
		seq += "?"
	}
	return seq + strconv.Itoa(m.Mode()) + "$p"
}

// DECRQM 是 [RequestMode] 的别名。
func DECRQM(m Mode) string {
	return RequestMode(m)
}

// ReportMode (DECRPM) 返回终端在响应模式请求 [DECRQM] 时发送给主机的序列。
//
// ANSI 格式：
//
//	CSI Pa ; Ps ; $ y
//
// DEC 格式：
//
//	CSI ? Pa ; Ps $ y
//
// 其中 Pa 是模式编号，Ps 是模式值。
//
//	0: 未识别
//	1: 已设置
//	2: 已重置
//	3: 永久设置
//	4: 永久重置
//
// 参见：https://vt100.net/docs/vt510-rm/DECRPM.html
func ReportMode(mode Mode, value ModeSetting) string {
	if value > 4 {
		value = 0
	}
	switch mode.(type) {
	case DECMode:
		return "\x1b[?" + strconv.Itoa(mode.Mode()) + ";" + strconv.Itoa(int(value)) + "$y"
	}
	return "\x1b[" + strconv.Itoa(mode.Mode()) + ";" + strconv.Itoa(int(value)) + "$y"
}

// DECRPM 是 [ReportMode] 的别名。
func DECRPM(mode Mode, value ModeSetting) string {
	return ReportMode(mode, value)
}

// ANSIMode 表示 ANSI 终端模式。
type ANSIMode int //nolint:revive

// Mode 返回 ANSI 模式的整数值。
func (m ANSIMode) Mode() int {
	return int(m)
}

// DECMode 表示专用 DEC 终端模式。
type DECMode int

// Mode 返回 DEC 模式的整数值。
func (m DECMode) Mode() int {
	return int(m)
}

// Keyboard Action Mode (KAM) 是控制键盘锁定的模式。
// 当键盘被锁定时，它不能向终端发送数据。
//
// 参见：https://vt100.net/docs/vt510-rm/KAM.html
const (
	ModeKeyboardAction = ANSIMode(2)
	KAM                = ModeKeyboardAction

	SetModeKeyboardAction     = "\x1b[2h"
	ResetModeKeyboardAction   = "\x1b[2l"
	RequestModeKeyboardAction = "\x1b[2$p"
)

// Insert/Replace Mode (IRM) 是确定字符在输入时是插入还是替换的模式。
//
// 启用时，字符在光标位置插入，将右侧字符向右推。禁用时，字符替换光标位置的字符。
//
// 参见：https://vt100.net/docs/vt510-rm/IRM.html
const (
	ModeInsertReplace = ANSIMode(4)
	IRM               = ModeInsertReplace

	SetModeInsertReplace     = "\x1b[4h"
	ResetModeInsertReplace   = "\x1b[4l"
	RequestModeInsertReplace = "\x1b[4$p"
)

// BiDirectional Support Mode (BDSM) 是确定终端是否支持双向文本的模式。启用时，终端支持双向文本并设置为隐式双向模式。禁用时，终端不支持双向文本。
//
// 参见 ECMA-48 7.2.1。
const (
	ModeBiDirectionalSupport = ANSIMode(8)
	BDSM                     = ModeBiDirectionalSupport

	SetModeBiDirectionalSupport     = "\x1b[8h"
	ResetModeBiDirectionalSupport   = "\x1b[8l"
	RequestModeBiDirectionalSupport = "\x1b[8$p"
)

// Send Receive Mode (SRM) 或 Local Echo Mode 是确定终端是否将字符回显给主机的模式。启用时，终端在输入字符时将其发送给主机。
//
// 参见：https://vt100.net/docs/vt510-rm/SRM.html
const (
	ModeSendReceive = ANSIMode(12)
	ModeLocalEcho   = ModeSendReceive
	SRM             = ModeSendReceive

	SetModeSendReceive     = "\x1b[12h"
	ResetModeSendReceive   = "\x1b[12l"
	RequestModeSendReceive = "\x1b[12$p"

	SetModeLocalEcho     = "\x1b[12h"
	ResetModeLocalEcho   = "\x1b[12l"
	RequestModeLocalEcho = "\x1b[12$p"
)

// Line Feed/New Line Mode (LNM) 是确定终端是否将换行符解释为新行的模式。
//
// 启用时，终端将换行符解释为新行。禁用时，终端将换行符解释为换行。
//
// 新行将光标移动到下一行的第一个位置。
// 换行将光标向下移动一行而不更改列，必要时滚动屏幕。
//
// 参见：https://vt100.net/docs/vt510-rm/LNM.html
const (
	ModeLineFeedNewLine = ANSIMode(20)
	LNM                 = ModeLineFeedNewLine

	SetModeLineFeedNewLine     = "\x1b[20h"
	ResetModeLineFeedNewLine   = "\x1b[20l"
	RequestModeLineFeedNewLine = "\x1b[20$p"
)

// Cursor Keys Mode (DECCKM) 是确定光标键发送 ANSI 光标序列还是应用程序序列的模式。
//
// 参见：https://vt100.net/docs/vt510-rm/DECCKM.html
const (
	ModeCursorKeys = DECMode(1)
	DECCKM         = ModeCursorKeys

	SetModeCursorKeys     = "\x1b[?1h"
	ResetModeCursorKeys   = "\x1b[?1l"
	RequestModeCursorKeys = "\x1b[?1$p"
)

// Origin Mode (DECOM) 是确定光标移动到起始位置还是边距位置的模式。
//
// 参见：https://vt100.net/docs/vt510-rm/DECOM.html
const (
	ModeOrigin = DECMode(6)
	DECOM      = ModeOrigin

	SetModeOrigin     = "\x1b[?6h"
	ResetModeOrigin   = "\x1b[?6l"
	RequestModeOrigin = "\x1b[?6$p"
)

// Auto Wrap Mode (DECAWM) 是确定光标到达右边界时是否换行到下一行的模式。
//
// 参见：https://vt100.net/docs/vt510-rm/DECAWM.html
const (
	ModeAutoWrap = DECMode(7)
	DECAWM       = ModeAutoWrap

	SetModeAutoWrap     = "\x1b[?7h"
	ResetModeAutoWrap   = "\x1b[?7l"
	RequestModeAutoWrap = "\x1b[?7$p"
)

// X10 Mouse Mode 是确定鼠标是否报告按钮按下的模式。
//
// 终端用以下编码响应：
//
//	CSI M CbCxCy
//
// 其中 Cb 是按钮-1，可以是 1、2 或 3。
// Cx 和 Cy 是鼠标事件的 x 和 y 坐标。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseX10 = DECMode(9)

	SetModeMouseX10     = "\x1b[?9h"
	ResetModeMouseX10   = "\x1b[?9l"
	RequestModeMouseX10 = "\x1b[?9$p"
)

// Text Cursor Enable Mode (DECTCEM) 是显示/隐藏光标的模式。
//
// 参见：https://vt100.net/docs/vt510-rm/DECTCEM.html
const (
	ModeTextCursorEnable = DECMode(25)
	DECTCEM              = ModeTextCursorEnable

	SetModeTextCursorEnable     = "\x1b[?25h"
	ResetModeTextCursorEnable   = "\x1b[?25l"
	RequestModeTextCursorEnable = "\x1b[?25$p"
)

// 这些是 [SetModeTextCursorEnable] 和 [ResetModeTextCursorEnable] 的别名。
const (
	ShowCursor = SetModeTextCursorEnable
	HideCursor = ResetModeTextCursorEnable
)

// Numeric Keypad Mode (DECNKM) 是确定小键盘发送应用程序序列还是数字序列的模式。
//
// 这与 [DECKPAM] 和 [DECKPNM] 类似，但使用不同的序列。
//
// 参见：https://vt100.net/docs/vt510-rm/DECNKM.html
const (
	ModeNumericKeypad = DECMode(66)
	DECNKM            = ModeNumericKeypad

	SetModeNumericKeypad     = "\x1b[?66h"
	ResetModeNumericKeypad   = "\x1b[?66l"
	RequestModeNumericKeypad = "\x1b[?66$p"
)

// Backarrow Key Mode (DECBKM) 是确定退格键发送退格字符还是删除字符的模式。默认禁用。
//
// 参见：https://vt100.net/docs/vt510-rm/DECBKM.html
const (
	ModeBackarrowKey = DECMode(67)
	DECBKM           = ModeBackarrowKey

	SetModeBackarrowKey     = "\x1b[?67h"
	ResetModeBackarrowKey   = "\x1b[?67l"
	RequestModeBackarrowKey = "\x1b[?67$p"
)

// Left Right Margin Mode (DECLRMM) 是确定是否可以使用 [DECSLRM] 设置左右边距的模式。
//
// 参见：https://vt100.net/docs/vt510-rm/DECLRMM.html
const (
	ModeLeftRightMargin = DECMode(69)
	DECLRMM             = ModeLeftRightMargin

	SetModeLeftRightMargin     = "\x1b[?69h"
	ResetModeLeftRightMargin   = "\x1b[?69l"
	RequestModeLeftRightMargin = "\x1b[?69$p"
)

// Normal Mouse Mode 是确定鼠标是否报告按钮按下、释放的模式。它还会报告修饰键、滚轮事件和额外按钮。
//
// 它使用与 [ModeMouseX10] 相同的编码，但有一些差异：
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseNormal = DECMode(1000)

	SetModeMouseNormal     = "\x1b[?1000h"
	ResetModeMouseNormal   = "\x1b[?1000l"
	RequestModeMouseNormal = "\x1b[?1000$p"
)

// Highlight Mouse Tracking 是确定鼠标是否报告按钮按下、释放和高亮单元格的模式。
//
// 它使用与 [ModeMouseNormal] 相同的编码，但有一些差异：
//
// 在高亮事件中，终端用以下编码响应：
//
//	CSI t CxCy
//	CSI T CxCyCxCyCxCy
//
// 其中参数是 startx、starty、endx、endy、mousex 和 mousey。
const (
	ModeMouseHighlight = DECMode(1001)

	SetModeMouseHighlight     = "\x1b[?1001h"
	ResetModeMouseHighlight   = "\x1b[?1001l"
	RequestModeMouseHighlight = "\x1b[?1001$p"
)

// VT Hilite Mouse Tracking 是确定鼠标是否报告按钮按下、释放和高亮单元格的模式。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
//

// Button Event Mouse Tracking 本质上与 [ModeMouseNormal] 相同，
// 但它还会在按钮按下时报告按钮移动事件。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseButtonEvent = DECMode(1002)

	SetModeMouseButtonEvent     = "\x1b[?1002h"
	ResetModeMouseButtonEvent   = "\x1b[?1002l"
	RequestModeMouseButtonEvent = "\x1b[?1002$p"
)

// Any Event Mouse Tracking 与 [ModeMouseButtonEvent] 相同，除了
// 即使没有按下鼠标按钮，也会报告所有移动事件。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseAnyEvent = DECMode(1003)

	SetModeMouseAnyEvent     = "\x1b[?1003h"
	ResetModeMouseAnyEvent   = "\x1b[?1003l"
	RequestModeMouseAnyEvent = "\x1b[?1003$p"
)

// Focus Event Mode 是确定终端是否报告焦点和失焦事件的模式。
//
// 终端发送以下编码：
//
//	CSI I // Focus In
//	CSI O // Focus Out
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Focus-Tracking
const (
	ModeFocusEvent = DECMode(1004)

	SetModeFocusEvent     = "\x1b[?1004h"
	ResetModeFocusEvent   = "\x1b[?1004l"
	RequestModeFocusEvent = "\x1b[?1004$p"
)

// SGR Extended Mouse Mode 是更改鼠标跟踪编码以使用 SGR 参数的模式。
//
// 终端用以下编码响应：
//
//	CSI < Cb ; Cx ; Cy M
//
// 其中 Cb 与 [ModeMouseNormal] 相同，Cx 和 Cy 是 x 和 y 坐标。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseExtSgr = DECMode(1006)

	SetModeMouseExtSgr     = "\x1b[?1006h"
	ResetModeMouseExtSgr   = "\x1b[?1006l"
	RequestModeMouseExtSgr = "\x1b[?1006$p"
)

// UTF-8 Extended Mouse Mode 是更改鼠标跟踪编码以使用 UTF-8 参数的模式。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseExtUtf8 = DECMode(1005)

	SetModeMouseExtUtf8     = "\x1b[?1005h"
	ResetModeMouseExtUtf8   = "\x1b[?1005l"
	RequestModeMouseExtUtf8 = "\x1b[?1005$p"
)

// URXVT Extended Mouse Mode 是更改鼠标跟踪编码以使用替代编码的模式。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseExtUrxvt = DECMode(1015)

	SetModeMouseExtUrxvt     = "\x1b[?1015h"
	ResetModeMouseExtUrxvt   = "\x1b[?1015l"
	RequestModeMouseExtUrxvt = "\x1b[?1015$p"
)

// SGR Pixel Extended Mouse Mode 是更改鼠标跟踪编码以使用带像素坐标的 SGR 参数的模式。
//
// 这与 [ModeMouseExtSgr] 类似，但还报告像素坐标。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
const (
	ModeMouseExtSgrPixel = DECMode(1016)

	SetModeMouseExtSgrPixel     = "\x1b[?1016h"
	ResetModeMouseExtSgrPixel   = "\x1b[?1016l"
	RequestModeMouseExtSgrPixel = "\x1b[?1016$p"
)

// Alternate Screen Mode 是确定备用屏幕缓冲区是否激活的模式。当启用此模式时，备用屏幕缓冲区被清除。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-The-Alternate-Screen-Buffer
const (
	ModeAltScreen = DECMode(1047)

	SetModeAltScreen     = "\x1b[?1047h"
	ResetModeAltScreen   = "\x1b[?1047l"
	RequestModeAltScreen = "\x1b[?1047$p"
)

// Save Cursor Mode 是保存光标位置的模式。
// 这相当于 [SaveCursor] 和 [RestoreCursor]。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-The-Alternate-Screen-Buffer
const (
	ModeSaveCursor = DECMode(1048)

	SetModeSaveCursor     = "\x1b[?1048h"
	ResetModeSaveCursor   = "\x1b[?1048l"
	RequestModeSaveCursor = "\x1b[?1048$p"
)

// Alternate Screen Save Cursor Mode 是保存光标位置（如 [ModeSaveCursor]）、切换到备用屏幕缓冲区（如 [ModeAltScreen]）
// 并在切换时清除屏幕的模式。
//
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-The-Alternate-Screen-Buffer
const (
	ModeAltScreenSaveCursor = DECMode(1049)

	SetModeAltScreenSaveCursor     = "\x1b[?1049h"
	ResetModeAltScreenSaveCursor   = "\x1b[?1049l"
	RequestModeAltScreenSaveCursor = "\x1b[?1049$p"
)

// Bracketed Paste Mode 是确定粘贴的文本是否用转义序列括起来的模式。
//
// 参见：https://cirw.in/blog/bracketed-paste
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Bracketed-Paste-Mode
const (
	ModeBracketedPaste = DECMode(2004)

	SetModeBracketedPaste     = "\x1b[?2004h"
	ResetModeBracketedPaste   = "\x1b[?2004l"
	RequestModeBracketedPaste = "\x1b[?2004$p"
)

// Synchronized Output Mode 是确定输出是否与终端同步的模式。
//
// 参见：https://gist.github.com/christianparpart/d8a62cc1ab659194337d73e399004036
const (
	ModeSynchronizedOutput = DECMode(2026)

	SetModeSynchronizedOutput     = "\x1b[?2026h"
	ResetModeSynchronizedOutput   = "\x1b[?2026l"
	RequestModeSynchronizedOutput = "\x1b[?2026$p"
)

// Unicode Core Mode 是确定终端是否应使用 Unicode 字形聚类来计算每个终端单元格的字形宽度的模式。
//
// 参见：https://github.com/contour-terminal/terminal-unicode-core
const (
	ModeUnicodeCore = DECMode(2027)

	SetModeUnicodeCore     = "\x1b[?2027h"
	ResetModeUnicodeCore   = "\x1b[?2027l"
	RequestModeUnicodeCore = "\x1b[?2027$p"
)

//

// ModeLightDark 是启用报告操作系统配色方案（亮色或暗色）偏好的模式。它将配色方案报告为 [DSR] 和 [LightDarkReport] 转义序列，编码如下：
//
//	CSI ? 997 ; 1 n   用于暗色模式
//	CSI ? 997 ; 2 n   用于亮色模式
//
// 还可以通过以下 [DSR] 和 [RequestLightDarkReport] 转义序列请求配色偏好：
//
//	CSI ? 996 n
//
// 参见：https://contour-terminal.org/vt-extensions/color-palette-update-notifications/
const (
	ModeLightDark = DECMode(2031)

	SetModeLightDark     = "\x1b[?2031h"
	ResetModeLightDark   = "\x1b[?2031l"
	RequestModeLightDark = "\x1b[?2031$p"
)

// ModeInBandResize 是将终端调整大小事件报告为转义序列的模式。这对不支持 [SIGWINCH] 的系统（如 Windows）很有用。
//
// 终端然后发送以下编码：
//
//	CSI 48 ; cellsHeight ; cellsWidth ; pixelHeight ; pixelWidth t
//
// 参见：https://gist.github.com/rockorager/e695fb2924d36b2bcf1fff4a3704bd83
const (
	ModeInBandResize = DECMode(2048)

	SetModeInBandResize     = "\x1b[?2048h"
	ResetModeInBandResize   = "\x1b[?2048l"
	RequestModeInBandResize = "\x1b[?2048$p"
)

// Win32Input 是确定输入是否由 Win32 控制台和 Conpty 处理的模式。
//
// 参见：https://github.com/microsoft/terminal/blob/main/doc/specs/%234999%20-%20Improved%20keyboard%20handling%20in%20Conpty.md
const (
	ModeWin32Input = DECMode(9001)

	SetModeWin32Input     = "\x1b[?9001h"
	ResetModeWin32Input   = "\x1b[?9001l"
	RequestModeWin32Input = "\x1b[?9001$p"
)
