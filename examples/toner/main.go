// Package main 演示使用方法。
package main

import (
	"io"
	"log"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/exp/toner"
)

func main() {
	bts, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("从标准输入读取失败: %v", err)
	}

	w := toner.Writer{Writer: os.Stdout}
	if _, err := w.Write(bts); err != nil {
		log.Fatalf("写入标准输出失败: %v", err)
	}
}
