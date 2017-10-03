package view

import "fmt"

type Accent int8

const (
	normal      Accent = 0
	highlighted Accent = 2
	accent1     Accent = 4
	accent2     Accent = 6
)

type Pixel struct {
	Ch rune
	Accent
}

type PixelGrid [][]Pixel

func newPixelGrid(b Rect) PixelGrid {
	fmt.Println(b.Y2, b.Y1)
	retGrid := make([][]Pixel, b.Y2-b.Y1)
	for i := range retGrid {
		retGrid[i] = make([]Pixel, b.X2-b.X1)
	}
	return retGrid
}

type Padding struct {
	Up    int
	Down  int
	Left  int
	Right int
}

type Rect struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

func (r Rect) resizeByPadding(p Padding) (ret Rect) {
	ret.X1 = r.X1 + p.Left
	ret.X2 = r.X2 - p.Right
	ret.Y1 = r.Y1 + p.Up
	ret.Y2 = r.Y2 - p.Down
	return
}

func (r Rect) resizeAllBy(i int) (ret Rect) {
	ret.X1 = r.X1 - i
	ret.X2 = r.X2 + i
	ret.Y1 = r.Y1 - i
	ret.Y2 = r.Y2 + i
	return
}

type Palette struct {
	normalFG      Color
	normalBG      Color
	highlightFG   Color
	highlightBG   Color
	accent1FG     Color
	accent1BG     Color
	accent2FG     Color
	accent2BG     Color
	normalAttr    Attribute
	highlightAttr Attribute
	accent1Attr   Attribute
	accent2Attr   Attribute
}

type Displayable struct {
	Padding Padding
	Bounds  Rect
	Grid    PixelGrid
	Palette Palette
}

// resize() should call content().resize() or similar
type TitledContainer struct {
	Title           string
	TitleVisibility bool
	Content         Containable
	bounds          Rect
	padding         Padding
}

func (c TitledContainer) Bounds() Rect {
	return c.bounds
}

func (c TitledContainer) Padding() Padding {
	return c.padding
}

func (c TitledContainer) Grid() PixelGrid {
	retGrid := newPixelGrid(c.bounds)
	fmt.Println(c.bounds)
	fmt.Println(len(retGrid))
	contentGrid := c.Content.Grid()
	for i := 0; i < c.bounds.X2-c.bounds.X1-(c.padding.Right+c.padding.Left); i++ {
		for j := 0; j < c.bounds.Y2-c.bounds.Y1-(c.padding.Down+c.padding.Up); j++ {
			fmt.Println("i:", i, "\nj:", j)
			retGrid[i+c.padding.Left-1][j+c.padding.Up-1] = contentGrid[i][j]
		}
	}
	return retGrid
}

type TextBox struct {
	Text    string
	bounds  Rect
	palette Palette
}

func (t TextBox) Grid() PixelGrid {
	retGrid := newPixelGrid(t.bounds)
	if len(t.Text) > 10 {
		panic("LOLOL")
	}
	for i, s := range t.Text {
		retGrid[0][i] = Pixel{s, normal}
	}
	return retGrid
}

func (t TextBox) SetPadding(p Padding) {}

func (t TextBox) Bounds() Rect {
	return t.bounds
}

func (t *TextBox) Resize(newBounds Rect) {
	t.bounds = newBounds
}

func (t TextBox) Padding() Padding {
	return Padding{0, 0, 0, 0}
}

func (t TextBox) Palette() Palette {
	return t.palette
}

func (t *TextBox) SetPalette(p Palette) {
	t.palette = p
}

func NewTextBox() *TextBox {
	ret := TextBox{"lololo", Rect{3, 19, 3, 19}, defaultPalette}
	return &ret
}

func SimpleTitledContainer() TitledContainer {
	ret := TitledContainer{"Lol", true, NewTextBox(), Rect{1, 20, 1, 20}, Padding{3, 3, 3, 3}}
	return ret
}
