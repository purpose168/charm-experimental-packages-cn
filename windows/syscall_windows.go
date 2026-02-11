package windows

import "golang.org/x/sys/windows"

// NewLazySystemDLL 是 windows.NewLazySystemDLL 的类型别名。
var NewLazySystemDLL = windows.NewLazySystemDLL

// Handle 是 windows.Handle 的类型别名。
type Handle = windows.Handle

//sys	ReadConsoleInput(console Handle, buf *InputRecord, toread uint32, read *uint32) (err error) = kernel32.ReadConsoleInputW
//sys	PeekConsoleInput(console Handle, buf *InputRecord, toread uint32, read *uint32) (err error) = kernel32.PeekConsoleInputW
//sys	GetNumberOfConsoleInputEvents(console Handle, numevents *uint32) (err error) = kernel32.GetNumberOfConsoleInputEvents
//sys	FlushConsoleInputBuffer(console Handle) (err error) = kernel32.FlushConsoleInputBuffer
