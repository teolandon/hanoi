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

func printStr(s string, x, y int, foregroundColor, backgroundColor colors.Color) {
	printHelper(s, x, y, foregroundColor.AsTermAttr(), backgroundColor.AsTermAttr())
	tb.Flush()
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}
