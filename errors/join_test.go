package errors

import (
	"fmt"
	"testing"
)

// TestJoin 测试 Join 函数的功能
func TestJoin(t *testing.T) {
	// 测试所有错误都是 nil 的情况
	t.Run("nil", func(t *testing.T) {
		err := Join(nil, nil, nil)
		if err != nil {
			t.Errorf("期望 nil，实际得到 %v", err)
		}
	})
	// 测试只有一个错误的情况
	t.Run("one err", func(t *testing.T) {
		expected := fmt.Errorf("fake")
		err := Join(nil, expected, nil)
		je := err.(*joinError)
		un := je.Unwrap()
		if len(un) != 1 {
			t.Fatalf("期望 1 个错误，实际得到 %d 个", len(un))
		}
		if s := un[0].Error(); s != expected.Error() {
			t.Errorf("期望 %v，实际得到 %v", expected, un[0])
		}
		if s := err.Error(); s != expected.Error() {
			t.Errorf("期望 %s，实际得到 %s", expected, err)
		}
	})
	// 测试多个错误的情况
	t.Run("many errs", func(t *testing.T) {
		expected1 := fmt.Errorf("fake 1")
		expected2 := fmt.Errorf("fake 2")
		err := Join(nil, expected1, nil, nil, expected2, nil)
		je := err.(*joinError)
		un := je.Unwrap()
		if len(un) != 2 {
			t.Fatalf("期望 2 个错误，实际得到 %d 个", len(un))
		}
		if s := un[0].Error(); s != expected1.Error() {
			t.Errorf("期望 %v，实际得到 %v", expected1, un[0])
		}
		if s := un[1].Error(); s != expected2.Error() {
			t.Errorf("期望 %v，实际得到 %v", expected2, un[1])
		}
		expectedS := expected1.Error() + "\n" + expected2.Error()
		if s := err.Error(); s != expectedS {
			t.Errorf("期望 %s，实际得到 %s", expectedS, err)
		}
	})
}
