package sixel

import (
	"fmt"
	"io"
	"strings"
)

// ErrInvalidRaster 在光栅属性无效时返回。
var ErrInvalidRaster = fmt.Errorf("无效光栅属性")

// WriteRaster 向写入器写入光栅属性。如果 ph 和 pv 为 0，则
// 省略它们。
func WriteRaster(w io.Writer, pan, pad, ph, pv int) (n int, err error) {
	if pad == 0 {
		return WriteRaster(w, 1, 1, ph, pv)
	}

	if ph <= 0 && pv <= 0 {
		return fmt.Fprintf(w, "%c%d;%d", RasterAttribute, pan, pad) //nolint:wrapcheck
	}

	return fmt.Fprintf(w, "%c%d;%d;%d;%d", RasterAttribute, pan, pad, ph, pv) //nolint:wrapcheck
}

// Raster 表示六像素光栅属性。
type Raster struct {
	Pan, Pad, Ph, Pv int
}

// WriteTo 向写入器写入光栅属性。
func (r Raster) WriteTo(w io.Writer) (int64, error) {
	n, err := WriteRaster(w, r.Pan, r.Pad, r.Ph, r.Pv)
	return int64(n), err
}

// String 将光栅返回为字符串。
func (r Raster) String() string {
	var b strings.Builder
	r.WriteTo(&b) //nolint:errcheck,gosec
	return b.String()
}

// DecodeRaster 从字节切片解码光栅。它返回光栅和
// 读取的字节数。
func DecodeRaster(data []byte) (r Raster, n int) {
	if len(data) == 0 || data[0] != RasterAttribute {
		return r, n
	}

	ptr := &r.Pan
	for n = 1; n < len(data); n++ {
		if data[n] == ';' { //nolint:nestif
			if ptr == &r.Pan {
				ptr = &r.Pad
			} else if ptr == &r.Pad {
				ptr = &r.Ph
			} else if ptr == &r.Ph {
				ptr = &r.Pv
			} else {
				n++
				break
			}
		} else if data[n] >= '0' && data[n] <= '9' {
			*ptr = (*ptr)*10 + int(data[n]-'0')
		} else {
			break
		}
	}

	return r, n
}
