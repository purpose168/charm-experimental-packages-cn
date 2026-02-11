package ansi

import (
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// State 表示 ANSI 转义序列解析器的状态，由 [DecodeSequence] 使用。
type State = byte

// ANSI 转义序列状态，由 [DecodeSequence] 使用。
const (
	NormalState State = iota
	PrefixState
	ParamsState
	IntermedState
	EscapeState
	StringState
)

// DecodeSequence 从给定数据中解码第一个 ANSI 转义序列或可打印的字形。它返回序列切片、读取的字节数、每个序列的单元格宽度以及新状态。
//
// 对于控制和转义序列，单元格宽度始终为 0；对于 ASCII 可打印字符，宽度为 1；对于其他 Unicode 字符，则为其占用的单元格数。
// 它使用 uniseg 包计算 Unicode 字形和字符的宽度，这意味着它始终会进行字形聚类（模式 2027）。
//
// 传递非 nil 的 [*Parser] 作为最后一个参数将允许解码器收集序列参数、数据和命令。解析器 cmd 将包含打包的命令值，其中包含中间和前缀字符。
// 在 OSC 序列的情况下，cmd 将是 OSC 命令编号。使用 [Cmd] 和 [Param] 类型来解包命令中间值和前缀以及参数。
//
// 零值 [Cmd] 表示 CSI、DCS 或 ESC 序列无效。此外，检查其他数据序列（如 OSC、DCS 等）的有效性需要检查返回的序列终止字节，如 ST (ESC \\) 和 BEL)。
//
// 我们在 [Cmd] 中按以下方式存储字节：最高有效字节存储命令字节，次高字节存储前缀字节，最低有效字节存储中间字节。
// 这样做是为了避免使用结构体来存储命令及其中间值和前缀。命令字节始终是最低有效字节，即 [Cmd & 0xff]。
// 使用 [Cmd] 类型来解包命令、中间值和前缀字节。请注意，我们只收集最后一个前缀字符和中间字节。
//
// [p.Params] 切片将包含序列的参数。任何子参数都会设置 [parser.HasMoreFlag]。使用 [Param] 类型来解包参数。
//
// 示例：
//
//	var state byte // 初始状态始终为零 [NormalState]
//	p := NewParser(32, 1024) // 创建一个新的解析器，带有 32 个参数缓冲区和 1024 个数据缓冲区（可选）
//	input := []byte("\x1b[31mHello, World!\x1b[0m")
//	for len(input) > 0 {
//		seq, width, n, newState := DecodeSequence(input, state, p)
//		log.Printf("seq: %q, width: %d", seq, width)
//		state = newState
//		input = input[n:]
//	}
//
// 此函数将文本视为字形聚类的序列。
func DecodeSequence[T string | []byte](b T, state byte, p *Parser) (seq T, width int, n int, newState byte) {
	return decodeSequence(GraphemeWidth, b, state, p)
}

// DecodeSequenceWc 从给定数据中解码第一个 ANSI 转义序列或可打印的字形。它返回序列切片、读取的字节数、每个序列的单元格宽度以及新状态。
//
// 对于控制和转义序列，单元格宽度始终为 0；对于 ASCII 可打印字符，宽度为 1；对于其他 Unicode 字符，则为其占用的单元格数。
// 它使用 uniseg 包计算 Unicode 字形和字符的宽度，这意味着它始终会进行字形聚类（模式 2027）。
//
// 传递非 nil 的 [*Parser] 作为最后一个参数将允许解码器收集序列参数、数据和命令。解析器 cmd 将包含打包的命令值，其中包含中间和前缀字符。
// 在 OSC 序列的情况下，cmd 将是 OSC 命令编号。使用 [Cmd] 和 [Param] 类型来解包命令中间值和前缀以及参数。
//
// 零值 [Cmd] 表示 CSI、DCS 或 ESC 序列无效。此外，检查其他数据序列（如 OSC、DCS 等）的有效性需要检查返回的序列终止字节，如 ST (ESC \\) 和 BEL)。
//
// 我们在 [Cmd] 中按以下方式存储字节：最高有效字节存储命令字节，次高字节存储前缀字节，最低有效字节存储中间字节。
// 这样做是为了避免使用结构体来存储命令及其中间值和前缀。命令字节始终是最低有效字节，即 [Cmd & 0xff]。
// 使用 [Cmd] 类型来解包命令、中间值和前缀字节。请注意，我们只收集最后一个前缀字符和中间字节。
//
// [p.Params] 切片将包含序列的参数。任何子参数都会设置 [parser.HasMoreFlag]。使用 [Param] 类型来解包参数。
//
// 示例：
//
//	var state byte // 初始状态始终为零 [NormalState]
//	p := NewParser(32, 1024) // 创建一个新的解析器，带有 32 个参数缓冲区和 1024 个数据缓冲区（可选）
//	input := []byte("\x1b[31mHello, World!\x1b[0m")
//	for len(input) > 0 {
//		seq, width, n, newState := DecodeSequenceWc(input, state, p)
//		log.Printf("seq: %q, width: %d", seq, width)
//		state = newState
//		input = input[n:]
//	}
//
// 此函数将文本视为宽字符和运行符的序列。
func DecodeSequenceWc[T string | []byte](b T, state byte, p *Parser) (seq T, width int, n int, newState byte) {
	return decodeSequence(WcWidth, b, state, p)
}

func decodeSequence[T string | []byte](m Method, b T, state State, p *Parser) (seq T, width int, n int, newState byte) {
	for i := 0; i < len(b); i++ {
		c := b[i]

		switch state {
		case NormalState:
			switch c {
			case ESC:
				if p != nil {
					if len(p.params) > 0 {
						p.params[0] = parser.MissingParam
					}
					p.cmd = 0
					p.paramsLen = 0
					p.dataLen = 0
				}
				state = EscapeState
				continue
			case CSI, DCS:
				if p != nil {
					if len(p.params) > 0 {
						p.params[0] = parser.MissingParam
					}
					p.cmd = 0
					p.paramsLen = 0
					p.dataLen = 0
				}
				state = PrefixState
				continue
			case OSC, APC, SOS, PM:
				if p != nil {
					p.cmd = parser.MissingCommand
					p.dataLen = 0
				}
				state = StringState
				continue
			}

			if p != nil {
				p.dataLen = 0
				p.paramsLen = 0
				p.cmd = 0
			}
			if c > US && c < DEL {
				// ASCII printable characters
				return b[i : i+1], 1, 1, NormalState
			}

			if c <= US || c == DEL || c < 0xC0 {
				// C0 & C1 control characters & DEL
				return b[i : i+1], 0, 1, NormalState
			}

			if utf8.RuneStart(c) {
				seq, width = FirstGraphemeCluster(b, m)
				i += len(seq)
				return b[:i], width, i, NormalState
			}

			// Invalid UTF-8 sequence
			return b[:i], 0, i, NormalState
		case PrefixState:
			if c >= '<' && c <= '?' {
				if p != nil {
					// We only collect the last prefix character.
					p.cmd &^= 0xff << parser.PrefixShift
					p.cmd |= int(c) << parser.PrefixShift
				}
				break
			}

			state = ParamsState
			fallthrough
		case ParamsState:
			if c >= '0' && c <= '9' {
				if p != nil {
					if p.params[p.paramsLen] == parser.MissingParam {
						p.params[p.paramsLen] = 0
					}

					p.params[p.paramsLen] *= 10
					p.params[p.paramsLen] += int(c - '0')
				}
				break
			}

			if c == ':' {
				if p != nil {
					p.params[p.paramsLen] |= parser.HasMoreFlag
				}
			}

			if c == ';' || c == ':' {
				if p != nil {
					p.paramsLen++
					if p.paramsLen < len(p.params) {
						p.params[p.paramsLen] = parser.MissingParam
					}
				}
				break
			}

			state = IntermedState
			fallthrough
		case IntermedState:
			if c >= ' ' && c <= '/' {
				if p != nil {
					p.cmd &^= 0xff << parser.IntermedShift
					p.cmd |= int(c) << parser.IntermedShift
				}
				break
			}

			if p != nil {
				// Increment the last parameter
				if p.paramsLen > 0 && p.paramsLen < len(p.params)-1 ||
					p.paramsLen == 0 && len(p.params) > 0 && p.params[0] != parser.MissingParam {
					p.paramsLen++
				}
			}

			if c >= '@' && c <= '~' {
				if p != nil {
					p.cmd &^= 0xff
					p.cmd |= int(c)
				}

				if HasDcsPrefix(b) {
					// Continue to collect DCS data
					if p != nil {
						p.dataLen = 0
					}
					state = StringState
					continue
				}

				return b[:i+1], 0, i + 1, NormalState
			}

			// Invalid CSI/DCS sequence
			return b[:i], 0, i, NormalState
		case EscapeState:
			switch c {
			case '[', 'P':
				if p != nil {
					if len(p.params) > 0 {
						p.params[0] = parser.MissingParam
					}
					p.paramsLen = 0
					p.cmd = 0
				}
				state = PrefixState
				continue
			case ']', 'X', '^', '_':
				if p != nil {
					p.cmd = parser.MissingCommand
					p.dataLen = 0
				}
				state = StringState
				continue
			}

			if c >= ' ' && c <= '/' {
				if p != nil {
					p.cmd &^= 0xff << parser.IntermedShift
					p.cmd |= int(c) << parser.IntermedShift
				}
				continue
			} else if c >= '0' && c <= '~' {
				if p != nil {
					p.cmd &^= 0xff
					p.cmd |= int(c)
				}
				return b[:i+1], 0, i + 1, NormalState
			}

			// Invalid escape sequence
			return b[:i], 0, i, NormalState
		case StringState:
			switch c {
			case BEL:
				if HasOscPrefix(b) {
					parseOscCmd(p)
					return b[:i+1], 0, i + 1, NormalState
				}
			case CAN, SUB:
				if HasOscPrefix(b) {
					// Ensure we parse the OSC command number
					parseOscCmd(p)
				}

				// Cancel the sequence
				return b[:i], 0, i, NormalState
			case ST:
				if HasOscPrefix(b) {
					// Ensure we parse the OSC command number
					parseOscCmd(p)
				}

				return b[:i+1], 0, i + 1, NormalState
			case ESC:
				if HasStPrefix(b[i:]) {
					if HasOscPrefix(b) {
						// Ensure we parse the OSC command number
						parseOscCmd(p)
					}

					// End of string 7-bit (ST)
					return b[:i+2], 0, i + 2, NormalState
				}

				// Otherwise, cancel the sequence
				return b[:i], 0, i, NormalState
			}

			if p != nil && p.dataLen < len(p.data) {
				p.data[p.dataLen] = c
				p.dataLen++

				// Parse the OSC command number
				if c == ';' && HasOscPrefix(b) {
					parseOscCmd(p)
				}
			}
		}
	}

	return b, 0, len(b), state
}

func parseOscCmd(p *Parser) {
	if p == nil || p.cmd != parser.MissingCommand {
		return
	}
	for j := range p.dataLen {
		d := p.data[j]
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

// Equal 如果给定的字节切片相等，则返回 true。
func Equal[T string | []byte](a, b T) bool {
	return string(a) == string(b)
}

// HasPrefix 如果给定的字节切片有前缀，则返回 true。
func HasPrefix[T string | []byte](b, prefix T) bool {
	return len(b) >= len(prefix) && Equal(b[0:len(prefix)], prefix)
}

// HasSuffix 如果给定的字节切片有后缀，则返回 true。
func HasSuffix[T string | []byte](b, suffix T) bool {
	return len(b) >= len(suffix) && Equal(b[len(b)-len(suffix):], suffix)
}

// HasCsiPrefix 如果给定的字节切片有 CSI 前缀，则返回 true。
func HasCsiPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == CSI) ||
		(len(b) > 1 && b[0] == ESC && b[1] == '[')
}

// HasOscPrefix 如果给定的字节切片有 OSC 前缀，则返回 true。
func HasOscPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == OSC) ||
		(len(b) > 1 && b[0] == ESC && b[1] == ']')
}

