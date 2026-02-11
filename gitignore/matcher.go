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
// 已被修改为提供一个简化的匹配器，用于处理单个模式。

package gitignore

// Matcher 定义了一个用于单个 gitignore 模式的匹配器
type Matcher struct {
	pattern Pattern
}

// NewMatcher 为单个模式构造一个新的简单匹配器
func NewMatcher(pattern Pattern) *Matcher {
	return &Matcher{pattern: pattern}
}

// Match 将给定路径与单个模式进行匹配
func (m *Matcher) Match(path []string, isDir bool) bool {
	match := m.pattern.Match(path, isDir)
	return match == Exclude
}
