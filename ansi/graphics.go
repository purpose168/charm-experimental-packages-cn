package ansi

import (
	"bytes"
	"strconv"
	"strings"
)

// SixelGraphics 返回一个将给定的 sixel 图像负载编码为 DCS sixel 序列的序列。
//
//	DCS p1; p2; p3; q [sixel payload] ST
//
// p1 = 像素宽高比，已废弃，被负载中的像素指标替换
//
// p2 = 这应该是用于透明度的 0，但终端似乎不能正确使用它。
// 在我尝试过的所有终端上，值 0 会留下难看的黑色条，而值 1 看起来是正确的。
//
// p3 = 水平网格大小参数。据我所知，所有人都忽略这个参数并使用固定的网格大小。
//
// 参见 https://shuford.invisible-island.net/all_about_sixels.txt
func SixelGraphics(p1, p2, p3 int, payload []byte) string {
	var buf bytes.Buffer

	buf.WriteString("\x1bP")
	if p1 >= 0 {
		buf.WriteString(strconv.Itoa(p1))
	}
	buf.WriteByte(';')
	if p2 >= 0 {
		buf.WriteString(strconv.Itoa(p2))
	}
	if p3 > 0 {
		buf.WriteByte(';')
		buf.WriteString(strconv.Itoa(p3))
	}
	buf.WriteByte('q')
	buf.Write(payload)
	buf.WriteString("\x1b\\")

	return buf.String()
}

// KittyGraphics 返回一个将给定图像编码为 Kitty 图形协议的序列。
//
//	APC G [逗号分隔的选项] ; [base64 编码的有效负载] ST
//
// 参见 https://sw.kovidgoyal.net/kitty/graphics-protocol/
func KittyGraphics(payload []byte, opts ...string) string {
	var buf bytes.Buffer
	buf.WriteString("\x1b_G")
	buf.WriteString(strings.Join(opts, ","))
	if len(payload) > 0 {
		buf.WriteString(";")
		buf.Write(payload)
	}
	buf.WriteString("\x1b\\")
	return buf.String()
}
