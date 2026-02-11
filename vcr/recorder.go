package vcr

import (
	"path/filepath"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// Recorder 是go-vcr Recorder的别名。
type Recorder = recorder.Recorder

type options struct {
	mode           recorder.Mode
	keepAllHeaders bool
}

// Option 定义了用于配置VCR记录器的函数选项。
type Option func(*options) error

// WithMode 设置记录器模式。
func WithMode(mode recorder.Mode) Option {
	return func(o *options) error {
		o.mode = mode
		return nil
	}
}

// WithKeepAllHeaders 配置记录器保留所有HTTP头。
func WithKeepAllHeaders() Option {
	return func(o *options) error {
		o.keepAllHeaders = true
		return nil
	}
}

// NewRecorder 使用提供的选项为给定测试创建一个新的VCR记录器。
func NewRecorder(t *testing.T, opts ...Option) *Recorder {
	o := options{
		mode:           recorder.ModeRecordOnce,
		keepAllHeaders: false,
	}
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			t.Fatalf("vcr: failed to apply option: %v", err)
		}
	}

	cassetteName := filepath.Join("testdata", t.Name())

	r, err := recorder.New(
		cassetteName,
		recorder.WithMode(o.mode),
		recorder.WithMatcher(customMatcher(t)),
		recorder.WithMarshalFunc(customMarshaler),
		recorder.WithSkipRequestLatency(true), // 禁用睡眠来模拟响应时间，使测试更快
		recorder.WithHook(hookRemoveHeaders(o.keepAllHeaders), recorder.AfterCaptureHook),
	)
	if err != nil {
		t.Fatalf("vcr: failed to create recorder: %v", err)
	}

	t.Cleanup(func() {
		if err := r.Stop(); err != nil {
			t.Errorf("vcr: failed to stop recorder: %v", err)
		}
	})

	return r
}
