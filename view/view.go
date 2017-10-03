package view

import "fmt"
import tb "github.com/nsf/termbox-go"
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

func drawOutline(coords Coords, size Size) {
	tb.SetCell(coords.X, coords.Y, topLeftCorner, foregroundColor, backgroundColor)
	tb.SetCell(coords.X+size.Width-1, coords.Y, topRightCorner, foregroundColor, backgroundColor)
	tb.SetCell(coords.X, coords.Y+size.Height-1, bottomLeftCorner, foregroundColor, backgroundColor)
	tb.SetCell(coords.X+size.Width-1, coords.Y+size.Height-1, bottomRightCorner, foregroundColor, backgroundColor)
	for i := coords.X + 1; i < coords.X+size.Width-1; i++ {
		tb.SetCell(i, coords.Y, hLine, foregroundColor, backgroundColor)
		tb.SetCell(i, coords.Y+size.Height-1, hLine, foregroundColor, backgroundColor)
	}
	for j := coords.Y + 1; j < coords.Y+size.Height-1; j++ {
		tb.SetCell(coords.X, j, vLine, foregroundColor, backgroundColor)
		tb.SetCell(coords.X+size.Width-1, j, vLine, foregroundColor, backgroundColor)
	}
}

func drawContainer(c TitledContainer) {
	drawOutline(c.Coords, c.Size)

	i := c.Coords.X + 1
	j := c.Coords.Y + 1
	if c.TitleVisibility {
		var ch rune
		for ; i < c.Coords.X+c.Size.Width-1; i++ {
			if i-c.Coords.X-1 < len(c.Title) {
				ch = []rune(c.Title)[i-c.Coords.X-1]
			} else {
				ch = ' '
			}
			tb.SetCell(i, j, ch, foregroundColor, backgroundColor)
			tb.SetCell(i, j+1, '=', foregroundColor, backgroundColor)
		}
		j += 2
	}
	for i = c.Coords.X + 1; i < c.Coords.X+c.Size.Height-1; i++ {
		for y := j; y < c.Coords.Y+c.Size.Height-1; y++ {
			tb.SetCell(i, y, ' ', foregroundColor, backgroundColor)
		}
	}
	drawContent(c.Coords.X+c.Padding.Left, j+c.Padding.Up, c.Content)
}

func drawContent(x, y int, c interface{}) {
	switch v := c.(type) {
	case *TextBox:
		printStr(v.Text, x, y)
	}
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		tb.HideCursor()
		text := SimpleTitledContainer()
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
