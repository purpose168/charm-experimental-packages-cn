// Package vttest 提供了一个用于测试终端应用程序的虚拟终端实现。它允许你创建一个带有伪终端（PTY）的终端实例，并在任何时刻捕获其状态，使你能够编写测试来验证终端应用程序的行为。
package vttest

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"maps"
	"os"
	"os/exec"
	"sync"
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/vt"
	"github.com/purpose168/charm-experimental-packages-cn/xpty"
)

// Terminal 表示一个带有 PTY 和状态的虚拟终端。
type Terminal struct {
	Emulator *vt.SafeEmulator
	tb       testing.TB

	cols, rows  int
	title       string
	altScreen   bool
	ansiModes   map[ansi.ANSIMode]ansi.ModeSetting
	decModes    map[ansi.DECMode]ansi.ModeSetting
	cursorPos   image.Point
	cursorVis   bool
	cursorColor color.Color
	cursorStyle vt.CursorStyle
	cursorBlink bool
	bgColor     color.Color
	fgColor     color.Color

	pty    xpty.Pty
	ptyIn  io.Reader
	ptyOut io.Writer

	mu sync.Mutex
}

// NewTerminal 创建一个具有给定大小的新虚拟终端，用于测试目的。在任何时刻，你可以通过在返回的 Terminal 实例上调用 [Terminal.Snapshot] 方法来拍摄终端状态的快照。
func NewTerminal(tb testing.TB, cols, rows int) (*Terminal, error) {
	pty, err := xpty.NewPty(cols, rows)
	if err != nil {
		return nil, fmt.Errorf("failed to create pty: %w", err)
	}

	term := new(Terminal)
	term.tb = tb
	term.cols = cols
	term.rows = rows
	term.ansiModes = make(map[ansi.ANSIMode]ansi.ModeSetting)
	term.decModes = make(map[ansi.DECMode]ansi.ModeSetting)

	switch p := pty.(type) {
	case *xpty.UnixPty:
		term.ptyIn = p.Slave()
		term.ptyOut = p.Slave()
	case *xpty.ConPty:
		inFile := os.NewFile(p.InPipeReadFd(), "|0")
		outFile := os.NewFile(p.OutPipeWriteFd(), "|1")
		term.ptyIn = inFile
		term.ptyOut = outFile
	}

	vterm := vt.NewSafeEmulator(cols, rows)
	vterm.SetCallbacks(vt.Callbacks{
		Title: func(title string) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.title = title
		},
		AltScreen: func(alt bool) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.altScreen = alt
		},
		EnableMode: func(mode ansi.Mode) {
			term.mu.Lock()
			defer term.mu.Unlock()
			switch m := mode.(type) {
			case ansi.ANSIMode:
				term.ansiModes[m] = ansi.ModeSet
			case ansi.DECMode:
				term.decModes[m] = ansi.ModeSet
			}
		},
		DisableMode: func(mode ansi.Mode) {
			term.mu.Lock()
			defer term.mu.Unlock()
			switch m := mode.(type) {
			case ansi.ANSIMode:
				term.ansiModes[m] = ansi.ModeReset
			case ansi.DECMode:
				term.decModes[m] = ansi.ModeReset
			}
		},
		CursorPosition: func(_, newpos uv.Position) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.cursorPos = newpos
		},
		CursorVisibility: func(visible bool) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.cursorVis = visible
		},
		CursorStyle: func(style vt.CursorStyle, blink bool) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.cursorStyle = style
			term.cursorBlink = blink
		},
		CursorColor: func(color color.Color) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.cursorColor = color
		},
		BackgroundColor: func(color color.Color) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.bgColor = color
		},
		ForegroundColor: func(color color.Color) {
			term.mu.Lock()
			defer term.mu.Unlock()
			term.fgColor = color
		},
	})

	term.Emulator = vterm
	term.pty = pty

	// Copy PTY input to terminal
	go io.Copy(vterm, pty) //nolint:errcheck
	// Copy terminal output to PTY
	go io.Copy(pty, vterm) //nolint:errcheck

	return term, nil
}

