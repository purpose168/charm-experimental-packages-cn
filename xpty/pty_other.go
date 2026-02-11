//go:build !linux && !darwin && !freebsd && !dragonfly && !netbsd && !openbsd && !solaris
// +build !linux,!darwin,!freebsd,!dragonfly,!netbsd,!openbsd,!solaris

package xpty

// setWinsize 设置窗口大小
// 此实现仅用于不支持的操作系统，返回 ErrUnsupported 错误
func (p *UnixPty) setWinsize(int, int, int, int) error {
	return ErrUnsupported
}

// size 获取终端大小
// 此实现仅用于不支持的操作系统，返回 ErrUnsupported 错误
func (*UnixPty) size() (int, int, error) {
	return 0, 0, ErrUnsupported
}
