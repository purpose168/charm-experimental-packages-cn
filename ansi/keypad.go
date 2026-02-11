package ansi

// 数字键盘应用模式 (DECKPAM) 是一种决定数字键盘发送应用序列还是 ANSI 序列的模式。
//
// 这类似于启用 [DECNKM]。
// 使用 [NumericKeypadMode] 来设置数字键盘模式。
//
//	ESC =
//
// 参见: https://vt100.net/docs/vt510-rm/DECKPAM.html
const (
	KeypadApplicationMode = "\x1b="
	DECKPAM               = KeypadApplicationMode
)

// 数字键盘数字模式 (DECKPNM) 是一种决定数字键盘发送应用序列还是 ANSI 序列的模式。
//
// 这与禁用 [DECNKM] 效果相同。
//
//	ESC >
//
// 参见: https://vt100.net/docs/vt510-rm/DECKPNM.html
const (
	KeypadNumericMode = "\x1b>"
	DECKPNM           = KeypadNumericMode
)
