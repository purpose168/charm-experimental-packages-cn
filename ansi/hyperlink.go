package ansi

import "strings"

// SetHyperlink 返回用于开始超链接的序列。
//
//	OSC 8 ; Params ; Uri ST
//	OSC 8 ; Params ; Uri BEL
//
// 要重置超链接，请省略 URI。
//
// 参考：https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func SetHyperlink(uri string, params ...string) string {
	var p string
	if len(params) > 0 {
		p = strings.Join(params, ":")
	}
	return "\x1b]8;" + p + ";" + uri + "\x07"
}

// ResetHyperlink 返回用于重置超链接的序列。
//
// 这等同于 SetHyperlink("", params...)。
//
// 参考：https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func ResetHyperlink(params ...string) string {
	return SetHyperlink("", params...)
}
