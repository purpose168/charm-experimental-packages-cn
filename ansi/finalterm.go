package ansi

import "strings"

// FinalTerm 返回一个用于 shell 集成的转义序列。
// 最初由 FinalTerm 设计了该协议，因此得名。
//
//	OSC 133 ; Ps ; Pm ST
//	OSC 133 ; Ps ; Pm BEL
//
// 请参阅：https://iterm2.com/documentation-shell-integration.html
func FinalTerm(pm ...string) string {
	return "\x1b]133;" + strings.Join(pm, ";") + "\x07"
}

// FinalTermPrompt 返回一个用于 shell 集成提示符标记的转义序列。
// 它在 shell 提示符开始之前发送。
//
// 这是 FinalTerm("A") 的别名。
func FinalTermPrompt(pm ...string) string {
	if len(pm) == 0 {
		return FinalTerm("A")
	}
	return FinalTerm(append([]string{"A"}, pm...)...)
}

// FinalTermCmdStart 返回一个用于 shell 集成命令开始标记的转义序列。
// 它在 shell 提示符结束之后、用户输入命令之前发送。
//
// 这是 FinalTerm("B") 的别名。
func FinalTermCmdStart(pm ...string) string {
	if len(pm) == 0 {
		return FinalTerm("B")
	}
	return FinalTerm(append([]string{"B"}, pm...)...)
}

// FinalTermCmdExecuted 返回一个用于 shell 集成命令执行标记的转义序列。
// 它在命令输出开始之前发送。
//
// 这是 FinalTerm("C") 的别名。
func FinalTermCmdExecuted(pm ...string) string {
	if len(pm) == 0 {
		return FinalTerm("C")
	}
	return FinalTerm(append([]string{"C"}, pm...)...)
}

// FinalTermCmdFinished 返回一个用于 shell 集成命令完成标记的转义序列。
//
// 如果在 [FinalTermCmdStart] 之后发送，表示命令已中止。
// 如果在 [FinalTermCmdExecuted] 之后发送，表示命令输出结束。
// 如果之前都没有发送，[FinalTermCmdFinished] 应该被忽略。
//
// 这是 FinalTerm("D") 的别名。
func FinalTermCmdFinished(pm ...string) string {
	if len(pm) == 0 {
		return FinalTerm("D")
	}
	return FinalTerm(append([]string{"D"}, pm...)...)
}
