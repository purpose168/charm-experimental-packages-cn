package pony

// 边框样式。
const (
	BorderNone    = "none"
	BorderNormal  = "normal"
	BorderRounded = "rounded"
	BorderThick   = "thick"
	BorderDouble  = "double"
	BorderHidden  = "hidden"
)

// 对齐常量。
const (
	AlignmentLeading  = "leading"  // 水平方向：左侧
	AlignmentCenter   = "center"   // 水平和垂直方向：居中
	AlignmentTrailing = "trailing" // 水平方向：右侧
	AlignmentTop      = "top"      // 垂直方向：顶部
	AlignmentBottom   = "bottom"   // 垂直方向：底部
)

// 尺寸约束单位。
const (
	UnitAuto    = "auto"
	UnitMin     = "min"
	UnitMax     = "max"
	UnitPercent = "%"
)

// 下划线样式（匹配 UV）。
const (
	UnderlineNone   = "none"
	UnderlineSingle = "single"
	UnderlineDouble = "double"
	UnderlineCurly  = "curly"
	UnderlineDotted = "dotted"
	UnderlineDashed = "dashed"
	UnderlineSolid  = "solid"
)

// 文本装饰样式。
const (
	DecorationUnderline     = "underline"
	DecorationStrikethrough = "strikethrough"
)

// 字重值。
const (
	FontWeightBold = "bold"
)

// 字体样式值。
const (
	FontStyleItalic = "italic"
)
