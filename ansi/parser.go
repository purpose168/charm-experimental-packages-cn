package ansi

import (
	"unicode/utf8"
	"unsafe"

	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// Parser 表示一个兼容 DEC ANSI 的序列解析器。
//
// 它使用状态机来解析 ANSI 转义序列和控制字符。该解析器设计用于终端模拟器或
// 需要解析 ANSI 转义序列和控制字符的类似应用程序。
// 有关更多信息，请参见 [parser] 包。
//
//go:generate go run ./gen.go
type Parser struct {
	handler Handler

	// params 包含序列的原始参数。
	// 这些参数用于构建 CSI 和 DCS 序列。
	params []int

	// data 包含序列的原始数据。
	// 这些数据用于构建 OSC、DCS、SOS、PM 和 APC 序列。
	data []byte

	// dataLen 跟踪数据缓冲区的长度。
	// 如果 dataLen 为 -1，则数据缓冲区是无限的，会根据需要增长。
	// 否则，dataLen 受数据缓冲区大小的限制。
	dataLen int

	// paramsLen 跟踪参数的数量。
	// 这受 params 缓冲区大小的限制。
	//
	// 这也用于在收集 UTF-8 符文时跟踪已收集的符文字节数。
	paramsLen int

	// cmd 包含原始命令以及序列的私有前缀和中间字节。
	// 第一个低字节包含命令字节，下一个字节包含私有前缀，
	// 再下一个字节包含中间字节。
	//
	// 这也用于在收集 UTF-8 符文时将其视为 4 字节的切片。
	cmd int

	// state 是解析器的当前状态。
	state byte
}

// NewParser 返回一个具有默认设置的新解析器。
// [Parser] 使用默认大小为 32 的参数缓冲区和 64KB 的数据缓冲区。
// 使用 [Parser.SetParamsSize] 和 [Parser.SetDataSize] 分别设置参数和数据缓冲区的大小。
func NewParser() *Parser {
	p := new(Parser)
	p.SetParamsSize(parser.MaxParamsSize)
	p.SetDataSize(1024 * 64) // 64KB data buffer
	return p
}

// SetParamsSize 设置参数缓冲区的大小。
// 这用于构建 CSI 和 DCS 序列。
func (p *Parser) SetParamsSize(size int) {
	p.params = make([]int, size)
}

// SetDataSize 设置数据缓冲区的大小。
// 这用于构建 OSC、DCS、SOS、PM 和 APC 序列。
// 如果 size 小于或等于 0，则数据缓冲区是无限的，会根据需要增长。
func (p *Parser) SetDataSize(size int) {
	if size <= 0 {
		size = 0
		p.dataLen = -1
	}
	p.data = make([]byte, size)
}

// Params 返回解析后的打包参数列表。
func (p *Parser) Params() Params {
	return unsafe.Slice((*Param)(unsafe.Pointer(&p.params[0])), p.paramsLen)
}

// Param 返回给定索引处的参数，如果参数缺失则返回默认值。
// 如果索引超出范围，则返回默认值和 false。
func (p *Parser) Param(i, def int) (int, bool) {
	if i < 0 || i >= p.paramsLen {
		return def, false
	}
	return Param(p.params[i]).Param(def), true
}

// Command 返回最后分发的序列的打包命令。
// 使用 [Cmd] 来解包命令。
func (p *Parser) Command() int {
	return p.cmd
}

// Rune 返回最后分发的序列作为一个 Unicode 符文。
func (p *Parser) Rune() rune {
	rw := utf8ByteLen(byte(p.cmd & 0xff))
	if rw == -1 {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRune((*[utf8.UTFMax]byte)(unsafe.Pointer(&p.cmd))[:rw])
	return r
}

// Control 返回最后分发的序列作为控制代码。
func (p *Parser) Control() byte {
	return byte(p.cmd & 0xff)
}

// Data 返回最后分发的序列的原始数据。
func (p *Parser) Data() []byte {
	return p.data[:p.dataLen]
}

// Reset 将解析器重置为初始状态。
func (p *Parser) Reset() {
	p.clear()
	p.state = parser.GroundState
}

// clear 清除解析器参数和命令。
func (p *Parser) clear() {
	if len(p.params) > 0 {
		p.params[0] = parser.MissingParam
	}
	p.paramsLen = 0
	p.cmd = 0
}

// State 返回解析器的当前状态。
func (p *Parser) State() parser.State {
	return p.state
}

// StateName 返回当前状态的名称。
func (p *Parser) StateName() string {
	return parser.StateNames[p.state]
}

// Parse 解析给定的字节缓冲区。
// 已弃用：请遍历缓冲区并调用 [Parser.Advance] 代替。
func (p *Parser) Parse(b []byte) {
	for i := range b {
		p.Advance(b[i])
	}
}

// Advance 使用给定的字节推进解析器。它返回解析器执行的操作。
func (p *Parser) Advance(b byte) parser.Action {
	switch p.state {
	case parser.Utf8State:
		// 处理 UTF-8 字符
		return p.advanceUtf8(b)
	default:
		return p.advance(b)
	}
}

// collectRune 收集 UTF-8 符文字节
func (p *Parser) collectRune(b byte) {
	if p.paramsLen >= utf8.UTFMax {
		return
	}

	shift := p.paramsLen * 8
	p.cmd &^= 0xff << shift
	p.cmd |= int(b) << shift
	p.paramsLen++
}

// advanceUtf8 处理 UTF-8 字符的解析
func (p *Parser) advanceUtf8(b byte) parser.Action {
	// 收集 UTF-8 符文字节
	p.collectRune(b)
	rw := utf8ByteLen(byte(p.cmd & 0xff))
	if rw == -1 {
		// 这里 panic 是因为第一个字节来自状态机，
		// 如果发生 panic，说明状态机存在 bug！
		panic("无效的符文") // 不可达
	}

	if p.paramsLen < rw {
		return parser.CollectAction
	}

	// 有足够的字节，可以使用 unsafe 解码符文
	if p.handler.Print != nil {
		p.handler.Print(p.Rune())
	}

	p.state = parser.GroundState
	p.paramsLen = 0

	return parser.PrintAction
}

// advance 处理非 UTF-8 字符的解析
func (p *Parser) advance(b byte) parser.Action {
	state, action := parser.Table.Transition(p.state, b)

	// 如果状态从 EscapeState 改变，需要清除解析器状态。
	// 这是因为当我们进入 EscapeState 时，我们没有机会清除解析器状态。
	// 例如，当序列以 ST (\x1b\\ 或 \x9c) 终止时，我们分发当前序列并转换到
	// EscapeState。然而，在这种情况下解析器状态没有被清除，
	// 我们需要在这里清除它，然后再分发 esc 序列。
	if p.state != state {
		if p.state == parser.EscapeState {
			p.performAction(parser.ClearAction, state, b)
		}
		if action == parser.PutAction &&
			p.state == parser.DcsEntryState && state == parser.DcsStringState {
			// XXX: 这是一个特殊情况，我们需要开始收集
			// 非字符串参数化数据，即不遵循 ECMA-48 § 5.4.1 字符串参数格式的数据。
			p.performAction(parser.StartAction, state, 0)
		}
	}

	// 处理特殊情况
	switch {
	case b == ESC && p.state == parser.EscapeState:
		// 连续两个 ESC
		p.performAction(parser.ExecuteAction, state, b)
	default:
		p.performAction(action, state, b)
	}

	p.state = state

	return action
}

func (p *Parser) parseStringCmd() {
	// 尝试解析命令
	datalen := len(p.data)
	if p.dataLen >= 0 {
		datalen = p.dataLen
	}
	for i := range datalen {
		d := p.data[i]
		if d < '0' || d > '9' {
			break
		}
		if p.cmd == parser.MissingCommand {
			p.cmd = 0
		}
		p.cmd *= 10
		p.cmd += int(d - '0')
	}
}

func (p *Parser) performAction(action parser.Action, state parser.State, b byte) {
	switch action {
	case parser.IgnoreAction:
		break

	case parser.ClearAction:
		p.clear()

	case parser.PrintAction:
		p.cmd = int(b)
		if p.handler.Print != nil {
			p.handler.Print(rune(b))
		}

	case parser.ExecuteAction:
		p.cmd = int(b)
		if p.handler.Execute != nil {
			p.handler.Execute(b)
		}

	case parser.PrefixAction:
		// 收集私有前缀
		// 我们只存储最后一个前缀
		p.cmd &^= 0xff << parser.PrefixShift
		p.cmd |= int(b) << parser.PrefixShift

	case parser.CollectAction:
		if state == parser.Utf8State {
			// 重置 UTF-8 计数器
			p.paramsLen = 0
			p.collectRune(b)
		} else {
			// 收集中间字节
			// 我们只存储最后一个中间字节
			p.cmd &^= 0xff << parser.IntermedShift
			p.cmd |= int(b) << parser.IntermedShift
		}

	case parser.ParamAction:
		// 收集参数
		if p.paramsLen >= len(p.params) {
			break
		}

		if b >= '0' && b <= '9' {
			if p.params[p.paramsLen] == parser.MissingParam {
				p.params[p.paramsLen] = 0
			}

			p.params[p.paramsLen] *= 10
			p.params[p.paramsLen] += int(b - '0')
		}

		if b == ':' {
			p.params[p.paramsLen] |= parser.HasMoreFlag
		}

		if b == ';' || b == ':' {
			p.paramsLen++
			if p.paramsLen < len(p.params) {
				p.params[p.paramsLen] = parser.MissingParam
			}
		}

	case parser.StartAction:
		if p.dataLen < 0 && p.data != nil {
			p.data = p.data[:0]
		} else {
			p.dataLen = 0
		}
		if p.state >= parser.DcsEntryState && p.state <= parser.DcsStringState {
			// 收集 DCS 的命令字节
			p.cmd |= int(b)
		} else {
			p.cmd = parser.MissingCommand
		}

	case parser.PutAction:
		switch p.state {
		case parser.OscStringState:
			if b == ';' && p.cmd == parser.MissingCommand {
				p.parseStringCmd()
			}
		}

		if p.dataLen < 0 {
			p.data = append(p.data, b)
		} else {
			if p.dataLen < len(p.data) {
				p.data[p.dataLen] = b
				p.dataLen++
			}
		}

	case parser.DispatchAction:
		// 增加最后一个参数
		if p.paramsLen > 0 && p.paramsLen < len(p.params)-1 ||
			p.paramsLen == 0 && len(p.params) > 0 && p.params[0] != parser.MissingParam {
			p.paramsLen++
		}

		if p.state == parser.OscStringState && p.cmd == parser.MissingCommand {
			// 确保我们有 OSC 的命令
			p.parseStringCmd()
		}

		data := p.data
		if p.dataLen >= 0 {
			data = data[:p.dataLen]
		}
		switch p.state {
		case parser.CsiEntryState, parser.CsiParamState, parser.CsiIntermediateState:
			p.cmd |= int(b)
			if p.handler.HandleCsi != nil {
				p.handler.HandleCsi(Cmd(p.cmd), p.Params())
			}
		case parser.EscapeState, parser.EscapeIntermediateState:
			p.cmd |= int(b)
			if p.handler.HandleEsc != nil {
				p.handler.HandleEsc(Cmd(p.cmd))
			}
		case parser.DcsEntryState, parser.DcsParamState, parser.DcsIntermediateState, parser.DcsStringState:
			if p.handler.HandleDcs != nil {
				p.handler.HandleDcs(Cmd(p.cmd), p.Params(), data)
			}
		case parser.OscStringState:
			if p.handler.HandleOsc != nil {
				p.handler.HandleOsc(p.cmd, data)
			}
		case parser.SosStringState:
			if p.handler.HandleSos != nil {
				p.handler.HandleSos(data)
			}
		case parser.PmStringState:
			if p.handler.HandlePm != nil {
				p.handler.HandlePm(data)
			}
		case parser.ApcStringState:
			if p.handler.HandleApc != nil {
				p.handler.HandleApc(data)
			}
		}
	}
}

func utf8ByteLen(b byte) int {
	if b <= 0b0111_1111 { // 0x00-0x7F
		return 1
	} else if b >= 0b1100_0000 && b <= 0b1101_1111 { // 0xC0-0xDF
		return 2
	} else if b >= 0b1110_0000 && b <= 0b1110_1111 { // 0xE0-0xEF
		return 3
	} else if b >= 0b1111_0000 && b <= 0b1111_0111 { // 0xF0-0xF7
		return 4
	}
	return -1
}
