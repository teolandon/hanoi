package structs

import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view/colors"
import displ "github.com/teolandon/hanoi/view/displayable"

// structs is a helper package to provide structs that fulfill all or
// most of the methods of some interfaces, in order to avoid rewriting
// the same methods

type Displayable struct {
	padding areas.Padding
	size    areas.Size
	palette colors.Palette
	layout  areas.Layout
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

func (d Displayable) Layout() areas.Layout {
	return d.layout
}

func (d *Displayable) SetLayout(l areas.Layout) {
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
	ret := Displayable{*new(areas.Padding), *new(areas.Size), colors.DefaultPalette, areas.Centered, nil, *new(pixel.SubGrid)}
	return ret
}

func DisplayableWithSize(size areas.Size) Displayable {
	ret := Displayable{*new(areas.Padding), size, colors.DefaultPalette, areas.Centered, nil, *new(pixel.SubGrid)}
	return ret
}

// Container

type SingleContainer struct {
	displ.Displayable
	content displ.Displayable
}

func (c SingleContainer) Content() displ.Displayable {
	return c.content
}

func (c *SingleContainer) SetContent(d displ.Displayable) {
	c.content = d
	d.SetParent(c.Displayable)
}

// SetGrid is to be called by the wrapper of SingleContainer providing the
// contentarea that is defined by the wrapper.  This will automate the setting
// of the grid of the content.
//
// NOTE: The contentArea field might be phased out to a function in the
// SingleContainer struct, so that the function is provided when the
// embedded struct is created in the wrapper's constructor.
func (c *SingleContainer) SetGrid(g pixel.SubGrid, contentArea areas.Area) {
	layoutedGrid := g.Layouted(c.Size(), c.Layout())
	c.Displayable.SetGrid(layoutedGrid)

	c.content.SetGrid(g.SubGrid(contentArea))
}

func NewContainer(width, height int) SingleContainer {
	d := DisplayableWithSize(areas.NewSize(width, height))
	return SingleContainer{&d, nil}
}

// MultiContainer

type MultiContainer struct {
	displ.Displayable
	children []displ.Displayable
}

func (c MultiContainer) Children() []displ.Displayable {
	return c.children
}
