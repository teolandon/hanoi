package view

import "errors"
import fmt "fmt"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/view/colors"
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

func chainOfFocus(f Displayable) {
	if f == nil {
		log.Log("End of chain")
		return
	}

	log.Log(f)
	chainOfFocus(f.Parent())
}

func SetFocused(f Displayable) {
	focused = f
	chainOfFocus(f)
	if initialized {
		focusChange <- true
		acceptInput(f)
	}
}

func SetRoot(f Displayable) {
	main.child = f
	f.SetParent(main)
}

func printGrid(grid colors.PixelGrid, x, y int) {
	for i, line := range grid {
		for j, pixel := range line {
			fg := pixel.Palette.GetFGTermAttr(pixel.Highlight)
			bg := pixel.Palette.GetBGTermAttr(pixel.Highlight)
			tb.SetCell(x+j, y+i, pixel.Char, fg, bg)
		}
	}
}

func drawDisplayable(parentArea area, d Displayable) {
	workArea := getWorkArea(parentArea, d.Size(), d.Layout())

	log.Log("Drawing:", d)
	log.Log("Area:", workArea)

	grid := d.PixelGrid(workArea)
	printGrid(grid, workArea.x1, workArea.y1)
}

// TODO: Handle areas larger than parent area, or in general outside it.
func getWorkArea(parentArea area, contentSize Size, layout Layout) area {
	var width, height, x, y int
	switch layout {
	case FitToParent:
		width = parentArea.width()
		height = parentArea.height()
		x = parentArea.x1
		y = parentArea.y1
	case Centered:
		width = contentSize.Width
		height = contentSize.Height
		x = parentArea.x1 + (parentArea.width()-width)/2
		y = parentArea.y1 + (parentArea.height()-height)/2
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

	// tb.SetInputMode(tb.InputEsc)
	// tb.SetOutputMode(tb.OutputNormal)
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
			initialized = false
			log.Close()
			tb.Close()
			close(StoppedChannel)
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
				log.Log("Uncovered event:", ev)
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

func (m mainContainer) String() string {
	return "Main Container"
}

func (m mainContainer) PixelGrid(a area) colors.PixelGrid {
	return m.child.PixelGrid(a)
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

func (mainContainer) Palette() colors.Palette {
	return colors.DefaultPalette
}

func (mainContainer) SetPalette(p colors.Palette) {}

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
	log.Log("Terminal size:", x, y)
	ret := newArea(0, 0, Size{x, y})
	log.Log("Terminal width:", ret.width())
	return ret
}
