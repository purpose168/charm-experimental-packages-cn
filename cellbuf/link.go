package cellbuf

import (
	"github.com/charmbracelet/colorprofile"
)

// ConvertLink 转换超链接以尊重给定的颜色配置文件。
func ConvertLink(h Link, p colorprofile.Profile) Link {
	if p == colorprofile.NoTTY {
		return Link{}
	}

	return h
}
