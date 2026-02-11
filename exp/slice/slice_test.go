package slice_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/purpose168/charm-experimental-packages-cn/exp/slice"
)

// TestGroupBy 测试 GroupBy 函数的功能
func TestGroupBy(t *testing.T) {
	expected := map[string][]string{
		"a": {"andrey", "ayman"},
		"b": {"bash"},
		"c": {"carlos", "christian"},
		"r": {"raphael"},
	}
	input := []string{
		"andrey",
		"ayman",
		"bash",
		"carlos",
		"christian",
		"raphael",
	}
	output := slice.GroupBy(input, func(s string) string { return string(s[0]) })

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("期望 %v, 得到 %v", expected, output)
	}
}

// TestTake 测试 Take 函数的功能
func TestTake(t *testing.T) {
	for i, tc := range []struct {
		input    []int
		take     int
		expected []int
	}{
		{
			input:    []int{1, 2, 3, 4, 5},
			take:     3,
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{1, 2, 3},
			take:     5,
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{},
			take:     2,
			expected: []int{},
		},
		{
			input:    []int{1, 2, 3},
			take:     0,
			expected: []int{},
		},
		{
			input:    nil,
			take:     2,
			expected: []int{},
		},
	} {
		actual := slice.Take(tc.input, tc.take)
		if len(actual) != len(tc.expected) {
			t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tc.expected, actual)
		}
	}
}

// TestLast 测试 Last 函数的功能
func TestLast(t *testing.T) {
	for i, tc := range []struct {
		input    []int
		ok       bool
		expected int
	}{
		{
			input:    []int{1, 2, 3, 4, 5},
			ok:       true,
			expected: 5,
		},
		{
			input:    []int{1, 2, 3},
			ok:       true,
			expected: 3,
		},
		{
			input:    []int{1},
			ok:       true,
			expected: 1,
		},
		{
			input:    []int{},
			ok:       false,
			expected: 0,
		},
	} {
		actual, ok := slice.Last(tc.input)
		if ok != tc.ok {
			t.Errorf("测试 %d: 期望 ok %v, 得到 %v", i, tc.ok, ok)
		}
		if actual != tc.expected {
			t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tc.expected, actual)
		}
	}
}

// TestUniq 测试 Uniq 函数的功能
func TestUniq(t *testing.T) {
	for i, tc := range []struct {
		input    []int
		expected []int
	}{
		{
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			input:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			input:    []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			input:    []int{},
			expected: []int{},
		},
	} {
		actual := slice.Uniq(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tc.expected, actual)
		}
	}
}

// TestIntersperse 测试 Intersperse 函数的功能
func TestIntersperse(t *testing.T) {
	for i, tc := range []struct {
		input    []string
		insert   string
		expected []string
	}{
		{
			input:    []string{},
			insert:   "-",
			expected: []string{},
		},
		{
			input:    []string{"a"},
			insert:   "-",
			expected: []string{"a"},
		},
		{
			input:    []string{"a", "b"},
			insert:   "-",
			expected: []string{"a", "-", "b"},
		},
		{
			input:    []string{"a", "b", "c"},
			insert:   "-",
			expected: []string{"a", "-", "b", "-", "c"},
		},
	} {
		actual := slice.Intersperse(tc.input, tc.insert)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tc.expected, actual)
		}
	}
}

// TestContainsAny 测试 ContainsAny 函数的功能
func TestContainsAny(t *testing.T) {
	for i, tc := range []struct {
		input    []string
		values   []string
		expected bool
	}{
		{
			input:    []string{"a", "b", "c"},
			values:   []string{"a", "b"},
			expected: true,
		},
		{
			input:    []string{"a", "b", "c"},
			values:   []string{"d", "e"},
			expected: false,
		},
		{
			input:    []string{"a", "b", "c"},
			values:   []string{"c", "d"},
			expected: true,
		},
		{
			input:    []string{},
			values:   []string{"d", "e"},
			expected: false,
		},
	} {
		actual := slice.ContainsAny(tc.input, tc.values...)
		if actual != tc.expected {
			t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tc.expected, actual)
		}
	}
}

// TestShift 测试 Shift 函数的功能
func TestShift(t *testing.T) {
	for i, tc := range []struct {
		input         []int
		ok            bool
		expectedVal   int
		expectedSlice []int
	}{
		{
			input:         []int{1, 2, 3, 4, 5},
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{2, 3, 4, 5},
		},
		{
			input:         []int{1, 2, 3},
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{2, 3},
		},
		{
			input:         []int{1},
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{},
		},
		{
			input:         []int{},
			ok:            false,
			expectedVal:   0,
			expectedSlice: []int{},
		},
	} {
		actual, newSlice, ok := slice.Shift(tc.input)
		if ok != tc.ok {
			t.Errorf("测试 %d: 期望 ok %v, 得到 %v", i, tc.ok, ok)
		}
		if actual != tc.expectedVal {
			t.Errorf("测试 %d: 期望 val %v, 得到 %v", i, tc.expectedVal, actual)
		}
		if !reflect.DeepEqual(newSlice, tc.expectedSlice) {
			t.Errorf("测试 %d: 期望 slice %v, 得到 %v", i, tc.expectedSlice, newSlice)
		}
	}
}

