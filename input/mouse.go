package input

import (
	"fmt"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// MouseButton 表示在鼠标消息期间按下的按钮。
type MouseButton = ansi.MouseButton

// 鼠标事件按钮
//
// 这基于 X11 鼠标按钮代码。
//
//	1 = 左键
//	2 = 中键（按下滚轮）
//	3 = 右键
//	4 = 向上滚动滚轮
//	5 = 向下滚动滚轮
//	6 = 向左推动滚轮
//	7 = 向右推动滚轮
//	8 = 第 4 个按钮（又称浏览器后退按钮）
//	9 = 第 5 个按钮（又称浏览器前进按钮）
//	10
//	11
//
// 不支持其他按钮。
const (
	MouseNone       = ansi.MouseNone
	MouseLeft       = ansi.MouseLeft
	MouseMiddle     = ansi.MouseMiddle
	MouseRight      = ansi.MouseRight
	MouseWheelUp    = ansi.MouseWheelUp
	MouseWheelDown  = ansi.MouseWheelDown
	MouseWheelLeft  = ansi.MouseWheelLeft
	MouseWheelRight = ansi.MouseWheelRight
	MouseBackward   = ansi.MouseBackward
	MouseForward    = ansi.MouseForward
	MouseButton10   = ansi.MouseButton10
	MouseButton11   = ansi.MouseButton11
)

// MouseEvent 表示鼠标消息。这是一个通用的鼠标消息，
// 可以表示任何类型的鼠标事件。
type MouseEvent interface {
	fmt.Stringer

	// Mouse 返回底层的鼠标事件。
	Mouse() Mouse
}

// Mouse 表示鼠标消息。使用 [MouseEvent] 表示所有鼠标消息。
//
// X 和 Y 坐标是从零开始的，(0,0) 是终端的左上角。
//
//	// 捕获所有鼠标事件
//	switch Event := Event.(type) {
//	case MouseEvent:
//	    m := Event.Mouse()
//	    fmt.Println("Mouse event:", m.X, m.Y, m)
//	}
//
//	// 只捕获鼠标点击事件
//	switch Event := Event.(type) {
//	case MouseClickEvent:
//	    fmt.Println("Mouse click event:", Event.X, Event.Y, Event)
//	}
type Mouse struct {
	X, Y   int
	Button MouseButton
	Mod    KeyMod
}

// String 返回鼠标消息的字符串表示形式。
func (m Mouse) String() (s string) {
	if m.Mod.Contains(ModCtrl) {
		s += "ctrl+"
	}
	if m.Mod.Contains(ModAlt) {
		s += "alt+"
	}
	if m.Mod.Contains(ModShift) {
		s += "shift+"
	}

	str := m.Button.String()
	if str == "" {
		s += "unknown"
	} else if str != "none" { // motion events don't have a button
		s += str
	}

	return s
}

// MouseClickEvent 表示鼠标按钮点击事件。
type MouseClickEvent Mouse

// String 返回鼠标点击事件的字符串表示形式。
func (e MouseClickEvent) String() string {
	return Mouse(e).String()
}

// Mouse 返回底层的鼠标事件。这是一个便捷方法和语法糖，用于满足 [MouseEvent] 接口，
// 并将鼠标事件转换为 [Mouse]。
func (e MouseClickEvent) Mouse() Mouse {
	return Mouse(e)
}

// MouseReleaseEvent 表示鼠标按钮释放事件。
type MouseReleaseEvent Mouse

// String 返回鼠标释放事件的字符串表示形式。
func (e MouseReleaseEvent) String() string {
	return Mouse(e).String()
}

// Mouse 返回底层的鼠标事件。这是一个便捷方法和语法糖，用于满足 [MouseEvent] 接口，
// 并将鼠标事件转换为 [Mouse]。
func (e MouseReleaseEvent) Mouse() Mouse {
	return Mouse(e)
}

// MouseWheelEvent 表示鼠标滚轮消息事件。
type MouseWheelEvent Mouse

// String 返回鼠标滚轮事件的字符串表示形式。
func (e MouseWheelEvent) String() string {
	return Mouse(e).String()
}

// Mouse 返回底层的鼠标事件。这是一个便捷方法和语法糖，用于满足 [MouseEvent] 接口，
// 并将鼠标事件转换为 [Mouse]。
func (e MouseWheelEvent) Mouse() Mouse {
	return Mouse(e)
}

// MouseMotionEvent 表示鼠标移动事件。
type MouseMotionEvent Mouse

// String 返回鼠标移动事件的字符串表示形式。
func (e MouseMotionEvent) String() string {
	m := Mouse(e)
	if m.Button != 0 {
		return m.String() + "+motion"
	}
	return m.String() + "motion"
}

// Mouse 返回底层的鼠标事件。这是一个便捷方法和语法糖，用于满足 [MouseEvent] 接口，
// 并将鼠标事件转换为 [Mouse]。
func (e MouseMotionEvent) Mouse() Mouse {
	return Mouse(e)
}

// 解析 SGR 编码的鼠标事件；SGR 扩展鼠标事件。SGR 鼠标事件看起来像：
//
//	ESC [ < Cb ; Cx ; Cy (M 或 m)
//
// 其中：
//
//	Cb 是编码的按钮代码
//	Cx 是鼠标的 x 坐标
//	Cy 是鼠标的 y 坐标
//	M 表示按钮按下，m 表示按钮释放
//
// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseSGRMouseEvent(cmd ansi.Cmd, params ansi.Params) Event {
	x, _, ok := params.Param(1, 1)
	if !ok {
		x = 1
	}
	y, _, ok := params.Param(2, 1)
	if !ok {
		y = 1
	}
	release := cmd.Final() == 'm'
	b, _, _ := params.Param(0, 0)
	mod, btn, _, isMotion := parseMouseButton(b)

	// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
	x--
	y--

	m := Mouse{X: x, Y: y, Button: btn, Mod: mod}

	// Wheel buttons don't have release events
	// Motion can be reported as a release event in some terminals (Windows Terminal)
	if isWheel(m.Button) {
		return MouseWheelEvent(m)
	} else if !isMotion && release {
		return MouseReleaseEvent(m)
	} else if isMotion {
		return MouseMotionEvent(m)
	}
	return MouseClickEvent(m)
}

