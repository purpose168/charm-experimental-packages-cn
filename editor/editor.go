// editor 包提供在外部编辑器中打开文件的功能
package editor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

const (
	defaultEditor        = "nano"         // 默认编辑器（非 Windows 平台）
	defaultEditorWindows = "notepad"      // 默认编辑器（Windows 平台）
)

// Option 定义编辑器选项
//
// Option 在某些编辑器中可能有不同的行为，或者在某些编辑器中不被支持
//
// 参数：
//   - editor: 编辑器名称
//   - filename: 文件名
//
// 返回值：
//   - []string: 编辑器参数
//   - bool: 路径是否已包含在参数中

type Option func(editor, filename string) (args []string, pathInArgs bool)

// OpenAtLine 在支持的编辑器中打开文件到指定行号
//
// 已废弃：请使用 LineNumber 代替
func OpenAtLine(n int) Option { return LineNumber(n) }

// LineNumber 在支持的编辑器中打开文件到指定行号。
// 如果 [number] 小于第 1 行，文件将在第 1 行打开
func LineNumber(number int) Option {
	number = max(1, number)
	plusLineEditors := []string{"vi", "vim", "nvim", "nano", "emacs", "kak", "gedit"}
	return func(editor, filename string) ([]string, bool) {
		if slices.Contains(plusLineEditors, editor) {
			return []string{fmt.Sprintf("+%d", number)}, false
		}
		switch editor {
		case "code":
			return []string{"--goto", fmt.Sprintf("%s:%d", filename, number)}, true
		case "hx":
			return []string{fmt.Sprintf("%s:%d", filename, number)}, true
		}
		return nil, false
	}
}

// EndOfLine 在支持的编辑器中打开文件到行尾
func EndOfLine() Option {
	return func(editor, _ string) (args []string, pathInArgs bool) {
		switch editor {
		case "vim", "nvim":
			return []string{"+norm! $"}, false
		}
		return nil, false
	}
}

// AtPosition 在支持的编辑器中打开文件到指定的行和列。
// 如果行或列小于 1，它们将被设置为 1。
// 如果编辑器只支持行号，列将被忽略。
func AtPosition(line, column int) Option {
	line = max(line, 1)
	column = max(column, 1)
	vimLike := []string{"vi", "vim", "nvim"}
	return func(editor, filename string) (args []string, pathInArgs bool) {
		if slices.Contains(vimLike, editor) {
			return []string{fmt.Sprintf("+call cursor(%d,%d)", line, column)}, false
		}
		switch editor {
		case "nano":
			return []string{fmt.Sprintf("+%d,%d", line, column)}, false
		case "emacs", "kak":
			return []string{fmt.Sprintf("+%d:%d", line, column)}, false
		case "gedit":
			return []string{fmt.Sprintf("+%d", line)}, false
		case "code":
			return []string{"--goto", fmt.Sprintf("%s:%d:%d", filename, line, column)}, true
		case "hx":
			return []string{fmt.Sprintf("%s:%d:%d", filename, line, column)}, true
		}
		return nil, false
	}
}

// Cmd 返回一个 *exec.Cmd，使用 $EDITOR 或 nano（如果未设置 $EDITOR）编辑给定路径
// 已废弃：请使用 Command 或 CommandContext 代替
func Cmd(app, path string, options ...Option) (*exec.Cmd, error) {
	return CommandContext(context.Background(), app, path, options...)
}

// Command 返回一个 *exec.Cmd，使用 $EDITOR 或 nano（如果未设置 $EDITOR）编辑给定路径
func Command(app, path string, options ...Option) (*exec.Cmd, error) {
	return CommandContext(context.Background(), app, path, options...)
}

// CommandContext 返回一个 *exec.Cmd，使用 $EDITOR 或 nano（如果未设置 $EDITOR）编辑给定路径
func CommandContext(ctx context.Context, app, path string, options ...Option) (*exec.Cmd, error) {
	if os.Getenv("SNAP_REVISION") != "" {
		return nil, fmt.Errorf("您是通过 Snap 安装的吗？%[1]s 被沙箱化，无法打开编辑器。请使用 Go 或其他包管理器安装 %[1]s 以启用编辑功能", app)
	}

	editor, args := getEditor()
	editorName := filepath.Base(editor)

	needsToAppendPath := true
	for _, opt := range options {
		optArgs, pathInArgs := opt(editorName, path)
		if pathInArgs {
			needsToAppendPath = false
		}
		args = append(args, optArgs...)
	}
	if needsToAppendPath {
		args = append(args, path)
	}

	return exec.CommandContext(ctx, editor, args...), nil
}

// getEditor 获取编辑器命令和参数
// 首先从 EDITOR 环境变量获取，如果未设置则使用默认编辑器
func getEditor() (string, []string) {
	editor := strings.Fields(os.Getenv("EDITOR"))
	if len(editor) > 1 {
		return editor[0], editor[1:]
	}
	if len(editor) == 1 {
		return editor[0], []string{}
	}
	switch runtime.GOOS {
	case "windows":
		return defaultEditorWindows, []string{}
	default:
		return defaultEditor, []string{}
	}
}
