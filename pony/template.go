package pony

import (
	"bytes"
	"fmt"
	"maps"
	"strings"
	"text/template"

	uv "github.com/charmbracelet/ultraviolet"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Template 是一个类型安全的 pony 模板，可以使用类型 T 的数据进行渲染。
type Template[T any] struct {
	markup   string
	goTmpl   *template.Template
	cacheKey string
}

// Parse 将 pony 标记解析为类型安全的模板。
// 标记可以包含 Go 模板语法，如 {{ .Variable }}。
func Parse[T any](markup string) (*Template[T], error) {
	return ParseWithFuncs[T](markup, nil)
}

// ParseWithFuncs 使用自定义模板函数解析 pony 标记。
func ParseWithFuncs[T any](markup string, funcs template.FuncMap) (*Template[T], error) {
	t := &Template[T]{
		markup:   markup,
		cacheKey: markup,
	}

	// 创建带有内置函数的 Go 模板
	tmplFuncs := defaultTemplateFuncs()
	maps.Copy(tmplFuncs, funcs)

	goTmpl, err := template.New("pony").Funcs(tmplFuncs).Parse(markup)
	if err != nil {
		return nil, fmt.Errorf("template parse: %w", err)
	}
	t.goTmpl = goTmpl

	return t, nil
}

// MustParse 解析 pony 标记，出错时会 panic。
func MustParse[T any](markup string) *Template[T] {
	t, err := Parse[T](markup)
	if err != nil {
		panic(err)
	}
	return t
}

// MustParseWithFuncs 使用自定义函数解析 pony 标记，出错时会 panic。
func MustParseWithFuncs[T any](markup string, funcs template.FuncMap) *Template[T] {
	t, err := ParseWithFuncs[T](markup, funcs)
	if err != nil {
		panic(err)
	}
	return t
}

// Render 使用给定的数据将模板渲染到指定的视口大小。
func (t *Template[T]) Render(data T, width, height int) string {
	scr, _ := t.RenderWithBounds(data, nil, width, height)
	str := scr.Render()
	return strings.ReplaceAll(str, "\r\n", "\n")
}

// RenderWithBounds 渲染模板并返回屏幕缓冲区和边界映射。
// 边界映射可用于事件处理程序中的鼠标点击测试。
func (t *Template[T]) RenderWithBounds(data T, slots map[string]Element, width, height int) (uv.ScreenBuffer, *BoundsMap) {
	// 首先执行 Go 模板
	var buf bytes.Buffer
	if err := t.goTmpl.Execute(&buf, data); err != nil {
		errScreen := uv.NewScreenBuffer(width, 1)
		return errScreen, NewBoundsMap()
	}

	processedMarkup := buf.String()

	// 解析处理后的标记
	root, err := parse(processedMarkup)
	if err != nil {
		errScreen := uv.NewScreenBuffer(width, 1)
		return errScreen, NewBoundsMap()
	}

	// 转换为元素树
	elem := root.toElement()
	if elem == nil {
		emptyScreen := uv.NewScreenBuffer(width, height)
		return emptyScreen, NewBoundsMap()
	}

	// 用提供的元素填充插槽
	if slots != nil {
		fillSlots(elem, slots)
	}

	// 布局元素
	constraints := Constraints{
		MinWidth:  0,
		MaxWidth:  width,
		MinHeight: 0,
		MaxHeight: height,
	}
	size := elem.Layout(constraints)

	// 使用计算大小和请求大小中的较小值
	if size.Width > width {
		size.Width = width
	}
	if size.Height > height {
		size.Height = height
	}

	// 创建缓冲区并渲染
	uvBuf := uv.NewScreenBuffer(size.Width, size.Height)
	area := uv.Rect(0, 0, size.Width, size.Height)
	elem.Draw(uvBuf, area)

	// 构建边界映射
	boundsMap := NewBoundsMap()
	walkAndRegister(elem, boundsMap)

	return uvBuf, boundsMap
}

// RenderWithSlots 使用数据和插槽元素渲染模板。
// 插槽允许将有状态组件注入到模板中。
func (t *Template[T]) RenderWithSlots(data T, slots map[string]Element, width, height int) string {
	scr, _ := t.RenderWithBounds(data, slots, width, height)
	str := scr.Render()
	return strings.ReplaceAll(str, "\r\n", "\n")
}

// fillSlots 递归地用对应的元素填充插槽元素。
func fillSlots(elem Element, slots map[string]Element) {
	if slot, ok := elem.(*Slot); ok {
		if slotElem, found := slots[slot.Name]; found {
			slot.setElement(slotElem)
		}
		return
	}

	// 递归地填充子元素中的插槽
	for _, child := range elem.Children() {
		fillSlots(child, slots)
	}
}

// defaultTemplateFuncs 返回默认的模板函数。
func defaultTemplateFuncs() template.FuncMap {
	titleCaser := cases.Title(language.English)
	return template.FuncMap{
		// 字符串函数
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": titleCaser.String,
		"trim":  strings.TrimSpace,
		"join":  strings.Join,

		// 数学函数
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},

		// 格式化
		"printf": fmt.Sprintf,
		"repeat": strings.Repeat,

		// 颜色辅助函数
		"colorHex": func(hex string) string {
			return fmt.Sprintf("fg:%s", hex)
		},
		"bgHex": func(hex string) string {
			return fmt.Sprintf("bg:%s", hex)
		},
	}
}
