package view

import "errors"
import "fmt"
import "time"
import tb "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/utils/log"
import "github.com/teolandon/hanoi/pixel"
import "github.com/teolandon/hanoi/areas"

var (
	initialized = false
	focused     Displayable
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

func printGrid(grid pixel.SubGrid, x, y int) {
	for i := 0; i < grid.Height(); i++ {
		line := grid.GetLine(i)
		for j, pixel := range line {
			fg := pixel.Palette.GetFGTermAttr(pixel.Highlight)
			bg := pixel.Palette.GetBGTermAttr(pixel.Highlight)
			tb.SetCell(x+j, y+i, pixel.Ch, fg, bg)
		}
	}
	tb.Sync()
}

func calculateGrids() {
	calculateGridsH(main, termGrid.TotalSubGrid())
}

func calculateGridsH(d Displayable, grid pixel.SubGrid) {
	if d == nil {
		return
	}

	d.setGrid(grid)
	log.Log("Set grid for", d, "to", grid)
	parent, ok := d.(Container)
	if !ok {
		return
	} // Displayable is a parent

	log.Log("Orig Grid:", grid)
	contentPadding := parent.ContentPadding()
	log.Log("ContentPadding:", contentPadding)
	contentGrid := grid.Padded(contentPadding)
	log.Log("ContentGrid:", contentGrid)
	layouted := layoutedArea(contentGrid.Size(), parent.Content())
	log.Log("layouted:", layouted)
	contentGrid = contentGrid.SubGrid(layouted)
	log.Log("ContentGrid:", contentGrid)
	calculateGridsH(parent.Content(), contentGrid)
}

func drawView() {
	drawDisplayable(main)
}

func layoutedArea(size areas.Size, d Displayable) areas.Area {
	var width, height, x, y int
	layout := d.Layout()
	contentSize := d.Size()
	log.Log("ContentSize for", d, ":", contentSize)
	switch layout {
	case FitToParent:
		log.Log("Displayable", d, "is Fit To parent")
		width = size.Width()
		height = size.Height()
		x = 0
		y = 0
	case Centered:
		log.Log("Displayable", d, "is Centered")
		width = contentSize.Width()
		height = contentSize.Height()
		x = (size.Width() - width) / 2
		y = (size.Height() - height) / 2
	default:
		panic("nooo")
	}
	ret := areas.NewFromSize(x, y, areas.NewSize(width, height))
	log.Log("Layout of size", size, "is", layout, "and results in", ret)
	return ret
}

func drawDisplayable(d Displayable) {
	log.Log("Drawing displayable", d)
	d.Draw()
	parent, ok := d.(Container)
	if ok {
		log.Log("Displayable", d, "is parent, drawing child")
		drawDisplayable(parent.Content())
	}
}

func getContentGrid(workingArea areas.Area, d Displayable) pixel.SubGrid {
	var width, height, x, y int
	switch d.Layout() {
	case FitToParent:
		width = workingArea.Width()
		height = workingArea.Height()
		x = workingArea.X1()
		y = workingArea.Y1()
	case Centered:
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

func terminalSize() areas.Size {
	x, y := tb.Size()
	return areas.NewSize(x, y)
}
