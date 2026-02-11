package pony

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"io"
	"strings"
)

// node 表示一个解析后的 XML 元素。
type node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  string     `xml:",chardata"`
	Children []*node    `xml:",any"`
}

// Props 将 XML 属性转换为 Props 映射。
func (n *node) Props() Props {
	props := make(Props)
	for _, attr := range n.Attrs {
		props[attr.Name.Local] = attr.Value
	}
	return props
}

// parse 将 XML 标记解析为节点树。
func parse(markup string) (*node, error) {
	// 如果未包装在根元素中，则添加包装
	wrapped := markup
	if !strings.HasPrefix(strings.TrimSpace(markup), "<") {
		wrapped = "<root>" + markup + "</root>"
	}

	decoder := xml.NewDecoder(strings.NewReader(wrapped))
	decoder.Strict = false // 对 XML 解析宽松处理

	var root node
	if err := decoder.Decode(&root); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("空标记")
		}
		return nil, fmt.Errorf("xml 解码: %w", err)
	}

	// 如果我们添加了包装器，则解包
	if strings.TrimSpace(markup) != wrapped {
		if len(root.Children) == 1 {
			return root.Children[0], nil
		}
		return &root, nil
	}

	return &root, nil
}

// toElement 将 XML 节点转换为 Element。
func (n *node) toElement() Element {
	if n == nil {
		return nil
	}

	// 获取标签名
	tagName := n.XMLName.Local
	props := n.Props()

	var elem Element

	// 首先检查自定义组件注册表
	if factory, ok := GetComponent(tagName); ok {
		elem = factory(props, n.childElements())
	} else {
		// 然后检查内置元素
		switch tagName {
		case "vstack":
			elem = n.toVStack(props)
		case "hstack":
			elem = n.toHStack(props)
		case "zstack":
			elem = n.toZStack(props)
		case "text":
			elem = n.toText(props)
		case "box":
			elem = n.toBox(props)
		case "spacer":
			elem = n.toSpacer(props)
		case "flex":
			elem = n.toFlex(props)
		case "positioned":
			elem = n.toPositioned(props)
		case "divider":
			elem = n.toDivider(props)
		case "slot":
			elem = n.toSlot(props)
		case "scrollview":
			elem = n.toScrollView(props)
		case "":
			// 匿名文本节点（无标签，只有内容）
			content := strings.TrimSpace(n.Content)
			if content != "" {
				elem = NewText(content)
			} else {
				return nil
			}
		default:
			// 未知元素，视为容器
			elem = NewVStack(n.childElements()...)
		}
	}

	// 如果提供了 ID，则设置
	if elem != nil && props.Has("id") {
		if setter, ok := elem.(interface{ SetID(string) }); ok {
			setter.SetID(props.Get("id"))
		}
	}

	return elem
}

// childElements 将子节点转换为 Elements。
func (n *node) childElements() []Element {
	var elements []Element

	for _, child := range n.Children {
		// 检查是否为文本节点
		if child.XMLName.Local == "" && strings.TrimSpace(child.Content) != "" {
			elements = append(elements, NewText(strings.TrimSpace(child.Content)))
			continue
		}

		if elem := child.toElement(); elem != nil {
			elements = append(elements, elem)
		}
	}

	return elements
}

// toVStack 将节点转换为 VStack 元素。
func (n *node) toVStack(props Props) Element {
	spacing := parseIntAttr(props, "spacing", 0)
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))
	alignment := props.GetOr("alignment", AlignmentLeading)

	vstack := NewVStack(n.childElements()...)
	if spacing > 0 {
		vstack = vstack.Spacing(spacing)
	}
	if !width.IsAuto() {
		vstack = vstack.Width(width)
	}
	if !height.IsAuto() {
		vstack = vstack.Height(height)
	}
	if alignment != AlignmentLeading {
		vstack = vstack.Alignment(alignment)
	}

	return vstack
}

// toHStack 将节点转换为 HStack 元素。
func (n *node) toHStack(props Props) Element {
	spacing := parseIntAttr(props, "spacing", 0)
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))
	alignment := props.GetOr("alignment", AlignmentTop)

	hstack := NewHStack(n.childElements()...)
	if spacing > 0 {
		hstack = hstack.Spacing(spacing)
	}
	if !width.IsAuto() {
		hstack = hstack.Width(width)
	}
	if !height.IsAuto() {
		hstack = hstack.Height(height)
	}
	if alignment != AlignmentTop {
		hstack = hstack.Alignment(alignment)
	}

	return hstack
}

