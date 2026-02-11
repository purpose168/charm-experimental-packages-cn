// Package ordered 为有序类型提供实用函数。
package ordered

import "cmp"

// Clamp 返回一个夹在给定的最小值和最大值之间的值。
func Clamp[T cmp.Ordered](n, low, high T) T {
	if low > high {
		low, high = high, low
	}
	return min(high, max(low, n))
}

// First 返回固定数量的 [cmp.Ordered] 类型参数中的第一个非默认值。
func First[T cmp.Ordered](x T, y ...T) T {
	var empty T
	if x != empty {
		return x
	}
	for _, s := range y {
		if s != empty {
			return s
		}
	}
	return empty
}
