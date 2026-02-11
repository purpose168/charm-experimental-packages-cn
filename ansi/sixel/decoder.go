package sixel

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"slices"
)

// Decoder 是一个 Sixel 图像解码器。它从 io.Reader 读取 Sixel 图像数据，
// 并将其解码为 image.Image。
type Decoder struct{}

// Decode 将 Sixel 图像数据解析为图像或返回错误。由于
// Sixel 图像格式没有可预测的大小，Sixel 图像数据的结尾只能在从读取器读取 ST、ESC 或 BEL 时识别。
// 为了避免逐字节读取读取器而错过结尾，此方法直接接受字节切片而不是读取器。调用者
// 应读取整个转义序列，并将序列的 Ps..Ps 部分传递给此方法。
func (d *Decoder) Decode(r io.Reader) (image.Image, error) {
	rd := bufio.NewReader(r)
	peeked, err := rd.Peek(1)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	var bounds image.Rectangle
	var raster Raster
	if peeked[0] == RasterAttribute {
		var read int
		n := 16
		for {
			peeked, err = rd.Peek(n) // 随机数字，只需要读取几个字节
			if err != nil {
				return nil, err //nolint:wrapcheck
			}

			raster, read = DecodeRaster(peeked)
			if read == 0 {
				return nil, ErrInvalidRaster
			}
			if read >= n {
				// 我们需要读取更多字节以获取完整的光栅
				n *= 2
				continue
			}

			rd.Discard(read) //nolint:errcheck,gosec,gosec
			break
		}

		bounds = image.Rect(0, 0, raster.Ph, raster.Pv)
	}

	if bounds.Max.X == 0 || bounds.Max.Y == 0 {
		// 我们正在解析没有像素指标的图像，所以为主读取循环取消读取该字节
		// 在开始解码之前，窥视整个缓冲区以获取图像大小
		var data []byte
		toPeak := 64 // 窥视的任意字节数
		for {
			data, err = rd.Peek(toPeak)
			if err != nil || len(data) < toPeak {
				break
			}
			toPeak *= 2
		}

		width, height := d.scanSize(data)
		bounds = image.Rect(0, 0, width, height)
	}

	img := image.NewRGBA(bounds)
	palette := DefaultPalette()
	var currentX, currentBandY, currentPaletteIndex int

	// 用于解码 Sixel 命令的数据缓冲区
	data := make([]byte, 0, 6) // 要读取的任意字节数
	// i := 0                     // 跟踪数据缓冲区索引
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return img, d.readError(err)
		}

		count := 1 // Sixel 命令的默认计数
		switch {
		case b == LineBreak: // LF
			currentBandY++
			currentX = 0
		case b == CarriageReturn: // CR
			currentX = 0
		case b == ColorIntroducer: // #
			data = data[:0]
			data = append(data, b)
			for {
				b, err = rd.ReadByte()
				if err != nil {
					return img, d.readError(err)
				}
				// 读取字节直到遇到非颜色字节，即非数字和非分号
				if (b < '0' || b > '9') && b != ';' {
					rd.UnreadByte() //nolint:errcheck,gosec
					break
				}

				data = append(data, b)
			}

			// 调色板操作
			c, n := DecodeColor(data)
			if n == 0 {
				return img, ErrInvalidColor
			}

			currentPaletteIndex = c.Pc
			if c.Pu > 0 {
				// 非零的 Pu 表示我们需要设置颜色定义。
				palette[currentPaletteIndex] = c
			}
		case b == RepeatIntroducer: // !
			data = data[:0]
			data = append(data, b)
			for {
				b, err = rd.ReadByte()
				if err != nil {
					return img, d.readError(err)
				}
				// 读取字节直到遇到非数字和非重复字节。
				if (b < '0' || b > '9') && (b < '?' || b > '~') {
					rd.UnreadByte() //nolint:errcheck,gosec
					break
				}

				data = append(data, b)
			}

			// RLE 操作
			r, n := DecodeRepeat(data)
			if n == 0 {
				return img, ErrInvalidRepeat
			}

			count = r.Count
			b = r.Char
			fallthrough
		case b >= '?' && b <= '~':
			color := palette[currentPaletteIndex]
			for range count {
				d.writePixel(currentX, currentBandY, b, color, img)
				currentX++
			}
		}
	}
}

