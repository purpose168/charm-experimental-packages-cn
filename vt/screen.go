package vt

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/exp/ordered"
)

// Screen 表示虚拟终端屏幕。
type Screen struct {
	// cb 是要使用的回调结构体。
	cb *Callbacks
	// 屏幕的缓冲区。
	buf uv.Buffer
	// 屏幕的光标和保存的光标。
	cur, saved Cursor
	// scroll 是滚动区域。
	scroll uv.Rectangle
}

// NewScreen 创建一个新屏幕。
func NewScreen(w, h int) *Screen {
	s := Screen{}
	s.Resize(w, h)
	return &s
}

// Reset 重置屏幕。
// 它清除屏幕，将光标设置到左上角，重置光标样式，并重置滚动区域。
func (s *Screen) Reset() {
	s.buf.Clear()
	s.cur = Cursor{}
	s.saved = Cursor{}
	s.scroll = s.buf.Bounds()
}

// Bounds 返回屏幕的边界。
func (s *Screen) Bounds() uv.Rectangle {
	return s.buf.Bounds()
}

// Touched 返回屏幕缓冲区中被修改的行。
func (s *Screen) Touched() []*uv.LineData {
	return s.buf.Touched
}

// CellAt 返回给定x, y位置的单元格。
func (s *Screen) CellAt(x int, y int) *uv.Cell {
	return s.buf.CellAt(x, y)
}

// SetCell 设置给定x, y位置的单元格。
func (s *Screen) SetCell(x, y int, c *uv.Cell) {
	s.buf.SetCell(x, y, c)
}

// Height 返回屏幕的高度。
func (s *Screen) Height() int {
	return s.buf.Height()
}

// Resize 调整屏幕的大小。
func (s *Screen) Resize(width int, height int) {
	s.buf.Resize(width, height)
	s.scroll = s.buf.Bounds()
}

// Width 返回屏幕的宽度。
func (s *Screen) Width() int {
	return s.buf.Width()
}

// Clear 用空白单元格清除屏幕。
func (s *Screen) Clear() {
	s.ClearArea(s.Bounds())
}

// ClearArea 清除给定区域。
func (s *Screen) ClearArea(area uv.Rectangle) {
	s.buf.ClearArea(area)
}

// Fill 填充屏幕或其部分。
func (s *Screen) Fill(c *uv.Cell) {
	s.FillArea(c, s.Bounds())
}

// FillArea 用给定的单元格填充给定区域。
func (s *Screen) FillArea(c *uv.Cell, area uv.Rectangle) {
	s.buf.FillArea(c, area)
}

// setHorizontalMargins 设置水平边距。
func (s *Screen) setHorizontalMargins(left, right int) {
	s.scroll.Min.X = left
	s.scroll.Max.X = right
}

// setVerticalMargins 设置垂直边距。
func (s *Screen) setVerticalMargins(top, bottom int) {
	s.scroll.Min.Y = top
	s.scroll.Max.Y = bottom
}

// setCursorX 设置光标X位置。如果margins为true，则光标仅在滚动边距内设置。
func (s *Screen) setCursorX(x int, margins bool) {
	s.setCursor(x, s.cur.Y, margins)
}

// setCursor 设置光标位置。如果margins为true，则光标仅在滚动边距内设置。这遵循[ansi.CUP]的工作方式。
func (s *Screen) setCursor(x, y int, margins bool) {
	old := s.cur.Position
	if !margins {
		y = ordered.Clamp(y, 0, s.buf.Height()-1) // 限制在屏幕边界内
		x = ordered.Clamp(x, 0, s.buf.Width()-1)  // 限制在屏幕边界内
	} else {
		y = ordered.Clamp(s.scroll.Min.Y+y, s.scroll.Min.Y, s.scroll.Max.Y-1) // 限制在滚动区域内
		x = ordered.Clamp(s.scroll.Min.X+x, s.scroll.Min.X, s.scroll.Max.X-1) // 限制在滚动区域内
	}
	s.cur.X, s.cur.Y = x, y

	// 如果光标位置发生变化，调用回调
	if s.cb.CursorPosition != nil && (old.X != x || old.Y != y) {
		s.cb.CursorPosition(old, uv.Pos(x, y))
	}
}

// moveCursor 按给定的x和y增量移动光标。如果光标位置在滚动区域内，则受滚动区域限制。
// 否则，受屏幕边界限制。
// 这遵循[ansi.CUU]、[ansi.CUD]、[ansi.CUF]、[ansi.CUB]、[ansi.CNL]、[ansi.CPL]的工作方式。
func (s *Screen) moveCursor(dx, dy int) {
	scroll := s.scroll
	old := s.cur.Position
	if old.X < scroll.Min.X {
		scroll.Min.X = 0
	}
	if old.X >= scroll.Max.X {
		scroll.Max.X = s.buf.Width()
	}

	pt := uv.Pos(s.cur.X+dx, s.cur.Y+dy)

	var x, y int
	if old.In(scroll) {
		y = ordered.Clamp(pt.Y, scroll.Min.Y, scroll.Max.Y-1) // 限制在滚动区域内
		x = ordered.Clamp(pt.X, scroll.Min.X, scroll.Max.X-1) // 限制在滚动区域内
	} else {
		y = ordered.Clamp(pt.Y, 0, s.buf.Height()-1) // 限制在屏幕边界内
		x = ordered.Clamp(pt.X, 0, s.buf.Width()-1)  // 限制在屏幕边界内
	}

	s.cur.X, s.cur.Y = x, y

	// 如果光标位置发生变化，调用回调
	if s.cb.CursorPosition != nil && (old.X != x || old.Y != y) {
		s.cb.CursorPosition(old, uv.Pos(x, y))
	}
}

