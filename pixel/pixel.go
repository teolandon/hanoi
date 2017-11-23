package pixel

import "fmt"
import "github.com/teolandon/hanoi/view/colors"

type Pixel struct {
	Ch        rune
	Palette   colors.Palette
	Highlight colors.Highlight
}

type PixelGrid struct {
	height int
	width  int
	grid   [][]Pixel
}

func NewGrid(h, w int) PixelGrid {
	grid := plainGrid(h, w)
	ret := PixelGrid{h, w, grid}

	return ret
}

func plainGrid(width int, height int) [][]Pixel {
	ret := make([][]Pixel, height)
	for i := range ret {
		ret[i] = make([]Pixel, width)
	}

	return ret
}

func (p PixelGrid) Height() int {
	return p.height
}

func (p PixelGrid) Width() int {
	return p.width
}

func (p PixelGrid) GetLine(y int) []Pixel {
	return p.grid[y]
}

func (p PixelGrid) Get(x, y int) Pixel {
	return p.grid[y][x]
}

func (pg PixelGrid) Set(x, y int, p Pixel) {
	pg.grid[y][x] = p
}

func (p PixelGrid) SubGrid(x1, x2, y1, y2 int) PixelGrid {
	sub := p.grid[y1:y2]
	for i := range sub {
		sub[i] = sub[i][x1:x2]
	}

	return PixelGrid{x2 - x1, y2 - y1, sub}
}

func (p PixelGrid) Check() bool {
	if len(p.grid) != p.height {
		return false
	}

	for i := range p.grid {
		if len(p.grid[i]) != p.width {
			return false
		}
	}

	return true
}

type PixelWriter struct {
	defaultPalette   colors.Palette
	defaultHighlight colors.Highlight
	grid             PixelGrid
}

func NewWriter(dP colors.Palette, dH colors.Highlight, grid PixelGrid) PixelWriter {
	return PixelWriter{dP, dH, grid}
}

func (p PixelWriter) Write(x, y int, r rune) error {
	return p.WriteWithHighlight(x, y, r, p.defaultHighlight)
}

func (p PixelWriter) WriteWithHighlight(x, y int, r rune, hi colors.Highlight) error {
	grid := p.grid.grid

	if y >= len(grid) {
		return fmt.Errorf("y-value out of bounds for PixelWriter %v's grid:\ny=%v,"+
			"len(grid)=%v", p, y, len(grid))
	}

	if x >= len(grid[y]) {
		return fmt.Errorf("x-value out of bounds for PixelWriter %v's grid:\nx=%v,"+
			"len(grid[%v])=%v", p, y, x, len(grid))
	}

	grid[y][x] = Pixel{r, p.defaultPalette, hi}

	return nil
}

func (p PixelWriter) WriteStr(x, y int, str string) {
	runes := []rune(str)
	for i := range runes {
		p.Write(x, y+i, runes[i])
	}
}
