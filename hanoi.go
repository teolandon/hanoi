package main

import "github.com/nsf/termbox-go"
import "github.com/teolandon/hanoi/menu"
import "fmt"

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

	menu.Init()

	return
}

func main() {
	c := testMe()
	fmt.Println(c[2])
}
