package sixel

import (
	"image"
	"image/color"
	"slices"
	"testing"
)

type testCase struct {
	maxColors       int
	expectedPalette []sixelColor
}

func TestPaletteCreationRedGreen(t *testing.T) {
	redGreen := image.NewRGBA(image.Rect(0, 0, 2, 2))
	redGreen.Set(0, 0, color.RGBA{255, 0, 0, 255})
	redGreen.Set(0, 1, color.RGBA{128, 0, 0, 255})
	redGreen.Set(1, 0, color.RGBA{0, 255, 0, 255})
	redGreen.Set(1, 1, color.RGBA{0, 128, 0, 255})

	testCases := map[string]testCase{
		"way too many colors": {
			maxColors: 16,
			expectedPalette: []sixelColor{
				{100, 0, 0, 100},
				{50, 0, 0, 100},
				{0, 100, 0, 100},
				{0, 50, 0, 100},
			},
		},
		"just the right number of colors": {
			maxColors: 4,
			expectedPalette: []sixelColor{
				{100, 0, 0, 100},
				{50, 0, 0, 100},
				{0, 100, 0, 100},
				{0, 50, 0, 100},
			},
		},
		"color reduction": {
			maxColors: 2,
			expectedPalette: []sixelColor{
				{75, 0, 0, 100},
				{0, 75, 0, 100},
			},
		},
	}

	runTests(t, redGreen, testCases)
}

func TestPaletteWithSemiTransparency(t *testing.T) {
	blueAlpha := image.NewRGBA(image.Rect(0, 0, 2, 2))
	blueAlpha.Set(0, 0, color.RGBA{0, 0, 255, 255})
	blueAlpha.Set(0, 1, color.RGBA{0, 0, 128, 255})
	blueAlpha.Set(1, 0, color.RGBA{0, 0, 255, 128})
	blueAlpha.Set(1, 1, color.RGBA{0, 0, 255, 0})

	testCases := map[string]testCase{
		"just the right number of colors": {
			maxColors: 4,
			expectedPalette: []sixelColor{
				{0, 0, 100, 100},
				{0, 0, 50, 100},
				{0, 0, 100, 50},
				{0, 0, 100, 0},
			},
		},
		"color reduction": {
			maxColors: 2,
			expectedPalette: []sixelColor{
				{0, 0, 75, 100},
				{0, 0, 100, 25},
			},
		},
	}
	runTests(t, blueAlpha, testCases)
}

func runTests(t *testing.T, img image.Image, testCases map[string]testCase) {
	for testName, test := range testCases {
		t.Run(testName, func(t *testing.T) {
			palette := newSixelPalette(img, test.maxColors)
			if len(palette.PaletteColors) != len(test.expectedPalette) {
				t.Errorf("期望调色板中有颜色 %+v，但得到 %+v", test.expectedPalette, palette.PaletteColors)
				return
			}

			for _, c := range test.expectedPalette {
				var foundColor bool
				if slices.Contains(palette.PaletteColors, c) {
					foundColor = true
				}

				if !foundColor {
					t.Errorf("期望调色板中有颜色 %+v，但得到 %+v", test.expectedPalette, palette.PaletteColors)
					return
				}
			}

			for lookupRawColor, lookupPaletteColor := range palette.colorConvert {
				paletteIndex, inReverseLookup := palette.paletteIndexes[lookupRawColor]
				if !inReverseLookup {
					t.Errorf("颜色 %+v 在 colorConvert 映射中映射到颜色 %+v，但 %+v 没有对应的调色板索引。", lookupRawColor, lookupPaletteColor, lookupPaletteColor)
					return
				}

				if paletteIndex >= len(palette.PaletteColors) {
					t.Errorf("图像颜色 %+v 映射到调色板索引 %d，但只有 %d 种调色板颜色。", lookupRawColor, paletteIndex, len(palette.PaletteColors))
					return
				}

				colorFromPalette := palette.PaletteColors[paletteIndex]
				if colorFromPalette != lookupPaletteColor {
					t.Errorf("图像颜色 %+v 映射到调色板颜色 %+v 和调色板索引 %d，但调色板索引 %d 实际上是调色板颜色 %+v", lookupRawColor, lookupPaletteColor, paletteIndex, paletteIndex, colorFromPalette)
					return
				}
			}
		})
	}
}
