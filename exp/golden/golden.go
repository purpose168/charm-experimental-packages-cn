// Package golden 提供一个辅助函数来断言测试的输出。
package golden

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/aymanbagabas/go-udiff"
)

var update = flag.Bool("update", false, "update .golden files")

// RequireEqual 是一个辅助函数，用于断言给定的输出与金色文件中的预期输出一致，
// 如果不一致则打印差异。
//
// 金色文件包含测试的原始预期输出，可以包含控制代码和转义序列。
// 当比较测试输出时，[RequireEqual] 会在将输出与金色文件比较之前
// 转义控制代码和序列。
//
// 您可以通过使用 -update 标志运行测试来更新金色文件。
func RequireEqual[T []byte | string](tb testing.TB, out T) {
	tb.Helper()

	golden := filepath.Join("testdata", tb.Name()+".golden")
	if *update {
		if err := os.MkdirAll(filepath.Dir(golden), 0o750); err != nil { //nolint: mnd
			tb.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(out), 0o600); err != nil { //nolint: mnd
			tb.Fatal(err)
		}
	}

	goldenBts, err := os.ReadFile(golden)
	if err != nil {
		tb.Fatal(err)
	}

	goldenStr := normalizeWindowsLineBreaks(string(goldenBts))
	goldenStr = escapeSeqs(goldenStr)
	outStr := escapeSeqs(string(out))

	diff := udiff.Unified("golden", "run", goldenStr, outStr)
	if diff != "" {
		tb.Fatalf("output does not match, expected:\n\n%s\n\ngot:\n\n%s\n\ndiff:\n\n%s", goldenStr, outStr, diff)
	}
}

// RequireEqualEscape 是一个辅助函数，用于断言给定的输出与金色文件中的预期输出一致，
// 如果不一致则打印差异。
//
// 已弃用: 使用 [RequireEqual] 代替。
func RequireEqualEscape(tb testing.TB, out []byte, escapes bool) { //nolint:revive
	RequireEqual(tb, out)
}

// escapeSeqs 转义给定字符串中的控制代码和转义序列。
// 唯一保留的例外是换行符。
func escapeSeqs(in string) string {
	s := strings.Split(in, "\n")
	for i, l := range s {
		q := strconv.Quote(l)
		q = strings.TrimPrefix(q, `"`)
		q = strings.TrimSuffix(q, `"`)
		s[i] = q
	}
	return strings.Join(s, "\n")
}

// normalizeWindowsLineBreaks 将所有 \r\n 替换为 \n。
// 这是必要的，因为 Git for Windows 默认以 \r\n 检出文件。
func normalizeWindowsLineBreaks(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, "\r\n", "\n")
	}
	return str
}
