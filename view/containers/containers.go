package containers

import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/structs"
import "github.com/teolandon/hanoi/utils/strutils"
import "github.com/teolandon/hanoi/view/colors"
import _ "github.com/teolandon/hanoi/utils/log"

type TextBox struct {
	Text string
	structs.Displayable
}

func (t TextBox) Draw() {
	pw := pixel.NewWriter(t.Palette(), colors.Normal, t.Grid())

	wrapped := strutils.WrapText(t.Text, t.Grid().Width(), t.Grid().Height())
	for i, str := range wrapped {
		pw.WriteStr(0, i, str)
	}
}

func New(text string) TextBox {
	ret := TextBox{text, structs.DisplayableWithSize(areas.NewSize(10, 5))}
	return ret
}
