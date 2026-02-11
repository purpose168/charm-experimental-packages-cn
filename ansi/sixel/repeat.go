package sixel

import (
	"fmt"
	"io"
	"strings"
)

// ErrInvalidRepeat 在重复无效时返回。
var ErrInvalidRepeat = fmt.Errorf("无效重复")

// WriteRepeat 向写入器写入重复。重复字符的范围是
// '?' (0x3F) 到 '~' (0x7E)。
func WriteRepeat(w io.Writer, count int, char byte) (int, error) {
	return fmt.Fprintf(w, "%c%d%c", RepeatIntroducer, count, char) //nolint:wrapcheck
}

// Repeat 表示六像素重复引导符。
type Repeat struct {
	Count int
	Char  byte
}

// WriteTo 向写入器写入重复。
func (r Repeat) WriteTo(w io.Writer) (int64, error) {
	n, err := WriteRepeat(w, r.Count, r.Char)
	return int64(n), err
}

// String 将重复返回为字符串。
func (r Repeat) String() string {
	var b strings.Builder
	r.WriteTo(&b) //nolint:errcheck,gosec
	return b.String()
}

// DecodeRepeat 从字节切片解码重复。它返回重复和
// 读取的字节数。
func DecodeRepeat(data []byte) (r Repeat, n int) {
	if len(data) == 0 || data[0] != RepeatIntroducer {
		return r, n
	}

	if len(data) < 3 { // 最小长度是 3：引导符、数字和字符。
		return r, n
	}

	for n = 1; n < len(data); n++ {
		if data[n] >= '0' && data[n] <= '9' {
			r.Count = r.Count*10 + int(data[n]-'0')
		} else {
			r.Char = data[n]
			n++ // 在计数中包含该字符。
			break
		}
	}

	return r, n
}
