package ansi

import "testing"

// TestSetProgress 测试设置进度条
func TestSetProgress(t *testing.T) {
	expect := "\x1b]9;4;1;50\x07"
	got := SetProgressBar(50)
	if expect != got {
		t.Errorf("SetProgressBar(50) = %q, want %q", got, expect)
	}
}

// TestSetProgressNegative 测试设置负数进度条
func TestSetProgressNegative(t *testing.T) {
	expect := "\x1b]9;4;1;0\x07"
	got := SetProgressBar(-2)
	if expect != got {
		t.Errorf("SetProgressBar(-2) = %q, want %q", got, expect)
	}
}

// TestSetProgressAbove100 测试设置超过100的进度条
func TestSetProgressAbove100(t *testing.T) {
	expect := "\x1b]9;4;1;100\x07"
	got := SetProgressBar(200)
	if expect != got {
		t.Errorf("SetProgressBar(200) = %q, want %q", got, expect)
	}
}

// TestSetErrorProgress 测试设置错误状态进度条
func TestSetErrorProgress(t *testing.T) {
	expect := "\x1b]9;4;2;50\x07"
	got := SetErrorProgressBar(50)
	if expect != got {
		t.Errorf("SetErrorProgressBar(50) = %q, want %q", got, expect)
	}
}

// TestSetErrorProgressNegative 测试设置错误状态负数进度条
func TestSetErrorProgressNegative(t *testing.T) {
	expect := "\x1b]9;4;2;0\x07"
	got := SetErrorProgressBar(-2)
	if expect != got {
		t.Errorf("SetErrorProgressBar(-2) = %q, want %q", got, expect)
	}
}

// TestSetErrorProgressAbove100 测试设置错误状态超过100的进度条
func TestSetErrorProgressAbove100(t *testing.T) {
	expect := "\x1b]9;4;2;100\x07"
	got := SetErrorProgressBar(200)
	if expect != got {
		t.Errorf("SetErrorProgressBar(200) = %q, want %q", got, expect)
	}
}

// TestSetWarningProgress 测试设置警告状态进度条
func TestSetWarningProgress(t *testing.T) {
	expect := "\x1b]9;4;4;50\x07"
	got := SetWarningProgressBar(50)
	if expect != got {
		t.Errorf("SetWarningProgressBar(50) = %q, want %q", got, expect)
	}
}

// TestSetWarningProgressNegative 测试设置警告状态负数进度条
func TestSetWarningProgressNegative(t *testing.T) {
	expect := "\x1b]9;4;4;0\x07"
	got := SetWarningProgressBar(-2)
	if expect != got {
		t.Errorf("SetWarningProgressBar(-2) = %q, want %q", got, expect)
	}
}

// TestSetWarningProgressAbove100 测试设置警告状态超过100的进度条
func TestSetWarningProgressAbove100(t *testing.T) {
	expect := "\x1b]9;4;4;100\x07"
	got := SetWarningProgressBar(200)
	if expect != got {
		t.Errorf("SetWarningProgressBar(200) = %q, want %q", got, expect)
	}
}
