package cellbuf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/colorprofile"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/term"
)

// ErrInvalidDimensions 当窗口尺寸对于操作无效时返回
var ErrInvalidDimensions = errors.New("invalid dimensions")

// notLocal 使用定义的阈值返回坐标是否不被视为本地移动
// 接收列数以及当前和目标位置的坐标
func notLocal(cols, fx, fy, tx, ty int) bool {
	// [ansi.CUP] 序列的典型距离。小于此值的被视为本地移动
	const longDist = 8 - 1
	return (tx > longDist) &&
		(tx < cols-1-longDist) &&
		(abs(ty-fy)+abs(tx-fx) > longDist)
}

// relativeCursorMove 返回使用以下序列之一或两个的相对光标移动序列：
// [ansi.CUU], [ansi.CUD], [ansi.CUF], [ansi.CUB], [ansi.VPA], [ansi.HPA]
// 当 overwrite 为 true 时，会尝试通过使用屏幕单元格值而不是转义序列来优化移动光标
func relativeCursorMove(s *Screen, fx, fy, tx, ty int, overwrite, useTabs, useBackspace bool) string {
	var seq strings.Builder

	width, height := s.newbuf.Width(), s.newbuf.Height()
	if ty != fy { //nolint:nestif
		var yseq string
		if s.caps.Contains(capVPA) && !s.opts.RelativeCursor {
			yseq = ansi.VerticalPositionAbsolute(ty + 1)
		}

		// 优化：使用 [ansi.LF] 和 [ansi.ReverseIndex] 作为优化

		if ty > fy {
			n := ty - fy
			if cud := ansi.CursorDown(n); yseq == "" || len(cud) < len(yseq) {
				yseq = cud
			}
			shouldScroll := !s.opts.AltScreen && fy+n >= s.scrollHeight
			if lf := strings.Repeat("\n", n); shouldScroll || (fy+n < height && len(lf) < len(yseq)) {
				// TODO: 确保我们不会无意中向下滚动屏幕
				yseq = lf
				s.scrollHeight = max(s.scrollHeight, fy+n)
				if s.opts.MapNL {
					fx = 0
				}
			}
		} else if ty < fy {
			n := fy - ty
			if cuu := ansi.CursorUp(n); yseq == "" || len(cuu) < len(yseq) {
				yseq = cuu
			}
			if n == 1 && fy-1 > 0 {
				// TODO: 确保我们不会无意中向上滚动屏幕
				yseq = ansi.ReverseIndex
			}
		}

		seq.WriteString(yseq)
	}

	if tx != fx { //nolint:nestif
		var xseq string
		if s.caps.Contains(capHPA) && !s.opts.RelativeCursor {
			xseq = ansi.HorizontalPositionAbsolute(tx + 1)
		}

		if tx > fx {
			n := tx - fx
			if useTabs {
				var tabs int
				var col int
				for col = fx; s.tabs.Next(col) <= tx; col = s.tabs.Next(col) {
					tabs++
					if col == s.tabs.Next(col) || col >= width-1 {
						break
					}
				}

				if tabs > 0 {
					cht := ansi.CursorHorizontalForwardTab(tabs)
					tab := strings.Repeat("\t", tabs)
					if false && s.caps.Contains(capCHT) && len(cht) < len(tab) {
						// Linux 控制台和一些终端如 Alacritty 不支持 [ansi.CHT]
						// 当我们有方法检测到这一点时，或者 5 年后当我们确定每个人都更新了终端时再启用 :P
						seq.WriteString(cht)
					} else {
						seq.WriteString(tab)
					}

					n = tx - col
					fx = col
				}
			}

			if cuf := ansi.CursorForward(n); xseq == "" || len(cuf) < len(xseq) {
				xseq = cuf
			}

			// 如果没有属性和样式更改，覆盖更便宜
			var ovw string
			if overwrite && ty >= 0 {
				for i := 0; i < n; i++ {
					cell := s.newbuf.Cell(fx+i, ty)
					if cell != nil && cell.Width > 0 {
						i += cell.Width - 1
						if !cell.Style.Equal(&s.cur.Style) || !cell.Link.Equal(&s.cur.Link) {
							overwrite = false
							break
						}
					}
				}
			}

			if overwrite && ty >= 0 {
				for i := 0; i < n; i++ {
					cell := s.newbuf.Cell(fx+i, ty)
					if cell != nil && cell.Width > 0 {
						ovw += cell.String()
						i += cell.Width - 1
					} else {
						ovw += " "
					}
				}
			}

			if overwrite && len(ovw) < len(xseq) {
				xseq = ovw
			}
		} else if tx < fx {
			n := fx - tx
			if useTabs && s.caps.Contains(capCBT) {
				// VT100 不支持向后制表符 [ansi.CBT]

				col := fx

				var cbt int // 光标向后制表符计数
				for s.tabs.Prev(col) >= tx {
					col = s.tabs.Prev(col)
					cbt++
					if col == s.tabs.Prev(col) || col <= 0 {
						break
					}
				}

				if cbt > 0 {
					seq.WriteString(ansi.CursorBackwardTab(cbt))
					n = col - tx
				}
			}

			if cub := ansi.CursorBackward(n); xseq == "" || len(cub) < len(xseq) {
				xseq = cub
			}

			if useBackspace && n < len(xseq) {
				xseq = strings.Repeat("\b", n)
			}
		}

		seq.WriteString(xseq)
	}

	return seq.String()
}

