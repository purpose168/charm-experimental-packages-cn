package cellbuf

import (
	"bytes"
	"image/color"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// ReadStyle 从参数列表中读取选择图形渲染（Select Graphic Rendition，SGR）转义序列。
func ReadStyle(params ansi.Params, pen *Style) {
	if len(params) == 0 {
		pen.Reset()
		return
	}

	for i := 0; i < len(params); i++ {
		param, hasMore, _ := params.Param(i, 0)
		switch param {
		case 0: // 重置
			pen.Reset()
		case 1: // 粗体
			pen.Bold(true)
		case 2: // 暗淡/微弱
			pen.Faint(true)
		case 3: // 斜体
			pen.Italic(true)
		case 4: // 下划线
			nextParam, _, ok := params.Param(i+1, 0)
			if hasMore && ok { // 仅接受子参数，即由":"分隔的参数
				switch nextParam {
				case 0, 1, 2, 3, 4, 5:
					i++
					switch nextParam {
					case 0: // 无下划线
						pen.UnderlineStyle(NoUnderline)
					case 1: // 单下划线
						pen.UnderlineStyle(SingleUnderline)
					case 2: // 双下划线
						pen.UnderlineStyle(DoubleUnderline)
					case 3: // 卷曲下划线
						pen.UnderlineStyle(CurlyUnderline)
					case 4: // 点下划线
						pen.UnderlineStyle(DottedUnderline)
					case 5: // 虚线下划线
						pen.UnderlineStyle(DashedUnderline)
					}
				}
			} else {
				// 单下划线
				pen.Underline(true)
			}
		case 5: // 慢速闪烁
			pen.SlowBlink(true)
		case 6: // 快速闪烁
			pen.RapidBlink(true)
		case 7: // 反显
			pen.Reverse(true)
		case 8: // 隐藏
			pen.Conceal(true)
		case 9: // 删除线
			pen.Strikethrough(true)
		case 22: // 正常强度（非粗体或暗淡）
			pen.Bold(false).Faint(false)
		case 23: // 非斜体，非哥特体
			pen.Italic(false)
		case 24: // 无下划线
			pen.Underline(false)
		case 25: // 关闭闪烁
			pen.SlowBlink(false).RapidBlink(false)
		case 27: // 正向（非反显）
			pen.Reverse(false)
		case 28: // 显示
			pen.Conceal(false)
		case 29: // 无删除线
			pen.Strikethrough(false)
		case 30, 31, 32, 33, 34, 35, 36, 37: // 设置前景色
			pen.Foreground(ansi.Black + ansi.BasicColor(param-30)) //nolint:gosec
		case 38: // 设置前景色 256 或真彩色
			var c color.Color
			n := ReadStyleColor(params[i:], &c)
			if n > 0 {
				pen.Foreground(c)
				i += n - 1
			}
		case 39: // 默认前景色
			pen.Foreground(nil)
		case 40, 41, 42, 43, 44, 45, 46, 47: // 设置背景色
			pen.Background(ansi.Black + ansi.BasicColor(param-40)) //nolint:gosec
		case 48: // 设置背景色 256 或真彩色
			var c color.Color
			n := ReadStyleColor(params[i:], &c)
			if n > 0 {
				pen.Background(c)
				i += n - 1
			}
		case 49: // 默认背景色
			pen.Background(nil)
		case 58: // 设置下划线颜色
			var c color.Color
			n := ReadStyleColor(params[i:], &c)
			if n > 0 {
				pen.UnderlineColor(c)
				i += n - 1
			}
		case 59: // 默认下划线颜色
			pen.UnderlineColor(nil)
		case 90, 91, 92, 93, 94, 95, 96, 97: // 设置亮前景色
			pen.Foreground(ansi.BrightBlack + ansi.BasicColor(param-90)) //nolint:gosec
		case 100, 101, 102, 103, 104, 105, 106, 107: // 设置亮背景色
			pen.Background(ansi.BrightBlack + ansi.BasicColor(param-100)) //nolint:gosec
		}
	}
}

// ReadLink 从数据缓冲区中读取超链接转义序列。
func ReadLink(p []byte, link *Link) {
	params := bytes.Split(p, []byte{';'})
	if len(params) != 3 {
		return
	}
	link.Params = string(params[1])
	link.URL = string(params[2])
}

// ReadStyleColor 从参数列表中读取颜色。
// 有关更多信息，请参见 [ansi.ReadStyleColor]。
func ReadStyleColor(params ansi.Params, c *color.Color) int {
	return ansi.ReadStyleColor(params, c)
}
