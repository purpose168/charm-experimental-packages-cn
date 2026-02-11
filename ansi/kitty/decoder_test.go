package kitty

import (
	"bytes"
	"compress/zlib"
	"image"
	"image/color"
	"image/png"
	"reflect"
	"testing"
)

// TestDecoder_Decode 测试解码器的解码功能。
func TestDecoder_Decode(t *testing.T) {
	// 辅助函数：创建压缩数据
	compress := func(data []byte) []byte {
		var buf bytes.Buffer
		w := zlib.NewWriter(&buf)
		w.Write(data)
		w.Close()
		return buf.Bytes()
	}

	tests := []struct {
		name    string
		decoder Decoder
		input   []byte
		want    image.Image
		wantErr bool
	}{
		{
			name: "RGBA format 2x2",
			decoder: Decoder{
				Format: RGBA,
				Width:  2,
				Height: 2,
			},
			input: []byte{
				255, 0, 0, 255, // 红色像素
				0, 0, 255, 255, // 蓝色像素
				0, 0, 255, 255, // 蓝色像素
				255, 0, 0, 255, // 红色像素
			},
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 2, 2))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(1, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				return img
			}(),
			wantErr: false,
		},
		{
			name: "RGB format 2x2",
			decoder: Decoder{
				Format: RGB,
				Width:  2,
				Height: 2,
			},
			input: []byte{
				255, 0, 0, // 红色像素
				0, 0, 255, // 蓝色像素
				0, 0, 255, // 蓝色像素
				255, 0, 0, // 红色像素
			},
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 2, 2))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(1, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				return img
			}(),
			wantErr: false,
		},
		{
			name: "RGBA with compression",
			decoder: Decoder{
				Format:     RGBA,
				Width:      2,
				Height:     2,
				Decompress: true,
			},
			input: compress([]byte{
				255, 0, 0, 255,
				0, 0, 255, 255,
				0, 0, 255, 255,
				255, 0, 0, 255,
			}),
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 2, 2))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
				img.Set(1, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				return img
			}(),
			wantErr: false,
		},
		{
			name: "PNG format",
			decoder: Decoder{
				Format: PNG,
				// 宽度和高度嵌入在 PNG 数据中并从中推断
			},
			input: func() []byte {
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				var buf bytes.Buffer
				png.Encode(&buf, img)
				return buf.Bytes()
			}(),
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				return img
			}(),
			wantErr: false,
		},
		{
			name: "invalid format",
			decoder: Decoder{
				Format: 999,
				Width:  2,
				Height: 2,
			},
			input:   []byte{0, 0, 0},
			want:    nil,
			wantErr: true,
		},
		{
			name: "incomplete RGBA data",
			decoder: Decoder{
				Format: RGBA,
				Width:  2,
				Height: 2,
			},
			input:   []byte{255, 0, 0}, // 不完整的像素数据
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid compressed data",
			decoder: Decoder{
				Format:     RGBA,
				Width:      2,
				Height:     2,
				Decompress: true,
			},
			input:   []byte{1, 2, 3}, // 无效的 zlib 数据
			want:    nil,
			wantErr: true,
		},
		{
			name: "default format (RGBA)",
			decoder: Decoder{
				Width:  1,
				Height: 1,
			},
			input: []byte{255, 0, 0, 255},
			want: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
				return img
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.decoder.Decode(bytes.NewReader(tt.input))

			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() 输出不匹配")
				if bounds := got.Bounds(); bounds != tt.want.Bounds() {
					t.Errorf("边界范围得到 %v，需要 %v", bounds, tt.want.Bounds())
				}

				// 比较像素
				bounds := got.Bounds()
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						gotColor := got.At(x, y)
						wantColor := tt.want.At(x, y)
						if !reflect.DeepEqual(gotColor, wantColor) {
							t.Errorf("像素位置 (%d,%d) = %v，需要 %v", x, y, gotColor, wantColor)
						}
					}
				}
			}
		})
	}
}

// TestDecoder_DecodeEdgeCases 测试解码器的边界情况处理。
func TestDecoder_DecodeEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		decoder Decoder
		input   []byte
		wantErr bool
	}{
		{
			name: "zero dimensions",
			decoder: Decoder{
				Format: RGBA,
				Width:  0,
				Height: 0,
			},
			input:   []byte{},
			wantErr: false,
		},
		{
			name: "negative width",
			decoder: Decoder{
				Format: RGBA,
				Width:  -1,
				Height: 1,
			},
			input:   []byte{255, 0, 0, 255},
			wantErr: false, // image 包优雅地处理了这种情况
		},
		{
			name: "very large dimensions",
			decoder: Decoder{
				Format: RGBA,
				Width:  1,
				Height: 1000000, // 非常大的高度
			},
			input:   []byte{255, 0, 0, 255}, // 数据不足
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.decoder.Decode(bytes.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