// moveCursor 移动并返回将光标移动到指定位置的光标移动序列
// 当 overwrite 为 true 时，会尝试通过使用屏幕单元格值而不是转义序列来优化移动光标
func moveCursor(s *Screen, x, y int, overwrite bool) (seq string) {
	fx, fy := s.cur.X, s.cur.Y

	if !s.opts.RelativeCursor {
		// 方法 #0：如果距离较长，使用 [ansi.CUP]
		seq = ansi.CursorPosition(x+1, y+1)
		if fx == -1 || fy == -1 || notLocal(s.newbuf.Width(), fx, fy, x, y) {
			return seq
		}
	}

	// 根据选项进行优化
	trials := 0
	if s.opts.HardTabs {
		trials |= 2 // 二进制 0b10
	}
	if s.opts.Backspace {
		trials |= 1 // 二进制 0b01
	}

	// 尝试硬制表符和退格键优化的所有可能组合
	for i := 0; i <= trials; i++ {
		// 跳过未启用的组合
		if i & ^trials != 0 {
			continue
		}

		useHardTabs := i&2 != 0
		useBackspace := i&1 != 0

		// 方法 #1：使用本地移动序列
		nseq := relativeCursorMove(s, fx, fy, x, y, overwrite, useHardTabs, useBackspace)
		if (i == 0 && len(seq) == 0) || len(nseq) < len(seq) {
			seq = nseq
		}

		// 方法 #2：使用 [ansi.CR] 和本地移动序列
		nseq = "\r" + relativeCursorMove(s, 0, fy, x, y, overwrite, useHardTabs, useBackspace)
		if len(nseq) < len(seq) {
			seq = nseq
		}

		if !s.opts.RelativeCursor {
			// 方法 #3：使用 [ansi.CursorHomePosition] 和本地移动序列
			nseq = ansi.CursorHomePosition + relativeCursorMove(s, 0, 0, x, y, overwrite, useHardTabs, useBackspace)
			if len(nseq) < len(seq) {
				seq = nseq
			}
		}
	}

	return seq
}

// moveCursor 将光标移动到指定位置
func (s *Screen) moveCursor(x, y int, overwrite bool) {
	if !s.opts.AltScreen && s.cur.X == -1 && s.cur.Y == -1 {
		// 内联模式下的第一次光标移动，在移动到目标位置之前，先将光标移动到第一列
		s.buf.WriteByte('\r')
		s.cur.X, s.cur.Y = 0, 0
	}
	s.buf.WriteString(moveCursor(s, x, y, overwrite))
	s.cur.X, s.cur.Y = x, y
}

func (s *Screen) move(x, y int) {
	// 确保在调整大小操作过程中使用缓冲区的最大高度和宽度
	width := max(s.newbuf.Width(), s.curbuf.Width())
	height := max(s.newbuf.Height(), s.curbuf.Height())

	if width > 0 && x >= width {
		// 处理自动换行
		y += (x / width)
		x %= width
	}

	// 如果有样式，禁用它们
	// 一些移动操作如 [ansi.LF] 会将样式应用到新的光标位置，因此我们需要在移动光标之前重置样式
	blank := s.clearBlank()
	resetPen := y != s.cur.Y && !blank.Equal(&BlankCell)
	if resetPen {
		s.updatePen(nil)
	}

	// 重置环绕（幻象光标）状态
	if s.atPhantom {
		s.cur.X = 0
		s.buf.WriteByte('\r')
		s.atPhantom = false // 重置幻象单元格状态
	}

	// TODO: 调查我们是否需要处理这种情况和/或是否需要以下代码
	//
	// if width > 0 && s.cur.X >= width {
	// 	l := (s.cur.X + 1) / width
	//
	// 	s.cur.Y += l
	// 	if height > 0 && s.cur.Y >= height {
	// 		l -= s.cur.Y - height - 1
	// 	}
	//
	// 	if l > 0 {
	// 		s.cur.X = 0
	// 		s.buf.WriteString("\r" + strings.Repeat("\n", l))
	// 	}
	// }

	if height > 0 {
		if s.cur.Y > height-1 {
			s.cur.Y = height - 1
		}
		if y > height-1 {
			y = height - 1
		}
	}

	if x == s.cur.X && y == s.cur.Y {
		// 我们稍后放弃，因为我们需要运行幻象单元格和其他检查，然后才能确定是否可以放弃
		return
	}

	// 我们在 [Screen.moveCursor] 中设置新光标
	s.moveCursor(x, y, true) // 尽可能覆盖单元格
}

// Cursor 表示终端光标
type Cursor struct {
	Style
	Link
	Position
}

// ScreenOptions 是屏幕的选项
type ScreenOptions struct {
	// Term 是写入屏幕时使用的终端类型。为空时，使用 [os.Getenv] 中的 `$TERM`
	Term string
	// Profile 是写入屏幕时使用的颜色配置文件
	Profile colorprofile.Profile
	// RelativeCursor 是否使用相对光标移动。当不使用替代屏幕或使用内联模式时很有用
	RelativeCursor bool
	// AltScreen 是否使用替代屏幕缓冲区
	AltScreen bool
	// ShowCursor 是否显示光标
	ShowCursor bool
	// HardTabs 是否使用硬制表符来优化光标移动
	HardTabs bool
	// Backspace 是否使用退格字符来移动光标
	Backspace bool
	// MapNL 是否启用了 ONLCR 映射。当我们将终端设置为原始模式时，ONLCR 模式会被禁用
	// ONLCR 将任何换行/换行符 (`\n`) 映射为回车 + 换行 (`\r\n`)
	MapNL bool
}

// lineData 表示一行的元数据
type lineData struct {
	// 第一个和最后一个更改的单元格索引
	firstCell, lastCell int
	// 用于滚动的旧索引
	oldIndex int //nolint:unused
}

// Screen 表示终端屏幕
type Screen struct {
	w                io.Writer
	buf              *bytes.Buffer // 用于写入屏幕的缓冲区
	curbuf           *Buffer       // 当前缓冲区
	newbuf           *Buffer       // 新缓冲区
	tabs             *TabStops
	touch            map[int]lineData
	queueAbove       []string  // 要在屏幕上方写入的字符串队列
	oldhash, newhash []uint64  // 每行的旧和新哈希值
	hashtab          []hashmap // 哈希表
	oldnum           []int     // 来自先前哈希的旧索引
	cur, saved       Cursor    // 当前和保存的光标
	opts             ScreenOptions
	mu               sync.Mutex
	method           ansi.Method
	scrollHeight     int          // 跟踪我们向下滚动的行数（内联模式）
	altScreenMode    bool         // 是否启用了替代屏幕模式
	cursorHidden     bool         // 是否启用了文本光标模式
	clear            bool         // 是否强制清除屏幕
	caps             capabilities // 终端控制序列功能
	queuedText       bool         // 是否有非零宽度文本排队
	atPhantom        bool         // 光标是否越界并位于幻象单元格
}

