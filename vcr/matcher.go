package vcr

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// customMatcher 创建一个自定义的请求匹配器，处理JSON请求中的键顺序问题。
func customMatcher(t *testing.T) recorder.MatcherFunc {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Body == nil || r.Body == http.NoBody {
			return cassette.DefaultMatcher(r, i)
		}
		if r.Method != i.Method || r.URL.String() != i.URL {
			return false
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("vcr: failed to read request body")
		}
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		// 一些提供商有时会生成键顺序不同的JSON请求，这意味着直接字符串比较会失败。
		// 如果没有匹配，我们会回退到反序列化内容。
		requestContent := normalizeLineEndings(reqBody)
		cassetteContent := normalizeLineEndings(i.Body)
		if requestContent == cassetteContent {
			return true
		}
		var content1, content2 any
		if err := json.Unmarshal([]byte(requestContent), &content1); err != nil {
			printDiff(t, requestContent, cassetteContent)
			return false
		}
		if err := json.Unmarshal([]byte(cassetteContent), &content2); err != nil {
			printDiff(t, requestContent, cassetteContent)
			return false
		}
		if isEqual := reflect.DeepEqual(content1, content2); !isEqual {
			printDiff(t, requestContent, cassetteContent)
			return false
		}
		return true
	}
}

// normalizeLineEndings 不仅将 `\r\n` 替换为 `\n`，
// 还将 `\\r\\n` 替换为 `\\n`。这是因为我们也希望替换JSON字符串中的内容。
func normalizeLineEndings[T string | []byte](s T) string {
	str := string(s)
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, `\r\n`, `\n`)
	return str
}

// printDiff 打印请求内容和盒式磁带内容之间的差异。
func printDiff(t *testing.T, requestContent, cassetteContent string) {
	t.Logf("Request interaction not found for %q.\nDiff:\n%s", t.Name(), cmp.Diff(cassetteContent, requestContent))
}
