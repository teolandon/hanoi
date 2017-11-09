package main

import "fmt"
import "github.com/teolandon/hanoi/view"

func testMe() {
	container := view.ButtonTitledContainer()
	view.SetRoot(container)
	view.SetFocused(container)

	err := view.Init()
	if err != nil {
		fmt.Println("LOLOLOL")
		return
	}

	view.SetFocused(container.Content())

	fmt.Println("before")
	<-view.StoppedChannel

	fmt.Println("after")
	return
}

func main() {
	testMe()
}
