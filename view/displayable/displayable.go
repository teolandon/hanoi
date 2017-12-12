package displayable

import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import _ "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"

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
	Grid() pixel.SubGrid
	SetGrid(grid pixel.SubGrid)
}

type Parent interface {
	Children() []Displayable
}

type Container interface {
	Content() Displayable
	SetContent(d Displayable)
	ContentPadding() areas.Padding
}
