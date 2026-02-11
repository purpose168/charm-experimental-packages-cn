// 包 json 提供了处理 JSON 的辅助函数。
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Reader 接收一个输入，将其序列化为 JSON 并返回一个 io.Reader。
func Reader[T any](v T) io.Reader {
	bts, err := json.Marshal(v)
	if err != nil {
		return &ErrorReader{err}
	}
	return bytes.NewReader(bts)
}

// ErrorReader 是一个始终返回给定错误的读取器。
type ErrorReader struct {
	err error
}

func (r *ErrorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

// From 解析包含 JSON 的 io.Reader。
func From[T any](r io.Reader, t T) (T, error) {
	bts, err := io.ReadAll(r)
	if err != nil {
		return t, fmt.Errorf("failed to read response: %w", err)
	}
	if err := json.Unmarshal(bts, &t); err != nil {
		return t, fmt.Errorf("failed to parse body: %w: %s", err, bts)
	}
	return t, nil
}

// Write 将给定数据以 JSON 形式写入。
func Write(w http.ResponseWriter, data any) error {
	bts, err := json.Marshal(data)
	if err != nil {
		return err //nolint:wrapcheck
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bts)
	return err //nolint:wrapcheck
}

// IsValid 检查给定数据是否为有效的 JSON。
func IsValid[T string | []byte](data T) bool {
	if len(data) == 0 { // hot path
		return false
	}
	var m json.RawMessage
	err := json.Unmarshal([]byte(data), &m)
	return err == nil
}
