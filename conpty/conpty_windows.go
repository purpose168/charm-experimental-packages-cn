//go:build windows
// +build windows

package conpty

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ConPty 表示 Windows 控制台伪终端
// https://learn.microsoft.com/zh-cn/windows/console/creating-a-pseudoconsole-session#preparing-the-communication-channels
type ConPty struct {
	// hpc 是伪控制台的句柄
	hpc *windows.Handle
	// inPipeWrite, inPipeRead 是输入管道的写入和读取句柄
	inPipeWrite, inPipeRead windows.Handle
	// outPipeWrite, outPipeRead 是输出管道的写入和读取句柄
	outPipeWrite, outPipeRead windows.Handle
	// attrList 是进程线程属性列表容器
	attrList *windows.ProcThreadAttributeListContainer
	// size 是伪控制台的大小
	size windows.Coord
	// closeOnce 确保 Close 方法只被调用一次
	closeOnce sync.Once
}

var (
	_ io.Writer = &ConPty{} // 确保 ConPty 实现了 io.Writer 接口
	_ io.Reader = &ConPty{} // 确保 ConPty 实现了 io.Reader 接口
)

// CreatePipes 是一个辅助函数，用于创建连接的输入和输出管道
func CreatePipes() (inPipeRead, inPipeWrite, outPipeRead, outPipeWrite uintptr, err error) {
	var inPipeReadHandle, inPipeWriteHandle windows.Handle
	var outPipeReadHandle, outPipeWriteHandle windows.Handle
	pSec := &windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(zeroSec)), InheritHandle: 1}

	if err := windows.CreatePipe(&inPipeReadHandle, &inPipeWriteHandle, pSec, 0); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("创建伪控制台输入管道失败: %w", err)
	}

	if err := windows.CreatePipe(&outPipeReadHandle, &outPipeWriteHandle, pSec, 0); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("创建伪控制台输出管道失败: %w", err)
	}

	return uintptr(inPipeReadHandle), uintptr(inPipeWriteHandle),
		uintptr(outPipeReadHandle), uintptr(outPipeWriteHandle),
		nil
}

