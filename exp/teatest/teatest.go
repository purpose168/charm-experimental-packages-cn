// Package teatest 提供测试 tea.Model 的辅助函数。
package teatest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	tea "github.com/purpose168/bubbletea-cn"
	"github.com/purpose168/charm-experimental-packages-cn/exp/golden"
)

// Program 定义了测试所需的 tea.Program API 子集。
type Program interface {
	Send(tea.Msg)
}

// TestModelOptions 定义了测试函数可用的所有选项。
type TestModelOptions struct {
	size tea.WindowSizeMsg
}

// TestOption 是一个函数式选项。
type TestOption func(opts *TestModelOptions)

// WithInitialTermSize 设置初始终端大小。
func WithInitialTermSize(x, y int) TestOption {
	return func(opts *TestModelOptions) {
		opts.size = tea.WindowSizeMsg{
			Width:  x,
			Height: y,
		}
	}
}

// WaitingForContext 是 WaitFor 的上下文。
type WaitingForContext struct {
	Duration      time.Duration
	CheckInterval time.Duration
}

// WaitForOption 更改 WaitFor 的行为。
type WaitForOption func(*WaitingForContext)

// WithCheckInterval 设置 WaitFor 在每次检查之间应睡眠的时间。
func WithCheckInterval(d time.Duration) WaitForOption {
	return func(wf *WaitingForContext) {
		wf.CheckInterval = d
	}
}

// WithDuration 设置 WaitFor 等待条件的时间。
func WithDuration(d time.Duration) WaitForOption {
	return func(wf *WaitingForContext) {
		wf.Duration = d
	}
}

// WaitFor 持续从 r 读取直到条件匹配。
// 默认持续时间为 1 秒，默认检查间隔为 50 毫秒。
// 这些默认值可以通过 WithDuration 和 WithCheckInterval 更改。
func WaitFor(
	tb testing.TB,
	r io.Reader,
	condition func(bts []byte) bool,
	options ...WaitForOption,
) {
	tb.Helper()
	if err := doWaitFor(r, condition, options...); err != nil {
		tb.Fatal(err)
	}
}

func doWaitFor(r io.Reader, condition func(bts []byte) bool, options ...WaitForOption) error {
	wf := WaitingForContext{
		Duration:      time.Second,
		CheckInterval: 50 * time.Millisecond, //nolint: mnd
	}

	for _, opt := range options {
		opt(&wf)
	}

	var b bytes.Buffer
	start := time.Now()
	for time.Since(start) <= wf.Duration {
		if _, err := io.ReadAll(io.TeeReader(r, &b)); err != nil {
			return fmt.Errorf("WaitFor: %w", err)
		}
		if condition(b.Bytes()) {
			return nil
		}
		time.Sleep(wf.CheckInterval)
	}
	return fmt.Errorf("WaitFor: 条件在 %s 后未满足。最后输出:\n%s", wf.Duration, b.String())
}

// TestModel 是正在被测试的模型。
type TestModel struct {
	program *tea.Program

	in  *bytes.Buffer
	out io.ReadWriter

	modelCh chan tea.Model
	model   tea.Model

	done   sync.Once
	doneCh chan bool
}

// NewTestModel 创建一个新的 TestModel，可用于测试。
func NewTestModel(tb testing.TB, m tea.Model, options ...TestOption) *TestModel {
	tm := &TestModel{
		in:      bytes.NewBuffer(nil),
		out:     safe(bytes.NewBuffer(nil)),
		modelCh: make(chan tea.Model, 1),
		doneCh:  make(chan bool, 1),
	}

	//nolint: staticcheck
	tm.program = tea.NewProgram(
		m,
		tea.WithInput(tm.in),
		tea.WithOutput(tm.out),
		tea.WithoutSignals(),
		tea.WithANSICompressor(), // 这有助于减少运行之间的差异
	)

	interruptions := make(chan os.Signal, 1)
	signal.Notify(interruptions, syscall.SIGINT)
	go func() {
		m, err := tm.program.Run()
		if err != nil {
			tb.Fatalf("应用失败: %s", err)
		}
		tm.modelCh <- m
		tm.doneCh <- true
	}()
	go func() {
		<-interruptions
		signal.Stop(interruptions)
		tb.Log("已中断")
		tm.program.Kill()
	}()

	var opts TestModelOptions
	for _, opt := range options {
		opt(&opts)
	}

	if opts.size.Width != 0 {
		tm.program.Send(opts.size)
	}
	return tm
}