// SetMethod 设置用于计算单元格宽度的方法
func (s *Screen) SetMethod(method ansi.Method) {
	s.method = method
}

// UseBackspaces 设置是否使用退格字符来移动光标
func (s *Screen) UseBackspaces(v bool) {
	s.opts.Backspace = v
}

// UseHardTabs 设置是否使用硬制表符来优化光标移动
func (s *Screen) UseHardTabs(v bool) {
	s.opts.HardTabs = v
}

// SetColorProfile 设置写入屏幕时使用的颜色配置文件
func (s *Screen) SetColorProfile(p colorprofile.Profile) {
	s.opts.Profile = p
}

// SetRelativeCursor 设置是否使用相对光标移动
func (s *Screen) SetRelativeCursor(v bool) {
	s.opts.RelativeCursor = v
}

// EnterAltScreen 进入替代屏幕缓冲区
func (s *Screen) EnterAltScreen() {
	s.opts.AltScreen = true
	s.clear = true
	s.saved = s.cur
}

// ExitAltScreen 退出替代屏幕缓冲区
func (s *Screen) ExitAltScreen() {
	s.opts.AltScreen = false
	s.clear = true
	s.cur = s.saved
}

// ShowCursor 显示光标
func (s *Screen) ShowCursor() {
	s.opts.ShowCursor = true
}

// HideCursor 隐藏光标
func (s *Screen) HideCursor() {
	s.opts.ShowCursor = false
}

// Bounds 实现 Window 接口
func (s *Screen) Bounds() Rectangle {
	// 始终返回新缓冲区的边界
	return s.newbuf.Bounds()
}

// Cell 实现 Window 接口
func (s *Screen) Cell(x int, y int) *Cell {
	return s.newbuf.Cell(x, y)
}

// Redraw 强制全屏重绘
func (s *Screen) Redraw() {
	s.mu.Lock()
	s.clear = true
	s.mu.Unlock()
}

// Clear 用空白单元格清除屏幕。这是 [Screen.Fill] 传入 nil 单元格的便捷方法
func (s *Screen) Clear() bool {
	return s.ClearRect(s.newbuf.Bounds())
}

// ClearRect 用空白单元格清除给定的矩形。这是 [Screen.FillRect] 传入 nil 单元格的便捷方法
func (s *Screen) ClearRect(r Rectangle) bool {
	return s.FillRect(nil, r)
}

// SetCell 实现 Window 接口
func (s *Screen) SetCell(x int, y int, cell *Cell) (v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cellWidth := 1
	if cell != nil {
		cellWidth = cell.Width
	}
	if prev := s.curbuf.Cell(x, y); !cellEqual(prev, cell) {
		chg, ok := s.touch[y]
		if !ok {
			chg = lineData{firstCell: x, lastCell: x + cellWidth}
		} else {
			chg.firstCell = min(chg.firstCell, x)
			chg.lastCell = max(chg.lastCell, x+cellWidth)
		}
		s.touch[y] = chg
	}

	return s.newbuf.SetCell(x, y, cell)
}

// Fill 实现 Window 接口
func (s *Screen) Fill(cell *Cell) bool {
	return s.FillRect(cell, s.newbuf.Bounds())
}

// FillRect 实现 Window 接口
func (s *Screen) FillRect(cell *Cell, r Rectangle) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.newbuf.FillRect(cell, r)
	for i := r.Min.Y; i < r.Max.Y; i++ {
		s.touch[i] = lineData{firstCell: r.Min.X, lastCell: r.Max.X}
	}
	return true
}

// capabilities 表示支持的 ANSI 转义序列的掩码
type capabilities uint

const (
	// 垂直位置绝对定位 [ansi.VPA]
	capVPA capabilities = 1 << iota
	// 水平位置绝对定位 [ansi.HPA]
	capHPA
	// 光标水平制表符 [ansi.CHT]
	capCHT
	// 光标向后制表符 [ansi.CBT]
	capCBT
	// 重复前一个字符 [ansi.REP]
	capREP
	// 擦除字符 [ansi.ECH]
	capECH
	// 插入字符 [ansi.ICH]
	capICH
	// 向下滚动 [ansi.SD]
	capSD
	// 向上滚动 [ansi.SU]
	capSU

	noCaps  capabilities = 0
	allCaps              = capVPA | capHPA | capCHT | capCBT | capREP | capECH | capICH |
		capSD | capSU
)

// Contains 返回 capabilities 是否包含给定的功能
func (v capabilities) Contains(c capabilities) bool {
	return v&c == c
}

// xtermCaps 返回终端是否类似 xterm。这意味着终端支持 ECMA-48 和 ANSI X3.64 转义序列
// xtermCaps 为给定的终端类型返回控制序列功能列表。这仅支持在不同终端之间可能不同的子集序列
// 注意：混合方法是支持 Terminfo 数据库以获得完整的功能集
func xtermCaps(termtype string) (v capabilities) {
	parts := strings.Split(termtype, "-")
	if len(parts) == 0 {
		return v
	}

	switch parts[0] {
	case
		"contour",
		"foot",
		"ghostty",
		"kitty",
		"rio",
		"st",
		"tmux",
		"wezterm",
		"xterm":
		v = allCaps
	case "alacritty":
		v = allCaps
		v &^= capCHT // 注意：alacritty 在 2024-12-28 #62d5b13 中添加了对 [ansi.CHT] 的支持
	case "screen":
		// 参见 https://www.gnu.org/software/screen/manual/screen.html#Control-Sequences-1
		v = allCaps
		v &^= capREP
	case "linux":
		// 参见 https://man7.org/linux/man-pages/man4/console_codes.4.html
		v = capVPA | capHPA | capECH | capICH
	}

	return v
}

