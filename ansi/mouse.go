package ansi

import (
	"fmt"
)

// MouseButton 表示鼠标消息期间按下的按钮。
type MouseButton byte

// 鼠标事件按钮
//
// 基于 X11 鼠标按钮代码。
//
//	1 = 左键
//	2 = 中键（按下滚轮）
//	3 = 右键
//	4 = 向上滚动滚轮
//	5 = 向下滚动滚轮
//	6 = 向左推动滚轮
//	7 = 向右推动滚轮
//	8 = 第 4 个按钮（即浏览器后退按钮）
//	9 = 第 5 个按钮（即浏览器前进按钮）
//	10
//	11
//
// 其他按钮不受支持。
const (
	MouseNone MouseButton = iota
	MouseButton1
	MouseButton2
	MouseButton3
	MouseButton4
	MouseButton5
	MouseButton6
	MouseButton7
	MouseButton8
	MouseButton9
	MouseButton10
	MouseButton11

	MouseLeft       = MouseButton1
	MouseMiddle     = MouseButton2
	MouseRight      = MouseButton3
	MouseWheelUp    = MouseButton4
	MouseWheelDown  = MouseButton5
	MouseWheelLeft  = MouseButton6
	MouseWheelRight = MouseButton7
	MouseBackward   = MouseButton8
	MouseForward    = MouseButton9
	MouseRelease    = MouseNone
)

var mouseButtons = map[MouseButton]string{
	MouseNone:       "none",
	MouseLeft:       "left",
	MouseMiddle:     "middle",
	MouseRight:      "right",
	MouseWheelUp:    "wheelup",
	MouseWheelDown:  "wheeldown",
	MouseWheelLeft:  "wheelleft",
	MouseWheelRight: "wheelright",
	MouseBackward:   "backward",
	MouseForward:    "forward",
	MouseButton10:   "button10",
	MouseButton11:   "button11",
}

// String 返回鼠标按钮的字符串表示。
func (b MouseButton) String() string {
	return mouseButtons[b]
}

// EncodeMouseButton 返回表示鼠标按钮的字节。
// 按钮是以下最左侧值的位掩码：
//
//   - 前两位是按钮编号：
//     0 = 左键、滚轮向上或第 8 号按钮（即后退）
//     1 = 中键、滚轮向下或第 9 号按钮（即前进）
//     2 = 右键、滚轮向左或第 10 号按钮
//     3 = 释放事件、滚轮向右或第 11 号按钮
//
//   - 第三位表示是否按下了 Shift 键。
//
//   - 第四位表示是否按下了 Alt 键。
//
//   - 第五位表示是否按下了 Control 键。
//
//   - 第六位表示移动事件。与按钮编号 3（即释放事件）组合，表示拖动事件。
//
//   - 第七位表示滚轮事件。
//
//   - 第八位表示附加按钮。
//
// 如果按钮是 [MouseNone]，且 motion 为 false，则返回释放事件。
// 如果按钮未定义，此函数返回 0xff。
func EncodeMouseButton(b MouseButton, motion, shift, alt, ctrl bool) (m byte) {
	// 鼠标位偏移
	const (
		bitShift  = 0b0000_0100
		bitAlt    = 0b0000_1000
		bitCtrl   = 0b0001_0000
		bitMotion = 0b0010_0000
		bitWheel  = 0b0100_0000
		bitAdd    = 0b1000_0000 // 附加按钮 8-11

		bitsMask = 0b0000_0011
	)

	if b == MouseNone {
		m = bitsMask
	} else if b >= MouseLeft && b <= MouseRight {
		m = byte(b - MouseLeft)
	} else if b >= MouseWheelUp && b <= MouseWheelRight {
		m = byte(b - MouseWheelUp)
		m |= bitWheel
	} else if b >= MouseBackward && b <= MouseButton11 {
		m = byte(b - MouseBackward)
		m |= bitAdd
	} else {
		m = 0xff // 无效按钮
	}

	if shift {
		m |= bitShift
	}
	if alt {
		m |= bitAlt
	}
	if ctrl {
		m |= bitCtrl
	}
	if motion {
		m |= bitMotion
	}

	return m
}

// x10Offset 是 X10 鼠标事件的偏移量。
// 请参阅 https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#Mouse%20Tracking
const x10Offset = 32

// MouseX10 返回表示 X10 模式下鼠标事件的转义序列。
// 请注意，这需要终端支持 X10 鼠标模式。
//
//	CSI M Cb Cx Cy
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#Mouse%20Tracking
func MouseX10(b byte, x, y int) string {
	return "\x1b[M" + string(b+x10Offset) + string(byte(x)+x10Offset+1) + string(byte(y)+x10Offset+1)
}

// MouseSgr 返回表示 SGR 模式下鼠标事件的转义序列。
//
//	CSI < Cb ; Cx ; Cy M
//	CSI < Cb ; Cx ; Cy m (释放)
//
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#Mouse%20Tracking
func MouseSgr(b byte, x, y int, release bool) string {
	s := 'M'
	if release {
		s = 'm'
	}
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	return fmt.Sprintf("\x1b[<%d;%d;%d%c", b, x+1, y+1, s)
}
