package input

import (
	"fmt"
	"strings"
)

// Event 表示终端事件。
type Event any

// UnknownEvent 表示未知事件。
type UnknownEvent string

// String 返回未知事件的字符串表示形式。
func (e UnknownEvent) String() string {
	return fmt.Sprintf("%q", string(e))
}

// MultiEvent 表示多个消息事件。
type MultiEvent []Event

// String 返回多个消息事件的字符串表示形式。
func (e MultiEvent) String() string {
	var sb strings.Builder
	for _, ev := range e {
		sb.WriteString(fmt.Sprintf("%v\n", ev))
	}
	return sb.String()
}

// WindowSizeEvent 用于报告终端大小。请注意，Windows 不支持通过 SIGWINCH 信号报告大小调整，
// 而是依赖 Windows 控制台 API 来报告窗口大小变化。
type WindowSizeEvent struct {
	Width  int
	Height int
}

// WindowOpEvent 是窗口操作（XTWINOPS）报告事件。这用于报告各种窗口操作，
// 如报告窗口大小或单元格大小。
type WindowOpEvent struct {
	Op   int
	Args []int
}
