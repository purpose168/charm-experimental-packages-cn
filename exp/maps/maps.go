// Package maps 提供用于处理映射的实用函数。
package maps

import (
	"cmp"
	"slices"
)

// SortedKeys 返回映射 m 的键。
// 键将被排序。
func SortedKeys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	r := Keys(m)
	slices.Sort(r)
	return r
}

// Keys 返回映射 m 的键。
func Keys[M ~map[K]V, K cmp.Ordered, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
