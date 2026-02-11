// Package iterm2 提供 iTerm2 特定的功能。
package iterm2

import (
	"strconv"
	"strings"
)

// Auto 是表示 "auto" 值的常量。
const Auto = "auto"

// Cells 返回表示单元格数量的字符串。这只是 strconv.Itoa 的包装器。
func Cells(n int) string {
	return strconv.Itoa(n)
}

// Pixels 返回表示像素数量的字符串。
func Pixels(n int) string {
	return strconv.Itoa(n) + "px"
}

// Percent 返回表示百分比的字符串。
func Percent(n int) string {
	return strconv.Itoa(n) + "%"
}

// file 表示 iTerm2 内联图像协议的可选参数。
//
// 请参阅 https://iterm2.com/documentation-images.html
type file struct {
	// Name 是文件的名称。如果为空，默认为 "Unnamed file"。
	Name string
	// Size 是文件大小（以字节为单位）。用于进度指示。这是可选的。
	Size int64
	// Width 是图像的宽度。这可以通过一个数字后跟单位或 "auto" 来指定。
	// 单位可以是无、"px" 或 "%". 无表示数字以单元格为单位。
	// 如果为空，默认为 "auto"。
	// 为方便起见，可以使用 [Pixels]、[Cells] 和 [Percent] 函数以及 [Auto]。
	Width string
	// Height 是图像的高度。这可以通过一个数字后跟单位或 "auto" 来指定。
	// 单位可以是无、"px" 或 "%". 无表示数字以单元格为单位。
	// 如果为空，默认为 "auto"。
	// 为方便起见，可以使用 [Pixels]、[Cells] 和 [Percent] 函数以及 [Auto]。
	Height string
	// IgnoreAspectRatio 是一个标志，指示图像是否应该被拉伸以适应指定的宽度和高度。
	// 默认为 false，表示保留宽高比。
	IgnoreAspectRatio bool
	// Inline 是一个标志，指示图像是否应该内联显示。
	// 否则，它会被下载到 Downloads 文件夹，在终端中没有视觉表示。默认为 false。
	Inline bool
	// DoNotMoveCursor 是一个标志，指示显示图像后光标是否不应移动。
	// 这是 WezTerm 引入的扩展，可能不适用于所有支持 iTerm2 协议的终端。
	// 默认为 false。
	DoNotMoveCursor bool
	// Content 是文件的 base64 编码数据。
	Content []byte
}

// String 实现 fmt.Stringer 接口。
func (f file) String() string {
	var opts []string
	if f.Name != "" {
		opts = append(opts, "name="+f.Name)
	}
	if f.Size != 0 {
		opts = append(opts, "size="+strconv.FormatInt(f.Size, 10))
	}
	if f.Width != "" {
		opts = append(opts, "width="+f.Width)
	}
	if f.Height != "" {
		opts = append(opts, "height="+f.Height)
	}
	if f.IgnoreAspectRatio {
		opts = append(opts, "preserveAspectRatio=0")
	}
	if f.Inline {
		opts = append(opts, "inline=1")
	}
	if f.DoNotMoveCursor {
		opts = append(opts, "doNotMoveCursor=1")
	}
	return strings.Join(opts, ";")
}

// File 表示 iTerm2 内联图像协议的可选参数。
type File file

// String 实现 fmt.Stringer 接口。
func (f File) String() string {
	var s strings.Builder
	s.WriteString("File=")
	s.WriteString(file(f).String())
	if len(f.Content) > 0 {
		s.WriteString(":")
		s.Write(f.Content)
	}

	return s.String()
}

// MultipartFile 表示 iTerm2 内联图像协议的可选参数。
type MultipartFile file

// String 实现 fmt.Stringer 接口。
func (f MultipartFile) String() string {
	return "MultipartFile=" + file(f).String()
}

// FilePart 表示 iTerm2 内联图像协议的可选参数。
type FilePart file

// String 实现 fmt.Stringer 接口。
func (f FilePart) String() string {
	return "FilePart=" + string(f.Content)
}

// FileEnd 表示 iTerm2 内联图像协议的可选参数。
type FileEnd struct{}

// String 实现 fmt.Stringer 接口。
func (f FileEnd) String() string {
	return "FileEnd"
}
