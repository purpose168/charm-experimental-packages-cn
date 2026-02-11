// Package kitty 提供 Kitty 终端图形协议功能。
package kitty

import (
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// Decoder 是 Kitty 图形协议的解码器。它支持解码 24 位 [RGB]、32 位 [RGBA] 和 [PNG] 格式的图像。
// 它还可以使用 zlib 解压缩数据。
// 默认格式是 32 位 [RGBA]。
type Decoder struct {
	// 使用 zlib 解压缩。
	Decompress bool

	// 可以是 [RGB]、[RGBA] 或 [PNG] 之一。
	Format int

	// 图像宽度（以像素为单位）。如果图像是 [PNG] 格式，则可以省略。
	Width int

	// 图像高度（以像素为单位）。如果图像是 [PNG] 格式，则可以省略。
	Height int
}

// Decode 从 r 中以指定格式解码图像数据。
func (d *Decoder) Decode(r io.Reader) (image.Image, error) {
	if d.Decompress {
		zr, err := zlib.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("failed to create zlib reader: %w", err)
		}

		defer zr.Close() //nolint:errcheck
		r = zr
	}

	if d.Format == 0 {
		d.Format = RGBA
	}

	switch d.Format {
	case RGBA, RGB:
		return d.decodeRGBA(r, d.Format == RGBA)

	case PNG:
		return png.Decode(r) //nolint:wrapcheck

	default:
		return nil, fmt.Errorf("unsupported format: %d", d.Format)
	}
}

// decodeRGBA 解码 32 位 RGBA 或 24 位 RGB 格式的图像数据。
func (d *Decoder) decodeRGBA(r io.Reader, alpha bool) (image.Image, error) {
	m := image.NewRGBA(image.Rect(0, 0, d.Width, d.Height))

	var buf []byte
	if alpha {
		buf = make([]byte, 4)
	} else {
		buf = make([]byte, 3)
	}

	for y := range d.Height {
		for x := range d.Width {
			if _, err := io.ReadFull(r, buf[:]); err != nil {
				return nil, fmt.Errorf("failed to read pixel data: %w", err)
			}
			if alpha {
				m.SetRGBA(x, y, color.RGBA{buf[0], buf[1], buf[2], buf[3]})
			} else {
				m.SetRGBA(x, y, color.RGBA{buf[0], buf[1], buf[2], 0xff})
			}
		}
	}

	return m, nil
}
