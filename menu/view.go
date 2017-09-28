package menu

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
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		tb.HideCursor()
		fmt.Println("Setting cell")
		text := containers.SimpleTitledContainer()
		for y, sl := range text.Grid() {
			for x, pixel := range sl {
				fmt.Println("Setting cell", x, y, "as", pixel.Ch)
				tb.SetCell(x, y, pixel.Ch, foregroundColor, backgroundColor)
			}
		}

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
