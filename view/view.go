package view

import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils"
import "strconv"
import "unicode/utf8"
import "time"

type KeyEvent struct {
	event    tb.Event
	consumed bool
}

var (
	initialized  = false
	focused      Focusable
	stopChannel  = make(chan bool)
	eventChannel = make(chan KeyEvent)
)

func setFocused(f Focusable) {
	stopChannel <- true
	focused = f
	acceptInput(f)
}

const (
	hLine             = '\u2500'
	vLine             = '\u2502'
	topLeftCorner     = '\u250C'
	topRightCorner    = '\u2510'
	bottomLeftCorner  = '\u2514'
	bottomRightCorner = '\u2518'
)

func paintArea(a area, fgColor, bgColor Color) {
	fg := fgColor.asTermAttr()
	bg := bgColor.asTermAttr()
	for i := a.x1; i < a.x2; i++ {
		for j := a.y1; j < a.y2; j++ {
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

func drawOutline(area area, foregroundColor, backgroundColor Color) {
	tbFG := foregroundColor.asTermAttr()
	tbBG := backgroundColor.asTermAttr()

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

func drawTitledContainer(parentArea area, c TitledContainer) {
	workingArea := getWorkArea(parentArea, c.Size, c.Layout)
	drawOutline(workingArea, c.Palette.normalFG, c.Palette.normalBG)
	tbFG := c.Palette.normalFG.asTermAttr()
	tbBG := c.Palette.normalBG.asTermAttr()

	i := workingArea.x1 + 1
	j := workingArea.y1 + 1
	if c.TitleVisibility {
		var ch rune
		for ; i < workingArea.x2-1; i++ {
			if i-workingArea.x1-1 < len(c.Title) {
				ch = []rune(c.Title)[i-workingArea.x1-1]
			} else {
				ch = ' '
			}
			tb.SetCell(i, j, ch, tbFG, tbBG)
			tb.SetCell(i, j+1, '=', tbFG, tbBG)
		}
		j += 2
	}
	for i = workingArea.x1 + 1; i < workingArea.x2-1; i++ {
		for y := j; y < workingArea.y2-1; y++ {
			tb.SetCell(i, y, ' ', tbFG, tbBG)
		}
	}
	childX, childY, s := c.DrawableArea()
	drawDisplayable(newArea(childX+workingArea.x1, childY+workingArea.y1, s), c.Content())
}

func getTitledContainerDrawableArea(totalArea area, t TitledContainer) area {
	offset := 0
	if t.TitleVisibility {
		offset = 2
	}
	return area{
		totalArea.x1 + 1 + t.Padding.Left,
		totalArea.y1 + 1 + offset + t.Padding.Up,
		totalArea.Width() - t.Padding.Left - t.Padding.Right - 2,
		totalArea.Height() - t.Padding.Up - t.Padding.Down - 2 - offset,
	}
}

func drawDisplayable(parentArea area, c interface{}) {
	switch v := c.(type) {
	case *TextBox:
		printTextBox(parentArea, *v)
	case TitledContainer:
		drawTitledContainer(parentArea, v)
	default:
		fmt.Printf("Type of c: %T\n", v)
	}
}

func printTextBox(parentArea area, t TextBox) {
	workingArea := getWorkArea(parentArea, t.Size, t.Layout)
	paintArea(parentArea, t.Palette.normalFG, t.Palette.normalBG)
	wrapped := utils.WrapText(t.Text, workingArea.Width(), workingArea.Height())
	for i, str := range wrapped {
		printStr(str, workingArea.x1, workingArea.y1+i, t.Palette.normalFG, t.Palette.normalBG)
	}
}

// TODO: Handle areas larger than parent area, or in general outside it.
func getWorkArea(parentArea area, contentSize Size, layout Layout) area {
	var width, height, x, y int
	switch layout {
	case FitToParent:
		width = parentArea.Width()
		height = parentArea.Height()
		x = parentArea.x1
		y = parentArea.y1
	case Centered:
		width = contentSize.Width
		height = contentSize.Height
		x = parentArea.x1 + (parentArea.Width()-contentSize.Width)/2
		fmt.Println("x =", parentArea.x1, "+ (", parentArea.Width(), "-", contentSize.Width, ")/2")
		y = parentArea.y1 + (parentArea.Height()-contentSize.Height)/2
	default:
		panic("nooo")
	}
	return newArea(x, y, Size{width, height})
}

func Init() {
	if !initialized {
		initialized = true
		tb.Flush()
		tb.HideCursor()
		container := SimpleTitledContainer()
		drawDisplayable(getTerminalArea(), container)
		textb := container.content.(*TextBox)
		focused = textb
		go pollEvents()

		for {
			tb.Sync()
			select {
			case ev := <-eventChannel:
				executeEvent(ev, focused)
			case <-stopChannel:
				return
			}
		}
	}
}

func pollEvents() {
	for {
		select {
		default:
			switch ev := tb.PollEvent(); ev.Type {
			case tb.EventKey:
				eventChannel <- KeyEvent{ev, false}
			default:
				fmt.Println("Uncovered event:", ev)
				time.Sleep(10 * time.Millisecond)
			}

		case <-stopChannel:
			return
		}

	}
}

func acceptInput(f Focusable) {
	go func() {
		for {
			select {
			case event := <-eventChannel:
				executeEvent(event, f)
			case <-stopChannel:
				return
			}
		}
	}()
}

func exit() {
	close(stopChannel)

}

func intCharSize(i int) int {
	return utf8.RuneCountInString(strconv.Itoa(i))
}

func getTerminalArea() area {
	y, x := tb.Size()
	fmt.Println("Terminal size:", x, y)
	return newArea(0, 0, Size{x, y})
}
