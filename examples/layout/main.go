// Package main 演示了使用方法。
package main

// 本示例演示了各种 Lip Gloss 样式和布局功能。

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/cellbuf"
	"github.com/purpose168/charm-experimental-packages-cn/input"
	"github.com/purpose168/charm-experimental-packages-cn/term"
	"github.com/rivo/uniseg"
)

const (
	// 在实际情况中，我们会根据检测到的宽度调整文档大小。
	// 在本示例中，我们硬编码宽度，稍后仅使用检测到的宽度进行截断，以避免锯齿状换行。
	width = 96

	// 布局中各列的渲染宽度。
	columnWidth = 30
)

var (
	// 检测到的背景颜色是否为深色。我们在 init() 中检测。
	hasDarkBG bool

	// 一个辅助函数，用于根据检测到的背景颜色选择亮色或暗色。我们在 init() 中创建。
	lightDark lipgloss.LightDarkFunc
)

func init() {
	// 检测背景颜色。
	hasDarkBG = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)

	// 创建一个新的辅助函数，用于根据检测到的背景颜色选择亮色或暗色。
	lightDark = lipgloss.LightDark(hasDarkBG)
}

func main() {
	// 样式定义。
	var (

		// 通用样式。

		subtle    = lightDark(lipgloss.Color("#D9DCCF"), lipgloss.Color("#383838"))
		highlight = lightDark(lipgloss.Color("#874BFD"), lipgloss.Color("#7D56F4"))
		special   = lightDark(lipgloss.Color("#43BF6D"), lipgloss.Color("#73F59F"))

		divider = lipgloss.NewStyle().
			SetString("•").
			Padding(0, 1).
			Foreground(subtle).
			String()

		url = lipgloss.NewStyle().Foreground(special).Render

		// 标签样式。

		activeTabBorder = lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┘",
			BottomRight: "└",
		}

		tabBorder = lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┴",
			BottomRight: "┴",
		}

		tab = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlight).
			Padding(0, 1)

		activeTab = tab.Border(activeTabBorder, true)

		tabGap = tab.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)

		// 标题样式。

		titleStyle = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(5).
				Padding(0, 1).
				Italic(true).
				Foreground(lipgloss.Color("#FFF7DB")).
				SetString("Lip Gloss")

		descStyle = lipgloss.NewStyle().MarginTop(1)

		infoStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				BorderForeground(subtle)

		// 对话框样式。

		dialogBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(1, 0).
				BorderTop(true).
				BorderLeft(true).
				BorderRight(true).
				BorderBottom(true)

		buttonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#888B7E")).
				Padding(0, 3).
				MarginTop(1)

		activeButtonStyle = buttonStyle.
					Foreground(lipgloss.Color("#FFF7DB")).
					Background(lipgloss.Color("#F25D94")).
					MarginRight(2).
					Underline(true)

		// 列表样式。

		list = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(subtle).
			MarginRight(2).
			Height(8).
			Width(columnWidth + 1)

		listHeader = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(subtle).
				MarginRight(2).
				Render

		listItem = lipgloss.NewStyle().PaddingLeft(2).Render

		checkMark = lipgloss.NewStyle().SetString("✓").
				Foreground(special).
				PaddingRight(1).
				String()

		listDone = func(s string) string {
			return checkMark + lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(lightDark(lipgloss.Color("#969B86"), lipgloss.Color("#696969"))).
				Render(s)
		}

		// 段落/历史记录样式。

		historyStyle = lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(highlight).
				Margin(1, 3, 0, 0).
				Padding(1, 2).
				Height(19).
				Width(columnWidth)

		// 状态栏样式。

		statusNugget = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFDF5")).
				Padding(0, 1)

		statusBarStyle = lipgloss.NewStyle().
				Foreground(lightDark(lipgloss.Color("#343433"), lipgloss.Color("#C1C6B2"))).
				Background(lightDark(lipgloss.Color("#D9DCCF"), lipgloss.Color("#353533")))

		statusStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
				Padding(0, 1).
				MarginRight(1)

		encodingStyle = statusNugget.
				Background(lipgloss.Color("#A550DF")).
				Align(lipgloss.Right)

		statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

		fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))

		// 页面样式。

		docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	)

	physicalWidth, physicalHeight, _ := term.GetSize(os.Stdout.Fd())
	doc := strings.Builder{}

	// 标签部分。
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Lip Gloss"),
			tab.Render("Blush"),
			tab.Render("Eye Shadow"),
			tab.Render("Mascara"),
			tab.Render("Foundation"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// 标题部分。
	{
		var (
			colors = colorGrid(1, 5)
			title  strings.Builder
		)

		for i, v := range colors {
			const offset = 2
			c := lipgloss.Color(v[0])
			fmt.Fprint(&title, titleStyle.MarginLeft(i*offset).Background(c))
			if i < len(colors)-1 {
				title.WriteRune('\n')
			}
		}

		desc := lipgloss.JoinVertical(lipgloss.Left,
			descStyle.Render("美观终端布局的样式定义"),
			infoStyle.Render("来自 Charm"+divider+url("https://github.com/charmbracelet/lipgloss")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n\n")
	}

	// 对话框部分。
	okButton := activeButtonStyle.Render("是")
	cancelButton := buttonStyle.Render("也许")

	grad := applyGradient(
		lipgloss.NewStyle(),
		"你确定要吃橘子酱吗？",
		lipgloss.Color("#EDFF82"),
		lipgloss.Color("#F25D94"),
	)

	question := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Center).
		Render(grad)

	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	dialogUI := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	dialog := lipgloss.Place(width, 9,
		lipgloss.Center, lipgloss.Center,
		"",
		// dialogBoxStyle.Render(dialogUi),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(subtle)),
	)

	doc.WriteString(dialog + "\n\n")

	// 颜色网格部分。
	colors := func() string {
		colors := colorGrid(14, 8)

		b := strings.Builder{}
		for _, x := range colors {
			for _, y := range x {
				s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
				b.WriteString(s.String())
			}
			b.WriteRune('\n')
		}

		return b.String()
	}()

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		list.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("尝试的柑橘类水果"),
				listDone("西柚"),
				listDone("柚子"),
				listItem("香橼"),
				listItem("金桔"),
				listItem("柚子"),
			),
		),
		list.Width(columnWidth).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listHeader("实际的唇彩供应商"),
				listItem("Glossier"),
				listItem("Claire‘s Boutique"),
				listDone("Nyx"),
				listItem("Mac"),
				listDone("Milk"),
			),
		),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists, colors))

	// 橘子酱历史部分。
	{
		const (
			historyA = "罗马人从希腊人那里了解到，用蜂蜜慢煮的榅桲在冷却后会\"凝固\"。《阿皮基乌斯》中记载了一种保存完整榅桲（连梗带叶）的方法，将其浸泡在用水稀释的蜂蜜中：这就是罗马橘子酱。榅桲和柠檬的蜜饯（以及玫瑰、苹果、李子和梨）出现在拜占庭皇帝君士坦丁七世·波菲罗格涅图斯的《礼仪书》中。"
			historyB = "中世纪的榅桲蜜饯，在法语中称为 cotignac，有清澈版本和果肉版本，在 16 世纪开始失去其中世纪的香料调味。在 17 世纪，拉瓦雷恩提供了制作浓稠和清澈 cotignac 的食谱。"
			historyC = "1524 年，英格兰国王亨利八世收到了埃克塞特的赫尔先生送的一盒\"橘子酱\"。这可能是 marmelada，一种来自葡萄牙的固体榅桲酱，至今仍在南欧制作和销售。它成为安妮·博林和她的侍女们的最爱。"
		)

		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			historyStyle.Align(lipgloss.Right).Render(historyA),
			historyStyle.Align(lipgloss.Center).Render(historyB),
			historyStyle.MarginRight(0).Render(historyC),
		))

		doc.WriteString("\n\n")
	}

	// 状态栏部分。
	{
		w := lipgloss.Width

		lightDarkState := "浅色"
		if hasDarkBG {
			lightDarkState = "深色"
		}

		statusKey := statusStyle.Render("状态")
		encoding := encodingStyle.Render("UTF-8")
		fishCake := fishCakeStyle.Render("🍥 鱼饼")
		statusVal := statusText.
			Width(width - w(statusKey) - w(encoding) - w(fishCake)).
			Render("令人陶醉的" + lightDarkState + "模式！")

		bar := lipgloss.JoinHorizontal(lipgloss.Top,
			statusKey,
			statusVal,
			encoding,
			fishCake,
		)

		doc.WriteString(statusBarStyle.Width(width).Render(bar))
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	termType := os.Getenv("TERM")
	scr := cellbuf.NewScreen(os.Stdout, physicalWidth, physicalHeight, &cellbuf.ScreenOptions{
		Term:      termType,
		Profile:   colorprofile.Detect(os.Stdout, os.Environ()),
		AltScreen: true,
	})

	defer scr.Close() //nolint:errcheck

	// 启用鼠标事件。
	modes := []ansi.Mode{
		ansi.ButtonEventMouseMode,
		ansi.SgrExtMouseMode,
	}

	os.Stdout.WriteString(ansi.SetMode(modes...))         //nolint:errcheck,gosec
	defer os.Stdout.WriteString(ansi.ResetMode(modes...)) //nolint:errcheck

	state, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		log.Fatalf("设置为原始模式: %v", err)
	}

	defer term.Restore(os.Stdin.Fd(), state) //nolint:errcheck

	drv, err := input.NewReader(os.Stdin, termType, 0)
	if err != nil {
		log.Fatalf("创建输入驱动: %v", err)
	}

	dialogWidth := lipgloss.Width(dialogUI) + dialogBoxStyle.GetHorizontalFrameSize()
	dialogHeight := lipgloss.Height(dialogUI) + dialogBoxStyle.GetVerticalFrameSize()
	dialogX, dialogY := physicalWidth/2-dialogWidth/2-docStyle.GetVerticalFrameSize()-1, 12
	scrw := cellbuf.NewScreenWriter(scr)
	render := func() {
		scr.Clear()
		scrw.SetContent(docStyle.Render(doc.String()))
		box := cellbuf.Rect(dialogX, dialogY, dialogWidth, dialogHeight)
		scrw.SetContentRect(dialogBoxStyle.Render(dialogUI), box)
		scr.Render()
		scr.Flush() //nolint:errcheck,gosec
	}

	// 首次渲染
	render()

	for {
		evs, err := drv.ReadEvents()
		if err != nil {
			log.Fatalf("读取事件: %v", err)
		}

		for _, ev := range evs {
			switch ev := ev.(type) {
			case input.WindowSizeEvent:
				scr.Resize(ev.Width, ev.Height)
			case input.MouseClickEvent:
				dialogX, dialogY = ev.X, ev.Y
			case input.KeyPressEvent:
				switch ev.String() {
				case "ctrl+c", "q":
					return
				case "left", "h":
					dialogX--
				case "down", "j":
					dialogY++
				case "up", "k":
					dialogY--
				case "right", "l":
					dialogX++
				}
			}
		}

		render()
	}
}

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	grid := make([][]string, ySteps)
	for x := range ySteps {
		y0 := x0[x]
		grid[x] = make([]string, xSteps)
		for y := range xSteps {
			grid[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return grid
}

// applyGradient 对给定的字符串应用渐变效果。
func applyGradient(base lipgloss.Style, input string, from, to color.Color) string {
	// 我们想要获取输入字符串的字形，即人类看到的字符数量。
	//
	// 我们绝对不想使用 len()，因为它返回的是字节数。
	// 符文计数会更接近，但在某些情况下，比如表情符号，符文计数会大于实际字符数。
	g := uniseg.NewGraphemes(input)
	var chars []string
	for g.Next() {
		chars = append(chars, g.Str())
	}

	// 生成混合色。
	a, _ := colorful.MakeColor(to)
	b, _ := colorful.MakeColor(from)
	var output strings.Builder
	var hex string
	for i := range chars {
		hex = a.BlendLuv(b, float64(i)/float64(len(chars)-1)).Hex()
		output.WriteString(base.Foreground(lipgloss.Color(hex)).Render(chars[i]))
	}

	return output.String()
}

func init() {
	f, err := os.OpenFile("layout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666) //nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
}