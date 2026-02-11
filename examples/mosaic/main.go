// 包 main 演示了使用方法。
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/purpose168/charm-experimental-packages-cn/mosaic"
	lipgloss "github.com/purpose168/lipgloss-cn"
)

func main() {
	dogImg, err := loadImage("./pekinas.jpg")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	m := mosaic.New().Width(80).Height(40)

	fmt.Println(lipgloss.JoinVertical(lipgloss.Right, lipgloss.JoinHorizontal(lipgloss.Center, m.Render(dogImg))))
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	defer f.Close() //nolint:errcheck,staticcheck
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return jpeg.Decode(f) //nolint:wrapcheck
}
