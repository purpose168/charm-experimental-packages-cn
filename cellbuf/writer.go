package cellbuf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// CellBuffer 是一个单元格缓冲区，代表屏幕或网格中的一组单元格。
type CellBuffer interface {
	// Cell 返回指定位置的单元格。
	Cell(x, y int) *Cell
	// SetCell 将指定位置的单元格设置为给定的单元格。
	// 返回单元格是否设置成功。
	SetCell(x, y int, c *Cell) bool
	// Bounds 返回单元格缓冲区的边界。
	Bounds() Rectangle
}

// FillRect 用给定的单元格填充单元格缓冲区中的矩形区域。
// 这不会填充单元格缓冲区边界外的单元格。
func FillRect(s CellBuffer, c *Cell, rect Rectangle) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			s.SetCell(x, y, c)
		}
	}
}

// Fill 用给定的单元格填充整个单元格缓冲区。
func Fill(s CellBuffer, c *Cell) {
	FillRect(s, c, s.Bounds())
}

// ClearRect 用空白单元格清除单元格缓冲区中的矩形区域。
func ClearRect(s CellBuffer, rect Rectangle) {
	FillRect(s, nil, rect)
}

// Clear 用空白单元格清除整个单元格缓冲区。
func Clear(s CellBuffer) {
	Fill(s, nil)
}

// SetContentRect 用空白单元格清除单元格缓冲区中的矩形区域，
// 并将给定的字符串设置为其内容。如果字符串的高度或宽度超过
// 单元格缓冲区的高度或宽度，它将被截断。
func SetContentRect(s CellBuffer, str string, rect Rectangle) {
	// 将所有 "\n" 替换为 "\r\n" 以确保光标重置到行首。
	// 确保我们不会将 "\r\n" 替换为 "\r\r\n"。
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\n", "\r\n")
	ClearRect(s, rect)
	printString(s, ansi.GraphemeWidth, rect.Min.X, rect.Min.Y, rect, str, true, "")
}

// SetContent 用空白单元格清除单元格缓冲区，并将给定的字符串
// 设置为其内容。如果字符串的高度或宽度超过单元格缓冲区的
// 高度或宽度，它将被截断。
func SetContent(s CellBuffer, str string) {
	SetContentRect(s, str, s.Bounds())
}

// Render 返回带有 ANSI 转义序列的网格字符串表示。
func Render(d CellBuffer) string {
	var buf bytes.Buffer
	height := d.Bounds().Dy()
	for y := range height {
		_, line := RenderLine(d, y)
		buf.WriteString(line)
		if y < height-1 {
			buf.WriteString("\r\n")
		}
	}
	return buf.String()
}

// RenderLine 返回网格第 y 行的字符串表示及其宽度。
func RenderLine(d CellBuffer, n int) (w int, line string) {
	var pen Style
	var link Link
	var buf bytes.Buffer
	var pendingLine string
	var pendingWidth int // 忽略空格单元格，直到遇到非空格单元格

	writePending := func() {
		// 如果没有待处理的行，我们不需要做任何事情。
		if len(pendingLine) == 0 {
			return
		}
		buf.WriteString(pendingLine)
		w += pendingWidth
		pendingWidth = 0
		pendingLine = ""
	}

	for x := range d.Bounds().Dx() {
		if cell := d.Cell(x, n); cell != nil && cell.Width > 0 { //nolint:nestif
			// 将单元格的样式和链接转换为给定的颜色配置文件。
			cellStyle := cell.Style
			cellLink := cell.Link
			if cellStyle.Empty() && !pen.Empty() {
				writePending()
				buf.WriteString(ansi.ResetStyle)
				pen.Reset()
			}
			if !cellStyle.Equal(&pen) {
				writePending()
				seq := cellStyle.DiffSequence(pen)
				buf.WriteString(seq)
				pen = cellStyle
			}

			// 写入 URL 转义序列
			if cellLink != link && link.URL != "" {
				writePending()
				buf.WriteString(ansi.ResetHyperlink())
				link.Reset()
			}
			if cellLink != link {
				writePending()
				buf.WriteString(ansi.SetHyperlink(cellLink.URL, cellLink.Params))
				link = cellLink
			}

			// 我们只在单元格内容不为空时写入。如果为空，
			// 将其添加到待处理行和宽度中，以便稍后评估。
			if cell.Equal(&BlankCell) {
				pendingLine += cell.String()
				pendingWidth += cell.Width
			} else {
				writePending()
				buf.WriteString(cell.String())
				w += cell.Width
			}
		}
	}
	if link.URL != "" {
		buf.WriteString(ansi.ResetHyperlink())
	}
	if !pen.Empty() {
		buf.WriteString(ansi.ResetStyle)
	}
	return w, strings.TrimRight(buf.String(), " ") // 修剪尾部空格
}

// ScreenWriter 表示一个写入到 [Screen] 的写入器，解析 ANSI
// 转义序列和 Unicode 字符，并将它们转换为可以写入到单元格
// [Buffer] 的单元格。
type ScreenWriter struct {
	*Screen
}

// NewScreenWriter 创建一个新的 ScreenWriter，写入到给定的 Screen。
// 这是创建 ScreenWriter 的便捷函数。
func NewScreenWriter(s *Screen) *ScreenWriter {
	return &ScreenWriter{s}
}

// Write 将给定的字节写入屏幕。
// 这将识别 ANSI [ansi.SGR] 样式和 [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) Write(p []byte) (n int, err error) {
	printString(s.Screen, s.method,
		s.cur.X, s.cur.Y, s.Bounds(),
		p, false, "")
	return len(p), nil
}

