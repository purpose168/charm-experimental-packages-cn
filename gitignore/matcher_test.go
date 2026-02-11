// Copyright 2019 The go-git authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// 此文件最初是 https://github.com/go-git/go-git/blob/main/plumbing/format/gitattributes 的一部分
// 已被修改为提供针对单一模式工作的简化匹配器的测试。
//
// 此文件最初是 https://github.com/go-git/go-git/blob/main/plumbing/format/gitattributes 的一部分
// 已被修改为无依赖项。

package gitignore

import (
	"testing"
)

func TestMatcher_Match(t *testing.T) {
	pattern := ParsePattern("**/middle/v[uo]l?ano", nil)
	matcher := NewMatcher(pattern)

	if !matcher.Match([]string{"head", "middle", "vulkano"}, false) {
		t.Error("期望匹配 'head/middle/vulkano'")
	}

	if matcher.Match([]string{"head", "middle", "other"}, false) {
		t.Error("期望不匹配 'head/middle/other'")
	}
}

func TestMatcher_Exclude(t *testing.T) {
	pattern := ParsePattern("!volcano", nil)
	matcher := NewMatcher(pattern)

	// 包含模式返回 false
	if matcher.Match([]string{"volcano"}, false) {
		t.Error("包含模式应该返回 false")
	}
}

// 测试展示如何使用 Matcher 处理排除模式
func TestMatcher_ExcludeHandling(t *testing.T) {
	// 对于排除模式，Matcher 将返回 false
	// 因为排除意味着 "不排除"，实际上就是 "包含"
	excludePattern := ParsePattern("!volcano", nil)
	matcher := NewMatcher(excludePattern)

	// 这返回 false 因为它是一个包含模式
	if matcher.Match([]string{"volcano"}, false) {
		t.Error("包含模式应该返回 false")
	}
}

// 测试来自 git 文档的 "排除所有除了..." 示例
// 注意：这是一个简化版本，测试各个单独的模式
func TestMatcher_EverythingExceptExample(t *testing.T) {
	// 测试 /* 模式（排除所有）
	pattern1 := ParsePattern("/*", nil)
	matcher1 := NewMatcher(pattern1)

	if !matcher1.Match([]string{"foo"}, true) { // 应该匹配并排除
		t.Error("期望匹配 'foo'")
	}

	if !matcher1.Match([]string{"baz"}, false) { // 应该匹配并排除
		t.Error("期望匹配 'baz'")
	}

	// 测试 !/foo 模式（但不排除 foo）
	pattern2 := ParsePattern("!/foo", nil)
	matcher2 := NewMatcher(pattern2)

	if matcher2.Match([]string{"foo"}, true) { // 应该匹配但包含（不排除）
		t.Error("包含模式应该返回 false")
	}

	// 测试 /foo/* 模式（排除 foo 目录中的文件）
	pattern3 := ParsePattern("/foo/*", nil)
	matcher3 := NewMatcher(pattern3)

	if !matcher3.Match([]string{"foo", "bar"}, false) { // 应该匹配并排除
		t.Error("期望匹配 'foo/bar'")
	}

	// 测试 !/foo/bar 模式（但不排除 foo/bar）
	pattern4 := ParsePattern("!/foo/bar", nil)
	matcher4 := NewMatcher(pattern4)

	if matcher4.Match([]string{"foo", "bar"}, false) { // 应该匹配但包含（不排除）
		t.Error("包含模式应该返回 false")
	}
}
