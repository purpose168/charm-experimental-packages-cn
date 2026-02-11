package term

import (
	"io"
)

// File 表示一个具有文件描述符并且可以读取、写入和关闭的文件。
type File interface {
	io.ReadWriteCloser
	Fd() uintptr
}