// toZStack 将节点转换为 ZStack 元素。
func (n *node) toZStack(props Props) Element {
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))
	alignment := props.GetOr("alignment", AlignmentCenter)
	verticalAlignment := props.GetOr("vertical-alignment", AlignmentCenter)

	zstack := NewZStack(n.childElements()...)
	if !width.IsAuto() {
		zstack = zstack.Width(width)
	}
	if !height.IsAuto() {
		zstack = zstack.Height(height)
	}
	if alignment != AlignmentCenter {
		zstack = zstack.Alignment(alignment)
	}
	if verticalAlignment != AlignmentCenter {
		zstack = zstack.VerticalAlignment(verticalAlignment)
	}

	return zstack
}

// toText 将节点转换为 Text 元素。
func (n *node) toText(props Props) Element {
	// 从内容和子节点收集文本
	var text string

	if n.Content != "" {
		text = strings.TrimSpace(n.Content)
	}

	// 也从子文本节点收集文本
	for _, child := range n.Children {
		if child.XMLName.Local == "" && child.Content != "" {
			if text != "" {
				text += " "
			}
			text += strings.TrimSpace(child.Content)
		}
	}

	textElem := NewText(text)

	// 解析精细样式属性
	if fontWeight := props.Get("font-weight"); fontWeight == FontWeightBold {
		textElem = textElem.Bold()
	}

	if fontStyle := props.Get("font-style"); fontStyle == FontStyleItalic {
		textElem = textElem.Italic()
	}

	if decoration := props.Get("text-decoration"); decoration != "" {
		switch decoration {
		case DecorationUnderline:
			textElem = textElem.Underline()
		case DecorationStrikethrough:
			textElem = textElem.Strikethrough()
		}
	}

	if fgColor := props.Get("foreground-color"); fgColor != "" {
		if c, err := parseColor(fgColor); err == nil {
			textElem = textElem.ForegroundColor(c)
		}
	}

	if bgColor := props.Get("background-color"); bgColor != "" {
		if c, err := parseColor(bgColor); err == nil {
			textElem = textElem.BackgroundColor(c)
		}
	}

	if wrap := parseBoolAttr(props, "wrap", false); wrap {
		textElem = textElem.Wrap(true)
	}

	if alignment := props.Get("alignment"); alignment != "" {
		textElem = textElem.Alignment(alignment)
	}

	return textElem
}

// toBox 将节点转换为 Box 元素。
func (n *node) toBox(props Props) Element {
	var child Element
	children := n.childElements()
	if len(children) > 0 {
		child = children[0]
	}

	// 如果存在边框颜色，则解析
	borderColorStr := props.Get("border-color")
	var borderColor color.Color
	if borderColorStr != "" {
		if c, err := parseColor(borderColorStr); err == nil {
			borderColor = c
		}
	}

	// 解析宽度和高度约束
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))

	// 解析内边距
	padding := parseIntAttr(props, "padding", 0)

	// 解析外边距
	margin := parseIntAttr(props, "margin", 0)
	marginTop := parseIntAttr(props, "margin-top", 0)
	marginRight := parseIntAttr(props, "margin-right", 0)
	marginBottom := parseIntAttr(props, "margin-bottom", 0)
	marginLeft := parseIntAttr(props, "margin-left", 0)

	box := NewBox(child)
	if border := props.Get("border"); border != "" {
		box = box.Border(border)
	}
	if borderColor != nil {
		box = box.BorderColor(borderColor)
	}
	if !width.IsAuto() {
		box = box.Width(width)
	}
	if !height.IsAuto() {
		box = box.Height(height)
	}
	if padding > 0 {
		box = box.Padding(padding)
	}
	if margin > 0 {
		box = box.Margin(margin)
	}
	if marginTop > 0 {
		box = box.MarginTop(marginTop)
	}
	if marginRight > 0 {
		box = box.MarginRight(marginRight)
	}
	if marginBottom > 0 {
		box = box.MarginBottom(marginBottom)
	}
	if marginLeft > 0 {
		box = box.MarginLeft(marginLeft)
	}

	return box
}

// toSpacer 将节点转换为 Spacer 元素。
func (n *node) toSpacer(props Props) Element {
	size := parseIntAttr(props, "size", 0)
	if size > 0 {
		return NewFixedSpacer(size)
	}
	return NewSpacer()
}

