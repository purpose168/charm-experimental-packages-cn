package kitty

import (
	"compress/zlib"
	"fmt"
	"image"
	"image/png"
	"io"
)

// Encoder 是 Kitty 图形协议的编码器。它支持将图像编码为 24 位 [RGB]、32 位 [RGBA] 和 [PNG] 格式，
// 并使用 zlib 压缩数据。
// 默认格式是 32 位 [RGBA]。
type Encoder struct {
	// 使用 zlib 压缩。
	Compress bool

	// 可以是 [RGBA]、[RGB] 或 [PNG] 之一。
	Format int
}

// Encode 将图像数据编码为指定格式并写入 w。
func (e *Encoder) Encode(w io.Writer, m image.Image) error {
	if m == nil {
		return nil
	}

	if e.Compress {
		zw := zlib.NewWriter(w)
		defer zw.Close() //nolint:errcheck
		w = zw
	}

	if e.Format == 0 {
		e.Format = RGBA
	}

	switch e.Format {
	case RGBA, RGB:
		bounds := m.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := m.At(x, y).RGBA()
				switch e.Format {
				case RGBA:
					w.Write([]byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}) //nolint:errcheck,gosec
				case RGB:
					w.Write([]byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)}) //nolint:errcheck,gosec
				}
			}
		}

	case PNG:
		if err := png.Encode(w, m); err != nil {
			return fmt.Errorf("failed to encode PNG: %w", err)
		}

	default:
		return fmt.Errorf("unsupported format: %d", e.Format)
	}

	return nil
}
