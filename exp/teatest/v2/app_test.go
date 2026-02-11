package teatest_test

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	tea "github.com/purpose168/bubbletea-cn/v2"
	"github.com/purpose168/charm-experimental-packages-cn/exp/teatest/v2"
)

// TestApp 测试基本的应用程序测试功能
func TestApp(t *testing.T) {
	m := model(10)
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(70, 30),
	)
	t.Cleanup(func() {
		if err := tm.Quit(); err != nil {
			t.Fatal(err)
		}
	})

	time.Sleep(time.Second + time.Millisecond*200)
	tm.Type("我在输入内容，但我的程序会忽略它")
	tm.Send("被忽略的消息")
	tm.Send(tea.KeyPressMsg{
		Code: tea.KeyEnter,
	})

	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}

	out := readBts(t, tm.FinalOutput(t, teatest.WithFinalTimeout(time.Second)))
	if !regexp.MustCompile(`This program will exit in \d+ seconds`).Match(out) {
		t.Fatalf("输出与给定的正则表达式不匹配: %s", string(out))
	}
	teatest.RequireEqualOutput(t, out)

	if tm.FinalModel(t).(model) != 9 {
		t.Errorf("期望模型为 10, 实际为 %d", m)
	}
}

// TestAppInteractive 测试交互式应用程序测试功能
func TestAppInteractive(t *testing.T) {
	m := model(10)
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(70, 30),
	)

	time.Sleep(time.Second + time.Millisecond*200)
	tm.Send("被忽略的消息")

	if bts := readBts(t, tm.Output()); !bytes.Contains(bts, []byte("9 seconds")) {
		t.Fatalf("输出不匹配: 期望 %q", string(bts))
	}

	teatest.WaitFor(t, tm.Output(), func(out []byte) bool {
		return bytes.Contains(out, []byte("7"))
	}, teatest.WithDuration(5*time.Second), teatest.WithCheckInterval(time.Millisecond*10))

	tm.Send(tea.KeyPressMsg{
		Code: tea.KeyEnter,
	})

	if err := tm.Quit(); err != nil {
		t.Fatal(err)
	}

	if tm.FinalModel(t).(model) != 7 {
		t.Errorf("期望模型为 7, 实际为 %d", m)
	}
}

// readBts 从 io.Reader 读取所有字节
func readBts(tb testing.TB, r io.Reader) []byte {
	tb.Helper()
	bts, err := io.ReadAll(r)
	if err != nil {
		tb.Fatal(err)
	}
	return bts
}

// model 可以是任何类型的数据。它保存程序的所有数据，
// 所以通常它是一个结构体。但在这个简单的例子中，
// 我们只需要一个简单的整数。
type model int

// Init 可选地返回我们应该运行的初始命令。在这种情况下，
// 我们想启动计时器。
func (m model) Init() tea.Cmd {
	return tick
}

// Update 在收到消息时被调用。其思想是检查消息并相应地返回更新后的模型。
// 你还可以返回一个命令，这是一个执行 I/O 并返回消息的函数。
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick
	}
	return m, nil
}

// View 根据模型中的数据返回一个字符串。该字符串将被渲染到终端。
func (m model) View() string {
	return fmt.Sprintf("Hi. This program will exit in %d seconds. To quit sooner press any key.\n", m)
}

// 消息是我们在 Update 函数中响应的事件。
// 这个特定的消息表示计时器已经滴答作响。
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
