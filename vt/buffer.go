package vt

import uv "github.com/charmbracelet/ultraviolet"

// Buffer 是一个终端单元格缓冲区。
type Buffer struct {
	uv.Buffer
}

// InsertLine 在指定的行位置插入 n 行，使用给定的可选单元格，在指定的矩形范围内。
// 如果未指定矩形，则在整个缓冲区中插入行。只有矩形水平边界内的单元格会受到影响。
// 行被推出矩形边界并丢失。这遵循终端 [ansi.IL] 行为。
// 返回被推出的行。
func (b *Buffer) InsertLine(y, n int, c *uv.Cell) {
	b.InsertLineRect(y, n, c, b.Bounds())
}

// InsertLineRect 在指定的行位置插入新行，使用给定的可选单元格，在矩形边界内。
// 只有矩形水平边界内的单元格会受到影响。行被推出矩形边界并丢失。
// 这遵循终端 [ansi.IL] 行为。
func (b *Buffer) InsertLineRect(y, n int, c *uv.Cell, rect uv.Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// 限制插入的行数为可用空间
	if y+n > rect.Max.Y {
		n = rect.Max.Y - y
	}

	// 在边界内将现有行向下移动
	for i := rect.Max.Y - 1; i >= y+n; i-- {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// 这里不需要克隆 c，因为我们只是将行向下移动。
			b.Lines[i][x] = b.Lines[i-n][x]
		}
	}

	// 清除边界内新插入的行
	for i := y; i < y+n; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.SetCell(x, i, c)
		}
	}
}

// DeleteLineRect 在指定的行位置删除 n 行，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响。行在边界内向上移动，并在底部创建新的空行。
// 这遵循终端 [ansi.DL] 行为。
func (b *Buffer) DeleteLineRect(y, n int, c *uv.Cell, rect uv.Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() {
		return
	}

	// 限制删除数量为滚动区域中的可用空间
	if n > rect.Max.Y-y {
		n = rect.Max.Y - y
	}

	// 在边界内将单元格向上移动
	for dst := y; dst < rect.Max.Y-n; dst++ {
		src := dst + n
		for x := rect.Min.X; x < rect.Max.X; x++ {
			// 这里不需要克隆 c，因为我们只是将单元格向上移动。
			b.Lines[dst][x] = b.Lines[src][x]
		}
	}

	// 用空单元格填充底部的 n 行
	for i := rect.Max.Y - n; i < rect.Max.Y; i++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.SetCell(x, i, c)
		}
	}
}

// DeleteLine 在指定的行位置删除 n 行，使用给定的可选单元格，在指定的矩形范围内。
// 如果未指定矩形，则在整个缓冲区中删除行。
func (b *Buffer) DeleteLine(y, n int, c *uv.Cell) {
	b.DeleteLineRect(y, n, c, b.Bounds())
}

// InsertCell 在指定的位置插入新单元格，使用给定的可选单元格，在指定的矩形范围内。
// 如果未指定矩形，则在整个缓冲区中插入单元格。这遵循终端 [ansi.ICH] 行为。
func (b *Buffer) InsertCell(x, y, n int, c *uv.Cell) {
	b.InsertCellRect(x, y, n, c, b.Bounds())
}

// InsertCellRect 在指定的位置插入新单元格，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响，遵循终端 [ansi.ICH] 行为。
func (b *Buffer) InsertCellRect(x, y, n int, c *uv.Cell, rect uv.Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// 限制插入的单元格数量为可用空间
	if x+n > rect.Max.X {
		n = rect.Max.X - x
	}

	// 在矩形边界内将现有单元格向右移动
	for i := rect.Max.X - 1; i >= x+n && i-n >= rect.Min.X; i-- {
		// 这里不需要克隆 c，因为我们只是将单元格向右移动。
		// b.lines[y][i] = b.lines[y][i-n]
		b.Lines[y][i] = b.Lines[y][i-n]
	}

	// 清除矩形边界内新插入的单元格
	for i := x; i < x+n && i < rect.Max.X; i++ {
		b.SetCell(i, y, c)
	}
}

// DeleteCell 在指定的位置删除单元格，使用给定的可选单元格，在指定的矩形范围内。
// 如果未指定矩形，则在整个缓冲区中删除单元格。这遵循终端 [ansi.DCH] 行为。
func (b *Buffer) DeleteCell(x, y, n int, c *uv.Cell) {
	b.DeleteCellRect(x, y, n, c, b.Bounds())
}

// DeleteCellRect 在指定的位置删除单元格，使用给定的可选单元格，在矩形边界内。
// 只有矩形边界内的单元格会受到影响，遵循终端 [ansi.DCH] 行为。
func (b *Buffer) DeleteCellRect(x, y, n int, c *uv.Cell, rect uv.Rectangle) {
	if n <= 0 || y < rect.Min.Y || y >= rect.Max.Y || y >= b.Height() ||
		x < rect.Min.X || x >= rect.Max.X || x >= b.Width() {
		return
	}

	// 计算我们实际可以删除的位置数量
	remainingCells := rect.Max.X - x
	if n > remainingCells {
		n = remainingCells
	}

	// 将剩余的单元格向左移动
	for i := x; i < rect.Max.X-n; i++ {
		if i+n < rect.Max.X {
			// 这里不需要克隆 c，因为我们只是将单元格向左移动。
			// b.lines[y][i] = b.lines[y][i+n]
			b.Lines[y][i] = b.Lines[y][i+n]
		}
	}

	// 用给定的单元格填充空出的位置
	for i := rect.Max.X - n; i < rect.Max.X; i++ {
		b.SetCell(i, y, c)
	}
}