// NewScreen creates a new Screen.
func NewScreen(w io.Writer, width, height int, opts *ScreenOptions) (s *Screen) {
	s = new(Screen)
	s.w = w
	if opts != nil {
		s.opts = *opts
	}

	if s.opts.Term == "" {
		s.opts.Term = os.Getenv("TERM")
	}

	if width <= 0 || height <= 0 {
		if f, ok := w.(term.File); ok {
			width, height, _ = term.GetSize(f.Fd())
		}
	}
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	s.buf = new(bytes.Buffer)
	s.caps = xtermCaps(s.opts.Term)
	s.curbuf = NewBuffer(width, height)
	s.newbuf = NewBuffer(width, height)
	s.cur = Cursor{Position: Pos(-1, -1)} // start at -1 to force a move
	s.saved = s.cur
	s.reset()

	return s
}

// Width returns the width of the screen.
func (s *Screen) Width() int {
	return s.newbuf.Width()
}

// Height returns the height of the screen.
func (s *Screen) Height() int {
	return s.newbuf.Height()
}

// cellEqual returns whether the two cells are equal. A nil cell is considered
// a [BlankCell].
func cellEqual(a, b *Cell) bool {
	if a == b {
		return true
	}
	if a == nil {
		a = &BlankCell
	}
	if b == nil {
		b = &BlankCell
	}
	return a.Equal(b)
}

// putCell draws a cell at the current cursor position.
func (s *Screen) putCell(cell *Cell) {
	width, height := s.newbuf.Width(), s.newbuf.Height()
	if s.opts.AltScreen && s.cur.X == width-1 && s.cur.Y == height-1 {
		s.putCellLR(cell)
	} else {
		s.putAttrCell(cell)
	}
}

// wrapCursor wraps the cursor to the next line.
//

func (s *Screen) wrapCursor() {
	const autoRightMargin = true
	if autoRightMargin {
		// Assume we have auto wrap mode enabled.
		s.cur.X = 0
		s.cur.Y++
	} else {
		s.cur.X--
	}
}

func (s *Screen) putAttrCell(cell *Cell) {
	if cell != nil && cell.Empty() {
		// XXX: Zero width cells are special and should not be written to the
		// screen no matter what other attributes they have.
		// Zero width cells are used for wide characters that are split into
		// multiple cells.
		return
	}

	if cell == nil {
		cell = s.clearBlank()
	}

	// We're at pending wrap state (phantom cell), incoming cell should
	// wrap.
	if s.atPhantom {
		s.wrapCursor()
		s.atPhantom = false
	}

	s.updatePen(cell)
	s.buf.WriteRune(cell.Rune)
	for _, c := range cell.Comb {
		s.buf.WriteRune(c)
	}

	s.cur.X += cell.Width

	if cell.Width > 0 {
		s.queuedText = true
	}

	if s.cur.X >= s.newbuf.Width() {
		s.atPhantom = true
	}
}

// putCellLR draws a cell at the lower right corner of the screen.
func (s *Screen) putCellLR(cell *Cell) {
	// Optimize for the lower right corner cell.
	curX := s.cur.X
	if cell == nil || !cell.Empty() {
		s.buf.WriteString(ansi.ResetModeAutoWrap)
		s.putAttrCell(cell)
		// Writing to lower-right corner cell should not wrap.
		s.atPhantom = false
		s.cur.X = curX
		s.buf.WriteString(ansi.SetModeAutoWrap)
	}
}

// updatePen updates the cursor pen styles.
func (s *Screen) updatePen(cell *Cell) {
	if cell == nil {
		cell = &BlankCell
	}

	if s.opts.Profile != 0 {
		// Downsample colors to the given color profile.
		cell.Style = ConvertStyle(cell.Style, s.opts.Profile)
		cell.Link = ConvertLink(cell.Link, s.opts.Profile)
	}

	if !cell.Style.Equal(&s.cur.Style) {
		seq := cell.Style.DiffSequence(s.cur.Style)
		if cell.Style.Empty() && len(seq) > len(ansi.ResetStyle) {
			seq = ansi.ResetStyle
		}
		s.buf.WriteString(seq)
		s.cur.Style = cell.Style
	}
	if !cell.Link.Equal(&s.cur.Link) {
		s.buf.WriteString(ansi.SetHyperlink(cell.Link.URL, cell.Link.Params))
		s.cur.Link = cell.Link
	}
}

// emitRange emits a range of cells to the buffer. It it equivalent to calling
// [Screen.putCell] for each cell in the range. This is optimized to use
// [ansi.ECH] and [ansi.REP].
// Returns whether the cursor is at the end of interval or somewhere in the
// middle.
func (s *Screen) emitRange(line Line, n int) (eoi bool) {
	for n > 0 {
		var count int
		for n > 1 && !cellEqual(line.At(0), line.At(1)) {
			s.putCell(line.At(0))
			line = line[1:]
			n--
		}

		cell0 := line[0]
		if n == 1 {
			s.putCell(cell0)
			return false
		}

		count = 2
		for count < n && cellEqual(line.At(count), cell0) {
			count++
		}

		ech := ansi.EraseCharacter(count)
		cup := ansi.CursorPosition(s.cur.X+count, s.cur.Y)
		rep := ansi.RepeatPreviousCharacter(count)
		if s.caps.Contains(capECH) && count > len(ech)+len(cup) && cell0 != nil && cell0.Clear() { //nolint:nestif
			s.updatePen(cell0)
			s.buf.WriteString(ech)

			// If this is the last cell, we don't need to move the cursor.
			if count < n {
				s.move(s.cur.X+count, s.cur.Y)
			} else {
				return true // cursor in the middle
			}
		} else if s.caps.Contains(capREP) && count > len(rep) &&
			(cell0 == nil || (len(cell0.Comb) == 0 && cell0.Rune < 256)) {
			// We only support ASCII characters. Most terminals will handle
			// non-ASCII characters correctly, but some might not, ahem xterm.
			//
			// NOTE: [ansi.REP] only repeats the last rune and won't work
			// if the last cell contains multiple runes.

			wrapPossible := s.cur.X+count >= s.newbuf.Width()
			repCount := count
			if wrapPossible {
				repCount--
			}

			s.updatePen(cell0)
			s.putCell(cell0)
			repCount-- // cell0 is a single width cell ASCII character

			s.buf.WriteString(ansi.RepeatPreviousCharacter(repCount))
			s.cur.X += repCount
			if wrapPossible {
				s.putCell(cell0)
			}
		} else {
			for i := range count {
				s.putCell(line.At(i))
			}
		}

		line = line[clamp(count, 0, len(line)):]
		n -= count
	}

	return eoi
}

