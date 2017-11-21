package main

import "github.com/teolandon/hanoi/view"

func testMe() {
	container := view.ButtonTitledContainer()
	view.SetRoot(container)
	view.SetFocused(container.Content())

	err := view.Init()
	if err != nil {
		return
	}

	<-view.StoppedChannel

	return
}

func main() {
	testMe()
}
