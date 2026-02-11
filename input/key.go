package input

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

const (
	// KeyExtended 是一个特殊的键码，用于表示一个键事件包含多个符文。
	KeyExtended = unicode.MaxRune + 1
)

// 特殊键符号。
const (

	// 特殊键。

	KeyUp rune = KeyExtended + iota + 1
	KeyDown
	KeyRight
	KeyLeft
	KeyBegin
	KeyFind
	KeyInsert
	KeyDelete
	KeySelect
	KeyPgUp
	KeyPgDown
	KeyHome
	KeyEnd

	// 键盘按键。

	KeyKpEnter
	KeyKpEqual
	KeyKpMultiply
	KeyKpPlus
	KeyKpComma
	KeyKpMinus
	KeyKpDecimal
	KeyKpDivide
	KeyKp0
	KeyKp1
	KeyKp2
	KeyKp3
	KeyKp4
	KeyKp5
	KeyKp6
	KeyKp7
	KeyKp8
	KeyKp9

	//nolint:godox
	// 以下是 Kitty 键盘协议中定义的键。
	// TODO: 调查这些键的名称。

	KeyKpSep
	KeyKpUp
	KeyKpDown
	KeyKpLeft
	KeyKpRight
	KeyKpPgUp
	KeyKpPgDown
	KeyKpHome
	KeyKpEnd
	KeyKpInsert
	KeyKpDelete
	KeyKpBegin

	// 功能键。

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyF26
	KeyF27
	KeyF28
	KeyF29
	KeyF30
	KeyF31
	KeyF32
	KeyF33
	KeyF34
	KeyF35
	KeyF36
	KeyF37
	KeyF38
	KeyF39
	KeyF40
	KeyF41
	KeyF42
	KeyF43
	KeyF44
	KeyF45
	KeyF46
	KeyF47
	KeyF48
	KeyF49
	KeyF50
	KeyF51
	KeyF52
	KeyF53
	KeyF54
	KeyF55
	KeyF56
	KeyF57
	KeyF58
	KeyF59
	KeyF60
	KeyF61
	KeyF62
	KeyF63

	//nolint:godox
	// 以下是 Kitty 键盘协议中定义的键。
	// TODO: 调查这些键的名称。

	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyMenu

	KeyMediaPlay
	KeyMediaPause
	KeyMediaPlayPause
	KeyMediaReverse
	KeyMediaStop
	KeyMediaFastForward
	KeyMediaRewind
	KeyMediaNext
	KeyMediaPrev
	KeyMediaRecord

	KeyLowerVol
	KeyRaiseVol
	KeyMute

	KeyLeftShift
	KeyLeftAlt
	KeyLeftCtrl
	KeyLeftSuper
	KeyLeftHyper
	KeyLeftMeta
	KeyRightShift
	KeyRightAlt
	KeyRightCtrl
	KeyRightSuper
	KeyRightHyper
	KeyRightMeta
	KeyIsoLevel3Shift
	KeyIsoLevel5Shift

	// C0 中的特殊名称。

	KeyBackspace = rune(ansi.DEL)
	KeyTab       = rune(ansi.HT)
	KeyEnter     = rune(ansi.CR)
	KeyReturn    = KeyEnter
	KeyEscape    = rune(ansi.ESC)
	KeyEsc       = KeyEscape

	// G0 中的特殊名称。

	KeySpace = rune(ansi.SP)
)

// KeyPressEvent 表示按键按下事件。
type KeyPressEvent Key

// String 实现了 [fmt.Stringer] 接口，对于匹配键事件非常有用。
// 有关其返回内容的详细信息，请参见 [Key.String]。
func (k KeyPressEvent) String() string {
	return Key(k).String()
}

// Keystroke 返回 [Key] 的按键表示形式。虽然在类型安全性上不如查看各个字段，
// 但在匹配按键时，使用此方法通常会更方便、更易读。
//
// 注意，修饰键始终按以下顺序打印：
//   - ctrl
//   - alt
//   - shift
//   - meta
//   - hyper
//   - super
//
// 例如，您将始终看到 "ctrl+shift+alt+a"，而不是 "shift+ctrl+alt+a"。
func (k KeyPressEvent) Keystroke() string {
	return Key(k).Keystroke()
}

