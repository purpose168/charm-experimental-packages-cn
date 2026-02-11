// Package cellbuf 提供终端单元格缓冲区功能。
package cellbuf

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// NewCell 返回一个新的单元格。这是一个便捷函数，使用给定内容初始化新单元格。
// 单元格的宽度通过 [runewidth.RuneWidth] 根据内容确定。
// 这只会考虑内容中的第一个组合符文。如果内容为空，将返回宽度为 0 的空单元格。
func NewCell(r rune, comb ...rune) (c *Cell) {
	c = new(Cell)
	c.Rune = r
	c.Width = runewidth.RuneWidth(r)
	for _, r := range comb {
		if runewidth.RuneWidth(r) > 0 {
			break
		}
		c.Comb = append(c.Comb, r)
	}
	c.Comb = comb
	c.Width = runewidth.StringWidth(string(append([]rune{r}, comb...)))
	return c
}

// NewCellString 返回一个带有给定字符串内容的新单元格。这是一个便捷函数，
// 使用给定内容初始化新单元格。单元格的宽度通过 [uniseg.FirstGraphemeClusterInString] 根据内容确定。
// 这只会使用字符串中的第一个字形集群。如果字符串为空，将返回宽度为 0 的空单元格。
func NewCellString(s string) (c *Cell) {
	c = new(Cell)
	for i, r := range s {
		if i == 0 {
			c.Rune = r
			// 我们只关心第一个符文的宽度
			c.Width = runewidth.RuneWidth(r)
		} else {
			if runewidth.RuneWidth(r) > 0 {
				break
			}
			c.Comb = append(c.Comb, r)
		}
	}
	return c
}

// NewGraphemeCell 返回一个新的单元格。这是一个便捷函数，使用给定内容初始化新单元格。
// 单元格的宽度通过 [uniseg.FirstGraphemeClusterInString] 根据内容确定。
// 当内容是一个字形集群（grapheme cluster）时使用，即形成单个视觉单元的符文序列。
// 这只会返回字符串中的第一个字形集群。如果字符串为空，将返回宽度为 0 的空单元格。
func NewGraphemeCell(s string) (c *Cell) {
	g, _, w, _ := uniseg.FirstGraphemeClusterInString(s, -1)
	return newGraphemeCell(g, w)
}

func newGraphemeCell(s string, w int) (c *Cell) {
	c = new(Cell)
	c.Width = w
	for i, r := range s {
		if i == 0 {
			c.Rune = r
		} else {
			c.Comb = append(c.Comb, r)
		}
	}
	return c
}

// Line 表示终端中的一行。
// nil 单元格表示空白单元格，即带有空格字符且宽度为 1 的单元格。
// 如果单元格没有内容且宽度为 0，则它是宽单元格的占位符。
type Line []*Cell

// Width 返回行的宽度。
func (l Line) Width() int {
	return len(l)
}

// Len 返回行的长度。
func (l Line) Len() int {
	return len(l)
}

// String 返回行的字符串表示形式。任何尾随空格都会被移除。
func (l Line) String() (s string) {
	for _, c := range l {
		if c == nil {
			s += " "
		} else if c.Empty() {
			continue
		} else {
			s += c.String()
		}
	}
	s = strings.TrimRight(s, " ")
	return s
}

// At 返回给定 x 位置的单元格。
// 如果单元格不存在，返回 nil。
func (l Line) At(x int) *Cell {
	if x < 0 || x >= len(l) {
		return nil
	}

	c := l[x]
	if c == nil {
		newCell := BlankCell
		return &newCell
	}

	return c
}

// Set 设置给定 x 位置的单元格。如果给定宽单元格，它会设置该单元格及其后续单元格为 [EmptyCell]。
// 如果单元格被设置，返回 true。
func (l Line) Set(x int, c *Cell) bool {
	return l.set(x, c, true)
}

func (l Line) set(x int, c *Cell, clone bool) bool {
	width := l.Width()
	if x < 0 || x >= width {
		return false
	}

	// 当宽单元格被部分覆盖时，我们需要
	// 用空格单元格填充剩余部分，以
	// 避免渲染问题。
	prev := l.At(x)
	if prev != nil && prev.Width > 1 {
		// 写入第一个宽单元格
		for j := 0; j < prev.Width && x+j < l.Width(); j++ {
			l[x+j] = prev.Clone().Blank()
		}
	} else if prev != nil && prev.Width == 0 {
		// 写入宽单元格占位符
		for j := 1; j < maxCellWidth && x-j >= 0; j++ {
			wide := l.At(x - j)
			if wide != nil && wide.Width > 1 && j < wide.Width {
				for k := range wide.Width {
					l[x-j+k] = wide.Clone().Blank()
				}
				break
			}
		}
	}

	if clone && c != nil {
		// 如果不为 nil，则克隆单元格。
		c = c.Clone()
	}

	if c != nil && x+c.Width > width {
		// 如果单元格太宽，我们用相同样式的空白填充。
		for i := 0; i < c.Width && x+i < width; i++ {
			l[x+i] = c.Clone().Blank()
		}
	} else {
		l[x] = c

		// 用宽度为 0 的空单元格标记宽单元格
		// 我们在下面设置宽单元格
		if c != nil && c.Width > 1 {
			for j := 1; j < c.Width && x+j < l.Width(); j++ {
				var wide Cell
				l[x+j] = &wide
			}
		}
	}

	return true
}

