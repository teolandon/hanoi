package containers

import tb "github.com/nsf/termbox-go"

type accent int8

const (
	normal  accent = 0
	accent1 accent = 2
	accent2 accent = 4
)

type Pixel struct {
	ch rune
	accent
}

type PixelGrid [][]Pixel

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
	padding() Padding
	bounds() Rect
	grid() PixelGrid
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

func (c TitledContainer) Resize(newBounds Rect) {
	c.content.Resize(newBounds.resizeByPadding(c.padding))
	c.bounds = newBounds
}
