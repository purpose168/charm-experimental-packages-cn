package ansi

import (
	"fmt"
	"strings"
)

// URxvtExt 返回一个用于调用URxvt Perl扩展的转义序列
// 参数包括扩展名称和参数
//
// 格式：
//
//	OSC 777 ; extension_name ; param1 ; param2 ; ... ST
//	OSC 777 ; extension_name ; param1 ; param2 ; ... BEL
//
// 请参阅：https://man.archlinux.org/man/extra/rxvt-unicode/urxvt.7.en#XTerm_Operating_System_Commands
func URxvtExt(extension string, params ...string) string {
	// 格式化为URxvt扩展调用序列
	return fmt.Sprintf("\x1b]777;%s;%s\x07", extension, strings.Join(params, ";"))
}