// putRange puts a range of cells from the old line to the new line.
// Returns whether the cursor is at the end of interval or somewhere in the
// middle.
func (s *Screen) putRange(oldLine, newLine Line, y, start, end int) (eoi bool) {
	inline := min(len(ansi.CursorPosition(start+1, y+1)),
		min(len(ansi.HorizontalPositionAbsolute(start+1)),
			len(ansi.CursorForward(start+1))))
	if (end - start + 1) > inline { //nolint:nestif
		var j, same int
		for j, same = start, 0; j <= end; j++ {
			oldCell, newCell := oldLine.At(j), newLine.At(j)
			if same == 0 && oldCell != nil && oldCell.Empty() {
				continue
			}
			if cellEqual(oldCell, newCell) {
				same++
			} else {
				if same > end-start {
					s.emitRange(newLine[start:], j-same-start)
					s.move(j, y)
					start = j
				}
				same = 0
			}
		}

		i := s.emitRange(newLine[start:], j-same-start)

		// Always return 1 for the next [Screen.move] after a [Screen.putRange] if
		// we found identical characters at end of interval.
		if same == 0 {
			return i
		}
		return true
	}

	return s.emitRange(newLine[start:], end-start+1)
}

// clearToEnd clears the screen from the current cursor position to the end of
// line.
func (s *Screen) clearToEnd(blank *Cell, force bool) { //nolint:unparam
	if s.cur.Y >= 0 {
		curline := s.curbuf.Line(s.cur.Y)
		for j := s.cur.X; j < s.curbuf.Width(); j++ {
			if j >= 0 {
				c := curline.At(j)
				if !cellEqual(c, blank) {
					curline.Set(j, blank)
					force = true
				}
			}
		}
	}

	if force {
		s.updatePen(blank)
		count := s.newbuf.Width() - s.cur.X
		if s.el0Cost() <= count {
			s.buf.WriteString(ansi.EraseLineRight)
		} else {
			for range count {
				s.putCell(blank)
			}
		}
	}
}

// clearBlank returns a blank cell based on the current cursor background color.
func (s *Screen) clearBlank() *Cell {
	c := BlankCell
	if !s.cur.Style.Empty() || !s.cur.Link.Empty() {
		c.Style = s.cur.Style
		c.Link = s.cur.Link
	}
	return &c
}

// insertCells inserts the count cells pointed by the given line at the current
// cursor position.
func (s *Screen) insertCells(line Line, count int) {
	supportsICH := s.caps.Contains(capICH)
	if supportsICH {
		// Use [ansi.ICH] as an optimization.
		s.buf.WriteString(ansi.InsertCharacter(count))
	} else {
		// Otherwise, use [ansi.IRM] mode.
		s.buf.WriteString(ansi.SetModeInsertReplace)
	}

	for i := 0; count > 0; i++ {
		s.putAttrCell(line[i])
		count--
	}

	if !supportsICH {
		s.buf.WriteString(ansi.ResetModeInsertReplace)
	}
}

// el0Cost returns the cost of using [ansi.EL] 0 i.e. [ansi.EraseLineRight]. If
// this terminal supports background color erase, it can be cheaper to use
// [ansi.EL] 0 i.e. [ansi.EraseLineRight] to clear
// trailing spaces.
func (s *Screen) el0Cost() int {
	if s.caps != noCaps {
		return 0
	}
	return len(ansi.EraseLineRight)
}

