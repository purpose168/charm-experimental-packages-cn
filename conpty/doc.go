// conpty 包实现 Windows 控制台伪终端支持。
//
// https://learn.microsoft.com/zh-cn/windows/console/creating-a-pseudoconsole-session
package conpty

import "errors"

// ErrUnsupported 当当前平台不支持时返回的错误
var ErrUnsupported = errors.New("conpty: 不支持的平台")

// 默认尺寸
const (
	DefaultWidth  = 80
	DefaultHeight = 25
)
