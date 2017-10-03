package view

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

func (c Color) asTermAttr() tb.Attribute {
	return tb.Attribute(c - 1)
}

func (c Color) withAttr(attr Attribute) tb.Attribute {
	return c.asTermAttr() | tb.Attribute(attr)
}
