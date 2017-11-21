package view

import "fmt"
import "github.com/teolandon/hanoi/utils"
import "github.com/teolandon/hanoi/view/colors"
import "github.com/teolandon/hanoi/view/colors/pixelwriter"

type Accent int8

const (
	normal      Accent = 0
	highlighted Accent = 2
	accent1     Accent = 4
	accent2     Accent = 6
)

// Layout is used as a field in elements to determine how they are
// positioned within a container.
type Layout uint8

// Relative parses the coordinates and size as percentages of the
// available area of the container.
//
// Absolute parses the coordinates as offsets from the top left corner
// of the available area of the container, and the size as absolute.
//
// FitToParent ignores the coordinates and size of the element, and
// instead assumes them to be the appropriate values to allow the element
// to fit its parent container.
//
// Centered ignored the coordinates, but takes the size into account to
// center the element in its parent container.
const (
	Relative Layout = iota
	Absolute
	FitToParent
	Centered
)

type Size struct {
	Width  int
	Height int
}

type area struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (a area) width() int {
	return a.x2 - a.x1
}

func (a area) height() int {
	return a.y2 - a.y1
}

func newArea(x, y int, size Size) area {
	return area{x, x + size.Width, y, y + size.Height}
}

type Padding struct {
	Up    int
	Down  int
	Left  int
	Right int
}

type TitledContainer struct {
	Title           string
	TitleVisibility bool
	content         Displayable
	displayable
}

func (t TitledContainer) String() string {
	return fmt.Sprintf("Titled container with title %s", t.Title)
}

func (t TitledContainer) Content() Displayable {
	return t.content
}

func (t *TitledContainer) SetContent(c Displayable) {
	t.content = c
	c.SetParent(t)
}

func (t TitledContainer) PixelGrid(a area) colors.PixelGrid {
	ret := utils.NewPixelGrid(a.width(), a.height())

	pw := pixelwriter.New(t.Palette(), colors.Normal, ret)

	pw.Write(0, 0, topLeftCorner)
	for i := range ret[0] {
		pw.Write(0, i, vLine)
	}
	pw.Write(0, len(ret[0])-1, topRightCorner)
	var i int
	for i = 1; i < len(ret)-1; i++ {
		pw.Write(i, 0, hLine)
		pw.Write(i, len(ret[i])-1, hLine)
	}
	pw.Write(i, 0, bottomLeftCorner)
	for j := range ret[0] {
		pw.Write(i, j, vLine)
	}
	pw.Write(i, len(ret[i])-1, bottomRightCorner)

	return ret
}

func (t TitledContainer) DrawableArea() (x, y int, s Size) {
	offset := 0
	if t.TitleVisibility {
		offset = 2
	}
	x = 1 + t.Padding().Left
	y = 1 + offset + t.Padding().Up
	s = Size{
		t.Size().Width - t.Padding().Left - t.Padding().Right - 2,
		t.Size().Height - t.Padding().Up - t.Padding().Down - 2 - offset,
	}
	return
}

type TextBox struct {
	Text string
	displayable
}

func (t TextBox) PixelGrid(a area) colors.PixelGrid {
	ret := utils.NewPixelGrid(a.width(), a.height())

	pw := pixelwriter.New(t.Palette(), colors.Normal, ret)

	// paintArea(parentArea, t.Palette().NormalFG, t.Palette().NormalBG)
	wrapped := utils.WrapText(t.Text, a.width(), a.height())
	for i, str := range wrapped {
		pw.WriteStr(a.x1, a.y1+i, str)
	}

	return ret
}

func NewTextBox(text string) TextBox {
	ret := TextBox{text, displayableWithSize(Size{10, 5})}
	return ret
}

func SimpleTitledContainer() TitledContainer {
	text := "Bigwordrighhere, butbigworderetoo, it, it has nice and small words, no long schlbberknockers to put you out of your lelelle"
	textBox := NewTextBox(text)
	textBox.SetLayout(FitToParent)
	ret := TitledContainer{"Title", true, nil, displayableWithSize(Size{20, 10})}
	ret.SetContent(&textBox)
	ret.SetLayout(Centered)
	return ret
}

func ButtonTitledContainer() *TitledContainer {
	ret := TitledContainer{"Test", true, nil, displayableWithSize(Size{20, 10})}
	button := NewButton("OK")
	ret.SetContent(&button)
	ret.SetLayout(Centered)
	return &ret
}