// Start 启动一个附加到终端 PTY 的进程。
func (t *Terminal) Start(cmd *exec.Cmd) error {
	if err := t.pty.Start(cmd); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	return nil
}

// Wait 等待附加到终端 PTY 的进程退出。
func (t *Terminal) Wait(cmd *exec.Cmd) error {
	if err := xpty.WaitProcess(t.tb.Context(), cmd); err != nil {
		return fmt.Errorf("process exited with error: %w", err)
	}
	return nil
}

// Close 关闭终端及其 PTY。
func (t *Terminal) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.Emulator.Close(); err != nil && !errors.Is(err, io.EOF) {
		_ = t.pty.Close()
		return fmt.Errorf("failed to close emulator: %w", err)
	}

	if err := t.pty.Close(); err != nil {
		return fmt.Errorf("failed to close pty: %w", err)
	}
	return nil
}

// Resize 调整终端及其 PTY 的大小。
func (t *Terminal) Resize(cols, rows int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.cols = cols
	t.rows = rows
	t.Emulator.Resize(cols, rows)
	if err := t.pty.Resize(cols, rows); err != nil {
		return fmt.Errorf("failed to resize pty: %w", err)
	}

	return nil
}

// SendText 将给定的原始文本发送到终端模拟器，就像用户输入的一样。
func (t *Terminal) SendText(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Emulator.SendText(text)
}

// SendKey 将给定的按键事件发送到终端模拟器，就像用户输入的一样。
func (t *Terminal) SendKey(k uv.KeyEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Emulator.SendKey(k)
}

// SendMouse 将给定的鼠标事件发送到终端模拟器，就像用户执行的一样。
func (t *Terminal) SendMouse(m uv.MouseEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Emulator.SendMouse(m)
}

// Paste 将给定的文本发送到终端模拟器，就像用户粘贴的一样。
func (t *Terminal) Paste(text string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Emulator.Paste(text)
}

// Input 返回终端 PTY 的输入端。
func (t *Terminal) Input() io.Reader {
	return t.ptyIn
}

// Output 返回终端 PTY 的输出端。
func (t *Terminal) Output() io.Writer {
	return t.ptyOut
}

// Snapshot 拍摄当前终端状态的快照。
// 返回的 [Snapshot] 可用于检查拍摄快照时的终端状态。它还可以序列化为 JSON 或 YAML，用于进一步分析或测试目的。
func (t *Terminal) Snapshot() Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()

	snap := Snapshot{
		Modes: Modes{
			ANSI: maps.Clone(t.ansiModes),
			DEC:  maps.Clone(t.decModes),
		},
		Title:     t.title,
		Rows:      t.rows,
		Cols:      t.cols,
		AltScreen: t.altScreen,
		Cursor: Cursor{
			Position: Position(t.cursorPos),
			Visible:  t.cursorVis,
			Color:    Color{t.cursorColor},
			Style:    t.cursorStyle,
			Blink:    t.cursorBlink,
		},
		BgColor: Color{t.bgColor},
		FgColor: Color{t.fgColor},
		Cells:   make([][]Cell, t.rows),
	}

	for r := 0; r < t.rows; r++ {
		snap.Cells[r] = make([]Cell, t.cols)
		for c := 0; c < t.cols; c++ {
			cell := t.Emulator.CellAt(c, r)
			snap.Cells[r][c] = Cell{
				Content: cell.Content,
				Style: Style{
					Fg:             Color{cell.Style.Fg},
					Bg:             Color{cell.Style.Bg},
					UnderlineColor: Color{cell.Style.UnderlineColor},
					Underline:      cell.Style.Underline,
					Attrs:          cell.Style.Attrs,
				},
				Link:  Link(cell.Link),
				Width: cell.Width,
			}
		}
	}

	return snap
}
