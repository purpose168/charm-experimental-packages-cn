package sixel

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestScanSize(t *testing.T) {
	testCases := map[string]struct {
		data           string
		expectedWidth  int
		expectedHeight int
	}{
		"two lines": {
			"~~~~~~-~~~~~~-", 6, 12,
		},
		"two lines no newline at end": {
			"~~~~~~-~~~~~~", 6, 12,
		},
		"no pixels": {
			"", 0, 0,
		},
		"smaller carriage returns": {
			"~$~~$~~~$~~~~$~~~~~$~~~~~~", 6, 6,
		},
		"transparent": {
			"??????", 6, 6,
		},
		"RLE": {
			"??!20?", 22, 6,
		},
		"Colors": {
			"#0;2;0;0;0~~~~~$#1;2;100;100;100;~~~~~~-#0~~~~~~-#1~~~~~~", 6, 18,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			decoder := &Decoder{}
			width, height := decoder.scanSize([]byte(testCase.data))
			if width != testCase.expectedWidth {
				t.Errorf("期望宽度为 %d，但接收到宽度为 %d", testCase.expectedWidth, width)
				return
			}

			if height != testCase.expectedHeight {
				t.Errorf("期望高度为 %d，但接收到高度为 %d", testCase.expectedHeight, height)
				return
			}
		})
	}
}

func TestFullImage(t *testing.T) {
	testCases := map[string]struct {
		imageWidth  int
		imageHeight int
		bandCount   int
		// 填充图像时，我们将使用索引到颜色的映射，并在当前索引在映射中时更改颜色
		// 这将防止连续出现多行相同颜色，使此测试稍微更易读
		colors map[int]color.RGBA
	}{
		"3x12 single color filled": {
			3, 12, 2,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
		},
		"3x12 two color filled": {
			3, 12, 2,
			map[int]color.RGBA{
				// 3 像素高的交替带
				0:  {0, 0, 255, 255},
				9:  {0, 255, 0, 255},
				18: {0, 0, 255, 255},
				27: {0, 255, 0, 255},
			},
		},
		"3x12 8 color with right gutter": {
			3, 12, 2,
			map[int]color.RGBA{
				0:  {255, 0, 0, 255},
				2:  {0, 255, 0, 255},
				3:  {255, 0, 0, 255},
				5:  {0, 255, 0, 255},
				6:  {255, 0, 0, 255},
				8:  {0, 255, 0, 255},
				9:  {0, 0, 255, 255},
				11: {128, 128, 0, 255},
				12: {0, 0, 255, 255},
				14: {128, 128, 0, 255},
				15: {0, 0, 255, 255},
				17: {128, 128, 0, 255},
				18: {0, 128, 128, 255},
				20: {128, 0, 128, 255},
				21: {0, 128, 128, 255},
				23: {128, 0, 128, 255},
				24: {0, 128, 128, 255},
				26: {128, 0, 128, 255},
				27: {64, 0, 0, 255},
				29: {0, 64, 0, 255},
				30: {64, 0, 0, 255},
				32: {0, 64, 0, 255},
				33: {64, 0, 0, 255},
				35: {0, 64, 0, 255},
			},
		},
		"3x12 single color with transparent band in the middle": {
			3, 12, 2,
			map[int]color.RGBA{
				0:  {255, 0, 0, 255},
				15: {0, 0, 0, 0},
				21: {255, 0, 0, 255},
			},
		},
		"3x5 single color": {
			3, 5, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
		},
		"12x4 single color use RLE": {
			12, 4, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
		},
		"12x1 two color use RLE": {
			12, 1, 1,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
				6: {0, 255, 0, 255},
			},
		},
		"12x12 single color use RLE": {
			12, 12, 2,
			map[int]color.RGBA{
				0: {255, 0, 0, 255},
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			img := image.NewRGBA(image.Rect(0, 0, testCase.imageWidth, testCase.imageHeight))

			currentColor := color.RGBA{0, 0, 0, 0}
			for y := range testCase.imageHeight {
				for x := range testCase.imageWidth {
					index := y*testCase.imageWidth + x
					newColor, changingColor := testCase.colors[index]
					if changingColor {
						currentColor = newColor
					}

					img.Set(x, y, currentColor)
				}
			}

			buffer := bytes.NewBuffer(nil)
			encoder := Encoder{}
			decoder := Decoder{}

			err := encoder.Encode(buffer, img)
			if err != nil {
				t.Errorf("意外错误: %+v", err)
				return
			}

			compareImg, err := decoder.Decode(buffer)
			if err != nil {
				t.Errorf("意外错误: %+v", err)
				return
			}

			expectedWidth := img.Bounds().Dx()
			expectedHeight := img.Bounds().Dy()
			actualWidth := compareImg.Bounds().Dx()
			actualHeight := compareImg.Bounds().Dy()

			if actualHeight != expectedHeight {
				t.Errorf("SixelImage 高度为 %d，但期望高度为 %d", actualHeight, expectedHeight)
				return
			}
			if actualWidth != expectedWidth {
				t.Errorf("SixelImage 宽度为 %d，但期望宽度为 %d", actualWidth, expectedWidth)
				return
			}

			for y := range expectedHeight {
				for x := range expectedWidth {
					r, g, b, a := compareImg.At(x, y).RGBA()
					expectedR, expectedG, expectedB, expectedA := img.At(x, y).RGBA()

					if r != expectedR || g != expectedG || b != expectedB || a != expectedA {
						t.Errorf("SixelImage 在坐标 (%d,%d) 处的颜色为 (%d,%d,%d,%d)，但期望颜色为 (%d,%d,%d,%d)",
							r, g, b, a, x, y, expectedR, expectedG, expectedB, expectedA)
						return
					}
				}
			}
		})
	}
}
