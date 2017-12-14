package containers

import "fmt"
import "github.com/teolandon/hanoi/structs"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view"
import "github.com/teolandon/hanoi/view/colors"
import "github.com/teolandon/hanoi/view/prints"
import "github.com/teolandon/hanoi/view/displayable"
import "github.com/teolandon/hanoi/utils/log"

type TitledContainer struct {
	Title           string
	TitleVisibility bool
	structs.SingleContainer
}

func (t *TitledContainer) SetGrid(g pixel.SubGrid) {
	t.SingleContainer.SetGrid(g, t.ContentPadding())
	log.Log("OGrid length:", g.Width())
	log.Log("grid length:", len(t.Grid().GetLine(0)))
}

func (t TitledContainer) ContentPadding() areas.Padding {
	pad := 0

	if t.TitleVisibility {
		pad = 2
	}

	return areas.Padding{1, 1, 1 + pad, 1}
}

func (t TitledContainer) String() string {
	return fmt.Sprintf("Titled container with title %s", t.Title)
}

func (t TitledContainer) Draw() {
	pw := pixel.NewWriter(t.Palette(), colors.Normal, t.Grid())

	pw.FillAll(' ')

	for i := range t.Grid().GetLine(0) {
		pw.Write(i, 0, prints.HLine)
	}
	var i int
	for i = 1; i < t.Grid().Height()-1; i++ {
		pw.Write(0, i, prints.VLine)
		pw.Write(t.Grid().Width()-1, i, prints.VLine)
	}
	for j := range t.Grid().GetLine(i) {
		pw.Write(j, i, prints.HLine)
	}

	// Corners
	pw.Write(0, 0, prints.TopLeftCorner)
	pw.Write(0, i, prints.BottomLeftCorner)
	pw.Write(t.Grid().Width()-1, 0, prints.TopRightCorner)
	pw.Write(t.Grid().Width()-1, i, prints.BottomRightCorner)

	if !t.TitleVisibility {
		return
	} // Title drawing

	log.Log("len(line2) =", len(t.Grid().GetLine(2)))
	log.Log("bounds:", 1, t.Grid().Width()-1)
	for i := 0; i < t.Grid().Width()-2; i++ {
		pw.Write(i+1, 1, ' ')
		pw.Write(i+1, 2, prints.HLine)
	}
	pw.WriteStr(1, 1, t.Title)
}

func NewTitled() *TitledContainer {
	text := "Bigwordrighhere, butbigworderetoo, it, it has nice and small words, no long schlbberknockers to put you out of your lelelle"
	textBox := New(text)
	textBox.SetLayout(displayable.FitToParent)
	ret := TitledContainer{"Title", true, structs.NewContainer(20, 10)}
	ret.SetContent(&textBox)
	ret.SetLayout(displayable.Centered)
	return &ret
}

func NewWithButton() *TitledContainer {
	ret := TitledContainer{"Test", true, structs.NewContainer(20, 10)}
	button := view.NewButton("OK")
	ret.SetContent(&button)
	log.Log("Content of buttoncontainer:", ret.Content())
	log.Log("parent of button:", ret.Content().Parent())
	ret.SetLayout(displayable.Centered)
	return &ret
}
