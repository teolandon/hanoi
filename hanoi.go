package main

import "fmt"
import "github.com/teolandon/hanoi/view"

func testMe() {
	container := view.SimpleTitledContainer()
	view.SetRoot(container)
	view.SetFocused(container)

	// termbox.SetInputMode(termbox.InputEsc)
	// termbox.SetOutputMode(termbox.OutputNormal)
	// termbox.Clear(termbox.ColorWhite, termbox.ColorRed)

	err := view.Init()
	if err != nil {
		fmt.Println("LOLOLOL")
		return
	}

	view.SetFocused(container.Content().(view.TextBox))

	fmt.Println("before")
	<-view.StoppedChannel

	fmt.Println("after")
	return
}

func main() {
	testMe()
}
