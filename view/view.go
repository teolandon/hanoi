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
	stopChannel chan bool
	focusChange = make(chan bool)
	// StoppedChannel is the channel that signals that
	// the view has been stopped and exited out of.
	StoppedChannel chan bool
	eventChannel   = make(chan KeyEvent)
	main           mainContainer
	termGrid       pixel.PixelGrid
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

func printGrid(grid pixel.PixelGrid, x, y int) {
	for i := 0; i < grid.Height(); i++ {
		line := grid.GetLine(i)
		for j, pixel := range line {
			fg := pixel.Palette.GetFGTermAttr(pixel.Highlight)
			bg := pixel.Palette.GetBGTermAttr(pixel.Highlight)
			tb.SetCell(x+j, y+i, pixel.Ch, fg, bg)
		}
	}
}

func calculateGrids() {
	calculateGridsH(main, 0, 0)
}

func calculateGridsH(d Displayable, x, y int) {
	// STUB
	// grid := utils.NewPixelGrid(d.Size.x, d.Size.y)
	// d.SetGrid()
}

func getContentGrid(workingArea areas.Area, d Displayable) pixel.PixelGrid {
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
	return termGrid.SubGridFromArea(area)
}

func drawDisplayable(parentArea areas.Area, d Displayable) {
	// workArea := getWorkArea(parentArea, d.Size(), d.Layout())

	// log.Log("Drawing:", d)
	// log.Log("Area:", workArea)

	// grid := d.FillPixelGrid(workArea)
	// printGrid(grid, workArea.x1, workArea.y1)
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

	termSize := terminalSize()
	termGrid = pixel.NewGrid(termSize.Width(), termSize.Height())

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

func terminalSize() areas.Size {
	x, y := tb.Size()
	return areas.NewSize(x, y)
}

func terminalArea() areas.Area {
	x, y := tb.Size()
	ret := areas.NewFromSize(0, 0, areas.NewSize(x, y))
	return ret
}
