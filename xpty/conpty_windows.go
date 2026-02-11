//go:build windows
// +build windows

package xpty

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func (c *ConPty) start(cmd *exec.Cmd) error {
	pid, proc, err := c.Spawn(cmd.Path, cmd.Args, &syscall.ProcAttr{
		Dir: cmd.Dir,
		Env: cmd.Env,
		Sys: cmd.SysProcAttr,
	})
	if err != nil {
		return err //nolint:wrapcheck
	}

	cmd.Process, err = os.FindProcess(pid)
	if err != nil {
		// 如果我们无法通过 os.FindProcess 找到进程，终止该进程
		// 因为我们依赖于该进程对象进行所有进一步的操作。
		if tErr := windows.TerminateProcess(windows.Handle(proc), 1); tErr != nil {
			return fmt.Errorf("在进程未找到后终止进程失败: %w", tErr)
		}
		return fmt.Errorf("启动后找不到进程: %w", err)
	}

	return nil
}
