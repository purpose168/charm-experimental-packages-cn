// 包 main 演示了使用方法。
package main

import (
	"io"
	"log"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/ansi"
	"github.com/purpose168/charm-experimental-packages-cn/cellbuf"
)

func main() {
	pw := cellbuf.NewPenWriter(os.Stdout)
	defer pw.Close() //nolint:errcheck

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(pw, ansi.Wrap(string(data), 10, "")) //nolint:errcheck,gosec
}
