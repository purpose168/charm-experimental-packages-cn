package ansi

import "strconv"

// ResetProgressBar 是一个将进度条重置为默认状态（隐藏）的序列。
//
// OSC 9 ; 4 ; 0 BEL
//
// 请参阅：https://learn.microsoft.com/zh-cn/windows/terminal/tutorials/progress-bar-sequences
const ResetProgressBar = "\x1b]9;4;0\x07"

// SetProgressBar 返回一个将进度条设置为特定百分比（0-100）的“默认”状态的序列。
//
// OSC 9 ; 4 ; 1 百分比 BEL
//
// 请参阅：https://learn.microsoft.com/zh-cn/windows/terminal/tutorials/progress-bar-sequences
func SetProgressBar(percentage int) string {
	return "\x1b]9;4;1;" + strconv.Itoa(min(max(0, percentage), 100)) + "\x07"
}

// SetErrorProgressBar 返回一个将进度条设置为特定百分比（0-100）的“错误”状态的序列。
//
// OSC 9 ; 4 ; 2 百分比 BEL
//
// 请参阅：https://learn.microsoft.com/zh-cn/windows/terminal/tutorials/progress-bar-sequences
func SetErrorProgressBar(percentage int) string {
	return "\x1b]9;4;2;" + strconv.Itoa(min(max(0, percentage), 100)) + "\x07"
}

// SetIndeterminateProgressBar 是一个将进度条设置为不确定状态的序列。
//
// OSC 9 ; 4 ; 3 BEL
//
// 请参阅：https://learn.microsoft.com/zh-cn/windows/terminal/tutorials/progress-bar-sequences
const SetIndeterminateProgressBar = "\x1b]9;4;3\x07"

// SetWarningProgressBar 是一个将进度条设置为“警告”状态的序列。
//
// OSC 9 ; 4 ; 4 百分比 BEL
//
// 请参阅：https://learn.microsoft.com/zh-cn/windows/terminal/tutorials/progress-bar-sequences
func SetWarningProgressBar(percentage int) string {
	return "\x1b]9;4;4;" + strconv.Itoa(min(max(0, percentage), 100)) + "\x07"
}
