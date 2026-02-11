package kitty

import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

var (
	_ encoding.TextMarshaler   = Options{}
	_ encoding.TextUnmarshaler = &Options{}
)

type Stringish interface{ string | []byte }

// Options 表示 Kitty 图形协议的选项。
type Options struct {
	// 通用选项。

	// Action (a=t) 是要对图像执行的操作。可以是以下之一：
	// [Transmit], [TransmitDisplay], [Query], [Put], [Delete], [Frame],
	// [Animate], [Compose]。
	Action byte

	// Quite mode (q=0) 是安静模式。可以是 0、1 或 2，
	// 其中 0 是默认值，1 抑制 OK 响应，2 抑制 OK 和错误响应。
	Quite byte

	// 传输选项。

	// ID (i=) 是图像 ID。ID 是图像的唯一标识符。
	// 必须是不超过 [math.MaxUint32] 的正整数。
	ID int

	// PlacementID (p=) 是放置 ID。放置 ID 是图像放置的唯一标识符。
	// 必须是不超过 [math.MaxUint32] 的正整数。
	PlacementID int

	// Number (I=0) 是要传输的图像数量。
	Number int

	// Format (f=32) 是图像格式。以下之一：[RGBA], [RGB], [PNG]。
	Format int

	// ImageWidth (s=0) 是传输的图像宽度。
	ImageWidth int

	// ImageHeight (v=0) 是传输的图像高度。
	ImageHeight int

	// Compression (o=) 是图像压缩类型。可以是 [Zlib] 或 0。
	Compression byte

	// Transmission (t=d) 是图像传输类型。可以是 [Direct], [File],
	// [TempFile], 或 [SharedMemory]。
	Transmission byte

	// File 是当传输类型为 [File] 时要使用的文件路径。
	// 如果 [Options.Transmission] 被省略（即 0）且此值非空，
	// 传输类型将设置为 [File]。
	File string

	// Size (S=0) 是要从传输介质读取的大小。
	Size int

	// Offset (O=0) 是开始从传输介质读取的偏移字节。
	Offset int

	// Chunk (m=) 表示图像是否以块方式传输。可以是 0 或 1。
	// 当为 true 时，图像以块方式传输。每个块必须是 4 的倍数，
	// 且最多 [MaxChunkSize] 字节。除最后一个块必须有 m=0 选项外，
	// 每个块都必须有 m=1 选项。
	Chunk bool

	// ChunkFormatter 是当 [Options.Chunk] 为 true 时用于格式化每个块的函数。
	// 如果为 nil，则块按原样发送。
	ChunkFormatter func(chunk string) string

	// 显示选项。

	// X (x=0) 是图像开始显示的像素 X 坐标。
	X int

	// Y (y=0) 是图像开始显示的像素 Y 坐标。
	Y int

	// Z (z=0) 是要显示的图像的 Z 坐标。
	Z int

	// Width (w=0) 是要显示的图像宽度。
	Width int

	// Height (h=0) 是要显示的图像高度。
	Height int

	// OffsetX (X=0) 是开始显示图像的光标单元格的 OffsetX 坐标。
	// OffsetX=0 是最左侧的单元格。这必须小于终端单元格宽度。
	OffsetX int

	// OffsetY (Y=0) 是开始显示图像的光标单元格的 OffsetY 坐标。
	// OffsetY=0 是最顶部的单元格。这必须小于终端单元格高度。
	OffsetY int

	// Columns (c=0) 是显示图像的列数。图像将被缩放以适应列数。
	Columns int

	// Rows (r=0) 是显示图像的行数。图像将被缩放以适应行数。
	Rows int

	// VirtualPlacement (U=0) 是否使用虚拟放置。这与 Unicode [Placeholder] 一起使用以显示图像。
	VirtualPlacement bool

	// DoNotMoveCursor (C=0) 是否在显示图像后移动光标。
	DoNotMoveCursor bool

	// ParentID (P=0) 是父图像 ID。父 ID 是当前图像的父图像的 ID。
	// 这与 Unicode [Placeholder] 一起使用以相对于父图像显示图像。
	ParentID int

	// ParentPlacementID (Q=0) 是父放置 ID。父放置 ID 是父图像放置的 ID。
	// 这与 Unicode [Placeholder] 一起使用以相对于父图像显示图像。
	ParentPlacementID int

	// 删除选项。

	// Delete (d=a) 是删除操作。可以是以下之一：[DeleteAll],
	// [DeleteID], [DeleteNumber], [DeleteCursor], [DeleteFrames],
	// [DeleteCell], [DeleteCellZ], [DeleteRange], [DeleteColumn], [DeleteRow],
	// [DeleteZ]。
	Delete byte

	// DeleteResources 指示是否删除与图像关联的资源。
	DeleteResources bool
}

