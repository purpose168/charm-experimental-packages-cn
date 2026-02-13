# API 接口文档

## 1. 接口定义概述

`charm-experimental-packages-cn` 项目是一个 Go 语言的库集合，提供了一系列终端相关的工具和功能。本文档详细描述了项目中各个模块的主要 API 接口，包括函数签名、参数说明、返回值说明等。

接口定义按照模块分类，每个模块的接口按照功能和用途进行组织。文档中的接口定义遵循 Go 语言的标准文档格式，便于开发者理解和使用。

## 2. 基础层模块接口

### 2.1 ansi 模块

#### 2.1.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Sequence` | 表示 ANSI 转义序列的结构 | `Prefix byte`, `Intermediates []byte`, `Parameters []int`, `Command byte` |
| `Command` | 表示 ANSI 命令的结构 | `Name string`, `Sequence Sequence`, `Handler func(*Command, ...interface{}) error` |
| `Parser` | ANSI 转义序列解析器 | `input []byte`, `pos int`, `sequences []Sequence` |

#### 2.1.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Parse` | `func(input []byte) []Sequence` | 解析 ANSI 转义序列 | input: 输入字节数组 | []Sequence: 解析出的转义序列 |
| `NewParser` | `func(input []byte) *Parser` | 创建新的 ANSI 解析器 | input: 输入字节数组 | *Parser: 解析器实例 |
| `(p *Parser) Next` | `func() (Sequence, error)` | 获取下一个转义序列 | 无 | Sequence: 转义序列, error: 错误信息 |
| `(p *Parser) ParseAll` | `func() []Sequence` | 解析所有转义序列 | 无 | []Sequence: 解析出的转义序列 |

### 2.2 term 模块

#### 2.2.1 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Size` | `func() (width, height int, err error)` | 获取终端大小 | 无 | width: 宽度, height: 高度, err: 错误信息 |
| `SetRaw` | `func() (oldState *Termios, err error)` | 设置终端为原始模式 | 无 | oldState: 原始终端状态, err: 错误信息 |
| `Restore` | `func(oldState *Termios)` | 恢复终端模式 | oldState: 原始终端状态 | 无 |
| `Fd` | `func() uintptr` | 获取终端文件描述符 | 无 | uintptr: 文件描述符 |

### 2.3 colors 模块

#### 2.3.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Color` | 表示颜色的结构 | `R, G, B uint8`, `A uint8` |
| `Profile` | 颜色配置文件 | `Name string`, `Colors map[string]Color` |

#### 2.3.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Parse` | `func(s string) (Color, error)` | 解析颜色字符串 | s: 颜色字符串 | Color: 解析出的颜色, error: 错误信息 |
| `Hex` | `func(hex string) (Color, error)` | 从十六进制字符串解析颜色 | hex: 十六进制颜色字符串 | Color: 解析出的颜色, error: 错误信息 |
| `RGB` | `func(r, g, b uint8) Color` | 创建 RGB 颜色 | r, g, b: 红绿蓝分量 | Color: 创建的颜色 |
| `RGBA` | `func(r, g, b, a uint8) Color` | 创建 RGBA 颜色 | r, g, b, a: 红绿蓝透明度分量 | Color: 创建的颜色 |

### 2.4 wcwidth 模块

#### 2.4.1 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `RuneWidth` | `func(r rune) int` | 计算单个字符的宽度 | r: 字符 | int: 字符宽度 |
| `StringWidth` | `func(s string) int` | 计算字符串的宽度 | s: 字符串 | int: 字符串宽度 |
| `Truncate` | `func(s string, width int) string` | 截断字符串到指定宽度 | s: 字符串, width: 目标宽度 | string: 截断后的字符串 |

## 3. 核心层模块接口

### 3.1 cellbuf 模块

#### 3.1.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Cell` | 表示终端单元格的结构 | `Rune rune`, `Comb []rune`, `Width int`, `Style Style` |
| `Line` | 表示终端行的结构 | `[]*Cell` |
| `Buffer` | 表示终端屏幕的缓冲区结构 | `Lines []Line` |
| `Rectangle` | 表示矩形区域的结构 | `Min, Max Point` |
| `Point` | 表示坐标点的结构 | `X, Y int` |

