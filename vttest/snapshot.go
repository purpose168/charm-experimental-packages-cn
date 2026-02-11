package vttest

import (
	"fmt"
	"image/color"
	"strconv"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/vt"
)

// Modes 表示终端模式。
type Modes struct {
	ANSI map[ansi.ANSIMode]ansi.ModeSetting `json:"ansi" yaml:"ansi"`
	DEC  map[ansi.DECMode]ansi.ModeSetting  `json:"dec" yaml:"dec"`
}

// Position 表示终端中的一个位置。
type Position struct {
	X int `json:"x" yaml:"x"`
	Y int `json:"y" yaml:"y"`
}

// Color 表示终端颜色，可以是以下之一：
// - 类型为 [ansi.BasicColor] 的 ANSI 16 色（0-15）。
// - 类型为 [ansi.IndexedColor] 的 ANSI 256 色（0-255）。
// - 或任何其他实现了 [color.Color] 的 24 位颜色。
type Color struct {
	Color color.Color `json:"color,omitempty" yaml:"color,omitempty"`
}

// MarshalText 为 Color 实现 [encoding.TextMarshaler] 接口。
func (c Color) MarshalText() ([]byte, error) {
	switch col := c.Color.(type) {
	case nil:
		return []byte{}, nil
	case ansi.BasicColor:
		return []byte(strconv.Itoa(int(col))), nil
	case ansi.IndexedColor:
		return []byte(strconv.Itoa(int(col))), nil
	default:
		r, g, b, _ := c.Color.RGBA()
		return fmt.Appendf(nil, "#%02x%02x%02x", r>>8, g>>8, b>>8), nil
	}
}

// UnmarshalText 为 Color 实现 [encoding.TextUnmarshaler] 接口。
func (c *Color) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" {
		return nil
	}
	if i, err := strconv.Atoi(s); err == nil {
		if i >= 0 && i <= 15 {
			c.Color = ansi.BasicColor(i)
			return nil
		} else if i >= 16 && i <= 255 {
			c.Color = ansi.IndexedColor(i)
			return nil
		}
	}

	col := ansi.XParseColor(s)
	if col == nil {
		return fmt.Errorf("invalid color: %s", s)
	}
	c.Color = col
	return nil
}

// Cursor 表示光标的状态。
type Cursor struct {
	Position Position       `json:"position," yaml:"position"`
	Visible  bool           `json:"visible" yaml:"visible"`
	Color    Color          `json:"color,omitzero" yaml:"color,omitzero"`
	Style    vt.CursorStyle `json:"style" yaml:"style"`
	Blink    bool           `json:"blink" yaml:"blink"`
}

// Style 表示单元格的样式。
type Style struct {
	Fg             Color        `json:"fg,omitzero" yaml:"fg,omitzero"`
	Bg             Color        `json:"bg,omitzero" yaml:"bg,omitzero"`
	UnderlineColor Color        `json:"underline_color,omitzero" yaml:"underline_color,omitzero"`
	Underline      uv.Underline `json:"underline,omitempty" yaml:"underline,omitempty"`
	Attrs          byte         `json:"attrs,omitempty" yaml:"attrs,omitempty"`
}

// Link 表示终端屏幕中的超链接。
type Link struct {
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`
	Params string `json:"params,omitempty" yaml:"params,omitempty"`
}

// Cell 表示终端屏幕中的单个单元格。
type Cell struct {
	// Content 是单元格的内容，由单个字形簇组成。大多数情况下，这也是单个符文，但也可以是形成字形簇的多个符文的组合。
	Content string `json:"content,omitempty" yaml:"content,omitempty"`

	// Style 是单元格的样式。零值表示重置序列。
	Style Style `json:"style,omitzero" yaml:"style,omitzero"`

	// Link 是单元格的超链接。
	Link Link `json:"link,omitzero" yaml:"link,omitzero"`

	// Width 是字形簇的等宽宽度。
	Width int `json:"width,omitzero" yaml:"width,omitzero"`
}

// Snapshot 表示给定时刻的终端状态快照。
type Snapshot struct {
	Modes     Modes    `json:"modes" yaml:"modes"`
	Title     string   `json:"title" yaml:"title"`
	Rows      int      `json:"rows" yaml:"rows"`
	Cols      int      `json:"cols" yaml:"cols"`
	AltScreen bool     `json:"alt_screen" yaml:"alt_screen"`
	Cursor    Cursor   `json:"cursor" yaml:"cursor"`
	BgColor   Color    `json:"bg_color,omitzero" yaml:"bg_color,omitzero"`
	FgColor   Color    `json:"fg_color,omitzero" yaml:"fg_color,omitzero"`
	Cells     [][]Cell `json:"cells" yaml:"cells"`
}
