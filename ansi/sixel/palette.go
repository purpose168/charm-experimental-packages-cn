package sixel

import (
	"container/heap"
	"image"
	"image/color"
	"math"
)

// sixelPalette 是一个最多包含 256 种颜色的调色板，列出了 SixelImage 将使用的颜色。
// 大多数图像，特别是 jpeg，有超过 256 种颜色，因此创建 sixelPalette 需要使用颜色量化。
// 为此，我们使用中位切割算法。
//
// 中位切割要求将图像中的所有像素定位在 4D 颜色立方体中，每个通道一个轴。
// 立方体沿其最长轴切成两半，使得立方体中的一半像素最终在一个子立方体中，
// 另一半在另一个子立方体中。我们继续沿最长轴将立方体切成两半，
// 直到有 256 个子立方体。然后，每个子立方体中所有像素的平均值用作该立方体的颜色。
//
// 颜色根据它们最接近的颜色转换为调色板颜色（不一定是它们立方体的颜色）。
//
// 此实现与官方算法有一些小的（但似乎非常常见的）差异：
//   - 在确定最长轴时，立方体中的像素数乘以轴长度
//     这在图像中有大量空间被相同颜色的不同色调占据的情况下，显著改善了颜色选择。
//   - 如果单个颜色位于切割线上，则该颜色的所有像素都分配给一个子立方体
//     而不是尝试在子立方体之间拆分它们。这允许我们使用唯一颜色的切片和像素计数映射，
//     而不是尝试单独表示每个像素。
type sixelPalette struct {
	// 用于将颜色从图像转换为调色板颜色的映射
	colorConvert map[sixelColor]sixelColor
	// 从图像颜色获取调色板索引的查找
	paletteIndexes map[sixelColor]int
	PaletteColors  []sixelColor
}

// quantizationChannel 是枚举类型，指示颜色立方体中的轴。用于指示
// 立方体中的哪个轴最长。
type quantizationChannel int

const (
	// MaxColors 是 sixelPalette 可以包含的最大颜色数。
	MaxColors int = 256

	quantizationRed quantizationChannel = iota
	quantizationGreen
	quantizationBlue
	quantizationAlpha
)

// quantizationCube 表示中位切割算法中的单个立方体。
type quantizationCube struct {
	// startIndex 是此立方体在 uniqueColors 切片中开始的索引
	startIndex int
	// length 是此立方体在 uniqueColors 切片中占用的元素数
	length int
	// sliceChannel 是如果此立方体被切成两半将被切割的轴
	sliceChannel quantizationChannel
	// score 是启发式值：越高表示此立方体更有可能被切割
	score uint64
	// pixelCount 是此立方体中包含的像素数
	pixelCount uint64
}

// cubePriorityQueue 是一个堆，用于对 quantizationCube 对象进行排序，以选择正确的
// 一个进行下一次切割。Pop 将移除具有最高分数的队列。
type cubePriorityQueue []any

func (p *cubePriorityQueue) Push(x any) {
	*p = append(*p, x)
}

func (p *cubePriorityQueue) Pop() any {
	popped := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return popped
}

func (p *cubePriorityQueue) Len() int {
	return len(*p)
}

func (p *cubePriorityQueue) Less(i, j int) bool {
	left := (*p)[i].(quantizationCube)
	right := (*p)[j].(quantizationCube)

	// 我们想要最大的通道方差
	return left.score > right.score
}

