//go:build windows
// +build windows

package conpty

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// 以下是用于处理 Windows CreateProcess 系列函数的辅助函数
// 这些大部分是从 Go 标准库中复制的相同工具函数

// lookExtensions 在指定目录中查找带有扩展名的可执行文件
func lookExtensions(path, dir string) (string, error) {
	if filepath.Base(path) == path {
		path = filepath.Join(".", path)
	}

	if dir == "" {
		return exec.LookPath(path)
	}

	if filepath.VolumeName(path) != "" {
		return exec.LookPath(path)
	}

	if len(path) > 1 && os.IsPathSeparator(path[0]) {
		return exec.LookPath(path)
	}

	dirandpath := filepath.Join(dir, path)

	// 我们假设 LookPath 只会添加文件扩展名
	lp, err := exec.LookPath(dirandpath)
	if err != nil {
		return "", err
	}

	ext := strings.TrimPrefix(lp, dirandpath)

	return path + ext, nil
}

// execEnvDefault 获取默认的环境变量
func execEnvDefault(sys *syscall.SysProcAttr) (env []string, err error) {
	if sys == nil || sys.Token == 0 {
		return syscall.Environ(), nil
	}

	var block *uint16
	err = windows.CreateEnvironmentBlock(&block, windows.Token(sys.Token), false)
	if err != nil {
		return nil, err
	}

	defer windows.DestroyEnvironmentBlock(block)
	blockp := uintptr(unsafe.Pointer(block))

	for {
		// 查找 NUL 终止符
		end := unsafe.Pointer(blockp)
		for *(*uint16)(end) != 0 {
			end = unsafe.Pointer(uintptr(end) + 2)
		}

		n := (uintptr(end) - uintptr(unsafe.Pointer(blockp))) / 2
		if n == 0 {
			// 环境块以空字符串结束
			break
		}

		entry := (*[(1 << 30) - 1]uint16)(unsafe.Pointer(blockp))[:n:n]
		env = append(env, string(utf16.Decode(entry)))
		blockp += 2 * (uintptr(len(entry)) + 1)
	}
	return env, err
}

// isSlash 检查字符是否为斜杠
func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

// normalizeDir 规范化目录路径
func normalizeDir(dir string) (name string, err error) {
	ndir, err := syscall.FullPath(dir)
	if err != nil {
		return "", err
	}
	if len(ndir) > 2 && isSlash(ndir[0]) && isSlash(ndir[1]) {
		// 目录不能有 \server\share\path 形式
		return "", syscall.EINVAL
	}
	return ndir, nil
}

// volToUpper 将卷号转换为大写
func volToUpper(ch int) int {
	if 'a' <= ch && ch <= 'z' {
		ch += 'A' - 'a'
	}
	return ch
}

// joinExeDirAndFName 连接可执行文件目录和文件名
func joinExeDirAndFName(dir, p string) (name string, err error) {
	if len(p) == 0 {
		return "", syscall.EINVAL
	}
	if len(p) > 2 && isSlash(p[0]) && isSlash(p[1]) {
		// \server\share\path 形式
		return p, nil
	}
	if len(p) > 1 && p[1] == ':' {
		// 有驱动器号
		if len(p) == 2 {
			return "", syscall.EINVAL
		}
		if isSlash(p[2]) {
			return p, nil
		} else {
			d, err := normalizeDir(dir)
			if err != nil {
				return "", err
			}
			if volToUpper(int(p[0])) == volToUpper(int(d[0])) {
				return syscall.FullPath(d + "\\" + p[2:])
			} else {
				return syscall.FullPath(p)
			}
		}
	} else {
		// 没有驱动器号
		d, err := normalizeDir(dir)
		if err != nil {
			return "", err
		}
		if isSlash(p[0]) {
			return windows.FullPath(d[:2] + p)
		} else {
			return windows.FullPath(d + "\\" + p)
		}
	}
}

// createEnvBlock 将环境字符串数组转换为
// CreateProcess 所需的表示形式：一系列 NUL 终止的字符串，后跟一个 nil
// 最后字节是两个 UCS-2 NUL，或四个 NUL 字节
func createEnvBlock(envv []string) *uint16 {
	if len(envv) == 0 {
		return &utf16.Encode([]rune("\x00\x00"))[0]
	}
	length := 0
	for _, s := range envv {
		length += len(s) + 1
	}
	length++

	b := make([]byte, length)
	i := 0
	for _, s := range envv {
		l := len(s)
		copy(b[i:i+l], []byte(s))
		copy(b[i+l:i+l+1], []byte{0})
		i = i + l + 1
	}
	copy(b[i:i+1], []byte{0})

	return &utf16.Encode([]rune(string(b)))[0]
}

// dedupEnvCase 是带有测试用例选项的 dedupEnv
// 如果 caseInsensitive 为 true，则忽略键的大小写
func dedupEnvCase(caseInsensitive bool, env []string) []string {
	out := make([]string, 0, len(env))
	saw := make(map[string]int, len(env)) // key => 索引到 out
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			out = append(out, kv)
			continue
		}
		k := kv[:eq]
		if caseInsensitive {
			k = strings.ToLower(k)
		}
		if dupIdx, isDup := saw[k]; isDup {
			out[dupIdx] = kv
			continue
		}
		saw[k] = len(out)
		out = append(out, kv)
	}
	return out
}

// addCriticalEnv 添加操作系统所需的任何关键环境变量
// （或者至少几乎总是需要的）
// 目前这仅用于 Windows
func addCriticalEnv(env []string) []string {
	for _, kv := range env {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			continue
		}
		k := kv[:eq]
		if strings.EqualFold(k, "SYSTEMROOT") {
			// 我们已经有了它
			return env
		}
	}
	return append(env, "SYSTEMROOT="+os.Getenv("SYSTEMROOT"))
}
