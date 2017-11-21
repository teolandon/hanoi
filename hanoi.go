package main

import "fmt"
import "github.com/teolandon/hanoi/view"

func testMe() {
	container := view.ButtonTitledContainer()
	view.SetRoot(container)

	fmt.Println(container.Parent())

	view.SetFocused(container.Content())
	fmt.Printf("container: %p\n", container)
	fmt.Printf("container.Content().Parent(): %p\n", (container.Content().Parent()))
	fmt.Println("end\n\n")

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