// Key 返回底层的键事件。这是将键事件转换为 [Key] 的语法糖。
func (k KeyPressEvent) Key() Key {
	return Key(k)
}

// KeyReleaseEvent 表示按键释放事件。
type KeyReleaseEvent Key

// String 实现了 [fmt.Stringer] 接口，对于匹配键事件非常有用。
// 有关其返回内容的详细信息，请参见 [Key.String]。
func (k KeyReleaseEvent) String() string {
	return Key(k).String()
}

// Keystroke 返回 [Key] 的按键表示形式。虽然在类型安全性上不如查看各个字段，
// 但在匹配按键时，使用此方法通常会更方便、更易读。
//
// 注意，修饰键始终按以下顺序打印：
//   - ctrl
//   - alt
//   - shift
//   - meta
//   - hyper
//   - super
//
// 例如，您将始终看到 "ctrl+shift+alt+a"，而不是 "shift+ctrl+alt+a"。
func (k KeyReleaseEvent) Keystroke() string {
	return Key(k).Keystroke()
}

// Key 返回底层的键事件。这是一个便捷方法和语法糖，用于满足 [KeyEvent] 接口，
// 并将键事件转换为 [Key]。
func (k KeyReleaseEvent) Key() Key {
	return Key(k)
}

// KeyEvent 表示键事件。这可以是按键按下或按键释放事件。
type KeyEvent interface {
	fmt.Stringer

	// Key 返回底层的键事件。
	Key() Key
}

// Key 表示按键按下或释放事件。它包含有关按下的键的信息，如符文、键的类型和按下的修饰键。
// 有几种通用模式可用于检查按键按下或释放：
//
//	// 根据键的字符串表示进行切换（更简短）
//	switch ev := ev.(type) {
//	case KeyPressEvent:
//	    switch ev.String() {
//	    case "enter":
//	        fmt.Println("you pressed enter!")
//	    case "a":
//	        fmt.Println("you pressed a!")
//	    }
//	}
//
//	// 根据键类型进行切换（更可靠）
//	switch ev := ev.(type) {
//	case KeyEvent:
//	    // 同时捕获 KeyPressEvent 和 KeyReleaseEvent
//	    switch key := ev.Key(); key.Code {
//	    case KeyEnter:
//	        fmt.Println("you pressed enter!")
//	    default:
//	        switch key.Text {
//	        case "a":
//	            fmt.Println("you pressed a!")
//	        }
//	    }
//	}
//
// 注意，对于特殊键（如 [KeyEnter]、[KeyTab]）以及不表示可打印字符的键（如带修饰键的组合键），
// [Key.Text] 将为空。换句话说，[Key.Text] 仅在键表示可打印字符（无论是否按下 Shift 键，如 'a'、'A'、'1'、'!' 等）时才会填充。
type Key struct {
	// Text 包含接收到的实际字符。这通常与 [Key.Code] 相同。当 [Key.Text] 非空时，
	// 表示按下的键代表可打印字符。
	Text string

	// Mod 表示修饰键，如 [ModCtrl]、[ModAlt] 等。
	Mod KeyMod

	// Code 表示按下的键。这通常是一个特殊键，如 [KeyTab]、[KeyEnter]、[KeyF1]，
	// 或一个可打印字符，如 'a'。
	Code rune

	// ShiftedCode 是用户实际按下的、经过 Shift 修饰的键。例如，
	// 如果用户按下 shift+a，或 Caps Lock 处于开启状态，[Key.ShiftedCode] 将是 'A'，
	// 而 [Key.Code] 将是 'a'。
	//
	// 对于非拉丁键盘（如阿拉伯语），[Key.ShiftedCode] 是键盘上未按下 Shift 的键。
	//
	// 这仅在 Kitty 键盘协议或 Windows 控制台 API 中可用。
	ShiftedCode rune

	// BaseCode 是根据标准 PC-101 键盘布局按下的键。
	// 在国际键盘上，这是当键盘设置为 US PC-101 布局时会按下的键。
	//
	// 例如，如果用户在法语 AZERTY 键盘上按下 'q'，[Key.BaseCode] 将是 'q'。
	//
	// 这仅在 Kitty 键盘协议或 Windows 控制台 API 中可用。
	BaseCode rune

	// IsRepeat 表示键是否被按住并重复发送事件。
	//
	// 这仅在 Kitty 键盘协议或 Windows 控制台 API 中可用。
	IsRepeat bool
}

