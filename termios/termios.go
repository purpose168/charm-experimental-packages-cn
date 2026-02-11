//go:build darwin || netbsd || freebsd || openbsd || linux || dragonfly || solaris
// +build darwin netbsd freebsd openbsd linux dragonfly solaris

// Package termios 为 Unix 和类 Unix 系统提供获取和设置 Termios 设置的统一接口。
package termios

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// SetWinsize 从 Winsize 设置文件描述符的窗口大小。
func SetWinsize(fd int, w *unix.Winsize) error {
	return unix.IoctlSetWinsize(fd, ioctlSetWinSize, w) //nolint:wrapcheck
}

// GetWinsize 获取文件描述符的窗口大小。
func GetWinsize(fd int) (*unix.Winsize, error) {
	return unix.IoctlGetWinsize(fd, ioctlGetWinSize) //nolint:wrapcheck
}

// GetTermios 获取给定文件描述符的 termios 设置。
func GetTermios(fd int) (*unix.Termios, error) {
	return unix.IoctlGetTermios(fd, ioctlGets) //nolint:wrapcheck
}

// SetTermios 在给定文件描述符的当前 termios 设置上设置给定的 termios。
func SetTermios(
	fd int,
	ispeed uint32,
	ospeed uint32,
	cc map[CC]uint8,
	iflag map[I]bool,
	oflag map[O]bool,
	cflag map[C]bool,
	lflag map[L]bool,
) error {
	term, err := unix.IoctlGetTermios(fd, ioctlGets)
	if err != nil {
		return err //nolint:wrapcheck
	}
	setSpeed(term, ispeed, ospeed)

	for key, value := range cc {
		call, ok := allCcOpts[key]
		if !ok {
			continue
		}
		term.Cc[call] = value
	}

	for key, value := range iflag {
		mask, ok := allInputOpts[key]
		if ok {
			if value {
				term.Iflag |= bit(mask)
			} else {
				term.Iflag &= ^bit(mask)
			}
		}
	}
	for key, value := range oflag {
		mask, ok := allOutputOpts[key]
		if ok {
			if value {
				term.Oflag |= bit(mask)
			} else {
				term.Oflag &= ^bit(mask)
			}
		}
	}
	for key, value := range cflag {
		mask, ok := allControlOpts[key]
		if ok {
			if value {
				term.Cflag |= bit(mask)
			} else {
				term.Cflag &= ^bit(mask)
			}
		}
	}
	for key, value := range lflag {
		mask, ok := allLineOpts[key]
		if ok {
			if value {
				term.Lflag |= bit(mask)
			} else {
				term.Lflag &= ^bit(mask)
			}
		}
	}
	return unix.IoctlSetTermios(fd, ioctlSets, term) //nolint:wrapcheck
}

// CC 是 termios 的 cc 字段。
//
// 它存储与终端 I/O 相关的特殊字符数组。
type CC uint8

// CC 可能的值。
const (
	INTR CC = iota // 中断字符
	QUIT          // 退出字符
	ERASE         // 擦除字符
	KILL          // 杀死字符
	EOF           // 文件结束字符
	EOL           // 行结束字符
	EOL2          // 备用行结束字符
	START         // 开始字符
	STOP          // 停止字符
	SUSP          // 挂起字符
	WERASE        // 字擦除字符
	RPRNT         // 重印字符
	LNEXT         // 字面下一个字符
	DISCARD       // 丢弃字符
	STATUS        // 状态字符
	SWTCH         // 切换字符
	DSUSP         // 延迟挂起字符
	FLUSH         // 刷新字符
)

// https://www.man7.org/linux/man-pages/man3/termios.3.html
var allCcOpts = map[CC]int{
	INTR:    syscall.VINTR,
	QUIT:    syscall.VQUIT,
	ERASE:   syscall.VERASE,
	KILL:    syscall.VQUIT,
	EOF:     syscall.VEOF,
	EOL:     syscall.VEOL,
	EOL2:    syscall.VEOL2,
	START:   syscall.VSTART,
	STOP:    syscall.VSTOP,
	SUSP:    syscall.VSUSP,
	WERASE:  syscall.VWERASE,
	RPRNT:   syscall.VREPRINT,
	LNEXT:   syscall.VLNEXT,
	DISCARD: syscall.VDISCARD,

	// XXX: these syscalls don't exist for any OS
	// FLUSH:  syscall.VFLUSH,
}

