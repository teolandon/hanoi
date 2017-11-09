package utils

func NewPixelGrid(width int, height int) [][]rune {
	ret := make([][]rune, height)
	for i := range ret {
		i := make([]rune, width)
	}
	return ret
}
