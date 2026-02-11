// Package term 提供与终端和 TTY 设备交互的平台无关接口。
package term

// State 包含终端的平台特定状态。
type State struct {
	state
}

// IsTerminal 返回给定的文件描述符是否为终端。
func IsTerminal(fd uintptr) bool {
	return isTerminal(fd)
}

// MakeRaw 将连接到给定文件描述符的终端置于原始模式，并返回终端的先前状态以便可以恢复。
func MakeRaw(fd uintptr) (*State, error) {
	return makeRaw(fd)
}

// GetState 返回终端的当前状态，这在信号后恢复终端可能很有用。
func GetState(fd uintptr) (*State, error) {
	return getState(fd)
}

// SetState 设置终端的给定状态。
func SetState(fd uintptr, state *State) error {
	return setState(fd, state)
}

// Restore 将连接到给定文件描述符的终端恢复到先前的状态。
func Restore(fd uintptr, oldState *State) error {
	return restore(fd, oldState)
}

// GetSize 返回给定终端的可见尺寸。
//
// 这些尺寸不包括任何回滚缓冲区高度。
func GetSize(fd uintptr) (width, height int, err error) {
	return getSize(fd)
}

// ReadPassword 从终端读取一行输入而不进行本地回显。
// 这通常用于输入密码和其他敏感数据。返回的切片不包含 \n。
func ReadPassword(fd uintptr) ([]byte, error) {
	return readPassword(fd)
}
