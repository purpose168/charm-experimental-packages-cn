// Package mosaic 提供Unicode图像渲染器。
package mosaic

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math"
	"strings"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	xdraw "golang.org/x/image/draw"
)

// 块定义。
var (
	halfBlocks = []block{
		{Char: '▀', Coverage: [4]bool{true, true, false, false}, CoverageMap: "██\n  "},   // 上半块。
		{Char: '▄', Coverage: [4]bool{false, false, true, true}, CoverageMap: "  \n██"},   // 下半块。
		{Char: ' ', Coverage: [4]bool{false, false, false, false}, CoverageMap: "  \n  "}, // 空格。
		{Char: '█', Coverage: [4]bool{true, true, true, true}, CoverageMap: "██\n██"},     // 全块。
	}
	quarterBlocks = []block{
		{Char: '▘', Coverage: [4]bool{true, false, false, false}, CoverageMap: "█ \n  "}, // 左上象限。
		{Char: '▝', Coverage: [4]bool{false, true, false, false}, CoverageMap: " █\n  "}, // 右上象限。
		{Char: '▖', Coverage: [4]bool{false, false, true, false}, CoverageMap: "  \n█ "}, // 左下象限。
		{Char: '▗', Coverage: [4]bool{false, false, false, true}, CoverageMap: "  \n █"}, // 右下象限。
		{Char: '▌', Coverage: [4]bool{true, false, true, false}, CoverageMap: "█ \n█ "},  // 左半块。
		{Char: '▐', Coverage: [4]bool{false, true, false, true}, CoverageMap: " █\n █"},  // 右半块。
		{Char: '▀', Coverage: [4]bool{true, true, false, false}, CoverageMap: "██\n  "},  // 上半块（已添加）。
		{Char: '▄', Coverage: [4]bool{false, false, true, true}, CoverageMap: "  \n██"},  // 下半块（已添加）。
	}
	complexBlocks = []block{
		{Char: '▙', Coverage: [4]bool{true, false, true, true}, CoverageMap: "█ \n██"},  // 左上象限和下半块。
		{Char: '▟', Coverage: [4]bool{false, true, true, true}, CoverageMap: " █\n██"},  // 右上象限和下半块。
		{Char: '▛', Coverage: [4]bool{true, true, true, false}, CoverageMap: "██\n█ "},  // 上半块和左下象限。
		{Char: '▜', Coverage: [4]bool{true, true, false, true}, CoverageMap: "██\n █"},  // 上半块和右下象限。
		{Char: '▚', Coverage: [4]bool{true, false, false, true}, CoverageMap: "█ \n █"}, // 左上象限和右下象限。
		{Char: '▞', Coverage: [4]bool{false, true, true, false}, CoverageMap: " █\n█ "}, // 右上象限和左下象限。
	}
)

// Block 表示不同的Unicode块字符。
type block struct {
	Char        rune
	Coverage    [4]bool // 块的哪些部分被填充（true = 填充）。
	CoverageMap string  // 用于调试的覆盖范围的可视化表示。
}

// Symbol 表示渲染图像时使用的符号类型。
type Symbol uint8

// 符号类型。
const (
	All     Symbol = iota // 所有符号
	Half                  // 半块符号
	Quarter               // 四分之一块符号
)

// 在许多情况下，默认阈值级别通常设置为0.5（或50%），
// 这意味着高于此阈值的值被视为正数，
// 而低于此阈值的值被视为负数。
// 值128表示0..255的0.5。
const middleThresholdLevel = 128

// 使用默认值渲染马赛克。
func Render(img image.Image, width int, height int) string {
	m := New().Width(width).Height(height)
	return m.Render(img)
}

// Mosaic 表示Unicode图像渲染器。
//
// 示例：
//
//	```go
//	art := mosaic.New().Width(100). // 限制为100个单元格
//	    Scale(mosaic.Fit).          // 适应宽度
//	    Render()
//	```
type Mosaic struct {
	outputWidth    int    // 输出宽度。
	outputHeight   int    // 输出高度（0表示自动）。
	thresholdLevel uint8  // 考虑像素为设置的阈值（0-255）。
	dither         bool   // 启用抖动（默认为false）。
	useFgBgOnly    bool   // 仅使用前景/背景颜色（无块符号）。
	invertColors   bool   // 反转颜色。
	scale          int    // 缩放级别
	symbols        Symbol // 使用哪些符号："half"（半块）, "quarter"（四分之一块）, "all"（所有）。
}

// New 创建并返回一个[Renderer]。
func New() Mosaic {
	return Mosaic{
		outputWidth:    0,                    // 覆盖宽度。
		outputHeight:   0,                    // 覆盖高度。
		thresholdLevel: middleThresholdLevel, // 中间阈值。
		dither:         false,                // 启用抖动。
		useFgBgOnly:    false,                // 使用块符号。
		invertColors:   false,                // 不反转。
		scale:          1,                    // 不缩放。
		symbols:        Half,                 // 使用半块。
	}
}

