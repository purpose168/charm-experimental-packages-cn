// Package xpty 提供了在 Go 中与伪终端（PTY）交互的平台无关接口。
// 它抽象了 Unix 和 Windows 系统之间的差异，同时支持 ConPTY 和经典 Unix PTY。
package xpty

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/creack/pty"
	"github.com/purpose168/charm-experimental-packages-cn/term"
)

// ErrUnsupported 当功能不被支持时返回。
var ErrUnsupported = pty.ErrUnsupported

// Pty 表示一个 PTY（伪终端）接口。
type Pty interface {
	term.File
	io.ReadWriteCloser

	// Resize 调整 PTY 的大小。
	Resize(width, height int) error

	// Size 返回 PTY 的大小。
	Size() (width, height int, err error)

	// Name 返回 PTY 的名称。
	Name() string

	// Start 在 PTY 上启动一个命令。
	// 启动的命令将其标准输入、输出和错误连接到 PTY。
	// 在 Windows 上，调用 Wait 不起作用，因为 Go 运行时不能正确处理 ConPTY 进程。
	// 参见 https://github.com/golang/go/pull/62710。
	Start(cmd *exec.Cmd) error
}

// Options 表示 PTY 选项。
type Options struct {
	Flags int
}

// PtyOption 是一个 PTY 选项。
type PtyOption func(o Options)

// NewPty 创建一个新的 PTY。
//
// 返回的 PTY 在 Unix 系统上是 Unix PTY，在 Windows 上是 ConPTY。
// width 和 height 参数指定 PTY 的初始大小。
// 你可以通过传递 PtyOptions 来为 PTY 传递额外的选项。
//
//	pty, err := xpty.NewPty(80, 24)
//	if err != nil {
//	   // 处理错误
//	}
//
// defer pty.Close() // 确保在完成后关闭 PTY。
// switch pty := pty.(type) {
// case xpty.UnixPty:
//
//	// Unix PTY
//
// case xpty.ConPty:
//
//	    // ConPTY
//	}
func NewPty(width, height int, opts ...PtyOption) (Pty, error) {
	if runtime.GOOS == "windows" {
		return NewConPty(width, height, opts...)
	}
	return NewUnixPty(width, height, opts...)
}

// WaitProcess 等待进程退出。
// 这个函数存在是因为在 Windows 上，cmd.Wait() 不能与 ConPty 一起工作。
// 当操作系统不是 Windows 时，它会简单地回退到 cmd.Wait()。
func WaitProcess(ctx context.Context, cmd *exec.Cmd) (err error) {
	if runtime.GOOS != "windows" {
		return cmd.Wait() //nolint:wrapcheck
	}

	if cmd.Process == nil {
		return errors.New("进程未启动")
	}

	type result struct {
		*os.ProcessState
		error
	}

	donec := make(chan result, 1)
	go func() {
		state, err := cmd.Process.Wait()
		donec <- result{state, err}
	}()

	select {
	case <-ctx.Done():
		err = cmd.Process.Kill()
	case r := <-donec:
		cmd.ProcessState = r.ProcessState
		err = r.error
	}

	return err
}
