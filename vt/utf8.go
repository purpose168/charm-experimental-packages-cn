package vt

import (
	"unicode/utf8"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// handlePrint 处理可打印字符。
func (e *Emulator) handlePrint(r rune) {
	if r >= ansi.SP && r < ansi.DEL {
		if len(e.grapheme) > 0 {
			// 如果我们有字形缓冲区，在处理ASCII字符之前先刷新它。
			e.flushGrapheme()
		}
		e.handleGrapheme(string(r), 1)
	} else {
		e.grapheme = append(e.grapheme, r)
	}
}

// flushGrapheme 刷新当前字形缓冲区（如果有），并将字形作为单个单元处理。
func (e *Emulator) flushGrapheme() {
	if len(e.grapheme) == 0 {
		return
	}

	// XXX: 我们在这里始终使用 [ansi.GraphemeWidth] 来报告准确的宽度，
	// 由调用者决定如何处理 Unicode 与非 Unicode 模式。
	method := ansi.GraphemeWidth
	graphemes := string(e.grapheme)
	for len(graphemes) > 0 {
		cluster, width := ansi.FirstGraphemeCluster(graphemes, method)
		e.handleGrapheme(cluster, width)
		graphemes = graphemes[len(cluster):]
	}
	e.grapheme = e.grapheme[:0] // 重置字形缓冲区。
}

// handleGrapheme 处理 UTF-8 字形。
func (e *Emulator) handleGrapheme(content string, width int) {
	awm := e.isModeSet(ansi.ModeAutoWrap)
	cell := uv.Cell{
		Content: content,
		Width:   width,
		Style:   e.scr.cursorPen(),
		Link:    e.scr.cursorLink(),
	}

	x, y := e.scr.CursorPosition()
	if e.atPhantom && awm {
		// 将光标向下移动，类似于 [Terminal.linefeed]，但不尊重 [ansi.LNM] 模式。
		// 这将重置幻影状态，即待换行状态。
		e.index()
		_, y = e.scr.CursorPosition()
		x = 0
	}

	// 处理字符集映射
	if len(content) == 1 { //nolint:nestif
		var charset CharSet
		c := content[0]
		if e.gsingle > 1 && e.gsingle < 4 {
			charset = e.charsets[e.gsingle]
			e.gsingle = 0
		} else if c < 128 {
			charset = e.charsets[e.gl]
		} else {
			charset = e.charsets[e.gr]
		}

		if charset != nil {
			if r, ok := charset[c]; ok {
				cell.Content = r
				cell.Width = 1
			}
		}
	}

	if cell.Width == 1 && len(content) == 1 {
		e.lastChar, _ = utf8.DecodeRuneInString(content)
	}

	e.scr.SetCell(x, y, &cell)

	// 处理行尾的幻影状态
	e.atPhantom = awm && x >= e.scr.Width()-1
	if !e.atPhantom {
		x += cell.Width
	}

	// 注意：我们不在这里重置幻影状态，而是在上面处理它。
	e.scr.setCursor(x, y, false)
}
