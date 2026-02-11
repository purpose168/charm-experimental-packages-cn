package vt

import (
	"image/color"
	"io"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/screen"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/ansi/parser"
)

// Logger 表示日志器接口。
type Logger interface {
	Printf(format string, v ...any) // 格式化输出日志
}

// Emulator 表示虚拟终端模拟器。
type Emulator struct {
	handlers

	// 终端的256个索引颜色。
	colors [256]color.Color

	// 主屏幕和备用屏幕，以及指向当前活动屏幕的指针。
	scrs [2]Screen
	scr  *Screen

	// 字符集
	charsets [4]CharSet

	// logger 是要使用的日志器。
	logger Logger

	// 终端默认颜色。
	defaultFg, defaultBg, defaultCur color.Color
	fgColor, bgColor, curColor       color.Color

	// 终端模式。
	modes ansi.Modes

	// 最后写入的字符。
	lastChar rune // 要么是ansi.Rune，要么是ansi.Grapheme
	// 用于组成字形的符文切片。
	grapheme []rune

	// 要使用的ANSI解析器。
	parser *ansi.Parser
	// 最后的解析器状态。
	lastState parser.State

	cb Callbacks

	// 终端的图标名称和标题。
	iconName, title string
	// 当前报告的工作目录。这不会被验证。
	cwd string

	// tabstops 是制表位列表。
	tabstops *uv.TabStops

	// I/O管道。
	pr *io.PipeReader
	pw *io.PipeWriter

	// GL和GR字符集标识符。
	gl, gr  int
	gsingle int // 临时选择GL或GR

	// 指示终端是否已关闭。
	closed bool

	// atPhantom 指示光标是否越界。
	// 当为true时，写入字符时，光标会移动到下一行。
	atPhantom bool
}

var _ Terminal = (*Emulator)(nil)

// NewEmulator 创建一个新的虚拟终端模拟器。
func NewEmulator(w, h int) *Emulator {
	t := new(Emulator)
	t.scrs[0] = *NewScreen(w, h) // 创建主屏幕
	t.scrs[1] = *NewScreen(w, h) // 创建备用屏幕
	t.scr = &t.scrs[0] // 默认使用主屏幕
	t.scrs[0].cb = &t.cb // 设置主屏幕的回调
	t.scrs[1].cb = &t.cb // 设置备用屏幕的回调
	t.parser = ansi.NewParser() // 创建ANSI解析器
	t.parser.SetParamsSize(parser.MaxParamsSize) // 设置参数大小
	t.parser.SetDataSize(1024 * 1024 * 4) // 4MB data buffer // 设置数据缓冲区大小
	t.parser.SetHandler(ansi.Handler{
		Print:     t.handlePrint,
		Execute:   t.handleControl,
		HandleCsi: t.handleCsi,
		HandleEsc: t.handleEsc,
		HandleDcs: t.handleDcs,
		HandleOsc: t.handleOsc,
		HandleApc: t.handleApc,
		HandlePm:  t.handlePm,
		HandleSos: t.handleSos,
	})
	t.pr, t.pw = io.Pipe() // 创建I/O管道
	t.resetModes() // 重置终端模式
	t.tabstops = uv.DefaultTabStops(w) // 设置默认制表位
	t.registerDefaultHandlers() // 注册默认处理器

	// 默认颜色
	t.defaultFg = color.White // 默认前景色为白色
	t.defaultBg = color.Black // 默认背景色为黑色
	t.defaultCur = color.White // 默认光标颜色为白色

	return t
}

// SetLogger 设置终端的日志器。
func (e *Emulator) SetLogger(l Logger) {
	e.logger = l
}

// SetCallbacks 设置终端的回调。
func (e *Emulator) SetCallbacks(cb Callbacks) {
	e.cb = cb
	e.scrs[0].cb = &e.cb
	e.scrs[1].cb = &e.cb
}

// Touched 返回当前屏幕缓冲区中被修改的行。
func (e *Emulator) Touched() []*uv.LineData {
	return e.scr.Touched()
}

// String 返回底层屏幕缓冲区的字符串表示。
func (e *Emulator) String() string {
	s := e.scr.buf.String()
	return uv.TrimSpace(s)
}

// Render 将终端屏幕的快照渲染为字符串，样式和链接编码为ANSI转义序列。
func (e *Emulator) Render() string {
	return e.scr.buf.Render()
}

var _ uv.Screen = (*Emulator)(nil)