func (tm *TestModel) waitDone(tb testing.TB, opts []FinalOpt) {
	tm.done.Do(func() {
		fopts := FinalOpts{}
		for _, opt := range opts {
			opt(&fopts)
		}
		if fopts.timeout > 0 {
			select {
			case <-time.After(fopts.timeout):
				if fopts.onTimeout == nil {
					tb.Fatalf("超时后 %s", fopts.timeout)
				}
				fopts.onTimeout(tb)
			case <-tm.doneCh:
			}
		} else {
			<-tm.doneCh
		}
	})
}

// FinalOpts 表示 FinalModel 和 FinalOutput 的选项。
type FinalOpts struct {
	timeout   time.Duration
	onTimeout func(tb testing.TB)
}

// FinalOpt 更改 FinalOpts。
type FinalOpt func(opts *FinalOpts)

// WithTimeoutFn 允许定义 WaitFinished 超时时发生的情况。
func WithTimeoutFn(fn func(tb testing.TB)) FinalOpt {
	return func(opts *FinalOpts) {
		opts.onTimeout = fn
	}
}

// WithFinalTimeout 允许设置 FinalModel 和 FinalOutput 等待程序完成的超时时间。
func WithFinalTimeout(d time.Duration) FinalOpt {
	return func(opts *FinalOpts) {
		opts.timeout = d
	}
}

// WaitFinished 等待应用程序完成。
// 此方法仅在程序完成运行或超时时返回。
func (tm *TestModel) WaitFinished(tb testing.TB, opts ...FinalOpt) {
	tm.waitDone(tb, opts)
}

// FinalModel 返回 program.Run() 产生的结果模型。
// 此方法仅在程序完成运行或超时时返回。
func (tm *TestModel) FinalModel(tb testing.TB, opts ...FinalOpt) tea.Model {
	tm.waitDone(tb, opts)
	select {
	case m := <-tm.modelCh:
		if m != nil {
			tm.model = m
		}
		return tm.model
	default:
		return tm.model
	}
}

// FinalOutput 返回程序的最终输出 io.Reader。
// 此方法仅在程序完成运行或超时时返回。
func (tm *TestModel) FinalOutput(tb testing.TB, opts ...FinalOpt) io.Reader {
	tm.waitDone(tb, opts)
	return tm.Output()
}

// Output 返回程序的当前输出 io.Reader。
func (tm *TestModel) Output() io.Reader {
	return tm.out
}

// Send 向底层程序发送消息。
func (tm *TestModel) Send(m tea.Msg) {
	tm.program.Send(m)
}

// Quit 退出程序并释放终端。
func (tm *TestModel) Quit() error {
	tm.program.Quit()
	return nil
}

// Type 将给定文本输入到给定程序中。
func (tm *TestModel) Type(s string) {
	for _, c := range []byte(s) {
		tm.Send(tea.KeyMsg{
			Runes: []rune{rune(c)},
			Type:  tea.KeyRunes,
		})
	}
}

// GetProgram 获取 TestModel 的程序。
func (tm *TestModel) GetProgram() *tea.Program {
	return tm.program
}

// RequireEqualOutput 是一个辅助函数，用于断言给定输出与黄金文件中的预期输出匹配，
// 如果不匹配则打印其差异。
//
// 重要：这使用系统的 `diff` 工具。
//
// 你可以通过使用 -update 标志运行测试来更新黄金文件。
func RequireEqualOutput(tb testing.TB, out []byte) {
	tb.Helper()
	golden.RequireEqual(tb, out)
}

func safe(rw io.ReadWriter) io.ReadWriter {
	return &safeReadWriter{rw: rw}
}

// safeReadWriter 实现 io.ReadWriter，但会锁定读写操作。
type safeReadWriter struct {
	rw io.ReadWriter
	m  sync.RWMutex
}

// Read 实现 io.ReadWriter。
func (s *safeReadWriter) Read(p []byte) (n int, err error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.rw.Read(p) //nolint: wrapcheck
}

// Write 实现 io.ReadWriter。
func (s *safeReadWriter) Write(p []byte) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	return s.rw.Write(p) //nolint: wrapcheck
}
