package view

import "errors"
import "fmt"
import "time"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"
import "github.com/teolandon/hanoi/view/displayable"

var (
	initialized = false
	focused     displayable.Displayable
	stopChannel chan bool // When closed, stops all goroutines running in hanoi
	focusChange = make(chan bool)

	// StoppedChannel is the channel that signals that
	// the view has been stopped and exited out of.
	StoppedChannel chan bool
	eventChannel   = make(chan KeyEvent)
	main           mainContainer
	termGrid       pixel.SquareGrid
)

type KeyEvent struct {
	event    tb.Event
	consumed bool
}

func SetFocused(f displayable.Displayable) {
	focused = f
	if initialized {
		focusChange <- true
		acceptInput(f)
	}
}

func SetRoot(f displayable.Displayable) {
	main.child = f
	f.SetParent(main)
}

func printGrid(grid pixel.SubGrid, x, y int) {
	for i := 0; i < grid.Height(); i++ {
		line := grid.GetLine(i)
		for j, pixel := range line {
			fg := pixel.Palette.GetFGTermAttr(pixel.Highlight)
			bg := pixel.Palette.GetBGTermAttr(pixel.Highlight)
			tb.SetCell(x+j, y+i, pixel.Ch, fg, bg)
		}
	}
	tb.Flush()
}

func calculateGrids() {
	calculateGridsH(main, termGrid.TotalSubGrid())
}

func calculateGridsH(d displayable.Displayable, grid pixel.SubGrid) {
	if d == nil {
		return
	}

	d.SetGrid(grid)
	log.Log("Set grid for", d, "to", grid)

	// The following code was commented out due to
	// the behavior of the default struct for containers,
	// which now automatically sets the grids for its children.

	/* COMMENTED OUT
		 -------------
			parent, ok := d.(displayable.SingleContainer)
			if !ok {
				return
			} // Displayable is a parent

			contentPadding := parent.ContentPadding()
			contentGrid := grid.Padded(contentPadding)
			layouted := utils.LayoutedArea(contentGrid.Size(), parent.Content())
			contentGrid = contentGrid.SubGrid(layouted)
			calculateGridsH(parent.Content(), contentGrid)
	   ------------
	*/
}

func drawView() {
	drawDisplayable(main)
}

func drawDisplayable(d displayable.Displayable) {
	log.Log("Drawing displayable", d)
	d.Draw()
	parent, ok := d.(displayable.SingleContainer)
	if ok {
		log.Log("Displayable", d, "is parent, drawing child")
		drawDisplayable(parent.Content())
	}
}

func getContentGrid(workingArea areas.Area, d displayable.Displayable) pixel.SubGrid {
	var width, height, x, y int
	switch d.Layout() {
	case areas.FitToParent:
		width = workingArea.Width()
		height = workingArea.Height()
		x = workingArea.X1()
		y = workingArea.Y1()
	case areas.Centered:
		width = d.Size().Width()
		height = d.Size().Height()
		x = workingArea.X1() + (workingArea.Width()-width)/2
		y = workingArea.Y1() + (workingArea.Height()-height)/2
	default:
		panic(fmt.Sprintf("Invalid layout for displayable ", d))
	}

	area := areas.New(x, x+width, y, y+height)
	return termGrid.SubGrid(area)
}

// IsInitialized returns the status of initialization
// of hanoi.
func IsInitialized() bool {
	return initialized
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
	tb.Clear(tb.ColorWhite, tb.ColorWhite)

	initialized = true
	StoppedChannel = make(chan bool)
	stopChannel = make(chan bool)

	tb.Flush()
	tb.HideCursor()

	setSize(terminalSize())

	acceptInput(focused)

	calculateGrids()
	drawView()
	printGrid(termGrid.TotalSubGrid(), 0, 0)

	go pollEvents()

	go func() {
		defer func() {
			initialized = false
			log.Close()
			tb.Close()
			close(StoppedChannel)
		}()
		for {
			tb.Flush()
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

func setSize(s areas.Size) {
	termGrid = pixel.NewGrid(s.Width(), s.Height())
}

func redraw() {
	tb.Flush()
	calculateGrids()
	drawView()
	printGrid(termGrid.TotalSubGrid(), 0, 0)
}

func pollEvents() {
	for {
		select {
		default:
			switch ev := tb.PollEvent(); ev.Type {
			case tb.EventKey:
				log.Log("Caught key event")
				eventChannel <- KeyEvent{ev, false}
			case tb.EventResize:
				log.Log("Caught resize event")
				setSize(areas.NewSize(ev.Width, ev.Height))
				redraw()
			default:
				log.Log("Uncovered event:", ev)
				time.Sleep(10 * time.Millisecond)
			}
		case <-stopChannel:
			return
		}
	}
}

func acceptInput(f displayable.Displayable) {
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
	fmt.Println("closing this time")
	close(stopChannel)
}

func terminalSize() areas.Size {
	x, y := tb.Size()
	return areas.NewSize(x, y)
}
