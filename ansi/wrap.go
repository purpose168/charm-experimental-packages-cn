package ansi

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// nbsp 是不换行空格
const nbsp = 0xA0

// Hardwrap 将字符串或文本块包装到给定的行长度，自动断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// 当 preserveSpace 为 true 时，行首的空格将被保留。
// 将文本视为字形（grapheme）序列。
func Hardwrap(s string, limit int, preserveSpace bool) string {
	return hardwrap(GraphemeWidth, s, limit, preserveSpace)
}

// HardwrapWc 将字符串或文本块包装到给定的行长度，自动断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// 当 preserveSpace 为 true 时，行首的空格将被保留。
// 将文本视为宽字符和符文（rune）序列。
func HardwrapWc(s string, limit int, preserveSpace bool) string {
	return hardwrap(WcWidth, s, limit, preserveSpace)
}

// hardwrap 是 Hardwrap 和 HardwrapWc 的通用实现
// m 是宽度计算方法，limit 是行最大长度，preserveSpace 是否保留行首空格
func hardwrap(m Method, s string, limit int, preserveSpace bool) string {
	if limit < 1 {
		return s
	}

	var (
		cluster      []byte               // 当前处理的字符簇
		buf          bytes.Buffer         // 输出缓冲区
		curWidth     int                  // 当前行宽度
		forceNewline bool                 // 是否强制换行
		pstate       = parser.GroundState // 解析器初始状态
		b            = []byte(s)          // 输入字符串的字节切片
	)

	// 添加换行符并重置当前宽度
	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
	}

	i := 0
	for i < len(b) {
		// 获取解析器的状态转换
		state, action := parser.Table.Transition(pstate, b[i])
		if state == parser.Utf8State {
			var width int
			// 获取第一个字形簇及其宽度
			cluster, width = FirstGraphemeCluster(b[i:], m)
			i += len(cluster)

			// 如果加上这个字形簇后超出宽度限制，则换行
			if curWidth+width > limit {
				addNewline()
			}
			// 如果是行首且不保留空格，跳过前导空格
			if !preserveSpace && curWidth == 0 && len(cluster) <= 4 {
				// 跳过行首的空格
				if r, _ := utf8.DecodeRune(cluster); r != utf8.RuneError && unicode.IsSpace(r) {
					pstate = parser.GroundState
					continue
				}
			}

			buf.Write(cluster)
			curWidth += width
			pstate = parser.GroundState
			continue
		}

		switch action {
		case parser.PrintAction, parser.ExecuteAction:
			// 换行符处理
			if b[i] == '\n' {
				addNewline()
				forceNewline = false
				break
			}

			// 如果超出宽度限制则换行
			if curWidth+1 > limit {
				addNewline()
				forceNewline = true
			}

			// 跳过行首的空格
			if curWidth == 0 {
				if !preserveSpace && forceNewline && unicode.IsSpace(rune(b[i])) {
					break
				}
				forceNewline = false
			}

			buf.WriteByte(b[i])
			// PrintAction 增加宽度，ExecuteAction 不增加
			if action == parser.PrintAction {
				curWidth++
			}
		default:
			// 其他动作直接写入字节
			buf.WriteByte(b[i])
		}

		// UTF8 状态在上面单独管理
		if pstate != parser.Utf8State {
			pstate = state
		}
		i++
	}

	return buf.String()
}

// Wordwrap 将字符串或文本块包装到给定的行长度，不断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// breakpoints 参数指定被视为断词点的字符列表。连字符 (-) 始终被视为断词点。
//
// 注意：断点必须是 1 个单元格宽度的符文字符组成的字符串。
//
// 将文本视为字形（grapheme）序列。
func Wordwrap(s string, limit int, breakpoints string) string {
	return wordwrap(GraphemeWidth, s, limit, breakpoints)
}

// WordwrapWc 将字符串或文本块包装到给定的行长度，不断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// breakpoints 参数指定被视为断词点的字符列表。连字符 (-) 始终被视为断词点。
//
// 注意：断点必须是 1 个单元格宽度的符文字符组成的字符串。
//
// 将文本视为宽字符和符文（rune）序列。
func WordwrapWc(s string, limit int, breakpoints string) string {
	return wordwrap(WcWidth, s, limit, breakpoints)
}

