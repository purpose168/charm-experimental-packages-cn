package vttest

import (
	"image"
	"image/color"
	"image/draw"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/math/fixed"
)

// DefaultDrawer 是用于创建终端屏幕图像的默认选项。
var DefaultDrawer = func() *Drawer {
	cellW, cellH := 10, 16
	regular, _ := freetype.ParseFont(gomono.TTF)
	bold, _ := freetype.ParseFont(gomonobold.TTF)
	italic, _ := freetype.ParseFont(gomonoitalic.TTF)
	boldItalic, _ := freetype.ParseFont(gomonobolditalic.TTF)
	faceOpts := &truetype.Options{
		Size:    14, // 字体大小 -2 以考虑内边距
		DPI:     72,
		Hinting: font.HintingFull,
	}
	regularFace := truetype.NewFace(regular, faceOpts)
	boldFace := truetype.NewFace(bold, faceOpts)
	italicFace := truetype.NewFace(italic, faceOpts)
	boldItalicFace := truetype.NewFace(boldItalic, faceOpts)

	return &Drawer{
		CellWidth:      cellW,
		CellHeight:     cellH,
		RegularFace:    regularFace,
		BoldFace:       boldFace,
		ItalicFace:     italicFace,
		BoldItalicFace: boldItalicFace,
	}
}()

// Drawer 包含用于将终端模拟器屏幕绘制到图像的选项。
type Drawer struct {
	// CellWidth 是每个单元格的宽度（以像素为单位）。默认为 10。
	CellWidth int
	// CellHeight 是每个单元格的高度（以像素为单位）。默认为 16。
	CellHeight int
	// RegularFace 是用于普通文本的字体。默认为 Go mono。
	RegularFace font.Face
	// BoldFace 是用于粗体文本的字体。如果为 nil，则使用 Go mono bold。
	BoldFace font.Face
	// ItalicFace 是用于斜体文本的字体。如果为 nil，则使用 Go mono italic。
	ItalicFace font.Face
	// BoldItalicFace 是用于粗斜体文本的字体。如果为 nil，则使用 Go mono bold italic。
	BoldItalicFace font.Face
}

// Draw 使用抽屉选项将 [uv.Screen] 绘制到图像中。
//
// 如果 s 实现了 [BackgroundColor]() 方法，则使用它来填充背景。否则，使用 [color.Black]。
func (d *Drawer) Draw(t uv.Screen) image.Image {
	opt := *d
	if opt.CellWidth <= 0 {
		opt.CellWidth = DefaultDrawer.CellWidth
	}
	if opt.CellHeight <= 0 {
		opt.CellHeight = DefaultDrawer.CellHeight
	}
	if opt.RegularFace == nil {
		opt.RegularFace = DefaultDrawer.RegularFace
	}
	if opt.BoldFace == nil {
		opt.BoldFace = DefaultDrawer.BoldFace
	}
	if opt.ItalicFace == nil {
		opt.ItalicFace = DefaultDrawer.ItalicFace
	}
	if opt.BoldItalicFace == nil {
		opt.BoldItalicFace = DefaultDrawer.BoldItalicFace
	}

	area := t.Bounds()
	width, height := area.Dx(), area.Dy()
	r := image.Rect(0, 0, width*opt.CellWidth, height*opt.CellHeight)
	img := image.NewRGBA(r)

	// 填充背景
	var bg color.Color = color.Black
	if tbg, ok := t.(interface {
		BackgroundColor() color.Color
	}); ok {
		if bgc := tbg.BackgroundColor(); bgc != nil {
			bg = bgc
		}
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	// 绘制单元格
	drawCell := func(x, y int, cell *uv.Cell) {
		px := x * opt.CellWidth
		py := y * opt.CellHeight
		dot := fixed.P(px, py+opt.CellHeight-4) // 从底部起 4 像素作为基线
		style := cell.Style
		attrs := style.Attrs
		fg := style.Fg
		if fg == nil {
			fg = color.White
		}
		face := opt.RegularFace
		if attrs&uv.AttrBold != 0 && attrs&uv.AttrItalic != 0 {
			face = opt.BoldItalicFace
		} else if attrs&uv.AttrBold != 0 {
			face = opt.BoldFace
		} else if attrs&uv.AttrItalic != 0 {
			face = opt.ItalicFace
		}

		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(fg),
			Face: face,
			Dot:  dot,
		}
		drawer.DrawString(cell.Content)

		// 处理下划线
		//nolint:godox
		// TODO: 实现更多下划线样式
		// 目前，我们只支持单下划线
		if cell.Style.Underline > uv.UnderlineNone {
			col := cell.Style.UnderlineColor
			if col == nil {
				col = fg
			}
			for i := range opt.CellWidth {
				img.Set(px+i, py+opt.CellHeight-2, col)
			}
		}
	}

	// 遍历屏幕单元格
	for y := range height {
		for x := 0; x < width; {
			cell := t.CellAt(x, y)
			if cell == nil {
				cell = &uv.EmptyCell
			}
			drawCell(x, y, cell)
			x += cell.Width
		}
	}

	return img
}

// Image 返回终端模拟器屏幕的图像。
func (t *Terminal) Image() image.Image {
	t.mu.Lock()
	defer t.mu.Unlock()
	return DefaultDrawer.Draw(t.Emulator)
}
