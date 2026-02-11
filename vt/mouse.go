package vt

import (
	"io"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// MouseButton 表示鼠标消息期间按下的按钮。
type MouseButton = uv.MouseButton

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
//	8 = 第4个按钮（又名浏览器后退按钮）
//	9 = 第5个按钮（又名浏览器前进按钮）
//	10
//	11
//
// 其他按钮不受支持。
const (
	MouseNone       = uv.MouseNone
	MouseLeft       = uv.MouseLeft
	MouseMiddle     = uv.MouseMiddle
	MouseRight      = uv.MouseRight
	MouseWheelUp    = uv.MouseWheelUp
	MouseWheelDown  = uv.MouseWheelDown
	MouseWheelLeft  = uv.MouseWheelLeft
	MouseWheelRight = uv.MouseWheelRight
	MouseBackward   = uv.MouseBackward
	MouseForward    = uv.MouseForward
	MouseButton10   = uv.MouseButton10
	MouseButton11   = uv.MouseButton11
)

// Mouse 表示鼠标事件。
type Mouse = uv.MouseEvent

// MouseClick 表示鼠标点击事件。
type MouseClick = uv.MouseClickEvent

// MouseRelease 表示鼠标释放事件。
type MouseRelease = uv.MouseReleaseEvent

// MouseWheel 表示鼠标滚轮事件。
type MouseWheel = uv.MouseWheelEvent

// MouseMotion 表示鼠标移动事件。
type MouseMotion = uv.MouseMotionEvent

// SendMouse 向终端发送鼠标事件。这可以是任何类型的鼠标事件，
// 例如 [MouseClick]、[MouseRelease]、[MouseWheel] 或 [MouseMotion]。
func (e *Emulator) SendMouse(m Mouse) {
	// XXX: 支持 [Utf8ExtMouseMode]、[UrxvtExtMouseMode] 和
	// [SgrPixelExtMouseMode]。
	var (
		enc  ansi.Mode
		mode ansi.Mode
	)

	for _, m := range []ansi.DECMode{
		ansi.X10MouseMode,         // 按钮按下
		ansi.NormalMouseMode,      // 按钮按下/释放
		ansi.HighlightMouseMode,   // 按钮按下/释放/高亮
		ansi.ButtonEventMouseMode, // 按钮按下/释放/单元格移动
		ansi.AnyEventMouseMode,    // 按钮按下/释放/所有移动
	} {
		if e.isModeSet(m) {
			mode = m
		}
	}

	if mode == nil {
		return
	}

	for _, mm := range []ansi.DECMode{
		// ansi.Utf8ExtMouseMode,
		ansi.SgrExtMouseMode,
		// ansi.UrxvtExtMouseMode,
		// ansi.SgrPixelExtMouseMode,
	} {
		if e.isModeSet(mm) {
			enc = mm
		}
	}

	// 编码按钮
	mouse := m.Mouse()
	_, isMotion := m.(MouseMotion)
	_, isRelease := m.(MouseRelease)
	b := ansi.EncodeMouseButton(mouse.Button, isMotion,
		mouse.Mod.Contains(ModShift),
		mouse.Mod.Contains(ModAlt),
		mouse.Mod.Contains(ModCtrl))

	switch enc {
	// XXX: 支持 [ansi.HighlightMouseMode]。
	// XXX: 支持 [ansi.Utf8ExtMouseMode]、[ansi.UrxvtExtMouseMode] 和
	// [ansi.SgrPixelExtMouseMode]。
	case nil: // X10 鼠标编码
		_, _ = io.WriteString(e.pw, ansi.MouseX10(b, mouse.X, mouse.Y))
	case ansi.SgrExtMouseMode: // SGR 鼠标编码
		_, _ = io.WriteString(e.pw, ansi.MouseSgr(b, mouse.X, mouse.Y, isRelease))
	}
}
