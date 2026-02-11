package parser

import "math"

// 用于序列参数和中间字节的移位和掩码。
const (
	PrefixShift    = 8
	IntermedShift  = 16
	FinalMask      = 0xff
	HasMoreFlag    = math.MinInt32
	ParamMask      = ^HasMoreFlag
	MissingParam   = ParamMask
	MissingCommand = MissingParam
	MaxParam       = math.MaxUint16 // 参数可以拥有的最大值
)

const (
	// MaxParamsSize 是序列可以拥有的最大参数数量。
	MaxParamsSize = 32

	// DefaultParamValue 是用于缺失参数的默认值。
	DefaultParamValue = 0
)

// Prefix 返回序列的前缀字节。
// 这始终是以下之一 '<' '=' '>' '?' 并且在 0x3C-0x3F 范围内。
// 如果序列没有前缀，则返回零。
func Prefix(cmd int) int {
	return (cmd >> PrefixShift) & FinalMask
}

// Intermediate 返回序列的中间字节。
// 中间字节在 0x20-0x2F 范围内。这包括这些字符：' ', '!', '"', '#', '$', '%', '&', ''', '(', ')', '*', '+',
// ',', '-', '.', '/'。
// 如果序列没有中间字节，则返回零。
func Intermediate(cmd int) int {
	return (cmd >> IntermedShift) & FinalMask
}

// Command 返回 CSI 序列的命令字节。
func Command(cmd int) int {
	return cmd & FinalMask
}

// Param 返回给定索引处的参数。
// 如果参数不存在，则返回 -1。
func Param(params []int, i int) int {
	if len(params) == 0 || i < 0 || i >= len(params) {
		return -1
	}

	p := params[i] & ParamMask
	if p == MissingParam {
		return -1
	}

	return p
}

// HasMore 如果参数有更多子参数，则返回 true。
func HasMore(params []int, i int) bool {
	if len(params) == 0 || i >= len(params) {
		return false
	}

	return params[i]&HasMoreFlag != 0
}

// Subparams 返回给定参数的子参数。
// 如果参数不存在，则返回 nil。
func Subparams(params []int, i int) []int {
	if len(params) == 0 || i < 0 || i >= len(params) {
		return nil
	}

	// 计算给定参数索引之前的参数数量。
	var count int
	var j int
	for j = range params {
		if count == i {
			break
		}
		if !HasMore(params, j) {
			count++
		}
	}

	if count > i || j >= len(params) {
		return nil
	}

	var subs []int
	for ; j < len(params); j++ {
		if !HasMore(params, j) {
			break
		}
		p := Param(params, j)
		if p == -1 {
			p = DefaultParamValue
		}
		subs = append(subs, p)
	}

	p := Param(params, j)
	if p == -1 {
		p = DefaultParamValue
	}

	return append(subs, p)
}

// Len 返回序列中的参数数量。
// 这将返回序列中的参数数量，不包括任何子参数。
func Len(params []int) int {
	var n int
	for i := range params {
		if !HasMore(params, i) {
			n++
		}
	}
	return n
}

// Range 遍历序列的参数，并为每个参数调用给定的函数。
// 函数应返回 false 以停止迭代。
func Range(params []int, fn func(i int, param int, hasMore bool) bool) {
	for i := range params {
		if !fn(i, Param(params, i), HasMore(params, i)) {
			break
		}
	}
}
