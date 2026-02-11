package ansi

import "io"

// Execute 是一个函数，通过将给定的转义序列写入到提供的输出写入器来"执行"该序列。
//
// 这是 [io.WriteString] 的语法糖。
func Execute(w io.Writer, s string) (int, error) {
	return io.WriteString(w, s) //nolint:wrapcheck
}
