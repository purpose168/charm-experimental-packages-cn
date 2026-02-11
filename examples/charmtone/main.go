// Package main 是一个用于渲染 CharmTone 调色板的简单命令行工具。
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/purpose168/charm-experimental-packages-cn/exp/charmtone"
	"github.com/spf13/cobra"
)

const (
	blackCircle = "●"
	whiteCircle = "○"
	rightArrow  = "→"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "charmtone",
		Short: "CharmTone 调色板工具",
		Long:  "一个用于以各种格式渲染 CharmTone 调色板的命令行工具",
		Run: func(_ *cobra.Command, _ []string) {
			renderGuide()
		},
	}

	cssCmd := &cobra.Command{
		Use:   "css",
		Short: "生成 CSS 变量",
		Long:  "为 CharmTone 调色板生成 CSS 自定义属性（变量）",
		Run: func(_ *cobra.Command, _ []string) {
			renderCSS()
		},
	}

	scssCmd := &cobra.Command{
		Use:   "scss",
		Short: "以 SCSS 变量形式打印",
		Long:  "打印 CharmTone 调色板的 SCSS 变量",
		Run: func(_ *cobra.Command, _ []string) {
			renderSCSS()
		},
	}

	vimCmd := &cobra.Command{
		Use:   "vim",
		Short: "生成 Vim 配色方案",
		Long:  "使用 CharmTone 调色板生成 Vim 配色方案",
		Run: func(_ *cobra.Command, _ []string) {
			renderVim()
		},
	}

	rootCmd.AddCommand(cssCmd, scssCmd, vimCmd)

	// 使用 Fang 执行命令，提供增强的样式和功能
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}

func renderSCSS() {
	for _, k := range charmtone.Keys() {
		name := strings.ToLower(strings.ReplaceAll(k.String(), " ", "-"))
		fmt.Printf("$%s: %s;\n", name, k.Hex())
	}
}

func renderVim() {
	for _, k := range charmtone.Keys() {
		name := strings.ToLower(strings.ReplaceAll(k.String(), " ", "-"))
		fmt.Printf("let %s = '%s'\n", name, k.Hex())
	}
}

func renderCSS() {
	for _, k := range charmtone.Keys() {
		name := strings.ToLower(strings.ReplaceAll(k.String(), " ", "-"))
		fmt.Printf("--charmtone-%s: %s;\n", name, k.Hex())
	}
}
