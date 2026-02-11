// errors 包提供错误处理工具
package errors

import "strings"

// Join 返回一个包装给定错误的错误。
// 任何 nil 错误值都会被丢弃。
// 如果 errs 中的每个值都是 nil，Join 返回 nil。
// 错误格式为调用 errs 中每个元素的 Error 方法所获得的字符串的连接，
// 每个字符串之间有一个换行符。
//
// Join 返回的非 nil 错误实现了 Unwrap() []error 方法。
//
// 这是从 Go 1.20 errors.Unwrap 复制的，做了一些调整以避免使用 unsafe。
// 主要目标是在较旧的 Go 版本中提供此功能。
func Join(errs ...error) error {
	var nonNil []error //nolint:prealloc
	for _, err := range errs {
		if err == nil {
			continue
		}
		nonNil = append(nonNil, err)
	}
	if len(nonNil) == 0 {
		return nil
	}
	return &joinError{
		errs: nonNil,
	}
}

// joinError 表示多个错误的连接
//
// 字段：
//   - errs: 错误列表

type joinError struct {
	errs []error
}

// Error 返回错误的字符串表示
//
// 返回值：
//   - string: 错误信息，由所有错误的 Error() 方法返回值连接而成，每个错误信息之间用换行符分隔
func (e *joinError) Error() string {
	strs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		strs = append(strs, err.Error())
	}
	return strings.Join(strs, "\n")
}

// Unwrap 返回包装的错误列表
//
// 返回值：
//   - []error: 错误列表
func (e *joinError) Unwrap() []error {
	return e.errs
}