// writePixel 接受一个六像素字节（从 ? 到 ~），定义 6 个垂直像素，
// 并将任何填充的像素写入图像。
func (d *Decoder) writePixel(x int, bandY int, sixel byte, color color.Color, img *image.RGBA) {
	maskedSixel := (sixel - '?') & 63
	yOffset := 0
	for maskedSixel != 0 {
		if maskedSixel&1 != 0 {
			img.Set(x, bandY*6+yOffset, color)
		}

		yOffset++
		maskedSixel >>= 1
	}
}

// scanSize 仅用于没有在标题附近定义像素指标的传统六像素图像（技术上是允许的）。
// 在这种情况下，我们需要快速扫描图像以确定高度和宽度。不同的终端
// 对图像边框周围未填充像素的处理方式不同，但在我们的情况下，
// 我们会将所有像素（包括空像素）视为图像的一部分。但是，
// 我们允许图像以 LF 码结尾而不增加图像大小。
//
// 为了速度起见，此方法不会以任何有意义的方式真正解析图像：像素码（? 到 ~），
// 以及 RLE、CR 和 LF 指示符（!、$, -）在六像素图像中不能以其他方式出现，
// 因此我们忽略所有其他内容。我们唯一花时间解析的是 RLE 指示符之后的数字，
// 以确定当前行要增加多少宽度。
func (d *Decoder) scanSize(data []byte) (int, int) {
	var maxWidth, bandCount int

	// 像素值从 ? 到 ~。遇到的每一个都会增加最大宽度。
	// - 表示 LF，并将最大波段数增加一。$ 表示 CR 并重置
	// 当前宽度。(char - '?') 将获得一个 6 位数，最高位是
	// 最低的 y 值，我们应该用它来增加 maxBandPixels。
	//
	// ! 是 RLE 指示符，我们应该将数字加到当前宽度
	var currentWidth int
	newBand := true
	for i := 0; i < len(data); i++ {
		b := data[i]
		switch {
		case b == LineBreak:
			// LF
			currentWidth = 0
			// 图像可能以 LF 结尾，因此我们在遇到像素之前不应增加波段数
			newBand = true
		case b == CarriageReturn:
			// CR
			currentWidth = 0
		case b == RepeatIntroducer || (b <= '~' && b >= '?'):
			count := 1
			if b == RepeatIntroducer {
				// 获取 RLE 操作的运行长度
				r, n := DecodeRepeat(data[i:])
				if n == 0 {
					return maxWidth, bandCount * 6
				}

				// 在循环中添加 1
				i += n - 1
				count = r.Count
			}

			currentWidth += count
			if newBand {
				newBand = false
				bandCount++
			}

			maxWidth = max(maxWidth, currentWidth)
		}
	}

	return maxWidth, bandCount * 6
}

// readError 接受从读取方法（ReadByte、FScanF 等）返回的任何错误，
// 并包装或忽略该错误。遇到的 EOF 表明是时候返回完成的图像了，
// 因此我们直接返回它。
func (d *Decoder) readError(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}

	return fmt.Errorf("读取六像素数据失败: %w", err)
}

