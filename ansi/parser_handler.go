package ansi

import "unsafe"

// Params 表示一组打包的参数。
type Params []Param

// Param 返回指定索引处的参数，以及它是否属于子参数的一部分。
// 如果参数缺失，则返回默认值。如果索引超出范围，则返回默认值和 false。
func (p Params) Param(i, def int) (int, bool, bool) {
	if i < 0 || i >= len(p) {
		return def, false, false
	}
	return p[i].Param(def), p[i].HasMore(), true
}

// ForEach 遍历参数并为每个参数调用给定的函数。
// 如果参数是子参数的一部分，则调用时 hasMore 会设置为 true。
// 使用 def 为缺失的参数设置默认值。
func (p Params) ForEach(def int, f func(i, param int, hasMore bool)) {
	for i := range p {
		f(i, p[i].Param(def), p[i].HasMore())
	}
}

// ToParams 将整数列表转换为参数列表。
func ToParams(params []int) Params {
	return unsafe.Slice((*Param)(unsafe.Pointer(&params[0])), len(params))
}

// Handler 处理解析器执行的操作。
// 用于处理 ANSI 转义序列、控制字符和符文。
type Handler struct {
	// Print 在遇到可打印符文时被调用。
	Print func(r rune)
	// Execute 在遇到控制字符时被调用。
	Execute func(b byte)
	// HandleCsi 在遇到 CSI 序列时被调用。
	HandleCsi func(cmd Cmd, params Params)
	// HandleEsc 在遇到 ESC 序列时被调用。
	HandleEsc func(cmd Cmd)
	// HandleDcs 在遇到 DCS 序列时被调用。
	HandleDcs func(cmd Cmd, params Params, data []byte)
	// HandleOsc 在遇到 OSC 序列时被调用。
	HandleOsc func(cmd int, data []byte)
	// HandlePm 在遇到 PM 序列时被调用。
	HandlePm func(data []byte)
	// HandleApc 在遇到 APC 序列时被调用。
	HandleApc func(data []byte)
	// HandleSos 在遇到 SOS 序列时被调用。
	HandleSos func(data []byte)
}

// SetHandler 为解析器设置处理器。
func (p *Parser) SetHandler(h Handler) {
	p.handler = h
}
