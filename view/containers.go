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
	Content         interface{}
	Displayable
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
}

func defaultDisplayable() Displayable {
	return Displayable{Padding{1, 1, 1, 1}, *new(Size), *new(Coords), defaultPalette, FitToParent}
}

func displayableWithSize(size Size) Displayable {
	return Displayable{Padding{1, 1, 1, 1}, size, *new(Coords), defaultPalette, FitToParent}
}

func NewTextBox() *TextBox {
	ret := TextBox{"This is my content", displayableWithSize(Size{10, 5})}
	ret.Palette.normalBG = Black
	ret.Palette.normalFG = White
	return &ret
}

func SimpleTitledContainer() TitledContainer {
	ret := TitledContainer{"Title", true, NewTextBox(), displayableWithSize(Size{20, 10})}
	return ret
}
