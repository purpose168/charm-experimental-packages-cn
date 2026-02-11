package conpty

import (
	"syscall"
)

// pty 接口定义了伪终端的基本操作
// 包含关闭、获取文件描述符、管道操作、读写、调整大小、获取大小和启动进程等方法
type pty interface {
	// Close 关闭伪终端
	Close() error
	// Fd 获取伪终端的文件描述符
	Fd() uintptr
	// InPipeReadFd 获取输入管道的读取文件描述符
	InPipeReadFd() uintptr
	// InPipeWriteFd 获取输入管道的写入文件描述符
	InPipeWriteFd() uintptr
	// OutPipeReadFd 获取输出管道的读取文件描述符
	OutPipeReadFd() uintptr
	// OutPipeWriteFd 获取输出管道的写入文件描述符
	OutPipeWriteFd() uintptr
	// Read 从伪终端读取数据
	Read(p []byte) (n int, err error)
	// Resize 调整伪终端的大小
	Resize(w int, h int) error
	// Size 获取伪终端的当前大小
	Size() (w int, h int, err error)
	// Spawn 在伪终端中启动一个新进程
	Spawn(name string, args []string, attr *syscall.ProcAttr) (pid int, handle uintptr, err error)
	// Write 向伪终端写入数据
	Write(p []byte) (n int, err error)
}

var _ pty = &ConPty{} // 确保 ConPty 实现了 pty 接口
