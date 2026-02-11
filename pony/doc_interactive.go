// Package pony 提供了构建交互式 TUI 组件的示例和模式。
//
// # 构建交互式组件
//
// 当创建渲染其他元素的自定义组件时，必须将组件的 ID 传递给返回的根元素。这确保了在组件中任何位置的鼠标点击都会返回组件的 ID，而不是子元素的 ID。
//
// ## 模式：传递组件 ID
//
//	type MyComponent struct {
//	    pony.BaseElement  // 提供 ID() 和 SetID() 方法
//	    // ... 你的字段
//	}
//
//	func (c *MyComponent) Render() pony.Element {
//	    // 构建你的 UI
//	    root := pony.NewVStack(
//	        pony.NewText(c.label),
//	        pony.NewBox(pony.NewText(c.value)),
//	    )
//
//	    // 关键：在根元素上设置组件的 ID
//	    root.SetID(c.ID())
//
//	    return root
//	}
//
// ## 为什么这很重要
//
// HitTest() 优先选择具有显式 ID 的元素，而不是自动生成的元素。
// 当多个元素在点击点重叠时：
//
//  1. 不使用 SetID：返回带有自动生成 ID（如 "elem_0x123..."）的子元素
//  2. 使用 SetID：返回带有有意义 ID（如 "name-input"）的组件
//
// 这允许你处理整个组件的点击，而不是单个子元素的点击。
//
// ## 示例用法
//
//	input := NewInput("Name:")
//	input.SetID("name-input")
//
//	// 在模板中
//	slots := map[string]pony.Element{
//	    "input": input.Render(),  // 带有 ID "name-input" 的 VStack
//	}
//
//	// 在回调中
//	view.Callback = func(msg tea.Msg) tea.Cmd {
//	    if click, ok := msg.(tea.MouseClickMsg); ok {
//	        elem := boundsMap.HitTest(click.X, click.Y)
//	        // elem.ID() 返回 "name-input" - 正是你想要的！
//	    }
//	}
//
// 请参阅 examples/interactive-form 获取完整的工作示例。
package pony
