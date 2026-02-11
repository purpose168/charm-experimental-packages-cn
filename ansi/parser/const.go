// Package parser 提供 ANSI 转义序列解析功能。
package parser

// Action 是 DEC ANSI 解析器动作。
type Action = byte

// 这些是解析器可以执行的动作。
const (
	NoneAction Action = iota
	ClearAction
	CollectAction
	PrefixAction
	DispatchAction
	ExecuteAction
	StartAction // 数据字符串的开始
	PutAction   // 放入数据字符串
	ParamAction
	PrintAction

	IgnoreAction = NoneAction
)

// ActionNames 为解析器动作提供字符串名称。
var ActionNames = []string{
	"NoneAction",
	"ClearAction",
	"CollectAction",
	"PrefixAction",
	"DispatchAction",
	"ExecuteAction",
	"StartAction",
	"PutAction",
	"ParamAction",
	"PrintAction",
}

// State 是 DEC ANSI 解析器状态。
type State = byte

// 这些是解析器可以处于的状态。
const (
	GroundState State = iota
	CsiEntryState
	CsiIntermediateState
	CsiParamState
	DcsEntryState
	DcsIntermediateState
	DcsParamState
	DcsStringState
	EscapeState
	EscapeIntermediateState
	OscStringState
	SosStringState
	PmStringState
	ApcStringState

	// Utf8State 不是 DEC ANSI 标准的一部分。它用于处理 UTF-8 序列。
	Utf8State
)

// StateNames 为解析器状态提供字符串名称。
var StateNames = []string{
	"GroundState",
	"CsiEntryState",
	"CsiIntermediateState",
	"CsiParamState",
	"DcsEntryState",
	"DcsIntermediateState",
	"DcsParamState",
	"DcsStringState",
	"EscapeState",
	"EscapeIntermediateState",
	"OscStringState",
	"SosStringState",
	"PmStringState",
	"ApcStringState",
	"Utf8State",
}
