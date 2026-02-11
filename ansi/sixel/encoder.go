package sixel

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/bits-and-blooms/bitset"
)

// Sixels 是一种通过写入大量 ANSI 转义数据来向终端写入图像的协议。
// 它们的工作原理是将 6 像素的列编码为单个字符（与 base64 每 6 位编码数据的方式非常相似）。
// Sixel 图像是调色板图像，在图像 blob 的开头建立调色板，像素在写入像素数据时通过索引标识调色板条目。
//
// Sixel 每次写入一个 6 像素高的带（band），一次写入一种颜色。对于每个带，先写入单一颜色的像素，
// 然后写入回车符将"光标"带回带的开头，在那里选择新颜色并写入像素。这将持续进行，直到绘制完整个带，
// 此时写入换行符开始绘制下一个带。

// Sixel 控制函数。
const (
	LineBreak        byte = '-'
	CarriageReturn   byte = '$'
	RepeatIntroducer byte = '!'
	ColorIntroducer  byte = '#'
	RasterAttribute  byte = '"'
)

// Encoder 是一个 Sixel 编码器。它将图像编码为 Sixel 数据格式。
type Encoder struct{}

// Encode 接受一个 Image 并向 Writer 写入 sixel 数据。Sixel 数据将是
// 结束 DCS 参数的 'q' 之后、结束序列的 ST 之前的所有内容。
// 这意味着它包括像素指标和调色板。
func (e *Encoder) Encode(w io.Writer, img image.Image) error {
	if img == nil {
		return nil
	}

	imageBounds := img.Bounds()

	// 如果未设置，则设置默认的光栅 1:1 宽高比
	if _, err := WriteRaster(w, 1, 1, imageBounds.Dx(), imageBounds.Dy()); err != nil {
		return fmt.Errorf("编码光栅时出错: %w", err)
	}

	palette := newSixelPalette(img, MaxColors)

	for paletteIndex, color := range palette.PaletteColors {
		e.encodePaletteColor(w, paletteIndex, color)
	}

	scratch := newSixelBuilder(imageBounds.Dx(), imageBounds.Dy(), palette)

	for y := range imageBounds.Dy() {
		for x := range imageBounds.Dx() {
			scratch.SetColor(x, y, img.At(x, y))
		}
	}

	pixels := scratch.GeneratePixels()
	io.WriteString(w, pixels) //nolint:errcheck,gosec

	return nil
}

func (e *Encoder) encodePaletteColor(w io.Writer, paletteIndex int, c sixelColor) {
	// 初始化调色板条目
	// #<a>;<b>;<c>;<d>;<e>
	// a = 调色板索引
	// b = 颜色类型，2 表示 RGB
	// c = R
	// d = G
	// e = B

	w.Write([]byte{ColorIntroducer})              //nolint:errcheck,gosec
	io.WriteString(w, strconv.Itoa(paletteIndex)) //nolint:errcheck,gosec
	io.WriteString(w, ";2;")                      //nolint:errcheck,gosec
	io.WriteString(w, strconv.Itoa(int(c.Red)))   //nolint:errcheck,gosec
	w.Write([]byte{';'})                          //nolint:errcheck,gosec
	io.WriteString(w, strconv.Itoa(int(c.Green))) //nolint:errcheck,gosec
	w.Write([]byte{';'})                          //nolint:errcheck,gosec
	io.WriteString(w, strconv.Itoa(int(c.Blue)))  //nolint:errcheck,gosec
}

// sixelBuilder 是一个临时结构，用于创建 SixelImage。它处理将像素分解为位，
// 然后将它们编码为 sixel 数据字符串。包括 RLE 处理。
//
// 使用 sixelBuilder 分两个阶段完成。首先，使用 SetColor 将所有像素写入内部 BitSet 数据。
// 然后，调用 GeneratePixels 检索以 sixel 格式编码的像素数据字符串。
type sixelBuilder struct {
	SixelPalette sixelPalette

	imageHeight int
	imageWidth  int

	pixelBands bitset.BitSet

	imageData   strings.Builder
	repeatByte  byte
	repeatCount int
}

// newSixelBuilder 创建一个 sixelBuilder 并准备写入。
func newSixelBuilder(width, height int, palette sixelPalette) sixelBuilder {
	scratch := sixelBuilder{
		imageWidth:   width,
		imageHeight:  height,
		SixelPalette: palette,
	}

	return scratch
}

// BandHeight 返回此图像由多少个 6 像素带组成。
func (s *sixelBuilder) BandHeight() int {
	bandHeight := s.imageHeight / 6
	if s.imageHeight%6 != 0 {
		bandHeight++
	}

	return bandHeight
}

