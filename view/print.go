package view

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/view/colors"

const (
	hLine             = '\u2500'
	vLine             = '\u2502'
	topLeftCorner     = '\u250C'
	topRightCorner    = '\u2510'
	bottomLeftCorner  = '\u2514'
	bottomRightCorner = '\u2518'
)

func drawOutline(area area, foregroundColor, backgroundColor colors.Color) {
	tbFG := foregroundColor.AsTermAttr()
	tbBG := backgroundColor.AsTermAttr()

	// Setting corners
	tb.SetCell(area.x1, area.y1, topLeftCorner, tbFG, tbBG)
	tb.SetCell(area.x2-1, area.y1, topRightCorner, tbFG, tbBG)
	tb.SetCell(area.x1, area.y2-1, bottomLeftCorner, tbFG, tbBG)
	tb.SetCell(area.x2-1, area.y2-1, bottomRightCorner, tbFG, tbBG)

	// Drawing sides
	for i := area.x1 + 1; i < area.x2-1; i++ {
		tb.SetCell(i, area.y1, hLine, tbFG, tbBG)
		tb.SetCell(i, area.y2-1, hLine, tbFG, tbBG)
	}
	for j := area.y1 + 1; j < area.y2-1; j++ {
		tb.SetCell(area.x1, j, vLine, tbFG, tbBG)
		tb.SetCell(area.x2-1, j, vLine, tbFG, tbBG)
	}
}

func paintArea(a area, fgColor, bgColor colors.Color) {
	fg := fgColor.AsTermAttr()
	bg := bgColor.AsTermAttr()
	for i := a.x1; i < a.x2; i++ {
		for j := a.y1; j < a.y2; j++ {
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