// I 代表输入控制。
type I uint8

// 输入可能的值。
const (
	IGNPAR I = iota // 忽略奇偶校验错误
	PARMRK          // 标记奇偶校验错误
	INPCK           // 启用奇偶校验
	ISTRIP          // 剥离第八位
	INLCR           // 将 NL 转换为 CR
	IGNCR           // 忽略 CR
	ICRNL           // 将 CR 转换为 NL
	IXON            // 启用输出流控制
	IXANY           // 允许任何字符重启输出
	IXOFF           // 启用输入流控制
	IMAXBEL         // 当输入队列满时响铃
	IUCLC           // 将大写字母转换为小写字母
)

var allInputOpts = map[I]uint32{
	IGNPAR:  syscall.IGNPAR,
	PARMRK:  syscall.PARMRK,
	INPCK:   syscall.INPCK,
	ISTRIP:  syscall.ISTRIP,
	INLCR:   syscall.INLCR,
	IGNCR:   syscall.IGNCR,
	ICRNL:   syscall.ICRNL,
	IXON:    syscall.IXON,
	IXANY:   syscall.IXANY,
	IXOFF:   syscall.IXOFF,
	IMAXBEL: syscall.IMAXBEL,
}

// O 代表输出控制。
type O uint8

// 输出可能的值。
const (
	OPOST O = iota // 启用输出处理
	ONLCR          // 将 NL 转换为 CR-NL
	OCRNL          // 将 CR 转换为 NL
	ONOCR          // 在第 0 列不输出 CR
	ONLRET         // NL 执行 CR 功能
	OLCUC          // 将小写字母转换为大写字母
)

var allOutputOpts = map[O]uint32{
	OPOST:  syscall.OPOST,
	ONLCR:  syscall.ONLCR,
	OCRNL:  syscall.OCRNL,
	ONOCR:  syscall.ONOCR,
	ONLRET: syscall.ONLRET,
}

// C 代表控制。
type C uint8

// 控制可能的值。
const (
	CS7 C = iota   // 7 位字符大小
	CS8            // 8 位字符大小
	PARENB         // 启用奇偶校验
	PARODD         // 使用奇校验而不是偶校验
)

var allControlOpts = map[C]uint32{
	CS7:    syscall.CS7,
	CS8:    syscall.CS8,
	PARENB: syscall.PARENB,
	PARODD: syscall.PARODD,
}

// L 代表行控制。
type L uint8

// 行可能的值。
const (
	ISIG L = iota // 启用信号
	ICANON        // 启用规范模式
	ECHO          // 启用回显
	ECHOE         // 擦除时回显
	ECHOK         // 杀死时回显
	ECHONL        // 回显 NL
	NOFLSH        // 不刷新
	TOSTOP        // 停止后台进程
	IEXTEN        // 启用扩展输入处理
	ECHOCTL       // 回显控制字符
	ECHOKE        // 杀死时回显
	PENDIN        // 挂起输入
	IUTF8         // 启用 UTF-8 输入
	XCASE         // 大小写映射
)

var allLineOpts = map[L]uint32{
	ISIG:    syscall.ISIG,
	ICANON:  syscall.ICANON,
	ECHO:    syscall.ECHO,
	ECHOE:   syscall.ECHOE,
	ECHOK:   syscall.ECHOK,
	ECHONL:  syscall.ECHONL,
	NOFLSH:  syscall.NOFLSH,
	TOSTOP:  syscall.TOSTOP,
	IEXTEN:  syscall.IEXTEN,
	ECHOCTL: syscall.ECHOCTL,
	ECHOKE:  syscall.ECHOKE,
	PENDIN:  syscall.PENDIN,
}
