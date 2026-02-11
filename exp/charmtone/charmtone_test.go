package charmtone

import (
	"strconv"
	"strings"
	"testing"
)

// TestValidateHexes 测试所有颜色键的十六进制颜色值是否有效
func TestValidateHexes(t *testing.T) {
	for _, key := range Keys() {
		hex := strings.TrimPrefix(key.Hex(), "#")
		if len(hex) != 6 && len(hex) != 3 {
			t.Errorf("颜色键 %s: 十六进制长度无效，长度为 %d，值为 %s", key, len(hex), key.Hex())
		}
		if _, err := strconv.ParseUint(hex, 16, 32); err != nil {
			t.Errorf("颜色键 %s: 十六进制值无效，值为 %s", key, key.Hex())
		}
	}
}