const x10MouseByteOffset = 32

// 解析 X10 编码的鼠标事件；最简单的一种。顺便说一下，X10 的最后一个版本是 1986 年 12 月。
// 原始的 X10 鼠标协议将 Cx 和 Cy 坐标限制为 223 (=255-032)。
//
// X10 鼠标事件看起来像：
//
//	ESC [M Cb Cx Cy
//
// 请参阅：http://www.xfree86.org/current/ctlseqs.html#Mouse%20Tracking
func parseX10MouseEvent(buf []byte) Event {
	v := buf[3:6]
	b := int(v[0])
	if b >= x10MouseByteOffset {
			// XXX: b < 32 应该是不可能的，但我们保持防御性。
			b -= x10MouseByteOffset
		}

	mod, btn, isRelease, isMotion := parseMouseButton(b)

	// (1,1) is the upper left. We subtract 1 to normalize it to (0,0).
	x := int(v[1]) - x10MouseByteOffset - 1
	y := int(v[2]) - x10MouseByteOffset - 1

	m := Mouse{X: x, Y: y, Button: btn, Mod: mod}
	if isWheel(m.Button) {
		return MouseWheelEvent(m)
	} else if isMotion {
		return MouseMotionEvent(m)
	} else if isRelease {
		return MouseReleaseEvent(m)
	}
	return MouseClickEvent(m)
}

// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Extended-coordinates
func parseMouseButton(b int) (mod KeyMod, btn MouseButton, isRelease bool, isMotion bool) {
	// 鼠标位偏移
	const (
		bitShift  = 0b0000_0100
		bitAlt    = 0b0000_1000
		bitCtrl   = 0b0001_0000
		bitMotion = 0b0010_0000
		bitWheel  = 0b0100_0000
		bitAdd    = 0b1000_0000 // 额外按钮 8-11

		bitsMask = 0b0000_0011
	)

	// 修饰键
	if b&bitAlt != 0 {
		mod |= ModAlt
	}
	if b&bitCtrl != 0 {
		mod |= ModCtrl
	}
	if b&bitShift != 0 {
		mod |= ModShift
	}

	if b&bitAdd != 0 {
		btn = MouseBackward + MouseButton(b&bitsMask)
	} else if b&bitWheel != 0 {
		btn = MouseWheelUp + MouseButton(b&bitsMask)
	} else {
		btn = MouseLeft + MouseButton(b&bitsMask)
		// X10 将按钮释放报告为 0b0000_0011 (3)
		if b&bitsMask == bitsMask {
			btn = MouseNone
			isRelease = true
		}
	}

	// 滚轮事件不会报告移动位。
	if b&bitMotion != 0 && !isWheel(btn) {
		isMotion = true
	}

	return mod, btn, isRelease, isMotion
}

// isWheel 如果鼠标事件是滚轮事件，则返回 true。
func isWheel(btn MouseButton) bool {
	return btn >= MouseWheelUp && btn <= MouseWheelRight
}
