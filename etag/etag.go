// Package etag 提供用于在 HTTP 请求和响应中生成和处理 ETag 头部的工具。
package etag

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
)

// Of 返回给定数据的 etag。
func Of(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf(`%x`, hash[:16])
}

// Request 在给定的请求中设置 `If-None-Match` 头部，适当地引用 etag 值。
func Request(req *http.Request, etag string) {
	if etag == "" {
		return
	}
	req.Header.Add("If-None-Match", fmt.Sprintf(`"%s"`, etag))
}

// Response 在给定的响应写入器中设置 `ETag` 头部，适当地引用 etag 值。
func Response(w http.ResponseWriter, etag string) {
	if etag == "" {
		return
	}
	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, etag))
}

// Matches 检查给定的请求是否有与给定 etag 匹配的 `If-None-Match` 头部。
func Matches(r *http.Request, etag string) bool {
	header := r.Header.Get("If-None-Match")
	if header == "" || etag == "" {
		return false
	}
	return unquote(header) == unquote(etag)
}

// unquote 移除字符串两端的引号。
func unquote(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
}
