package view

import "errors"
import "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils"
import "time"

var (
	initialized = false
	focused     Displayable
	stopChannel chan bool
	focusChange = make(chan bool)
	// StoppedChannel is the channel that signals that
	// the view has been stopped and exited out of.
	StoppedChannel chan bool
	eventChannel   = make(chan KeyEvent)
	main           mainContainer
)

type KeyEvent struct {
	event    tb.Event
	consumed bool
}

func SetFocused(f Displayable) {
	focused = f
	if initialized {
		focusChange <- true
		acceptInput(f)
	}
}

func SetRoot(f Displayable) {
	main.child = f
	f.SetParent(main)
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
	case mainContainer:
		drawDisplayable(terminalArea(), v.child)
	default:
		fmt.Printf("Type of c: %T\n", v)
	}
}

func drawTextBox(parentArea area, t TextBox) {
	workingArea := getWorkArea(parentArea, t.Size(), t.Layout())
	fmt.Println("parent area for textbox:", parentArea)
	fmt.Println("working area for textbox:", workingArea)
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
		x = parentArea.x1 + (parentArea.Width()-width)/2
		fmt.Println("x =", parentArea.x1, "+ (", parentArea.Width(), "-", width, ")/2")
		y = parentArea.y1 + (parentArea.Height()-height)/2
	default:
		panic("nooo")
	}
	return newArea(x, y, Size{width, height})
}

func Init() error {
	if tb.IsInit {
		return errors.New("Termbox has already been initialized.")
	}
	if initialized {
		return errors.New("Hanoi has already been initialized.")
	}

	err := tb.Init()
	if err != nil {
		return err
	}

	tb.SetInputMode(tb.InputEsc)
	tb.SetOutputMode(tb.OutputNormal)
	tb.Clear(tb.ColorWhite, tb.ColorRed)

	initialized = true
	StoppedChannel = make(chan bool)
	stopChannel = make(chan bool)
	tb.Flush()
	tb.HideCursor()

	acceptInput(focused)

	drawDisplayable(terminalArea(), main)

	go pollEvents()

	go func() {
		defer func() {
			close(StoppedChannel)
			initialized = false
			tb.Close()
		}()
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

	return nil
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
			case <-focusChange:
				return
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

type mainContainer struct {
	child Displayable
}

func (mainContainer) HandleKey(e KeyEvent) {
	if e.event.Ch == 'q' {
		Exit()
		e.consumed = true
	}
}

func (m mainContainer) Children() []Displayable {
	ret := make([]Displayable, 1)
	ret[0] = m.child
	return ret
}

func (mainContainer) Padding() Padding {
	return Padding{0, 0, 0, 0}
}

func (mainContainer) SetPadding(p Padding) {}

func (mainContainer) Size() Size {
	x, y := tb.Size()
	return Size{x, y}
}

func (mainContainer) SetSize(s Size) {}

func (mainContainer) Palette() Palette {
	return defaultPalette
}

func (mainContainer) SetPalette(p Palette) {}

func (mainContainer) Layout() Layout {
	return FitToParent
}

func (mainContainer) SetLayout(l Layout) {}

func (mainContainer) Parent() Displayable {
	return nil
}

func (mainContainer) SetParent(d Displayable) {}

func terminalArea() area {
	x, y := tb.Size()
	fmt.Println("Terminal size:", x, y)
	ret := newArea(0, 0, Size{x, y})
	fmt.Println("Terminal width:", ret.Width())
	return ret
}