// Bounds 返回终端的边界。
func (e *Emulator) Bounds() uv.Rectangle {
	return e.scr.Bounds()
}

// CellAt 返回给定x, y位置的当前焦点屏幕单元格。
// 如果单元格越界，则返回nil。
func (e *Emulator) CellAt(x, y int) *uv.Cell {
	return e.scr.CellAt(x, y)
}

// SetCell 设置给定x, y位置的当前焦点屏幕单元格。
func (e *Emulator) SetCell(x, y int, c *uv.Cell) {
	e.scr.SetCell(x, y, c)
}

// WidthMethod 返回终端使用的宽度计算方法。
func (e *Emulator) WidthMethod() uv.WidthMethod {
	if e.isModeSet(ansi.ModeUnicodeCore) {
		return ansi.GraphemeWidth
	}
	return ansi.WcWidth
}

// Draw 实现[uv.Drawable]接口。
func (e *Emulator) Draw(scr uv.Screen, area uv.Rectangle) {
	bg := uv.EmptyCell
	bg.Style.Bg = e.BackgroundColor()
	screen.FillArea(scr, &bg, area) // 填充背景
	for y := range e.Touched() {
		if y < 0 || y >= e.Height() {
			continue
		}
		for x := 0; x < e.Width(); {
			w := 1
			cell := e.CellAt(x, y)
			if cell != nil {
				cell = cell.Clone()
				if cell.Width > 1 {
					w = cell.Width
				}
				if cell.Style.Bg == nil && e.bgColor != nil {
					cell.Style.Bg = e.bgColor
				}
				if cell.Style.Fg == nil && e.fgColor != nil {
					cell.Style.Fg = e.fgColor
				}
				scr.SetCell(x+area.Min.X, y+area.Min.Y, cell)
			}
			x += w
		}
	}
}

// Height 返回终端的高度。
func (e *Emulator) Height() int {
	return e.scr.Height()
}

// Width 返回终端的宽度。
func (e *Emulator) Width() int {
	return e.scr.Width()
}

// CursorPosition 返回终端的光标位置。
func (e *Emulator) CursorPosition() uv.Position {
	x, y := e.scr.CursorPosition()
	return uv.Pos(x, y)
}

// Resize 调整终端的大小。
func (e *Emulator) Resize(width int, height int) {
	x, y := e.scr.CursorPosition()
	if e.atPhantom {
		if x < width-1 {
			e.atPhantom = false
			x++
		}
	}

	// 确保光标位置在新的边界内
	if y < 0 {
		y = 0
	}
	if y >= height {
		y = height - 1
	}
	if x < 0 {
		x = 0
	}
	if x >= width {
		x = width - 1
	}

	// 调整屏幕大小
	e.scrs[0].Resize(width, height)
	e.scrs[1].Resize(width, height)
	e.tabstops = uv.DefaultTabStops(width) // 重置制表位

	e.setCursor(x, y) // 设置光标位置

	// 如果启用了带内调整大小模式，发送调整大小事件
	if e.isModeSet(ansi.ModeInBandResize) {
		_, _ = io.WriteString(e.pw, ansi.InBandResize(e.Height(), e.Width(), 0, 0))
	}
}

// Read 从终端输入缓冲区读取数据。
func (e *Emulator) Read(p []byte) (n int, err error) {
	if e.closed {
		return 0, io.EOF
	}

	return e.pr.Read(p) //nolint:wrapcheck
}

// Close 关闭终端。
func (e *Emulator) Close() error {
	if e.closed {
		return nil
	}

	e.closed = true
	return e.pw.CloseWithError(io.EOF)
}

// Write 将数据写入终端输出缓冲区。
func (e *Emulator) Write(p []byte) (n int, err error) {
	if e.closed {
		return 0, io.ErrClosedPipe
	}

	for i := range p {
		e.parser.Advance(p[i])
		state := e.parser.State()
		// 如果我们转换到非utf8状态或已写入整个字节切片，则刷新字形
		if len(e.grapheme) > 0 {
			if (e.lastState == parser.GroundState && state != parser.Utf8State) || i == len(p)-1 {
				e.flushGrapheme()
			}
		}
		e.lastState = state
	}
	return len(p), nil
}

// WriteString 将字符串写入终端输出缓冲区。
func (e *Emulator) WriteString(s string) (n int, err error) {
	return e.Write([]byte(s)) //nolint:wrapcheck
}

