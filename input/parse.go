package input

import (
	"bytes"
	"encoding/base64"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
	"github.com/rivo/uniseg"
)

// 控制解析器行为的标志。
const (
	// 当设置此标志时，驱动程序会将 Ctrl+Space 和 Ctrl+@ 视为相同的按键序列。
	//
	// 历史上，ANSI 规范在 Ctrl+Space 和 Ctrl+@ 按键序列上都会生成 NUL (0x00)。
	// 此标志允许驱动程序将两者视为相同的按键序列。
	FlagCtrlAt = 1 << iota

	// 当设置此标志时，驱动程序会将 Tab 键和 Ctrl+I 视为相同的按键序列。
	//
	// 历史上，ANSI 规范在 Tab 键和 Ctrl+I 上都会生成 HT (0x09)。
	// 此标志允许驱动程序将两者视为相同的按键序列。
	FlagCtrlI

	// 当设置此标志时，驱动程序会将 Enter 键和 Ctrl+M 视为相同的按键序列。
	//
	// 历史上，ANSI 规范在 Enter 键和 Ctrl+M 上都会生成 CR (0x0D)。
	// 此标志允许驱动程序将两者视为相同的按键。
	FlagCtrlM

	// 当设置此标志时，驱动程序会将 Escape 和 Ctrl+[ 视为相同的按键序列。
	//
	// 历史上，ANSI 规范在 Escape 键和 Ctrl+[ 上都会生成 ESC (0x1B)。
	// 此标志允许驱动程序将两者视为相同的按键序列。
	FlagCtrlOpenBracket

	// 当设置此标志时，驱动程序会在按下 Backspace 键时发送 BS (0x08 字节) 字符
	// 而不是 DEL (0x7F 字节) 字符。
	//
	// VT100 终端同时有 Backspace 和 Delete 键。VT220 终端去掉了 Backspace 键，
	// 用 Delete 键取而代之。两个终端在按下 Delete 键时都会发送 DEL 字符。
	// 现代终端和 PC 后来重新添加了 Delete 键，但使用了不同的按键序列，
	// 而 Backspace 键被标准化为发送 DEL 字符。
	FlagBackspace

	// 当设置此标志时，驱动程序会识别 Find 键，而不是将其视为 Home 键。
	//
	// Find 键是 VT220 键盘的一部分，在现代 PC 中不再使用。
	FlagFind

	// 当设置此标志时，驱动程序会识别 Select 键，而不是将其视为 End 键。
	//
	// Symbol 键是 VT220 键盘的一部分，在现代 PC 中不再使用。
	FlagSelect

	// 当设置此标志时，驱动程序会使用 Terminfo 数据库覆盖默认按键序列。
	FlagTerminfo

	// 当设置此标志时，驱动程序会将功能键 (F13-F63) 保留为符号。
	//
	// 由于这些键不是当今标准 20 世纪键盘的一部分，
	// 我们将它们视为 F1-F12 修饰键，即 ctrl/shift/alt + Fn 组合键。
	// 键定义来自 Terminfo，此标志仅在未设置 FlagTerminfo 时有用。
	FlagFKeys

	// 当设置此标志时，驱动程序会在 Windows 上启用鼠标模式。
	// 这仅在 Windows 上有用，在其他平台上没有效果。
	FlagMouseMode
)

// Parser 是输入转义序列的解析器。
type Parser struct {
	flags int
}

// NewParser 返回一个新的输入解析器。这是一个低级解析器，将转义序列解析为人类可读的事件。
// 这与 [ansi.Parser] 和 [ansi.DecodeSequence] 不同，因为它能识别一些终端可能发送的不正确序列。
//
// 例如，X10 鼠标协议发送一个 `CSI M` 序列，后跟 3 个字节。如果解析器不识别这 3 个字节，
// 它们可能会被回显到终端输出，造成混乱。
//
// 另一个例子是 URxvt 如何使用无效的 CSI 最终字符（如 '$'）为修饰键发送无效序列。
//
// 使用标志来控制模糊按键序列的行为。
func NewParser(flags int) *Parser {
	return &Parser{flags: flags}
}

