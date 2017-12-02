package view

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view/colors"

type mainContainer struct {
	child Displayable
}

func (m mainContainer) String() string {
	return "Main Container"
}

func (m mainContainer) Draw() {}

func (mainContainer) HandleKey(e KeyEvent) {
	if e.event.Ch == 'q' {
		Exit()
		e.consumed = true
	}
}

func (m mainContainer) Children() []Displayable {
	ret := make([]Displayable, 1)
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
	return colors.DefaultPalette
}

func (mainContainer) SetPalette(p colors.Palette) {}

func (mainContainer) Layout() Layout {
	return FitToParent
}

func (mainContainer) SetLayout(l Layout) {}

func (mainContainer) Parent() Displayable {
	return nil
}

func (mainContainer) SetParent(d Displayable) {}

func (mainContainer) setGrid(g pixel.SubGrid) {}

func (m mainContainer) Content() Displayable {
	return m.child
}

func (m mainContainer) SetContent(d Displayable) {
	m.child = d
	d.SetParent(m)
}

func (m mainContainer) ContentPadding() areas.Padding {
	return areas.Padding{0, 0, 0, 0}
}