// Options 返回作为键值对切片的选项。
func (o *Options) Options() (opts []string) {
	opts = []string{}
	if o.Format == 0 {
		o.Format = RGBA
	}

	if o.Action == 0 {
		o.Action = Transmit
	}

	if o.Delete == 0 {
		o.Delete = DeleteAll
	}

	if o.Transmission == 0 {
		if len(o.File) > 0 {
			o.Transmission = File
		} else {
			o.Transmission = Direct
		}
	}

	if o.Format != RGBA {
		opts = append(opts, fmt.Sprintf("f=%d", o.Format))
	}

	if o.Quite > 0 {
		opts = append(opts, fmt.Sprintf("q=%d", o.Quite))
	}

	if o.ID > 0 {
		opts = append(opts, fmt.Sprintf("i=%d", o.ID))
	}

	if o.PlacementID > 0 {
		opts = append(opts, fmt.Sprintf("p=%d", o.PlacementID))
	}

	if o.Number > 0 {
		opts = append(opts, fmt.Sprintf("I=%d", o.Number))
	}

	if o.ImageWidth > 0 {
		opts = append(opts, fmt.Sprintf("s=%d", o.ImageWidth))
	}

	if o.ImageHeight > 0 {
		opts = append(opts, fmt.Sprintf("v=%d", o.ImageHeight))
	}

	if o.Transmission != Direct {
		opts = append(opts, fmt.Sprintf("t=%c", o.Transmission))
	}

	if o.Size > 0 {
		opts = append(opts, fmt.Sprintf("S=%d", o.Size))
	}

	if o.Offset > 0 {
		opts = append(opts, fmt.Sprintf("O=%d", o.Offset))
	}

	if o.Compression == Zlib {
		opts = append(opts, fmt.Sprintf("o=%c", o.Compression))
	}

	if o.VirtualPlacement {
		opts = append(opts, "U=1")
	}

	if o.DoNotMoveCursor {
		opts = append(opts, "C=1")
	}

	if o.ParentID > 0 {
		opts = append(opts, fmt.Sprintf("P=%d", o.ParentID))
	}

	if o.ParentPlacementID > 0 {
		opts = append(opts, fmt.Sprintf("Q=%d", o.ParentPlacementID))
	}

	if o.X > 0 {
		opts = append(opts, fmt.Sprintf("x=%d", o.X))
	}

	if o.Y > 0 {
		opts = append(opts, fmt.Sprintf("y=%d", o.Y))
	}

	if o.Z > 0 {
		opts = append(opts, fmt.Sprintf("z=%d", o.Z))
	}

	if o.Width > 0 {
		opts = append(opts, fmt.Sprintf("w=%d", o.Width))
	}

	if o.Height > 0 {
		opts = append(opts, fmt.Sprintf("h=%d", o.Height))
	}

	if o.OffsetX > 0 {
		opts = append(opts, fmt.Sprintf("X=%d", o.OffsetX))
	}

	if o.OffsetY > 0 {
		opts = append(opts, fmt.Sprintf("Y=%d", o.OffsetY))
	}

	if o.Columns > 0 {
		opts = append(opts, fmt.Sprintf("c=%d", o.Columns))
	}

	if o.Rows > 0 {
		opts = append(opts, fmt.Sprintf("r=%d", o.Rows))
	}

	if o.Delete != DeleteAll || o.DeleteResources {
		da := o.Delete
		if o.DeleteResources {
			da = da - ' ' // 转为大写
		}

		opts = append(opts, fmt.Sprintf("d=%c", da))
	}

	if o.Action != Transmit {
		opts = append(opts, fmt.Sprintf("a=%c", o.Action))
	}

	return opts // 具有多个返回值的复杂函数
}

// String 返回选项的字符串表示。
func (o Options) String() string {
	return strings.Join(o.Options(), ",")
}

// MarshalText 返回选项的字符串表示。
func (o Options) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

// UnmarshalText 从给定字符串解析选项。
func (o *Options) UnmarshalText(text []byte) error {
	opts := strings.Split(string(text), ",")
	for _, opt := range opts {
		ps := strings.SplitN(opt, "=", 2)
		if len(ps) != 2 || len(ps[1]) == 0 {
			continue
		}

		switch ps[0] {
		case "a":
			o.Action = ps[1][0]
		case "o":
			o.Compression = ps[1][0]
		case "t":
			o.Transmission = ps[1][0]
		case "d":
			d := ps[1][0]
			if d >= 'A' && d <= 'Z' {
				o.DeleteResources = true
				d = d + ' ' // 转为小写
			}
			o.Delete = d
		case "i", "q", "p", "I", "f", "s", "v", "S", "O", "m", "x", "y", "z", "w", "h", "X", "Y", "c", "r", "U", "P", "Q":
			v, err := strconv.Atoi(ps[1])
			if err != nil {
				continue
			}

			switch ps[0] {
			case "i":
				o.ID = v
			case "q":
				o.Quite = byte(v)
			case "p":
				o.PlacementID = v
			case "I":
				o.Number = v
			case "f":
				o.Format = v
			case "s":
				o.ImageWidth = v
			case "v":
				o.ImageHeight = v
			case "S":
				o.Size = v
			case "O":
				o.Offset = v
			case "m":
				o.Chunk = v == 0 || v == 1
			case "x":
				o.X = v
			case "y":
				o.Y = v
			case "z":
				o.Z = v
			case "w":
				o.Width = v
			case "h":
				o.Height = v
			case "X":
				o.OffsetX = v
			case "Y":
				o.OffsetY = v
			case "c":
				o.Columns = v
			case "r":
				o.Rows = v
			case "U":
				o.VirtualPlacement = v == 1
			case "P":
				o.ParentID = v
			case "Q":
				o.ParentPlacementID = v
			}
		}
	}

	return nil
}
