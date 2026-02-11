package main

import (
	"fmt"
	"image/color"
	"io"
	"iter"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/purpose168/charm-experimental-packages-cn/exp/charmtone"
)

func renderGuide() {
	// 找到最长的键名。
	var widestKeyName int
	for _, k := range charmtone.Keys() {
		if w := lipgloss.Width(k.String()); w > widestKeyName {
			widestKeyName = w
		}
	}

	// 样式。
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := lipgloss.LightDark(hasDarkBG)
	logo := lipgloss.NewStyle().
		Foreground(charmtone.Ash).
		Background(charmtone.Charple).
		Padding(0, 1).
		SetString("Charm™")
	title := lipgloss.NewStyle().
		Foreground(lightDark(charmtone.Charcoal, charmtone.Smoke))
	subdued := lipgloss.NewStyle().
		Foreground(lightDark(charmtone.Squid, charmtone.Oyster))
	fg := lipgloss.NewStyle().
		MarginLeft(2).
		Width(widestKeyName).
		Align(lipgloss.Right)
	bg := lipgloss.NewStyle().
		Width(8)
	hex := lipgloss.NewStyle().
		Foreground(lightDark(charmtone.Smoke, charmtone.Charcoal))
	legend := subdued.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(charmtone.Charcoal).
		Padding(0, 2).
		MarginLeft(2)
	primaryMark := lipgloss.NewStyle().
		Foreground(lightDark(charmtone.Squid, charmtone.Smoke)).
		SetString(blackCircle)
	secondaryMark := primaryMark.
		Foreground(lightDark(charmtone.Squid, charmtone.Oyster)).
		SetString(blackCircle)
	tertiaryMark := primaryMark.
		Foreground(lightDark(charmtone.Squid, charmtone.Oyster)).
		SetString(whiteCircle)
	rightArrowMark := lipgloss.NewStyle().
		Foreground(lightDark(charmtone.Squid, charmtone.Oyster)).
		Margin(0, 1).
		SetString(rightArrow)

	var b strings.Builder

	// 渲染标题和描述。
	fmt.Fprintf(
		&b,
		"\n  %s %s %s\n\n",
		logo.String(),
		title.Render("CharmTone"),
		subdued.Render("• 配色指南"),
	)

	// 渲染色板及其元数据。
	renderSwatch := func(w io.Writer, k charmtone.Key) {
		mark := " "
		switch {
		case k.IsPrimary():
			mark = primaryMark.String()
		case k.IsSecondary():
			mark = secondaryMark.String()
		case k.IsTertiary():
			mark = tertiaryMark.String()
		}
		_, _ = fmt.Fprintf(w,
			"%s %s %s %s",
			fg.Foreground(k).Render(k.String()),
			mark,
			bg.Background(k).Render(),
			hex.Render(k.Hex()),
		)
	}

	// 渲染主色块。
	for i := charmtone.Cumin; i < charmtone.Pepper; i++ {
		k := charmtone.Keys()[i]
		renderSwatch(&b, k)
		if i%3 == 2 {
			b.WriteRune('\n')
		} else {
			b.WriteRune(' ')
		}
	}

	// 获取到目前为止的总块宽度。
	var totalWidth int
	for l := range SplitSeq(b.String(), "\n") {
		if w := lipgloss.Width(l); w > totalWidth {
			totalWidth = w
		}
	}

	// 灰度块。
	var grays strings.Builder
	for i := charmtone.Pepper; i <= charmtone.Butter; i++ {
		k := charmtone.Keys()[i]
		renderSwatch(&grays, k)
		if i < charmtone.Butter {
			grays.WriteRune('\n')
		}
	}

	// 获取灰度块的宽度。
	var grayWidth int
	for l := range SplitSeq(grays.String(), "\n") {
		if w := lipgloss.Width(l); w > grayWidth {
			grayWidth = w
		}
	}

	fmt.Fprint(&b, "\n")

	// 构建图例。
	legendBlock := legend.Render(
		strings.Join([]string{
			primaryMark.String() + subdued.Render(" 主色"),
			secondaryMark.String() + subdued.Render(" 次色"),
			tertiaryMark.String() + subdued.Render(" 第三色"),
		}, "  "),
	)

	// 构建渐变。
	var grads strings.Builder
	gap := "  "
	gapWidth := lipgloss.Width(gap)
	{
		fullWidth := (totalWidth - grayWidth) - lipgloss.Width(gap)
		if fullWidth%2 != 0 {
			fullWidth--
		}
		halfWidth := fullWidth / gapWidth

		block := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(halfWidth)
		s := subdued.
			Foreground(charmtone.Squid)

		// 渐变：
		// Hazy -> Blush, Bok -> Zest
		left := blendKeys(halfWidth, charmtone.Hazy, charmtone.Blush)
		left += "\n" + block.Render(s.Render("Hazy")+rightArrowMark.String()+s.Render("Blush"))
		right := blendKeys(halfWidth, charmtone.Bok, charmtone.Zest)
		right += "\n" + block.Render(s.Render("Bok")+rightArrowMark.String()+s.Render("Zest"))
		fmt.Fprint(&grads, "\n", lipgloss.JoinHorizontal(lipgloss.Top, gap, left, right))

		// 渐变：
		// Uni -> Coral -> Tuna -> Violet -> Malibu -> Turtle
		block = block.Width(fullWidth)
		buf := strings.Builder{}
		fmt.Fprint(&buf, blendKeys(fullWidth, charmtone.Uni,
			charmtone.Coral, charmtone.Tuna, charmtone.Violet,
			charmtone.Malibu, charmtone.Turtle,
		))
		fmt.Fprint(&buf, "\n",
			block.Render(
				s.Render("Uni")+rightArrowMark.String()+
					s.Render("Coral")+rightArrowMark.String()+
					s.Render("Tuna")+rightArrowMark.String()+
					s.Render("Violet")+rightArrowMark.String()+
					s.Render("Malibu")+rightArrowMark.String()+
					s.Render("Turtle"),
			),
		)
		fmt.Fprint(&grads, "\n\n", lipgloss.JoinHorizontal(lipgloss.Top, gap, buf.String()))
	}

	// 连接灰度和图例。
	fmt.Fprint(&b, lipgloss.JoinHorizontal(lipgloss.Top, grays.String(), " ", grads.String()))

	fmt.Fprint(&b, "\n\n", legendBlock, "\n\n")

	// 输出。
	_, err := lipgloss.Print(b.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "打印错误: %v\n", err)
		os.Exit(1)
	}
}

func blendKeys(size int, keys ...charmtone.Key) string {
	stops := make([]color.Color, len(keys))
	for i := range keys {
		stops[i] = color.Color(keys[i])
	}

	var w strings.Builder
	for _, c := range lipgloss.Blend1D(size, stops...) {
		fmt.Fprint(&w, lipgloss.NewStyle().Background(c).Render(" "))
	}
	return w.String()
}

// SplitSeq 返回一个由分隔符分隔的字符串子串的迭代器。
//
// 这是 strings.Split 的 Go 1.23 兼容版本。一旦我们支持 Go 1.24，
// 就可以删除它，转而使用标准库中的 strings.SplitSeq。
func SplitSeq(s, sep string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if sep == "" {
			for _, r := range s {
				if !yield(string(r)) {
					return
				}
			}
			return
		}

		start := 0
		for {
			i := strings.Index(s[start:], sep)
			if i == -1 {
				if !yield(s[start:]) {
					return
				}
				break
			}
			if !yield(s[start : start+i]) {
				return
			}
			start += i + len(sep)
		}
	}
}
