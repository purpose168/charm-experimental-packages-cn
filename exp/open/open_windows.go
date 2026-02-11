//go:build windows
// +build windows

package open

import (
	"context"
	"os/exec"
)

// buildCmd 构建打开文件的命令。
// 在 Windows 系统上，它会：
// 1. 如果指定了应用程序，使用 cmd /C start 命令启动该应用程序
// 2. 如果未指定应用程序，使用 rundll32 url.dll,FileProtocolHandler 打开路径
func buildCmd(ctx context.Context, app, path string) *exec.Cmd {
	if app != "" {
		return exec.Command("cmd", "/C", "start", "", app, path)
	}
	return exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", path)
}
