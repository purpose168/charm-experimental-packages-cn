package xpty

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// UnixPty 表示一个经典的 Unix PTY（伪终端）。
type UnixPty struct {
	master, slave *os.File
}

var _ Pty = &UnixPty{}

// NewUnixPty 创建一个新的 Unix PTY。
func NewUnixPty(width, height int, _ ...PtyOption) (*UnixPty, error) {
	ptm, pts, err := pty.Open()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	p := &UnixPty{
		master: ptm,
		slave:  pts,
	}

	if width >= 0 && height >= 0 {
		if err := p.Resize(width, height); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Close 实现 XPTY 接口。
func (p *UnixPty) Close() (err error) {
	defer func() {
		serr := p.slave.Close()
		if err == nil {
			err = serr
		}
	}()
	if err := p.master.Close(); err != nil {
		return err //nolint:wrapcheck
	}
	return err
}

// Fd 实现 XPTY 接口。
func (p *UnixPty) Fd() uintptr {
	return p.master.Fd()
}

// Name 实现 XPTY 接口。
func (p *UnixPty) Name() string {
	return p.master.Name()
}

// SlaveName 返回从 PTY 的名称。
// 这通常用于远程会话中识别正在运行的 TTY。您可以在 SSH 会话中找到它，定义为 $SSH_TTY。
func (p *UnixPty) SlaveName() string {
	return p.slave.Name()
}

// Read 实现 XPTY 接口。
func (p *UnixPty) Read(b []byte) (n int, err error) {
	return p.master.Read(b) //nolint:wrapcheck
}

// Resize 实现 XPTY 接口。
func (p *UnixPty) Resize(width int, height int) (err error) {
	return p.setWinsize(width, height, 0, 0)
}

// SetWinsize 设置 PTY 的窗口大小。
func (p *UnixPty) SetWinsize(width, height, x, y int) error {
	return p.setWinsize(width, height, x, y)
}

// Size 返回 PTY 的大小。
func (p *UnixPty) Size() (width, height int, err error) {
	return p.size()
}

// Start 实现 XPTY 接口。
func (p *UnixPty) Start(c *exec.Cmd) error {
	if c.Stdout == nil {
		c.Stdout = p.slave
	}
	if c.Stderr == nil {
		c.Stderr = p.slave
	}
	if c.Stdin == nil {
		c.Stdin = p.slave
	}
	if err := c.Start(); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

// Write 实现 XPTY 接口。
func (p *UnixPty) Write(b []byte) (n int, err error) {
	return p.master.Write(b) //nolint:wrapcheck
}

// Master 返回 PTY 的主端。
func (p *UnixPty) Master() *os.File {
	return p.master
}

// Slave 返回 PTY 的从端。
func (p *UnixPty) Slave() *os.File {
	return p.slave
}

// Control 使用主 PTY 的文件描述符运行给定的函数。
func (p *UnixPty) Control(fn func(fd uintptr)) error {
	conn, err := p.master.SyscallConn()
	if err != nil {
		return err //nolint:wrapcheck
	}

	return conn.Control(fn) //nolint:wrapcheck
}
