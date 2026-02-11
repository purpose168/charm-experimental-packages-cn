//go:build !windows
// +build !windows

package conpty

import (
	"syscall"
)

// ConPty 表示 Windows 控制台伪终端
// https://learn.microsoft.com/zh-cn/windows/console/creating-a-pseudoconsole-session#preparing-the-communication-channels
type ConPty struct{}

// New 创建一个新的 ConPty
// 此函数在非 Windows 平台上不支持
func New(int, int, int) (*ConPty, error) {
	return nil, ErrUnsupported
}

// Size 返回 ConPty 的大小
func (*ConPty) Size() (int, int, error) {
	return 0, 0, ErrUnsupported
}

// Close 关闭 ConPty
func (*ConPty) Close() error {
	return ErrUnsupported
}

// Fd 返回 ConPty 的文件描述符
func (*ConPty) Fd() uintptr {
	return 0
}

// Read 实现 io.Reader 接口
func (*ConPty) Read([]byte) (int, error) {
	return 0, ErrUnsupported
}

// Write 实现 io.Writer 接口
func (*ConPty) Write([]byte) (int, error) {
	return 0, ErrUnsupported
}

// Resize 调整 ConPty 的大小
func (*ConPty) Resize(int, int) error {
	return ErrUnsupported
}

// InPipeReadFd 返回输入管道的读取文件描述符
func (*ConPty) InPipeReadFd() uintptr {
	return 0
}

// InPipeWriteFd 返回输入管道的写入文件描述符
func (*ConPty) InPipeWriteFd() uintptr {
	return 0
}

// OutPipeReadFd 返回输出管道的读取文件描述符
func (*ConPty) OutPipeReadFd() uintptr {
	return 0
}

// OutPipeWriteFd 返回输出管道的写入文件描述符
func (*ConPty) OutPipeWriteFd() uintptr {
	return 0
}

// Spawn 实现 pty 接口
func (c *ConPty) Spawn(name string, args []string, attr *syscall.ProcAttr) (pid int, handle uintptr, err error) {
	return 0, 0, ErrUnsupported
}
