package displayable

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
	Layout() areas.Layout
	SetLayout(p areas.Layout)
	Parent() Displayable
	SetParent(p Displayable)
	Draw()
	Grid() pixel.SubGrid
	SetGrid(grid pixel.SubGrid)
}

type Parent interface {
	Children() []Displayable
}

type SingleContainer interface {
	Content() Displayable
	SetContent(d Displayable)
	ContentArea() areas.Area
}

type MultiContainer interface {
	Children() []Displayable
	// SetChildren(children []Displayable)
	ContentAreas() []areas.Area
}
