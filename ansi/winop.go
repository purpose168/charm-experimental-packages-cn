package ansi

import (
	"strconv"
	"strings"
)

const (
	// ResizeWindowWinOp 是一个调整终端窗口大小的窗口操作
	//
	// 已弃用：请直接使用常量数字与[WindowOp]一起使用
	ResizeWindowWinOp = 4

	// RequestWindowSizeWinOp 是一个请求报告终端窗口像素大小的窗口操作
	// 响应格式为：
	//  CSI 4 ; height ; width t
	//
	// 已弃用：请直接使用常量数字与[WindowOp]一起使用
	RequestWindowSizeWinOp = 14

	// RequestCellSizeWinOp 是一个请求报告终端单元格像素大小的窗口操作
	// 响应格式为：
	//  CSI 6 ; height ; width t
	//
	// 已弃用：请直接使用常量数字与[WindowOp]一起使用
	RequestCellSizeWinOp = 16
)

// WindowOp (XTWINOPS) 是一个操作终端窗口的序列
//
// 格式：
//
//	CSI Ps ; Ps ; Ps t
//
// Ps 是一个用分号分隔的参数列表
// 请参阅：https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h4-Functions-using-CSI-_-ordered-by-the-final-character-lparen-s-rparen:CSI-Ps;Ps;Ps-t.1EB0
func WindowOp(p int, ps ...int) string {
	if p <= 0 {
		return "" // 如果参数p小于等于0，返回空字符串
	}

	if len(ps) == 0 {
		return "\x1b[" + strconv.Itoa(p) + "t" // 只有一个参数的情况
	}

	// 创建参数列表
	params := make([]string, 0, len(ps)+1)
	params = append(params, strconv.Itoa(p)) // 添加第一个参数
	for _, p := range ps {
		if p >= 0 {
			params = append(params, strconv.Itoa(p)) // 添加后续参数
		}
	}

	// 拼接成完整的ANSI序列
	return "\x1b[" + strings.Join(params, ";") + "t"
}

// XTWINOPS 是[WindowOp]的别名
func XTWINOPS(p int, ps ...int) string {
	return WindowOp(p, ps...) // 直接调用WindowOp函数
}
