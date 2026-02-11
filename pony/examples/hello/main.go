package main

import (
	"fmt"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/pony"
	"github.com/purpose168/charm-experimental-packages-cn/term"
)

func main() {
	const tmpl = `
<vstack spacing="1">
	<box border="rounded">
		<text>Hello, World!</text>
	</box>
	<text>Welcome to pony - a declarative markup language for terminal UIs.</text>
	<divider />
	<hstack spacing="2">
		<text>Left</text>
		<text>Right</text>
	</hstack>
</vstack>
`

	// 获取终端大小
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width, height = 80, 24 // 备用值
	}

	t := pony.MustParse[any](tmpl)
	output := t.Render(nil, width, height)
	fmt.Print(output)
}
