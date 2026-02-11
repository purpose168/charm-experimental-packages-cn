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
		log.Fatalf("failed to read from stdin: %v", err)
	}

	w := toner.Writer{Writer: os.Stdout}
	if _, err := w.Write(bts); err != nil {
		log.Fatalf("failed to write to stdout: %v", err)
	}
}
