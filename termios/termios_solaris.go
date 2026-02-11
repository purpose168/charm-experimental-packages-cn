//go:build solaris
// +build solaris

package termios

import "golang.org/x/sys/unix"

// 参见 https://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libc/port/gen/isatty.c
// 参见 https://github.com/omniti-labs/illumos-omnios/blob/master/usr/src/uts/common/sys/termios.h
const (
	ioctlSets       = unix.TCSETA
	ioctlGets       = unix.TCGETA
	ioctlSetWinSize = (int('T') << 8) | 103
	ioctlGetWinSize = (int('T') << 8) | 104
)

func setSpeed(*unix.Termios, uint32, uint32) {
	// TODO: 支持在 Solaris 上设置速度？
	// 参见 cfgetospeed(3C) 和 cfsetospeed(3C)
	// 参见 cfgetispeed(3C) 和 cfsetispeed(3C)
	// https://github.com/omniti-labs/illumos-omnios/blob/master/usr/src/uts/common/sys/termios.h#L103
}

func getSpeed(*unix.Termios) (uint32, uint32) {
	return 0, 0
}
