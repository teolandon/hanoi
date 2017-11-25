package view

import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view/colors"

type Displayable interface {
	Padding() Padding
	SetPadding(p Padding)
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
	setGrid(grid pixel.PixelGrid)
}

type Parent interface {
	Children() []Displayable
}

type displayable struct {
	padding Padding
	size    areas.Size
	palette colors.Palette
	layout  Layout
	parent  Displayable
	grid    pixel.PixelGrid
}

func (d displayable) Padding() Padding {
	return d.padding
}

func (d *displayable) SetPadding(p Padding) {
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

func (d *displayable) setGrid(g pixel.PixelGrid) {
	d.grid = g
}

func defaultDisplayable() displayable {
	ret := displayable{*new(Padding), *new(areas.Size), colors.DefaultPalette, Centered, nil, *new(pixel.PixelGrid)}
	return ret
}

func displayableWithSize(size areas.Size) displayable {
	ret := displayable{*new(Padding), size, colors.DefaultPalette, Centered, nil, *new(pixel.PixelGrid)}
	return ret
}

type Container interface {
	Content() Displayable
	SetContent(d Displayable)
	ContentArea() areas.Area
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

func (c *container) setGrid(g pixel.PixelGrid, contentArea areas.Area) {
	c.displayable.setGrid(g)
	c.content.setGrid(g.SubGrid(contentArea.X1(), contentArea.X2(), contentArea.Y1(), contentArea.Y2()))
}