// transformLine transforms the given line in the current window to the
// corresponding line in the new window. It uses [ansi.ICH] and [ansi.DCH] to
// insert or delete characters.
func (s *Screen) transformLine(y int) {
	var firstCell, oLastCell, nLastCell int // first, old last, new last index
	oldLine := s.curbuf.Line(y)
	newLine := s.newbuf.Line(y)

	// Find the first changed cell in the line
	var lineChanged bool
	for i := range s.newbuf.Width() {
		if !cellEqual(newLine.At(i), oldLine.At(i)) {
			lineChanged = true
			break
		}
	}

	const ceolStandoutGlitch = false
	if ceolStandoutGlitch && lineChanged { //nolint:nestif
		s.move(0, y)
		s.clearToEnd(nil, false)
		s.putRange(oldLine, newLine, y, 0, s.newbuf.Width()-1)
	} else {
		blank := newLine.At(0)

		// It might be cheaper to clear leading spaces with [ansi.EL] 1 i.e.
		// [ansi.EraseLineLeft].
		if blank == nil || blank.Clear() {
			var oFirstCell, nFirstCell int
			for oFirstCell = range s.curbuf.Width() {
				if !cellEqual(oldLine.At(oFirstCell), blank) {
					break
				}
			}
			for nFirstCell = range s.newbuf.Width() {
				if !cellEqual(newLine.At(nFirstCell), blank) {
					break
				}
			}

			if nFirstCell == oFirstCell {
				firstCell = nFirstCell

				// Find the first differing cell
				for firstCell < s.newbuf.Width() &&
					cellEqual(oldLine.At(firstCell), newLine.At(firstCell)) {
					firstCell++
				}
			} else if oFirstCell > nFirstCell {
				firstCell = nFirstCell
			} else if oFirstCell < nFirstCell {
				firstCell = oFirstCell
				el1Cost := len(ansi.EraseLineLeft)
				if el1Cost < nFirstCell-oFirstCell {
					if nFirstCell >= s.newbuf.Width() {
						s.move(0, y)
						s.updatePen(blank)
						s.buf.WriteString(ansi.EraseLineRight)
					} else {
						s.move(nFirstCell-1, y)
						s.updatePen(blank)
						s.buf.WriteString(ansi.EraseLineLeft)
					}

					for firstCell < nFirstCell {
						oldLine.Set(firstCell, blank)
						firstCell++
					}
				}
			}
		} else {
			// Find the first differing cell
			for firstCell < s.newbuf.Width() && cellEqual(newLine.At(firstCell), oldLine.At(firstCell)) {
				firstCell++
			}
		}

		// If we didn't find one, we're done
		if firstCell >= s.newbuf.Width() {
			return
		}

		blank = newLine.At(s.newbuf.Width() - 1)
		if blank != nil && !blank.Clear() {
			// Find the last differing cell
			nLastCell = s.newbuf.Width() - 1
			for nLastCell > firstCell && cellEqual(newLine.At(nLastCell), oldLine.At(nLastCell)) {
				nLastCell--
			}

			if nLastCell >= firstCell {
				s.move(firstCell, y)
				s.putRange(oldLine, newLine, y, firstCell, nLastCell)
				if firstCell < len(oldLine) && firstCell < len(newLine) {
					copy(oldLine[firstCell:], newLine[firstCell:])
				} else {
					copy(oldLine, newLine)
				}
			}

			return
		}

		// Find last non-blank cell in the old line.
		oLastCell = s.curbuf.Width() - 1
		for oLastCell > firstCell && cellEqual(oldLine.At(oLastCell), blank) {
			oLastCell--
		}

		// Find last non-blank cell in the new line.
		nLastCell = s.newbuf.Width() - 1
		for nLastCell > firstCell && cellEqual(newLine.At(nLastCell), blank) {
			nLastCell--
		}

		if nLastCell == firstCell && s.el0Cost() < oLastCell-nLastCell {
			s.move(firstCell, y)
			if !cellEqual(newLine.At(firstCell), blank) {
				s.putCell(newLine.At(firstCell))
			}
			s.clearToEnd(blank, false)
		} else if nLastCell != oLastCell &&
			!cellEqual(newLine.At(nLastCell), oldLine.At(oLastCell)) {
			s.move(firstCell, y)
			if oLastCell-nLastCell > s.el0Cost() {
				if s.putRange(oldLine, newLine, y, firstCell, nLastCell) {
					s.move(nLastCell+1, y)
				}
				s.clearToEnd(blank, false)
			} else {
				n := max(nLastCell, oLastCell)
				s.putRange(oldLine, newLine, y, firstCell, n)
			}
		} else {
			nLastNonBlank := nLastCell
			oLastNonBlank := oLastCell

			// Find the last cells that really differ.
			// Can be -1 if no cells differ.
			for cellEqual(newLine.At(nLastCell), oldLine.At(oLastCell)) {
				if !cellEqual(newLine.At(nLastCell-1), oldLine.At(oLastCell-1)) {
					break
				}
				nLastCell--
				oLastCell--
				if nLastCell == -1 || oLastCell == -1 {
					break
				}
			}

			n := min(oLastCell, nLastCell)
			if n >= firstCell {
				s.move(firstCell, y)
				s.putRange(oldLine, newLine, y, firstCell, n)
			}

			if oLastCell < nLastCell {
				m := max(nLastNonBlank, oLastNonBlank)
				if n != 0 {
					for n > 0 {
						wide := newLine.At(n + 1)
						if wide == nil || !wide.Empty() {
							break
						}
						n--
						oLastCell--
					}
				} else if n >= firstCell && newLine.At(n) != nil && newLine.At(n).Width > 1 {
					next := newLine.At(n + 1)
					for next != nil && next.Empty() {
						n++
						oLastCell++
					}
				}

				s.move(n+1, y)
				ichCost := 3 + nLastCell - oLastCell
				if s.caps.Contains(capICH) && (nLastCell < nLastNonBlank || ichCost > (m-n)) {
					s.putRange(oldLine, newLine, y, n+1, m)
				} else {
					s.insertCells(newLine[n+1:], nLastCell-oLastCell)
				}
			} else if oLastCell > nLastCell {
				s.move(n+1, y)
				dchCost := 3 + oLastCell - nLastCell
				if dchCost > len(ansi.EraseLineRight)+nLastNonBlank-(n+1) {
					if s.putRange(oldLine, newLine, y, n+1, nLastNonBlank) {
						s.move(nLastNonBlank+1, y)
					}
					s.clearToEnd(blank, false)
				} else {
					s.updatePen(blank)
					s.deleteCells(oLastCell - nLastCell)
				}
			}
		}
	}

	// Update the old line with the new line
	if firstCell < len(oldLine) && firstCell < len(newLine) {
		copy(oldLine[firstCell:], newLine[firstCell:])
	} else {
		copy(oldLine, newLine)
	}
}

// deleteCells deletes the count cells at the current cursor position and moves
// the rest of the line to the left. This is equivalent to [ansi.DCH].
func (s *Screen) deleteCells(count int) {
	// [ansi.DCH] will shift in cells from the right margin so we need to
	// ensure that they are the right style.
	s.buf.WriteString(ansi.DeleteCharacter(count))
}

