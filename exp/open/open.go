// Package open 提供打开文件和 URL 的功能。
package open

import (
	"context"
	"errors"
	"fmt"
)

// ErrNotSupported 当找不到打开文件的方法时发生。
var ErrNotSupported = errors.New("not supported")

// Open 打开给定的输入。
func Open(ctx context.Context, input string) error {
	return With(ctx, "", input)
}

// With 使用给定的应用程序打开给定的输入。
func With(ctx context.Context, app, input string) error {
	cmd := buildCmd(ctx, app, input)
	if cmd == nil {
		return ErrNotSupported
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("open: %w: %s", err, string(out))
	}
	return nil
}
