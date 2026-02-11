// 包 main 演示 cellbuf 的使用方法。
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/cellbuf"
	"github.com/purpose168/charm-experimental-packages-cn/input"
	"github.com/purpose168/charm-experimental-packages-cn/term"
)

func main() {
	w, h, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		log.Fatalf("获取终端大小: %v", err)
	}

	state, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		log.Fatalf("设置为原始模式: %v", err)
	}

	defer term.Restore(os.Stdin.Fd(), state) //nolint:errcheck

	const altScreen = true
	if !altScreen {
		h = 10
	}

	termType := os.Getenv("TERM")
	scr := cellbuf.NewScreen(os.Stdout, w, h, &cellbuf.ScreenOptions{
		Term:           termType,
		RelativeCursor: !altScreen,
		AltScreen:      altScreen,
	})

	defer scr.Close() //nolint:errcheck

	drv, err := input.NewReader(os.Stdin, termType, 0)
	if err != nil {
		log.Fatalf("创建输入驱动: %v", err)
	}

	modes := []ansi.Mode{
		ansi.ButtonEventMouseMode,
		ansi.SgrExtMouseMode,
	}

	os.Stdout.WriteString(ansi.SetMode(modes...))         //nolint:errcheck,gosec
	defer os.Stdout.WriteString(ansi.ResetMode(modes...)) //nolint:errcheck

	x, y := (w/2)-10, h/2

	text := ansi.SetHyperlink("https://charm.sh") +
		ansi.Style{}.Reverse(true).Styled(" !你好，世界! ") +
		ansi.ResetHyperlink()
	scrw := cellbuf.NewScreenWriter(scr)
	render := func() {
		scr.Fill(cellbuf.NewCell('你'))
		scrw.PrintCropAt(x, y, text, "")
		scr.Render()
		scr.Flush() //nolint:errcheck,gosec
	}

	resize := func(nw, nh int) {
		if !altScreen {
			nh = h
			w = nw
		}
		scr.Resize(nw, nh)
		render()
	}

	if runtime.GOOS != "windows" {
		// 监听窗口大小调整事件
		go listenForResize(func() {
			nw, nh, _ := term.GetSize(os.Stdout.Fd())
			resize(nw, nh)
		})
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
				resize(ev.Width, ev.Height)
			case input.MouseClickEvent:
				x, y = ev.X, ev.Y
			case input.KeyPressEvent:
				switch ev.String() {
				case "ctrl+c", "q":
					return
				case "left", "h":
					x--
				case "down", "j":
					y++
				case "up", "k":
					y--
				case "right", "l":
					x++
				}
			}
		}

		render()
	}
}

func init() {
	f, err := os.OpenFile("cellbuf.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666) //nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
}