#### 3.1.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `NewCell` | `func(r rune, comb ...rune) *Cell` | 创建新的单元格 | r: 主字符, comb: 组合字符 | *Cell: 创建的单元格 |
| `NewCellString` | `func(s string) *Cell` | 创建带有字符串内容的单元格 | s: 字符串内容 | *Cell: 创建的单元格 |
| `NewGraphemeCell` | `func(s string) *Cell` | 创建带有字形集群的单元格 | s: 字符串内容 | *Cell: 创建的单元格 |
| `NewBuffer` | `func(width int, height int) *Buffer` | 创建新的缓冲区 | width: 宽度, height: 高度 | *Buffer: 创建的缓冲区 |
| `(b *Buffer) Cell` | `func(x int, y int) *Cell` | 获取指定位置的单元格 | x, y: 坐标 | *Cell: 单元格 |
| `(b *Buffer) SetCell` | `func(x, y int, c *Cell) bool` | 设置指定位置的单元格 | x, y: 坐标, c: 单元格 | bool: 是否设置成功 |
| `(b *Buffer) Width` | `func() int` | 获取缓冲区宽度 | 无 | int: 宽度 |
| `(b *Buffer) Height` | `func() int` | 获取缓冲区高度 | 无 | int: 高度 |
| `(b *Buffer) Resize` | `func(width int, height int)` | 调整缓冲区大小 | width: 宽度, height: 高度 | 无 |
| `(b *Buffer) Clear` | `func()` | 清除整个缓冲区 | 无 | 无 |
| `(b *Buffer) ClearRect` | `func(rect Rectangle)` | 清除指定矩形区域 | rect: 矩形区域 | 无 |
| `(b *Buffer) Fill` | `func(c *Cell)` | 用指定单元格填充整个缓冲区 | c: 单元格 | 无 |
| `(b *Buffer) FillRect` | `func(c *Cell, rect Rectangle)` | 用指定单元格填充指定矩形区域 | c: 单元格, rect: 矩形区域 | 无 |
| `(b *Buffer) InsertLine` | `func(y, n int, c *Cell)` | 在指定位置插入行 | y: 行位置, n: 行数, c: 单元格 | 无 |
| `(b *Buffer) DeleteLine` | `func(y, n int, c *Cell)` | 在指定位置删除行 | y: 行位置, n: 行数, c: 单元格 | 无 |
| `(b *Buffer) InsertCell` | `func(x, y, n int, c *Cell)` | 在指定位置插入单元格 | x, y: 坐标, n: 单元格数, c: 单元格 | 无 |
| `(b *Buffer) DeleteCell` | `func(x, y, n int, c *Cell)` | 在指定位置删除单元格 | x, y: 坐标, n: 单元格数, c: 单元格 | 无 |

### 3.2 input 模块

#### 3.2.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Input` | 输入处理器接口 | 无 |
| `Event` | 表示输入事件的结构 | `Type EventType`, `Key Key`, `Mouse MouseEvent`, `Rune rune`, `String string` |
| `EventType` | 事件类型枚举 | `None`, `Key`, `Mouse`, `Resize`, `Error` |
| `Key` | 表示按键的结构 | `Name string`, `Modifiers Modifier` |
| `MouseEvent` | 表示鼠标事件的结构 | `Type MouseEventType`, `X, Y int`, `Button MouseButton`, `Modifiers Modifier` |

#### 3.2.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `New` | `func() *Input` | 创建新的输入处理器 | 无 | *Input: 输入处理器 |
| `(i *Input) Start` | `func()` | 开始处理输入 | 无 | 无 |
| `(i *Input) Stop` | `func()` | 停止处理输入 | 无 | 无 |
| `(i *Input) Events` | `func() <-chan Event` | 获取事件通道 | 无 | <-chan Event: 事件通道 |
| `(i *Input) SetCursorMode` | `func(mode CursorMode)` | 设置光标模式 | mode: 光标模式 | 无 |
| `(i *Input) EnableMouse` | `func()` | 启用鼠标输入 | 无 | 无 |
| `(i *Input) DisableMouse` | `func()` | 禁用鼠标输入 | 无 | 无 |

### 3.3 xpty 模块

#### 3.3.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `PTY` | 表示伪终端的接口 | 无 |

#### 3.3.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `New` | `func() (PTY, error)` | 创建新的 PTY | 无 | PTY: 伪终端, error: 错误信息 |
| `(p PTY) Master` | `func() io.ReadWriteCloser` | 获取主设备 | 无 | io.ReadWriteCloser: 主设备 |
| `(p PTY) Slave` | `func() io.ReadWriteCloser` | 获取从设备 | 无 | io.ReadWriteCloser: 从设备 |
| `(p PTY) Close` | `func() error` | 关闭 PTY | 无 | error: 错误信息 |

### 3.4 conpty 模块

#### 3.4.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `ConPTY` | 表示 Windows ConPTY 的结构 | `hPC HWND`, `hIn HANDLE`, `hOut HANDLE`, `pid uint32` |

