package pixel

import "fmt"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"

type Pixel struct {
	Ch        rune
	Palette   colors.Palette
	Highlight colors.Highlight
}

type PixelGrid struct {
	width  int
	height int
	grid   [][]Pixel
}

func NewGrid(w, h int) PixelGrid {
	grid := plainGrid(w, h)
	ret := PixelGrid{w, h, grid}

	return ret
}

func (p PixelGrid) String() string {
	return fmt.Sprintf("[Grid of height %d, width %d]", p.height, p.width)
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
	log.Log("Returning line from pixel grid. Grid height:", len(p.grid))
	return p.grid[y]
}

func (p PixelGrid) Get(x, y int) Pixel {
	return p.grid[y][x]
}

func (pg PixelGrid) Set(x, y int, p Pixel) {
	pg.grid[y][x] = p
}

func (p PixelGrid) SubGrid(x1, x2, y1, y2 int) PixelGrid {
	if !p.Check() {
		panic("badbadbad")
	}
	sub := make([][]Pixel, y2-y1)
	copy(sub, p.grid[y1:y2]) // Essential, so as to not modify the original length
	for i := range sub {
		sub[i] = sub[i][x1:x2]
	}

	return PixelGrid{x2 - x1, y2 - y1, sub}
}

func (p PixelGrid) SubGridFromArea(a areas.Area) PixelGrid {
	return p.SubGrid(a.X1(), a.X2(), a.Y1(), a.Y2())
}

func (p PixelGrid) Padded(pad areas.Padding) PixelGrid {
	width := p.Width()
	height := p.Height()
	return p.SubGrid(pad.Left, width-pad.Right, pad.Up, height-pad.Down)
}

func (p PixelGrid) Check() bool {
	if len(p.grid) != p.height {
		log.Log("Height mismatched")
		return false
	}

	for i := range p.grid {
		if len(p.grid[i]) != p.width {
			log.Log("Some width mismatched", len(p.grid[i]), "instead of", p.width, "at", i)
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
