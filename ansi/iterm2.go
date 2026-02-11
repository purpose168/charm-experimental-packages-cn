package ansi

import "fmt"

// ITerm2 返回一个使用 iTerm2 专有协议的序列。使用 iterm2 包可以获得更便捷的 API。
//
//	OSC 1337 ; key = value ST
//
// 示例：
//
//	ITerm2(iterm2.File{...})
//
// 参见 https://iterm2.com/documentation-escape-codes.html
// 参见 https://iterm2.com/documentation-images.html
func ITerm2(data any) string {
	return "\x1b]1337;" + fmt.Sprint(data) + "\x07"
}
