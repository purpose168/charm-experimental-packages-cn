package ansi

import (
	"strings"

	"github.com/clipperhouse/displaywidth"
	"github.com/clipperhouse/uax29/v2/graphemes"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// Cut 切割字符串，不添加任何前缀或尾部字符串。此函数
// 可识别 ANSI 转义码且不会破坏它们，并考虑了
// 宽字符（如东亚字符和表情符号）。
// 此函数将文本视为字素序列。
func Cut(s string, left, right int) string {
	return cut(GraphemeWidth, s, left, right)
}

// CutWc 切割字符串，不添加任何前缀或尾部字符串。此函数
// 可识别 ANSI 转义码且不会破坏它们，并考虑了
// 宽字符（如东亚字符和表情符号）。
// 注意 [left] 参数是包含的，而 [right] 不是，
// 也就是说它会返回 `[left, right)` 区间的内容。
//
// 此函数将文本视为宽字符和符文序列。
func CutWc(s string, left, right int) string {
	return cut(WcWidth, s, left, right)
}

func cut(m Method, s string, left, right int) string {
	if right <= left {
		return ""
	}

	truncate := Truncate
	truncateLeft := TruncateLeft
	if m == WcWidth {
		truncate = TruncateWc
		truncateLeft = TruncateWc
	}

	if left == 0 {
		return truncate(s, right, "")
	}
	return truncateLeft(truncate(s, right, ""), left, "")
}

// Truncate 将字符串截断到指定长度，如果字符串长于指定长度，则在末尾添加尾部字符串。
// 此函数可识别 ANSI 转义码且不会破坏它们，并考虑了宽字符（如东亚字符和表情符号）。
// 此函数将文本视为字素序列。
func Truncate(s string, length int, tail string) string {
	return truncate(GraphemeWidth, s, length, tail)
}

// TruncateWc 将字符串截断到指定长度，如果字符串长于指定长度，则在末尾添加尾部字符串。
// 此函数可识别 ANSI 转义码且不会破坏它们，并考虑了宽字符（如东亚字符和表情符号）。
// 此函数将文本视为宽字符和符文序列。
func TruncateWc(s string, length int, tail string) string {
	return truncate(WcWidth, s, length, tail)
}

func truncate(m Method, s string, length int, tail string) string {
	if sw := StringWidth(s); sw <= length {
		return s
	}

	tw := StringWidth(tail)
	length -= tw
	if length < 0 {
		return ""
	}

	var cluster string
	var buf strings.Builder
	curWidth := 0
	ignoring := false
	pstate := parser.GroundState // initial state
	i := 0

	// 这里我们遍历字符串的字节并收集可打印字符和符文。我们还会跟踪字符串在单元格中的宽度。
	//
	// 一旦达到给定长度，我们开始忽略字符，只收集 ANSI 转义码，直到到达字符串末尾。
	for i < len(s) {
		state, action := parser.Table.Transition(pstate, s[i])
		if state == parser.Utf8State {
			// 当我们转换到 Utf8State 时会发生此操作。
			var width int
			cluster, width = FirstGraphemeCluster(s[i:], m)
			// 将索引增加聚类的长度
			i += len(cluster)
			curWidth += width

			// 我们是否在忽略？跳到下一个字节
			if ignoring {
				continue
			}

			// 这会太宽吗？
			// 如果是，写入尾部并停止收集。
			if curWidth > length && !ignoring {
				ignoring = true
				buf.WriteString(tail)
			}

			if curWidth > length {
				continue
			}

			buf.WriteString(cluster)

			// 收集完成，现在我们回到地面状态。
			pstate = parser.GroundState
			continue
		}

		switch action {
		case parser.PrintAction:
			// 这会太宽吗？
			// 如果是，写入尾部并停止收集。
			if curWidth >= length && !ignoring {
				ignoring = true
				buf.WriteString(tail)
			}

			// 如果我们在忽略，跳到下一个字节
			if ignoring {
				i++
				continue
			}

			// 收集可打印的 ASCII
			curWidth++
			fallthrough
		case parser.ExecuteAction:
			// execute action 会是像 \n 这样的字符，如果在切割范围外，
			// 应该被忽略。
			if ignoring {
				i++
				continue
			}
			fallthrough
		default:
			buf.WriteByte(s[i])
			i++
		}

		// 转换到下一个状态。
		pstate = state

		// 一旦达到给定长度，我们开始忽略符文并将尾部写入缓冲区。
		if curWidth > length && !ignoring {
			ignoring = true
			buf.WriteString(tail)
		}
	}

	return buf.String()
}

// TruncateLeft 从左侧截断字符串，移除 n 个字符，
// 如果字符串长于 n，则在开头添加前缀。
// 此函数可识别 ANSI 转义码且不会破坏它们，并考虑了
// 宽字符（如东亚字符和表情符号）。
// 此函数将文本视为字素序列。
func TruncateLeft(s string, n int, prefix string) string {
	return truncateLeft(GraphemeWidth, s, n, prefix)
}

// TruncateLeftWc 从左侧截断字符串，移除 n 个字符，
// 如果字符串长于 n，则在开头添加前缀。
// 此函数可识别 ANSI 转义码且不会破坏它们，并考虑了
// 宽字符（如东亚字符和表情符号）。
// 此函数将文本视为宽字符和符文序列。
func TruncateLeftWc(s string, n int, prefix string) string {
	return truncateLeft(WcWidth, s, n, prefix)
}

func truncateLeft(m Method, s string, n int, prefix string) string {
	if n <= 0 {
		return s
	}

	var cluster string
	var buf strings.Builder
	curWidth := 0
	ignoring := true
	pstate := parser.GroundState
	i := 0

	for i < len(s) {
		if !ignoring {
			buf.WriteString(s[i:])
			break
		}

		state, action := parser.Table.Transition(pstate, s[i])
		if state == parser.Utf8State {
			var width int
			cluster, width = FirstGraphemeCluster(s[i:], m)

			i += len(cluster)
			curWidth += width

			if curWidth > n && ignoring {
				ignoring = false
				buf.WriteString(prefix)
			}

			if curWidth > n {
				buf.WriteString(cluster)
			}

			if ignoring {
				continue
			}

			pstate = parser.GroundState
			continue
		}

		switch action {
		case parser.PrintAction:
			curWidth++

			if curWidth > n && ignoring {
				ignoring = false
				buf.WriteString(prefix)
			}

			if ignoring {
				i++
				continue
			}

			fallthrough
		case parser.ExecuteAction:
			// execute action 会是像 \n 这样的字符，如果在切割范围外，
			// 应该被忽略。
			if ignoring {
				i++
				continue
			}
			fallthrough
		default:
			buf.WriteByte(s[i])
			i++
		}

		pstate = state
		if curWidth > n && ignoring {
			ignoring = false
			buf.WriteString(prefix)
		}
	}

	return buf.String()
}

// ByteToGraphemeRange 接收起始和结束字节位置，并将它们转换为
// 字素感知的字符位置。
// 您可以将此函数与 [Truncate]、[TruncateLeft] 和 [Cut] 一起使用。
func ByteToGraphemeRange(str string, byteStart, byteStop int) (charStart, charStop int) {
	bytePos, charPos := 0, 0
	gr := graphemes.FromString(str)
	for byteStart > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Value())
		charPos += max(1, displaywidth.String(gr.Value()))
	}
	charStart = charPos
	for byteStop > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Value())
		charPos += max(1, displaywidth.String(gr.Value()))
	}
	charStop = charPos
	return charStart, charStop
}
