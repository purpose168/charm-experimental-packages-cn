//nolint:gosec
package windows

import (
	"encoding/binary"

	"golang.org/x/sys/windows"
)

// FocusEventRecord 对应于 Windows 控制台 API 中的 FocusEventRecord 结构。
// https://docs.microsoft.com/en-us/windows/console/focus-event-record-str
type FocusEventRecord struct {
	// SetFocus 是保留的，不应使用。
	SetFocus bool
}

// MenuEventRecord 对应于 Windows 控制台 API 中的 MenuEventRecord 结构。
// https://docs.microsoft.com/en-us/windows/console/menu-event-record-str
type MenuEventRecord struct {
	CommandID uint32
}

// MouseEventRecord 对应于 Windows 控制台 API 中的 MouseEventRecord 结构。
// https://docs.microsoft.com/en-us/windows/console/mouse-event-record-str
type MouseEventRecord struct {
	// MousePosition 包含光标的位置，以控制台屏幕缓冲区的字符单元格坐标表示。
	MousePositon windows.Coord

	// ButtonState 保存鼠标按钮的状态。
	ButtonState uint32

	// ControlKeyState 保存控制键的状态。
	ControlKeyState uint32

	// EventFlags 指定鼠标事件的类型。
	EventFlags uint32
}

// WindowBufferSizeRecord 对应于 Windows 控制台 API 中的 WindowBufferSizeRecord 结构。
// https://docs.microsoft.com/en-us/windows/console/window-buffer-size-record-str
type WindowBufferSizeRecord struct {
	// Size 包含控制台屏幕缓冲区的大小，以字符单元格列和行表示。
	Size windows.Coord
}

// FocusEvent 将事件作为 FOCUS_EVENT_RECORD 返回。
func (ir InputRecord) FocusEvent() FocusEventRecord {
	return FocusEventRecord{SetFocus: ir.Event[0] > 0}
}

// KeyEvent 将事件作为 KEY_EVENT_RECORD 返回。
func (ir InputRecord) KeyEvent() KeyEventRecord {
	return KeyEventRecord{
		KeyDown:         binary.LittleEndian.Uint32(ir.Event[0:4]) > 0,
		RepeatCount:     binary.LittleEndian.Uint16(ir.Event[4:6]),
		VirtualKeyCode:  binary.LittleEndian.Uint16(ir.Event[6:8]),
		VirtualScanCode: binary.LittleEndian.Uint16(ir.Event[8:10]),
		Char:            rune(binary.LittleEndian.Uint16(ir.Event[10:12])),
		ControlKeyState: binary.LittleEndian.Uint32(ir.Event[12:16]),
	}
}

// MouseEvent 将事件作为 MOUSE_EVENT_RECORD 返回。
func (ir InputRecord) MouseEvent() MouseEventRecord {
	return MouseEventRecord{
		MousePositon: windows.Coord{
			X: int16(binary.LittleEndian.Uint16(ir.Event[0:2])),
			Y: int16(binary.LittleEndian.Uint16(ir.Event[2:4])),
		},
		ButtonState:     binary.LittleEndian.Uint32(ir.Event[4:8]),
		ControlKeyState: binary.LittleEndian.Uint32(ir.Event[8:12]),
		EventFlags:      binary.LittleEndian.Uint32(ir.Event[12:16]),
	}
}

// WindowBufferSizeEvent 将事件作为 WINDOW_BUFFER_SIZE_RECORD 返回。
func (ir InputRecord) WindowBufferSizeEvent() WindowBufferSizeRecord {
	return WindowBufferSizeRecord{
		Size: windows.Coord{
			X: int16(binary.LittleEndian.Uint16(ir.Event[0:2])),
			Y: int16(binary.LittleEndian.Uint16(ir.Event[2:4])),
		},
	}
}

// MenuEvent 将事件作为 MENU_EVENT_RECORD 返回。
func (ir InputRecord) MenuEvent() MenuEventRecord {
	return MenuEventRecord{
		CommandID: binary.LittleEndian.Uint32(ir.Event[0:4]),
	}
}
