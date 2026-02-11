package cellbuf

import (
	"io"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// PenWriter 是一个写入器，用于写入缓冲区并跟踪当前的画笔样式和链接状态，以便使用换行符进行包装。
type PenWriter struct {
	w     io.Writer
	p     *ansi.Parser
	style Style
	link  Link
}

// NewPenWriter 返回一个新的 PenWriter。
func NewPenWriter(w io.Writer) *PenWriter {
	pw := &PenWriter{w: w}
	pw.p = ansi.GetParser()
	handleCsi := func(cmd ansi.Cmd, params ansi.Params) {
		if cmd == 'm' {
			ReadStyle(params, &pw.style)
		}
	}
	handleOsc := func(cmd int, data []byte) {
		if cmd == 8 {
			ReadLink(data, &pw.link)
		}
	}
	pw.p.SetHandler(ansi.Handler{
		HandleCsi: handleCsi,
		HandleOsc: handleOsc,
	})
	return pw
}

// Style 返回当前的画笔样式。
func (w *PenWriter) Style() Style {
	return w.style
}

// Link 返回当前的画笔链接。
func (w *PenWriter) Link() Link {
	return w.link
}

// Write 写入到缓冲区。
func (w *PenWriter) Write(p []byte) (int, error) {
	for i := range p {
		b := p[i]
		w.p.Advance(b)
		if b == '\n' {
			if !w.style.Empty() {
				_, _ = w.w.Write([]byte(ansi.ResetStyle))
			}
			if !w.link.Empty() {
				_, _ = w.w.Write([]byte(ansi.ResetHyperlink()))
			}
		}

		_, _ = w.w.Write([]byte{b})
		if b == '\n' {
			if !w.link.Empty() {
				_, _ = w.w.Write([]byte(ansi.SetHyperlink(w.link.URL, w.link.Params)))
			}
			if !w.style.Empty() {
				_, _ = w.w.Write([]byte(w.style.Sequence()))
			}
		}
	}

	return len(p), nil
}

// Close 关闭写入器，必要时重置样式和链接，并释放其解析器。调用它对性能至关重要，但忘记调用不会导致安全问题或内存泄漏。
func (w *PenWriter) Close() error {
	if !w.style.Empty() {
		_, _ = w.w.Write([]byte(ansi.ResetStyle))
	}
	if !w.link.Empty() {
		_, _ = w.w.Write([]byte(ansi.ResetHyperlink()))
	}
	if w.p != nil {
		ansi.PutParser(w.p)
		w.p = nil
	}
	return nil
}
