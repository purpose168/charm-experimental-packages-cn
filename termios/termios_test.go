//go:build !windows
// +build !windows

package termios

import (
	"os"
	"runtime"
	"testing"
)

// 这个测试主要是为了让 ./.github/workflows/termios.yml 能够为我们想要支持的平台构建测试，
// 并验证每个平台上的所有功能是否可用。
func TestTermios(t *testing.T) {
	if runtime.GOOS != "linux" {
		// 我们下面打开 pty 的方式是 Linux 方式。
		t.Skip()
	}
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() { _ = p.Close() })
	fd := int(p.Fd())
	w, err := GetWinsize(fd)
	if err != nil {
		t.Error(err)
	}
	if err := SetWinsize(fd, w); err != nil {
		t.Error(err)
	}

	term, err := GetTermios(fd)
	if err != nil {
		t.Error(err)
	}

	ispeed, ospeed := getSpeed(term)
	if err := SetTermios(
		fd,
		ispeed,
		ospeed,
		map[CC]uint8{
			ERASE:  1,
			CC(50): 12, // invalid, should be ignored
		},
		map[I]bool{
			IGNCR: true,
			IXOFF: false,
			I(50): true, // invalid, should be ignored
		},
		map[O]bool{
			OCRNL: true,
			ONLCR: false,
			O(50): true, // invalid, should be ignored
		},
		map[C]bool{
			CS7:    true,
			PARODD: false,
			C(50):  true, // invalid, should be ignored
		},
		map[L]bool{
			ECHO:  true,
			ECHOE: false,
			L(50): true, // invalid, should be ignored
		},
	); err != nil {
		t.Error(err)
	}

	term, err = GetTermios(fd)
	if err != nil {
		t.Error(err)
	}
	if v := term.Cc[allCcOpts[ERASE]]; v != 1 {
		t.Errorf("Cc.ERROR should be 1, was %d", v)
	}
	if v := term.Iflag & bit(allInputOpts[IGNCR]); v == 0 {
		t.Errorf("I.IGNCR should be true, was %d", v)
	}
	if v := term.Iflag & bit(allInputOpts[IXOFF]); v != 0 {
		t.Errorf("I.IGNCR should be false, was %d", v)
	}
	if v := term.Oflag & bit(allOutputOpts[OCRNL]); v == 0 {
		t.Errorf("O.OCRNL should be true, was %d", v)
	}
	if v := term.Oflag & bit(allOutputOpts[ONLCR]); v != 0 {
		t.Errorf("O.ONLCR should be false, was %d", v)
	}
	if v := term.Cflag & bit(allControlOpts[CS7]); v == 0 {
		t.Errorf("C.CS7 should be true, was %d", v)
	}
	if v := term.Cflag & bit(allControlOpts[PARODD]); v != 0 {
		t.Errorf("C.PARODD should be false, was %d", v)
	}
	if v := term.Lflag & bit(allLineOpts[ECHO]); v == 0 {
		t.Errorf("L.ECHO should be true, was %d", v)
	}
	if v := term.Lflag & bit(allLineOpts[ECHOE]); v != 0 {
		t.Errorf("L.ECHOE should be false, was %d", v)
	}
}
