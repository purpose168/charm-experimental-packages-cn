package ansi

import (
	"net/url"
	"path"
)

// NotifyWorkingDirectory 返回一个通知终端当前工作目录的序列。
//
//	OSC 7 ; Pt BEL
//
// 其中 Pt 是格式为 "file://[host]/[path]" 的 URL。
// 如果这是本地计算机上的路径，请将 host 设置为 "localhost"。
//
// 参见: https://wezfurlong.org/wezterm/shell-integration.html#osc-7-escape-sequence-to-set-the-working-directory
// 参见: https://iterm2.com/documentation-escape-codes.html#:~:text=RemoteHost%20and%20CurrentDir%3A-,OSC%207,-%3B%20%5BPs%5D%20ST
func NotifyWorkingDirectory(host string, paths ...string) string {
	path := path.Join(paths...)
	u := &url.URL{
		Scheme: "file",
		Host:   host,
		Path:   path,
	}
	return "\x1b]7;" + u.String() + "\x07"
}
