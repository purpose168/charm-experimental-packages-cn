//go:build !windows
// +build !windows

package xpty

import "os/exec"

// start 方法在非 Windows 平台上返回不支持错误
// 此文件仅在非 Windows 平台构建时使用
func (c *ConPty) start(*exec.Cmd) error {
	return ErrUnsupported
}