// clearToBottom clears the screen from the current cursor position to the end
// of the screen.
func (s *Screen) clearToBottom(blank *Cell) {
	row, col := s.cur.Y, s.cur.X
	if row < 0 {
		row = 0
	}

	s.updatePen(blank)
	s.buf.WriteString(ansi.EraseScreenBelow)
	// Clear the rest of the current line
	s.curbuf.ClearRect(Rect(col, row, s.curbuf.Width()-col, 1))
	// Clear everything below the current line
	s.curbuf.ClearRect(Rect(0, row+1, s.curbuf.Width(), s.curbuf.Height()-row-1))
}

// clearBottom tests if clearing the end of the screen would satisfy part of
// the screen update. Scan backwards through lines in the screen checking if
// each is blank and one or more are changed.
// It returns the top line.
func (s *Screen) clearBottom(total int) (top int) {
	if total <= 0 {
		return top
	}

	top = total
	last := s.newbuf.Width()
	blank := s.clearBlank()
	canClearWithBlank := blank == nil || blank.Clear()

	if canClearWithBlank { //nolint:nestif
		var row int
		for row = total - 1; row >= 0; row-- {
			oldLine := s.curbuf.Line(row)
			newLine := s.newbuf.Line(row)

			var col int
			ok := true
			for col = 0; ok && col < last; col++ {
				ok = cellEqual(newLine.At(col), blank)
			}
			if !ok {
				break
			}

			for col = 0; ok && col < last; col++ {
				ok = len(oldLine) == last && cellEqual(oldLine.At(col), blank)
			}
			if !ok {
				top = row
			}
		}

		if top < total {
			s.move(0, top-1) // top is 1-based
			s.clearToBottom(blank)
			if s.oldhash != nil && s.newhash != nil &&
				row < len(s.oldhash) && row < len(s.newhash) {
				for row := top; row < s.newbuf.Height(); row++ {
					s.oldhash[row] = s.newhash[row]
				}
			}
		}
	}

	return top
}

// clearScreen clears the screen and put cursor at home.
func (s *Screen) clearScreen(blank *Cell) {
	s.updatePen(blank)
	s.buf.WriteString(ansi.CursorHomePosition)
	s.buf.WriteString(ansi.EraseEntireScreen)
	s.cur.X, s.cur.Y = 0, 0
	s.curbuf.Fill(blank)
}

// clearBelow clears everything below and including the row.
func (s *Screen) clearBelow(blank *Cell, row int) {
	s.move(0, row)
	s.clearToBottom(blank)
}

// clearUpdate forces a screen redraw.
func (s *Screen) clearUpdate() {
	blank := s.clearBlank()
	var nonEmpty int
	if s.opts.AltScreen {
		// XXX: We're using the maximum height of the two buffers to ensure
		// we write newly added lines to the screen in [Screen.transformLine].
		nonEmpty = max(s.curbuf.Height(), s.newbuf.Height())
		s.clearScreen(blank)
	} else {
		nonEmpty = s.newbuf.Height()
		s.clearBelow(blank, 0)
	}
	nonEmpty = s.clearBottom(nonEmpty)
	for i := range nonEmpty {
		s.transformLine(i)
	}
}

// Flush flushes the buffer to the screen.
func (s *Screen) Flush() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flush()
}

func (s *Screen) flush() (err error) {
	// Write the buffer
	if s.buf.Len() > 0 {
		_, err = s.w.Write(s.buf.Bytes())
		if err == nil {
			s.buf.Reset()
		}
	}

	return err //nolint:wrapcheck
}

// Render renders changes of the screen to the internal buffer. Call
// [Screen.Flush] to flush pending changes to the screen.
func (s *Screen) Render() {
	s.mu.Lock()
	s.render()
	s.mu.Unlock()
}

