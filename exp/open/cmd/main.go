// Package main 演示 open 包的使用方法。
package main

import (
	"context"
	"fmt"

	"github.com/purpose168/charm-experimental-packages-cn/exp/open"
)

func main() {
	fmt.Println(open.Open(context.Background(), "https://charm.sh"))
}