// wordwrap 是 Wordwrap 和 WordwrapWc 的通用实现
// m 是宽度计算方法，limit 是行最大长度，breakpoints 是断词点字符
func wordwrap(m Method, s string, limit int, breakpoints string) string {
	if limit < 1 {
		return s
	}

	var (
		cluster  []byte               // 当前处理的字符簇
		buf      bytes.Buffer         // 输出缓冲区
		word     bytes.Buffer         // 当前单词缓冲区
		space    bytes.Buffer         // 空格缓冲区
		curWidth int                  // 当前行宽度
		wordLen  int                  // 单词长度（不含 ANSI 转义码）
		pstate   = parser.GroundState // 解析器初始状态
		b        = []byte(s)          // 输入字符串的字节切片
	)

	// 添加空格到输出
	addSpace := func() {
		curWidth += space.Len()
		buf.Write(space.Bytes())
		space.Reset()
	}

	// 添加单词到输出
	addWord := func() {
		if word.Len() == 0 {
			return
		}

		addSpace()
		curWidth += wordLen
		buf.Write(word.Bytes())
		word.Reset()
		wordLen = 0
	}

	// 添加换行符
	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
		space.Reset()
	}

	i := 0
	for i < len(b) {
		state, action := parser.Table.Transition(pstate, b[i])
		if state == parser.Utf8State { //nolint:nestif
			var width int
			cluster, width = FirstGraphemeCluster(b[i:], m)
			i += len(cluster)

			r, _ := utf8.DecodeRune(cluster)
			// 如果是空格（但不是不换行空格），则结束当前单词并保存空格
			if r != utf8.RuneError && unicode.IsSpace(r) && r != nbsp {
				addWord()
				space.WriteRune(r)
			} else if bytes.ContainsAny(cluster, breakpoints) {
				// 如果遇到断点字符，添加空格和单词，然后写入断点
				addSpace()
				addWord()
				buf.Write(cluster)
				curWidth++
			} else {
				// 将字符添加到当前单词
				word.Write(cluster)
				wordLen += width
				// 如果当前行已满，换行
				if curWidth+space.Len()+wordLen > limit &&
					wordLen < limit {
					addNewline()
				}
			}

			pstate = parser.GroundState
			continue
		}

		switch action {
		case parser.PrintAction, parser.ExecuteAction:
			r := rune(b[i])
			switch {
			case r == '\n':
				// 换行符处理
				if wordLen == 0 {
					if curWidth+space.Len() > limit {
						curWidth = 0
					} else {
						buf.Write(space.Bytes())
					}
					space.Reset()
				}

				addWord()
				addNewline()
			case unicode.IsSpace(r):
				// 空格处理，结束当前单词并保存空格
				addWord()
				space.WriteByte(b[i])
			case r == '-':
				// 连字符处理，作为断点
				fallthrough
			case runeContainsAny(r, breakpoints):
				// 断点字符处理
				addSpace()
				addWord()
				buf.WriteByte(b[i])
				curWidth++
			default:
				// 其他字符添加到单词
				word.WriteByte(b[i])
				wordLen++
				// 如果超出宽度则换行
				if curWidth+space.Len()+wordLen > limit &&
					wordLen < limit {
					addNewline()
				}
			}

		default:
			// 其他动作添加到单词
			word.WriteByte(b[i])
		}

		// UTF8 状态在上面单独管理
		if pstate != parser.Utf8State {
			pstate = state
		}
		i++
	}

	// 处理最后一个单词
	addWord()

	return buf.String()
}

// Wrap 将字符串或文本块包装到给定的行长度，根据需要断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// breakpoints 参数指定被视为断词点的字符列表。连字符 (-) 始终被视为断词点。
//
// 注意：断点必须是 1 个单元格宽度的符文字符组成的字符串。
//
// 将文本视为字形（grapheme）序列。
func Wrap(s string, limit int, breakpoints string) string {
	return wrap(GraphemeWidth, s, limit, breakpoints)
}

// WrapWc 将字符串或文本块包装到给定的行长度，根据需要断词。
// 这将保留 ANSI 转义码，并会考虑字符串中的宽字符。
// breakpoints 参数指定被视为断词点的字符列表。连字符 (-) 始终被视为断词点。
//
// 注意：断点必须是 1 个单元格宽度的符文字符组成的字符串。
//
// 将文本视为宽字符和符文（rune）序列。
func WrapWc(s string, limit int, breakpoints string) string {
	return wrap(WcWidth, s, limit, breakpoints)
}