// 六像素图像的默认调色板。这是六像素默认颜色和 xterm 颜色的组合。
var colorPalette = [256]color.Color{
	// Sixel-specific default colors
	0:  color.RGBA{0, 0, 0, 255},
	1:  color.RGBA{51, 51, 204, 255},
	2:  color.RGBA{204, 36, 36, 255},
	3:  color.RGBA{51, 204, 51, 255},
	4:  color.RGBA{204, 51, 204, 255},
	5:  color.RGBA{51, 204, 204, 255},
	6:  color.RGBA{204, 204, 51, 255},
	7:  color.RGBA{120, 120, 120, 255},
	8:  color.RGBA{69, 69, 69, 255},
	9:  color.RGBA{87, 87, 153, 255},
	10: color.RGBA{153, 69, 69, 255},
	11: color.RGBA{87, 153, 87, 255},
	12: color.RGBA{153, 87, 153, 255},
	13: color.RGBA{87, 153, 153, 255},
	14: color.RGBA{153, 153, 87, 255},
	15: color.RGBA{204, 204, 204, 255},

	// xterm colors
	16:  color.RGBA{0, 0, 0, 255},       // Black1
	17:  color.RGBA{0, 0, 95, 255},      // DarkBlue2
	18:  color.RGBA{0, 0, 135, 255},     // DarkBlue1
	19:  color.RGBA{0, 0, 175, 255},     // DarkBlue
	20:  color.RGBA{0, 0, 215, 255},     // Blue3
	21:  color.RGBA{0, 0, 255, 255},     // Blue2
	22:  color.RGBA{0, 95, 0, 255},      // DarkGreen4
	23:  color.RGBA{0, 95, 95, 255},     // DarkGreenBlue5
	24:  color.RGBA{0, 95, 135, 255},    // DarkGreenBlue4
	25:  color.RGBA{0, 95, 175, 255},    // DarkGreenBlue3
	26:  color.RGBA{0, 95, 215, 255},    // GreenBlue8
	27:  color.RGBA{0, 95, 255, 255},    // GreenBlue7
	28:  color.RGBA{0, 135, 0, 255},     // DarkGreen3
	29:  color.RGBA{0, 135, 95, 255},    // DarkGreen2
	30:  color.RGBA{0, 135, 0, 255},     // DarkGreenBlue2
	31:  color.RGBA{0, 135, 175, 255},   // DarkGreenBlue1
	32:  color.RGBA{0, 125, 215, 255},   // GreenBlue6
	33:  color.RGBA{0, 135, 255, 255},   // GreenBlue5
	34:  color.RGBA{0, 175, 0, 255},     // DarkGreen1
	35:  color.RGBA{0, 175, 95, 255},    // DarkGreen
	36:  color.RGBA{0, 175, 135, 255},   // DarkBlueGreen
	37:  color.RGBA{0, 175, 175, 255},   // DarkGreenBlue
	38:  color.RGBA{0, 175, 215, 255},   // GreenBlue4
	39:  color.RGBA{0, 175, 255, 255},   // GreenBlue3
	40:  color.RGBA{0, 215, 0, 255},     // Green7
	41:  color.RGBA{0, 215, 95, 255},    // Green6
	42:  color.RGBA{0, 215, 135, 255},   // Green5
	43:  color.RGBA{0, 215, 175, 255},   // BlueGreen1
	44:  color.RGBA{0, 215, 215, 255},   // GreenBlue2
	45:  color.RGBA{0, 215, 255, 255},   // GreenBlue1
	46:  color.RGBA{0, 255, 0, 255},     // Green4
	47:  color.RGBA{0, 255, 95, 255},    // Green3
	48:  color.RGBA{0, 255, 135, 255},   // Green2
	49:  color.RGBA{0, 255, 175, 255},   // Green1
	50:  color.RGBA{0, 255, 215, 255},   // BlueGreen
	51:  color.RGBA{0, 255, 255, 255},   // GreenBlue
	52:  color.RGBA{95, 0, 0, 255},      // DarkRed2
	53:  color.RGBA{95, 0, 95, 255},     // DarkPurple4
	54:  color.RGBA{95, 0, 135, 255},    // DarkBluePurple2
	55:  color.RGBA{95, 0, 175, 255},    // DarkBluePurple1
	56:  color.RGBA{95, 0, 215, 255},    // PurpleBlue
	57:  color.RGBA{95, 0, 255, 255},    // Blue1
	58:  color.RGBA{95, 95, 0, 255},     // DarkYellow4
	59:  color.RGBA{95, 95, 95, 255},    // Gray3
	60:  color.RGBA{95, 95, 135, 255},   // PlueBlue8
	61:  color.RGBA{95, 95, 175, 255},   // PaleBlue7
	62:  color.RGBA{95, 95, 215, 255},   // PaleBlue6
	63:  color.RGBA{95, 95, 255, 255},   // PaleBlue5
	64:  color.RGBA{95, 135, 0, 255},    // DarkYellow3
	65:  color.RGBA{95, 135, 95, 255},   // PaleGreen12
	66:  color.RGBA{95, 135, 135, 255},  // PaleGreen11
	67:  color.RGBA{95, 135, 175, 255},  // PaleGreenBlue10
	68:  color.RGBA{95, 135, 215, 255},  // PaleGreenBlue9
	69:  color.RGBA{95, 135, 255, 255},  // PaleBlue4
	70:  color.RGBA{95, 175, 0, 255},    // DarkGreenYellow
	71:  color.RGBA{95, 175, 95, 255},   // PaleGreen11
	72:  color.RGBA{95, 175, 135, 255},  // PaleGreen10
	73:  color.RGBA{95, 175, 175, 255},  // PaleGreenBlue8
	74:  color.RGBA{95, 175, 215, 255},  // PaleGreenBlue7
	75:  color.RGBA{95, 175, 255, 255},  // PaleGreenBlue6
	76:  color.RGBA{95, 215, 0, 255},    // YellowGreen1
	77:  color.RGBA{95, 215, 95, 255},   // PaleGreen9
	78:  color.RGBA{95, 215, 135, 255},  // PaleGreen8
	79:  color.RGBA{95, 215, 175, 255},  // PaleGreen7
	80:  color.RGBA{95, 215, 215, 255},  // PaleGreenBlue5
	81:  color.RGBA{95, 215, 255, 255},  // PaleGreenBlue4
	82:  color.RGBA{95, 255, 0, 255},    // YellowGreen
	83:  color.RGBA{95, 255, 95, 255},   // PaleGreen6
	84:  color.RGBA{95, 255, 135, 255},  // PaleGreen5
	85:  color.RGBA{95, 255, 175, 255},  // PaleGreen4
	86:  color.RGBA{95, 255, 215, 255},  // PaleGreen3
	87:  color.RGBA{95, 255, 255, 255},  // PaleGreenBlue3
	88:  color.RGBA{135, 0, 0, 255},     // DarkRed1
	89:  color.RGBA{135, 0, 95, 255},    // DarkPurple3
	90:  color.RGBA{135, 0, 135, 255},   // DarkPurple2
	91:  color.RGBA{135, 0, 175, 255},   // DarkBluePurple
	92:  color.RGBA{135, 0, 215, 255},   // BluePurple4
	93:  color.RGBA{135, 0, 255, 255},   // BluePurple3
	94:  color.RGBA{135, 95, 0, 255},    // DarkOrange1
	95:  color.RGBA{135, 95, 95, 255},   // PaleRed5
	96:  color.RGBA{135, 95, 135, 255},  // PalePurple7
	97:  color.RGBA{135, 95, 175, 255},  // PalePurpleBlue
	98:  color.RGBA{135, 95, 215, 255},  // PaleBlue3
	99:  color.RGBA{135, 95, 255, 255},  // PaleBlue2
	100: color.RGBA{135, 135, 0, 255},   // DarkYellow2
	101: color.RGBA{135, 135, 95, 255},  // PaleYellow7
	102: color.RGBA{135, 135, 135, 255}, // Gray2
	103: color.RGBA{135, 135, 175, 255}, // PaleBlue1
	104: color.RGBA{135, 135, 215, 255}, // PaleBlue
	105: color.RGBA{135, 135, 255, 255}, // LightPaleBlue4
	106: color.RGBA{135, 175, 0, 255},   // DarkYellow1
	107: color.RGBA{135, 175, 95, 255},  // PaleYellowGreen3
	108: color.RGBA{135, 175, 135, 255}, // PaleGreen2
	109: color.RGBA{135, 175, 175, 255}, // PaleGreenBlue2
	110: color.RGBA{135, 175, 215, 255}, // PaleGreenBlue1
	111: color.RGBA{135, 175, 255, 255}, // LightPaleGreenBlue6
	112: color.RGBA{135, 215, 0, 255},   // Yellow6
	113: color.RGBA{135, 215, 95, 255},  // PaleYellowGreen2
	114: color.RGBA{135, 215, 135, 255}, // PaleGreen1
	115: color.RGBA{135, 215, 175, 255}, // PaleGreen
	116: color.RGBA{135, 215, 215, 255}, // PaleGreenBlue
	117: color.RGBA{135, 215, 255, 255}, // LightPaleGreenBlue5
	118: color.RGBA{135, 255, 0, 255},   // GreenYellow
	119: color.RGBA{135, 255, 95, 255},  // PaleYellowGreen1
	120: color.RGBA{135, 255, 135, 255}, // LightPaleGreen6
	121: color.RGBA{135, 255, 175, 255}, // LightPaleGreen5
	122: color.RGBA{135, 255, 215, 255}, // LightPaleGreen4
	123: color.RGBA{135, 255, 255, 255}, // LightPaleGreenBlue4
	124: color.RGBA{175, 0, 0, 255},     // DarkRed
	125: color.RGBA{175, 0, 95, 255},    // DarkRedPurple
	126: color.RGBA{175, 0, 135, 255},   // DarkPurple1
	127: color.RGBA{175, 0, 175, 255},   // DarkPurple
	128: color.RGBA{175, 0, 215, 255},   // BluePurple2
	129: color.RGBA{175, 0, 255, 255},   // BluePurple1
	130: color.RGBA{175, 95, 0, 255},    // DarkOrange
	131: color.RGBA{175, 95, 95, 255},   // PaleRed4
	132: color.RGBA{175, 95, 135, 255},  // PalePurpleRed3
	133: color.RGBA{175, 95, 175, 255},  // PalePurple6
	134: color.RGBA{175, 95, 215, 255},  // PaleBluePurple3
	135: color.RGBA{175, 95, 255, 255},  // PaleBluePurple2
	136: color.RGBA{175, 135, 0, 255},   // DarkYellowOrange
	137: color.RGBA{175, 135, 95, 255},  // PaleRedOrange3
	138: color.RGBA{175, 135, 135, 255}, // PaleRed3
	139: color.RGBA{175, 135, 175, 255}, // PalePurple5
	140: color.RGBA{175, 135, 215, 255}, // PaleBluePurple1
	141: color.RGBA{175, 135, 255, 255}, // LightPaleBlue3
	142: color.RGBA{175, 175, 0, 255},   // DarkYellow
	143: color.RGBA{175, 175, 95, 255},  // PaleYellow6
	144: color.RGBA{175, 175, 135, 255}, // PaleYellow5
	145: color.RGBA{175, 175, 175, 255}, // Gray1
	146: color.RGBA{175, 175, 215, 255}, // LightPaleBlue2
	147: color.RGBA{175, 175, 255, 255}, // LightPaleBlue1
	148: color.RGBA{175, 215, 0, 255},   // Yellow5
	149: color.RGBA{175, 215, 95, 255},  // PaleYellow4
	150: color.RGBA{175, 215, 135, 255}, // PaleGreenYellow
	151: color.RGBA{175, 215, 175, 255}, // LightPaleGreen3
	152: color.RGBA{175, 215, 215, 255}, // LightPaleGreenBlue3
	153: color.RGBA{175, 215, 255, 255}, // LightPaleGreenBlue2
	154: color.RGBA{175, 255, 0, 255},   // Yellow4
	155: color.RGBA{175, 255, 95, 255},  // PaleYellowGreen
	156: color.RGBA{175, 255, 135, 255}, // LightPaleYellowGreen1
	157: color.RGBA{175, 255, 215, 255}, // LightPaleGreen2
	158: color.RGBA{175, 255, 215, 255}, // LightPaleGreen1
	159: color.RGBA{175, 255, 255, 255}, // LightPaleGreenBlue1
	160: color.RGBA{215, 0, 0, 255},     // Red2
	161: color.RGBA{215, 0, 95, 255},    // PurpleRed1
	162: color.RGBA{215, 0, 135, 255},   // Purple6
	163: color.RGBA{215, 0, 175, 255},   // Purple5
	164: color.RGBA{215, 0, 215, 255},   // Purple4
	165: color.RGBA{215, 0, 255, 255},   // BluePurple
	166: color.RGBA{215, 95, 0, 255},    // RedOrange1
	167: color.RGBA{215, 95, 95, 255},   // PaleRed2
	168: color.RGBA{215, 95, 135, 255},  // PalePurpleRed2
	169: color.RGBA{215, 95, 175, 255},  // PalePurple4
	170: color.RGBA{215, 95, 215, 255},  // PalePurple3
	171: color.RGBA{215, 95, 255, 255},  // PaleBluePurple
	172: color.RGBA{215, 135, 0, 255},   // Orange2
	173: color.RGBA{215, 135, 95, 255},  // PaleRedOrange2
	174: color.RGBA{215, 135, 135, 255}, // PaleRed1
	175: color.RGBA{215, 135, 175, 255}, // PaleRedPurple
	176: color.RGBA{215, 135, 215, 255}, // PalePurple2
	177: color.RGBA{215, 135, 255, 255}, // LightPaleBluePurple
	178: color.RGBA{215, 175, 0, 255},   // OrangeYellow1
	179: color.RGBA{215, 175, 95, 255},  // PaleOrange1
	180: color.RGBA{215, 175, 135, 255}, // PaleRedOrange1
	181: color.RGBA{215, 175, 175, 255}, // LightPaleRed3
	182: color.RGBA{215, 175, 215, 255}, // LightPalePurple4
	183: color.RGBA{215, 175, 255, 255}, // LightPalePurpleBlue
	184: color.RGBA{215, 215, 0, 255},   // Yellow3
	185: color.RGBA{215, 215, 95, 255},  // PaleYellow3
	186: color.RGBA{215, 215, 135, 255}, // PaleYellow2
	187: color.RGBA{215, 215, 175, 255}, // LightPaleYellow4
	188: color.RGBA{215, 215, 215, 255}, // LightGray
	189: color.RGBA{215, 215, 255, 255}, // LightPaleBlue
	190: color.RGBA{215, 255, 0, 255},   // Yellow2
	191: color.RGBA{215, 255, 95, 255},  // PaleYellow1
	192: color.RGBA{215, 255, 135, 255}, // LightPaleYellow3
	193: color.RGBA{215, 255, 175, 255}, // LightPaleYellowGreen
	194: color.RGBA{215, 255, 215, 255}, // LightPaleGreen
	195: color.RGBA{215, 255, 255, 255}, // LightPaleGreenBlue
	196: color.RGBA{255, 0, 0, 255},     // Red1
	197: color.RGBA{255, 0, 95, 255},    // PurpleRed
	198: color.RGBA{255, 0, 135, 255},   // RedPurple
	199: color.RGBA{255, 0, 175, 255},   // Purple3
	200: color.RGBA{255, 0, 215, 255},   // Purple2
	201: color.RGBA{255, 0, 255, 255},   // Purple1
	202: color.RGBA{255, 95, 0, 255},    // RedOrange
	203: color.RGBA{255, 95, 95, 255},   // PaleRed
	204: color.RGBA{255, 95, 135, 255},  // PalePurpleRed1
	205: color.RGBA{255, 95, 175, 255},  // PalePurpleRed
	206: color.RGBA{255, 95, 215, 255},  // PalePurple1
	207: color.RGBA{255, 95, 255, 255},  // PalePurple
	208: color.RGBA{255, 135, 0, 255},   // Orange1
	209: color.RGBA{255, 135, 95, 255},  // PaleOrangeRed
	210: color.RGBA{255, 135, 135, 255}, // LightPaleRed2
	211: color.RGBA{255, 135, 175, 255}, // LightPalePurpleRed1
	212: color.RGBA{255, 135, 215, 255}, // LightPalePurple3
	213: color.RGBA{255, 135, 255, 255}, // LightPalePurple2
	214: color.RGBA{255, 175, 0, 255},   // Orange
	215: color.RGBA{255, 175, 95, 255},  // PaleRedOrange
	216: color.RGBA{255, 175, 135, 255}, // LightPaleRedOrange1
	217: color.RGBA{255, 175, 175, 255}, // LightPaleRed1
	218: color.RGBA{255, 175, 215, 255}, // LightPalePurpleRed
	219: color.RGBA{255, 175, 255, 255}, // LightPalePurple1
	220: color.RGBA{255, 215, 0, 255},   // OrangeYellow
	221: color.RGBA{255, 215, 95, 255},  // PaleOrange
	222: color.RGBA{255, 215, 135, 255}, // LightPaleOrange
	223: color.RGBA{255, 215, 175, 255}, // LightPaleRedOrange
	224: color.RGBA{255, 215, 215, 255}, // LightPaleRed
	225: color.RGBA{255, 215, 255, 255}, // LightPalePurple
	226: color.RGBA{255, 255, 0, 255},   // Yellow1
	227: color.RGBA{255, 255, 95, 255},  // PaleYellow
	228: color.RGBA{255, 255, 135, 255}, // LightPaleYellow2
	229: color.RGBA{255, 255, 175, 255}, // LightPaleYellow1
	230: color.RGBA{255, 255, 215, 255}, // LightPaleYellow
	231: color.RGBA{255, 255, 255, 255}, // White1
	232: color.RGBA{8, 8, 8, 255},       // Gray4
	233: color.RGBA{18, 18, 18, 255},    // Gray8
	234: color.RGBA{28, 28, 28, 255},    // Gray11
	235: color.RGBA{38, 38, 38, 255},    // Gray15
	236: color.RGBA{48, 48, 48, 255},    // Gray19
	237: color.RGBA{58, 58, 58, 255},    // Gray23
	238: color.RGBA{68, 68, 68, 255},    // Gray27
	239: color.RGBA{78, 78, 78, 255},    // Gray31
	240: color.RGBA{88, 88, 88, 255},    // Gray35
	241: color.RGBA{98, 98, 98, 255},    // Gray39
	242: color.RGBA{108, 108, 108, 255}, // Gray43
	243: color.RGBA{118, 118, 118, 255}, // Gray47
	244: color.RGBA{128, 128, 128, 255}, // Gray51
	245: color.RGBA{138, 138, 138, 255}, // Gray55
	246: color.RGBA{148, 148, 148, 255}, // Gray59
	247: color.RGBA{158, 158, 158, 255}, // Gray62
	248: color.RGBA{168, 168, 168, 255}, // Gray66
	249: color.RGBA{178, 178, 178, 255}, // Gray70
	250: color.RGBA{188, 188, 188, 255}, // Gray74
	251: color.RGBA{198, 198, 198, 255}, // Gray78
	252: color.RGBA{208, 208, 208, 255}, // Gray82
	253: color.RGBA{218, 218, 218, 255}, // Gray86
	254: color.RGBA{228, 228, 228, 255}, // Gray90
	255: color.RGBA{238, 238, 238, 255}, // Gray94
}

// DefaultPalette is the default palette used when decoding a Sixel image.
// It contains the 256 colors defined by the xterm 256-color palette.
func DefaultPalette() color.Palette {
	// Undefined colors in sixel images use a set of default colors: 0-15
	// are sixel-specific, 16-255 are the same as the xterm 256-color values
	palette := slices.Clone(colorPalette[:])
	return palette[:]
}