// SetColor 将单个像素写入 sixelBuilder 的内部位集数据，供 GeneratePixels 使用。
func (s *sixelBuilder) SetColor(x int, y int, color color.Color) {
	bandY := y / 6
	paletteIndex := s.SixelPalette.ColorIndex(sixelConvertColor(color))

	bit := s.BandHeight()*s.imageWidth*6*paletteIndex + bandY*s.imageWidth*6 + (x * 6) + (y % 6)
	s.pixelBands.Set(uint(bit)) //nolint:gosec
}

// GeneratePixels 用于将像素数据写入内部的 imageData 字符串构建器。
// 在调用此方法之前，必须使用 SetColor 将图像中的所有像素写入 sixelBuilder。
// 此方法返回一个表示像素数据的字符串。Sixel 字符串由五部分组成：
// ISC <header> <palette> <pixels> ST
// header 包含一些任意选项，指示如何绘制 sixel 图像。
// palette 将调色板索引映射到 RGB 颜色
// pixels 指示使用哪些调色板颜色绘制哪些像素。
//
// GeneratePixels 仅生成字符串的 <pixels> 部分。其余部分由 Style.RenderSixelImage 写入。
func (s *sixelBuilder) GeneratePixels() string {
	s.imageData = strings.Builder{}
	bandHeight := s.BandHeight()

	for bandY := range bandHeight {
		if bandY > 0 {
			s.writeControlRune(LineBreak)
		}

		hasWrittenAColor := false

		for paletteIndex := range s.SixelPalette.PaletteColors {
			if s.SixelPalette.PaletteColors[paletteIndex].Alpha < 1 {
				// 不为完全透明的像素绘制任何内容
				continue
			}

			firstColorBit := uint(s.BandHeight()*s.imageWidth*6*paletteIndex + bandY*s.imageWidth*6) //nolint:gosec
			nextColorBit := firstColorBit + uint(s.imageWidth*6)                                     //nolint:gosec

			firstSetBitInBand, anySet := s.pixelBands.NextSet(firstColorBit)
			if !anySet || firstSetBitInBand >= nextColorBit {
				// 此行中不出现该颜色
				continue
			}

			if hasWrittenAColor {
				s.writeControlRune(CarriageReturn)
			}
			hasWrittenAColor = true

			s.writeControlRune(ColorIntroducer)
			s.imageData.WriteString(strconv.Itoa(paletteIndex))
			for x := 0; x < s.imageWidth; x += 4 {
				bit := firstColorBit + uint(x*6) //nolint:gosec
				word := s.pixelBands.GetWord64AtBit(bit)

				pixel1 := byte((word & 63) + '?')
				pixel2 := byte(((word >> 6) & 63) + '?')
				pixel3 := byte(((word >> 12) & 63) + '?')
				pixel4 := byte(((word >> 18) & 63) + '?')

				s.writeImageRune(pixel1)

				if x+1 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel2)

				if x+2 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel3)

				if x+3 >= s.imageWidth {
					continue
				}
				s.writeImageRune(pixel4)
			}
		}
	}

	s.writeControlRune('-')
	return s.imageData.String()
}

// writeImageRune 将单个像素行（6 个像素）写入像素数据。数据不会直接写入 imageData，
// 而是会被缓冲以用于 RLE 处理。
func (s *sixelBuilder) writeImageRune(r byte) {
	if r == s.repeatByte {
		s.repeatCount++
		return
	}

	s.flushRepeats()
	s.repeatByte = r
	s.repeatCount = 1
}

// writeControlRune 将特殊字符（如换行符或回车符）写入。如果需要，它将首先调用 flushRepeats。
func (s *sixelBuilder) writeControlRune(r byte) {
	if s.repeatCount > 0 {
		s.flushRepeats()
		s.repeatCount = 0
		s.repeatByte = 0
	}

	s.imageData.WriteByte(r)
}

// flushRepeats 用于在实际更改时将当前的 repeatByte 写入 imageData。
// 此缓冲用于管理 sixelBuilder 中的 RLE。
func (s *sixelBuilder) flushRepeats() {
	if s.repeatCount == 0 {
		return
	}

	// 仅在实际提供空间节省时才使用 RLE 形式写入
	if s.repeatCount > 3 {
		WriteRepeat(&s.imageData, s.repeatCount, s.repeatByte) //nolint:errcheck,gosec
		return
	}

	for range s.repeatCount {
		s.imageData.WriteByte(s.repeatByte)
	}
}