// toFlex 将节点转换为 Flex 元素。
func (n *node) toFlex(props Props) Element {
	var child Element
	children := n.childElements()
	if len(children) > 0 {
		child = children[0]
	}

	grow := parseIntAttr(props, "grow", 0)
	shrink := parseIntAttr(props, "shrink", 1)
	basis := parseIntAttr(props, "basis", 0)

	flex := NewFlex(child)
	if grow > 0 {
		flex = flex.Grow(grow)
	}
	if shrink != 1 {
		flex = flex.Shrink(shrink)
	}
	if basis > 0 {
		flex = flex.Basis(basis)
	}

	return flex
}

// toPositioned 将节点转换为 Positioned 元素。
func (n *node) toPositioned(props Props) Element {
	var child Element
	children := n.childElements()
	if len(children) > 0 {
		child = children[0]
	}

	x := parseIntAttr(props, "x", 0)
	y := parseIntAttr(props, "y", 0)
	right := parseIntAttr(props, "right", -1)
	bottom := parseIntAttr(props, "bottom", -1)
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))

	positioned := NewPositioned(child, x, y)
	if right >= 0 {
		positioned = positioned.Right(right)
	}
	if bottom >= 0 {
		positioned = positioned.Bottom(bottom)
	}
	if !width.IsAuto() {
		positioned = positioned.Width(width)
	}
	if !height.IsAuto() {
		positioned = positioned.Height(height)
	}

	return positioned
}

// toDivider 将节点转换为 Divider 元素。
func (n *node) toDivider(props Props) Element {
	vertical := parseBoolAttr(props, "vertical", false)
	char := props.Get("char")

	var divider *Divider
	if vertical {
		divider = NewVerticalDivider()
	} else {
		divider = NewDivider()
	}

	// 解析前景色
	if fgColor := props.Get("foreground-color"); fgColor != "" {
		if c, err := parseColor(fgColor); err == nil {
			divider = divider.ForegroundColor(c)
		}
	}

	if char != "" {
		divider = divider.Char(char)
	}

	return divider
}

// toSlot 将节点转换为 Slot 元素。
func (n *node) toSlot(props Props) Element {
	name := props.Get("name")
	if name == "" {
		// Slot 需要一个名称
		name = "unnamed"
	}

	return NewSlot(name)
}

// toScrollView 将节点转换为 ScrollView 元素。
func (n *node) toScrollView(props Props) Element {
	var child Element
	children := n.childElements()
	if len(children) > 0 {
		child = children[0]
	}

	// 解析尺寸
	width := parseSizeConstraint(props.Get("width"))
	height := parseSizeConstraint(props.Get("height"))

	// 解析滚动选项
	offsetX := parseIntAttr(props, "offset-x", 0)
	offsetY := parseIntAttr(props, "offset-y", 0)
	showScrollbar := parseBoolAttr(props, "scrollbar", true)
	vertical := parseBoolAttr(props, "vertical", true)
	horizontal := parseBoolAttr(props, "horizontal", false)

	// 解析滚动条颜色
	var scrollbarColor color.Color
	if colorStr := props.Get("scrollbar-color"); colorStr != "" {
		if c, err := parseColor(colorStr); err == nil {
			scrollbarColor = c
		}
	}

	scrollView := NewScrollView(child)
	if offsetX != 0 || offsetY != 0 {
		scrollView = scrollView.Offset(offsetX, offsetY)
	}
	if !width.IsAuto() {
		scrollView = scrollView.Width(width)
	}
	if !height.IsAuto() {
		scrollView = scrollView.Height(height)
	}
	if !showScrollbar {
		scrollView = scrollView.Scrollbar(false)
	}
	if !vertical {
		scrollView = scrollView.Vertical(false)
	}
	if horizontal {
		scrollView = scrollView.Horizontal(true)
	}
	if scrollbarColor != nil {
		scrollView = scrollView.ScrollbarColor(scrollbarColor)
	}

	return scrollView
}

// 解析属性的辅助函数

func parseIntAttr(props Props, key string, defaultValue int) int {
	if val := props.Get(key); val != "" {
		var i int
		if _, err := fmt.Sscanf(val, "%d", &i); err == nil {
			return i
		}
	}
	return defaultValue
}

func parseBoolAttr(props Props, key string, defaultValue bool) bool {
	val := props.Get(key)
	if val == "" {
		return defaultValue
	}
	return val == "true" || val == "1" || val == "yes"
}