// parseSequence 查找第一个被识别的事件序列并返回它及其长度。
//
// 如果没有识别到序列或缓冲区为空，它将返回零和 nil。如果序列不受支持，将返回 UnknownEvent。
func (p *Parser) parseSequence(buf []byte) (n int, Event Event) {
	if len(buf) == 0 {
		return 0, nil
	}

	switch b := buf[0]; b {
	case ansi.ESC:
		if len(buf) == 1 {
			// Escape key
			return 1, KeyPressEvent{Code: KeyEscape}
		}

		switch bPrime := buf[1]; bPrime {
		case 'O': // Esc-prefixed SS3
			return p.parseSs3(buf)
		case 'P': // Esc-prefixed DCS
			return p.parseDcs(buf)
		case '[': // Esc-prefixed CSI
			return p.parseCsi(buf)
		case ']': // Esc-prefixed OSC
			return p.parseOsc(buf)
		case '_': // Esc-prefixed APC
			return p.parseApc(buf)
		case '^': // Esc-prefixed PM
			return p.parseStTerminated(ansi.PM, '^', nil)(buf)
		case 'X': // Esc-prefixed SOS
			return p.parseStTerminated(ansi.SOS, 'X', nil)(buf)
		default:
			n, e := p.parseSequence(buf[1:])
			if k, ok := e.(KeyPressEvent); ok {
				k.Text = ""
				k.Mod |= ModAlt
				return n + 1, k
			}

			// Not a key sequence, nor an alt modified key sequence. In that
			// case, just report a single escape key.
			return 1, KeyPressEvent{Code: KeyEscape}
		}
	case ansi.SS3:
		return p.parseSs3(buf)
	case ansi.DCS:
		return p.parseDcs(buf)
	case ansi.CSI:
		return p.parseCsi(buf)
	case ansi.OSC:
		return p.parseOsc(buf)
	case ansi.APC:
		return p.parseApc(buf)
	case ansi.PM:
		return p.parseStTerminated(ansi.PM, '^', nil)(buf)
	case ansi.SOS:
		return p.parseStTerminated(ansi.SOS, 'X', nil)(buf)
	default:
		if b <= ansi.US || b == ansi.DEL || b == ansi.SP {
			return 1, p.parseControl(b)
		} else if b >= ansi.PAD && b <= ansi.APC {
			// C1 control code
			// UTF-8 never starts with a C1 control code
			// Encode these as Ctrl+Alt+<code - 0x40>
			code := rune(b) - 0x40
			return 1, KeyPressEvent{Code: code, Mod: ModCtrl | ModAlt}
		}
		return p.parseUtf8(buf)
	}
}