// Buffer 是表示屏幕或终端的二维单元格网格。
type Buffer struct {
	// Lines 保存缓冲区的行。
	Lines []Line
}

// NewBuffer 创建具有给定宽度和高度的新缓冲区。
// 这是一个便捷函数，初始化新缓冲区并调整其大小。
func NewBuffer(width int, height int) *Buffer {
	b := new(Buffer)
	b.Resize(width, height)
	return b
}

// String 返回缓冲区的字符串表示形式。
func (b *Buffer) String() (s string) {
	for i, l := range b.Lines {
		s += l.String()
		if i < len(b.Lines)-1 {
			s += "\r\n"
		}
	}
	return s
}

// Line 返回给定 y 位置的行指针。
// 如果行不存在，返回 nil。
func (b *Buffer) Line(y int) Line {
	if y < 0 || y >= len(b.Lines) {
		return nil
	}
	return b.Lines[y]
}

// Cell 实现 Screen 接口。
func (b *Buffer) Cell(x int, y int) *Cell {
	if y < 0 || y >= len(b.Lines) {
		return nil
	}
	return b.Lines[y].At(x)
}

// maxCellWidth 是终端单元格可以获得的最大宽度。
const maxCellWidth = 4

// SetCell 设置给定 x, y 位置的单元格。
func (b *Buffer) SetCell(x, y int, c *Cell) bool {
	return b.setCell(x, y, c, true)
}

// setCell 设置给定 x, y 位置的单元格。如果 c 不为 nil，这将始终克隆并分配新单元格。
func (b *Buffer) setCell(x, y int, c *Cell, clone bool) bool {
	if y < 0 || y >= len(b.Lines) {
		return false
	}
	return b.Lines[y].set(x, c, clone)
}

// Height 实现 Screen 接口。
func (b *Buffer) Height() int {
	return len(b.Lines)
}

// Width 实现 Screen 接口。
func (b *Buffer) Width() int {
	if len(b.Lines) == 0 {
		return 0
	}
	return b.Lines[0].Width()
}

// Bounds 返回缓冲区的边界。
func (b *Buffer) Bounds() Rectangle {
	return Rect(0, 0, b.Width(), b.Height())
}

// Resize 将缓冲区调整为给定的宽度和高度。
func (b *Buffer) Resize(width int, height int) {
	if width == 0 || height == 0 {
		b.Lines = nil
		return
	}

	if width > b.Width() {
		line := make(Line, width-b.Width())
		for i := range b.Lines {
			b.Lines[i] = append(b.Lines[i], line...)
		}
	} else if width < b.Width() {
		for i := range b.Lines {
			b.Lines[i] = b.Lines[i][:width]
		}
	}

	if height > len(b.Lines) {
		for i := len(b.Lines); i < height; i++ {
			b.Lines = append(b.Lines, make(Line, width))
		}
	} else if height < len(b.Lines) {
		b.Lines = b.Lines[:height]
	}
}

// FillRect 用给定的单元格和矩形填充缓冲区。
func (b *Buffer) FillRect(c *Cell, rect Rectangle) {
	cellWidth := 1
	if c != nil && c.Width > 1 {
		cellWidth = c.Width
	}
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x += cellWidth {
			b.setCell(x, y, c, false)
		}
	}
}

// Fill 用给定的单元格填充整个缓冲区。
func (b *Buffer) Fill(c *Cell) {
	b.FillRect(c, b.Bounds())
}

// Clear 用空格单元格清除整个缓冲区。
func (b *Buffer) Clear() {
	b.ClearRect(b.Bounds())
}

// ClearRect 在指定矩形内用空格单元格清除缓冲区。只有矩形边界内的单元格会受到影响。
func (b *Buffer) ClearRect(rect Rectangle) {
	b.FillRect(nil, rect)
}

