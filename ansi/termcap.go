package ansi

import (
	"encoding/hex"
	"strings"
)

// XTGETTCAP (RequestTermcap) 请求 Termcap/Terminfo 字符串。
//
//	DCS + q <Pt> ST
//
// 其中 <Pt> 是 Termcap/Terminfo 功能列表，以两位十六进制编码，用分号分隔。
//
// 参见：https://man7.org/linux/man-pages/man5/terminfo.5.html
// 参见：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands
func XTGETTCAP(caps ...string) string {
	if len(caps) == 0 {
		return ""
	}

	s := "\x1bP+q"
	for i, c := range caps {
		if i > 0 {
			s += ";"
		}
		s += strings.ToUpper(hex.EncodeToString([]byte(c)))
	}

	return s + "\x1b\\"
}

// RequestTermcap 是 [XTGETTCAP] 的别名。
func RequestTermcap(caps ...string) string {
	return XTGETTCAP(caps...)
}

// RequestTerminfo 是 [XTGETTCAP] 的别名。
func RequestTerminfo(caps ...string) string {
	return XTGETTCAP(caps...)
}
