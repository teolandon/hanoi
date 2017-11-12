package view

import "github.com/teolandon/hanoi/view/colors"

type Displayable interface {
	Padding() Padding
	SetPadding(p Padding)
	Size() Size
	SetSize(s Size)
	Palette() colors.Palette
	SetPalette(p colors.Palette)
	Layout() Layout
	SetLayout(p Layout)
	Parent() Displayable
	SetParent(p Displayable)
	PixelGrid(workingArea area) colors.PixelGrid
}

type Parent interface {
	Children() []Displayable
}

type displayable struct {
	padding Padding
	size    Size
	palette colors.Palette
	layout  Layout
	parent  Displayable
}

func (d displayable) Padding() Padding {
	return d.padding
}

func (d *displayable) SetPadding(p Padding) {
	d.padding = p
}

func (d displayable) Size() Size {
	return d.size
}

func (d *displayable) SetSize(s Size) {
	d.size = s
}

func (d displayable) Palette() colors.Palette {
	return d.palette
}

func (d *displayable) SetPalette(p colors.Palette) {
	d.palette = p
}

func (d displayable) Layout() Layout {
	return d.layout
}

func (d *displayable) SetLayout(l Layout) {
	d.layout = l
}

func (d displayable) Parent() Displayable {
	return d.parent
}

func (d *displayable) SetParent(dis Displayable) {
	d.parent = dis
}

func defaultDisplayable() displayable {
	ret := displayable{*new(Padding), *new(Size), colors.DefaultPalette, Centered, nil}
	return ret
}

func displayableWithSize(size Size) displayable {
	ret := displayable{*new(Padding), size, colors.DefaultPalette, Centered, nil}
	return ret
}
