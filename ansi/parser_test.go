package ansi

import (
	"os"
	"reflect"
	"slices"
	"testing"
)

// csiSequence 表示 CSI 序列

type csiSequence struct {
	Cmd    Cmd    // 命令
	Params Params // 参数
}

// dcsSequence 表示 DCS 序列

type dcsSequence struct {
	Cmd    Cmd    // 命令
	Params Params // 参数
	Data   []byte // 数据
}

// testCase 表示测试用例

type testCase struct {
	name     string // 测试用例名称
	input    string // 输入字符串
	expected []any  // 期望结果
}

// testDispatcher 表示测试调度器

type testDispatcher struct {
	dispatched []any // 已分发的序列
}

// dispatchRune 分发符文

func (d *testDispatcher) dispatchRune(r rune) {
	d.dispatched = append(d.dispatched, r)
}

// dispatchControl 分发控制字符

func (d *testDispatcher) dispatchControl(b byte) {
	d.dispatched = append(d.dispatched, b)
}

// dispatchEsc 分发 ESC 序列

func (d *testDispatcher) dispatchEsc(cmd Cmd) {
	d.dispatched = append(d.dispatched, cmd)
}

// dispatchCsi 分发 CSI 序列

func (d *testDispatcher) dispatchCsi(cmd Cmd, params Params) {
	params = slices.Clone(params)
	d.dispatched = append(d.dispatched, csiSequence{Cmd: cmd, Params: params})
}

// dispatchDcs 分发 DCS 序列

func (d *testDispatcher) dispatchDcs(cmd Cmd, params Params, data []byte) {
	params = slices.Clone(params)
	data = slices.Clone(data)
	d.dispatched = append(d.dispatched, dcsSequence{Cmd: cmd, Params: params, Data: data})
}

// dispatchOsc 分发 OSC 序列

func (d *testDispatcher) dispatchOsc(cmd int, data []byte) {
	data = slices.Clone(data)
	d.dispatched = append(d.dispatched, data)
}

// dispatchApc 分发 APC 序列

func (d *testDispatcher) dispatchApc(data []byte) {
	data = slices.Clone(data)
	d.dispatched = append(d.dispatched, data)
}

// testParser 创建一个测试用的解析器

func testParser(d *testDispatcher) *Parser {
	p := NewParser()
	p.SetHandler(Handler{
		Print:     d.dispatchRune,
		Execute:   d.dispatchControl,
		HandleEsc: d.dispatchEsc,
		HandleCsi: d.dispatchCsi,
		HandleDcs: d.dispatchDcs,
		HandleOsc: d.dispatchOsc,
		HandleApc: d.dispatchApc,
	})
	p.SetParamsSize(16)
	p.SetDataSize(0)
	return p
}

// TestControlSequence 测试控制序列的解析

func TestControlSequence(t *testing.T) {
	cases := []testCase{
		{
			name:     "仅ESC",
			input:    "\x1b",
			expected: []any{},
		},
		{
			name:  "双重ESC",
			input: "\x1b\x1b",
			expected: []any{
				byte(0x1b),
			},
		},
		// {
		// 	name:  "esc_bracket",
		// 	input: "\x1b[",
		// 	expected: []Sequence{
		// 		EscSequence('['),
		// 	},
		// },
		// {
		// 	name:  "csi_rune_esc_bracket",
		// 	input: "\x1b[1;2;3mabc\x1b\x1bP",
		// 	expected: []Sequence{
		// 		CsiSequence{
		// 			Params: []Parameter{1, 2, 3},
		// 			Cmd:    'm',
		// 		},
		// 		Rune('a'),
		// 		Rune('b'),
		// 		Rune('c'),
		// 		ControlCode(0x1b),
		// 		EscSequence('P'),
		// 	},
		// },
		{
			name:  "CSI加文本",
			input: "Hello, \x1b[31mWorld!\x1b[0m",
			expected: []any{
				rune('H'),
				rune('e'),
				rune('l'),
				rune('l'),
				rune('o'),
				rune(','),
				rune(' '),
				csiSequence{
					Params: Params{31},
					Cmd:    'm',
				},
				rune('W'),
				rune('o'),
				rune('r'),
				rune('l'),
				rune('d'),
				rune('!'),
				csiSequence{
					Params: Params{0},
					Cmd:    'm',
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dispatcher := &testDispatcher{}
			parser := testParser(dispatcher)
			parser.Parse([]byte(c.input))
			assertEqual(t, len(c.expected), len(dispatcher.dispatched))
			for i := range c.expected {
				assertEqual(t, c.expected[i], dispatcher.dispatched[i])
			}
		})
	}
}

// parsers 定义了不同配置的解析器

var parsers = []struct {
	name   string  // 解析器名称
	parser *Parser // 解析器实例
}{
	{
		name:   "simple",
		parser: &Parser{},
	},
	{
		name: "params",
		parser: func() *Parser {
			p := NewParser()
			p.SetDataSize(0)
			p.SetParamsSize(16)
			return p
		}(),
	},
	{
		name: "params and data",
		parser: func() *Parser {
			p := NewParser()
			p.SetDataSize(1024)
			p.SetParamsSize(16)
			return p
		}(),
	},
}

// BenchmarkParser 基准测试解析器性能

func BenchmarkParser(b *testing.B) {
	bts, err := os.ReadFile("./fixtures/demo.vte")
	if err != nil {
		b.Fatalf("错误: %v", err)
	}

	for _, p := range parsers {
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.parser.Parse(bts)
			}
		})
	}
}

// BenchmarkParserUTF8 基准测试解析器处理UTF-8的性能

func BenchmarkParserUTF8(b *testing.B) {
	bts, err := os.ReadFile("./fixtures/UTF-8-demo.txt")
	if err != nil {
		b.Fatalf("Error: %v", err)
	}

	for _, p := range parsers {
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.parser.Parse(bts)
			}
		})
	}
}

func BenchmarkParserStateChanges(b *testing.B) {
	input := []byte("\x1b]2;X\x1b\\こんにちは\x1b[0m \x1bP0@\x1b\\")

	for _, p := range parsers {
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.parser.Parse(input)
			}
		})
	}
}

func assertEqual[T any](t *testing.T, expected, got T) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected:\n  %#v, got:\n  %#v", expected, got)
	}
}
