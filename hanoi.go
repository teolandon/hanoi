package main

import "github.com/teolandon/hanoi/view"
import "github.com/teolandon/hanoi/view/containers"
import "github.com/teolandon/hanoi/utils/log"

func testMe() {
	container := containers.NewWithButton()
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
	log.Init()
	testMe()
}