func (p *Parser) parseCsi(b []byte) (int, Event) {
	if len(b) == 2 && b[0] == ansi.ESC {
		// short cut if this is an alt+[ key
		return 2, KeyPressEvent{Text: string(rune(b[1])), Mod: ModAlt}
	}

	var cmd ansi.Cmd
	var params [parser.MaxParamsSize]ansi.Param
	var paramsLen int

	var i int
	if b[i] == ansi.CSI || b[i] == ansi.ESC {
		i++
	}
	if i < len(b) && b[i-1] == ansi.ESC && b[i] == '[' {
		i++
	}

	// 初始 CSI 字节
	if i < len(b) && b[i] >= '<' && b[i] <= '?' {
		cmd |= ansi.Cmd(b[i]) << parser.PrefixShift
	}

	// 扫描参数字节在 0x30-0x3F 范围内
	var j int
	for j = 0; i < len(b) && paramsLen < len(params) && b[i] >= 0x30 && b[i] <= 0x3F; i, j = i+1, j+1 {
		if b[i] >= '0' && b[i] <= '9' {
			if params[paramsLen] == parser.MissingParam {
				params[paramsLen] = 0
			}
			params[paramsLen] *= 10
			params[paramsLen] += ansi.Param(b[i]) - '0'
		}
		if b[i] == ':' {
			params[paramsLen] |= parser.HasMoreFlag
		}
		if b[i] == ';' || b[i] == ':' {
			paramsLen++
			if paramsLen < len(params) {
				// 不要溢出 params 切片
				params[paramsLen] = parser.MissingParam
			}
		}
	}

	if j > 0 && paramsLen < len(params) {
		// 有参数
		paramsLen++
	}

	// 扫描中间字节在 0x20-0x2F 范围内
	var intermed byte
	for ; i < len(b) && b[i] >= 0x20 && b[i] <= 0x2F; i++ {
		intermed = b[i]
	}

	// 设置中间字节
	cmd |= ansi.Cmd(intermed) << parser.IntermedShift

	// 扫描最终字节在 0x40-0x7E 范围内
	if i >= len(b) || b[i] < 0x40 || b[i] > 0x7E {
		// URxvt 键的特殊情况
		// CSI <number> $ 是无效序列，但 URxvt 用它来表示 shift 修饰的键。
		if b[i-1] == '$' {
			n, ev := p.parseCsi(append(b[:i-1], '~'))
			if k, ok := ev.(KeyPressEvent); ok {
				k.Mod |= ModShift
				return n, k
			}
		}
		return i, UnknownEvent(b[:i-1])
	}

	// 添加最终字节
	cmd |= ansi.Cmd(b[i])
	i++

	pa := ansi.Params(params[:paramsLen])
	switch cmd {
	case 'y' | '?'<<parser.PrefixShift | '$'<<parser.IntermedShift:
		// 报告模式 (DECRPM)
		mode, _, ok := pa.Param(0, -1)
		if !ok || mode == -1 {
			break
		}
		value, _, ok := pa.Param(1, -1)
		if !ok || value == -1 {
			break
		}
		return i, ModeReportEvent{Mode: ansi.DECMode(mode), Value: ansi.ModeSetting(value)}
	case 'c' | '?'<<parser.PrefixShift:
		// 主要设备属性
		return i, parsePrimaryDevAttrs(pa)
	case 'u' | '?'<<parser.PrefixShift:
		// Kitty 键盘标志
		flags, _, ok := pa.Param(0, -1)
		if !ok || flags == -1 {
			break
		}
		return i, KittyEnhancementsEvent(flags)
	case 'R' | '?'<<parser.PrefixShift:
		// 此报告可能会返回表示页码的第三个参数，但我们实际上不需要它。
		row, _, ok := pa.Param(0, 1)
		if !ok {
			break
		}
		col, _, ok := pa.Param(1, 1)
		if !ok {
			break
		}
		return i, CursorPositionEvent{Y: row - 1, X: col - 1}
	case 'm' | '<'<<parser.PrefixShift, 'M' | '<'<<parser.PrefixShift:
		// 处理 SGR 鼠标
		if paramsLen == 3 {
			return i, parseSGRMouseEvent(cmd, pa)
		}
	case 'm' | '>'<<parser.PrefixShift:
		// XTerm modifyOtherKeys
		mok, _, ok := pa.Param(0, 0)
		if !ok || mok != 4 {
			break
		}
		val, _, ok := pa.Param(1, -1)
		if !ok || val == -1 {
			break
		}
		return i, ModifyOtherKeysEvent(val) //nolint:gosec
	case 'I':
		return i, FocusEvent{}
	case 'O':
		return i, BlurEvent{}
	case 'R':
		// 光标位置报告或修改的 F3
		row, _, rok := pa.Param(0, 1)
		col, _, cok := pa.Param(1, 1)
		if paramsLen == 2 && rok && cok {
			m := CursorPositionEvent{Y: row - 1, X: col - 1}
			if row == 1 && col-1 <= int(ModMeta|ModShift|ModAlt|ModCtrl) {
			// XXX: 当光标在第 1 行时，我们无法区分光标位置报告和 CSI 1 ; <mod> R（修改的 F3）。
			// 在这种情况下，我们报告两种消息。
			//
			// 对于无歧义的光标位置报告，请改用 [ansi.RequestExtendedCursorPosition] (DECXCPR)。
				return i, MultiEvent{KeyPressEvent{Code: KeyF3, Mod: KeyMod(col - 1)}, m}
			}

			return i, m
		}

		if paramsLen != 0 {
			break
		}

		// 未修改的键 F3 (CSI R)
		fallthrough
	case 'a', 'b', 'c', 'd', 'A', 'B', 'C', 'D', 'E', 'F', 'H', 'P', 'Q', 'S', 'Z':
		var k KeyPressEvent
		switch cmd {
		case 'a', 'b', 'c', 'd':
			k = KeyPressEvent{Code: KeyUp + rune(cmd-'a'), Mod: ModShift}
		case 'A', 'B', 'C', 'D':
			k = KeyPressEvent{Code: KeyUp + rune(cmd-'A')}
		case 'E':
			k = KeyPressEvent{Code: KeyBegin}
		case 'F':
			k = KeyPressEvent{Code: KeyEnd}
		case 'H':
			k = KeyPressEvent{Code: KeyHome}
		case 'P', 'Q', 'R', 'S':
			k = KeyPressEvent{Code: KeyF1 + rune(cmd-'P')}
		case 'Z':
			k = KeyPressEvent{Code: KeyTab, Mod: ModShift}
		}
		id, _, _ := pa.Param(0, 1)
		if id == 0 {
			id = 1
		}
		mod, _, _ := pa.Param(1, 1)
		if mod == 0 {
			mod = 1
		}
		if paramsLen > 1 && id == 1 && mod != -1 {
				// CSI 1 ; <修饰键> A
				k.Mod |= KeyMod(mod - 1)
			}
			// 不要忘记处理 Kitty 键盘协议
		return i, parseKittyKeyboardExt(pa, k)
	case 'M':
		// 处理 X10 鼠标
		if i+3 > len(b) {
			return i, UnknownEvent(b[:i])
		}
		return i + 3, parseX10MouseEvent(append(b[:i], b[i:i+3]...))
	case 'y' | '$'<<parser.IntermedShift:
		// 报告模式 (DECRPM)
		mode, _, ok := pa.Param(0, -1)
		if !ok || mode == -1 {
			break
		}
		val, _, ok := pa.Param(1, -1)
		if !ok || val == -1 {
			break
		}
		return i, ModeReportEvent{Mode: ansi.ANSIMode(mode), Value: ansi.ModeSetting(val)}
	case 'u':
		// Kitty 键盘协议 & CSI u (fixterms)
		if paramsLen == 0 {
			return i, UnknownEvent(b[:i])
		}
		return i, parseKittyKeyboard(pa)
	case '_':
		// Win32 输入模式
		if paramsLen != 6 {
			return i, UnknownEvent(b[:i])
		}

		vrc, _, _ := pa.Param(5, 0)
		rc := uint16(vrc) //nolint:gosec
		if rc == 0 {
			rc = 1
		}

		vk, _, _ := pa.Param(0, 0)
		sc, _, _ := pa.Param(1, 0)
		uc, _, _ := pa.Param(2, 0)
		kd, _, _ := pa.Param(3, 0)
		cs, _, _ := pa.Param(4, 0)
		event := p.parseWin32InputKeyEvent(
			nil,
			uint16(vk), //nolint:gosec // Vk wVirtualKeyCode
			uint16(sc), //nolint:gosec // Sc wVirtualScanCode
			rune(uc),   // Uc UnicodeChar
			kd == 1,    // Kd bKeyDown
			uint32(cs), //nolint:gosec // Cs dwControlKeyState
			rc,         // Rc wRepeatCount
		)

		if event == nil {
			return i, UnknownEvent(b[:])
		}

		return i, event
	case '@', '^', '~':
		if paramsLen == 0 {
			return i, UnknownEvent(b[:i])
		}

		param, _, _ := pa.Param(0, 0)
		switch cmd {
		case '~':
			switch param {
			case 27:
				// XTerm modifyOtherKeys 2
				if paramsLen != 3 {
					return i, UnknownEvent(b[:i])
				}
				return i, parseXTermModifyOtherKeys(pa)
			case 200:
					// 括号粘贴开始
					return i, PasteStartEvent{}
				case 201:
					// 括号粘贴结束
					return i, PasteEndEvent{}
			}
		}

		switch param {
		case 1, 2, 3, 4, 5, 6, 7, 8,
			11, 12, 13, 14, 15,
			17, 18, 19, 20, 21,
			23, 24, 25, 26,
			28, 29, 31, 32, 33, 34:
			var k KeyPressEvent
			switch param {
			case 1:
				if p.flags&FlagFind != 0 {
					k = KeyPressEvent{Code: KeyFind}
				} else {
					k = KeyPressEvent{Code: KeyHome}
				}
			case 2:
				k = KeyPressEvent{Code: KeyInsert}
			case 3:
				k = KeyPressEvent{Code: KeyDelete}
			case 4:
				if p.flags&FlagSelect != 0 {
					k = KeyPressEvent{Code: KeySelect}
				} else {
					k = KeyPressEvent{Code: KeyEnd}
				}
			case 5:
				k = KeyPressEvent{Code: KeyPgUp}
			case 6:
				k = KeyPressEvent{Code: KeyPgDown}
			case 7:
				k = KeyPressEvent{Code: KeyHome}
			case 8:
				k = KeyPressEvent{Code: KeyEnd}
			case 11, 12, 13, 14, 15:
				k = KeyPressEvent{Code: KeyF1 + rune(param-11)}
			case 17, 18, 19, 20, 21:
				k = KeyPressEvent{Code: KeyF6 + rune(param-17)}
			case 23, 24, 25, 26:
				k = KeyPressEvent{Code: KeyF11 + rune(param-23)}
			case 28, 29:
				k = KeyPressEvent{Code: KeyF15 + rune(param-28)}
			case 31, 32, 33, 34:
				k = KeyPressEvent{Code: KeyF17 + rune(param-31)}
			}

			// 修饰键
			mod, _, _ := pa.Param(1, -1)
			if paramsLen > 1 && mod != -1 {
				k.Mod |= KeyMod(mod - 1)
			}

			// 处理 URxvt 奇怪的键
			switch cmd {
			case '~':
				// 不要忘记处理 Kitty 键盘协议
				return i, parseKittyKeyboardExt(pa, k)
			case '^':
				k.Mod |= ModCtrl
			case '@':
				k.Mod |= ModCtrl | ModShift
			}

			return i, k
		}

	case 't':
		param, _, ok := pa.Param(0, 0)
		if !ok {
			break
		}

		var winop WindowOpEvent
		winop.Op = param
		for j := 1; j < paramsLen; j++ {
			val, _, ok := pa.Param(j, 0)
			if ok {
				winop.Args = append(winop.Args, val)
			}
		}

		return i, winop
	}
	return i, UnknownEvent(b[:i])
}

