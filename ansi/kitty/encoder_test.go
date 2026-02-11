package kitty

import (
	"bytes"
	"compress/zlib"
	"image"
	"image/color"
	"io"
	"testing"
)

// 取自 "image/png" 包
const pngHeader = "\x89PNG\r\n\x1a\n"

// testImage 创建一个带有红色和蓝色图案的简单测试图像
func testImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 红色
	img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // 蓝色
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // 蓝色
	img.Set(1, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // 红色
	return img
}

func TestEncoder_Encode(t *testing.T) {
	tests := []struct {
		name    string
		encoder Encoder
		img     image.Image
		wantErr bool
		verify  func([]byte) error
	}{
		{
			name: "nil 图像",
			encoder: Encoder{
				Format: RGBA,
			},
			img:     nil,
			wantErr: false,
			verify: func(got []byte) error {
				if len(got) != 0 {
					t.Errorf("预期 nil 图像输出为空，得到 %d 字节", len(got))
				}
				return nil
			},
		},
		{
			name: "RGBA 格式",
			encoder: Encoder{
				Format: RGBA,
			},
			img:     testImage(),
			wantErr: false,
			verify: func(got []byte) error {
				expected := []byte{
					255, 0, 0, 255, // 红色像素
					0, 0, 255, 255, // 蓝色像素
					0, 0, 255, 255, // 蓝色像素
					255, 0, 0, 255, // 红色像素
				}
				if !bytes.Equal(got, expected) {
					t.Errorf("意外的 RGBA 输出\n得到:  %v\n预期: %v", got, expected)
				}
				return nil
			},
		},
		{
			name: "RGB 格式",
			encoder: Encoder{
				Format: RGB,
			},
			img:     testImage(),
			wantErr: false,
			verify: func(got []byte) error {
				expected := []byte{
					255, 0, 0, // 红色像素
					0, 0, 255, // 蓝色像素
					0, 0, 255, // 蓝色像素
					255, 0, 0, // 红色像素
				}
				if !bytes.Equal(got, expected) {
					t.Errorf("意外的 RGB 输出\n得到:  %v\n预期: %v", got, expected)
				}
				return nil
			},
		},
		{
			name: "PNG 格式",
			encoder: Encoder{
				Format: PNG,
			},
			img:     testImage(),
			wantErr: false,
			verify: func(got []byte) error {
				// 验证 PNG 头部
				if len(got) < 8 || !bytes.Equal(got[:8], []byte(pngHeader)) {
					t.Error("无效的 PNG 头部")
				}
				return nil
			},
		},
		{
			name: "无效格式",
			encoder: Encoder{
				Format: 999, // 无效格式
			},
			img:     testImage(),
			wantErr: true,
			verify:  nil,
		},
		{
			name: "RGBA 带压缩",
			encoder: Encoder{
				Format:   RGBA,
				Compress: true,
			},
			img:     testImage(),
			wantErr: false,
			verify: func(got []byte) error {
				// 解压缩数据
				r, err := zlib.NewReader(bytes.NewReader(got))
				if err != nil {
					return err //nolint:wrapcheck
				}
				defer r.Close()

				decompressed, err := io.ReadAll(r)
				if err != nil {
					return err //nolint:wrapcheck
				}

				expected := []byte{
					255, 0, 0, 255, // 红色像素
					0, 0, 255, 255, // 蓝色像素
					0, 0, 255, 255, // 蓝色像素
					255, 0, 0, 255, // 红色像素
				}
				if !bytes.Equal(decompressed, expected) {
					t.Errorf("意外的解压缩输出\n得到:  %v\n预期: %v", decompressed, expected)
				}
				return nil
			},
		},
		{
			name: "零格式默认为 RGBA",
			encoder: Encoder{
				Format: 0,
			},
			img:     testImage(),
			wantErr: false,
			verify: func(got []byte) error {
				expected := []byte{
					255, 0, 0, 255, // 红色像素
					0, 0, 255, 255, // 蓝色像素
					0, 0, 255, 255, // 蓝色像素
					255, 0, 0, 255, // 红色像素
				}
				if !bytes.Equal(got, expected) {
					t.Errorf("意外的 RGBA 输出\n得到:  %v\n预期: %v", got, expected)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tt.encoder.Encode(&buf, tt.img)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.verify != nil {
				if err := tt.verify(buf.Bytes()); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestEncoder_EncodeWithDifferentImageTypes(t *testing.T) {
	// 创建不同类型的图像用于测试
	rgba := image.NewRGBA(image.Rect(0, 0, 1, 1))
	rgba.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	gray := image.NewGray(image.Rect(0, 0, 1, 1))
	gray.Set(0, 0, color.Gray{Y: 128})

	tests := []struct {
		name    string
		img     image.Image
		format  int
		wantLen int
	}{
		{
			name:    "RGBA 图像转换为 RGBA 格式",
			img:     rgba,
			format:  RGBA,
			wantLen: 4, // 每像素 4 字节
		},
		{
			name:    "Gray 图像转换为 RGBA 格式",
			img:     gray,
			format:  RGBA,
			wantLen: 4, // 每像素 4 字节
		},
		{
			name:    "RGBA 图像转换为 RGB 格式",
			img:     rgba,
			format:  RGB,
			wantLen: 3, // 每像素 3 字节
		},
		{
			name:    "Gray 图像转换为 RGB 格式",
			img:     gray,
			format:  RGB,
			wantLen: 3, // 每像素 3 字节
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := Encoder{Format: tt.format}

			err := enc.Encode(&buf, tt.img)
			if err != nil {
				t.Errorf("Encode() error = %v", err)
				return
			}

			if got := buf.Len(); got != tt.wantLen {
				t.Errorf("Encode() output length = %v, want %v", got, tt.wantLen)
			}
		})
	}
}
