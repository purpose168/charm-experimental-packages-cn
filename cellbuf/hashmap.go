package cellbuf

import (
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
)

// hash 返回 [Line] 的哈希值。
func hash(l Line) (h uint64) {
	for _, c := range l {
		var r rune
		if c == nil {
			r = ansi.SP
		} else {
			r = c.Rune
		}
		h += (h << 5) + uint64(r)
	}
	return h
}

// hashmap 表示单个 [Line] 的哈希。
type hashmap struct {
	value              uint64
	oldcount, newcount int
	oldindex, newindex int
}

// 用于指示通过插入和滚动创建的行的值。
const newIndex = -1

// updateHashmap 用新的哈希值更新哈希表。
func (s *Screen) updateHashmap() {
	height := s.newbuf.Height()
	if len(s.oldhash) >= height && len(s.newhash) >= height {
		// 重新哈希已更改的行
		for i := range height {
			_, ok := s.touch[i]
			if ok {
				s.oldhash[i] = hash(s.curbuf.Line(i))
				s.newhash[i] = hash(s.newbuf.Line(i))
			}
		}
	} else {
		// 重新哈希所有行
		if len(s.oldhash) != height {
			s.oldhash = make([]uint64, height)
		}
		if len(s.newhash) != height {
			s.newhash = make([]uint64, height)
		}
		for i := range height {
			s.oldhash[i] = hash(s.curbuf.Line(i))
			s.newhash[i] = hash(s.newbuf.Line(i))
		}
	}

	s.hashtab = make([]hashmap, height*2)
	for i := range height {
		hashval := s.oldhash[i]

		// 查找匹配的哈希或空槽
		idx := 0
		for idx < len(s.hashtab) && s.hashtab[idx].value != 0 {
			if s.hashtab[idx].value == hashval {
				break
			}
			idx++
		}

		s.hashtab[idx].value = hashval // 以防这是新哈希
		s.hashtab[idx].oldcount++
		s.hashtab[idx].oldindex = i
	}
	for i := range height {
		hashval := s.newhash[i]

		// 查找匹配的哈希或空槽
		idx := 0
		for idx < len(s.hashtab) && s.hashtab[idx].value != 0 {
			if s.hashtab[idx].value == hashval {
				break
			}
			idx++
		}

		s.hashtab[idx].value = hashval // 以防这是新哈希
		s.hashtab[idx].newcount++
		s.hashtab[idx].newindex = i

		s.oldnum[i] = newIndex // 初始化旧索引切片
	}

	// 标记对应于唯一哈希对的行对。
	for i := 0; i < len(s.hashtab) && s.hashtab[i].value != 0; i++ {
		hsp := &s.hashtab[i]
		if hsp.oldcount == 1 && hsp.newcount == 1 && hsp.oldindex != hsp.newindex {
			s.oldnum[hsp.newindex] = hsp.oldindex
		}
	}

	s.growHunks()

	// 消除不良或不可能的偏移。这包括移除那些因冲突而无法增长的块，
	// 以及那些要移动太远的块，它们可能会破坏多于携带的内容。
	for i := 0; i < height; {
		var start, shift, size int
		for i < height && s.oldnum[i] == newIndex {
			i++
		}
		if i >= height {
			break
		}
		start = i
		shift = s.oldnum[i] - i
		i++
		for i < height && s.oldnum[i] != newIndex && s.oldnum[i]-i == shift {
			i++
		}
		size = i - start
		if size < 3 || size+min(size/8, 2) < abs(shift) {
			for start < i {
				s.oldnum[start] = newIndex
				start++
			}
		}
	}

	// 清除无效块后，尝试增长其余块。
	s.growHunks()
}

// scrollOldhash 滚动旧哈希。
func (s *Screen) scrollOldhash(n, top, bot int) {
	if len(s.oldhash) == 0 {
		return
	}

	size := bot - top + 1 - abs(n)
	if n > 0 {
		// 将现有哈希向上移动
		copy(s.oldhash[top:], s.oldhash[top+n:top+n+size])
		// 重新计算新移入行的哈希
		for i := bot; i > bot-n; i-- {
			s.oldhash[i] = hash(s.curbuf.Line(i))
		}
	} else {
		// 将现有哈希向下移动
		copy(s.oldhash[top-n:], s.oldhash[top:top+size])
		// 重新计算新移入行的哈希
		for i := top; i < top-n; i++ {
			s.oldhash[i] = hash(s.curbuf.Line(i))
		}
	}
}

