package ansi

import (
	"bytes"
)

// ScreenPassthrough 将给定的 ANSI 序列包装在 DCS 直通序列中，发送到外部终端。当在 GNU Screen 中运行时，用于将原始转义序列发送到外部终端。
//
//	DCS <data> ST
//
// 注意：Screen 限制字符串序列的长度为 768 字节（自 2014 年起）。使用零表示无限制，否则将返回的字符串分成指定大小的块。
//
// 请参阅：https://www.gnu.org/software/screen/manual/screen.html#String-Escapes
// 请参阅：https://git.savannah.gnu.org/cgit/screen.git/tree/src/screen.h?id=c184c6ec27683ff1a860c45be5cf520d896fd2ef#n44
func ScreenPassthrough(seq string, limit int) string {
	var b bytes.Buffer
	b.WriteString("\x1bP")
	if limit > 0 {
		for i := 0; i < len(seq); i += limit {
			end := min(i+limit, len(seq))
			b.WriteString(seq[i:end])
			if end < len(seq) {
				b.WriteString("\x1b\\\x1bP")
			}
		}
	} else {
		b.WriteString(seq)
	}
	b.WriteString("\x1b\\")
	return b.String()
}

// TmuxPassthrough 将给定的 ANSI 序列包装在特殊的 DCS 直通序列中，发送到外部终端。当在 Tmux 中运行时，用于将原始转义序列发送到外部终端。
//
//	DCS tmux ; <escaped-data> ST
//
// 其中 <escaped-data> 是将所有 ESC (0x1b) 出现次数加倍的序列，即替换为 ESC ESC (0x1b 0x1b)。
//
// 注意：需要将 `allow-passthrough` 选项设置为 `on`。
//
// 请参阅：https://github.com/tmux/tmux/wiki/FAQ#what-is-the-passthrough-escape-sequence-and-how-do-i-use-it
func TmuxPassthrough(seq string) string {
	var b bytes.Buffer
	b.WriteString("\x1bPtmux;")
	for i := range len(seq) {
		if seq[i] == ESC {
			b.WriteByte(ESC)
		}
		b.WriteByte(seq[i])
	}
	b.WriteString("\x1b\\")
	return b.String()
}