func (p *cubePriorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// createCube 用于初始化包含 uniqueColors 切片区域的新 quantizationCube。
func (p *sixelPalette) createCube(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, startIndex, bucketLength int) quantizationCube {
	minRed, minGreen, minBlue, minAlpha := uint32(0xffff), uint32(0xffff), uint32(0xffff), uint32(0xffff)
	maxRed, maxGreen, maxBlue, maxAlpha := uint32(0), uint32(0), uint32(0), uint32(0)
	totalWeight := uint64(0)

	// 找出哪个通道具有最大的方差
	for i := startIndex; i < startIndex+bucketLength; i++ {
		r, g, b, a := uniqueColors[i].Red, uniqueColors[i].Green, uniqueColors[i].Blue, uniqueColors[i].Alpha
		totalWeight += pixelCounts[uniqueColors[i]]

		if r < minRed {
			minRed = r
		}
		if r > maxRed {
			maxRed = r
		}
		if g < minGreen {
			minGreen = g
		}
		if g > maxGreen {
			maxGreen = g
		}
		if b < minBlue {
			minBlue = b
		}
		if b > maxBlue {
			maxBlue = b
		}
		if a < minAlpha {
			minAlpha = a
		}
		if a > maxAlpha {
			maxAlpha = a
		}
	}

	dRed := maxRed - minRed
	dGreen := maxGreen - minGreen
	dBlue := maxBlue - minBlue
	dAlpha := maxAlpha - minAlpha

	cube := quantizationCube{
		startIndex: startIndex,
		length:     bucketLength,
		pixelCount: totalWeight,
	}

	if dRed >= dGreen && dRed >= dBlue && dRed >= dAlpha {
		cube.sliceChannel = quantizationRed
		cube.score = uint64(dRed)
	} else if dGreen >= dBlue && dGreen >= dAlpha {
		cube.sliceChannel = quantizationGreen
		cube.score = uint64(dGreen)
	} else if dBlue >= dAlpha {
		cube.sliceChannel = quantizationBlue
		cube.score = uint64(dBlue)
	} else {
		cube.sliceChannel = quantizationAlpha
		cube.score = uint64(dAlpha)
	}

	// 提高包含更多像素的立方体的分数
	cube.score *= totalWeight

	return cube
}

// quantize 是一个方法，将初始化调色板的颜色和查找，提供一组
// 唯一颜色和包含这些颜色的像素计数的映射。
func (p *sixelPalette) quantize(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, maxColors int) {
	p.colorConvert = make(map[sixelColor]sixelColor)
	p.paletteIndexes = make(map[sixelColor]int)

	// 如果我们甚至没有超过最大颜色数，则不需要量化，实际上，如果我们有少于最大颜色数，此代码会爆炸
	if len(uniqueColors) <= maxColors {
		p.PaletteColors = uniqueColors
		return
	}

	cubeHeap := make(cubePriorityQueue, 0, maxColors)

	// 从包含所有颜色的立方体开始
	heap.Init(&cubeHeap)
	heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, 0, len(uniqueColors)))

	// 将最好的立方体切成两个立方体，直到我们有最大颜色数，然后我们有了调色板
	for cubeHeap.Len() < maxColors {
		cubeToSplit := heap.Pop(&cubeHeap).(quantizationCube)

		//nolint:godox
		// TODO: 将来使用 slices.SortFunc 和 cmp.Compare (>=1.24)
		// 然后可以删除 palette_sort.go
		sortFunc(uniqueColors[cubeToSplit.startIndex:cubeToSplit.startIndex+cubeToSplit.length],
			func(left sixelColor, right sixelColor) int {
				switch cubeToSplit.sliceChannel { //nolint:exhaustive // alpha channel not used
				case quantizationRed:
					return compare(left.Red, right.Red)
				case quantizationGreen:
					return compare(left.Green, right.Green)
				case quantizationBlue:
					return compare(left.Blue, right.Blue)
				default:
					return compare(left.Alpha, right.Alpha)
				}
			})

		// 我们需要拆分此立方体中的颜色，使得像素在两个立方体之间均匀分配，
		// 或者至少尽可能接近。我们所做的是在遍历时计算像素，
		// 并在大约一半的像素在左侧时放置切割点
		countSoFar := pixelCounts[uniqueColors[cubeToSplit.startIndex]]
		targetCount := cubeToSplit.pixelCount / 2
		leftLength := 1

		for i := cubeToSplit.startIndex + 1; i < cubeToSplit.startIndex+cubeToSplit.length; i++ {
			c := uniqueColors[i]
			weight := pixelCounts[c]
			if countSoFar+weight > targetCount {
				break
			}
			leftLength++
			countSoFar += weight
		}

		rightLength := cubeToSplit.length - leftLength
		rightIndex := cubeToSplit.startIndex + leftLength
		heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, cubeToSplit.startIndex, leftLength))
		heap.Push(&cubeHeap, p.createCube(uniqueColors, pixelCounts, rightIndex, rightLength))
	}

	// 一旦我们在堆中有最大立方体，就将它们全部取出并加载到调色板中
	for cubeHeap.Len() > 0 {
		bucketToLoad := heap.Pop(&cubeHeap).(quantizationCube)
		p.loadColor(uniqueColors, pixelCounts, bucketToLoad.startIndex, bucketToLoad.length)
	}
}

// ColorIndex 接受原始图像颜色（不是调色板颜色）并提供该颜色的调色板索引。
func (p *sixelPalette) ColorIndex(c sixelColor) int {
	return p.paletteIndexes[c]
}