// parseSs3 解析 SS3 序列。
// 请参阅 https://vt100.net/docs/vt220-rm/chapter4.html#S4.4.4.2
func (p *Parser) parseSs3(b []byte) (int, Event) {
	if len(b) == 2 && b[0] == ansi.ESC {
		// 如果这是 alt+O 键的快捷键
		return 2, KeyPressEvent{Code: rune(b[1]), Mod: ModAlt}
	}

	var i int
	if b[i] == ansi.SS3 || b[i] == ansi.ESC {
		i++
	}
	if i < len(b) && b[i-1] == ansi.ESC && b[i] == 'O' {
		i++
	}

	// 扫描 0-9 的数字
	var mod int
	for ; i < len(b) && b[i] >= '0' && b[i] <= '9'; i++ {
		mod *= 10
		mod += int(b[i]) - '0'
	}

	// 扫描 GL 字符
	// GL 字符是范围在 0x21-0x7E 之间的单个字节
	// 请参阅 https://vt100.net/docs/vt220-rm/chapter2.html#S2.3.2
	if i >= len(b) || b[i] < 0x21 || b[i] > 0x7E {
		return i, UnknownEvent(b[:i])
	}

	// GL 字符
	gl := b[i]
	i++

	var k KeyPressEvent
	switch gl {
	case 'a', 'b', 'c', 'd':
		k = KeyPressEvent{Code: KeyUp + rune(gl-'a'), Mod: ModCtrl}
	case 'A', 'B', 'C', 'D':
		k = KeyPressEvent{Code: KeyUp + rune(gl-'A')}
	case 'E':
		k = KeyPressEvent{Code: KeyBegin}
	case 'F':
		k = KeyPressEvent{Code: KeyEnd}
	case 'H':
		k = KeyPressEvent{Code: KeyHome}
	case 'P', 'Q', 'R', 'S':
		k = KeyPressEvent{Code: KeyF1 + rune(gl-'P')}
	case 'M':
		k = KeyPressEvent{Code: KeyKpEnter}
	case 'X':
		k = KeyPressEvent{Code: KeyKpEqual}
	case 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y':
		k = KeyPressEvent{Code: KeyKpMultiply + rune(gl-'j')}
	default:
		return i, UnknownEvent(b[:i])
	}

	// 处理奇怪的 SS3 <修饰键> 功能
	if mod > 0 {
		k.Mod |= KeyMod(mod - 1)
	}

	return i, k
}

