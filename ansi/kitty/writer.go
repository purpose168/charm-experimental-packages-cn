package kitty

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

var (
	// GraphicsTempDir 是存储临时文件的目录。
	// 这在 [WriteKittyGraphics] 中与 [os.CreateTemp] 一起使用。
	GraphicsTempDir = ""

	// GraphicsTempPattern 是用于创建临时文件的模式。
	// 这在 [WriteKittyGraphics] 中与 [os.CreateTemp] 一起使用。
	// Kitty 图形协议要求文件路径包含子字符串 "tty-graphics-protocol"。
	GraphicsTempPattern = "tty-graphics-protocol-*"
)

// EncodeGraphics 使用 Kitty 图形协议和给定的选项将图像写入 w。如果 o.Chunk 为 true，则分块写入数据。
//
// 当从文件渲染图像时，您可以省略 m 并使用 nil。在这种情况下，
// 您必须在 o.File 中提供文件路径，并使用 o.Transmission = [File]。
// 您还可以使用 o.Transmission = [TempFile] 将图像写入临时文件。
// 在这种情况下，文件路径将被忽略，图像将被写入由终端自动删除的临时文件。
//
// 请参阅 https://sw.kovidgoyal.net/kitty/graphics-protocol/
func EncodeGraphics(w io.Writer, m image.Image, o *Options) error {
	if o == nil {
		o = &Options{}
	}

	if o.Transmission == 0 && len(o.File) != 0 {
		o.Transmission = File
	}

	var data bytes.Buffer // 要编码为 base64 的数据
	e := &Encoder{
		Compress: o.Compression == Zlib,
		Format:   o.Format,
	}

	switch o.Transmission {
	case Direct:
		if err := e.Encode(&data, m); err != nil {
			return fmt.Errorf("failed to encode direct image: %w", err)
		}

	case SharedMemory:
		//nolint:godox
		// TODO: Implement shared memory
		return fmt.Errorf("shared memory transmission is not yet implemented")

	case File:
		if len(o.File) == 0 {
			return ErrMissingFile
		}

		f, err := os.Open(o.File)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		defer f.Close() //nolint:errcheck

		stat, err := f.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}

		mode := stat.Mode()
		if !mode.IsRegular() {
			return fmt.Errorf("file is not a regular file")
		}

		// 将文件路径写入缓冲区
		if _, err := data.WriteString(f.Name()); err != nil {
			return fmt.Errorf("failed to write file path to buffer: %w", err)
		}

	case TempFile:
		f, err := os.CreateTemp(GraphicsTempDir, GraphicsTempPattern)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		defer f.Close() //nolint:errcheck

		if err := e.Encode(f, m); err != nil {
			return fmt.Errorf("failed to encode image to file: %w", err)
		}

		// 将文件路径写入缓冲区
		if _, err := data.WriteString(f.Name()); err != nil {
			return fmt.Errorf("failed to write file path to buffer: %w", err)
		}
	}

	// 将图像编码为 base64
	var payload bytes.Buffer // 要写入 w 的 base64 编码图像
	b64 := base64.NewEncoder(base64.StdEncoding, &payload)
	if _, err := data.WriteTo(b64); err != nil {
		return fmt.Errorf("failed to write base64 encoded image to payload: %w", err)
	}
	if err := b64.Close(); err != nil {
		return err //nolint:wrapcheck
	}

	// 如果不分块，一次性写入所有内容
	if !o.Chunk {
		_, err := io.WriteString(w, ansi.KittyGraphics(payload.Bytes(), o.Options()...))
		return err //nolint:wrapcheck
	}

	// 分块写入
	var (
		err error
		n   int
	)
	chunk := make([]byte, MaxChunkSize)
	isFirstChunk := true
	chunkFormatter := o.ChunkFormatter
	if chunkFormatter == nil {
		// 默认不进行格式化
		chunkFormatter = func(s string) string { return s }
	}

	for {
		// 如果读取的大小小于块大小 [MaxChunkSize]，则停止。
		n, err = io.ReadFull(&payload, chunk)
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read chunk: %w", err)
		}

		opts := buildChunkOptions(o, isFirstChunk, false)
		if _, err := io.WriteString(w,
			chunkFormatter(ansi.KittyGraphics(chunk[:n], opts...))); err != nil {
			return err //nolint:wrapcheck
		}

		isFirstChunk = false
	}

	// 写入最后一个块
	opts := buildChunkOptions(o, isFirstChunk, true)
	_, err = io.WriteString(w, chunkFormatter(ansi.KittyGraphics(chunk[:n], opts...)))
	return err //nolint:wrapcheck
}

// buildChunkOptions 为块创建选项切片。
func buildChunkOptions(o *Options, isFirstChunk, isLastChunk bool) []string {
	var opts []string
	if isFirstChunk {
		opts = o.Options()
	} else {
		// 这些选项在后续块中是允许的
		if o.Quite > 0 {
			opts = append(opts, fmt.Sprintf("q=%d", o.Quite))
		}
		if o.Action == Frame {
			opts = append(opts, "a=f")
		}
	}

	if !isFirstChunk || !isLastChunk {
		// 当我们只有一个块时，不需要编码 (m=) 选项。
		if isLastChunk {
			opts = append(opts, "m=0")
		} else {
			opts = append(opts, "m=1")
		}
	}
	return opts
}