// String 实现了 [fmt.Stringer] 接口，对于匹配键事件非常有用。
// 如果 [Key] 有文本表示，它将返回该表示，否则将回退到 [Key.Keystroke]。
//
// 例如，在 US ANSI 键盘上，您将始终得到 "?" 而不是 "shift+/"。
func (k Key) String() string {
	if len(k.Text) > 0 && k.Text != " " {
		return k.Text
	}
	return k.Keystroke()
}

// Keystroke 返回 [Key] 的按键表示形式。虽然在类型安全性上不如查看各个字段，
// 但在匹配按键时，使用此方法通常会更方便、更易读。
//
// 注意，修饰键始终按以下顺序打印：
//   - ctrl
//   - alt
//   - shift
//   - meta
//   - hyper
//   - super
//
// 例如，您将始终看到 "ctrl+shift+alt+a"，而不是 "shift+ctrl+alt+a"。
func (k Key) Keystroke() string {
	var sb strings.Builder
	if k.Mod.Contains(ModCtrl) && k.Code != KeyLeftCtrl && k.Code != KeyRightCtrl {
		sb.WriteString("ctrl+")
	}
	if k.Mod.Contains(ModAlt) && k.Code != KeyLeftAlt && k.Code != KeyRightAlt {
		sb.WriteString("alt+")
	}
	if k.Mod.Contains(ModShift) && k.Code != KeyLeftShift && k.Code != KeyRightShift {
		sb.WriteString("shift+")
	}
	if k.Mod.Contains(ModMeta) && k.Code != KeyLeftMeta && k.Code != KeyRightMeta {
		sb.WriteString("meta+")
	}
	if k.Mod.Contains(ModHyper) && k.Code != KeyLeftHyper && k.Code != KeyRightHyper {
		sb.WriteString("hyper+")
	}
	if k.Mod.Contains(ModSuper) && k.Code != KeyLeftSuper && k.Code != KeyRightSuper {
		sb.WriteString("super+")
	}

	if kt, ok := keyTypeString[k.Code]; ok {
		sb.WriteString(kt)
	} else {
		code := k.Code
		if k.BaseCode != 0 {
			// 如果存在 [Key.BaseCode]，使用它来表示标准 PC-101 键盘布局上的键。
			code = k.BaseCode
		}

		switch code {
		case KeySpace:
			// 空格是唯一不可见的可打印字符。
			sb.WriteString("space")
		case KeyExtended:
			// 当键包含多个符文时，写入键的实际文本。
			sb.WriteString(k.Text)
		default:
			sb.WriteRune(code)
		}
	}

	return sb.String()
}

