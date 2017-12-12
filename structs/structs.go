package structs

import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import _ "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"
import displ "github.com/teolandon/hanoi/view/displayable"

// structs is a helper package to provide structs that fulfill all or
// most of the methods of some interfaces, in order to avoid rewriting
// the same methods

type Displayable struct {
	padding areas.Padding
	size    areas.Size
	palette colors.Palette
	layout  displ.Layout
	parent  displ.Displayable
	grid    pixel.SubGrid
}

func (d Displayable) Padding() areas.Padding {
	return d.padding
}

func (d *Displayable) SetPadding(p areas.Padding) {
	d.padding = p
}

func (d Displayable) Size() areas.Size {
	return d.size
}

func (d Displayable) RealSize() areas.Size {
	return areas.NewSize(d.grid.Width(), d.grid.Height())
}

func (d *Displayable) SetSize(s areas.Size) {
	d.size = s
}

func (d Displayable) Palette() colors.Palette {
	return d.palette
}

func (d *Displayable) SetPalette(p colors.Palette) {
	d.palette = p
}

func (d Displayable) Layout() displ.Layout {
	return d.layout
}

func (d *Displayable) SetLayout(l displ.Layout) {
	d.layout = l
}

func (d Displayable) Parent() displ.Displayable {
	return d.parent
}

func (d *Displayable) SetParent(dis displ.Displayable) {
	d.parent = dis
}

func (d Displayable) Grid() pixel.SubGrid {
	return d.grid
}

func (d *Displayable) SetGrid(g pixel.SubGrid) {
	d.grid = g
}

func (Displayable) Draw() {
	// STUB
}

func DefaultDisplayable() Displayable {
	ret := Displayable{*new(areas.Padding), *new(areas.Size), colors.DefaultPalette, displ.Centered, nil, *new(pixel.SubGrid)}
	return ret
}

func DisplayableWithSize(size areas.Size) Displayable {
	ret := Displayable{*new(areas.Padding), size, colors.DefaultPalette, displ.Centered, nil, *new(pixel.SubGrid)}
	return ret
}

// Container

type Container struct {
	displ.Displayable
	content displ.Displayable
}

func (c Container) Content() displ.Displayable {
	return c.content
}

func (c *Container) SetContent(d displ.Displayable) {
	c.content = d
	d.SetParent(c.Displayable)
}

func (c *Container) SetGrid(g pixel.SubGrid, contentPadding areas.Padding) {
	c.Displayable.SetGrid(g)
	c.content.SetGrid(g.Padded(contentPadding))
}

func NewContainer() Container {
	d := DefaultDisplayable()
	return Container{&d, nil}
}