#### 3.4.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `New` | `func(cols, rows int) (*ConPTY, error)` | 创建新的 ConPTY | cols, rows: 列数和行数 | *ConPTY: ConPTY 实例, error: 错误信息 |
| `(c *ConPTY) Read` | `func(b []byte) (int, error)` | 从 ConPTY 读取数据 | b: 缓冲区 | int: 读取的字节数, error: 错误信息 |
| `(c *ConPTY) Write` | `func(b []byte) (int, error)` | 向 ConPTY 写入数据 | b: 缓冲区 | int: 写入的字节数, error: 错误信息 |
| `(c *ConPTY) Close` | `func() error` | 关闭 ConPTY | 无 | error: 错误信息 |

## 4. 应用层模块接口

### 4.1 pony 模块

#### 4.1.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Node` | 表示 AST 节点的接口 | 无 |
| `TextNode` | 表示文本节点的结构 | `Value string` |
| `BoldNode` | 表示粗体节点的结构 | `Children []Node` |
| `ItalicNode` | 表示斜体节点的结构 | `Children []Node` |
| `UnderlineNode` | 表示下划线节点的结构 | `Children []Node` |
| `StrikethroughNode` | 表示删除线节点的结构 | `Children []Node` |
| `ColorNode` | 表示颜色节点的结构 | `Color string`, `Children []Node` |
| `LinkNode` | 表示链接节点的结构 | `URL string`, `Text string` |

#### 4.1.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Parse` | `func(input string) []Node` | 解析 pony 标记 | input: 输入字符串 | []Node: AST 节点 |
| `Render` | `func(nodes []Node) string` | 渲染 pony 标记为终端输出 | nodes: AST 节点 | string: 渲染后的字符串 |
| `RenderNode` | `func(node Node) string` | 渲染单个节点 | node: AST 节点 | string: 渲染后的字符串 |

### 4.2 vt 模块

#### 4.2.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Terminal` | 表示虚拟终端的结构 | `buffer *cellbuf.Buffer`, `parser *ansi.Parser`, `state *State` |
| `State` | 表示终端状态的结构 | `cursor Point`, `mode Mode`, `scrollRegion Rectangle` |
| `Mode` | 表示终端模式的结构 | `Insert bool`, `AutoWrap bool`, `CursorVisible bool` |

#### 4.2.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `New` | `func(width, height int) *Terminal` | 创建新的虚拟终端 | width, height: 宽度和高度 | *Terminal: 虚拟终端实例 |
| `(t *Terminal) Write` | `func(b []byte) (int, error)` | 向终端写入数据 | b: 缓冲区 | int: 写入的字节数, error: 错误信息 |
| `(t *Terminal) Read` | `func(b []byte) (int, error)` | 从终端读取数据 | b: 缓冲区 | int: 读取的字节数, error: 错误信息 |
| `(t *Terminal) Resize` | `func(width, height int)` | 调整终端大小 | width, height: 宽度和高度 | 无 |
| `(t *Terminal) Buffer` | `func() *cellbuf.Buffer` | 获取终端缓冲区 | 无 | *cellbuf.Buffer: 缓冲区 |
| `(t *Terminal) Cursor` | `func() Point` | 获取光标位置 | 无 | Point: 光标位置 |
| `(t *Terminal) SetCursor` | `func(x, y int)` | 设置光标位置 | x, y: 坐标 | 无 |

### 4.3 editor 模块

#### 4.3.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Editor` | 表示编辑器的结构 | `buffer *Buffer`, `cursor *Cursor`, `input *input.Input` |
| `Buffer` | 表示编辑缓冲区的结构 | `lines []string` |
| `Cursor` | 表示光标的结构 | `x, y int` |

#### 4.3.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `New` | `func() *Editor` | 创建新的编辑器 | 无 | *Editor: 编辑器实例 |
| `(e *Editor) Insert` | `func(r rune)` | 插入字符 | r: 字符 | 无 |
| `(e *Editor) InsertString` | `func(s string)` | 插入字符串 | s: 字符串 | 无 |
| `(e *Editor) Delete` | `func()` | 删除字符 | 无 | 无 |
| `(e *Editor) Backspace` | `func()` | 退格删除 | 无 | 无 |
| `(e *Editor) CursorUp` | `func()` | 光标上移 | 无 | 无 |
| `(e *Editor) CursorDown` | `func()` | 光标下移 | 无 | 无 |
| `(e *Editor) CursorLeft` | `func()` | 光标左移 | 无 | 无 |
| `(e *Editor) CursorRight` | `func()` | 光标右移 | 无 | 无 |
| `(e *Editor) Content` | `func() string` | 获取编辑器内容 | 无 | string: 编辑器内容 |
| `(e *Editor) SetContent` | `func(s string)` | 设置编辑器内容 | s: 字符串 | 无 |

