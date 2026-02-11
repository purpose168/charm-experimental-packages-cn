package teatest

import (
	"fmt"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	tea "github.com/purpose168/bubbletea-cn/v2"
)

// TestWaitForErrorReader 测试 WaitFor 函数在读取器返回错误时的行为
func TestWaitForErrorReader(t *testing.T) {
	err := doWaitFor(iotest.ErrReader(fmt.Errorf("fake")), func(bts []byte) bool {
		return true
	}, WithDuration(time.Millisecond), WithCheckInterval(10*time.Microsecond))
	if err == nil {
		t.Fatal("期望一个错误，得到 nil")
	}
	if err.Error() != "WaitFor: fake" {
		t.Fatalf("意外的错误: %s", err.Error())
	}
}

// TestWaitForTimeout 测试 WaitFor 函数在超时时的行为
func TestWaitForTimeout(t *testing.T) {
	err := doWaitFor(strings.NewReader("nope"), func(bts []byte) bool {
		return false
	}, WithDuration(time.Millisecond), WithCheckInterval(10*time.Microsecond))
	if err == nil {
		t.Fatal("期望一个错误，得到 nil")
	}
	if err.Error() != "WaitFor: condition not met after 1ms. Last output:\nnope" {
		t.Fatalf("意外的错误: %s", err.Error())
	}
}

// m 是一个简单的 tea.Model 实现
type m string

func (m m) Init() tea.Cmd                       { return nil }
func (m m) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m m) View() string                        { return string(m) }

// TestWaitFinishedWithTimeoutFn 测试 WaitFinished 函数带有超时函数时的行为
func TestWaitFinishedWithTimeoutFn(t *testing.T) {
	tm := NewTestModel(t, m("a"))
	var timedOut bool
	tm.WaitFinished(t, WithFinalTimeout(time.Nanosecond), WithTimeoutFn(func(testing.TB) {
		timedOut = true
	}))
	if !timedOut {
		t.Fatal("期望 timedOut 被设置")
	}
}
