// Package higherorder 为 Go 提供高阶函数。
package higherorder

// Foldl 从左侧开始，对列表的每个元素应用一个函数。
// 返回单个值。
func Foldl[A any](f func(x, y A) A, start A, list []A) A {
	for _, v := range list {
		start = f(start, v)
	}
	return start
}

// Foldr 从右侧开始，对列表的每个元素应用一个函数。
// 返回单个值。
func Foldr[A any](f func(x, y A) A, start A, list []A) A {
	for i := len(list) - 1; i >= 0; i-- {
		start = f(start, list[i])
	}
	return start
}

// Map 对列表的每个元素应用给定的函数，返回一个新列表。
func Map[A, B any](f func(A) B, list []A) []B {
	res := make([]B, len(list))
	for i, v := range list {
		res[i] = f(v)
	}
	return res
}

// Filter 对列表的每个元素应用一个函数，如果函数返回 false，则移除这些元素，返回一个新列表。
func Filter[A any](f func(A) bool, list []A) []A {
	res := make([]A, 0)
	for _, v := range list {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
