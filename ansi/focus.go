package ansi

// Focus 是一个转义序列，用于通知终端它已获得焦点。
// 这与 [FocusEventMode] 一起使用。
const Focus = "\x1b[I"

// Blur 是一个转义序列，用于通知终端它已失去焦点。
// 这与 [FocusEventMode] 一起使用。
const Blur = "\x1b[O"
