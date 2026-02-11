// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sixel

import (
	"math/bits"
)

const (
	unknownHint sortedHint = iota
	increasingHint
	decreasingHint
)

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Compare 返回
//
//	-1 如果 x 小于 y，
//	 0 如果 x 等于 y，
//	+1 如果 x 大于 y。
//
// 对于浮点类型，NaN 被认为小于任何非 NaN 值，
// NaN 被认为等于 NaN，-0.0 等于 0.0。
func compare[T ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN {
		if yNaN {
			return 0
		}
		return -1
	}
	if yNaN {
		return +1
	}
	if x < y {
		return -1
	}
	if x > y {
		return +1
	}
	return 0
}

// isNaN 报告 x 是否为 NaN，不需要 math 包。
// 如果 T 不是浮点类型，这将始终返回 false。
func isNaN[T ordered](x T) bool {
	return x != x
}

func sortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	n := len(x)
	pdqsortCmpFunc(x, 0, n, bits.Len(uint(n)), cmp)
}

type sortedHint int // pdqsort 选择枢轴时的提示

// xorshift 论文: https://www.jstatsoft.org/article/view/v008i14/xorshift.pdf
type xorshift uint64

func (r *xorshift) Next() uint64 {
	*r ^= *r << 13
	*r ^= *r >> 7
	*r ^= *r << 17
	return uint64(*r)
}

func nextPowerOfTwo(length int) uint {
	return 1 << bits.Len(uint(length)) //nolint:gosec
}

