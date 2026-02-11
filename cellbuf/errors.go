package cellbuf

import "errors"

// ErrOutOfBounds 当给定的 x, y 位置超出边界时返回。
var ErrOutOfBounds = errors.New("超出边界")