func (s *Screen) render() {
	// Do we need to render anything?
	if s.opts.AltScreen == s.altScreenMode &&
		!s.opts.ShowCursor == s.cursorHidden &&
		!s.clear &&
		len(s.touch) == 0 &&
		len(s.queueAbove) == 0 {
		return
	}

	//nolint:godox
	// TODO: Investigate whether this is necessary. Theoretically, terminals
	// can add/remove tab stops and we should be able to handle that. We could
	// use [ansi.DECTABSR] to read the tab stops, but that's not implemented in
	// most terminals :/
	// // Are we using hard tabs? If so, ensure tabs are using the
	// // default interval using [ansi.DECST8C].
	// if s.opts.HardTabs && !s.initTabs {
	// 	s.buf.WriteString(ansi.SetTabEvery8Columns)
	// 	s.initTabs = true
	// }

	// Do we need alt-screen mode?
	if s.opts.AltScreen != s.altScreenMode {
		if s.opts.AltScreen {
			s.buf.WriteString(ansi.SetModeAltScreenSaveCursor)
		} else {
			s.buf.WriteString(ansi.ResetModeAltScreenSaveCursor)
		}
		s.altScreenMode = s.opts.AltScreen
	}

	// Do we need text cursor mode?
	if !s.opts.ShowCursor != s.cursorHidden {
		s.cursorHidden = !s.opts.ShowCursor
		if s.cursorHidden {
			s.buf.WriteString(ansi.HideCursor)
		}
	}

	// Do we have queued strings to write above the screen?
	if len(s.queueAbove) > 0 {
		//nolint:godox
		// TODO: Use scrolling region if available.
		//nolint:godox
		// TODO: Use [Screen.Write] [io.Writer] interface.

		// We need to scroll the screen up by the number of lines in the queue.
		// We can't use [ansi.SU] because we want the cursor to move down until
		// it reaches the bottom of the screen.
		s.move(0, s.newbuf.Height()-1)
		s.buf.WriteString(strings.Repeat("\n", len(s.queueAbove)))
		s.cur.Y += len(s.queueAbove)
		// XXX: Now go to the top of the screen, insert new lines, and write
		// the queued strings. It is important to use [Screen.moveCursor]
		// instead of [Screen.move] because we don't want to perform any checks
		// on the cursor position.
		s.moveCursor(0, 0, false)
		s.buf.WriteString(ansi.InsertLine(len(s.queueAbove)))
		for _, line := range s.queueAbove {
			s.buf.WriteString(line + "\r\n")
		}

		// Clear the queue
		s.queueAbove = s.queueAbove[:0]
	}

	var nonEmpty int

	// XXX: In inline mode, after a screen resize, we need to clear the extra
	// lines at the bottom of the screen. This is because in inline mode, we
	// don't use the full screen height and the current buffer size might be
	// larger than the new buffer size.
	partialClear := !s.opts.AltScreen && s.cur.X != -1 && s.cur.Y != -1 &&
		s.curbuf.Width() == s.newbuf.Width() &&
		s.curbuf.Height() > 0 &&
		s.curbuf.Height() > s.newbuf.Height()

	if !s.clear && partialClear {
		s.clearBelow(nil, s.newbuf.Height()-1)
	}

	if s.clear { //nolint:nestif
		s.clearUpdate()
		s.clear = false
	} else if len(s.touch) > 0 {
		if s.opts.AltScreen {
			// Optimize scrolling for the alternate screen buffer.
			//nolint:godox
			// TODO: Should we optimize for inline mode as well? If so, we need
			// to know the actual cursor position to use [ansi.DECSTBM].
			s.scrollOptimize()
		}

		var changedLines int
		var i int

		if s.opts.AltScreen {
			nonEmpty = min(s.curbuf.Height(), s.newbuf.Height())
		} else {
			nonEmpty = s.newbuf.Height()
		}

		nonEmpty = s.clearBottom(nonEmpty)
		for i = range nonEmpty {
			_, ok := s.touch[i]
			if ok {
				s.transformLine(i)
				changedLines++
			}
		}
	}

	// Sync windows and screen
	s.touch = make(map[int]lineData, s.newbuf.Height())

	if s.curbuf.Width() != s.newbuf.Width() || s.curbuf.Height() != s.newbuf.Height() {
		// Resize the old buffer to match the new buffer.
		_, oldh := s.curbuf.Width(), s.curbuf.Height()
		s.curbuf.Resize(s.newbuf.Width(), s.newbuf.Height())
		// Sync new lines to old lines
		for i := oldh - 1; i < s.newbuf.Height(); i++ {
			copy(s.curbuf.Line(i), s.newbuf.Line(i))
		}
	}

	s.updatePen(nil) // nil indicates a blank cell with no styles

	// Do we have enough changes to justify toggling the cursor?
	if s.buf.Len() > 1 && s.opts.ShowCursor && !s.cursorHidden && s.queuedText {
		nb := new(bytes.Buffer)
		nb.Grow(s.buf.Len() + len(ansi.HideCursor) + len(ansi.ShowCursor))
		nb.WriteString(ansi.HideCursor)
		nb.Write(s.buf.Bytes())
		nb.WriteString(ansi.ShowCursor)
		*s.buf = *nb
	}

	s.queuedText = false
}

// Close writes the final screen update and resets the screen.
func (s *Screen) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.render()
	s.updatePen(nil)
	// Go to the bottom of the screen
	s.move(0, s.newbuf.Height()-1)

	if s.altScreenMode {
		s.buf.WriteString(ansi.ResetModeAltScreenSaveCursor)
		s.altScreenMode = false
	}

	if s.cursorHidden {
		s.buf.WriteString(ansi.ShowCursor)
		s.cursorHidden = false
	}

	// Write the buffer
	err = s.flush()
	if err != nil {
		return err
	}

	s.reset()
	return err
}

// reset resets the screen to its initial state.
func (s *Screen) reset() {
	s.scrollHeight = 0
	s.cursorHidden = false
	s.altScreenMode = false
	s.touch = make(map[int]lineData, s.newbuf.Height())
	if s.curbuf != nil {
		s.curbuf.Clear()
	}
	if s.newbuf != nil {
		s.newbuf.Clear()
	}
	s.buf.Reset()
	s.tabs = DefaultTabStops(s.newbuf.Width())
	s.oldhash, s.newhash = nil, nil

	// We always disable HardTabs when termtype is "linux".
	if strings.HasPrefix(s.opts.Term, "linux") {
		s.opts.HardTabs = false
	}
}

// Resize resizes the screen.
func (s *Screen) Resize(width, height int) bool {
	oldw := s.newbuf.Width()
	oldh := s.newbuf.Height()

	if s.opts.AltScreen || width != oldw {
		// We only clear the whole screen if the width changes. Adding/removing
		// rows is handled by the [Screen.render] and [Screen.transformLine]
		// methods.
		s.clear = true
	}

	// Clear new columns and lines
	if width > oldh {
		s.ClearRect(Rect(max(oldw-1, 0), 0, width-oldw, height))
	} else if width < oldw {
		s.ClearRect(Rect(max(width-1, 0), 0, oldw-width, height))
	}

	if height > oldh {
		s.ClearRect(Rect(0, max(oldh, 0), width, height-oldh))
	} else if height < oldh {
		s.ClearRect(Rect(0, max(height, 0), width, oldh-height))
	}

	s.mu.Lock()
	s.newbuf.Resize(width, height)
	s.tabs.Resize(width)
	s.oldhash, s.newhash = nil, nil
	s.scrollHeight = 0 // reset scroll lines
	s.mu.Unlock()

	return true
}

// MoveTo moves the cursor to the given position.
func (s *Screen) MoveTo(x, y int) {
	s.mu.Lock()
	s.move(x, y)
	s.mu.Unlock()
}

// InsertAbove inserts string above the screen. The inserted string is not
// managed by the screen. This does nothing when alternate screen mode is
// enabled.
func (s *Screen) InsertAbove(str string) {
	if s.opts.AltScreen {
		return
	}
	s.mu.Lock()
	for _, line := range strings.Split(str, "\n") {
		s.queueAbove = append(s.queueAbove, s.method.Truncate(line, s.Width(), ""))
	}
	s.mu.Unlock()
}
