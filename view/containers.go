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
const (
	Relative Layout = iota
	Absolute
	FitToParent
)

type Size struct {
	Width  int
	Height int
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
	ret := TextBox{"lololo", defaultDisplayable()}
	return &ret
}

func SimpleTitledContainer() TitledContainer {
	ret := TitledContainer{"Lol", true, NewTextBox(), displayableWithSize(Size{20, 10})}
	return ret
}
