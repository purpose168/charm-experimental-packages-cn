package cellbuf

import (
	"bytes"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

const nbsp = '\u00a0' // 非断空格字符

// Wrap 返回一个根据指定宽度限制进行换行的字符串，同时保留字符串中的ANSI转义序列。
// 它会尝试在单词边界处换行，但必要时会截断单词。
//
// breakpoints 参数是一个字符串，包含被视为单词换行断点的字符。连字符 (-) 始终被视为断点。
//
// 注意：breakpoints 必须是由1单元格宽度的字符组成的字符串。
func Wrap(s string, limit int, breakpoints string) string {
	//nolint:godox
	// TODO: 一旦 https://github.com/purpose168/lipgloss-cn/pull/489 发布后，使用 [PenWriter]
	// 问题是 [ansi.Wrap] 不跟踪样式和链接状态，所以组合使用时会破坏带样式的空格单元格。
	// 为了解决这个问题，我们对填充和带样式的空白单元格使用非断空格单元格。
	// 由于两种换行方法都尊重非断空格，我们可以使用它们来保留输出中的带样式空格。

	if len(s) == 0 {
		return ""
	}

	if limit < 1 {
		return s
	}

	p := ansi.GetParser()
	defer ansi.PutParser(p)

	var (
		buf             bytes.Buffer    // 结果缓冲区
		word            bytes.Buffer    // 当前单词缓冲区
		space           bytes.Buffer    // 当前空格缓冲区
		style, curStyle Style           // 当前样式和累积样式
		link, curLink   Link            // 当前链接和累积链接
		curWidth        int             // 当前行宽度
		wordLen         int             // 当前单词长度
	)

	// hasBlankStyle 检查当前样式是否为空白样式（仅考虑反向属性、背景色和下划线样式）
	hasBlankStyle := func() bool {
		// 仅跟踪反向属性、背景色和下划线样式
		return !style.Attrs.Contains(ReverseAttr) && style.Bg == nil && style.UlStyle == NoUnderline
	}

	// addSpace 将空格缓冲区内容添加到结果中
	addSpace := func() {
		curWidth += space.Len()
		buf.Write(space.Bytes())
		space.Reset()
	}

	// addWord 将单词缓冲区内容添加到结果中
	addWord := func() {
		if word.Len() == 0 {
			return
		}

		curLink = link
		curStyle = style

		addSpace()
		curWidth += wordLen
		buf.Write(word.Bytes())
		word.Reset()
		wordLen = 0
	}

	// addNewline 添加换行符并重置相关状态
	addNewline := func() {
		if !curStyle.Empty() {
			buf.WriteString(ansi.ResetStyle)
		}
		if !curLink.Empty() {
			buf.WriteString(ansi.ResetHyperlink())
		}
		buf.WriteByte('\n')
		if !curLink.Empty() {
			buf.WriteString(ansi.SetHyperlink(curLink.URL, curLink.Params))
		}
		if !curStyle.Empty() {
			buf.WriteString(curStyle.Sequence())
		}
		curWidth = 0
		space.Reset()
	}

	var state byte
	for len(s) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(s, state, p)
		switch width {
		case 0:
			if ansi.Equal(seq, "\t") { //nolint:nestif
				addWord()
				space.WriteString(seq)
				break
			} else if ansi.Equal(seq, "\n") {
				if wordLen == 0 {
					if curWidth+space.Len() > limit {
						curWidth = 0
					} else {
						// 保留空白字符
						buf.Write(space.Bytes())
					}
					space.Reset()
				}

				addWord()
				addNewline()
				break
			} else if ansi.HasCsiPrefix(seq) && p.Command() == 'm' {
				// SGR 样式序列 [ansi.SGR]
				ReadStyle(p.Params(), &style)
			} else if ansi.HasOscPrefix(seq) && p.Command() == 8 {
				// 超链接序列 [ansi.SetHyperlink]
				ReadLink(p.Data(), &link)
			}

			word.WriteString(seq)
		default:
			if len(seq) == 1 {
				// ASCII 字符
				r, _ := utf8.DecodeRuneInString(seq)
				if r != nbsp && unicode.IsSpace(r) && hasBlankStyle() {
					addWord()
					space.WriteRune(r)
					break
				} else if r == '-' || runeContainsAny(r, breakpoints) {
					addSpace()
					if curWidth+wordLen+width <= limit {
						addWord()
						buf.WriteString(seq)
						curWidth += width
						break
					}
				}
			}

			if wordLen+width > limit {
				// 如果单词太长，强制换行
				addWord()
			}

			word.WriteString(seq)
			wordLen += width

			if curWidth+wordLen+space.Len() > limit {
				addNewline()
			}
		}

		s = s[n:]
		state = newState
	}

	if wordLen == 0 {
		if curWidth+space.Len() > limit {
			curWidth = 0
		} else {
			// 保留空白字符
			buf.Write(space.Bytes())
		}
		space.Reset()
	}

	addWord()

	if !curLink.Empty() {
		buf.WriteString(ansi.ResetHyperlink())
	}
	if !curStyle.Empty() {
		buf.WriteString(ansi.ResetStyle)
	}

	return buf.String()
}

// runeContainsAny 检查符文 r 是否包含在序列 s 中
func runeContainsAny[T string | []rune](r rune, s T) bool {
	return slices.Contains([]rune(s), r)
}