// HasApcPrefix 如果给定的字节切片有 APC 前缀，则返回 true。
func HasApcPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == APC) ||
		(len(b) > 1 && b[0] == ESC && b[1] == '_')
}

// HasDcsPrefix 如果给定的字节切片有 DCS 前缀，则返回 true。
func HasDcsPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == DCS) ||
		(len(b) > 1 && b[0] == ESC && b[1] == 'P')
}

// HasSosPrefix 如果给定的字节切片有 SOS 前缀，则返回 true。
func HasSosPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == SOS) ||
		(len(b) > 1 && b[0] == ESC && b[1] == 'X')
}

// HasPmPrefix 如果给定的字节切片有 PM 前缀，则返回 true。
func HasPmPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == PM) ||
		(len(b) > 1 && b[0] == ESC && b[1] == '^')
}

// HasStPrefix 如果给定的字节切片有 ST 前缀，则返回 true。
func HasStPrefix[T string | []byte](b T) bool {
	return (len(b) > 0 && b[0] == ST) ||
		(len(b) > 1 && b[0] == ESC && b[1] == '\\')
}

// HasEscPrefix 如果给定的字节切片有 ESC 前缀，则返回 true。
func HasEscPrefix[T string | []byte](b T) bool {
	return len(b) > 0 && b[0] == ESC
}