// InsertLine 在给定行位置插入 n 行，使用给定的可选单元格，在指定矩形内。
// 如果未指定矩形，则在整个缓冲区中插入行。只有矩形水平边界内的单元格会受到影响。
// 行被推出矩形边界并丢失。这遵循终端 [ansi.IL] 行为。
// 返回被推出的行。
func (b *Buffer) InsertLine(y, n int, c *Cell) {
	b.InsertLineRect(y, n, c, b.Bounds())
}

// InsertLineRect 在给定行位置插入新行，使用给定的可选单元格，在矩形边界内。
// 只有矩形水平边界内的单元格会受到影响。行被推出矩形边界并丢失。
// 这遵循终端 [ansi.IL] 行为。
func (b *Buffer) InsertLineRect(y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// 限制要插入的行数到可用空间
	if y+n > rect.Max.Y {
		n = rect.Max.Y - y
	}

	// 在边界内将现有行向下移动
	for i := rect.Max.Y - 1; i >= y+n; i-- {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// 我们不需要在这里克隆 c，因为我们只是将行向下移动。
			b.setCell(x, i, b.Lines[i-n][x], false)
		}
	}

	// 清除边界内新插入的行
	for i := y; i < y+n; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.setCell(x, i, c, true)
		}
	}
}

// DeleteLineRect 在给定行位置删除行，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响。行在边界内向上移动，底部创建新的空白行。
// 这遵循终端 [ansi.DL] 行为。
func (b *Buffer) DeleteLineRect(y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// 限制删除计数到滚动区域的可用空间
	if n > rect.Max.Y-y {
		n = rect.Max.Y - y
	}

	// 在边界内将单元格向上移动
	for dst := y; dst < rect.Max.Y-n; dst++ {
		src := dst + n
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// 我们不需要在这里克隆 c，因为我们只是将单元格向上移动。
			// b.lines[dst][x] = b.lines[src][x]
			b.setCell(x, dst, b.Lines[src][x], false)
		}
	}

	// 用空白单元格填充底部 n 行
	for i := rect.Max.Y - n; i < rect.Max.Y; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.setCell(x, i, c, true)
		}
	}
}

// DeleteLine 在给定行位置删除 n 行，使用给定的可选单元格，在指定矩形内。
// 如果未指定矩形，则在整个缓冲区中删除行。
func (b *Buffer) DeleteLine(y, n int, c *Cell) {
	b.DeleteLineRect(y, n, c, b.Bounds())
}

// InsertCell 在给定位置插入新单元格，使用给定的可选单元格，在指定矩形内。
// 如果未指定矩形，则在整个缓冲区中插入单元格。这遵循终端 [ansi.ICH] 行为。
func (b *Buffer) InsertCell(x, y, n int, c *Cell) {
	b.InsertCellRect(x, y, n, c, b.Bounds())
}

// InsertCellRect 在给定位置插入新单元格，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响，遵循终端 [ansi.ICH] 行为。
func (b *Buffer) InsertCellRect(x, y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// 限制要插入的单元格数到可用空间
	if x+n > rect.Max.X {
		n = rect.Max.X - x
	}

	// 在矩形边界内将现有单元格向右移动
	for i := rect.Max.X - 1; i >= x+n && i-n >= rect.Min.X; i-- {
		// 我们不需要在这里克隆 c，因为我们只是将单元格向右移动。
		// b.lines[y][i] = b.lines[y][i-n]
		b.setCell(i, y, b.Lines[y][i-n], false)
	}

	// 清除矩形边界内新插入的单元格
	for i := x; i < x+n && i < rect.Max.X; i++ {
		b.setCell(i, y, c, true)
	}
}

// DeleteCell 在给定位置删除单元格，使用给定的可选单元格，在指定矩形内。
// 如果未指定矩形，则在整个缓冲区中删除单元格。这遵循终端 [ansi.DCH] 行为。
func (b *Buffer) DeleteCell(x, y, n int, c *Cell) {
	b.DeleteCellRect(x, y, n, c, b.Bounds())
}

// DeleteCellRect 在给定位置删除单元格，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响，遵循终端 [ansi.DCH] 行为。
func (b *Buffer) DeleteCellRect(x, y, n int, c *Cell, rect Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// 计算我们实际可以删除的位置数
	remainingCells := rect.Max.X - x
	if n > remainingCells {
		n = remainingCells
	}

	// 将剩余的单元格向左移动
	for i := x; i < rect.Max.X-n; i++ {
		if i+n < rect.Max.X {
			// 我们不需要在这里克隆 c，因为我们只是将单元格向左移动。
			// b.lines[y][i] = b.lines[y][i+n]
			b.setCell(i, y, b.Lines[y][i+n], false)
		}
	}

	// 用给定的单元格填充空出的位置
	for i := rect.Max.X - n; i < rect.Max.X; i++ {
		b.setCell(i, y, c, true)
	}
}
