package containers

import _ "fmt"
import "github.com/teolandon/hanoi/structs"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import _ "github.com/teolandon/hanoi/view"
import _ "github.com/teolandon/hanoi/view/colors"
import _ "github.com/teolandon/hanoi/view/prints"
import "github.com/teolandon/hanoi/view/displayable"
import _ "github.com/teolandon/hanoi/utils/log"

type Orientation bool

const (
	Vertical   Orientation = true
	Horizontal Orientation = false
)

type Split struct {
	structs.Displayable
	separators  []int
	children    []displayable.Displayable
	orientation Orientation
}

func (t *Split) SetGrid(g pixel.SubGrid) {
	// Stub
}

func (t Split) Children() []displayable.Displayable {
	return t.children
}

func (c Split) ContentAreas() []areas.Area {
	if c.orientation == Horizontal {
		return c.hContentAreas()
	} else {
		return c.vContentAreas()
	}
}

// The two helpers could be refactored even more,
// but I feel like it would get in the way of
// readability.

func (c Split) hContentAreas() []areas.Area {
	length := len(c.separators)
	ret := make([]areas.Area, length-1)

	for i := 0; i < length-1; i++ {
		sep1 := c.separators[i]
		sep2 := c.separators[i+1]

		ret[i] = areas.New(sep1, 0, sep2, c.Size().Height())
	}

	return ret
}

func (c Split) vContentAreas() []areas.Area {
	length := len(c.separators)
	ret := make([]areas.Area, length-1)

	for i := 0; i < length-1; i++ {
		sep1 := c.separators[i]
		sep2 := c.separators[i+1]

		ret[i] = areas.New(0, sep1, c.Size().Width(), sep2)
	}

	return ret
}

// Separators returns the user-settable separators
// of a Split c. The first (0) and last (MAX) separators
// are ommited.
func (c Split) Separators() []int {
	length := len(c.separators)
	return c.separators[1 : length-1]
}

// Separator returns a single separator
// at the given index i of a Split c.
// Returns -1 if index is invalid.
func (c Split) Separator(i int) int {
	if i < 1 || i >= len(c.separators)-1 {
		return -1
	}

	return c.Separators()[i]
}

func (c *Split) SetSeparator(i, pos int) {
	if i < 0 || i >= len(c.separators) {
		return
	}
	c.separators[i] = pos
}