// FirstGraphemeCluster 返回给定字符串或字节切片中的第一个字形聚类及其等宽显示宽度。
func FirstGraphemeCluster[T string | []byte](b T, m Method) (T, int) {
	switch b := any(b).(type) {
	case string:
		cluster := graphemes.FromString(b).First()
		if m == WcWidth {
			return T(cluster), wcOptions.StringWidth(cluster)
		}
		return T(cluster), dwOptions.String(cluster)
	case []byte:
		cluster := graphemes.FromBytes(b).First()
		if m == WcWidth {
			return T(cluster), wcOptions.StringWidth(string(cluster))
		}
		return T(cluster), dwOptions.Bytes(cluster)
	}
	panic("unreachable")
}

// Cmd 表示序列命令。这用于打包/解包带有中间和前缀字符的序列命令。这些通常出现在 CSI 和 DCS 序列中。
type Cmd int

// Prefix 返回 CSI 序列的解包前缀字节。
// 这始终是以下字符之一：'<' '=' '>' '?'，范围在 0x3C-0x3F 之间。
// 如果序列没有前缀，则返回零。
func (c Cmd) Prefix() byte {
	return byte(parser.Prefix(int(c)))
}

// Intermediate 返回 CSI 序列的解包中间字节。
// 中间字节的范围是 0x20-0x2F，包括以下字符：' ', '!', '"', '#', '$', '%', '&', '\”, '(', ')', '*', '+', ',', '-', '.', '/'.
// 如果序列没有中间字节，则返回零。
func (c Cmd) Intermediate() byte {
	return byte(parser.Intermediate(int(c)))
}

