package vcr

import (
	"strings"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// headersToKeep 定义了需要保留的HTTP头部。
var headersToKeep = map[string]struct{}{
	"accept":       {},       // 接受的内容类型
	"content-type": {},       // 内容类型
	"user-agent":   {},       // 用户代理
}

// hookRemoveHeaders 创建一个钩子函数，用于移除HTTP请求和响应中的非必要头部信息。
// keepAll 参数为true时，保留所有头部信息。
func hookRemoveHeaders(keepAll bool) recorder.HookFunc {
	return func(i *cassette.Interaction) error {
		if keepAll {
			return nil // 如果keepAll为true，保留所有头部
		}
		// 移除请求中的非必要头部
		for k := range i.Request.Headers {
			if _, ok := headersToKeep[strings.ToLower(k)]; !ok {
				delete(i.Request.Headers, k)
			}
		}
		// 移除响应中的非必要头部
		for k := range i.Response.Headers {
			if _, ok := headersToKeep[strings.ToLower(k)]; !ok {
				delete(i.Response.Headers, k)
			}
		}
		return nil
	}
}
