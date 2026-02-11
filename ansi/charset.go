package ansi

// SelectCharacterSet 将 G 集字符指定符设置为指定的字符集。
//
//	ESC Ps Pd
//
// 其中 Ps 是 G 集字符指定符，Pd 是标识符。
// 对于 94 字符集，指定符可以是以下之一：
//   - ( G0
//   - ) G1
//   - * G2
//   - + G3
//
// 对于 96 字符集，指定符可以是以下之一：
//   - - G1
//   - . G2
//   - / G3
//
// 一些常见的 94 字符集：
//   - 0 DEC 特殊绘图集
//   - A 英国 (UK)
//   - B 美国 (USASCII)
//
// 示例：
//
//	ESC ( B  选择字符集 G0 = 美国 (USASCII)
//	ESC ( 0  选择字符集 G0 = 特殊字符和线条绘图集
//	ESC ) 0  选择字符集 G1 = 特殊字符和线条绘图集
//	ESC * A  选择字符集 G2 = 英国 (UK)
//
// 参见: https://vt100.net/docs/vt510-rm/SCS.html
func SelectCharacterSet(gset byte, charset byte) string {
	return "\x1b" + string(gset) + string(charset)
}

// SCS 是 SelectCharacterSet 的别名。
func SCS(gset byte, charset byte) string {
	return SelectCharacterSet(gset, charset)
}

// LS1R (右锁定移位 1) 将 G1 移入 GR 字符集。
const LS1R = "\x1b~"

// LS2 (锁定移位 2) 将 G2 移入 GL 字符集。
const LS2 = "\x1bn"

// LS2R (右锁定移位 2) 将 G2 移入 GR 字符集。
const LS2R = "\x1b}"

// LS3 (锁定移位 3) 将 G3 移入 GL 字符集。
const LS3 = "\x1bo"

// LS3R (右锁定移位 3) 将 G3 移入 GR 字符集。
const LS3R = "\x1b|"
