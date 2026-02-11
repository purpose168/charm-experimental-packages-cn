package higherorder

import (
	"reflect"
	"strings"
	"testing"
)

func Test_Foldl(t *testing.T) {
	x := Foldl(func(a, b int) int {
		return a + b
	}, 0, []int{1, 2, 3})

	const expect = 6
	if x != expect {
		t.Errorf("Expected %d, got %d", expect, x)
	}
}

func Test_Foldr(t *testing.T) {
	x := Foldl(func(a, b int) int {
		return a - b
	}, 6, []int{1, 2, 3})

	const expect = 0
	if x != expect {
		t.Errorf("Expected %d, got %d", expect, x)
	}
}

func Test_Map(t *testing.T) {
	{
		// 对整数进行映射，返回每个整数的平方。
		// （接受整数，返回整数。）
		x := Map(func(a int) int {
			return a * a
		}, []int{2, 3, 4})

		expected := []int{4, 9, 16}
		for i, v := range x {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	}
	{
		// 对字符串进行映射，返回每个字符串的长度。
		// （接受字符串，返回整数。）
		x := Map(func(a string) int {
			return len([]rune(a))
		}, []string{"one", "two", "three"})

		expected := []int{3, 3, 5}
		for i, v := range x {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	}
}

func Test_Filter(t *testing.T) {
	t.Run("with string slices", func(t *testing.T) {
		got := Filter(func(a string) bool {
			return strings.HasPrefix(a, "t")
		}, []string{"one", "two", "three"})

		want := []string{"two", "three"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("with int slices", func(t *testing.T) {
		got := Filter(func(a int) bool {
			return a%2 == 0
		}, []int{1, 2, 3, 4, 5})

		want := []int{2, 4}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
