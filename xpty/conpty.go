package xpty

import (
	"os/exec"

	"github.com/purpose168/charm-experimental-packages-cn/conpty"
)

// ConPty 是一个 Windows 控制台伪终端。
type ConPty struct {
	*conpty.ConPty
}

var _ Pty = &ConPty{}

// NewConPty 创建一个新的 ConPty。
func NewConPty(width, height int, opts ...PtyOption) (*ConPty, error) {
	var opt Options
	for _, o := range opts {
		o(opt)
	}

	c, err := conpty.New(width, height, opt.Flags)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &ConPty{c}, nil
}

// Name 返回 ConPty 的名称。
func (c *ConPty) Name() string {
	return "windows-pty"
}

// Start 在 ConPty 上启动一个命令。
// 这是对 conpty.Spawn 的封装。
func (c *ConPty) Start(cmd *exec.Cmd) error {
	return c.start(cmd)
}