func (p *Parser) parseOsc(b []byte) (int, Event) {
	defaultKey := func() KeyPressEvent {
		return KeyPressEvent{Code: rune(b[1]), Mod: ModAlt}
	}
	if len(b) == 2 && b[0] == ansi.ESC {
		// 如果这是 alt+] 键的快捷键
		return 2, defaultKey()
	}

	var i int
	if b[i] == ansi.OSC || b[i] == ansi.ESC {
		i++
	}
	if i < len(b) && b[i-1] == ansi.ESC && b[i] == ']' {
		i++
	}

	// 解析 OSC 命令
	// OSC 序列由 BEL、ESC 或 ST 字符终止
	var start, end int
	cmd := -1
	for ; i < len(b) && b[i] >= '0' && b[i] <= '9'; i++ {
		if cmd == -1 {
			cmd = 0
		} else {
			cmd *= 10
		}
		cmd += int(b[i]) - '0'
	}

	if i < len(b) && b[i] == ';' {
		// 标记序列数据的开始
		i++
		start = i
	}

	for ; i < len(b); i++ {
		// 前进到序列的末尾
		if slices.Contains([]byte{ansi.BEL, ansi.ESC, ansi.ST, ansi.CAN, ansi.SUB}, b[i]) {
			break
		}
	}

	if i >= len(b) {
		return i, UnknownEvent(b[:i])
	}

	end = i // 序列数据的末尾
	i++

	// 检查 7 位 ST（字符串终止符）字符
	switch b[i-1] {
	case ansi.CAN, ansi.SUB:
		return i, UnknownEvent(b[:i])
	case ansi.ESC:
		if i >= len(b) || b[i] != '\\' {
			if cmd == -1 || (start == 0 && end == 2) {
				return 2, defaultKey()
			}

			// 如果我们没有有效的 ST 终止符，那么这是一个
			// 已取消的序列，应该被忽略。
			return i, UnknownEvent(b[:i])
		}

		i++
	}

	if end <= start {
		return i, UnknownEvent(b[:i])
	}

	data := string(b[start:end])
	switch cmd {
	case 10:
		return i, ForegroundColorEvent{ansi.XParseColor(data)}
	case 11:
		return i, BackgroundColorEvent{ansi.XParseColor(data)}
	case 12:
		return i, CursorColorEvent{ansi.XParseColor(data)}
	case 52:
		parts := strings.Split(data, ";")
		if len(parts) == 0 {
			return i, ClipboardEvent{}
		}
		if len(parts) != 2 || len(parts[0]) < 1 {
			break
		}

		b64 := parts[1]
		bts, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			break
		}

		sel := ClipboardSelection(parts[0][0]) //nolint:unconvert
		return i, ClipboardEvent{Selection: sel, Content: string(bts)}
	}

	return i, UnknownEvent(b[:i])
}

