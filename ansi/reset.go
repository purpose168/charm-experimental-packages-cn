package ansi

// ResetInitialState (RIS) 将终端重置为初始状态。
//
//	ESC c
//
// 参见：https://vt100.net/docs/vt510-rm/RIS.html
const (
	ResetInitialState = "\x1bc"
	RIS               = ResetInitialState
)