// New 创建一个新的 ConPty 设备
// 接受自定义宽度、高度和将传递给 windows.CreatePseudoConsole 的标志
func New(w int, h int, flags int) (*ConPty, error) {
	inPipeRead, inPipeWrite, outPipeRead, outPipeWrite, err := CreatePipes()
	if err != nil {
		return nil, fmt.Errorf("为伪控制台创建管道失败: %w", err)
	}

	c, err := NewWithPipes(inPipeRead, inPipeWrite, outPipeRead, outPipeWrite, w, h, flags)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewWithPipes 使用提供的管道句柄创建一个新的 ConPty 设备
// 当你想要使用现有的管道时非常有用，例如当你想要将 ConPty 与已经创建的进程一起使用，
// 或者当你想要使用特定的管道集进行输入和输出时。
//
// 管道的 PTY 从端（输入读取和输出写入）在 ConPty 创建后可以关闭，
// 因为 ConPty 会接管这些句柄并为将要生成的新进程复制它们。
// 管道的 PTY 主端将用于与伪控制台通信。
func NewWithPipes(inPipeRead, inPipeWrite, outPipeRead, outPipeWrite uintptr, w int, h int, flags int) (c *ConPty, err error) {
	if w <= 0 {
		w = DefaultWidth
	}
	if h <= 0 {
		h = DefaultHeight
	}

	c = &ConPty{
		hpc: new(windows.Handle),
		size: windows.Coord{
			X: int16(w), Y: int16(h),
		},
		inPipeWrite:  windows.Handle(inPipeWrite),
		inPipeRead:   windows.Handle(inPipeRead),
		outPipeWrite: windows.Handle(outPipeWrite),
		outPipeRead:  windows.Handle(outPipeRead),
	}

	if err := windows.CreatePseudoConsole(c.size, windows.Handle(inPipeRead), windows.Handle(outPipeWrite), uint32(flags), c.hpc); err != nil {
		return nil, fmt.Errorf("创建伪控制台失败: %w", err)
	}

	// 分配一个足够大的属性列表来执行我们关心的操作
	// 1. 伪控制台设置
	c.attrList, err = windows.NewProcThreadAttributeList(1)
	if err != nil {
		return nil, err
	}

	if err := c.attrList.Update(
		windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		unsafe.Pointer(*c.hpc),
		unsafe.Sizeof(*c.hpc),
	); err != nil {
		return nil, fmt.Errorf("更新伪控制台的进程线程属性失败: %w", err)
	}

	return c, err
}

// Fd 返回 ConPty 句柄
func (p *ConPty) Fd() uintptr {
	return uintptr(*p.hpc)
}

// Close 关闭 ConPty 设备
func (p *ConPty) Close() error {
	var err error
	p.closeOnce.Do(func() {
		// 确保我们关闭了管道的 PTY 端
		_ = windows.CloseHandle(p.inPipeRead)
		_ = windows.CloseHandle(p.outPipeWrite)
		if p.attrList != nil {
			p.attrList.Delete()
		}
		windows.ClosePseudoConsole(*p.hpc)
		err = errors.Join(
			windows.CloseHandle(p.inPipeWrite),
			windows.CloseHandle(p.outPipeRead),
		)
	})
	return err
}

// InPipeReadFd 返回 ConPty 输入管道的读取文件描述符句柄
func (p *ConPty) InPipeReadFd() uintptr {
	return uintptr(p.inPipeRead)
}

// InPipeWriteFd 返回 ConPty 输入管道的写入文件描述符句柄
func (p *ConPty) InPipeWriteFd() uintptr {
	return uintptr(p.inPipeWrite)
}

// OutPipeReadFd 返回 ConPty 输出管道的读取文件描述符句柄
func (p *ConPty) OutPipeReadFd() uintptr {
	return uintptr(p.outPipeRead)
}

// OutPipeWriteFd 返回 ConPty 输出管道的写入文件描述符句柄
func (p *ConPty) OutPipeWriteFd() uintptr {
	return uintptr(p.outPipeWrite)
}

// Write 安全地向 ConPty 的主端写入字节
func (c *ConPty) Write(p []byte) (n int, err error) {
	var l uint32
	err = windows.WriteFile(c.inPipeWrite, p, &l, nil)
	return int(l), err
}

// Read 安全地从 ConPty 的主端读取字节
func (c *ConPty) Read(p []byte) (n int, err error) {
	var l uint32
	err = windows.ReadFile(c.outPipeRead, p, &l, nil)
	return int(l), err
}

// Resize 调整伪控制台的大小
func (c *ConPty) Resize(w int, h int) error {
	size := windows.Coord{X: int16(w), Y: int16(h)}
	if err := windows.ResizePseudoConsole(*c.hpc, size); err != nil {
		return fmt.Errorf("调整伪控制台大小失败: %w", err)
	}
	c.size = size
	return nil
}

// Size 返回当前伪控制台的大小
func (c *ConPty) Size() (w int, h int, err error) {
	w = int(c.size.X)
	h = int(c.size.Y)
	return w, h, err
}

var (
	zeroAttr syscall.ProcAttr
	zeroSec  windows.SecurityAttributes
)

// Spawn 在伪控制台中启动一个新进程
func (c *ConPty) Spawn(name string, args []string, attr *syscall.ProcAttr) (pid int, handle uintptr, err error) {
	if attr == nil {
		attr = &zeroAttr
	}

	argv0, err := lookExtensions(name, attr.Dir)
	if err != nil {
		return 0, 0, err
	}
	if len(attr.Dir) != 0 {
		// Windows CreateProcess 会在当前目录中查找 argv0，
		// 并且只有在新进程启动后，才会执行 Chdir(attr.Dir)。
		// 我们通过使 argv0 成为绝对路径来调整这种差异
		var err error
		argv0, err = joinExeDirAndFName(attr.Dir, argv0)
		if err != nil {
			return 0, 0, err
		}
	}

	argv0p, err := windows.UTF16PtrFromString(argv0)
	if err != nil {
		return 0, 0, err
	}

	var cmdline string
	if attr.Sys != nil && attr.Sys.CmdLine != "" {
		cmdline = attr.Sys.CmdLine
	} else {
		cmdline = windows.ComposeCommandLine(args)
	}
	argvp, err := windows.UTF16PtrFromString(cmdline)
	if err != nil {
		return 0, 0, err
	}

	var dirp *uint16
	if len(attr.Dir) != 0 {
		dirp, err = windows.UTF16PtrFromString(attr.Dir)
		if err != nil {
			return 0, 0, err
		}
	}

	if attr.Env == nil {
		attr.Env, err = execEnvDefault(attr.Sys)
		if err != nil {
			return 0, 0, err
		}
	}

	siEx := new(windows.StartupInfoEx)
	siEx.Flags = windows.STARTF_USESTDHANDLES

	pi := new(windows.ProcessInformation)

	// 需要 EXTENDED_STARTUPINFO_PRESENT，因为我们正在使用属性列表字段
	flags := uint32(windows.CREATE_UNICODE_ENVIRONMENT) | windows.EXTENDED_STARTUPINFO_PRESENT
	if attr.Sys != nil && attr.Sys.CreationFlags != 0 {
		flags |= attr.Sys.CreationFlags
	}

	pSec := &windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(zeroSec)), InheritHandle: 1}
	if attr.Sys != nil && attr.Sys.ProcessAttributes != nil {
		pSec = &windows.SecurityAttributes{
			Length:        attr.Sys.ProcessAttributes.Length,
			InheritHandle: attr.Sys.ProcessAttributes.InheritHandle,
		}
	}
	tSec := &windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(zeroSec)), InheritHandle: 1}
	if attr.Sys != nil && attr.Sys.ThreadAttributes != nil {
		tSec = &windows.SecurityAttributes{
			Length:        attr.Sys.ThreadAttributes.Length,
			InheritHandle: attr.Sys.ThreadAttributes.InheritHandle,
		}
	}

	siEx.ProcThreadAttributeList = c.attrList.List() //nolint:govet // unusedwrite: ProcThreadAttributeList 将在系统调用中读取
	siEx.Cb = uint32(unsafe.Sizeof(*siEx))
	if attr.Sys != nil && attr.Sys.Token != 0 {
		err = windows.CreateProcessAsUser(
			windows.Token(attr.Sys.Token),
			argv0p,
			argvp,
			pSec,
			tSec,
			false,
			flags,
			createEnvBlock(addCriticalEnv(dedupEnvCase(true, attr.Env))),
			dirp,
			&siEx.StartupInfo,
			pi,
		)
	} else {
		err = windows.CreateProcess(
			argv0p,
			argvp,
			pSec,
			tSec,
			false,
			flags,
			createEnvBlock(addCriticalEnv(dedupEnvCase(true, attr.Env))),
			dirp,
			&siEx.StartupInfo,
			pi,
		)
	}
	if err != nil {
		return 0, 0, fmt.Errorf("创建进程失败: %w", err)
	}

	defer windows.CloseHandle(pi.Thread)

	return int(pi.ProcessId), uintptr(pi.Process), nil
}
