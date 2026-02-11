package ansi

import (
	"os"
	"strconv"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
)

var wcOptions = &runewidth.Condition{
	EastAsianWidth:     false,
	StrictEmojiNeutral: true,
}

var dwOptions = &displaywidth.Options{
	EastAsianWidth: false,
}

func init() {
	if ea, err := strconv.ParseBool(os.Getenv("RUNEWIDTH_EASTASIAN")); err == nil && ea {
		wcOptions.EastAsianWidth = true
		dwOptions.EastAsianWidth = true
	}
}

// Method 是一个类型，表示渲染器应该如何计算单元格的显示宽度。
type Method uint8

// 显示宽度模式。
const (
	WcWidth Method = iota
	GraphemeWidth
)

// StringWidth 返回字符串在单元格中的宽度。这是字符串在终端中打印时将占用的单元格数量。
// ANSI 转义码会被忽略，宽字符（如东亚字符和表情符号）会被计算在内。
func (m Method) StringWidth(s string) int {
	return stringWidth(m, s)
}

// Truncate 将字符串截断到指定长度，如果字符串长于指定长度，则在末尾添加尾部。
// 此函数可识别 ANSI 转义码，不会破坏它们，并会考虑宽字符（如东亚字符和表情符号）。
func (m Method) Truncate(s string, length int, tail string) string {
	return truncate(m, s, length, tail)
}

// TruncateLeft 将字符串截断到指定长度，如果字符串长于指定长度，则在开头添加前缀。
// 此函数可识别 ANSI 转义码，不会破坏它们，并会考虑宽字符（如东亚字符和表情符号）。
func (m Method) TruncateLeft(s string, length int, prefix string) string {
	return truncateLeft(m, s, length, prefix)
}

// Cut 切割字符串，不添加任何前缀或尾部字符串。此函数可识别 ANSI 转义码，不会破坏它们，
// 并会考虑宽字符（如东亚字符和表情符号）。注意，[left] 参数是包含的，而 [right] 不是。
func (m Method) Cut(s string, left, right int) string {
	return cut(m, s, left, right)
}

// Hardwrap 将字符串或文本块换行到指定的行长度，打破单词边界。这将保留 ANSI 转义码，
// 并会考虑字符串中的宽字符。
// 当 preserveSpace 为 true 时，行首的空格将被保留。
// 这将文本视为 grapheme 序列。
func (m Method) Hardwrap(s string, length int, preserveSpace bool) string {
	return hardwrap(m, s, length, preserveSpace)
}

// Wordwrap 将字符串或文本块换行到指定的行长度，不打破单词边界。这将保留 ANSI 转义码，
// 并会考虑字符串中的宽字符。
// breakpoints 字符串是被视为单词换行断点的字符列表。连字符 (-) 始终被视为断点。
//
// 注意：breakpoints 必须是 1 单元格宽的符文字符组成的字符串。
func (m Method) Wordwrap(s string, length int, breakpoints string) string {
	return wordwrap(m, s, length, breakpoints)
}

// Wrap 将字符串或文本块换行到指定的行长度，必要时打破单词边界。这将保留 ANSI 转义码，
// 并会考虑字符串中的宽字符。breakpoints 字符串是被视为单词换行断点的字符列表。
// 连字符 (-) 始终被视为断点。
//
// 注意：breakpoints 必须是 1 单元格宽的符文字符组成的字符串。
func (m Method) Wrap(s string, length int, breakpoints string) string {
	return wrap(m, s, length, breakpoints)
}