var keyTypeString = map[rune]string{
	KeyEnter:      "enter",
	KeyTab:        "tab",
	KeyBackspace:  "backspace",
	KeyEscape:     "esc",
	KeySpace:      "space",
	KeyUp:         "up",
	KeyDown:       "down",
	KeyLeft:       "left",
	KeyRight:      "right",
	KeyBegin:      "begin",
	KeyFind:       "find",
	KeyInsert:     "insert",
	KeyDelete:     "delete",
	KeySelect:     "select",
	KeyPgUp:       "pgup",
	KeyPgDown:     "pgdown",
	KeyHome:       "home",
	KeyEnd:        "end",
	KeyKpEnter:    "kpenter",
	KeyKpEqual:    "kpequal",
	KeyKpMultiply: "kpmul",
	KeyKpPlus:     "kpplus",
	KeyKpComma:    "kpcomma",
	KeyKpMinus:    "kpminus",
	KeyKpDecimal:  "kpperiod",
	KeyKpDivide:   "kpdiv",
	KeyKp0:        "kp0",
	KeyKp1:        "kp1",
	KeyKp2:        "kp2",
	KeyKp3:        "kp3",
	KeyKp4:        "kp4",
	KeyKp5:        "kp5",
	KeyKp6:        "kp6",
	KeyKp7:        "kp7",
	KeyKp8:        "kp8",
	KeyKp9:        "kp9",

	// Kitty keyboard extension
	KeyKpSep:    "kpsep",
	KeyKpUp:     "kpup",
	KeyKpDown:   "kpdown",
	KeyKpLeft:   "kpleft",
	KeyKpRight:  "kpright",
	KeyKpPgUp:   "kppgup",
	KeyKpPgDown: "kppgdown",
	KeyKpHome:   "kphome",
	KeyKpEnd:    "kpend",
	KeyKpInsert: "kpinsert",
	KeyKpDelete: "kpdelete",
	KeyKpBegin:  "kpbegin",

	KeyF1:  "f1",
	KeyF2:  "f2",
	KeyF3:  "f3",
	KeyF4:  "f4",
	KeyF5:  "f5",
	KeyF6:  "f6",
	KeyF7:  "f7",
	KeyF8:  "f8",
	KeyF9:  "f9",
	KeyF10: "f10",
	KeyF11: "f11",
	KeyF12: "f12",
	KeyF13: "f13",
	KeyF14: "f14",
	KeyF15: "f15",
	KeyF16: "f16",
	KeyF17: "f17",
	KeyF18: "f18",
	KeyF19: "f19",
	KeyF20: "f20",
	KeyF21: "f21",
	KeyF22: "f22",
	KeyF23: "f23",
	KeyF24: "f24",
	KeyF25: "f25",
	KeyF26: "f26",
	KeyF27: "f27",
	KeyF28: "f28",
	KeyF29: "f29",
	KeyF30: "f30",
	KeyF31: "f31",
	KeyF32: "f32",
	KeyF33: "f33",
	KeyF34: "f34",
	KeyF35: "f35",
	KeyF36: "f36",
	KeyF37: "f37",
	KeyF38: "f38",
	KeyF39: "f39",
	KeyF40: "f40",
	KeyF41: "f41",
	KeyF42: "f42",
	KeyF43: "f43",
	KeyF44: "f44",
	KeyF45: "f45",
	KeyF46: "f46",
	KeyF47: "f47",
	KeyF48: "f48",
	KeyF49: "f49",
	KeyF50: "f50",
	KeyF51: "f51",
	KeyF52: "f52",
	KeyF53: "f53",
	KeyF54: "f54",
	KeyF55: "f55",
	KeyF56: "f56",
	KeyF57: "f57",
	KeyF58: "f58",
	KeyF59: "f59",
	KeyF60: "f60",
	KeyF61: "f61",
	KeyF62: "f62",
	KeyF63: "f63",

	// Kitty keyboard extension
	KeyCapsLock:         "capslock",
	KeyScrollLock:       "scrolllock",
	KeyNumLock:          "numlock",
	KeyPrintScreen:      "printscreen",
	KeyPause:            "pause",
	KeyMenu:             "menu",
	KeyMediaPlay:        "mediaplay",
	KeyMediaPause:       "mediapause",
	KeyMediaPlayPause:   "mediaplaypause",
	KeyMediaReverse:     "mediareverse",
	KeyMediaStop:        "mediastop",
	KeyMediaFastForward: "mediafastforward",
	KeyMediaRewind:      "mediarewind",
	KeyMediaNext:        "medianext",
	KeyMediaPrev:        "mediaprev",
	KeyMediaRecord:      "mediarecord",
	KeyLowerVol:         "lowervol",
	KeyRaiseVol:         "raisevol",
	KeyMute:             "mute",
	KeyLeftShift:        "leftshift",
	KeyLeftAlt:          "leftalt",
	KeyLeftCtrl:         "leftctrl",
	KeyLeftSuper:        "leftsuper",
	KeyLeftHyper:        "lefthyper",
	KeyLeftMeta:         "leftmeta",
	KeyRightShift:       "rightshift",
	KeyRightAlt:         "rightalt",
	KeyRightCtrl:        "rightctrl",
	KeyRightSuper:       "rightsuper",
	KeyRightHyper:       "righthyper",
	KeyRightMeta:        "rightmeta",
	KeyIsoLevel3Shift:   "isolevel3shift",
	KeyIsoLevel5Shift:   "isolevel5shift",
}
