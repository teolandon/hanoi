package view

import "fmt"
import "github.com/teolandon/hanoi/utils"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import _ "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"

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

type TitledContainer struct {
	Title           string
	TitleVisibility bool
	container
}

func (t *TitledContainer) setGrid(g pixel.SubGrid) {
	t.container.setGrid(g, t.ContentPadding())
}

func (t TitledContainer) ContentPadding() areas.Padding {
	pad := 0

	if t.TitleVisibility {
		pad = 2
	}

	return areas.Padding{1, 1, 1 + pad, 1}
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

func (t TitledContainer) Draw() {
	pw := pixel.NewWriter(t.Palette(), colors.Normal, t.grid)

	pw.Write(0, 0, topLeftCorner)
	for i := range t.grid.GetLine(0) {
		pw.Write(0, i, vLine)
	}
	pw.Write(0, t.grid.Height()-1, bottomLeftCorner)
	var i int
	for i = 1; i < t.grid.Height()-1; i++ {
		pw.Write(i, 0, hLine)
		pw.Write(i, t.grid.Width()-1, hLine)
	}
	pw.Write(i, 0, topRightCorner)
	for j := range t.grid.GetLine(0) {
		pw.Write(i, j, vLine)
	}
	pw.Write(i, t.grid.Width()-1, bottomRightCorner)
}

func (t TitledContainer) DrawableArea() (x, y int, s areas.Size) {
	offset := 0
	if t.TitleVisibility {
		offset = 2
	}
	x = 1 + t.Padding().Left
	y = 1 + offset + t.Padding().Up
	s = areas.NewSize(
		t.Size().Width()-t.Padding().Left-t.Padding().Right-2,
		t.Size().Height()-t.Padding().Up-t.Padding().Down-2-offset,
	)
	return
}

type TextBox struct {
	Text string
	displayable
}

func (t TextBox) Draw() {
	pw := pixel.NewWriter(t.Palette(), colors.Normal, t.grid)

	// paintArea(parentArea, t.Palette().NormalFG, t.Palette().NormalBG)
	wrapped := utils.WrapText(t.Text, t.grid.Width(), t.grid.Height())
	for i, str := range wrapped {
		pw.WriteStr(0, i, str)
	}
}

func NewTextBox(text string) TextBox {
	ret := TextBox{text, displayableWithSize(areas.NewSize(10, 5))}
	return ret
}

func SimpleTitledContainer() TitledContainer {
	text := "Bigwordrighhere, butbigworderetoo, it, it has nice and small words, no long schlbberknockers to put you out of your lelelle"
	textBox := NewTextBox(text)
	textBox.SetLayout(FitToParent)
	ret := TitledContainer{"Title", true, container{displayableWithSize(areas.NewSize(20, 10)), nil}}
	ret.SetContent(&textBox)
	ret.SetLayout(Centered)
	return ret
}

func ButtonTitledContainer() *TitledContainer {
	ret := TitledContainer{"Test", true, container{displayableWithSize(areas.NewSize(20, 10)), nil}}
	button := NewButton("OK")
	ret.SetContent(&button)
	ret.SetLayout(Centered)
	return &ret
}