func (s *Screen) growHunks() {
	var (
		backLimit    int // 要填充的单元格的限制
		backRefLimit int // 引用的限制
		i            int
		nextHunk     int
	)

	height := s.newbuf.Height()
	for i < height && s.oldnum[i] == newIndex {
		i++
	}
	for ; i < height; i = nextHunk {
		var (
			forwardLimit    int
			forwardRefLimit int
			end             int
			start           = i
			shift           = s.oldnum[i] - i
		)

		// 获取前向限制
		i = start + 1
		for i < height &&
			s.oldnum[i] != newIndex &&
			s.oldnum[i]-i == shift {
			i++
		}

		end = i
		for i < height && s.oldnum[i] == newIndex {
			i++
		}

		nextHunk = i
		forwardLimit = i
		if i >= height || s.oldnum[i] >= i {
			forwardRefLimit = i
		} else {
			forwardRefLimit = s.oldnum[i]
		}

		i = start - 1

		// 向后增长
		if shift < 0 {
			backLimit = backRefLimit + (-shift)
		}
		for i >= backLimit {
			if s.newhash[i] == s.oldhash[i+shift] ||
				s.costEffective(i+shift, i, shift < 0) {
				s.oldnum[i] = i + shift
			} else {
				break
			}
			i--
		}

		i = end
		// grow forward
		if shift > 0 {
			forwardLimit = forwardRefLimit - shift
		}
		for i < forwardLimit {
			if s.newhash[i] == s.oldhash[i+shift] ||
				s.costEffective(i+shift, i, shift > 0) {
				s.oldnum[i] = i + shift
			} else {
				break
			}
			i++
		}

		backLimit = i
		backRefLimit = backLimit
		if shift > 0 {
			backRefLimit += shift
		}
	}
}

// costEffective 如果将行 'from' 移动到行 'to' 的成本似乎是有效的，则返回 true。
// 'blank' 指示行 'to' 是否会变为空白。
func (s *Screen) costEffective(from, to int, blank bool) bool {
	if from == to {
		return false
	}

	newFrom := s.oldnum[from]
	if newFrom == newIndex {
		newFrom = from
	}

	// >= 左侧是移动前的成本。右侧是移动后的成本。

	// 计算移动前的成本。
	var costBeforeMove int
	if blank {
		// 更新目标处空白行的成本。
		costBeforeMove = s.updateCostBlank(s.newbuf.Line(to))
	} else {
		// 更新目标处现有行的成本。
		costBeforeMove = s.updateCost(s.curbuf.Line(to), s.newbuf.Line(to))
	}

	// 添加更新源行的成本
	costBeforeMove += s.updateCost(s.curbuf.Line(newFrom), s.newbuf.Line(from))

	// 计算移动后的成本。
	var costAfterMove int
	if newFrom == from {
		// 移动后源变为空白
		costAfterMove = s.updateCostBlank(s.newbuf.Line(from))
	} else {
		// 源从另一行获取更新
		costAfterMove = s.updateCost(s.curbuf.Line(newFrom), s.newbuf.Line(from))
	}

	// 添加将源行移动到目标的成本
	costAfterMove += s.updateCost(s.curbuf.Line(from), s.newbuf.Line(to))

	// 如果移动成本有效（成本更低或相等），则返回 true
	return costBeforeMove >= costAfterMove
}

func (s *Screen) updateCost(from, to Line) (cost int) {
	var fidx, tidx int
	for i := s.newbuf.Width() - 1; i > 0; i, fidx, tidx = i-1, fidx+1, tidx+1 {
		if !cellEqual(from.At(fidx), to.At(tidx)) {
			cost++
		}
	}
	return cost
}

func (s *Screen) updateCostBlank(to Line) (cost int) {
	var tidx int
	for i := s.newbuf.Width() - 1; i > 0; i, tidx = i-1, tidx+1 {
		if !cellEqual(nil, to.At(tidx)) {
			cost++
		}
	}
	return cost
}
