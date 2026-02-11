package ansi

import (
	"bytes"

	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// Strip 从字符串中移除 ANSI 转义码。
func Strip(s string) string {
	var (
		buf    bytes.Buffer         // 用于收集可打印字符的缓冲区
		ri     int                  // 符文索引
		rw     int                  // 符文宽度
		pstate = parser.GroundState // 初始状态
	)

	// 这里实现了解析器的一个子集，仅收集符文和可打印字符。
	for i := range len(s) {
		if pstate == parser.Utf8State {
			// 在这个状态下，收集 rw 个字节以在缓冲区中形成有效的符文。
			// 将所有符文字节放入缓冲区后，转换到 GroundState 并重置计数器。
			buf.WriteByte(s[i])
			ri++
			if ri < rw {
				continue
			}
			pstate = parser.GroundState
			ri = 0
			rw = 0
			continue
		}

		state, action := parser.Table.Transition(pstate, s[i])
		switch action {
		case parser.CollectAction:
			if state == parser.Utf8State {
				// 当我们转换到 Utf8State 时会发生此操作。
				rw = utf8ByteLen(s[i])
				buf.WriteByte(s[i])
				ri++
			}
		case parser.PrintAction, parser.ExecuteAction:
			// 收集可打印的 ASCII 和不可打印的字符
			buf.WriteByte(s[i])
		}

		// 转换到下一个状态。
		// Utf8State 在上面单独管理。
		if pstate != parser.Utf8State {
			pstate = state
		}
	}

	return buf.String()
}

// StringWidth 返回字符串在单元格中的宽度。这是字符串在终端中打印时将占用的单元格数。
// 忽略 ANSI 转义码，并考虑宽字符（如东亚字符和表情符号）。
// 此函数将文本视为字形集群的序列。
func StringWidth(s string) int {
	return stringWidth(GraphemeWidth, s)
}

// StringWidthWc 返回字符串在单元格中的宽度。这是字符串在终端中打印时将占用的单元格数。
// 忽略 ANSI 转义码，并考虑宽字符（如东亚字符和表情符号）。
// 此函数将文本视为宽字符和符文的序列。
func StringWidthWc(s string) int {
	return stringWidth(WcWidth, s)
}

func stringWidth(m Method, s string) int {
	if s == "" {
		return 0
	}

	var (
		pstate = parser.GroundState // initial state
		width  int
	)

	for i := 0; i < len(s); i++ {
		state, action := parser.Table.Transition(pstate, s[i])
		if state == parser.Utf8State {
			cluster, w := FirstGraphemeCluster(s[i:], m)
			width += w

			i += len(cluster) - 1
			pstate = parser.GroundState
			continue
		}

		if action == parser.PrintAction {
			width++
		}

		pstate = state
	}

	return width
}