## 5. 工具层模块接口

### 5.1 etag 模块

#### 5.1.1 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Generate` | `func(data []byte) string` | 生成 ETag | data: 数据 | string: 生成的 ETag |
| `GenerateWeak` | `func(data []byte) string` | 生成弱 ETag | data: 数据 | string: 生成的弱 ETag |
| `Validate` | `func(etag string, data []byte) bool` | 验证 ETag | etag: ETag, data: 数据 | bool: 是否有效 |
| `IsWeak` | `func(etag string) bool` | 检查是否为弱 ETag | etag: ETag | bool: 是否为弱 ETag |

### 5.2 gitignore 模块

#### 5.2.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Pattern` | 表示忽略模式的结构 | `Pattern string`, `Negate bool`, `Directory bool` |
| `Matcher` | 表示匹配器的结构 | `Patterns []Pattern` |

#### 5.2.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `Parse` | `func(s string) []Pattern` | 解析忽略规则 | s: 规则字符串 | []Pattern: 解析出的模式 |
| `ParseFile` | `func(path string) ([]Pattern, error)` | 解析忽略文件 | path: 文件路径 | []Pattern: 解析出的模式, error: 错误信息 |
| `NewMatcher` | `func(patterns []Pattern) *Matcher` | 创建新的匹配器 | patterns: 模式 | *Matcher: 匹配器实例 |
| `(m *Matcher) Match` | `func(path string, isDir bool) bool` | 匹配路径 | path: 路径, isDir: 是否为目录 | bool: 是否匹配 |

### 5.3 vcr 模块

#### 5.3.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Recorder` | 表示请求录制器的结构 | `Mode Mode`, `Cassette *Cassette`, `Matcher Matcher` |
| `Cassette` | 表示录制磁带的结构 | `Interactions []Interaction` |
| `Interaction` | 表示交互的结构 | `Request Request`, `Response Response` |
| `Request` | 表示请求的结构 | `Method string`, `URL string`, `Headers map[string][]string`, `Body []byte` |
| `Response` | 表示响应的结构 | `Status int`, `Headers map[string][]string`, `Body []byte` |

#### 5.3.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `NewRecorder` | `func() *Recorder` | 创建新的录制器 | 无 | *Recorder: 录制器实例 |
| `(r *Recorder) Start` | `func(cassettePath string, mode Mode)` | 开始录制 | cassettePath: 磁带路径, mode: 模式 | error: 错误信息 |
| `(r *Recorder) Stop` | `func() error` | 停止录制 | 无 | error: 错误信息 |
| `(r *Recorder) RoundTrip` | `func(req *http.Request) (*http.Response, error)` | 执行 HTTP 请求 | req: HTTP 请求 | *http.Response: HTTP 响应, error: 错误信息 |

### 5.4 powernap 模块

#### 5.4.1 核心类型

| 类型名 | 说明 | 字段 |
|--------|------|------|
| `Monitor` | 表示睡眠监视器的结构 | `onWake []func()` |

#### 5.4.2 核心函数

| 函数名 | 签名 | 说明 | 参数 | 返回值 |
|--------|------|------|------|--------|
| `NewMonitor` | `func() *Monitor` | 创建新的监视器 | 无 | *Monitor: 监视器实例 |
| `(m *Monitor) Start` | `func() error` | 开始监控 | 无 | error: 错误信息 |
| `(m *Monitor) Stop` | `func()` | 停止监控 | 无 | 无 |
| `(m *Monitor) OnWake` | `func(f func())` | 注册唤醒事件回调 | f: 回调函数 | 无 |

## 6. 接口使用示例

### 6.1 基础层模块示例

#### 6.1.1 ansi 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/ansi"

// 解析 ANSI 转义序列
input := []byte("\x1b[31mHello\x1b[0m")
sequences := ansi.Parse(input)

for _, seq := range sequences {
    fmt.Printf("Prefix: %c\n", seq.Prefix)
    fmt.Printf("Parameters: %v\n", seq.Parameters)
    fmt.Printf("Command: %c\n", seq.Command)
}
```

#### 6.1.2 term 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/term"

// 获取终端大小
width, height, err := term.Size()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Terminal size: %dx%d\n", width, height)

// 设置终端为原始模式
oldState, err := term.SetRaw()
if err != nil {
    log.Fatal(err)
}
defer term.Restore(oldState)

// 读取单个字符
var b [1]byte
_, err = os.Stdin.Read(b[:])
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Read character: %c\n", b[0])
```

