package view

import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/containers"
import "strconv"
import "unicode/utf8"
import "time"

var (
	initialized = false
)

const (
	backgroundColor tb.Attribute = tb.ColorBlue
	foregroundColor tb.Attribute = tb.ColorDefault
	ambientColor    tb.Attribute = tb.ColorDefault
	selectedFGColor tb.Attribute = foregroundColor | tb.AttrReverse
	selectedBGColor tb.Attribute = backgroundColor | tb.AttrReverse
)

const (
	hLine             = '\u2500'
	vLine             = '\u2502'
	topLeftCorner     = '\u250C'
	topRightCorner    = '\u2510'
	bottomLeftCorner  = '\u2514'
	bottomRightCorner = '\u2518'
)

func eraseRange(start, end, y int) {
	for ; start < end; start++ {
		eraseHelper(start, y)
	}
	tb.Sync()
}

// Does NOT Sync terminal
func eraseHelper(x, y int) {
	tb.SetCell(x, y, ' ', ambientColor, ambientColor)
}

func printStr(s string, x, y int) {
	printHelper(s, x, y, foregroundColor, backgroundColor)
	tb.Sync()
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}

func drawOutline(box containers.Rect) {
	tb.SetCell(box.X1, box.Y1, topLeftCorner, foregroundColor, backgroundColor)
	tb.SetCell(box.X2-1, box.Y1, topRightCorner, foregroundColor, backgroundColor)
	tb.SetCell(box.X1, box.Y2-1, bottomLeftCorner, foregroundColor, backgroundColor)
	tb.SetCell(box.X2-1, box.Y2-1, bottomRightCorner, foregroundColor, backgroundColor)
	for i := box.X1 + 1; i < box.X2-1; i++ {
		tb.SetCell(i, box.Y1, hLine, foregroundColor, backgroundColor)
		tb.SetCell(i, box.Y2-1, hLine, foregroundColor, backgroundColor)
	}
	for j := box.Y1 + 1; j < box.Y2-1; j++ {
		tb.SetCell(box.X1, j, vLine, foregroundColor, backgroundColor)
		tb.SetCell(box.X2-1, j, vLine, foregroundColor, backgroundColor)
	}
}

func drawContainer(c containers.TitledContainer) {
	drawOutline(c.Bounds())

	i := c.Bounds().X1 + 1
	j := c.Bounds().Y1 + 1
	if c.TitleVisibility {
		var ch rune
		for ; i < c.Bounds().X2-1; i++ {
			if i-c.Bounds().X1-1 < len(c.Title) {
				ch = []rune(c.Title)[i-c.Bounds().X1-1]
			} else {
				ch = ' '
			}
			tb.SetCell(i, j, ch, foregroundColor, backgroundColor)
			tb.SetCell(i, j+1, '=', foregroundColor, backgroundColor)
		}
		j += 2
	}
	for i = c.Bounds().X1 + 1; i < c.Bounds().X2-1; i++ {
		for y := j; y < c.Bounds().Y2-1; y++ {
			tb.SetCell(i, y, ' ', foregroundColor, backgroundColor)
		}
	}
	drawContainable(c.Bounds().X1+c.Padding().Left, c.Bounds().Y1+c.Padding().Up, c.Content)
}

func drawContainable(x, y int, c containers.Containable) {
	switch v := c.(type) {
	case *containers.TextBox:
		printStr(v.Text, x, y)
	}
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		tb.HideCursor()
		text := containers.SimpleTitledContainer()
		drawContainer(text)

	loop:
		for {
			tb.Sync()
			switch ev := tb.PollEvent(); ev.Type {
			case tb.EventKey:
				if ev.Key == tb.KeyEsc || ev.Ch == 'q' {
					break loop
				} else {
					fmt.Println("Uncovered event:", ev)
				}
			default:
				fmt.Println("Uncovered event:", ev)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func intCharSize(i int) int {
	return utf8.RuneCountInString(strconv.Itoa(i))
}