// parseStTerminated 解析由 ST 字符终止的控制序列。
func (p *Parser) parseStTerminated(intro8, intro7 byte, fn func([]byte) Event) func([]byte) (int, Event) {
	defaultKey := func(b []byte) (int, Event) {
		switch intro8 {
		case ansi.SOS:
			return 2, KeyPressEvent{Code: 'x', Mod: ModShift | ModAlt}
		case ansi.PM, ansi.APC:
			return 2, KeyPressEvent{Code: rune(b[1]), Mod: ModAlt}
		}
		return 0, nil
	}
	return func(b []byte) (int, Event) {
		if len(b) == 2 && b[0] == ansi.ESC {
			return defaultKey(b)
		}

		var i int
		if b[i] == intro8 || b[i] == ansi.ESC {
			i++
		}
		if i < len(b) && b[i-1] == ansi.ESC && b[i] == intro7 {
			i++
		}

		// 扫描控制序列
		// 最常见的控制序列由 ST 字符终止
		// ST 是 7 位字符串终止符字符 (ESC \)
		start := i
		for ; i < len(b); i++ {
			if slices.Contains([]byte{ansi.ESC, ansi.ST, ansi.CAN, ansi.SUB}, b[i]) {
				break
			}
		}

		if i >= len(b) {
			return i, UnknownEvent(b[:i])
		}

		end := i // 序列数据的末尾
		i++

		// 检查 7 位 ST（字符串终止符）字符
		switch b[i-1] {
		case ansi.CAN, ansi.SUB:
			return i, UnknownEvent(b[:i])
		case ansi.ESC:
			if i >= len(b) || b[i] != '\\' {
				if start == end {
					return defaultKey(b)
				}

				// 如果我们没有有效的 ST 终止符，那么这是一个
				// 已取消的序列，应该被忽略。
				return i, UnknownEvent(b[:i])
			}

			i++
		}

		// 调用函数解析序列并返回结果
		if fn != nil {
			if e := fn(b[start:end]); e != nil {
				return i, e
			}
		}

		return i, UnknownEvent(b[:i])
	}
}

