package containers

import tb "github.com/nsf/termbox-go"

type Accent int8

const (
	normal  Accent = 0
	accent1 Accent = 2
	accent2 Accent = 4
)

type Pixel struct {
	Ch rune
	Accent
}

type PixelGrid [][]Pixel

func newPixelGrid(b Rect) PixelGrid {
	retGrid := make([][]Pixel, b.y2-b.y1)
	for i := range retGrid {
		retGrid[i] = make([]Pixel, b.x2-b.x1)
	}
	return retGrid
}

type Padding struct {
	up    int
	down  int
	left  int
	right int
}

type Rect struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (r Rect) resizeByPadding(p Padding) (ret Rect) {
	ret.x1 = r.x1 + p.left
	ret.x2 = r.x2 - p.right
	ret.y1 = r.y1 + p.up
	ret.y2 = r.y2 - p.down
	return
}

func (r Rect) resizeAllBy(i int) (ret Rect) {
	ret.x1 = r.x1 - i
	ret.x2 = r.x2 + i
	ret.y1 = r.y1 - i
	ret.y2 = r.y2 + i
	return
}

type palette [6]tb.Attribute

func newPalette(fg, bg, accent1fg, accent1bg, accent2fg, accent2bg tb.Attribute) palette {
	return [6]tb.Attribute{fg, bg, accent1fg, accent1bg, accent2fg, accent2bg}
}

type Resizable interface {
	Resize(newBounds Rect)
}

type Displayable interface {
	Padding() Padding
	Bounds() Rect
	Grid() PixelGrid
}

type Containable interface {
	Resizable
	Displayable
}

// resize() should call content().resize() or similar
type TitledContainer struct {
	title           string
	titleVisibility bool
	content         Containable
	bounds          Rect
	padding         Padding
}

func (c *TitledContainer) Resize(newBounds Rect) {
	c.content.Resize(newBounds.resizeByPadding(c.padding))
	c.bounds = newBounds
}

func (c TitledContainer) Bounds() Rect {
	return c.bounds
}

func (c TitledContainer) Padding() Padding {
	return c.padding
}

func (c TitledContainer) Grid() PixelGrid {
	retGrid := make([][]Pixel, c.bounds.y2-c.bounds.y1)
	for i := range retGrid {
		retGrid[i] = make([]Pixel, c.bounds.x2-c.bounds.x1)
	}
	contentGrid := c.content.Grid()
	for i := 0; i < c.bounds.x2-c.bounds.x1-(c.padding.right+c.padding.left); i++ {
		for j := 0; j < c.bounds.y2-c.bounds.y1-(c.padding.down+c.padding.up); j++ {
			retGrid[i+c.padding.left][j+c.padding.up] = contentGrid[i][j]
		}
	}
	return retGrid
}

type TextBox struct {
	text   string
	bounds Rect
}

func (t TextBox) Grid() PixelGrid {
	retGrid := newPixelGrid(t.bounds)
	if len(t.text) > 10 {
		panic("LOLOL")
	}
	for i, s := range t.text {
		retGrid[0][i] = Pixel{s, normal}
	}
	return retGrid
}

func (t TextBox) Bounds() Rect {
	return t.bounds
}

func (t *TextBox) Resize(newBounds Rect) {
	t.bounds = newBounds
}

func (t TextBox) Padding() Padding {
	return Padding{0, 0, 0, 0}
}

func NewTextBox() TextBox {
	return TextBox{"lololo", Rect{1, 20, 1, 20}}
}
