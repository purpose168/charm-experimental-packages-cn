package input

import "github.com/purpose168/charm-experimental-packages-cn/ansi"

// ClipboardSelection 表示剪贴板选择。最常见的剪贴板选择是 "system"（系统）和 "primary"（主）选择。
type ClipboardSelection = byte

// 剪贴板选择。
const (
	SystemClipboard  ClipboardSelection = ansi.SystemClipboard
	PrimaryClipboard ClipboardSelection = ansi.PrimaryClipboard
)

// ClipboardEvent 是剪贴板读取消息事件。当终端收到 OSC52 剪贴板读取消息事件时，会发出此消息。
type ClipboardEvent struct {
	Content   string
	Selection ClipboardSelection
}

// String 返回剪贴板消息的字符串表示形式。
func (e ClipboardEvent) String() string {
	return e.Content
}
