package pixelwriter

import "fmt"
import "github.com/teolandon/hanoi/view/colors"

type PixelWriter struct {
	defaultPalette   colors.Palette
	defaultHighlight colors.Highlight
	grid             colors.PixelGrid
}

func New(dP colors.Palette, dH colors.Highlight, grid colors.PixelGrid) PixelWriter {
	return PixelWriter{dP, dH, grid}
}

func (p PixelWriter) Write(x, y int, r rune) error {
	if y >= len(p.grid) {
		return fmt.Errorf("y-value out of bounds for PixelWriter %v's grid:\ny=%v,"+
			"len(grid)=%v", p, y, len(p.grid))
	}

	if x >= len(p.grid[y]) {
		return fmt.Errorf("x-value out of bounds for PixelWriter %v's grid:\nx=%v,"+
			"len(grid[%v])=%v", p, y, x, len(p.grid))
	}

	p.grid[y][x] = colors.Pixel{r, p.defaultPalette, p.defaultHighlight}

	return nil
}

func (p PixelWriter) WriteStr(x, y int, str string) {
	runes := []rune(str)
	for i := range runes {
		p.Write(x, y+i, runes[i])
	}
}
