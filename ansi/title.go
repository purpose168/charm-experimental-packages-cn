package ansi

// SetIconNameWindowTitle 返回用于设置图标名称和窗口标题的序列。
//
//	OSC 0 ; title ST
//	OSC 0 ; title BEL
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Operating-System-Commands
func SetIconNameWindowTitle(s string) string {
	return "\x1b]0;" + s + "\x07"
}

// SetIconName 返回用于设置图标名称的序列。
//
//	OSC 1 ; title ST
//	OSC 1 ; title BEL
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Operating-System-Commands
func SetIconName(s string) string {
	return "\x1b]1;" + s + "\x07"
}

// SetWindowTitle 返回用于设置窗口标题的序列。
//
//	OSC 2 ; title ST
//	OSC 2 ; title BEL
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Operating-System-Commands
func SetWindowTitle(s string) string {
	return "\x1b]2;" + s + "\x07"
}

// DECSWT 是用于设置窗口标题的序列。
//
// 这是 [SetWindowTitle]("1;<name>") 的别名。
// 请参阅：EK-VT520-RM 5–156 https://vt100.net/dec/ek-vt520-rm.pdf
func DECSWT(name string) string {
	return SetWindowTitle("1;" + name)
}

// DECSIN 是用于设置图标名称的序列。
//
// 这是 [SetWindowTitle]("L;<name>") 的别名。
// 请参阅：EK-VT520-RM 5–134 https://vt100.net/dec/ek-vt520-rm.pdf
func DECSIN(name string) string {
	return SetWindowTitle("L;" + name)
}