func (p *Parser) parseDcs(b []byte) (int, Event) {
	if len(b) == 2 && b[0] == ansi.ESC {
		// 如果这是 alt+P 键的快捷键
		return 2, KeyPressEvent{Code: 'p', Mod: ModShift | ModAlt}
	}

	var params [16]ansi.Param
	var paramsLen int
	var cmd ansi.Cmd

	// DCS 序列由 DCS (0x90) 或 ESC P (0x1b 0x50) 引入
	var i int
	if b[i] == ansi.DCS || b[i] == ansi.ESC {
		i++
	}
	if i < len(b) && b[i-1] == ansi.ESC && b[i] == 'P' {
		i++
	}

	// 初始 DCS 字节
	if i < len(b) && b[i] >= '<' && b[i] <= '?' {
		cmd |= ansi.Cmd(b[i]) << parser.PrefixShift
	}

	// 扫描参数字节在 0x30-0x3F 范围内
	var j int
	for j = 0; i < len(b) && paramsLen < len(params) && b[i] >= 0x30 && b[i] <= 0x3F; i, j = i+1, j+1 {
		if b[i] >= '0' && b[i] <= '9' {
			if params[paramsLen] == parser.MissingParam {
				params[paramsLen] = 0
			}
			params[paramsLen] *= 10
			params[paramsLen] += ansi.Param(b[i]) - '0'
		}
		if b[i] == ':' {
			params[paramsLen] |= parser.HasMoreFlag
		}
		if b[i] == ';' || b[i] == ':' {
			paramsLen++
			if paramsLen < len(params) {
				// Don't overflow the params slice
				params[paramsLen] = parser.MissingParam
			}
		}
	}

	if j > 0 && paramsLen < len(params) {
		// 有参数
		paramsLen++
	}

	// 扫描中间字节在 0x20-0x2F 范围内
	var intermed byte
	for j := 0; i < len(b) && b[i] >= 0x20 && b[i] <= 0x2F; i, j = i+1, j+1 {
		intermed = b[i]
	}

	// 设置中间字节
	cmd |= ansi.Cmd(intermed) << parser.IntermedShift

	// 扫描最终字节在 0x40-0x7E 范围内
	if i >= len(b) || b[i] < 0x40 || b[i] > 0x7E {
		return i, UnknownEvent(b[:i])
	}

	// Add the final byte
	cmd |= ansi.Cmd(b[i])
	i++

	start := i // start of the sequence data
	for ; i < len(b); i++ {
		if b[i] == ansi.ST || b[i] == ansi.ESC {
			break
		}
	}

	if i >= len(b) {
		return i, UnknownEvent(b[:i])
	}

	end := i // end of the sequence data
	i++

	// Check 7-bit ST (string terminator) character
	if i < len(b) && b[i-1] == ansi.ESC && b[i] == '\\' {
		i++
	}

	pa := ansi.Params(params[:paramsLen])
	switch cmd {
	case 'r' | '+'<<parser.IntermedShift:
		// XTGETTCAP 响应
		param, _, _ := pa.Param(0, 0)
		switch param {
		case 1: // 1 表示有效响应，0 表示无效响应
			tc := parseTermcap(b[start:end])
			// XXX: 一些终端如 KiTTY 在查询时会报告无效响应，例如使用 "\x1bP+q5463\x1b\\"
			// 发送 "Tc" 的查询会返回 "\x1bP0+r5463\x1b\\"。
			// 规范说无效响应应该采用 DCS 0 + r ST "\x1bP0+r\x1b\\" 的形式
			// 我们忽略无效响应，只将有效响应发送给程序。
			//
			// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
			return i, tc
		}
	case '|' | '>'<<parser.PrefixShift:
		// XTVersion 响应
		return i, TerminalVersionEvent(b[start:end])
	}

	return i, UnknownEvent(b[:i])
}