// Final 返回 CSI 序列的解包命令字节。
func (c Cmd) Final() byte {
	return byte(parser.Command(int(c)))
}

// Command 使用给定的前缀、中间值和最终值打包命令。零字节表示序列没有前缀或中间值。
//
// 前缀的范围是 0x3C-0x3F，即 `<=>?` 之一。
//
// 中间值的范围是 0x20-0x2F，即 `!"#$%&'()*+,-./` 中的任何字符。
//
// 最终字节的范围是 0x40-0x7E，即 `@A–Z[\]^_`a–z{|}~` 范围内的任何字符。
func Command(prefix, inter, final byte) (c int) {
	c = int(final)
	c |= int(prefix) << parser.PrefixShift
	c |= int(inter) << parser.IntermedShift
	return c
}

// Param 表示序列参数。带有子参数的序列参数会设置 HasMoreFlag。这用于从 CSI 和 DCS 序列中解包参数。
type Param int

// Param 返回给定索引处的解包参数。
// 如果参数缺失，则返回默认值。
func (s Param) Param(def int) int {
	p := int(s) & parser.ParamMask
	if p == parser.MissingParam {
		return def
	}
	return p
}

// HasMore 从参数中解包 HasMoreFlag。
func (s Param) HasMore() bool {
	return s&parser.HasMoreFlag != 0
}

// Parameter 使用给定的参数和该参数是否有后续子参数来打包转义码参数。
func Parameter(p int, hasMore bool) (s int) {
	s = p & parser.ParamMask
	if hasMore {
		s |= parser.HasMoreFlag
	}
	return s
}
