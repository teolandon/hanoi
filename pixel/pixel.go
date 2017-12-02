package pixel

import "fmt"
import "errors"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"

type Pixel struct {
	Ch        rune
	Palette   colors.Palette
	Highlight colors.Highlight
}

type SquareGrid struct {
	width  int
	height int
	grid   [][]Pixel
}

func plainGrid(width int, height int) [][]Pixel {
	ret := make([][]Pixel, height)
	for i := range ret {
		ret[i] = make([]Pixel, width)
	}

	return ret
}

func NewGrid(width, height int) SquareGrid {
	return SquareGrid{width, height, plainGrid(width, height)}
}

func (p SquareGrid) Height() int {
	return p.height
}

func (p SquareGrid) Width() int {
	return p.width
}

func (p SquareGrid) GetLine(y int) []Pixel {
	log.Log("Returning line from pixel grid. Grid height:", len(p.grid))
	if y >= p.Height() || y < 0 {
		return nil
	}

	return p.grid[y]
}

func (p SquareGrid) Get(x, y int) Pixel {
	if x >= p.Width() || x < 0 {
		return *new(Pixel)
	}
	if y >= p.Height() || y < 0 {
		return *new(Pixel)
	}

	return p.grid[y][x]
}

func (pg SquareGrid) Set(x, y int, p Pixel) {
	if x >= pg.Width() || x < 0 {
		return
	}
	if y >= pg.Height() || y < 0 {
		return
	}

	pg.grid[y][x] = p
}

func (p SquareGrid) SubGrid(area areas.Area) SubGrid {
	return SubGrid{area, p}
}

func (p SquareGrid) TotalSubGrid() SubGrid {
	return p.SubGrid(areas.New(0, p.Width(), 0, p.Height()))
}

type SubGrid struct {
	area areas.Area
	grid SquareGrid
}

func (p SubGrid) String() string {
	return fmt.Sprintf("[Subgrid of area %v]", p.area)
}

func (p SubGrid) SubGrid(a areas.Area) SubGrid {
	return SubGrid{p.area.SubArea(a), p.grid}
}

func (p SubGrid) Padded(pad areas.Padding) SubGrid {
	return p.SubGrid(p.area.Padded(pad))
}

func (p SubGrid) Width() int {
	return p.area.Width()
}

func (p SubGrid) Height() int {
	return p.area.Height()
}

func (p SubGrid) Get(x, y int) (pixel Pixel, err error) {
	gx := x + p.area.X1()
	gy := y + p.area.Y1()

	return p.grid.Get(gx, gy), nil
}

func (p SubGrid) GetLine(y int) []Pixel {
	gy := y + p.area.Y1()
	ret := p.grid.GetLine(gy)
	return ret
}

func (p SubGrid) Set(x, y int, pixel Pixel) error {
	if x >= p.Width() || x < 0 {
		return errors.New("X out of bounds")
	}
	if y >= p.Height() || y < 0 {
		return errors.New("Y out of bounds")
	}

	gx := x + p.area.X1()
	gy := y + p.area.Y1()

	p.grid.Set(gx, gy, pixel)

	return nil
}

type PixelWriter struct {
	defaultPalette   colors.Palette
	defaultHighlight colors.Highlight
	grid             SubGrid
}

func NewWriter(dP colors.Palette, dH colors.Highlight, grid SubGrid) PixelWriter {
	return PixelWriter{dP, dH, grid}
}

func (p PixelWriter) Write(x, y int, r rune) error {
	return p.WriteWithHighlight(x, y, r, p.defaultHighlight)
}

func (p PixelWriter) WriteWithHighlight(x, y int, r rune, hi colors.Highlight) error {
	grid := p.grid.grid

	if y >= grid.Height() {
		return fmt.Errorf("y-value out of bounds for PixelWriter %v's grid:\ny=%v,"+
			"len(grid)=%v", p, y, grid.Height())
	}

	if x >= grid.Width() {
		return fmt.Errorf("x-value out of bounds for PixelWriter %v's grid:\nx=%v,"+
			"len(grid[%v])=%v", p, y, x, grid.Width())
	}

	grid.Set(x, y, Pixel{r, p.defaultPalette, hi})

	return nil
}

func (p PixelWriter) WriteStr(x, y int, str string) {
	runes := []rune(str)
	for i := range runes {
		p.Write(x, y+i, runes[i])
	}
}