// wrap 是 Wrap 和 WrapWc 的通用实现
// m 是宽度计算方法，limit 是行最大长度，breakpoints 是断词点字符
func wrap(m Method, s string, limit int, breakpoints string) string {
	if limit < 1 {
		return s
	}

	var (
		cluster    string               // 当前处理的字符簇
		buf        bytes.Buffer         // 输出缓冲区
		word       bytes.Buffer         // 当前单词缓冲区
		space      bytes.Buffer         // 空格缓冲区
		spaceWidth int                  // 空格缓冲区的宽度
		curWidth   int                  // 当前行的已写宽度
		wordLen    int                  // 单词长度（不含 ANSI 转义码）
		pstate     = parser.GroundState // 解析器初始状态
	)

	// 添加空格到输出
	addSpace := func() {
		if spaceWidth == 0 && space.Len() == 0 {
			return
		}
		curWidth += spaceWidth
		buf.Write(space.Bytes())
		space.Reset()
		spaceWidth = 0
	}

	// 添加单词到输出
	addWord := func() {
		if word.Len() == 0 {
			return
		}

		addSpace()
		curWidth += wordLen
		buf.Write(word.Bytes())
		word.Reset()
		wordLen = 0
	}

	// 添加换行符
	addNewline := func() {
		buf.WriteByte('\n')
		curWidth = 0
		space.Reset()
		spaceWidth = 0
	}

	i := 0
	for i < len(s) {
		state, action := parser.Table.Transition(pstate, s[i])
		if state == parser.Utf8State { //nolint:nestif
			var width int
			cluster, width = FirstGraphemeCluster(s[i:], m)
			i += len(cluster)

			r, _ := utf8.DecodeRuneInString(cluster)
			switch {
			case r != utf8.RuneError && unicode.IsSpace(r) && r != nbsp: // nbsp 是不换行空格
				// 如果是空格（但不是不换行空格），则结束当前单词并保存空格
				addWord()
				space.WriteRune(r)
				spaceWidth += width
			case strings.ContainsAny(cluster, breakpoints):
				// 断点字符处理
				addSpace()
				if curWidth+wordLen+width > limit {
					// 如果放不下当前断点，将其作为单词的一部分
					word.WriteString(cluster)
					wordLen += width
				} else {
					addWord()
					buf.WriteString(cluster)
					curWidth += width
				}
			default:
				// 如果单词太长，进行硬换行
				if wordLen+width > limit {
					// 如果单词太长，对其进行硬换行
					addWord()
				}

				word.WriteString(cluster)
				wordLen += width

				// 如果超出宽度则换行
				if curWidth+wordLen+spaceWidth > limit {
					addNewline()
				}

				// 如果单词正好等于宽度限制，进行硬换行
				if wordLen == limit {
					// 如果单词太长，对其进行硬换行
					addWord()
				}
			}

			pstate = parser.GroundState
			continue
		}

		switch action {
		case parser.PrintAction, parser.ExecuteAction:
			switch r := rune(s[i]); {
			case r == '\n':
				// 换行符处理
				if wordLen == 0 {
					if curWidth+spaceWidth > limit {
						curWidth = 0
					} else {
						// 保留空白字符
						buf.Write(space.Bytes())
					}
					space.Reset()
					spaceWidth = 0
				}

				addWord()
				addNewline()
			case unicode.IsSpace(r):
				// 空格处理
				addWord()
				space.WriteRune(r)
				spaceWidth++
			case r == '-':
				// 连字符处理
				fallthrough
			case runeContainsAny(r, breakpoints):
				// 断点字符处理
				addSpace()
				if curWidth+wordLen >= limit {
					// 无法在当前行容纳断点字符，将其作为单词的一部分
					word.WriteRune(r)
					wordLen++
				} else {
					addWord()
					buf.WriteRune(r)
					curWidth++
				}
			default:
				// 如果已达到宽度限制，换行
				if curWidth == limit {
					addNewline()
				}

				word.WriteRune(r)
				wordLen++

				// 如果单词正好等于宽度限制，进行硬换行
				if wordLen == limit {
					// 如果单词太长，对其进行硬换行
					addWord()
				}

				// 如果超出宽度则换行
				if curWidth+wordLen+spaceWidth > limit {
					addNewline()
				}
			}

		default:
			// 其他动作添加到单词
			word.WriteByte(s[i])
		}

		// UTF8 状态在上面单独管理
		if pstate != parser.Utf8State {
			pstate = state
		}
		i++
	}

	// 处理剩余内容
	if wordLen == 0 {
		if curWidth+spaceWidth > limit {
			curWidth = 0
		} else {
			// 保留空白字符
			buf.Write(space.Bytes())
		}
		space.Reset()
		spaceWidth = 0
	}

	// 处理最后一个单词
	addWord()

	return buf.String()
}

// runeContainsAny 检查符文 r 是否包含在字符串 s 中
func runeContainsAny(r rune, s string) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
