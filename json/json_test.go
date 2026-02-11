package json

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

// TestReader 测试Reader函数是否能正确将map转换为JSON
func TestReader(t *testing.T) {
	// 创建一个Reader，传入一个map
	r := Reader(map[string]int{
		"foo": 2,
	})
	// 读取所有数据
	bts, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	// 验证生成的JSON是否正确
	if string(bts) != `{"foo":2}` {
		t.Fatalf("生成的JSON错误: %s", string(bts))
	}
}

// TestFrom 测试From函数是否能正确从Reader中读取JSON并转换为map
func TestFrom(t *testing.T) {
	// 定义输入map
	in := map[string]int{"foo": 10, "bar": 20}
	// 从Reader中读取并转换为map
	m, err := From(Reader(in), map[string]int{})
	if err != nil {
		t.Fatalf("意外错误: %v", err)
	}
	// 验证转换后的map是否与原map相等
	if !reflect.DeepEqual(m, in) {
		t.Fatalf("两个map应该相等: %v vs %v", in, m)
	}
}

// TestErrReader 测试ErrorReader是否能正确传递错误
func TestErrReader(t *testing.T) {
	// 创建一个错误
	err := fmt.Errorf("foo")
	// 读取ErrorReader的数据，应该返回相同的错误
	_, err2 := io.ReadAll(&ErrorReader{err})
	if err != err2 {
		t.Fatalf("应该返回相同的错误")
	}
}

// TestIsValid 测试IsValid函数是否能正确判断JSON是否有效
func TestIsValid(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name string // 测试用例名称
		data any    // 测试数据
		want bool   // 期望结果
	}{
		{
			name: "空字符串",
			data: "",
			want: false,
		},
		{
			name: "空字节数组",
			data: []byte(""),
			want: false,
		},
		{
			name: "有效的JSON字符串",
			data: `{"foo": 2}`,
			want: true,
		},
		{
			name: "有效的JSON字节数组",
			data: []byte(`{"foo": 2}`),
			want: true,
		},
		{
			name: "无效的JSON字符串",
			data: `{"foo": 2`,
			want: false,
		},
		{
			name: "无效的JSON字节数组",
			data: []byte(`{"foo": 2`),
			want: false,
		},
	}
	// 运行所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			// 根据数据类型调用不同的IsValid函数
			switch v := tt.data.(type) {
			case string:
				got = IsValid(v)
			case []byte:
				got = IsValid(v)
			default:
				t.Fatalf("不支持的类型: %T", tt.data)
			}
			// 验证结果是否符合预期
			if got != tt.want {
				t.Errorf("IsValid() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}