// DecodeSequence 从给定数据中解码第一个 ANSI 转义序列或可打印的 grapheme。
// 它返回序列切片、读取的字节数、每个序列的单元格宽度以及新状态。
//
// 对于控制和转义序列，单元格宽度始终为 0；对于 ASCII 可打印字符，宽度为 1；
// 对于其他 Unicode 字符，则为它们占用的单元格数。它使用 uniseg 包来计算 Unicode 
// grapheme 和字符的宽度。这意味着它总是会进行 grapheme 聚类（模式 2027）。
//
// 作为最后一个参数传递非 nil 的 [*Parser] 将允许解码器收集序列参数、数据和命令。
// 解析器 cmd 将包含打包的命令值，其中包含中间和前缀字符。在 OSC 序列的情况下，
// cmd 将是 OSC 命令编号。使用 [Cmd] 和 [Param] 类型来解包命令中间件和前缀以及参数。
//
// 零 [Cmd] 表示 CSI、DCS 或 ESC 序列无效。此外，检查其他数据序列（如 OSC、DCS 等）
// 的有效性需要检查返回的序列终止字节，如 ST (ESC \\) 和 BEL)。
//
// 我们将命令字节存储在 [Cmd] 的最高有效字节中，前缀字节存储在下一个字节中，
// 中间字节存储在最低有效字节中。这样做是为了避免使用结构体来存储命令及其中间件和前缀。
// 命令字节始终是最低有效字节，即 [Cmd & 0xff]。使用 [Cmd] 类型来解包命令、中间和前缀字节。
// 注意，我们只收集最后一个前缀字符和中间字节。
//
// [p.Params] 切片将包含序列的参数。任何子参数都将设置 [parser.HasMoreFlag]。
// 使用 [Param] 类型来解包参数。
//
// 示例：
//
//	var state byte // 初始状态始终为零 [NormalState]
//	p := NewParser(32, 1024) // 创建一个新的解析器，具有 32 个参数缓冲区和 1024 个数据缓冲区（可选）
//	input := []byte("\x1b[31mHello, World!\x1b[0m")
//	for len(input) > 0 {
//		seq, width, n, newState := DecodeSequence(input, state, p)
//		log.Printf("seq: %q, width: %d", seq, width)
//		state = newState
//		input = input[n:]
//	}
func (m Method) DecodeSequence(data []byte, state byte, p *Parser) (seq []byte, width, n int, newState byte) {
	return decodeSequence(m, data, state, p)
}

// DecodeSequenceInString 从给定数据中解码第一个 ANSI 转义序列或可打印的 grapheme。
// 它返回序列字符串、读取的字节数、每个序列的单元格宽度以及新状态。
//
// 对于控制和转义序列，单元格宽度始终为 0；对于 ASCII 可打印字符，宽度为 1；
// 对于其他 Unicode 字符，则为它们占用的单元格数。它使用 uniseg 包来计算 Unicode 
// grapheme 和字符的宽度。这意味着它总是会进行 grapheme 聚类（模式 2027）。
//
// 作为最后一个参数传递非 nil 的 [*Parser] 将允许解码器收集序列参数、数据和命令。
// 解析器 cmd 将包含打包的命令值，其中包含中间和前缀字符。在 OSC 序列的情况下，
// cmd 将是 OSC 命令编号。使用 [Cmd] 和 [Param] 类型来解包命令中间件和前缀以及参数。
//
// 零 [Cmd] 表示 CSI、DCS 或 ESC 序列无效。此外，检查其他数据序列（如 OSC、DCS 等）
// 的有效性需要检查返回的序列终止字节，如 ST (ESC \\) 和 BEL)。
//
// 我们将命令字节存储在 [Cmd] 的最高有效字节中，前缀字节存储在下一个字节中，
// 中间字节存储在最低有效字节中。这样做是为了避免使用结构体来存储命令及其中间件和前缀。
// 命令字节始终是最低有效字节，即 [Cmd & 0xff]。使用 [Cmd] 类型来解包命令、中间和前缀字节。
// 注意，我们只收集最后一个前缀字符和中间字节。
//
// [p.Params] 切片将包含序列的参数。任何子参数都将设置 [parser.HasMoreFlag]。
// 使用 [Param] 类型来解包参数。
//
// 示例：
//
//	var state byte // 初始状态始终为零 [NormalState]
//	p := NewParser(32, 1024) // 创建一个新的解析器，具有 32 个参数缓冲区和 1024 个数据缓冲区（可选）
//	input := []byte("\x1b[31mHello, World!\x1b[0m")
//	for len(input) > 0 {
//		seq, width, n, newState := DecodeSequenceInString(input, state, p)
//		log.Printf("seq: %q, width: %d", seq, width)
//		state = newState
//		input = input[n:]
//	}
func (m Method) DecodeSequenceInString(data string, state byte, p *Parser) (seq string, width, n int, newState byte) {
	return decodeSequence(m, data, state, p)
}
