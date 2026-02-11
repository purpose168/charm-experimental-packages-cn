package ansi

import (
	"strconv"
	"strings"
)

// RequestNameVersion (XTVERSION) 是一个控制序列，用于请求终端的名称和版本。
// 它响应一个 DSR 序列，标识终端信息。
//
//	CSI > 0 q
//	DCS > | text ST
//
// 参见 https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys
const (
	RequestNameVersion = "\x1b[>q"
	XTVERSION          = RequestNameVersion
)

// RequestXTVersion 是一个控制序列，用于请求终端的 XTVERSION。它响应一个 DSR 序列来标识版本信息。
//
//	CSI > Ps q
//	DCS > | text ST
//
// 参见 https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys
//
// 已废弃：请改用 [RequestNameVersion]。
const RequestXTVersion = RequestNameVersion

// PrimaryDeviceAttributes (DA1) 是一个控制序列，用于报告终端的主设备属性。
//
//	CSI c
//	CSI 0 c
//	CSI ? Ps ; ... c
//
// 如果未给出属性，或者属性为 0，此函数返回请求序列。否则，它返回响应序列。
//
// 常见属性包括：
//   - 1	132 列
//   - 2	打印机端口
//   - 4	Sixel
//   - 5	选择性擦除
//   - 6	软字符集 (DRCS)
//   - 7	用户定义键 (UDKs)
//   - 8	国家替换字符集 (NRCS)（仅国际终端）
//   - 9	南斯拉夫语 (SCS)
//   - 12	技术字符集
//   - 15	窗口功能
//   - 18	水平滚动
//   - 21	希腊语
//   - 23	土耳其语
//   - 42	ISO Latin-2 字符集
//   - 44	PCTerm
//   - 45	软键映射
//   - 46	ASCII 仿真
//
// 参见 https://vt100.net/docs/vt510-rm/DA1.html
func PrimaryDeviceAttributes(attrs ...int) string {
	if len(attrs) == 0 {
		return RequestPrimaryDeviceAttributes
	} else if len(attrs) == 1 && attrs[0] == 0 {
		return "\x1b[0c"
	}

	as := make([]string, len(attrs))
	for i, a := range attrs {
		as[i] = strconv.Itoa(a)
	}
	return "\x1b[?" + strings.Join(as, ";") + "c"
}

// DA1 是 [PrimaryDeviceAttributes] 的别名。
func DA1(attrs ...int) string {
	return PrimaryDeviceAttributes(attrs...)
}

// RequestPrimaryDeviceAttributes 是一个控制序列，用于请求终端的主设备属性 (DA1)。
//
//	CSI c
//
// 参见 https://vt100.net/docs/vt510-rm/DA1.html
const RequestPrimaryDeviceAttributes = "\x1b[c"

// SecondaryDeviceAttributes (DA2) 是一个控制序列，用于报告终端的次要设备属性。
//
//	CSI > c
//	CSI > 0 c
//	CSI > Ps ; ... c
//
// 参见 https://vt100.net/docs/vt510-rm/DA2.html
func SecondaryDeviceAttributes(attrs ...int) string {
	if len(attrs) == 0 {
		return RequestSecondaryDeviceAttributes
	}

	as := make([]string, len(attrs))
	for i, a := range attrs {
		as[i] = strconv.Itoa(a)
	}
	return "\x1b[>" + strings.Join(as, ";") + "c"
}

// DA2 是 [SecondaryDeviceAttributes] 的别名。
func DA2(attrs ...int) string {
	return SecondaryDeviceAttributes(attrs...)
}

// RequestSecondaryDeviceAttributes 是一个控制序列，用于请求终端的次要设备属性 (DA2)。
//
//	CSI > c
//
// 参见 https://vt100.net/docs/vt510-rm/DA2.html
const RequestSecondaryDeviceAttributes = "\x1b[>c"

// TertiaryDeviceAttributes (DA3) 是一个控制序列，用于报告终端的第三级设备属性。
//
//	CSI = c
//	CSI = 0 c
//	DCS ! | Text ST
//
// 其中 Text 是终端的单元 ID。
//
// 如果未给出单元 ID，或者单元 ID 为 0，此函数返回请求序列。否则，它返回响应序列。
//
// 参见 https://vt100.net/docs/vt510-rm/DA3.html
func TertiaryDeviceAttributes(unitID string) string {
	switch unitID {
	case "":
		return RequestTertiaryDeviceAttributes
	case "0":
		return "\x1b[=0c"
	}

	return "\x1bP!|" + unitID + "\x1b\\"
}

// DA3 是 [TertiaryDeviceAttributes] 的别名。
func DA3(unitID string) string {
	return TertiaryDeviceAttributes(unitID)
}

// RequestTertiaryDeviceAttributes 是一个控制序列，用于请求终端的第三级设备属性 (DA3)。
//
//	CSI = c
//
// 参见 https://vt100.net/docs/vt510-rm/DA3.html
const RequestTertiaryDeviceAttributes = "\x1b[=c"
