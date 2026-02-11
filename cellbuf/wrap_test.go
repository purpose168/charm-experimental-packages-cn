package cellbuf

import (
	"fmt"
	"testing"
)

// wrapCases 包含了各种文本换行测试用例
var wrapCases = []struct {
	name     string // 测试用例名称
	input    string // 输入文本
	expected string // 期望的输出文本
	width    int    // 换行宽度
}{
	{
		name:     "简单文本",
		input:    "I really \x1B[38;2;249;38;114mlove the\x1B[0m Go language!",
		expected: "I really \x1B[38;2;249;38;114mlove\x1b[m\n\x1B[38;2;249;38;114mthe\x1B[0m Go\nlanguage!",
		width:    14,
	},
	{
		name:     "直接通过",
		input:    "hello world",
		expected: "hello world",
		width:    11,
	},
	{
		name:     "亚洲语言",
		input:    "こんにち",
		expected: "こんに\nち",
		width:    7,
	},
	{
		name:     "表情符号",
		input:    "😃👰🏻‍♀️🫧",
		expected: "😃\n👰🏻‍♀️\n🫧",
		width:    2,
	},
	{
		name:     "长样式文本",
		input:    "\x1B[38;2;249;38;114ma really long string\x1B[0m",
		expected: "\x1B[38;2;249;38;114ma really\x1b[m\n\x1B[38;2;249;38;114mlong\x1b[m\n\x1B[38;2;249;38;114mstring\x1B[0m",
		width:    10,
	},
	{
		name:     "长样式文本带非断空格",
		input:    "\x1B[38;2;249;38;114ma really\u00a0long string\x1B[0m",
		expected: "\x1b[38;2;249;38;114ma\x1b[m\n\x1b[38;2;249;38;114mreally\u00a0lon\x1b[m\n\x1b[38;2;249;38;114mg string\x1b[0m",
		width:    10,
	},
	{
		name:     "更长的文本",
		input:    "the quick brown foxxxxxxxxxxxxxxxx jumped over the lazy dog.",
		expected: "the quick brown\nfoxxxxxxxxxxxxxx\nxx jumped over\nthe lazy dog.",
		width:    16,
	},
	{
		name:     "更长的亚洲文本",
		input:    "猴 猴 猴猴 猴猴猴猴猴猴猴猴猴 猴猴猴 猴猴 猴’ 猴猴 猴.",
		expected: "猴 猴 猴猴\n猴猴猴猴猴猴猴猴\n猴 猴猴猴 猴猴\n猴’ 猴猴 猴.",
		width:    16,
	},
	{
		name:     "长输入文本",
		input:    "Rotated keys for a-good-offensive-cheat-code-incorporated/animal-like-law-on-the-rocks.",
		expected: "Rotated keys for a-good-offensive-cheat-code-incorporated/animal-like-law-\non-the-rocks.",
		width:    76,
	},
	{
		name:     "长输入文本2",
		input:    "Rotated keys for a-good-offensive-cheat-code-incorporated/crypto-line-operating-system.",
		expected: "Rotated keys for a-good-offensive-cheat-code-incorporated/crypto-line-\noperating-system.",
		width:    76,
	},
	{
		name:     "连字符断点",
		input:    "a-good-offensive-cheat-code",
		expected: "a-good-\noffensive-\ncheat-code",
		width:    10,
	},
	{
		name:     "精确宽度",
		input:    "\x1b[91mfoo\x1b[0m",
		expected: "\x1b[91mfoo\x1b[0m",
		width:    3,
	},
	{
		// XXX: 我们是否应该在文本换行时保留空格？
		name:     "额外空格",
		input:    "foo ",
		expected: "foo",
		width:    3,
	},
	{
		name:     "带样式的额外空格",
		input:    "\x1b[mfoo \x1b[m",
		expected: "\x1b[mfoo\x1b[m",
		width:    3,
	},
	{
		name:     "带样式的段落",
		input:    "Lorem ipsum dolor \x1b[1msit\x1b[m amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. \x1b[31mUt enim\x1b[m ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea \x1b[38;5;200mcommodo consequat\x1b[m. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. \x1b[1;2;33mExcepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\x1b[m",
		expected: "Lorem ipsum dolor \x1b[1msit\x1b[m amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor\nincididunt ut labore et dolore\nmagna aliqua. \x1b[31mUt enim\x1b[m ad minim\nveniam, quis nostrud\nexercitation ullamco laboris\nnisi ut aliquip ex ea \x1b[38;5;200mcommodo\x1b[m\n\x1b[38;5;200mconsequat\x1b[m. Duis aute irure\ndolor in reprehenderit in\nvoluptate velit esse cillum\ndolore eu fugiat nulla\npariatur. \x1b[1;2;33mExcepteur sint\x1b[m\n\x1b[1;2;33moccaecat cupidatat non\x1b[m\n\x1b[1;2;33mproident, sunt in culpa qui\x1b[m\n\x1b[1;2;33mofficia deserunt mollit anim\x1b[m\n\x1b[1;2;33mid est laborum.\x1b[m",
		width:    30,
	},
	{"连字符换行", "foo-bar", "foo-\nbar", 5},
	{"双空格", "f  bar foobaz", "f  bar\nfoobaz", 6},
	{"直接通过", "foobar\n ", "foobar\n ", 0},
	{"通过", "foo", "foo", 3},
	{"过长文本", "foobarfoo", "foob\narfo\no", 4},
	{"空白字符", "foo bar foo", "foo\nbar\nfoo", 4},
	{"按空格分割", "foo bars foobars", "foo\nbars\nfoob\nars", 4},
	{"连字符", "foob-foobar", "foob\n-foo\nbar", 4},
	{"宽表情符号断点", "foo🫧 foobar", "foo\n🫧\nfoob\nar", 4},
	{"空格断点", "foo --bar", "foo --bar", 9},
	{"简单", "foo bars foobars", "foo\nbars\nfoob\nars", 4},
	{"限制宽度", "foo bar", "foo\nbar", 5},
	{"移除空白字符", "foo    \nb   ar   ", "foo\nb\nar", 4},
	{"空白字符尾部宽度", "foo\nb\ta\n bar", "foo\nb\ta\n bar", 4},
	{"显式换行", "foo bar foo\n", "foo\nbar\nfoo\n", 4},
	{"多个显式换行", "\nfoo bar\n\n\nfoo\n", "\nfoo\nbar\n\n\nfoo\n", 4},
	{"示例", " This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* foo  \nbar    ", " This\nis a\nlist: \n\n\t* foo\n\t* bar\n\n\n\t* foo\nbar", 6},
	{"样式代码不影响长度", "\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m", "\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m", 7},
	{"样式代码不被换行", "\x1B[38;2;249;38;114m(\x1B[0m\x1B[38;2;248;248;242mjust another test\x1B[38;2;249;38;114m)\x1B[0m", "\x1b[38;2;249;38;114m(\x1b[0m\x1b[38;2;248;248;242mjust\x1b[m\n\x1b[38;2;248;248;242manother\x1b[m\n\x1b[38;2;248;248;242mtest\x1b[38;2;249;38;114m)\x1b[0m", 7},
	{"OSC8 链接包装", "สวัสดีสวัสดี\x1b]8;;https://example.com\x1b\\ สวัสดีสวัสดี\x1b]8;;\x1b\\", "สวัสดีสวัสดี\x1b]8;;https://example.com\x1b\\\x1b]8;;\x07\n\x1b]8;;https://example.com\x07สวัสดีสวัสดี\x1b]8;;\x1b\\", 8},
	{"制表符", "foo\tbar", "foo\nbar", 3},
	{"包装样式示例", "", "", 10},
	{
		name:     "带格式的单词后带空格和标点",
		input:    "\x1b[38;5;203;48;5;236m arm64 \x1b[0m, \x1b[38;5;203;48;5;236m amd64 \x1b[0m, \x1b[38;5;203;48;5;236m i386 \x1b[0m",
		expected: "\x1b[38;5;203;48;5;236m arm64 \x1b[0m,\n\x1b[38;5;203;48;5;236m amd64 \x1b[0m, \x1b[38;5;203;48;5;236m i386 \x1b[0m",
		width:    15,
	},
}

// TestWrap 测试 Wrap 函数的换行功能
func TestWrap(t *testing.T) {
	for i, tc := range wrapCases {
		t.Run(tc.name, func(t *testing.T) {
			output := Wrap(tc.input, tc.width, "")
			if output != tc.expected {
				t.Errorf("测试用例 %d, 输入:\n%q\n期望输出:\n%q\n%s\n\n实际输出:\n%q\n%s", i+1, tc.input, tc.expected, tc.expected, output, output)
			}
		})
	}
}

// ExampleWrap 展示 Wrap 函数的使用示例
func ExampleWrap() {
	fmt.Println(Wrap("The quick brown fox jumped over the lazy dog.", 20, ""))
	// 输出:
	// The quick brown fox
	// jumped over the lazy
	// dog.
}
