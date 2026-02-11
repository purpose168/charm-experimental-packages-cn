// Package snapshot 提供了用于处理终端快照的辅助函数。
package snapshot

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/purpose168/charm-experimental-packages-cn/vttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// update 指示是否更新测试数据文件。
var update = flag.Bool("update", false, "更新测试数据文件")

// Snapshotter 是一个接口，用于可以生成其状态快照的类型。
type Snapshotter interface {
	Snapshot() vttest.Snapshot
}

// Imager 表示可以生成其状态图像表示的类型。
type Imager interface {
	Image() image.Image
}

// TestdataEqualf 将给定 [Snapshotter] 的快照与存储在 "testdata" 目录中的预期快照进行比较，
// 使用提供的格式和参数构造文件名。
//
// 如果快照不匹配，它会在 testing.TB 上报告错误。
func TestdataEqualf(tb testing.TB, expectedNameSuffix string, actual Snapshotter, format string, args ...any) {
	tb.Helper()
	testdataEq(tb, expectedNameSuffix, actual, func(expected, actual vttest.Snapshot) {
		assert.Equalf(tb, expected, actual, format, args...)
	})
}

// TestdataEqual 将给定 [Snapshotter] 的快照与存储在 "testdata" 目录中的预期快照进行比较，
// 使用测试名称和提供的 expectedNameSuffix 构造文件名。
//
// 如果快照不匹配，它会在 testing.TB 上报告错误。
func TestdataEqual(tb testing.TB, expectedNameSuffix string, actual Snapshotter, msgAndArgs ...any) {
	tb.Helper()
	testdataEq(tb, expectedNameSuffix, actual, func(expected, actual vttest.Snapshot) {
		assert.Equal(tb, expected, actual, msgAndArgs...)
	})
}

// TestdataRequireEqualf 将给定 [Snapshotter] 的快照与存储在 "testdata" 目录中的预期快照进行比较，
// 使用提供的格式和参数构造文件名。
//
// 如果快照不匹配，它会立即使测试失败。
func TestdataRequireEqualf(tb testing.TB, expectedNameSuffix string, actual Snapshotter, format string, args ...any) {
	tb.Helper()
	testdataEq(tb, expectedNameSuffix, actual, func(expected, actual vttest.Snapshot) {
		require.Equalf(tb, expected, actual, format, args...)
	})
}

// TestdataRequireEqual 将给定 [Snapshotter] 的快照与存储在 "testdata" 目录中的预期快照进行比较，
// 使用测试名称和提供的 expectedNameSuffix 构造文件名。
//
// 如果快照不匹配，它会立即使测试失败。
func TestdataRequireEqual(tb testing.TB, expectedNameSuffix string, actual Snapshotter, msgAndArgs ...any) {
	tb.Helper()
	testdataEq(tb, expectedNameSuffix, actual, func(expected, actual vttest.Snapshot) {
		require.Equal(tb, expected, actual, msgAndArgs...)
	})
}

func testdataEq(tb testing.TB, expectedNameSuffix string, actual Snapshotter, cb func(expected, actual vttest.Snapshot)) {
	tb.Helper()

	actualSnap := actual.Snapshot()
	fp := filepath.Join("testdata", fmt.Sprintf("%s_%s.json", tb.Name(), expectedNameSuffix))
	if *update {
		if err := os.MkdirAll(filepath.Dir(fp), 0o750); err != nil { //nolint: mnd
			tb.Fatal(err)
		}

		f, err := os.Create(fp)
		if err != nil {
			tb.Fatalf("创建快照文件失败: %v", err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(actualSnap); err != nil {
			tb.Fatalf("编码快照失败: %v", err)
		}

		if imgSnap, ok := actual.(Imager); ok {
			// 创建图像表示
			img := imgSnap.Image()
			fp := filepath.Join("testdata", fmt.Sprintf("%s_%s.png", tb.Name(), expectedNameSuffix))
			imgFile, err := os.Create(fp)
			if err != nil {
				tb.Fatalf("创建图像文件失败: %v", err)
			}
			defer imgFile.Close()

			if err := png.Encode(imgFile, img); err != nil {
				tb.Fatalf("编码图像失败: %v", err)
			}
		}
	}

	expectedSnapFile, err := os.Open(fp)
	if err != nil {
		tb.Fatalf("读取快照文件失败: %v", err)
	}

	var expectedSnap vttest.Snapshot
	if err := json.NewDecoder(expectedSnapFile).Decode(&expectedSnap); err != nil {
		tb.Fatalf("解码快照失败: %v", err)
	}

	cb(expectedSnap, actualSnap)
}
