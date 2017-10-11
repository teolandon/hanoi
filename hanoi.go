package main

import "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/view"

func testMe() (ret []termbox.Cell) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	ret = termbox.CellBuffer()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.Clear(termbox.ColorWhite, termbox.ColorRed)

	view.Init()

	return
}

func main() {
	testMe()
}