// SetContent 用空白单元格清除屏幕，并将给定的字符串设置为
// 其内容。如果字符串的高度或宽度超过屏幕的高度或宽度，
// 它将被截断。
//
// 这将识别 ANSI [ansi.SGR] 样式和 [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) SetContent(str string) {
	s.SetContentRect(str, s.Bounds())
}

// SetContentRect 用空白单元格清除屏幕中的矩形区域，并
// 将给定的字符串设置为其内容。如果字符串的高度或宽度
// 超过屏幕的高度或宽度，它将被截断。
//
// 这将识别 ANSI [ansi.SGR] 样式和 [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) SetContentRect(str string, rect Rectangle) {
	// 将所有 "\n" 替换为 "\r\n" 以确保光标重置到行首。
	// 确保我们不会将 "\r\n" 替换为 "\r\r\n"。
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\n", "\r\n")
	s.ClearRect(rect)
	printString(s.Screen, s.method,
		rect.Min.X, rect.Min.Y, rect,
		str, true, "")
}

// Print 在当前光标位置打印字符串。如果字符串超过屏幕宽度，
// 它将换行到屏幕宽度。这将识别 ANSI [ansi.SGR] 样式和
// [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) Print(str string, v ...any) {
	if len(v) > 0 {
		str = fmt.Sprintf(str, v...)
	}
	printString(s.Screen, s.method,
		s.cur.X, s.cur.Y, s.Bounds(),
		str, false, "")
}

// PrintAt 在给定位置打印字符串。如果字符串超过屏幕宽度，
// 它将换行到屏幕宽度。这将识别 ANSI [ansi.SGR] 样式和
// [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) PrintAt(x, y int, str string, v ...any) {
	if len(v) > 0 {
		str = fmt.Sprintf(str, v...)
	}
	printString(s.Screen, s.method,
		x, y, s.Bounds(),
		str, false, "")
}

// PrintCrop 在当前光标位置打印字符串，如果文本超过屏幕宽度则截断。
// 如果字符串被截断，使用 tail 指定要追加的字符串。
// 这将识别 ANSI [ansi.SGR] 样式和 [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) PrintCrop(str string, tail string) {
	printString(s.Screen, s.method,
		s.cur.X, s.cur.Y, s.Bounds(),
		str, true, tail)
}

// PrintCropAt 在给定位置打印字符串，如果文本超过屏幕宽度则截断。
// 如果字符串被截断，使用 tail 指定要追加的字符串。
// 这将识别 ANSI [ansi.SGR] 样式和 [ansi.SetHyperlink] 转义序列。
func (s *ScreenWriter) PrintCropAt(x, y int, str string, tail string) {
	printString(s.Screen, s.method,
		x, y, s.Bounds(),
		str, true, tail)
}

// printString 从给定位置开始绘制字符串。
func printString[T []byte | string](
	s CellBuffer,
	m ansi.Method,
	x, y int,
	bounds Rectangle, str T,
	truncate bool, tail string,
) {
	p := ansi.GetParser()
	defer ansi.PutParser(p)

	var tailc Cell
	if truncate && len(tail) > 0 {
		if m == ansi.WcWidth {
			tailc = *NewCellString(tail)
		} else {
			tailc = *NewGraphemeCell(tail)
		}
	}

	decoder := ansi.DecodeSequenceWc[T]
	if m == ansi.GraphemeWidth {
		decoder = ansi.DecodeSequence[T]
	}

	var cell Cell
	var style Style
	var link Link
	var state byte
	for len(str) > 0 {
		seq, width, n, newState := decoder(str, state, p)

		switch width {
		case 1, 2, 3, 4: // 宽单元格最多可以宽达 4 个单元格
			cell.Width += width
			cell.Append([]rune(string(seq))...)

			if !truncate && x+cell.Width > bounds.Max.X && y+1 < bounds.Max.Y {
				// 将字符串换行到窗口宽度
				x = bounds.Min.X
				y++
			}
			if Pos(x, y).In(bounds) {
				if truncate && tailc.Width > 0 && x+cell.Width > bounds.Max.X-tailc.Width {
					// 截断字符串并在需要时追加尾部
					cell := tailc
					cell.Style = style
					cell.Link = link
					s.SetCell(x, y, &cell)
					x += tailc.Width
				} else {
					// 将单元格打印到屏幕
					cell.Style = style
					cell.Link = link
					s.SetCell(x, y, &cell)
					x += width
				}
			}

			// 字符串太长，超出行宽，截断它。
			// 确保我们为下一次迭代重置单元格。
			cell.Reset()
		default:
			// 有效的序列总是有非零的 Cmd。
			//nolint:godox
			// TODO: 处理光标移动和其他序列
			switch {
			case ansi.HasCsiPrefix(seq) && p.Command() == 'm':
				// SGR - 选择图形渲染
				ReadStyle(p.Params(), &style)
			case ansi.HasOscPrefix(seq) && p.Command() == 8:
				// 超链接
				ReadLink(p.Data(), &link)
			case ansi.Equal(seq, T("\n")):
				y++
			case ansi.Equal(seq, T("\r")):
				x = bounds.Min.X
			default:
				cell.Append([]rune(string(seq))...)
			}
		}

		// 推进状态和数据
		state = newState
		str = str[n:]
	}

	// 如果最后一个单元格不为空，确保设置它。
	if !cell.Empty() {
		s.SetCell(x, y, &cell)
		cell.Reset()
	}
}
