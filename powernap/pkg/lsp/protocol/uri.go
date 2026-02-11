// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

// 此文件声明了 URI、DocumentUri 及其方法。
//
// 有关这些类型的 LSP 定义，请参阅
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#uri

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

// DocumentURI 是客户端编辑器文档的 URI。
//
// 根据 LSP 规范：
//
//	应注意处理 URI 中的编码。例如，
//	一些客户端（如 VS Code）可能会对驱动器号中的冒号进行编码，
//	而其他客户端则不会。下面的 URI 都是有效的，
//	但客户端和服务器应在自己使用的形式上保持一致，
//	以确保另一方不会将它们解释为不同的 URI。
//	客户端和服务器不应假设对方以相同的方式编码
//	（例如，编码驱动器号中冒号的客户端不能假设
//	服务器响应会有编码的冒号）。同样适用于驱动器号的大小写 - 
//	一方不应假设另一方会返回与自己大小写相同的驱动器号路径。
//
//	file:///c:/project/readme.md
//	file:///C%3A/project/readme.md
//
// 这是在 JSON 反序列化期间完成的；
// 有关详细信息，请参阅 [DocumentURI.UnmarshalText]。
type DocumentURI string

// URI 是任意 URL（例如 https），不一定是文件。
type URI = string

// UnmarshalText 实现 DocumentUri 值的解码。
//
// 特别是，它实现了对 LSP 规范中 DocumentUri 定义的各种奇怪特性的系统纠正，
// 这些特性似乎是为了解决 VS Code 中的错误。例如，它可能会对 URI 本身进行 URI 编码，
// 使冒号变为 %3A，并且可能会发送只有两个斜杠（不是三个）且没有主机名的 file://foo.go URI。
//
// 我们使用 UnmarshalText 而不是 UnmarshalJSON，因为即使对于不可寻址的值，
// 如 map[K]V 的键和值，也会调用它，而这些值没有可调用 UnmarshalJSON 的 *K 或 *V 类型指针。
// （有关更多详细信息，请参阅 Go 问题 #28189。）
//
// 非空的 DocumentUris 是有效的 "file" 方案 URI。
// 空的 DocumentUri 是有效的。
func (uri *DocumentURI) UnmarshalText(data []byte) (err error) {
	*uri, err = ParseDocumentURI(string(data))
	return err
}

// Path 返回给定 URI 的文件路径。
//
// DocumentUri("").Path() 返回空字符串。
//
// 如果对不是有效文件名的 URI 调用 Path，会导致恐慌。
func (uri DocumentURI) Path() (string, error) {
	filename, err := filename(uri)
	if err != nil {
		// 例如，ParseRequestURI 失败。
		//
		// 这只会影响通过直接字符串操作创建的 DocumentUris；
		// 从客户端接收的所有 DocumentUris 都经过 ParseRequestURI，
		// 这确保了有效性。
		return "", fmt.Errorf("invalid URI %q: %w", uri, err)
	}
	return filepath.FromSlash(filename), nil
}

func filename(uri DocumentURI) (string, error) {
	if uri == "" {
		return "", nil
	}

	// 这种对简单非空绝对 POSIX 文件名常见情况的保守检查
	// 避免了分配 net.URL。
	if strings.HasPrefix(string(uri), "file:///") {
		rest := string(uri)[len("file://"):] // leave one slash
		for i := range len(rest) {
			b := rest[i]
			// 拒绝这些情况：
			if b < ' ' || b == 0x7f || // 控制字符
				b == '%' || b == '+' || // URI 转义
				b == ':' || // Windows 驱动器号
				b == '@' || b == '&' || b == '?' { // 权限或查询
				goto slow
			}
		}
		return rest, nil
	}
slow:

	u, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return "", fmt.Errorf("parsing URI %q: %w", uri, err)
	}
	if u.Scheme != fileScheme {
		return "", fmt.Errorf("only file URIs are supported, got %q from %q", u.Scheme, uri)
	}
	// If the URI is a Windows URI, we trim the leading "/" and uppercase
	// the drive letter, which will never be case sensitive.
	if isWindowsDrivePath(u.Path) {
		u.Path = strings.ToUpper(string(u.Path[1])) + u.Path[2:]
	}

	return u.Path, nil
}

// ParseDocumentURI 将字符串解释为 DocumentUri，应用 VS Code 变通方法；
// 有关详细信息，请参阅 [DocumentURI.UnmarshalText]。
func ParseDocumentURI(s string) (DocumentURI, error) {
	if s == "" {
		return "", nil
	}

	if !strings.HasPrefix(s, "file://") {
		return "", fmt.Errorf("DocumentUri scheme is not 'file': %s", s)
	}

	// VS Code 发送只有两个斜杠的 URL，这是无效的。golang/go#39789。
	if !strings.HasPrefix(s, "file:///") {
		s = "file:///" + s[len("file://"):]
	}

	// 尽管输入是 URI，但它可能不是规范形式。特别是 VS Code 会过度转义 :、@ 等字符。
	// 解转义并重新编码以规范化。
	path, err := url.PathUnescape(s[len("file://"):])
	if err != nil {
		return "", fmt.Errorf("unescaping URI path %q: %w", s, err)
	}

	// 来自 Windows 的文件 URI 可能有小写的驱动器号。
	// 由于驱动器号保证不区分大小写，我们将它们更改为大写以保持一致性。
	// 例如，file:///c:/x/y/z 变为 file:///C:/x/y/z。
	if isWindowsDrivePath(path) {
		path = path[:1] + strings.ToUpper(string(path[1])) + path[2:]
	}
	u := url.URL{Scheme: fileScheme, Path: path}
	return DocumentURI(u.String()), nil
}

// URIFromPath 为提供的文件路径返回 DocumentUri。
// 给定 ""，它返回 ""。
func URIFromPath(path string) DocumentURI {
	if path == "" {
		return ""
	}
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	if isWindowsDrivePath(path) {
		path = "/" + strings.ToUpper(string(path[0])) + path[1:]
	}
	path = filepath.ToSlash(path)
	filepath.Clean(path)
	u := url.URL{
		Scheme: fileScheme,
		Path:   path,
	}
	return DocumentURI(u.String())
}

const fileScheme = "file"

// isWindowsDrivePath 如果文件路径采用 Windows 使用的形式，则返回 true。
// 我们检查路径是否以驱动器号开头，后跟 ":"。
// 例如：C:/x/y/z。
func isWindowsDrivePath(path string) bool {
	return filepath.VolumeName(path) != ""
}
