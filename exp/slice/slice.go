// Package slice 提供在 Go 中处理切片的实用函数。
package slice

import (
	"iter"
	"slices"
)

// GroupBy 根据键函数对切片中的项目进行分组。
func GroupBy[T any, K comparable](list []T, key func(T) K) map[K][]T {
	groups := make(map[K][]T)

	for _, item := range list {
		k := key(item)
		groups[k] = append(groups[k], item)
	}

	return groups
}

// Take 返回给定切片的前 n 个元素。如果切片中的元素不足，返回整个切片。
func Take[A any](slice []A, n int) []A {
	if n > len(slice) {
		return slice
	}
	return slice[:n]
}

// Last 返回切片的最后一个元素和 true。如果切片为空，返回零值和 false。
func Last[T any](list []T) (T, bool) {
	if len(list) == 0 {
		var zero T
		return zero, false
	}
	return list[len(list)-1], true
}

// Uniq 返回一个去除了所有重复项的新切片。
func Uniq[T comparable](list []T) []T {
	seen := make(map[T]struct{}, len(list))
	uniqList := make([]T, 0, len(list))

	for _, item := range list {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			uniqList = append(uniqList, item)
		}
	}

	return uniqList
}

// Intersperse 在切片的每个元素之间插入一个项目，返回一个新切片。
func Intersperse[T any](slice []T, insert T) []T {
	if len(slice) <= 1 {
		return slice
	}

	// 创建一个具有所需容量的新切片。
	result := make([]T, len(slice)*2-1)

	for i := range slice {
		// 用原始元素和插入项填充新切片。
		result[i*2] = slice[i]

		// 在项目之间添加插入项（最后一个除外）。
		if i < len(slice)-1 {
			result[i*2+1] = insert
		}
	}

	return result
}

// ContainsAny 检查列表中是否存在任何给定的值。
func ContainsAny[T comparable](list []T, values ...T) bool {
	return slices.ContainsFunc(list, func(v T) bool {
		return slices.Contains(values, v)
	})
}

// Shift 移除并返回切片的第一个元素。
// 返回被移除的元素和修改后的切片。
// 第三个返回值 (ok) 表示是否移除了元素。
func Shift[T any](slice []T) (element T, newSlice []T, ok bool) {
	if len(slice) == 0 {
		var zero T
		return zero, slice, false
	}
	return slice[0], slice[1:], true
}

// Pop 移除并返回切片的最后一个元素。
// 返回被移除的元素和修改后的切片。
// 第三个返回值 (ok) 表示是否移除了元素。
func Pop[T any](slice []T) (element T, newSlice []T, ok bool) {
	if len(slice) == 0 {
		var zero T
		return zero, slice, false
	}
	lastIdx := len(slice) - 1
	return slice[lastIdx], slice[:lastIdx], true
}

// DeleteAt 移除并返回指定索引处的元素。
// 返回被移除的元素和修改后的切片。
// 第三个返回值 (ok) 表示是否移除了元素。
func DeleteAt[T any](slice []T, index int) (element T, newSlice []T, ok bool) {
	if index < 0 || index >= len(slice) {
		var zero T
		return zero, slice, false
	}

	element = slice[index]
	newSlice = slices.Delete(slices.Clone(slice), index, index+1)

	return element, newSlice, true
}

// IsSubset 检查切片 a 的所有元素是否都存在于切片 b 中。
func IsSubset[T comparable](a, b []T) bool {
	if len(a) > len(b) {
		return false
	}
	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}
	for _, item := range a {
		if _, exists := set[item]; !exists {
			return false
		}
	}
	return true
}

// Map 接受一个类型为 E 的迭代器和一个映射函数，返回一个类型为 F 的迭代器。
func Map[E any, F any](seq iter.Seq[E], fn func(e E) F) iter.Seq[F] {
	return func(yield func(F) bool) {
		for e := range seq {
			yield(fn(e))
		}
	}
}