// loadColor 接受表示单个中位切割立方体的颜色范围。它计算
// 立方体中的平均颜色并将其添加到调色板中。
func (p *sixelPalette) loadColor(uniqueColors []sixelColor, pixelCounts map[sixelColor]uint64, startIndex, cubeLen int) {
	totalRed, totalGreen, totalBlue, totalAlpha := uint64(0), uint64(0), uint64(0), uint64(0)
	totalCount := uint64(0)
	for i := startIndex; i < startIndex+cubeLen; i++ {
		count := pixelCounts[uniqueColors[i]]
		totalRed += uint64(uniqueColors[i].Red) * count
		totalGreen += uint64(uniqueColors[i].Green) * count
		totalBlue += uint64(uniqueColors[i].Blue) * count
		totalAlpha += uint64(uniqueColors[i].Alpha) * count
		totalCount += count
	}

	averageColor := sixelColor{
		Red:   uint32(totalRed / totalCount),   //nolint:gosec
		Green: uint32(totalGreen / totalCount), //nolint:gosec
		Blue:  uint32(totalBlue / totalCount),  //nolint:gosec
		Alpha: uint32(totalAlpha / totalCount), //nolint:gosec
	}

	p.PaletteColors = append(p.PaletteColors, averageColor)
}

// sixelColor 是一个扁平结构，包含单个颜色：所有通道都是 0-100
// 而不是任何合理的值。
type sixelColor struct {
	Red   uint32
	Green uint32
	Blue  uint32
	Alpha uint32
}

// sixelConvertColor 接受普通的 Go 颜色并将其转换为 sixelColor，其
// 通道范围从 0-100。
func sixelConvertColor(c color.Color) sixelColor {
	r, g, b, a := c.RGBA()
	return sixelColor{
		Red:   sixelConvertChannel(r),
		Green: sixelConvertChannel(g),
		Blue:  sixelConvertChannel(b),
		Alpha: sixelConvertChannel(a),
	}
}

// sixelConvertChannel 将单个颜色通道从 go 的标准 0-0xffff 范围转换为
// sixel 的 0-100 范围。
func sixelConvertChannel(channel uint32) uint32 {
	// 我们加 328 是因为那大约是 sixel 0-100 颜色范围内的 0.5，我们试图
	// 四舍五入到最近的值
	return (channel + 328) * 100 / 0xffff
}

// newSixelPalette 接受一个图像并使用中位切割算法生成 N 颜色量化调色板。
// 生成的 sixelPalette 可以在 O(1) 时间内将颜色从图像转换为量化调色板。
func newSixelPalette(image image.Image, maxColors int) sixelPalette {
	pixelCounts := make(map[sixelColor]uint64)

	height := image.Bounds().Dy()
	width := image.Bounds().Dx()

	// 记录每种颜色的像素计数，同时获取图像中所有唯一颜色的集合
	for y := range height {
		for x := range width {
			c := sixelConvertColor(image.At(x, y))
			count := pixelCounts[c]
			count++

			pixelCounts[c] = count
		}
	}

	p := sixelPalette{}
	uniqueColors := make([]sixelColor, 0, len(pixelCounts))
	for c := range pixelCounts {
		uniqueColors = append(uniqueColors, c)
	}

	// 使用中位切割算法构建 p.PaletteColors
	p.quantize(uniqueColors, pixelCounts, maxColors)

	// 立方体中颜色的平均值并不总是最近的调色板颜色。因此，
	// 我们需要使用这个非常令人不安的双重循环来查找图像中每个
	// 唯一颜色的查找调色板颜色。
	for _, c := range uniqueColors {
		var bestColor sixelColor
		var bestColorIndex int
		bestScore := uint32(math.MaxUint32)

		for paletteIndex, paletteColor := range p.PaletteColors {
			redDiff := c.Red - paletteColor.Red
			greenDiff := c.Green - paletteColor.Green
			blueDiff := c.Blue - paletteColor.Blue
			alphaDiff := c.Alpha - paletteColor.Alpha

			score := (redDiff * redDiff) + (greenDiff * greenDiff) + (blueDiff * blueDiff) + (alphaDiff * alphaDiff)
			if score < bestScore {
				bestColor = paletteColor
				bestColorIndex = paletteIndex
				bestScore = score
			}
		}

		p.paletteIndexes[c] = bestColorIndex
		p.colorConvert[c] = bestColor
	}

	return p
}
