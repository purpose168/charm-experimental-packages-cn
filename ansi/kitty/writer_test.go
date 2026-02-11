package kitty

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteKittyGraphics(t *testing.T) {
	// 创建测试图像
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// 创建大型测试图像（大于 [MaxChunkSize] 4096 字节）
	imgLarge := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			imgLarge.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	// 创建临时测试文件
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-image")
	if err := os.WriteFile(tmpFile, []byte("test image data"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		img       image.Image
		opts      *Options
		wantError bool
		check     func(t *testing.T, output string)
	}{
		{
			name: "直接传输",
			img:  img,
			opts: &Options{
				Transmission: Direct,
				Format:       RGB,
			},
			wantError: false,
			check: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "\x1b_G") {
					t.Error("输出应以 ESC 序列开头")
				}
				if !strings.HasSuffix(output, "\x1b\\") {
					t.Error("输出应以 ST 序列结尾")
				}
				if !strings.Contains(output, "f=24") {
					t.Error("输出应包含格式规范")
				}
			},
		},
		{
			name: "分块传输",
			img:  imgLarge,
			opts: &Options{
				Transmission: Direct,
				Format:       RGB,
				Chunk:        true,
			},
			wantError: false,
			check: func(t *testing.T, output string) {
				chunks := strings.Split(output, "\x1b\\")
				if len(chunks) < 2 {
					t.Error("输出应包含多个块")
				}

				chunks = chunks[:len(chunks)-1] // 移除最后一个空块
				for i, chunk := range chunks {
					if i == len(chunks)-1 {
						if !strings.Contains(chunk, "m=0") {
							t.Errorf("输出应包含块 %d 的数据结束指示符 %q", i, chunk)
						}
					} else {
						if !strings.Contains(chunk, "m=1") {
							t.Errorf("输出应包含块 %d 的块指示符 %q", i, chunk)
						}
					}
				}
			},
		},
		{
			name: "文件传输",
			img:  img,
			opts: &Options{
				Transmission: File,
				File:         tmpFile,
			},
			wantError: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, base64.StdEncoding.EncodeToString([]byte(tmpFile))) {
					t.Error("输出应包含编码的文件路径")
				}
			},
		},
		{
			name: "临时文件传输",
			img:  img,
			opts: &Options{
				Transmission: TempFile,
			},
			wantError: false,
			check: func(t *testing.T, output string) {
				output = strings.TrimPrefix(output, "\x1b_G")
				output = strings.TrimSuffix(output, "\x1b\\")
				payload := strings.Split(output, ";")[1]
				fn, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					t.Error("输出应包含 base64 编码的临时文件路径")
				}
				if !strings.Contains(string(fn), "tty-graphics-protocol") {
					t.Error("输出应包含临时文件路径")
				}
				if !strings.Contains(output, "t=t") {
					t.Error("输出应包含传输规范")
				}
			},
		},
		{
			name: "启用压缩",
			img:  img,
			opts: &Options{
				Transmission: Direct,
				Compression:  Zlib,
			},
			wantError: false,
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "o=z") {
					t.Error("输出应包含压缩规范")
				}
			},
		},
		{
			name: "无效的文件路径",
			img:  img,
			opts: &Options{
				Transmission: File,
				File:         "/nonexistent/file",
			},
			wantError: true,
			check:     nil,
		},
		{
			name:      "空选项",
			img:       img,
			opts:      nil,
			wantError: false,
			check: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "\x1b_G") {
					t.Error("输出应以 ESC 序列开头")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := EncodeGraphics(&buf, tt.img, tt.opts)

			if (err != nil) != tt.wantError {
				t.Errorf("WriteKittyGraphics() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.check != nil {
				tt.check(t, buf.String())
			}
		})
	}
}

func TestWriteKittyGraphicsEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		img       image.Image
		opts      *Options
		wantError bool
	}{
		{
			name: "零尺寸图像",
			img:  image.NewRGBA(image.Rect(0, 0, 0, 0)),
			opts: &Options{
				Transmission: Direct,
			},
			wantError: false,
		},
		{
			name: "共享内存传输",
			img:  image.NewRGBA(image.Rect(0, 0, 1, 1)),
			opts: &Options{
				Transmission: SharedMemory,
			},
			wantError: true, // 未实现
		},
		{
			name: "无文件路径的文件传输",
			img:  image.NewRGBA(image.Rect(0, 0, 1, 1)),
			opts: &Options{
				Transmission: File,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := EncodeGraphics(&buf, tt.img, tt.opts)

			if (err != nil) != tt.wantError {
				t.Errorf("WriteKittyGraphics() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
