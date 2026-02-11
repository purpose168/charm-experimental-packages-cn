package main

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/purpose168/charm-experimental-packages-cn/pony"
)

// TemplateData 表示传递给模板的数据
type TemplateData struct {
	Title       string
	Count       int
	LastClicked string
	HoveredID   string
}

// 定义带有交互按钮的模板
const tmpl = `
<vstack spacing="1">
	<box border="double" border-color="cyan">
		<text font-weight="bold" foreground-color="yellow" alignment="center">{{ .Title }}</text>
	</box>

	<divider foreground-color="gray" />

	<vstack spacing="1">
		<text font-weight="bold">Click Counter: {{ .Count }}</text>
		<text foreground-color="gray">Last clicked: {{ .LastClicked }}</text>
	</vstack>

	<divider foreground-color="gray" />

	<vstack spacing="1">
		<text font-weight="bold">Interactive Buttons:</text>
		<hstack spacing="2">
			<button id="increment-btn" text="Increment" border="rounded" padding="1" />
			<button id="decrement-btn" text="Decrement" border="rounded" padding="1" />
			<button id="reset-btn" text="Reset" border="rounded" padding="1" />
		</hstack>
	</vstack>

	<divider foreground-color="gray" />

	<vstack spacing="1">
		<text font-weight="bold">More Buttons:</text>
		<button id="hello-btn" text="Say Hello" border="thick" padding="1" foreground-color="green" />
		<button id="quit-btn" text="Quit Application" border="thick" padding="1" foreground-color="red" />
	</vstack>

	{{ if ne .HoveredID "" }}
	<divider foreground-color="gray" />
	<text font-style="italic" foreground-color="cyan">Hovering: {{ .HoveredID }}</text>
	{{ end }}

	<text font-style="italic" foreground-color="gray">Click buttons with mouse or press 'q' to quit</text>
</vstack>
`

type model struct {
	template    *pony.Template[TemplateData]
	count       int
	lastClicked string
	hoveredID   string
	width       int
	height      int
}

func initialModel() model {
	return model{
		template:    pony.MustParse[TemplateData](tmpl),
		count:       0,
		lastClicked: "none",
		width:       80,
		height:      24,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestWindowSize,
	)
}

// 按钮点击的自定义消息
type buttonClickMsg string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case buttonClickMsg:
		// 处理按钮点击
		m.lastClicked = string(msg)
		switch msg {
		case "increment-btn":
			m.count++
		case "decrement-btn":
			m.count--
		case "reset-btn":
			m.count = 0
		case "hello-btn":
			m.lastClicked = "hello-btn (Hello!)"
		case "quit-btn":
			return m, tea.Quit
		}

	case hoverMsg:
		m.hoveredID = string(msg)
	}

	return m, nil
}

type hoverMsg string

func (m model) View() tea.View {
	// 准备模板数据
	data := TemplateData{
		Title:       "pony Mouse Click Demo",
		Count:       m.count,
		LastClicked: m.lastClicked,
		HoveredID:   m.hoveredID,
	}

	// 渲染并包含碰撞测试的边界
	scr, boundsMap := m.template.RenderWithBounds(data, nil, m.width, m.height)

	// 创建带有鼠标事件回调的视图
	view := tea.NewView(scr.Render())
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion

	// 设置回调以使用边界映射处理鼠标事件
	view.Callback = func(msg tea.Msg) tea.Cmd {
		switch msg := msg.(type) {
		case tea.MouseClickMsg:
			mouse := msg.Mouse()
			// 碰撞测试以找出被点击的元素
			if elem := boundsMap.HitTest(mouse.X, mouse.Y); elem != nil {
				// 返回带有按钮 ID 的命令
				return func() tea.Msg {
					return buttonClickMsg(elem.ID())
				}
			}

		case tea.MouseMotionMsg:
			mouse := msg.Mouse()
			// 跟踪悬停状态
			if elem := boundsMap.HitTest(mouse.X, mouse.Y); elem != nil {
				return func() tea.Msg {
					return hoverMsg(elem.ID())
				}
			} else {
				return func() tea.Msg {
					return hoverMsg("")
				}
			}
		}
		return nil
	}

	return view
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nThanks for trying pony mouse interactions!")
}
