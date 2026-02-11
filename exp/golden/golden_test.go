package golden

import "testing"

// TestRequireEqualUpdate 测试更新模式下的RequireEqual

func TestRequireEqualUpdate(t *testing.T) {
	*update = true
	RequireEqual(t, []byte("test"))
}

// TestRequireEqualNoUpdate 测试非更新模式下的RequireEqual

func TestRequireEqualNoUpdate(t *testing.T) {
	*update = false
	RequireEqual(t, []byte("test"))
}

// TestRequireWithLineBreaks 测试包含换行符的RequireEqual

func TestRequireWithLineBreaks(t *testing.T) {
	*update = false
	RequireEqual(t, []byte("foo\nbar\nbaz\n"))
}

// TestTypes 测试不同类型的RequireEqual

func TestTypes(t *testing.T) {
	*update = false

	t.Run("字节切片", func(t *testing.T) {
		RequireEqual(t, []byte("test"))
	})
	t.Run("字符串", func(t *testing.T) {
		RequireEqual(t, "test")
	})
}
