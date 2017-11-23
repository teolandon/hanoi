package utils

import "github.com/teolandon/hanoi/pixel"

func SubGrid(grid [][]pixel.Pixel, x1, x2, y1, y2 int) [][]pixel.Pixel {
	ret := grid[y1:y2]
	for i := range ret {
		ret[i] = ret[i][x1:x2]
	}
	return ret
}