// TestPop 测试 Pop 函数的功能
func TestPop(t *testing.T) {
	for i, tc := range []struct {
		input         []int
		ok            bool
		expectedVal   int
		expectedSlice []int
	}{
		{
			input:         []int{1, 2, 3, 4, 5},
			ok:            true,
			expectedVal:   5,
			expectedSlice: []int{1, 2, 3, 4},
		},
		{
			input:         []int{1, 2, 3},
			ok:            true,
			expectedVal:   3,
			expectedSlice: []int{1, 2},
		},
		{
			input:         []int{1},
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{},
		},
		{
			input:         []int{},
			ok:            false,
			expectedVal:   0,
			expectedSlice: []int{},
		},
	} {
		actual, newSlice, ok := slice.Pop(tc.input)
		if ok != tc.ok {
			t.Errorf("测试 %d: 期望 ok %v, 得到 %v", i, tc.ok, ok)
		}
		if actual != tc.expectedVal {
			t.Errorf("测试 %d: 期望 val %v, 得到 %v", i, tc.expectedVal, actual)
		}
		if !reflect.DeepEqual(newSlice, tc.expectedSlice) {
			t.Errorf("测试 %d: 期望 slice %v, 得到 %v", i, tc.expectedSlice, newSlice)
		}
	}
}

// TestDeleteAt 测试 DeleteAt 函数的功能
func TestDeleteAt(t *testing.T) {
	for i, tc := range []struct {
		input         []int
		index         int
		ok            bool
		expectedVal   int
		expectedSlice []int
	}{
		{
			input:         []int{1, 2, 3, 4, 5},
			index:         2,
			ok:            true,
			expectedVal:   3,
			expectedSlice: []int{1, 2, 4, 5},
		},
		{
			input:         []int{1, 2, 3},
			index:         0,
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{2, 3},
		},
		{
			input:         []int{1, 2, 3},
			index:         2,
			ok:            true,
			expectedVal:   3,
			expectedSlice: []int{1, 2},
		},
		{
			input:         []int{1},
			index:         0,
			ok:            true,
			expectedVal:   1,
			expectedSlice: []int{},
		},
		{
			input:         []int{},
			index:         0,
			ok:            false,
			expectedVal:   0,
			expectedSlice: []int{},
		},
	} {
		actual, newSlice, ok := slice.DeleteAt(tc.input, tc.index)
		if ok != tc.ok {
			t.Errorf("测试 %d: 期望 ok %v, 得到 %v", i, tc.ok, ok)
		}
		if actual != tc.expectedVal {
			t.Errorf("测试 %d: 期望 val %v, 得到 %v", i, tc.expectedVal, actual)
		}
		if !reflect.DeepEqual(newSlice, tc.expectedSlice) {
			t.Errorf("测试 %d: 期望 slice %v, 得到 %v", i, tc.expectedSlice, newSlice)
		}
	}
}

// TestIsSubset 测试 IsSubset 函数的功能
func TestIsSubset(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected bool
	}{
		// 基本子集情况
		{
			name:     "空集是空集的子集",
			a:        []string{},
			b:        []string{},
			expected: true,
		},
		{
			name:     "空集是非空集的子集",
			a:        []string{},
			b:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "非空集不是空集的子集",
			a:        []string{"a"},
			b:        []string{},
			expected: false,
		},
		{
			name:     "单个元素是子集",
			a:        []string{"b"},
			b:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "单个元素不是子集",
			a:        []string{"d"},
			b:        []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "多个元素是子集",
			a:        []string{"a", "c"},
			b:        []string{"a", "b", "c", "d"},
			expected: true,
		},
		{
			name:     "多个元素不是子集",
			a:        []string{"a", "e"},
			b:        []string{"a", "b", "c", "d"},
			expected: false,
		},
		{
			name:     "相等的集合是子集",
			a:        []string{"a", "b", "c"},
			b:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "较大的集合不是较小集合的子集",
			a:        []string{"a", "b", "c", "d"},
			b:        []string{"a", "b"},
			expected: false,
		},

		// 顺序无关
		{
			name:     "不同顺序的子集",
			a:        []string{"c", "a"},
			b:        []string{"b", "a", "d", "c"},
			expected: true,
		},

		// 重复元素处理
		{
			name:     "子集有重复元素",
			a:        []string{"a", "a", "b"},
			b:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "超集有重复元素",
			a:        []string{"a", "b"},
			b:        []string{"a", "a", "b", "b", "c"},
			expected: true,
		},
		{
			name:     "两者都有重复元素",
			a:        []string{"a", "a", "b"},
			b:        []string{"a", "a", "b", "b", "c"},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slice.IsSubset(tt.a, tt.b)
			if actual != tt.expected {
				t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tt.expected, actual)
			}
		})
	}
}

// TestIsSubsetWithInts 测试 IsSubset 函数处理整数的功能
func TestIsSubsetWithInts(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected bool
	}{
		{
			name:     "整数子集",
			a:        []int{1, 3},
			b:        []int{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "整数不是子集",
			a:        []int{1, 5},
			b:        []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "空整数子集",
			a:        []int{},
			b:        []int{1, 2, 3},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slice.IsSubset(tt.a, tt.b)
			if actual != tt.expected {
				t.Errorf("测试 %d: 期望 %v, 得到 %v", i, tt.expected, actual)
			}
		})
	}
}

// TestMapStringToInt 测试 Map 函数将字符串映射为整数的功能
func TestMapStringToInt(t *testing.T) {
	seq := slices.Values([]string{"a", "ab", "abc", "abcd"})
	mapped := slice.Map(seq, func(s string) int { return len(s) })
	expected := []int{1, 2, 3, 4}

	result := slices.Collect(mapped)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("期望 %v, 得到 %v", expected, result)
	}
}
