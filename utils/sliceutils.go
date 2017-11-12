package utils

import "github.com/teolandon/hanoi/view/colors"

func NewPixelGrid(width int, height int) [][]colors.Pixel {
	ret := make([][]colors.Pixel, height)
	for i := range ret {
		ret[i] = make([]colors.Pixel, width)
	}
	return ret
}

func WritePixels(line []colors.Pixel, str string, color colors.Color, attr colors.Attribute) {
	runes := []rune(str)
	for i := range runes {
		if i >= len(line) {
			return
		}
		line[i] = colors.Pixel{runes[i], color, attr}
	}
}

type pxFunc func(i int) (c colors.Color, a colors.Attribute)

func WritePixelsWithFunc(line []colors.Pixel, str string, f pxFunc) {
	runes := []rune(str)
	for i := range runes {
		if i >= len(line) {
			return
		}

		color, attr := f(i)
		line[i] = colors.Pixel{runes[i], color, attr}
	}
}
