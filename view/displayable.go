package view

import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import _ "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"

type Displayable interface {
	Padding() areas.Padding
	SetPadding(p areas.Padding)
	Size() areas.Size
	RealSize() areas.Size
	SetSize(s areas.Size)
	Palette() colors.Palette
	SetPalette(p colors.Palette)
	Layout() Layout
	SetLayout(p Layout)
	Parent() Displayable
	SetParent(p Displayable)
	Draw()
	setGrid(grid pixel.SubGrid)
}

type Parent interface {
	Children() []Displayable
}

type displayable struct {
	padding areas.Padding
	size    areas.Size
	palette colors.Palette
	layout  Layout
	parent  Displayable
	grid    pixel.SubGrid
}

func (d displayable) Padding() areas.Padding {
	return d.padding
}

func (d *displayable) SetPadding(p areas.Padding) {
	d.padding = p
}

func (d displayable) Size() areas.Size {
	return d.size
}

func (d displayable) RealSize() areas.Size {
	return areas.NewSize(d.grid.Width(), d.grid.Height())
}

func (d *displayable) SetSize(s areas.Size) {
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

func (d *displayable) setGrid(g pixel.SubGrid) {
	d.grid = g
}

func defaultDisplayable() displayable {
	ret := displayable{*new(areas.Padding), *new(areas.Size), colors.DefaultPalette, Centered, nil, *new(pixel.SubGrid)}
	return ret
}

func displayableWithSize(size areas.Size) displayable {
	ret := displayable{*new(areas.Padding), size, colors.DefaultPalette, Centered, nil, *new(pixel.SubGrid)}
	return ret
}

type Container interface {
	Content() Displayable
	SetContent(d Displayable)
	ContentPadding() areas.Padding
}

type container struct {
	displayable
	content Displayable
}

func (c container) Content() Displayable {
	return c.content
}

func (c *container) SetContent(d Displayable) {
	c.content = d
}

func (c *container) setGrid(g pixel.SubGrid, contentPadding areas.Padding) {
	c.displayable.setGrid(g)
	c.content.setGrid(g.Padded(contentPadding))
}
