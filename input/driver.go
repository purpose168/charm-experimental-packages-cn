//nolint:unused,revive,nolintlint
package input

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/muesli/cancelreader"
)

// Logger 是一个简单的日志记录器接口。
type Logger interface {
	Printf(format string, v ...any)
}

// win32InputState 是一个状态机，用于解析来自 Windows 控制台 API 的按键事件为转义序列和 UTF-8 符文，
// 并跟踪上一次控制键状态以确定修饰键的变化。它还跟踪上一次鼠标按钮状态和窗口大小变化，
// 以确定哪些鼠标按钮被释放，并防止多个大小事件触发。
type win32InputState struct {
	ansiBuf                    [256]byte
	ansiIdx                    int
	utf16Buf                   [2]rune
	utf16Half                  bool
	lastCks                    uint32 // 上一个事件的控制键状态
	lastMouseBtns              uint32 // 上一个事件的鼠标按钮状态
	lastWinsizeX, lastWinsizeY int16  // 上一个事件的窗口大小，用于防止多个大小事件触发
}

// Reader 表示输入事件读取器。它从终端输入缓冲区读取输入事件并解析转义序列，
// 将它们转换为人类可读的事件。
type Reader struct {
	rd    cancelreader.CancelReader
	table map[string]Key // table 是按键序列的查找表。

	term string // term 是终端名称 $TERM。

	// paste 是括号粘贴模式缓冲区。
	// 为 nil 时，括号粘贴模式被禁用。
	paste []byte

	buf [256]byte // 我们需要更大的缓冲区吗？

	// keyState 跟踪当前 Windows 控制台 API 按键事件状态。
	// 它用于解码 ANSI 转义序列和 UTF-16 序列。
	keyState win32InputState

	parser Parser
	logger Logger
}

// NewReader 返回一个新的输入事件读取器。该读取器从终端读取输入事件，
// 并将转义序列解析为人类可读的事件。它支持读取 Terminfo 数据库。
// 有关更多信息，请参阅 [Parser]。
//
// 示例：
//
//	r, _ := input.NewReader(os.Stdin, os.Getenv("TERM"), 0)
//	defer r.Close()
//	events, _ := r.ReadEvents()
//	for _, ev := range events {
//	  log.Printf("%v", ev)
//	}
func NewReader(r io.Reader, termType string, flags int) (*Reader, error) {
	d := new(Reader)
	cr, err := newCancelreader(r, flags)
	if err != nil {
		return nil, err
	}

	d.rd = cr
	d.table = buildKeysTable(flags, termType)
	d.term = termType
	d.parser.flags = flags
	return d, nil
}

// SetLogger 为读取器设置日志记录器。
func (d *Reader) SetLogger(l Logger) {
	d.logger = l
}

// Read 实现 [io.Reader] 接口。
func (d *Reader) Read(p []byte) (int, error) {
	return d.rd.Read(p) //nolint:wrapcheck
}

// Cancel 取消底层读取器。
func (d *Reader) Cancel() bool {
	return d.rd.Cancel()
}

// Close 关闭底层读取器。
func (d *Reader) Close() error {
	return d.rd.Close() //nolint:wrapcheck
}

func (d *Reader) readEvents() ([]Event, error) {
	nb, err := d.rd.Read(d.buf[:])
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var events []Event
	buf := d.buf[:nb]

	// 首先查找表
	if bytes.HasPrefix(buf, []byte{'\x1b'}) {
		if k, ok := d.table[string(buf)]; ok {
			if d.logger != nil {
				d.logger.Printf("input: %q", buf)
			}
			events = append(events, KeyPressEvent(k))
			return events, nil
		}
	}

	var i int
	for i < len(buf) {
		nb, ev := d.parser.parseSequence(buf[i:])
		if d.logger != nil {
			d.logger.Printf("input: %q", buf[i:i+nb])
		}

		// 处理括号粘贴
		if d.paste != nil {
			if _, ok := ev.(PasteEndEvent); !ok {
				d.paste = append(d.paste, buf[i])
				i++
				continue
			}
		}

		switch ev.(type) {
		case UnknownEvent:
			// 如果解析器不识别该序列，尝试在查找表中查找。
			if k, ok := d.table[string(buf[i:i+nb])]; ok {
				ev = KeyPressEvent(k)
			}
		case PasteStartEvent:
			d.paste = []byte{}
		case PasteEndEvent:
			// 将捕获的数据解码为符文。
			var paste []rune
			for len(d.paste) > 0 {
				r, w := utf8.DecodeRune(d.paste)
				if r != utf8.RuneError {
					paste = append(paste, r)
				}
				d.paste = d.paste[w:]
			}
			d.paste = nil // 重置缓冲区
			events = append(events, PasteEvent(paste))
		case nil:
			i++
			continue
		}

		if mevs, ok := ev.(MultiEvent); ok {
			events = append(events, []Event(mevs)...)
		} else {
			events = append(events, ev)
		}
		i += nb
	}

	return events, nil
}