// Cursor 返回光标。
func (s *Screen) Cursor() Cursor {
	return s.cur
}

// CursorPosition 返回光标位置。
func (s *Screen) CursorPosition() (x, y int) {
	return s.cur.X, s.cur.Y
}

// ScrollRegion 返回滚动区域。
func (s *Screen) ScrollRegion() uv.Rectangle {
	return s.scroll
}

// SaveCursor 保存光标。
func (s *Screen) SaveCursor() {
	s.saved = s.cur
}

// RestoreCursor 恢复光标。
func (s *Screen) RestoreCursor() {
	old := s.cur.Position
	s.cur = s.saved

	// 如果光标位置发生变化，调用回调
	if s.cb.CursorPosition != nil && (old.X != s.cur.X || old.Y != s.cur.Y) {
		s.cb.CursorPosition(old, s.cur.Position)
	}
}

// setCursorHidden 设置光标隐藏。
func (s *Screen) setCursorHidden(hidden bool) {
	changed := s.cur.Hidden != hidden
	s.cur.Hidden = hidden
	if changed && s.cb.CursorVisibility != nil {
		s.cb.CursorVisibility(!hidden)
	}
}

// setCursorStyle 设置光标样式。
func (s *Screen) setCursorStyle(style CursorStyle, blink bool) {
	changed := s.cur.Style != style || s.cur.Steady != !blink
	s.cur.Style = style
	s.cur.Steady = !blink
	if changed && s.cb.CursorStyle != nil {
		s.cb.CursorStyle(style, !blink)
	}
}

// cursorPen 返回光标笔。
func (s *Screen) cursorPen() uv.Style {
	return s.cur.Pen
}

// cursorLink 返回光标链接。
func (s *Screen) cursorLink() uv.Link {
	return s.cur.Link
}

// ShowCursor 显示光标。
func (s *Screen) ShowCursor() {
	s.setCursorHidden(false)
}

// HideCursor 隐藏光标。
func (s *Screen) HideCursor() {
	s.setCursorHidden(true)
}

// InsertCell 在光标位置插入n个空白字符，将右侧的单元格向右推并移出屏幕。
func (s *Screen) InsertCell(n int) {
	if n <= 0 {
		return
	}

	x, y := s.cur.X, s.cur.Y
	s.buf.InsertCellArea(x, y, n, s.blankCell(), s.scroll)
}

// DeleteCell 删除光标位置的n个单元格，将左侧的单元格向左移动。
// 如果光标在滚动区域外，这没有效果。
func (s *Screen) DeleteCell(n int) {
	if n <= 0 {
		return
	}

	x, y := s.cur.X, s.cur.Y
	s.buf.DeleteCellArea(x, y, n, s.blankCell(), s.scroll)
}

// ScrollUp 在给定区域内向上滚动内容n行。超过上边缘滚动的行将丢失。
// 这相当于[ansi.SU]，它将光标移动到上边缘并执行[ansi.DL]操作。
func (s *Screen) ScrollUp(n int) {
	x, y := s.CursorPosition()
	s.setCursor(s.cur.X, 0, true) // 移动到上边缘
	s.DeleteLine(n)                // 删除n行
	s.setCursor(x, y, false)       // 恢复光标位置
}

// ScrollDown 在给定区域内向下滚动内容n行。超过下边缘滚动的行将丢失。
// 这相当于[ansi.SD]，它将光标移动到上边缘并执行[ansi.IL]操作。
func (s *Screen) ScrollDown(n int) {
	x, y := s.CursorPosition()
	s.setCursor(s.cur.X, 0, true) // 移动到上边缘
	s.InsertLine(n)                // 插入n行
	s.setCursor(x, y, false)       // 恢复光标位置
}

// InsertLine 在光标位置Y坐标处插入n个空白行。
// 仅当光标在滚动区域内时操作。光标Y下方的行向下移动，超过下边缘的行被丢弃。
// 如果操作成功，返回true。
func (s *Screen) InsertLine(n int) bool {
	if n <= 0 {
		return false
	}

	x, y := s.cur.X, s.cur.Y

	// 仅当光标Y在滚动区域内时操作
	if y < s.scroll.Min.Y || y >= s.scroll.Max.Y ||
		x < s.scroll.Min.X || x >= s.scroll.Max.X {
		return false
	}

	s.buf.InsertLineArea(y, n, s.blankCell(), s.scroll)

	return true
}

// DeleteLine 在光标位置Y坐标处删除n行。
// 仅当光标在滚动区域内时操作。光标Y下方的行向上移动，滚动区域底部插入空白行。
// 如果操作成功，返回true。
func (s *Screen) DeleteLine(n int) bool {
	if n <= 0 {
		return false
	}

	scroll := s.scroll
	x, y := s.cur.X, s.cur.Y

	// 仅当光标Y在滚动区域内时操作
	if y < scroll.Min.Y || y >= scroll.Max.Y ||
		x < scroll.Min.X || x >= scroll.Max.X {
		return false
	}

	s.buf.DeleteLineArea(y, n, s.blankCell(), scroll)

	return true
}

// blankCell 返回光标空白单元格，背景颜色设置为当前笔背景颜色。
// 如果笔背景颜色为nil，返回值为nil。
func (s *Screen) blankCell() *uv.Cell {
	if s.cur.Pen.Bg == nil {
		return nil
	}

	c := uv.EmptyCell
	c.Style.Bg = s.cur.Pen.Bg
	return &c
}
