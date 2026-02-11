// Package pony 提供一种声明式、类型安全的标记语言，用于构建终端用户界面，使用 Ultraviolet 作为渲染引擎。
//
// ⚠️ 实验性：这是一个实验性项目，主要由 AI 生成，用于探索声明式 TUI 框架。使用风险自负。
//
// pony 允许您使用熟悉的类似 XML 的标记语法定义 TUI 布局，并结合 Go 模板处理动态内容。它与 Bubble Tea 无缝集成，用于应用程序生命周期管理，同时利用 Ultraviolet 的高效基于单元格的渲染。
//
// # 基本示例
//
//	type ViewData struct {
//	    Title   string
//	    Content string
//	}
//
//	const tmpl = `
//	<vstack spacing="1">
//	  <box border="rounded">
//	    <text font-weight="bold" foreground-color="cyan">{{ .Title }}</text>
//	  </box>
//	  <text>{{ .Content }}</text>
//	</vstack>
//	`
//
//	t := pony.MustParse[ViewData](tmpl)
//	data := ViewData{
//	    Title:   "Hello World",
//	    Content: "Welcome to pony!",
//	}
//	output := t.Render(data, 80, 24)
//
// # 元素
//
//   - vstack: 垂直堆叠容器，带有间距和对齐方式
//   - hstack: 水平堆叠容器，带有间距和对齐方式
//   - text: 文本内容，带有样式和对齐方式
//   - box: 带有边框和内边距的容器
//   - scrollview: 带有滚动条的可滚动视口
//   - divider: 水平或垂直分隔线
//   - spacer: 灵活或固定的空白空间
//   - slot: 动态内容的占位符
//
// # 样式
//
// 文本元素支持精细的样式属性：
//
//	<text foreground-color="cyan" background-color="#1a1b26" font-weight="bold" font-style="italic">Styled text</text>
//
// 对于编程式样式设置，使用 Text 方法：
//
//	text := pony.NewText("Hello").
//	    ForegroundColor(pony.Hex("#FF5555")).
//	    Bold().
//	    Italic()
//
// # 自定义组件
//
// 向组件注册表注册自定义组件：
//
//	pony.Register("card", func(props Props, children []Element) Element {
//	    return pony.NewBox(
//	        pony.NewVStack(children...),
//	    ).Border("rounded").Padding(1)
//	})
//
// 在标记中使用：
//
//	<card><text>Content</text></card>
//
// # 有状态组件
//
// 对管理自己状态的有状态组件使用插槽：
//
//	type Input struct {
//	    value string
//	}
//
//	func (i *Input) Update(msg tea.Msg) { /* handle events */ }
//
//	func (i *Input) Render() pony.Element {
//	    return pony.NewBox(pony.NewText(i.value)).Border("rounded")
//	}
//
// 带插槽的模板：
//
//	<vstack>
//	  <text>Enter name:</text>
//	  <slot name="input" />
//	</vstack>
//
// 使用插槽渲染：
//
//	slots := map[string]pony.Element{
//	    "input": m.inputComp.Render(),
//	}
//	output := tmpl.RenderWithSlots(data, slots, width, height)
//
// # Bubble Tea 集成
//
//	type model struct {
//	    template *pony.Template[ViewData]
//	    width    int
//	    height   int
//	}
//
//	func (m model) Init() tea.Cmd {
//	    return tea.RequestWindowSize
//	}
//
//	func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	    switch msg := msg.(type) {
//	    case tea.WindowSizeMsg:
//	        m.width = msg.Width
//	        m.height = msg.Height
//	    }
//	    return m, nil
//	}
//
//	func (m model) View() tea.View {
//	    data := ViewData{...}
//	    output := m.template.Render(data, m.width, m.height)
//	    return tea.NewView(output)
//	}
package pony
