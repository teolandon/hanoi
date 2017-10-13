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
	initialized = false
	focused     Displayable
	stopChannel = make(chan bool)
	// StoppedChannel is the channel that signals that
	// the view has been stopped and exited out of.
	StoppedChannel = make(chan bool)
	eventChannel   = make(chan KeyEvent)
)

func setFocused(f Displayable) {
	stopChannel <- true
	focused = f
	acceptInput(f)
}

func drawTitledContainer(parentArea area, c TitledContainer) {
	workingArea := getWorkArea(parentArea, c.Size(), c.Layout())
	drawOutline(workingArea, c.Palette().normalFG, c.Palette().normalBG)
	tbFG := c.Palette().normalFG.asTermAttr()
	tbBG := c.Palette().normalBG.asTermAttr()

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

func drawDisplayable(parentArea area, c interface{}) {
	switch v := c.(type) {
	case TextBox:
		drawTextBox(parentArea, v)
	case TitledContainer:
		drawTitledContainer(parentArea, v)
	default:
		fmt.Printf("Type of c: %T\n", v)
	}
}

func drawTextBox(parentArea area, t TextBox) {
	workingArea := getWorkArea(parentArea, t.Size(), t.Layout())
	paintArea(parentArea, t.Palette().normalFG, t.Palette().normalBG)
	wrapped := utils.WrapText(t.Text, workingArea.Width(), workingArea.Height())
	for i, str := range wrapped {
		printStr(str, workingArea.x1, workingArea.y1+i, t.Palette().normalFG, t.Palette().normalBG)
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
		textb := container.content.(TextBox)
		focused = textb
		go pollEvents()

		go func() {
			defer close(StoppedChannel)
			for {
				tb.Sync()
				select {
				case ev := <-eventChannel:
					executeEvent(ev, focused)
				case <-stopChannel:
					return
				}
			}
		}()
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

func acceptInput(f Displayable) {
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

// Stops all goroutines that run the view, and quits the view,
// reverting the terminal to normal mode
func Exit() {
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
