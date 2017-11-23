package colors

import tb "github.com/nsf/termbox-go"

// A Color maps to a termbox-go color, but also supports
// inheriting.
type Color uint16

// Valid Colors
const (
	Default Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	Inherit Color = Color(^uint16(0))
)

type Highlight uint8

const (
	Normal Highlight = iota
	Highlighted
	Accent1
	Accent2
)

// Term Attributes that map to termbox-go attrs, but also support
// inheriting and a default attribute. Should not be used directly
// with termbox-go, but instead should obtain the correct termbox.Attribute
// type using the asTermAttr() or withAttr(attr) funcs.
type Attribute uint16

// Valid Attributes
const (
	AttrBold Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
	AttrDefault Attribute = 0
	AttrInherit Attribute = Attribute(^uint16(0))
)

func (c Color) AsTermAttr() tb.Attribute {
	return tb.Attribute(c - 1)
}

func (c Color) WithAttr(attr Attribute) tb.Attribute {
	return c.AsTermAttr() | tb.Attribute(attr)
}

type Palette struct {
	colors [8]Color
	attrs  [4]Attribute
}

var (
	DefaultPalette = Palette{[8]Color{Black, Blue, Black, Yellow, Green, White,
		Magenta, Cyan}, [4]Attribute{AttrDefault, AttrDefault, AttrDefault, AttrDefault}}
	InheritAll = Palette{[8]Color{Inherit, Inherit, Inherit, Inherit, Inherit, Inherit,
		Inherit, Inherit}, [4]Attribute{AttrInherit, AttrInherit, AttrInherit, AttrInherit}}
	AlternatePalette = Palette{[8]Color{Green, Black, Yellow, Black, Magenta, Cyan,
		Black, White}, [4]Attribute{AttrDefault, AttrDefault, AttrDefault, AttrDefault}}
)

func (p Palette) GetFGColor(t Highlight) Color {
	return p.colors[2*t]
}

func (p Palette) GetBGColor(t Highlight) Color {
	return p.colors[2*t+1]
}

func (p Palette) GetAttr(t Highlight) Attribute {
	return p.attrs[t]
}

func (p Palette) GetFGTermAttr(t Highlight) tb.Attribute {
	return p.GetFGColor(t).WithAttr(p.GetAttr(t))
}

func (p Palette) GetBGTermAttr(t Highlight) tb.Attribute {
	return p.GetBGColor(t).AsTermAttr()
}