// PixelBlock 表示图像中的2x2像素块。
type pixelBlock struct {
	Pixels      [2][2]color.Color // 2x2像素网格。
	AvgFg       color.Color       // 平均前景颜色。
	AvgBg       color.Color       // 平均背景颜色。
	BestSymbol  rune              // 最佳匹配字符。
	BestFgColor color.Color       // 最佳前景颜色。
	BestBgColor color.Color       // 最佳背景颜色。
}

// 表示255。
const u8MaxValue = 0xff

type shiftable interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

func shift[T shiftable](x T) T {
	if x > u8MaxValue {
		x >>= 8
	}
	return x
}

// Scale 设置[Mosaic]上的[ScaleMode]。
func (m Mosaic) Scale(scale int) Mosaic {
	m.scale = scale
	return m
}

// IgnoreBlockSymbols 设置[Mosaic]上的UseFgBgOnly。
func (m Mosaic) IgnoreBlockSymbols(fgOnly bool) Mosaic {
	m.useFgBgOnly = fgOnly
	return m
}

// Dither 设置[Mosaic]上的抖动级别。
func (m Mosaic) Dither(dither bool) Mosaic {
	m.dither = dither
	return m
}

// Threshold 设置[Mosaic]上的阈值级别。
// 它期望一个0-255之间的值，其他值将被忽略。
func (m Mosaic) Threshold(threshold int) Mosaic {
	if threshold >= 0 && threshold <= u8MaxValue {
		m.thresholdLevel = uint8(threshold)
	}

	return m
}

// InvertColors 是否反转马赛克图像的颜色。
func (m Mosaic) InvertColors(invertColors bool) Mosaic {
	m.invertColors = invertColors
	return m
}

// Width 设置图像可以拥有的最大宽度。默认为图像宽度。
func (m Mosaic) Width(width int) Mosaic {
	m.outputWidth = width
	return m
}

// Height 设置图像可以拥有的最大高度。默认为图像高度。
func (m Mosaic) Height(height int) Mosaic {
	m.outputHeight = height
	return m
}

// Symbol 设置马赛克符号类型。
func (m Mosaic) Symbol(symbol Symbol) Mosaic {
	m.symbols = symbol
	return m
}

// Render 将图像渲染为字符串。
func (m *Mosaic) Render(img image.Image) string {
	// Calculate dimensions.
	bounds := img.Bounds()
	srcWidth := bounds.Max.X - bounds.Min.X
	srcHeight := bounds.Max.Y - bounds.Min.Y

	// Determine output dimensions.
	outWidth := srcWidth
	if m.outputWidth > 0 {
		outWidth = m.outputWidth
	}

	outHeight := srcHeight
	if m.outputHeight > 0 {
		outHeight = m.outputHeight
	}

	if outHeight <= 0 {
		// Calculate height based on aspect ratio and character cell proportions.
		// Terminal characters are roughly twice as tall as wide, so we divide by 2.
		const divider = 2
		outHeight = int(float64(outWidth) * float64(srcHeight) / float64(srcWidth) / divider)
		if outHeight < 1 {
			outHeight = 1
		}
	}

	// Scale image according to the scale.
	scaledImg := m.applyScaling(img, outWidth*m.scale, outHeight*m.scale)

	// Apply dithering if enabled.
	if m.dither {
		scaledImg = m.applyDithering(scaledImg)
	}

	// Invert colors if needed.
	if m.invertColors {
		scaledImg = m.invertImage(scaledImg)
	}

	// Generate terminal outpum.
	var output strings.Builder

	// Process the image by 2x2 blocks (representing one character cell).
	imageBounds := scaledImg.Bounds()

	// Set initial blocks based on symbols value (initial/default is half)
	blocks := halfBlocks

	// Quarter blocks.
	if m.symbols == Quarter || m.symbols == All {
		blocks = append(blocks, quarterBlocks...)
	}

	// All block elements (including complex combinations).
	if m.symbols == All {
		blocks = append(blocks, complexBlocks...)
	}

	for y := 0; y < imageBounds.Max.Y; y += 2 {
		for x := 0; x < imageBounds.Max.X; x += 2 {
			// Create and analyze the 2x2 pixel block.
			block := m.createPixelBlock(scaledImg, x, y)

			// Determine best symbol and colors.
			m.findBestRepresentation(block, blocks)

			// Append to output.
			output.WriteString(
				ansi.Style{}.ForegroundColor(block.BestFgColor).BackgroundColor(block.BestBgColor).Styled(string(block.BestSymbol)),
			)
		}
		output.WriteString("\n")
	}

	return output.String()
}

// createPixelBlock 从图像中提取2x2像素块。
func (m *Mosaic) createPixelBlock(img image.Image, x, y int) *pixelBlock {
	block := &pixelBlock{}

	// Extract the 2x2 pixel grid.
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			block.Pixels[dy][dx] = m.getPixelSafe(img, x+dx, y+dy)
		}
	}

	return block
}

