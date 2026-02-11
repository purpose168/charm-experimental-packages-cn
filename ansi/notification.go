package ansi

import (
	"fmt"
	"strings"
)

// Notify 使用 iTerm 的 OSC 9 发送桌面通知。
//
//	OSC 9 ; Mc ST
//	OSC 9 ; Mc BEL
//
// 其中 Mc 是通知正文。
//
// See: https://iterm2.com/documentation-escape-codes.html
func Notify(s string) string {
	return "\x1b]9;" + s + "\x07"
}

// DesktopNotification 基于可扩展的 OSC 99 转义码发送桌面通知。
//
//	OSC 99 ; <metadata> ; <payload> ST
//	OSC 99 ; <metadata> ; <payload> BEL
//
// 其中 <metadata> 是冒号分隔的键值对列表，<payload> 是通知正文。
//
// See: https://sw.kovidgoyal.net/kitty/desktop-notifications/
func DesktopNotification(payload string, metadata ...string) string {
	return fmt.Sprintf("\x1b]99;%s;%s\x07", strings.Join(metadata, ":"), payload)
}
