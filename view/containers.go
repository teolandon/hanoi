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

type area struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (a area) Width() int {
	return a.y2 - a.y1
}

func (a area) Height() int {
	return a.x2 - a.x1
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

type TitledContainer struct {
	Title           string
	TitleVisibility bool
	content         interface{}
	Displayable
}

func (t TitledContainer) Content() interface{} {
	return t.content
}

func (t *TitledContainer) SetContent(c interface{}) {
	t.content = c
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
	Displayable
}

func (TextBox) HandleKey(e KeyEvent) {
	if e.event.Ch == 'q' {
		Exit()
		e.consumed = true
	}
}

func defaultDisplayable() Displayable {
	ret := displayable{*new(Padding), *new(Size), defaultPalette, FitToParent, nil}
	return &ret
}

func displayableWithSize(size Size) Displayable {
	ret := displayable{*new(Padding), size, defaultPalette, FitToParent, nil}
	return &ret
}

func NewTextBox(text string) TextBox {
	ret := TextBox{text, displayableWithSize(Size{10, 5})}
	return ret
}

func SimpleTitledContainer() TitledContainer {
	text := "Bigwordrighhere, butbigworderetoo, it, it has nice and small words, no long schlbberknockers to put you out of your lelelle"
	ret := TitledContainer{"Title", true, NewTextBox(text), displayableWithSize(Size{20, 10})}
	textb := ret.content.(TextBox)
	textb.SetParent(&ret)
	ret.SetLayout(Centered)
	return ret
}
