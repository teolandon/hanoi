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
	NormalFG      Color
	NormalBG      Color
	HighlightFG   Color
	HighlightBG   Color
	Accent1FG     Color
	Accent1BG     Color
	Accent2FG     Color
	Accent2BG     Color
	NormalAttr    Attribute
	HighlightAttr Attribute
	Accent1Attr   Attribute
	Accent2Attr   Attribute
}

var (
	DefaultPalette = Palette{Black, Blue, Black, Yellow, Green, White,
		Magenta, Cyan, AttrDefault, AttrDefault, AttrDefault, AttrDefault}
	InheritAll = Palette{Inherit, Inherit, Inherit, Inherit, Inherit, Inherit,
		Inherit, Inherit, AttrInherit, AttrInherit, AttrInherit, AttrInherit}
	AlternatePalette = Palette{Green, Black, Yellow, Black, Magenta, Cyan,
		Black, White, AttrDefault, AttrDefault, AttrDefault, AttrDefault}
)

type Pixel struct {
	Char  rune
	Color Color
	Attr  Attribute
}

type PixelGrid [][]Pixel