// InputPipe 返回终端的输入管道。
// 这可用于向终端发送输入。
func (e *Emulator) InputPipe() io.Writer {
	return e.pw
}

// Paste 将文本粘贴到终端。
// 如果启用了括号粘贴模式，文本会用适当的转义序列括起来。
func (e *Emulator) Paste(text string) {
	if e.isModeSet(ansi.ModeBracketedPaste) {
		_, _ = io.WriteString(e.pw, ansi.BracketedPasteStart)
		defer io.WriteString(e.pw, ansi.BracketedPasteEnd) //nolint:errcheck
	}

	_, _ = io.WriteString(e.pw, text)
}

// SendText 向终端发送任意文本。
func (e *Emulator) SendText(text string) {
	_, _ = io.WriteString(e.pw, text)
}

// SendKeys 向终端发送多个按键。
func (e *Emulator) SendKeys(keys ...uv.KeyEvent) {
	for _, k := range keys {
		e.SendKey(k)
	}
}

// ForegroundColor 返回终端的前景颜色。如果前景颜色未设置，这将返回nil，
// 这意味着使用外部终端颜色。
func (e *Emulator) ForegroundColor() color.Color {
	if e.fgColor == nil {
		return e.defaultFg
	}
	return e.fgColor
}

// SetForegroundColor 设置终端的前景颜色。
func (e *Emulator) SetForegroundColor(c color.Color) {
	if c == nil {
		c = e.defaultFg
	}
	e.fgColor = c
	if e.cb.ForegroundColor != nil {
		e.cb.ForegroundColor(c)
	}
}

// SetDefaultForegroundColor 设置终端的默认前景颜色。
func (e *Emulator) SetDefaultForegroundColor(c color.Color) {
	if c == nil {
		c = color.White
	}
	e.defaultFg = c
}

// BackgroundColor 返回终端的背景颜色。如果背景颜色未设置，这将返回nil，
// 这意味着使用外部终端颜色。
func (e *Emulator) BackgroundColor() color.Color {
	if e.bgColor == nil {
		return e.defaultBg
	}
	return e.bgColor
}

// SetBackgroundColor 设置终端的背景颜色。
func (e *Emulator) SetBackgroundColor(c color.Color) {
	if c == nil {
		c = e.defaultBg
	}
	e.bgColor = c
	if e.cb.BackgroundColor != nil {
		e.cb.BackgroundColor(c)
	}
}

// SetDefaultBackgroundColor 设置终端的默认背景颜色。
func (e *Emulator) SetDefaultBackgroundColor(c color.Color) {
	if c == nil {
		c = color.Black
	}
	e.defaultBg = c
}

// CursorColor 返回终端的光标颜色。如果光标颜色未设置，这将返回nil，
// 这意味着使用外部终端颜色。
func (e *Emulator) CursorColor() color.Color {
	if e.curColor == nil {
		return e.defaultCur
	}
	return e.curColor
}

// SetCursorColor 设置终端的光标颜色。
func (e *Emulator) SetCursorColor(c color.Color) {
	if c == nil {
		c = e.defaultCur
	}
	e.curColor = c
	if e.cb.CursorColor != nil {
		e.cb.CursorColor(c)
	}
}

// SetDefaultCursorColor 设置终端的默认光标颜色。
func (e *Emulator) SetDefaultCursorColor(c color.Color) {
	if c == nil {
		c = color.White
	}
	e.defaultCur = c
}

// IndexedColor 返回终端的索引颜色。索引颜色是介于0和255之间的颜色。
func (e *Emulator) IndexedColor(i int) color.Color {
	if i < 0 || i > 255 {
		return nil
	}

	c := e.colors[i]
	if c == nil {
		// 返回默认颜色。
		return ansi.IndexedColor(i) //nolint:gosec
	}

	return c
}

// SetIndexedColor 设置终端的索引颜色。
// 索引必须介于0和255之间。
func (e *Emulator) SetIndexedColor(i int, c color.Color) {
	if i < 0 || i > 255 {
		return
	}

	e.colors[i] = c
}

// resetTabStops 将终端制表位重置为默认设置。
func (e *Emulator) resetTabStops() {
	e.tabstops = uv.DefaultTabStops(e.Width())
}

// logf 记录日志。
func (e *Emulator) logf(format string, v ...any) {
	if e.logger != nil {
		e.logger.Printf(format, v...)
	}
}
