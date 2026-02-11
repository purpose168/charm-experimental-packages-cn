package cellbuf

// DefaultTabInterval 是默认的制表符间隔。
const DefaultTabInterval = 8

// TabStops 表示水平行的制表位。
type TabStops struct {
	stops    []int // 存储制表位的数组
	interval int   // 制表符间隔
	width    int   // 宽度
}

// NewTabStops 从列数和间隔创建一组新的制表位。
func NewTabStops(width, interval int) *TabStops {
	ts := new(TabStops)
	ts.interval = interval
	ts.width = width
	// 计算所需的制表位数组大小
	ts.stops = make([]int, (width+(interval-1))/interval)
	// 初始化制表位
	ts.init(0, width)
	return ts
}

// DefaultTabStops 使用默认间隔创建一组新的制表位。
func DefaultTabStops(cols int) *TabStops {
	return NewTabStops(cols, DefaultTabInterval)
}

// Resize 将制表位调整到给定的宽度。
func (ts *TabStops) Resize(width int) {
	if width == ts.width {
		return
	}

	if width < ts.width {
		// 缩小数组大小
		size := (width + (ts.interval - 1)) / ts.interval
		ts.stops = ts.stops[:size]
	} else {
		// 扩展数组大小
		size := (width - ts.width + (ts.interval - 1)) / ts.interval
		ts.stops = append(ts.stops, make([]int, size)...)
	}

	// 初始化新的制表位
	ts.init(ts.width, width)
	ts.width = width
}

// IsStop 如果给定列是制表位，则返回 true。
func (ts TabStops) IsStop(col int) bool {
	mask := ts.mask(col)
	i := col >> 3 // 相当于 col / 8
	if i < 0 || i >= len(ts.stops) {
		return false
	}
	return ts.stops[i]&mask != 0
}

// Next 返回给定列之后的下一个制表位。
func (ts TabStops) Next(col int) int {
	return ts.Find(col, 1)
}

// Prev 返回给定列之前的上一个制表位。
func (ts TabStops) Prev(col int) int {
	return ts.Find(col, -1)
}

// Find 根据给定的列和增量返回前一个/下一个制表位。
// 如果增量为正，则返回给定列之后的下一个制表位。
// 如果增量为负，则返回给定列之前的上一个制表位。
// 如果增量为零，则返回给定的列。
func (ts TabStops) Find(col, delta int) int {
	if delta == 0 {
		return col
	}

	var prev bool
	count := delta
	if count < 0 {
		count = -count
		prev = true
	}

	for count > 0 {
		if !prev {
			// 向前查找
			if col >= ts.width-1 {
				return col
			}

			col++
		} else {
			// 向后查找
			if col < 1 {
				return col
			}

			col--
		}

		// 检查是否是制表位
		if ts.IsStop(col) {
			count--
		}
	}

	return col
}

// Set 在给定列添加一个制表位。
func (ts *TabStops) Set(col int) {
	mask := ts.mask(col)
	ts.stops[col>>3] |= mask
}

// Reset 移除给定列的制表位。
func (ts *TabStops) Reset(col int) {
	mask := ts.mask(col)
	ts.stops[col>>3] &= ^mask
}

// Clear 移除所有制表位。
func (ts *TabStops) Clear() {
	ts.stops = make([]int, len(ts.stops))
}

// mask 返回给定列的掩码。
func (ts *TabStops) mask(col int) int {
	return 1 << (col & (ts.interval - 1))
}

// init 从 col 开始初始化制表位，直到 width。
func (ts *TabStops) init(col, width int) {
	for x := col; x < width; x++ {
		if x%ts.interval == 0 {
			// 在间隔位置设置制表位
			ts.Set(x)
		} else {
			// 移除非间隔位置的制表位
			ts.Reset(x)
		}
	}
}
