package view

import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils"
import "strconv"
import "unicode/utf8"
import "time"

var (
	initialized = false
)

const (
	hLine             = '\u2500'
	vLine             = '\u2502'
	topLeftCorner     = '\u250C'
	topRightCorner    = '\u2510'
	bottomLeftCorner  = '\u2514'
	bottomRightCorner = '\u2518'
)

func paintArea(area Area, fgColor, bgColor Color) {
	fg := fgColor.asTermAttr()
	bg := bgColor.asTermAttr()
	for i := area.X1(); i < area.X2(); i++ {
		for j := area.Y1(); j < area.Y2(); j++ {
			fmt.Println("erasing at", i, ":", j)
			tb.SetCell(i, j, ' ', fg, bg)
		}
	}
	tb.Sync()
}

func printStr(s string, x, y int, foregroundColor, backgroundColor Color) {
	printHelper(s, x, y, foregroundColor.asTermAttr(), backgroundColor.asTermAttr())
	tb.Sync()
}

// Does NOT Sync Terminal
func printHelper(s string, x, y int, fg, bg tb.Attribute) {
	for i, r := range s {
		tb.SetCell(x+i, y, r, fg, bg)
	}
}

func drawOutline(area Area, foregroundColor, backgroundColor Color) {
	tbFG := foregroundColor.asTermAttr()
	tbBG := backgroundColor.asTermAttr()

	// Setting corners
	tb.SetCell(area.X1(), area.Y1(), topLeftCorner, tbFG, tbBG)
	tb.SetCell(area.X2()-1, area.Y1(), topRightCorner, tbFG, tbBG)
	tb.SetCell(area.X1(), area.Y2()-1, bottomLeftCorner, tbFG, tbBG)
	tb.SetCell(area.X2()-1, area.Y2()-1, bottomRightCorner, tbFG, tbBG)

	// Drawing sides
	for i := area.X1() + 1; i < area.X2()-1; i++ {
		tb.SetCell(i, area.Y1(), hLine, tbFG, tbBG)
		tb.SetCell(i, area.Y2()-1, hLine, tbFG, tbBG)
	}
	for j := area.Y1() + 1; j < area.Y2()-1; j++ {
		tb.SetCell(area.X1(), j, vLine, tbFG, tbBG)
		tb.SetCell(area.X2()-1, j, vLine, tbFG, tbBG)
	}
}

func drawTitledContainer(parentArea Area, c TitledContainer) {
	workingArea := getWorkArea(parentArea, Area{c.Size, c.Coords}, c.Layout)
	drawOutline(workingArea, c.Palette.normalFG, c.Palette.normalBG)
	tbFG := c.Palette.normalFG.asTermAttr()
	tbBG := c.Palette.normalBG.asTermAttr()

	i := workingArea.X1() + 1
	j := workingArea.Y1() + 1
	if c.TitleVisibility {
		var ch rune
		for ; i < workingArea.X2()-1; i++ {
			if i-workingArea.X1()-1 < len(c.Title) {
				ch = []rune(c.Title)[i-workingArea.X1()-1]
			} else {
				ch = ' '
			}
			tb.SetCell(i, j, ch, tbFG, tbBG)
			tb.SetCell(i, j+1, '=', tbFG, tbBG)
		}
		j += 2
	}
	for i = workingArea.X1() + 1; i < workingArea.X2()-1; i++ {
		for y := j; y < workingArea.Y2()-1; y++ {
			tb.SetCell(i, y, ' ', tbFG, tbBG)
		}
	}
	drawDisplayable(getTitledContainerDrawableArea(workingArea, c), c.Content)
}

func getTitledContainerDrawableArea(totalArea Area, t TitledContainer) Area {
	offset := 0
	if t.TitleVisibility {
		offset = 2
	}
	return Area{Size{totalArea.Width() - t.Padding.Left - t.Padding.Right - 2,
		totalArea.Height() - t.Padding.Up - t.Padding.Down - 2 - offset},
		Coords{totalArea.X1() + 1 + t.Padding.Left, totalArea.Y1() + 1 + offset + t.Padding.Up}}
}

func drawDisplayable(parentArea Area, c interface{}) {
	fmt.Println("Drawing in area {", parentArea.X1(), "-", parentArea.X2(), ":",
		parentArea.Y1(), "-", parentArea.Y2(), "}")
	switch v := c.(type) {
	case *TextBox:
		printTextBox(parentArea, *v)
	case TitledContainer:
		drawTitledContainer(parentArea, v)
	default:
		fmt.Printf("Type of c: %T\n", v)
	}
}

func printTextBox(parentArea Area, t TextBox) {
	workingArea := getWorkArea(parentArea, Area{t.Size, t.Coords}, t.Layout)
	paintArea(parentArea, t.Palette.normalFG, t.Palette.normalBG)
	wrapped := strutil.WrapText(t.Text, workingArea.Width(), workingArea.Height())
	for i, str := range wrapped {
		fmt.Println("printing "+str+" at", workingArea.X1(), workingArea.Y1()+i)
		printStr(str, workingArea.X1(), workingArea.Y1()+i, t.Palette.normalFG, t.Palette.normalBG)
	}
}

// TODO: Handle areas larger than parent area, or in general outside it.
func getWorkArea(parentArea Area, contentArea Area, layout Layout) Area {
	var width, height, x, y int
	switch layout {
	case FitToParent:
		width = parentArea.Width()
		height = parentArea.Height()
		x = parentArea.X1()
		y = parentArea.Y1()
	case Centered:
		width = contentArea.Width()
		height = contentArea.Height()
		x = parentArea.X1() + (parentArea.Width()-contentArea.Width())/2
		y = parentArea.Y1() + (parentArea.Height()-contentArea.Height())/2
	case Absolute:
		width = contentArea.Width()
		height = contentArea.Height()
		x = parentArea.X1() + contentArea.X1()
		y = parentArea.Y1() + contentArea.Y1()
	default:
		panic("nooo")
	}
	return Area{Size{width, height}, Coords{x, y}}
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		tb.HideCursor()
		container := SimpleTitledContainer()
		drawDisplayable(getTerminalArea(), container)

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

func getTerminalArea() Area {
	x, y := tb.Size()
	return Area{Size{x, y}, Coords{0, 0}}
}
