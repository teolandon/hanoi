package view

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/view/colors"
import "github.com/teolandon/hanoi/areas"

const (
	hLine             = '\u2500'
	vLine             = '\u2502'
	topLeftCorner     = '\u250C'
	topRightCorner    = '\u2510'
	bottomLeftCorner  = '\u2514'
	bottomRightCorner = '\u2518'
)

func paintArea(a areas.Area, fgColor, bgColor colors.Color) {
	fg := fgColor.AsTermAttr()
	bg := bgColor.AsTermAttr()
	for i := a.X1(); i < a.X2(); i++ {
		for j := a.Y1(); j < a.Y2(); j++ {
			tb.SetCell(i, j, ' ', fg, bg)
		}
	}
	tb.Sync()
}

func printStr(s string, x, y int, foregroundColor, backgroundColor colors.Color) {
	printHelper(s, x, y, foregroundColor.AsTermAttr(), backgroundColor.AsTermAttr())
	tb.Sync()
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}
