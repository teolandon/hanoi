package utils

import "github.com/teolandon/hanoi/view/colors"

func NewPixelGrid(width int, height int) [][]colors.Pixel {
	ret := make([][]colors.Pixel, height)
	for i := range ret {
		ret[i] = make([]colors.Pixel, width)
	}
	return ret
}
