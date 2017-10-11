package view

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

type Area struct {
	s Size
	c Coords
}

func (a Area) Width() int {
	return a.s.Width
}

func (a Area) Height() int {
	return a.s.Height
}

func (a Area) X1() int {
	return a.c.X
}

func (a Area) Y1() int {
	return a.c.Y
}

func (a Area) X2() int {
	return a.c.X + a.Width()
}

func (a Area) Y2() int {
	return a.c.Y + a.Height()
}

type Coords struct {
	X int
	Y int
}

type Padding struct {
	Up    int
	Down  int
	Left  int
	Right int
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

var (
	defaultPalette = Palette{Black, Blue, Black, Yellow, Green, White,
		Magenta, Cyan, AttrDefault, AttrDefault, AttrDefault, AttrDefault}
	inheritAll = Palette{Inherit, Inherit, Inherit, Inherit, Inherit, Inherit,
		Inherit, Inherit, AttrInherit, AttrInherit, AttrInherit, AttrInherit}
)

type Displayable struct {
	Padding Padding
	Size    Size
	Coords  Coords
	Palette Palette
	Layout  Layout
}

type TitledContainer struct {
	Title           string
	TitleVisibility bool
	content         interface{}
	Displayable
}

func (TitledContainer) Parent() Focusable {
	return nil
}

func (t TitledContainer) Content() interface{} {
	return t.content
}

func (t *TitledContainer) SetContent(c interface{}) {
	t.content = c
}

func (t TitledContainer) DrawableArea() Area {
	offset := 0
	if t.TitleVisibility {
		offset = 2
	}
	return Area{Size{t.Size.Width - t.Padding.Left - t.Padding.Right - 2,
		t.Size.Height - t.Padding.Up - t.Padding.Down - 2 - offset},
		Coords{t.Coords.X + 1 + t.Padding.Left, t.Coords.Y + 1 + offset + t.Padding.Up}}
}

type TextBox struct {
	Text string
	Displayable
	parent Focusable
}

func (t TextBox) Parent() Focusable {
	return t.parent
}

func (TextBox) HandleKey(e KeyEvent) {
	if e.event.Ch == 'q' {
		exit()
		e.consumed = true
	}
}

func defaultDisplayable() Displayable {
	return Displayable{*new(Padding), *new(Size), *new(Coords), defaultPalette, FitToParent}
}

func displayableWithSize(size Size) Displayable {
	return Displayable{*new(Padding), size, *new(Coords), defaultPalette, FitToParent}
}

func NewTextBox(text string) *TextBox {
	ret := TextBox{text, displayableWithSize(Size{10, 5}), Parent{nil}}
	return &ret
}

func SimpleTitledContainer() TitledContainer {
	text := "Bigwordrighhere, butbigworderetoo, it, it has nice and small words, no long schlbberknockers to put you out of your lelelle"
	ret := TitledContainer{"Title", true, NewTextBox(text), displayableWithSize(Size{20, 10})}
	textb := ret.content.(*TextBox)
	textb.parent = Parent{ret}
	ret.Layout = Centered
	return ret
}
