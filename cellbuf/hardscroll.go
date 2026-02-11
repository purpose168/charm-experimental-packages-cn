package cellbuf

import (
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// scrollOptimize 优化屏幕以将旧缓冲区转换为新缓冲区。
func (s *Screen) scrollOptimize() {
	height := s.newbuf.Height()
	if s.oldnum == nil || len(s.oldnum) < height {
		s.oldnum = make([]int, height)
	}

	// Calculate the indices
	s.updateHashmap()
	if len(s.hashtab) < height {
		return
	}

	// 第一轮 - 从上到下向上滚动
	for i := 0; i < height; {
		for i < height && (s.oldnum[i] == newIndex || s.oldnum[i] <= i) {
			i++
		}
		if i >= height {
			break
		}

		shift := s.oldnum[i] - i // shift > 0
		start := i

		i++
		for i < height && s.oldnum[i] != newIndex && s.oldnum[i]-i == shift {
			i++
		}
		end := i - 1 + shift

		if !s.scrolln(shift, start, end, height-1) {
			continue
		}
	}

	// 第二轮 - 从下到上向下滚动
	for i := height - 1; i >= 0; {
		for i >= 0 && (s.oldnum[i] == newIndex || s.oldnum[i] >= i) {
			i--
		}
		if i < 0 {
			break
		}

		shift := s.oldnum[i] - i // shift < 0
		end := i

		i--
		for i >= 0 && s.oldnum[i] != newIndex && s.oldnum[i]-i == shift {
			i--
		}

		start := i + 1 - (-shift)
		if !s.scrolln(shift, start, end, height-1) {
			continue
		}
	}
}

// scrolln 向上滚动屏幕 n 行。
func (s *Screen) scrolln(n, top, bot, maxY int) (v bool) { //nolint:unparam
	const (
		nonDestScrollRegion = false
		memoryBelow         = false
	)

	blank := s.clearBlank()
	if n > 0 { //nolint:nestif
		// 向上滚动（向前）
		v = s.scrollUp(n, top, bot, 0, maxY, blank)
		if !v {
			s.buf.WriteString(ansi.SetTopBottomMargins(top+1, bot+1))

			// XXX: 在不使用替代屏幕的内联模式下，我们应该如何处理这个问题？
			s.cur.X, s.cur.Y = -1, -1
			v = s.scrollUp(n, top, bot, top, bot, blank)
			s.buf.WriteString(ansi.SetTopBottomMargins(1, maxY+1))
			s.cur.X, s.cur.Y = -1, -1
		}

		if !v {
			v = s.scrollIdl(n, top, bot-n+1, blank)
		}

		// 清除新移入的行。
		if v &&
			(nonDestScrollRegion || (memoryBelow && bot == maxY)) {
			if bot == maxY {
				s.move(0, bot-n+1)
				s.clearToBottom(nil)
			} else {
				for i := range n {
					s.move(0, bot-i)
					s.clearToEnd(nil, false)
				}
			}
		}
	} else if n < 0 {
		// 向下滚动（向后）
		v = s.scrollDown(-n, top, bot, 0, maxY, blank)
		if !v {
			s.buf.WriteString(ansi.SetTopBottomMargins(top+1, bot+1))

			// XXX: 在不使用替代屏幕的内联模式下，我们应该如何处理这个问题？
			s.cur.X, s.cur.Y = -1, -1
			v = s.scrollDown(-n, top, bot, top, bot, blank)
			s.buf.WriteString(ansi.SetTopBottomMargins(1, maxY+1))
			s.cur.X, s.cur.Y = -1, -1

			if !v {
				v = s.scrollIdl(-n, bot+n+1, top, blank)
			}

			// 清除新移入的行。
			if v &&
				(nonDestScrollRegion || (memoryBelow && top == 0)) {
				for i := range -n {
					s.move(0, top+i)
					s.clearToEnd(nil, false)
				}
			}
		}
	}

	if !v {
		return v
	}

	s.scrollBuffer(s.curbuf, n, top, bot, blank)

	// 也移动哈希值，它们可以被重用
	s.scrollOldhash(n, top, bot)

	return true
}

// scrollBuffer 滚动缓冲区 n 行。
func (s *Screen) scrollBuffer(b *Buffer, n, top, bot int, blank *Cell) {
	if top < 0 || bot < top || bot >= b.Height() {
		// 没有内容需要滚动
		return
	}

	if n < 0 {
		// 向下移动 n 行
		limit := top - n
		for line := bot; line >= limit && line >= 0 && line >= top; line-- {
			copy(b.Lines[line], b.Lines[line+n])
		}
		for line := top; line < limit && line <= b.Height()-1 && line <= bot; line++ {
			b.FillRect(blank, Rect(0, line, b.Width(), 1))
		}
	}

	if n > 0 {
		// 向上移动 n 行
		limit := bot - n
		for line := top; line <= limit && line <= b.Height()-1 && line <= bot; line++ {
			copy(b.Lines[line], b.Lines[line+n])
		}
		for line := bot; line > limit && line >= 0 && line >= top; line-- {
			b.FillRect(blank, Rect(0, line, b.Width(), 1))
		}
	}

	s.touchLine(b.Width(), b.Height(), top, bot-top+1, true)
}

// touchLine 将行标记为已触摸。
func (s *Screen) touchLine(width, height, y, n int, changed bool) {
	if n < 0 || y < 0 || y >= height {
		return // 没有内容需要触摸
	}

	for i := y; i < y+n && i < height; i++ {
		if changed {
			s.touch[i] = lineData{firstCell: 0, lastCell: width - 1}
		} else {
			delete(s.touch, i)
		}
	}
}

// scrollUp 向上滚动屏幕 n 行。
func (s *Screen) scrollUp(n, top, bot, minY, maxY int, blank *Cell) bool {
	if n == 1 && top == minY && bot == maxY { //nolint:nestif
		s.move(0, bot)
		s.updatePen(blank)
		s.buf.WriteByte('\n')
	} else if n == 1 && bot == maxY {
		s.move(0, top)
		s.updatePen(blank)
		s.buf.WriteString(ansi.DeleteLine(1))
	} else if top == minY && bot == maxY {
		supportsSU := s.caps.Contains(capSU)
		if supportsSU {
			s.move(0, bot)
		} else {
			s.move(0, top)
		}
		s.updatePen(blank)
		if supportsSU {
			s.buf.WriteString(ansi.ScrollUp(n))
		} else {
			s.buf.WriteString(strings.Repeat("\n", n))
		}
	} else if bot == maxY {
		s.move(0, top)
		s.updatePen(blank)
		s.buf.WriteString(ansi.DeleteLine(n))
	} else {
		return false
	}
	return true
}

// scrollDown 向下滚动屏幕 n 行。
func (s *Screen) scrollDown(n, top, bot, minY, maxY int, blank *Cell) bool {
	if n == 1 && top == minY && bot == maxY { //nolint:nestif
		s.move(0, top)
		s.updatePen(blank)
		s.buf.WriteString(ansi.ReverseIndex)
	} else if n == 1 && bot == maxY {
		s.move(0, top)
		s.updatePen(blank)
		s.buf.WriteString(ansi.InsertLine(1))
	} else if top == minY && bot == maxY {
		s.move(0, top)
		s.updatePen(blank)
		if s.caps.Contains(capSD) {
			s.buf.WriteString(ansi.ScrollDown(n))
		} else {
			s.buf.WriteString(strings.Repeat(ansi.ReverseIndex, n))
		}
	} else if bot == maxY {
		s.move(0, top)
		s.updatePen(blank)
		s.buf.WriteString(ansi.InsertLine(n))
	} else {
		return false
	}
	return true
}

// scrollIdl 通过在 del 处使用 [ansi.DL] 和在 ins 处使用 [ansi.IL] 来滚动屏幕 n 行。
func (s *Screen) scrollIdl(n, del, ins int, blank *Cell) bool {
	if n < 0 {
		return false
	}

	// Delete lines
	s.move(0, del)
	s.updatePen(blank)
	s.buf.WriteString(ansi.DeleteLine(n))

	// Insert lines
	s.move(0, ins)
	s.updatePen(blank)
	s.buf.WriteString(ansi.InsertLine(n))

	return true
}
