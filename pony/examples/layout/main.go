package main

import (
	"fmt"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/pony"
	"github.com/purpose168/charm-experimental-packages-cn/term"
)

// getSize 获取终端的宽度和高度
// 如果获取失败，返回默认值 80x24
func getSize() (int, int) {
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 80, 24
	}
	return width, height
}

// main 函数是程序的入口点
// 它解析 pony 模板并渲染布局展示
func main() {
	const tmpl = `
<vstack spacing="1">
	<box border="double">
		<text>pony 布局展示</text>
	</box>

	<divider />

	<vstack spacing="0">
		<text>垂直堆栈演示：</text>
		<box border="normal">
			<text>项目 1</text>
		</box>
		<box border="normal">
			<text>项目 2</text>
		</box>
		<box border="normal">
			<text>项目 3</text>
		</box>
	</vstack>

	<divider />

	<vstack spacing="0">
		<text>水平堆栈演示：</text>
		<hstack spacing="2">
			<box border="rounded">
				<text>左侧框</text>
			</box>
			<box border="rounded">
				<text>中间框</text>
			</box>
			<box border="rounded">
				<text>右侧框</text>
			</box>
		</hstack>
	</vstack>

	<divider />

	<vstack spacing="0">
		<text>宽度/高度属性演示：</text>
		<hstack spacing="1">
			<box border="normal" width="30%">
				<text>30% 宽度</text>
			</box>
			<box border="normal" width="70%">
				<text>70% 宽度</text>
			</box>
		</hstack>
	</vstack>

	<divider />

	<vstack spacing="0">
		<text>固定大小演示：</text>
		<hstack spacing="1">
			<box border="normal" width="20">
				<text>固定 20</text>
			</box>
			<box border="normal" width="30">
				<text>固定 30 单元格宽</text>
			</box>
		</hstack>
	</vstack>

	<divider />

	<vstack spacing="0">
		<text>嵌套布局演示：</text>
		<box border="thick">
			<vstack spacing="1">
				<hstack spacing="1">
					<box border="normal">
						<text>A</text>
					</box>
					<box border="normal">
						<text>B</text>
					</box>
				</hstack>
				<hstack spacing="1">
					<box border="normal">
						<text>C</text>
					</box>
					<box border="normal">
						<text>D</text>
					</box>
				</hstack>
			</vstack>
		</box>
	</vstack>

	<divider />

	<text>边框样式：</text>
	<hstack spacing="1">
		<box border="normal">
			<text>普通</text>
		</box>
		<box border="rounded">
			<text>圆角</text>
		</box>
		<box border="thick">
			<text>粗体</text>
		</box>
		<box border="double">
			<text>双边框</text>
		</box>
	</hstack>
</vstack>
`

	t := pony.MustParse[any](tmpl)
	w, h := getSize()
	output := t.Render(nil, w, h)
	fmt.Print(output)
}
