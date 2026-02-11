package vt

import (
	"image/color"
	"io"

	uv "github.com/charmbracelet/ultraviolet"
)

// Terminal 表示虚拟终端接口。
type Terminal interface {
	BackgroundColor() color.Color           // 获取背景颜色
	Blur()                                  // 失去焦点
	Bounds() uv.Rectangle                   // 获取边界矩形
	CellAt(x int, y int) *uv.Cell           // 获取指定位置的单元格
	Close() error                           // 关闭终端
	CursorColor() color.Color               // 获取光标颜色
	CursorPosition() uv.Position            // 获取光标位置
	Draw(scr uv.Screen, area uv.Rectangle)  // 绘制到屏幕
	Focus()                                 // 获取焦点
	ForegroundColor() color.Color           // 获取前景颜色
	Height() int                            // 获取高度
	IndexedColor(i int) color.Color         // 获取索引颜色
	InputPipe() io.Writer                   // 获取输入管道
	Paste(text string)                      // 粘贴文本
	Read(p []byte) (n int, err error)       // 读取数据
	RegisterApcHandler(handler ApcHandler)  // 注册APC处理器
	RegisterCsiHandler(cmd int, handler CsiHandler)  // 注册CSI处理器
	RegisterDcsHandler(cmd int, handler DcsHandler)  // 注册DCS处理器
	RegisterEscHandler(cmd int, handler EscHandler)  // 注册ESC处理器
	RegisterOscHandler(cmd int, handler OscHandler)  // 注册OSC处理器
	RegisterPmHandler(handler PmHandler)    // 注册PM处理器
	RegisterSosHandler(handler SosHandler)  // 注册SOS处理器
	Render() string                         // 渲染为字符串
	Resize(width int, height int)           // 调整大小
	SendKey(k uv.KeyEvent)                  // 发送按键事件
	SendKeys(keys ...uv.KeyEvent)           // 发送多个按键事件
	SendMouse(m Mouse)                      // 发送鼠标事件
	SendText(text string)                   // 发送文本
	SetBackgroundColor(c color.Color)       // 设置背景颜色
	SetCallbacks(cb Callbacks)              // 设置回调
	SetCell(x int, y int, c *uv.Cell)       // 设置单元格
	SetCursorColor(c color.Color)           // 设置光标颜色
	SetDefaultBackgroundColor(c color.Color) // 设置默认背景颜色
	SetDefaultCursorColor(c color.Color)    // 设置默认光标颜色
	SetDefaultForegroundColor(c color.Color) // 设置默认前景颜色
	SetForegroundColor(c color.Color)       // 设置前景颜色
	SetIndexedColor(i int, c color.Color)   // 设置索引颜色
	SetLogger(l Logger)                     // 设置日志器
	String() string                         // 转换为字符串
	Touched() []*uv.LineData                // 获取被修改的行数据
	Width() int                             // 获取宽度
	WidthMethod() uv.WidthMethod            // 获取宽度计算方法
	Write(p []byte) (n int, err error)       // 写入数据
	WriteString(s string) (n int, err error) // 写入字符串
}