// findBestRepresentation 为2x2像素块找到最佳的块字符和颜色。
func (m *Mosaic) findBestRepresentation(block *pixelBlock, availableBlocks []block) {
	// Simple case: use only foreground/background colors.
	if m.useFgBgOnly {
		// Just use the upper half block with top pixels as background and bottom as foreground.
		block.BestSymbol = '▀'
		block.BestBgColor = m.averageColors(block.Pixels[0][0], block.Pixels[0][1])
		block.BestFgColor = m.averageColors(block.Pixels[1][0], block.Pixels[1][1])
		return
	}

	// Determine which pixels are "set" based on threshold.
	pixelMask := [2][2]bool{}
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			// Calculate luminance.
			luma := rgbaToLuminance(block.Pixels[y][x])
			pixelMask[y][x] = luma >= m.thresholdLevel
		}
	}

	// Find the best matching block character.
	bestChar := ' '
	bestScore := math.MaxFloat64

	for _, blockChar := range availableBlocks {
		score := 0.0
		for i := 0; i < 4; i++ {
			y, x := i/2, i%2 //nolint:mnd
			if blockChar.Coverage[i] != pixelMask[y][x] {
				score += 1.0
			}
		}

		if score < bestScore {
			bestScore = score
			bestChar = blockChar.Char
		}
	}

	// Determine foreground and background colors based on the best character.
	var fgPixels, bgPixels []color.Color

	// Get the coverage pattern for the selected character.
	var coverage [4]bool
	for _, b := range availableBlocks {
		if b.Char == bestChar {
			coverage = b.Coverage
			break
		}
	}

	// Assign pixels to foreground or background based on the character's coverage.
	for i := 0; i < 4; i++ {
		y, x := i/2, i%2 //nolint:mnd
		if coverage[i] {
			fgPixels = append(fgPixels, block.Pixels[y][x])
		} else {
			bgPixels = append(bgPixels, block.Pixels[y][x])
		}
	}

	// Calculate average colors.
	if len(fgPixels) > 0 {
		block.BestFgColor = m.averageColors(fgPixels...)
	} else {
		// Default to black if no foreground pixels.
		block.BestFgColor = color.Black
	}

	if len(bgPixels) > 0 {
		block.BestBgColor = m.averageColors(bgPixels...)
	} else {
		// Default to black if no background pixels.
		block.BestBgColor = color.Black
	}

	block.BestSymbol = bestChar
}

// averageColors 计算颜色切片的平均颜色。
func (m *Mosaic) averageColors(colors ...color.Color) color.Color {
	if len(colors) == 0 {
		return color.Black
	}

	var sumR, sumG, sumB, sumA uint32

	for _, c := range colors {
		r, g, b, a := c.RGBA()
		r, g, b, a = shift(r), shift(g), shift(b), shift(a)
		sumR += r
		sumG += g
		sumB += b
		sumA += a
	}

	count := uint32(len(colors)) //nolint:gosec
	return color.RGBA{
		R: uint8(sumR / count), //nolint:gosec
		G: uint8(sumG / count), //nolint:gosec
		B: uint8(sumB / count), //nolint:gosec
		A: uint8(sumA / count), //nolint:gosec
	}
}

// getPixelSafe 返回(x,y)处的颜色，如果越界则返回黑色。
func (m *Mosaic) getPixelSafe(img image.Image, x, y int) color.RGBA {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return color.RGBA{0, 0, 0, 255}
	}

	r8, g8, b8, a8 := img.At(x, y).RGBA()
	return color.RGBA{
		R: uint8(r8 >> 8), //nolint:gosec,mnd
		G: uint8(g8 >> 8), //nolint:gosec,mnd
		B: uint8(b8 >> 8), //nolint:gosec,mnd
		A: uint8(a8 >> 8), //nolint:gosec,mnd
	}
}

// applyScaling 将图像调整为指定的尺寸。
func (m *Mosaic) applyScaling(img image.Image, width, height int) image.Image {
	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)
	xdraw.ApproxBiLinear.Scale(dst, rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

// applyDithering 应用Floyd-Steinberg抖动。
func (m *Mosaic) applyDithering(img image.Image) image.Image {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

// invertImage 反转图像的颜色。
func (m *Mosaic) invertImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	result := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r8, g8, b8, a8 := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			result.Set(x, y, color.RGBA{
				R: uint8(255 - (r8 >> 8)), //nolint:gosec,mnd
				G: uint8(255 - (g8 >> 8)), //nolint:gosec,mnd
				B: uint8(255 - (b8 >> 8)), //nolint:gosec,mnd
				A: uint8(a8 >> 8),         //nolint:gosec,mnd
			})
		}
	}

	return result
}

// rgbaToLuminance 将RGBA颜色转换为亮度（brightness）。
func rgbaToLuminance(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	r, g, b = shift(r), shift(g), shift(b)
	// 加权RGB以考虑人类感知
	// 来源：https://www.w3.org/TR/AERT/#color-contrast
	// 上下文：https://stackoverflow.com/questions/596216/formula-to-determine-perceived-brightness-of-rgb-color
	return uint8(float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114)
}
