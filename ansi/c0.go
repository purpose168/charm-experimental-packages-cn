package ansi

// C0 控制字符。
//
// 这些是 ISO 646 (ASCII) 中定义的字符范围 (0x00-0x1F)。
// 参见: https://en.wikipedia.org/wiki/C0_and_C1_control_codes
const (
	// NUL 是空字符 (Caret: ^@, Char: \0)。
	NUL = 0x00
	// SOH 是标题开始字符 (Caret: ^A)。
	SOH = 0x01
	// STX 是文本开始字符 (Caret: ^B)。
	STX = 0x02
	// ETX 是文本结束字符 (Caret: ^C)。
	ETX = 0x03
	// EOT 是传输结束字符 (Caret: ^D)。
	EOT = 0x04
	// ENQ 是询问字符 (Caret: ^E)。
	ENQ = 0x05
	// ACK 是确认字符 (Caret: ^F)。
	ACK = 0x06
	// BEL 是响铃字符 (Caret: ^G, Char: \a)。
	BEL = 0x07
	// BS 是退格字符 (Caret: ^H, Char: \b)。
	BS = 0x08
	// HT 是水平制表符 (Caret: ^I, Char: \t)。
	HT = 0x09
	// LF 是换行符 (Caret: ^J, Char: \n)。
	LF = 0x0A
	// VT 是垂直制表符 (Caret: ^K, Char: \v)。
	VT = 0x0B
	// FF 是换页符 (Caret: ^L, Char: \f)。
	FF = 0x0C
	// CR 是回车符 (Caret: ^M, Char: \r)。
	CR = 0x0D
	// SO 是移出字符 (Caret: ^N)。
	SO = 0x0E
	// SI 是移入字符 (Caret: ^O)。
	SI = 0x0F
	// DLE 是数据链路转义字符 (Caret: ^P)。
	DLE = 0x10
	// DC1 是设备控制 1 字符 (Caret: ^Q)。
	DC1 = 0x11
	// DC2 是设备控制 2 字符 (Caret: ^R)。
	DC2 = 0x12
	// DC3 是设备控制 3 字符 (Caret: ^S)。
	DC3 = 0x13
	// DC4 是设备控制 4 字符 (Caret: ^T)。
	DC4 = 0x14
	// NAK 是否定确认字符 (Caret: ^U)。
	NAK = 0x15
	// SYN 是同步空闲字符 (Caret: ^V)。
	SYN = 0x16
	// ETB 是传输块结束字符 (Caret: ^W)。
	ETB = 0x17
	// CAN 是取消字符 (Caret: ^X)。
	CAN = 0x18
	// EM 是介质结束字符 (Caret: ^Y)。
	EM = 0x19
	// SUB 是替换字符 (Caret: ^Z)。
	SUB = 0x1A
	// ESC 是转义字符 (Caret: ^[, Char: \e)。
	ESC = 0x1B
	// FS 是文件分隔符 (Caret: ^\)。
	FS = 0x1C
	// GS 是组分隔符 (Caret: ^])。
	GS = 0x1D
	// RS 是记录分隔符 (Caret: ^^)。
	RS = 0x1E
	// US 是单元分隔符 (Caret: ^_)。
	US = 0x1F

	// LS0 是锁定移位 0 字符。
	// 这是 [SI] 的别名。
	LS0 = SI
	// LS1 是锁定移位 1 字符。
	// 这是 [SO] 的别名。
	LS1 = SO
)
