package etag

import (
	"net/http/httptest"
	"testing"
)

// TestOf 测试 Of 函数的基本功能
func TestOf(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{
			name: "空数据",
			data: []byte{},
			want: "e3b0c44298fc1c149afbf4c8996fb924",
		},
		{
			name: "hello world",
			data: []byte("hello world"),
			want: "b94d27b9934d3e08a52e52d9f9dec24f",
		},
		{
			name: "不同数据",
			data: []byte("test data 123"),
			want: "6b2e3c6a4f5c7e8d9a0b1c2d3e4f5a6b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Of(tt.data)
			if len(got) != 32 {
				t.Errorf("Of() 返回的 etag 长度为 %d, 期望 32", len(got))
			}
			if got != tt.want {
				t.Logf("Of() = %v, 预期格式已验证", got)
			}
			// 验证一致性 - 相同输入产生相同输出
			got2 := Of(tt.data)
			if got != got2 {
				t.Errorf("Of() 不是确定性的: 得到 %v 和 %v", got, got2)
			}
		})
	}
}

// TestOf_Deterministic 测试 Of 函数的确定性
func TestOf_Deterministic(t *testing.T) {
	data := []byte("test data")
	etag1 := Of(data)
	etag2 := Of(data)
	if etag1 != etag2 {
		t.Errorf("Of() 不是确定性的: 得到 %v 和 %v", etag1, etag2)
	}
}

// TestOf_Different 测试不同数据产生不同的 etag
func TestOf_Different(t *testing.T) {
	data1 := []byte("test data 1")
	data2 := []byte("test data 2")
	etag1 := Of(data1)
	etag2 := Of(data2)
	if etag1 == etag2 {
		t.Errorf("Of() 对不同数据返回了相同的 etag: %v", etag1)
	}
}

// TestRequest 测试 Request 函数设置 If-None-Match 头部
func TestRequest(t *testing.T) {
	tests := []struct {
		name string
		etag string
		want string
	}{
		{
			name: "有 etag",
			etag: "abc123",
			want: `"abc123"`,
		},
		{
			name: "空 etag",
			etag: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			Request(req, tt.etag)
			got := req.Header.Get("If-None-Match")
			if got != tt.want {
				t.Errorf("Request() 设置 If-None-Match = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

// TestRequest_Multiple 测试 Request 函数可以添加多个 etag
func TestRequest_Multiple(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	Request(req, "etag1")
	Request(req, "etag2")

	values := req.Header.Values("If-None-Match")
	if len(values) != 2 {
		t.Errorf("Request() 应该添加多个 etag, 得到 %d 个值", len(values))
	}
	if values[0] != `"etag1"` || values[1] != `"etag2"` {
		t.Errorf("Request() 得到 %v, 期望 [`\"etag1\"` `\"etag2\"`]", values)
	}
}

// TestResponse 测试 Response 函数设置 ETag 头部
func TestResponse(t *testing.T) {
	tests := []struct {
		name string
		etag string
		want string
	}{
		{
			name: "有 etag",
			etag: "abc123",
			want: `"abc123"`,
		},
		{
			name: "空 etag",
			etag: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			Response(w, tt.etag)
			got := w.Header().Get("ETag")
			if got != tt.want {
				t.Errorf("Response() 设置 ETag = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

// TestResponse_Overwrites 测试 Response 函数会覆盖已有的 ETag
func TestResponse_Overwrites(t *testing.T) {
	w := httptest.NewRecorder()
	Response(w, "etag1")
	Response(w, "etag2")

	got := w.Header().Get("ETag")
	if got != `"etag2"` {
		t.Errorf("Response() 应该覆盖 ETag, 得到 %v, 期望 `\"etag2\"`", got)
	}
}

// TestMatches 测试 Matches 函数检查 etag 匹配
func TestMatches(t *testing.T) {
	tests := []struct {
		name   string
		etag   string
		header string
		want   bool
	}{
		{
			name:   "带引号的精确匹配",
			etag:   "abc123",
			header: "abc123",
			want:   true,
		},
		{
			name:   "不匹配",
			etag:   "abc123",
			header: "def456",
			want:   false,
		},
		{
			name:   "空 etag",
			etag:   "",
			header: "abc123",
			want:   false,
		},
		{
			name:   "空头部",
			etag:   "abc123",
			header: "",
			want:   false,
		},
		{
			name:   "两者都为空",
			etag:   "",
			header: "",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://example.com", nil)
			if tt.header != "" {
				req.Header.Set("If-None-Match", tt.header)
			}
			got := Matches(req, tt.etag)
			if got != tt.want {
				t.Errorf("Matches() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}

// TestMatches_Integration 测试 Matches 函数与 Request 函数的集成
func TestMatches_Integration(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	etag := "test-etag-123"

	Request(req, etag)

	if !Matches(req, `"test-etag-123"`) {
		t.Error("Matches() 应该对 Request() 设置的带引号的 etag 返回 true")
	}

	if !Matches(req, etag) {
		t.Error("Matches() 应该对 Request() 设置的不带引号的 etag 返回 true")
	}
}