// insertionSortCmpFunc 使用插入排序对 data[a:b] 进行排序。
func insertionSortCmpFunc[E any](data []E, a, b int, cmp func(a, b E) int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (cmp(data[j], data[j-1]) < 0); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

// siftDownCmpFunc 在 data[lo:hi] 上实现堆属性。
// first 是数组中堆根所在的偏移量。
func siftDownCmpFunc[E any](data []E, lo, hi, first int, cmp func(a, b E) int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && (cmp(data[first+child], data[first+child+1]) < 0) {
			child++
		}
		if !(cmp(data[first+root], data[first+child]) < 0) { //nolint:staticcheck
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func heapSortCmpFunc[E any](data []E, a, b int, cmp func(a, b E) int) {
	first := a
	lo := 0
	hi := b - a

	// 构建堆，最大元素在顶部。
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownCmpFunc(data, i, hi, first, cmp)
	}

	// 弹出元素，最大的先弹出，到 data 的末尾。
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDownCmpFunc(data, lo, i, first, cmp)
	}
}

// pdqsortCmpFunc 对 data[a:b] 进行排序。
// 该算法基于模式击败快速排序（pdqsort），但没有 BlockQuicksort 的优化。
// pdqsort 论文: https://arxiv.org/pdf/2106.05123.pdf
// C++ 实现: https://github.com/orlp/pdqsort
// Rust 实现: https://docs.rs/pdqsort/latest/pdqsort/
// limit 是在回退到堆排序之前允许的坏（非常不平衡）枢轴的数量。
func pdqsortCmpFunc[E any](data []E, a, b, limit int, cmp func(a, b E) int) {
	const maxInsertion = 12

	var (
		wasBalanced    = true // 上一次分区是否合理平衡
		wasPartitioned = true // 切片是否已经分区
	)

	for {
		length := b - a

		if length <= maxInsertion {
			insertionSortCmpFunc(data, a, b, cmp)
			return
		}

		// 如果做出了太多错误选择，则回退到堆排序。
		if limit == 0 {
			heapSortCmpFunc(data, a, b, cmp)
			return
		}

		// 如果上一次分区不平衡，我们需要打破模式。
		if !wasBalanced {
			breakPatternsCmpFunc(data, a, b)
			limit--
		}

		pivot, hint := choosePivotCmpFunc(data, a, b, cmp)
		if hint == decreasingHint {
			reverseRangeCmpFunc(data, a, b)
			// 选择的枢轴是数组开始后 pivot-a 个元素。
			// 反转后，它是数组结束前 pivot-a 个元素。
			// 这个想法来自 Rust 的实现。
			pivot = (b - 1) - (pivot - a)
			hint = increasingHint
		}

		// 切片可能已经排序。
		if wasBalanced && wasPartitioned && hint == increasingHint {
			if partialInsertionSortCmpFunc(data, a, b, cmp) {
				return
			}
		}

		// 可能切片包含许多重复元素，将切片分区为
		// 等于枢轴的元素和大于枢轴的元素。
		if a > 0 && !(cmp(data[a-1], data[pivot]) < 0) { //nolint:staticcheck
			mid := partitionEqualCmpFunc(data, a, b, pivot, cmp)
			a = mid
			continue
		}

		mid, alreadyPartitioned := partitionCmpFunc(data, a, b, pivot, cmp)
		wasPartitioned = alreadyPartitioned

		leftLen, rightLen := mid-a, b-mid
		balanceThreshold := length / 8
		if leftLen < rightLen {
			wasBalanced = leftLen >= balanceThreshold
			pdqsortCmpFunc(data, a, mid, limit, cmp)
			a = mid + 1
		} else {
			wasBalanced = rightLen >= balanceThreshold
			pdqsortCmpFunc(data, mid+1, b, limit, cmp)
			b = mid
		}
	}
}

// partitionCmpFunc 执行一次快速排序分区。
// 设 p = data[pivot]
// 移动 data[a:b] 中的元素，使得 data[i]<p 且 data[j]>=p 对于 i<newpivot 和 j>newpivot。
// 返回时，data[newpivot] = p。
func partitionCmpFunc[E any](data []E, a, b, pivot int, cmp func(a, b E) int) (newpivot int, alreadyPartitioned bool) {
	data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i 和 j 包含待分区的元素

	for i <= j && (cmp(data[i], data[a]) < 0) {
		i++
	}
	for i <= j && !(cmp(data[j], data[a]) < 0) { //nolint:staticcheck
		j--
	}
	if i > j {
		data[j], data[a] = data[a], data[j]
		return j, true
	}
	data[i], data[j] = data[j], data[i]
	i++
	j--

	for {
		for i <= j && (cmp(data[i], data[a]) < 0) {
			i++
		}
		for i <= j && !(cmp(data[j], data[a]) < 0) { //nolint:staticcheck
			j--
		}
		if i > j {
			break
		}
		data[i], data[j] = data[j], data[i]
		i++
		j--
	}
	data[j], data[a] = data[a], data[j]
	return j, false
}

// partitionEqualCmpFunc 将 data[a:b] 分区为等于 data[pivot] 的元素，后跟大于 data[pivot] 的元素。
// 假设 data[a:b] 不包含小于 data[pivot] 的元素。
func partitionEqualCmpFunc[E any](data []E, a, b, pivot int, cmp func(a, b E) int) (newpivot int) {
	data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i 和 j 包含待分区的元素

	for {
		for i <= j && !(cmp(data[a], data[i]) < 0) { //nolint:staticcheck
			i++
		}
		for i <= j && (cmp(data[a], data[j]) < 0) {
			j--
		}
		if i > j {
			break
		}
		data[i], data[j] = data[j], data[i]
		i++
		j--
	}
	return i
}

// partialInsertionSortCmpFunc 部分排序切片，如果切片在末尾已排序则返回 true。
func partialInsertionSortCmpFunc[E any](data []E, a, b int, cmp func(a, b E) int) bool {
	const (
		maxSteps         = 5  // 将被移动的最大相邻无序对数
		shortestShifting = 50 // 不移动短数组上的任何元素
	)
	i := a + 1
	for range maxSteps {
		for i < b && !(cmp(data[i], data[i-1]) < 0) { //nolint:staticcheck
			i++
		}

		if i == b {
			return true
		}

		if b-a < shortestShifting {
			return false
		}

		data[i], data[i-1] = data[i-1], data[i]

		// 将较小的向左移动。
		if i-a >= 2 {
			for j := i - 1; j >= 1; j-- {
				if !(cmp(data[j], data[j-1]) < 0) { //nolint:staticcheck
					break
				}
				data[j], data[j-1] = data[j-1], data[j]
			}
		}
		// 将较大的向右移动。
		if b-i >= 2 {
			for j := i + 1; j < b; j++ {
				if !(cmp(data[j], data[j-1]) < 0) { //nolint:staticcheck
					break
				}
				data[j], data[j-1] = data[j-1], data[j]
			}
		}
	}
	return false
}

// breakPatternsCmpFunc 分散一些元素以试图打破一些模式
// 这些模式可能导致快速排序中的不平衡分区。
func breakPatternsCmpFunc[E any](data []E, a, b int) {
	length := b - a
	if length >= 8 {
		random := xorshift(length)
		modulus := nextPowerOfTwo(length)

		for idx := a + (length/4)*2 - 1; idx <= a+(length/4)*2+1; idx++ {
			other := int(uint(random.Next()) & (modulus - 1)) //nolint:gosec
			if other >= length {
				other -= length
			}
			data[idx], data[a+other] = data[a+other], data[idx]
		}
	}
}

// choosePivotCmpFunc 在 data[a:b] 中选择枢轴。
//
// [0,8): 选择静态枢轴。
// [8,shortestNinther): 使用简单的三位数中值法。
// [shortestNinther,∞): 使用 Tukey 九分位数法。
func choosePivotCmpFunc[E any](data []E, a, b int, cmp func(a, b E) int) (pivot int, hint sortedHint) {
	const (
		shortestNinther = 50
		maxSwaps        = 4 * 3
	)

	l := b - a

	var (
		swaps int
		i     = a + l/4*1
		j     = a + l/4*2
		k     = a + l/4*3
	)

	if l >= 8 {
		if l >= shortestNinther {
			// Tukey 九分位数法，这个想法来自 Rust 的实现。
			i = medianAdjacentCmpFunc(data, i, &swaps, cmp)
			j = medianAdjacentCmpFunc(data, j, &swaps, cmp)
			k = medianAdjacentCmpFunc(data, k, &swaps, cmp)
		}
		// 找到 i、j、k 的中值并将其存储到 j 中。
		j = medianCmpFunc(data, i, j, k, &swaps, cmp)
	}

	switch swaps {
	case 0:
		return j, increasingHint
	case maxSwaps:
		return j, decreasingHint
	default:
		return j, unknownHint
	}
}

// order2CmpFunc 返回 x,y，其中 data[x] <= data[y]，x,y=a,b 或 x,y=b,a。
func order2CmpFunc[E any](data []E, a, b int, swaps *int, cmp func(a, b E) int) (int, int) {
	if cmp(data[b], data[a]) < 0 {
		*swaps++
		return b, a
	}
	return a, b
}

// medianCmpFunc 返回 x，其中 data[x] 是 data[a],data[b],data[c] 的中值，x 是 a、b 或 c。
func medianCmpFunc[E any](data []E, a, b, c int, swaps *int, cmp func(a, b E) int) int {
	a, b = order2CmpFunc(data, a, b, swaps, cmp)
	b, _ = order2CmpFunc(data, b, c, swaps, cmp)
	_, b = order2CmpFunc(data, a, b, swaps, cmp)
	return b
}

// medianAdjacentCmpFunc 找到 data[a - 1], data[a], data[a + 1] 的中值并将索引存储到 a 中。
func medianAdjacentCmpFunc[E any](data []E, a int, swaps *int, cmp func(a, b E) int) int {
	return medianCmpFunc(data, a-1, a, a+1, swaps, cmp)
}

func reverseRangeCmpFunc[E any](data []E, a, b int) {
	i := a
	j := b - 1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i++
		j--
	}
}
