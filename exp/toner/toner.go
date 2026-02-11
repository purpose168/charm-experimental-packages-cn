package toner

import (
	"bytes"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/exp/charmtone"
)

// Strings 返回输入字节切片或字符串的彩色字符串表示。
func Strings(s string) string {
	var buf bytes.Buffer
	_, _ = writeColorize(&buf, s)
	return buf.String()
}

// Bytes 返回输入字节切片的彩色字节切片表示。
func Bytes(b []byte) []byte {
	var buf bytes.Buffer
	_, _ = writeColorize(&buf, b)
	return buf.Bytes()
}

// Writer 封装了一个 [io.Writer]，使用 charm 色调对输出进行着色。
type Writer struct {
	io.Writer
}

// Write 将彩色输出写入到底层写入器。
func (w Writer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	return writeColorize(w, p)
}

// WriteString 将输入字符串的彩色输出写入到底层写入器。
func (w Writer) WriteString(s string) (n int, err error) {
	if len(s) == 0 {
		return 0, nil
	}

	return writeColorize(w, s)
}

var colors = func() []charmtone.Key {
	cols := charmtone.Keys()
	// Filter out these colors.
	filterOut := []charmtone.Key{
		charmtone.Cumin,
		charmtone.Tang,
		charmtone.Paprika,
		charmtone.Pepper,
		charmtone.Charcoal,
		charmtone.Iron,
		charmtone.Oyster,
		charmtone.Squid,
		charmtone.Smoke,
		charmtone.Ash,
		charmtone.Salt,
		charmtone.Butter,
	}
	for _, k := range filterOut {
		cols = slices.DeleteFunc(cols, func(c charmtone.Key) bool {
			return c == k
		})
	}
	return cols
}()

func writeColorize[T []byte | string](w io.Writer, p T) (n int, err error) {
	pa := ansi.NewParser()

	var state byte
	for len(p) > 0 {
		seq, width, nr, newState := ansi.DecodeSequence(p, state, pa)
		cmd := pa.Command()

		var st ansi.Style
		var s string
		if cmd == 0 && width > 0 {
			s = string(seq)
		} else {
			st = st.ForegroundColor(colors[cmd%len(colors)])
			s = string(seq)
		}

		s = strconv.Quote(s)
		s = strings.TrimPrefix(s, "\"")
		s = strings.TrimSuffix(s, "\"")
		if len(st) > 0 {
			s = st.Styled(s)
		}

		m, err := io.WriteString(w, s)
		if err != nil {
			return n, err
		}

		n += m

		p = p[nr:]
		state = newState
	}

	return n, nil
}
