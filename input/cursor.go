package input

import "image"

// CursorPositionEvent 表示光标位置事件。其中 X 是从零开始的列，Y 是从零开始的行。
type CursorPositionEvent image.Point