func (p *Parser) parseApc(b []byte) (int, Event) {
	if len(b) == 2 && b[0] == ansi.ESC {
		// 如果这是 alt+_ 键的快捷键
		return 2, KeyPressEvent{Code: rune(b[1]), Mod: ModAlt}
	}

	// APC 序列由 APC (0x9f) 或 ESC _ (0x1b 0x5f) 引入
	return p.parseStTerminated(ansi.APC, '_', func(b []byte) Event {
		if len(b) == 0 {
			return nil
		}

		switch b[0] {
		case 'G': // Kitty Graphics Protocol
			var g KittyGraphicsEvent
			parts := bytes.Split(b[1:], []byte{';'})
			g.Options.UnmarshalText(parts[0]) //nolint:errcheck,gosec
			if len(parts) > 1 {
				g.Payload = parts[1]
			}
			return g
		}

		return nil
	})(b)
}

func (p *Parser) parseUtf8(b []byte) (int, Event) {
	if len(b) == 0 {
		return 0, nil
	}

	c := b[0]
	if c <= ansi.US || c == ansi.DEL || c == ansi.SP {
		// 控制代码由 parseControl 处理
		return 1, p.parseControl(c)
	} else if c > ansi.US && c < ansi.DEL {
		// ASCII 可打印字符
		code := rune(c)
		k := KeyPressEvent{Code: code, Text: string(code)}
		if unicode.IsUpper(code) {
			// 将大写字母转换为小写 + shift 修饰键
			k.Code = unicode.ToLower(code)
			k.ShiftedCode = code
			k.Mod |= ModShift
		}

		return 1, k
	}

	code, _ := utf8.DecodeRune(b)
	if code == utf8.RuneError {
		return 1, UnknownEvent(b[0])
	}

	cluster, _, _, _ := uniseg.FirstGraphemeCluster(b, -1)
	text := string(cluster)
	for i := range text {
		if i > 0 {
			// 对多符文图形集群使用 [KeyExtended]
			code = KeyExtended
			break
		}
	}

	return len(cluster), KeyPressEvent{Code: code, Text: text}
}

func (p *Parser) parseControl(b byte) Event {
	switch b {
	case ansi.NUL:
		if p.flags&FlagCtrlAt != 0 {
			return KeyPressEvent{Code: '@', Mod: ModCtrl}
		}
		return KeyPressEvent{Code: KeySpace, Mod: ModCtrl}
	case ansi.BS:
		return KeyPressEvent{Code: 'h', Mod: ModCtrl}
	case ansi.HT:
		if p.flags&FlagCtrlI != 0 {
			return KeyPressEvent{Code: 'i', Mod: ModCtrl}
		}
		return KeyPressEvent{Code: KeyTab}
	case ansi.CR:
		if p.flags&FlagCtrlM != 0 {
			return KeyPressEvent{Code: 'm', Mod: ModCtrl}
		}
		return KeyPressEvent{Code: KeyEnter}
	case ansi.ESC:
		if p.flags&FlagCtrlOpenBracket != 0 {
			return KeyPressEvent{Code: '[', Mod: ModCtrl}
		}
		return KeyPressEvent{Code: KeyEscape}
	case ansi.DEL:
		if p.flags&FlagBackspace != 0 {
			return KeyPressEvent{Code: KeyDelete}
		}
		return KeyPressEvent{Code: KeyBackspace}
	case ansi.SP:
		return KeyPressEvent{Code: KeySpace, Text: " "}
	default:
		if b >= ansi.SOH && b <= ansi.SUB {
			// Use lower case letters for control codes
			code := rune(b + 0x60)
			return KeyPressEvent{Code: code, Mod: ModCtrl}
		} else if b >= ansi.FS && b <= ansi.US {
			code := rune(b + 0x40)
			return KeyPressEvent{Code: code, Mod: ModCtrl}
		}
		return UnknownEvent(b)
	}
}