### 6.2 核心层模块示例

#### 6.2.1 cellbuf 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/cellbuf"

// 创建缓冲区
width, height := 80, 24
buf := cellbuf.NewBuffer(width, height)

// 创建单元格
cell := cellbuf.NewCell('H')

// 设置单元格
buf.SetCell(0, 0, cell)

// 获取单元格
c := buf.Cell(0, 0)
fmt.Printf("Cell at (0,0): %c\n", c.Rune)

// 清除缓冲区
buf.Clear()

// 插入行
buf.InsertLine(10, 2, nil)

// 删除行
buf.DeleteLine(10, 1, nil)
```

#### 6.2.2 input 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/input"

// 创建输入处理器
in := input.New()

// 开始处理输入
in.Start()
defer in.Stop()

// 接收输入事件
fmt.Println("Press any key...")
event := <-in.Events()

switch event.Type {
case input.Key:
    fmt.Printf("Key pressed: %s\n", event.Key.Name)
case input.Mouse:
    fmt.Printf("Mouse event: %v at (%d,%d)\n", event.Mouse.Type, event.Mouse.X, event.Mouse.Y)
case input.Resize:
    fmt.Println("Terminal resized")
}
```

### 6.3 应用层模块示例

#### 6.3.1 pony 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/pony"

// 解析 pony 标记
input := `
# Hello, Pony!

This is **bold** and *italic* text.

[Click here](https://example.com) to visit example.com.

`

nodes := pony.Parse(input)

// 渲染为终端输出
output := pony.Render(nodes)
fmt.Println(output)
```

### 6.4 工具层模块示例

#### 6.4.1 etag 模块示例

```go
import "github.com/purpose168/charm-experimental-packages-cn/etag"

// 生成 ETag
data := []byte("Hello, World!")
etagValue := etag.Generate(data)
fmt.Printf("Generated ETag: %s\n", etagValue)

// 验证 ETag
valid := etag.Validate(etagValue, data)
fmt.Printf("ETag valid: %t\n", valid)

// 生成弱 ETag
weakEtag := etag.GenerateWeak(data)
fmt.Printf("Generated weak ETag: %s\n", weakEtag)
fmt.Printf("Is weak: %t\n", etag.IsWeak(weakEtag))
```

## 7. 接口设计原则

### 7.1 设计理念

1. **简洁明了**：接口设计简洁明了，易于理解和使用
2. **一致性**：接口命名和使用方式保持一致
3. **可扩展性**：接口设计考虑未来扩展，避免破坏性变更
4. **错误处理**：合理的错误处理机制，提供清晰的错误信息
5. **跨平台兼容**：接口在不同平台上保持一致的行为

### 7.2 命名规范

1. **类型名**：使用 PascalCase，如 `Buffer`、`Input`
2. **函数名**：使用 PascalCase，如 `NewBuffer`、`Parse`
3. **方法名**：使用 PascalCase，如 `SetCell`、`Start`
4. **变量名**：使用 camelCase，如 `width`、`height`
5. **常量名**：使用 ALL_CAPS，如 `MaxCellWidth`、`DefaultWidth`

### 7.3 最佳实践

1. **接口隔离**：每个接口只负责单一功能
2. **依赖注入**：通过参数传递依赖，而不是硬编码
3. **返回错误**：函数通过返回值传递错误，而不是通过全局变量或 panic
4. **文档注释**：为每个接口提供详细的文档注释
5. **示例代码**：为重要接口提供使用示例

## 8. 总结

`charm-experimental-packages-cn` 项目提供了一套完整、一致、易用的 API 接口，涵盖了终端处理的各个方面。这些接口设计遵循了 Go 语言的最佳实践，具有以下特点：

1. **模块化**：每个模块提供独立的功能，通过清晰的接口与其他模块交互
2. **一致性**：接口命名和使用方式保持一致，便于理解和使用
3. **可扩展性**：接口设计考虑未来扩展，避免破坏性变更
4. **跨平台兼容**：接口在不同平台上保持一致的行为
5. **详细的文档**：为每个接口提供详细的文档注释和使用示例

通过这些 API 接口，开发者可以方便地构建各种终端应用和命令行工具，无需关心底层的实现细节。同时，项目的模块化设计也使得开发者可以根据需要只导入所需的模块，减少不必要的依赖。

随着项目的不断发展和完善，API 接口也会不断演进，以适应新的需求和技术趋势。但项目会保持 API 的向后兼容性，确保现有代码不会因为 API 变更而受到影响。