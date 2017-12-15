package view

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"
import "github.com/teolandon/hanoi/view/displayable"

type mainContainer struct {
	grid  pixel.SubGrid
	child displayable.Displayable
}

func (m mainContainer) String() string {
	return "Main Container"
}

func (m mainContainer) Draw() {
	pw := pixel.NewWriter(m.Palette(), colors.Normal, termGrid.TotalSubGrid())

	log.Log("Drawing main container:")
	pw.FillAll(' ')
	log.Log("Finished drawing main container")
}

func (mainContainer) HandleKey(e KeyEvent) {
	if e.event.Ch == 'q' {
		Exit()
		e.consumed = true
	}
}

func (m mainContainer) Children() []displayable.Displayable {
	ret := make([]displayable.Displayable, 1)
	ret[0] = m.child
	return ret
}

func (mainContainer) Padding() areas.Padding {
	return areas.Padding{0, 0, 0, 0}
}

func (mainContainer) SetPadding(p areas.Padding) {}

func (mainContainer) Size() areas.Size {
	x, y := tb.Size()
	return areas.NewSize(x, y)
}

func (mainContainer) RealSize() areas.Size {
	x, y := tb.Size()
	return areas.NewSize(x, y)
}

func (mainContainer) SetSize(s areas.Size) {}

func (mainContainer) Palette() colors.Palette {
	return colors.MainPalette
}

func (mainContainer) SetPalette(p colors.Palette) {}

func (mainContainer) Layout() areas.Layout {
	return areas.FitToParent
}

func (mainContainer) SetLayout(l areas.Layout) {}

func (mainContainer) Parent() displayable.Displayable {
	return nil
}

func (mainContainer) SetParent(d displayable.Displayable) {}

func (m mainContainer) Grid() pixel.SubGrid {
	return m.grid
}

func (m mainContainer) SetGrid(g pixel.SubGrid) {
	m.child.SetGrid(g)
}

func (m mainContainer) Content() displayable.Displayable {
	return m.child
}

func (m mainContainer) SetContent(d displayable.Displayable) {
	m.child = d
	d.SetParent(m)
}

func (m mainContainer) ContentPadding() areas.Padding {
	return areas.Padding{0, 0, 0, 0}
}
