package prints

import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/view/colors"

const (
	HLine             = '\u2500'
	VLine             = '\u2502'
	TopLeftCorner     = '\u250C'
	TopRightCorner    = '\u2510'
	BottomLeftCorner  = '\u2514'
	BottomRightCorner = '\u2518'
)

func PrintStr(s string, x, y int, foregroundColor, backgroundColor colors.Color) {
	printHelper(s, x, y, foregroundColor.AsTermAttr(), backgroundColor.AsTermAttr())
	tb.Flush()
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}
