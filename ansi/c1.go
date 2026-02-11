package ansi

// C1 控制字符。
//
// 这些是 ISO 6429 (ECMA-48) 中定义的字符范围 (0x80-0x9F)。
// 参见: https://en.wikipedia.org/wiki/C0_and_C1_control_codes
const (
	// PAD 是填充字符。
	PAD = 0x80
	// HOP 是高字节预置字符。
	HOP = 0x81
	// BPH 是此处允许换行字符。
	BPH = 0x82
	// NBH 是此处禁止换行字符。
	NBH = 0x83
	// IND 是索引字符。
	IND = 0x84
	// NEL 是下一行字符。
	NEL = 0x85
	// SSA 是选定区域开始字符。
	SSA = 0x86
	// ESA 是选定区域结束字符。
	ESA = 0x87
	// HTS 是水平制表设置字符。
	HTS = 0x88
	// HTJ 是带对齐的水平制表符。
	HTJ = 0x89
	// VTS 是垂直制表设置字符。
	VTS = 0x8A
	// PLD 是部分行前进字符。
	PLD = 0x8B
	// PLU 是部分行后退字符。
	PLU = 0x8C
	// RI 是反向索引字符。
	RI = 0x8D
	// SS2 是单移位 2 字符。
	SS2 = 0x8E
	// SS3 是单移位 3 字符。
	SS3 = 0x8F
	// DCS 是设备控制字符串字符。
	DCS = 0x90
	// PU1 是私有使用 1 字符。
	PU1 = 0x91
	// PU2 是私有使用 2 字符。
	PU2 = 0x92
	// STS 是设置传输状态字符。
	STS = 0x93
	// CCH 是取消字符。
	CCH = 0x94
	// MW 是消息等待字符。
	MW = 0x95
	// SPA 是受保护区域开始字符。
	SPA = 0x96
	// EPA 是受保护区域结束字符。
	EPA = 0x97
	// SOS 是字符串开始字符。
	SOS = 0x98
	// SGCI 是单图形字符引入符。
	SGCI = 0x99
	// SCI 是单字符引入符。
	SCI = 0x9A
	// CSI 是控制序列引入符。
	CSI = 0x9B
	// ST 是字符串终止符。
	ST = 0x9C
	// OSC 是操作系统命令字符。
	OSC = 0x9D
	// PM 是隐私消息字符。
	PM = 0x9E
	// APC 是应用程序命令字符。
	APC = 0x9F
)
