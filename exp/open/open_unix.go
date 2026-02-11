//go:build !windows
// +build !windows

package open

import (
	"context"
	"os/exec"
)

// buildCmd 构建打开文件的命令。
// 在 Unix 系统上，它会尝试使用以下命令：
// 1. 如果存在 "open" 命令（macOS），使用它
// 2. 如果指定了应用程序，直接使用该应用程序
// 3. 如果存在 "xdg-open" 命令（Linux），使用它
// 4. 如果以上都不存在，返回 nil
func buildCmd(ctx context.Context, app, path string) *exec.Cmd {
	if _, err := exec.LookPath("open"); err == nil {
		var arg []string
		if app != "" {
			arg = append(arg, "-a", app)
		}
		arg = append(arg, path)
		return exec.CommandContext(ctx, "open", arg...)
	}
	if app != "" {
		return exec.CommandContext(ctx, app, path)
	}
	if _, err := exec.LookPath("xdg-open"); err == nil {
		return exec.CommandContext(ctx, "xdg-open", path)
	}
	return nil
}
